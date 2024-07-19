package main

import (
	"context"
	"log"
	"time"

	pb "github.com/synexio/pos_node/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("35.241.224.46:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBlockchainClient(conn)
	log.Printf("Connected to server '35.241.224.46:50051'")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	/*res, err := client.Register(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}

	uuid := res.GetUuid()
	log.Printf("Registered with UUID: %s, Reputation: %d", res.GetUuid(), res.GetReputation())

	res2, err := client.Subscribe(ctx, &pb.SubscribeRequest{Uuid: uuid})
	if err != nil {
		log.Fatalf("Could not subscribe: %v", err)
	}
	log.Printf("Subscribed ! Message: %s", res2.Message)

	res3, err := client.GetLastBlock(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Could not get last block: %v", err)
	}
	log.Printf("Last Block ! Hash: %s, Previous Hash: %s, Block Number: %d, Data: %s", res3.BlockHash, res3.PreviousBlockHash, res3.BlockNumber, res3.Data)

	_, err = client.AddTransaction(ctx, &pb.Transaction{
		Sender:   "Alexandre",
		Receiver: "Paul",
		Amount:   1,
		Data:     "Pour le mcdo",
	})
	if err != nil {
		log.Fatalf("Could not add transaction: %v", err)
	}
	log.Printf("Transaction Added !")*/

	uuid := "5773415196899605267"

	res4, err := client.BakeBlock(ctx, &pb.BakeRequest{Uuid: uuid})
	if err != nil {
		log.Fatalf("Could not bake block: %v", err)
	}
	log.Printf("Block Baked ! UUID : %s, Message: %s", res4.Uuid, res4.Message)

	_, err = client.ConfirmBake(ctx, &pb.ConfirmRequest{Uuid: uuid})
	if err != nil {
		log.Fatalf("Could not confirm request: %v", err)
	}
	log.Printf("Request confirmed !")

}
