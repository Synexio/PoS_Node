package main

import (
	"context"
	"fmt"
	pb "github.com/synexio/pos_node/proto"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
)

type Block struct {
}

type server struct {
	pb.UnimplementedBlockchainServer
}

func (s *server) Register(ctx context.Context, in *pb.Empty) (*pb.RegisterResponse, error) {
	uuid := generateUUID()
	reputation := int32(100)
	log.Printf("Client %s registered with reputation %d", uuid, reputation)
	return &pb.RegisterResponse{Uuid: uuid, Reputation: reputation}, nil
}

func generateUUID() string {
	return fmt.Sprintf("%d", rand.Int63())
}

func main() {
	lis, err := net.Listen("tcp", ":55001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Server listening on port 55001")

	s := grpc.NewServer()

	pb.RegisterBlockchainServer(s, &server{})
	log.Println("Server started successfully")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
