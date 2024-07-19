package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/synexio/pos_node/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define flags
	register := flag.Bool("register", false, "Register a new client")
	subscribe := flag.Bool("subscribe", false, "Subscribe using UUID")
	getLastBlock := flag.Bool("getLastBlock", false, "Get the last block")
	addTransaction := flag.Bool("addTransaction", false, "Add a transaction")
	addTransactions := flag.Bool("addTransactions", false, "Add multiple transactions")
	bakeBlock := flag.Bool("bakeBlock", false, "Bake a block")
	confirmBake := flag.Bool("confirmBake", false, "Confirm bake")
	uuid := flag.String("uuid", "", "UUID to use for subscribe, bakeBlock and confirmBake")
	number := flag.Int("number", 1, "Number of transactions")
	sender := flag.String("sender", "Alexandre", "Transaction sender")
	receiver := flag.String("receiver", "Paul", "Transaction receiver")
	amount := flag.Int("amount", 1, "Transaction amount")
	data := flag.String("data", "Pour le McDo", "Transaction data")
	help := flag.Bool("help", false, "Display help")
	cheat := flag.Bool("cheat", false, "Can't say")

	// Parse the flags
	flag.Parse()

	// Connect to the server
	conn, err := grpc.NewClient("35.241.224.46:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBlockchainClient(conn)
	log.Printf("Connected to server '35.241.224.46:50051'")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Flag management
	switch {
	case *help:
		log.Printf("You must add a flag to run a precise function !\n" +
			"-register : Function to register a new client \n" +
			"-subscribe : Function to subscribe, must use -uuid flag that you got from registering \n" +
			"-addTransaction : Function to add a transaction, must use -sender and -receiver flags as uuids \n" +
			"-addTransactions : Function to add multiple transactions, must use -number as the number of transaction and -sender and -receiver flags as uuids \n" +
			"-bakeBlock : Function to bake a block, must use -uuid flag \n" +
			"-confirmBake : Function to confirm bake request, must use -uuid flag \n" +
			"-uuid : String parameter for functions \n" +
			"-sender : String parameter for transaction functions \n" +
			"-receiver : String parameter for transaction functions \n" +
			"-number : Int parameter for -addTransactions function")
	case *register:
		res, err := client.Register(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalf("Could not register: %v", err)
		}
		log.Printf("Registered with UUID: %s, Reputation: %d", res.GetUuid(), res.GetReputation())

	case *subscribe:
		if *uuid == "" {
			log.Fatalf("UUID is required for subscribe")
		}
		res, err := client.Subscribe(ctx, &pb.SubscribeRequest{Uuid: *uuid})
		if err != nil {
			log.Fatalf("Could not subscribe: %v", err)
		}
		log.Printf("Subscribed! Message: %s", res.Message)

	case *getLastBlock:
		res, err := client.GetLastBlock(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalf("Could not get last block: %v", err)
		}
		log.Printf("Last Block! Hash: %s, Previous Hash: %s, Block Number: %d, Data: %s", res.BlockHash, res.PreviousBlockHash, res.BlockNumber, res.Data)

	case *cheat:
		// Crée 10 comptes et les 9 derniers donnent tous leur fond au premier
		res, err := client.Register(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalf("Could not register: %v", err)
		}
		uuid := res.GetUuid()
		log.Printf("Registered with UUID: %s, Reputation: %d", res.GetUuid(), res.GetReputation())

		res2, err := client.Subscribe(ctx, &pb.SubscribeRequest{Uuid: uuid})
		if err != nil {
			log.Fatalf("Could not subscribe: %v", err)
		}
		log.Printf("Subscribed! Message: %s", res2.Message)

		for i := 0; i < 9; i++ {
			var uuid2 string
			res, err := client.Register(ctx, &pb.Empty{})
			if err != nil {
				log.Fatalf("Could not register: %v", err)
			}
			uuid2 = res.GetUuid()

			_, err = client.Subscribe(ctx, &pb.SubscribeRequest{Uuid: uuid2})
			if err != nil {
				log.Fatalf("Could not subscribe: %v", err)
			}

			_, err = client.AddTransaction(ctx, &pb.Transaction{
				Sender:   uuid2,
				Receiver: uuid,
				Amount:   100,
				Data:     *data,
			})
			if err != nil {
				log.Fatalf("Could not add transaction: %v", err)
			}
			log.Printf("Transaction Added!")

		}

	case *addTransaction:
		if *sender == "" || *receiver == "" {
			log.Fatalf("Sender and receiver are required for a transaction ! Add the flags")
		}
		_, err := client.AddTransaction(ctx, &pb.Transaction{
			Sender:   *sender,
			Receiver: *receiver,
			Amount:   int32(*amount),
			Data:     *data,
		})
		if err != nil {
			log.Fatalf("Could not add transaction: %v", err)
		}
		log.Printf("Transaction Added!")

	case *addTransactions:
		if *number <= 0 {
			log.Fatalf("Can't make 0 transactions, at least 1 !")
		}
		if *sender == "" || *receiver == "" {
			log.Fatalf("Sender and receiver are required for a transaction ! Add the flags")
		}
		for i := 1; i < *number+1; i++ {
			/* Trying to abuse balance bug
			var s, r string
			if i%2 == 0 {
				s = *sender
				r = *receiver
			} else {
				r = *sender
				s = *receiver
			}*/
			_, err := client.AddTransaction(ctx, &pb.Transaction{
				Sender:   *sender,
				Receiver: *receiver,
				Amount:   int32(*amount),
				Data:     *data,
			})
			if err != nil {
				log.Fatalf("Could not add transaction n°%d: %v", i, err)
			}
			log.Printf("Transaction n°%d Added!", i)
		}

	case *bakeBlock:
		if *uuid == "" {
			log.Fatalf("UUID is required for baking a block")
		}
		res, err := client.BakeBlock(ctx, &pb.BakeRequest{Uuid: *uuid})
		if err != nil {
			log.Fatalf("Could not bake block: %v", err)
		}
		log.Printf("Block Baked! UUID: %s, Message: %s", res.Uuid, res.Message)

	case *confirmBake:
		if *uuid == "" {
			log.Fatalf("UUID is required for confirming bake request")
		}
		_, err := client.ConfirmBake(ctx, &pb.ConfirmRequest{Uuid: *uuid})
		if err != nil {
			log.Fatalf("Could not confirm request: %v", err)
		}
		log.Printf("Request confirmed!")

	default:
		log.Printf("You must add a flag to run a function ! Try -help")
	}
}
