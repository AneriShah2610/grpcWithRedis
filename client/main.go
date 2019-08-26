package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"log"
	"test/grpcWithRedis/helloWorld"
	"time"
)

const (
	address     = "localhost:4000"
	defaultName = "Aneri Shah"
)

func main() {
	// create new redis client connection
	client := newClient()

	// get value
	err := hGetValue(client)
	if err != nil {
		//connect grpc server
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("GRPC connection error: %v", err)
		}
		defer conn.Close()

		c := helloWorld.NewGreeterClient(conn)

		name := defaultName
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.SayHello(ctx, &helloWorld.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("Error to get value from grpc %v", err)
			return
		}
		log.Printf("Greeting from client: %v", r.Name)
	}
}

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

type User struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	MobileNo string `json:"mobileNo"`
}

func hGetValue(client *redis.Client) error {
	m, err := client.HGetAll("userSampleKey:AneriShah").Result()
	if err != nil {
		return err
	}
	user := User{}
	for key, value := range m {
		switch key {
		case "ID":
			user.ID = value
		case "UserName":
			user.UserName = value
		case "Email":
			user.Email = value
		case "MobileNo":
			user.MobileNo = value
		}
	}
	fmt.Printf("%v \n", user)
	return nil
}
