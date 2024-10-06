package main

import (
	"context"
	"log"
	"net"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
 	"go.mongodb.org/mongo-driver/mongo/options"
 	"google.golang.org/grpc"
	
 	pb "github.com/akilesh1706/crud-clientapps/proto"
 	"github.com/akilesh1706/crud-clientapps/server/clientservice"
)

const (
	port = ":50051"
	//mongoURI = "mongodb://localhost:27017"
	mongoURI = "mongodb+srv://SpiderAdmin:SpiderRootAdmin2024@spiderproductions.aivs0.mongodb.net"
	//dbName = "clientdb"
	dbName = "lynxDev"
)

func main() {
	//Connection with MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel();

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	} ()

	//Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	//Create a new ClientService instance
	//clientCollection := client.Database(dbName).Collection("clients")
	clientCollection := client.Database(dbName).Collection("clientapps")
	clientService := clientservice.NewClientService(clientCollection)

	//Create a new gRPC server
	s := grpc.NewServer()

	//Register our service with the gRPC server
	pb.RegisterClientServiceServer(s, clientService)

	//Start listening on the specified port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Server listening on %s", port)

	//Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}