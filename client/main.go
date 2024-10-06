package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "github.com/akilesh1706/crud-clientapps/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new client
	createResp, err := client.CreateClient(ctx, &pb.CreateClientRequest{
		ClientId: "112122060",
		ClientName:  "Client Akilesh",
		ClientLogo:  "https://example.com/logo12345.png",
		ClientSecret: "lolpog123",
		Domains:     []string{"amazon.com", "flipkart.com"},
		Permissions: &pb.Permissions{
			RollNo:       true,
			Name:         true,
			Gender:       false,
			DepartmentId: true,
			Nationality:  false,
			MobileNo:     true,
			Address:      true,
			PictureUrl:   false,
			Dob:          false,
		},
		StudentList: []*pb.Student{
			{
				RollNo: 150,
				Oid:    "unique-student-id-150",
			},
			{
				RollNo: 200,
				Oid:    "unique-student-id-200",
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}
	fmt.Printf("Created client: %v\n", createResp)

	// Get the created client
	getResp, err := client.GetClient(ctx, &pb.GetClientRequest{ClientId: createResp.ClientId})
	if err != nil {
		log.Fatalf("Could not get client: %v", err)
	}
	fmt.Printf("Got client: %v\n", getResp)

	// Update the client
	updateResp, err := client.UpdateClient(ctx, &pb.UpdateClientRequest{
		ClientId:    createResp.Id,  // Use the correct MongoDB document ID, which is stored in 'createResp.Id'
		ClientName:  "Updated Client Name",
		ClientLogo:  "https://example.com/updated-logo.png",
		ClientSecret: "updatedsecret",
		Domains:     []string{"updated.com", "newdomain.com"},
		Permissions: &pb.Permissions{
			RollNo:       true,
			Name:         true,
			Gender:       true,
			DepartmentId: false,
			Nationality:  true,
			MobileNo:     true,
			Address:      false,
			PictureUrl:   true,
			Dob:          true,
		},
		StudentList: []*pb.Student{
			{
				RollNo: 201,
				Oid:    "updated-student-id-1",
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not update client: %v", err)
	}
	fmt.Printf("Updated client: %v\n", updateResp)
	

	// List clients
	listResp, err := client.ListClients(ctx, &pb.ListClientsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("Could not list clients: %v", err)
	}
	fmt.Printf("Listed %d clients, total count: %d\n", len(listResp.Clients), listResp.TotalCount)
	for _, client := range listResp.Clients {
		fmt.Printf("- %s: %s\n", client.ClientId, client.ClientName)
	}

	// Delete the client
	deleteResp, err := client.DeleteClient(ctx, &pb.DeleteClientRequest{ClientId: createResp.Id})
if err != nil {
	log.Fatalf("Could not delete client: %v", err)
}
fmt.Printf("Deleted client, success: %v\n", deleteResp.Success)

}
