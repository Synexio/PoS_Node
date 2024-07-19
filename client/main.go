package main

import (
	"context"
	"log"
	"time"

	protobuf "github.com/synexio/pos_node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.NewClient("localhost:55001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	grpc_client := protobuf.NewBlockchainClient(connection)
	ctx, CtxCancel := context.WithTimeout(context.Background(), 10*time.Second)

	r, err := grpc_client.SayHello(ctx, &protobuf.InputRequest{Name: "Toto"})
	if err != nil {
		log.Fatalf("SayHello Method Error: %v", err)
	}

	log.Printf("Bien reçu : %s", r.GetMessage())

	defer connection.Close()
	defer CtxCancel()
}
