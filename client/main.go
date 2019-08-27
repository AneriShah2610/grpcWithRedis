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

var grpcClient helloWorld.GreeterClient

func main() {
	// create new redis client connection
	redisClient := newRedisClient()

	// grpc connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error to connect grpc server")
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	grpcClient = helloWorld.NewGreeterClient(conn)
	// set connection instance
	redisGRPCInstance := NewRedisGRPC(redisClient, grpcClient, ctx)
	err = redisGRPCInstance.getDataFromRedis()
	if err != nil {
		name := defaultName
		redisGRPCInstance.getDataFromGRPC(name)
	}
}

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

type RedisGRPCInterface interface {
	getDataFromRedis() error
	getDataFromGRPC(name string)
}

type RedisGRPCStruct struct {
	RedisClient *redis.Client
	GrpcClient  helloWorld.GreeterClient
	Ctx         context.Context
}

func NewRedisGRPC(redisClient *redis.Client, grpcClient helloWorld.GreeterClient, ctx context.Context) RedisGRPCInterface {
	return &RedisGRPCStruct{
		RedisClient: redisClient,
		GrpcClient:  grpcClient,
		Ctx:         ctx,
	}
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

func (rgi *RedisGRPCStruct) getDataFromRedis() error {
	return hGetValue(rgi.RedisClient)
}

func (rgi *RedisGRPCStruct) getDataFromGRPC(name string) {
	r, err := rgi.GrpcClient.SayHello(rgi.Ctx, &helloWorld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Error to get value from grpc %v", err)
		return
	}
	log.Printf("Greeting from client: %v", r.Name)
}
