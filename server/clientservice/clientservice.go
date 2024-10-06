package clientservice

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/akilesh1706/crud-clientapps/proto"
)

type ClientService struct {
	pb.UnimplementedClientServiceServer
	collection *mongo.Collection
}

func NewClientService(collection *mongo.Collection) *ClientService {
	return &ClientService{collection: collection}
}

func (s *ClientService) CreateClient (ctx context.Context, req *pb.CreateClientRequest) (*pb.Client, error) {
	client := &pb.Client{
		ClientId:    req.ClientId,
		ClientName:  req.ClientName,
		ClientLogo:  req.ClientLogo,
		ClientSecret: req.ClientSecret,
		Domains:     req.Domains,
		Permissions: &pb.Permissions{
			RollNo:       req.Permissions.RollNo,
			Name:         req.Permissions.Name,
			Gender:       req.Permissions.Gender,
			DepartmentId: req.Permissions.DepartmentId,
			Nationality:  req.Permissions.Nationality,
			MobileNo:     req.Permissions.MobileNo,
			Address:      req.Permissions.Address,
			PictureUrl:   req.Permissions.PictureUrl,
			Dob:          req.Permissions.Dob,
		},
		StudentList: req.StudentList,
	}
	res, err := s.collection.InsertOne(ctx, client)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create client: %v", err)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to convert InsertedID to ObjectID")
	}
	
	client.Id = oid.Hex()
	return client, nil
}

func (s *ClientService) GetClient(ctx context.Context, req *pb.GetClientRequest) (*pb.Client, error) {

    oid := req.ClientId

    filter := bson.M{"clientid": oid}

    var client pb.Client
    err := s.collection.FindOne(ctx, filter).Decode(&client)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, status.Errorf(codes.NotFound, "Client not found")
        }
        return nil, status.Errorf(codes.Internal, "Failed to get client: %v", err)
    }

    client.Id = oid
    return &client, nil
}

func (s *ClientService) UpdateClient(ctx context.Context, req *pb.UpdateClientRequest) (*pb.Client, error) {
	// Convert ClientId to MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(req.ClientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid ID format: %v", err)
	}

	// Define the updated client document
	updatedDocument := bson.M{
		"_id":          oid,
		"clientId":     req.ClientId,
		"clientName":   req.ClientName,
		"clientLogo":   req.ClientLogo,
		"clientSecret": req.ClientSecret,
		"domains":      req.Domains,
		"permissions":  req.Permissions,
		"student_list": req.StudentList,
	}

	// Perform the replace operation
	_, err = s.collection.ReplaceOne(
		ctx,
		bson.M{"_id": oid},  // Use the correct MongoDB filter with the ObjectID
		updatedDocument,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Client not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to update client: %v", err)
	}

	// Convert the updated document back to the response type
	var updatedClient pb.Client
	err = mapstructure.Decode(updatedDocument, &updatedClient)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to parse updated client: %v", err)
	}

	// Set the ID of the updated client
	updatedClient.Id = req.ClientId
	return &updatedClient, nil
}


func (s *ClientService) DeleteClient(ctx context.Context, req *pb.DeleteClientRequest) (*pb.DeleteClientResponse, error) {
	// Convert ClientId to MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(req.ClientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid ID format: %v", err)
	}

	// Perform the deletion
	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete client: %v", err)
	}

	// Check if a document was deleted
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Client not found")
	}

	return &pb.DeleteClientResponse{Success: true}, nil
}


func (s *ClientService) ListClients (ctx context.Context, req *pb.ListClientsRequest) (*pb.ListClientsResponse, error) {
	var clients []*pb.Client
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))
	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list clients: %v", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var client pb.Client
		if err := cursor.Decode(&client); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to decode client: %v", err)
		}
		//client.Id = client.Id
		clients = append(clients, &client)
	}
	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to count clients: %v", err)
	}
	totalCount, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to count clients: %v", err)

	}
	return &pb.ListClientsResponse{
		Clients: clients,
		TotalCount: int32(totalCount),
	}, nil
}