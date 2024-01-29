package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	gchat "oliva-back-chat/internal/gen"
	client2 "oliva-back-chat/pkg/handlers/client"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Faile to conncet to gRPC server :: %v", err)
	}
	defer conn.Close()

	//call ChatService to create a stream
	client := gchat.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	// implement communication with gRPC server
	ch := client2.ClientHandle{Stream: stream}
	ch.ClientConfig()
	go ch.SendMessage()
	go ch.ReceiveMessage()

	//blocker
	bl := make(chan bool)
	<-bl

}
