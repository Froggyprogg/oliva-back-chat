package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	gchat "oliva-back-chat/internal/gen"
	client2 "oliva-back-chat/pkg/handlers/chat"
)

func RunClient() {
	const serverID = "localhost:5000"

	log.Println("Connecting : " + serverID)
	conn, err := grpc.Dial(serverID)

	if err != nil {
		log.Fatalf("Failed to connect gRPC server :: %v", err)
	}
	defer conn.Close()

	client := gchat.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to get response from gRPC server :: %v", err)
	}
	ch := client2.ClientHandle{Stream: stream}
	ch.ClientConfig()
	go ch.SendMessage()
	go ch.ReceiveMessage()

	// block main
	bl := make(chan bool)
	<-bl
}
