package main

import (
	"log"
	"net"
	"service/mail-server/config"
	"service/mail-server/register"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.New()

	lis, err := net.Listen("tcp", ":"+cfg.PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	register.Register(server, cfg)

	log.Println("gRPC Server running on :" + cfg.PORT)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
