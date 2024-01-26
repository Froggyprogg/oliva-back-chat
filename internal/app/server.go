package app

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"oliva-back-chat/config"
	chat "oliva-back-chat/pkg/handlers/chat"
)

func Run(cfg *config.Config) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Server is startup", "Server listen port:", cfg.GRPC.Port)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	chat.Register(grpcServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}

}
