package main

import (
	"context"
	"log"
	"time"

	message "github.com/bottomhalf/btc-grpc-shared/genproto/message"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Starting gRPC Verification Check...")

	// 1. Connect to the Message Service gRPC server
	address := "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := message.NewMessageProcessorClient(conn)

	// 2. Prepare a test event
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Sending test event: 'test_handoff' for user 'BOT00035'...")
	req := &message.WsEventRequest{
		Event:   "test_handoff",
		Payload: []byte("{\"message\": \"Hello from verification script\"}"),
		UserId:  "BOT00035",
	}

	// 3. Call the service
	resp, err := c.ProcessEvent(ctx, req)
	if err != nil {
		log.Fatalf("could not process event: %v (Is the Message Service running?)", err)
	}

	if resp.Success {
		log.Println("SUCCESS: Message Service acknowledged the gRPC event!")
	} else {
		log.Printf("FAILURE: Message Service returned error: %s", resp.Error)
	}
}
