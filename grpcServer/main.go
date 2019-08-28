package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"test/grpcWithRedis/helloWorld"
)

type GRPCServer struct{}

var client helloWorld.GreeterClient

const (
	port = ":4000"
)

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error to listen %v", err)
	}
	fmt.Println("Listen on server", port)
	server := grpc.NewServer()
	helloWorld.RegisterGreeterServer(server, &GRPCServer{})
	if err = server.Serve(listen); err != nil {
		log.Fatalf("Failed to sreve %v", err)
	}
}

func (s *GRPCServer) SayHello(ctx context.Context, in *helloWorld.HelloRequest) (*helloWorld.HelloResponse, error) {
	log.Printf("Received: %v", in.UserName)
	return &helloWorld.HelloResponse{UserName: in.UserName}, nil
}
