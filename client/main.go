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
	address         = "localhost:4000"
	redisGrpcCtxKey = "redisGrpcCtxKey"
)

var grpcClient helloWorld.GreeterClient

func main() {
	start := time.Now()
	// create new redis client connection
	redisClient := newRedisClient()

	// grpc connection
	grpcConnection, err := newGrpcConnection()
	if err != nil {
		fmt.Printf("Error to connect grpc connection: %v", err)
	}
	defer grpcConnection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	grpcClient = helloWorld.NewGreeterClient(grpcConnection)

	// set both redis and grpc client instances
	redisGRPCInstance := NewRedisGRPC(redisClient, grpcClient, ctx)
	var keys []interface{}
	keys = append(keys, "userSampleKey:AneriShah", "userSampleKey:AbcXyz")

	userData, err := redisGRPCInstance.getData(keys)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Println("userData:", userData)
	log.Printf("Code execution time %s", time.Since(start))

	//commonMethod([]string{"userData:", "userData:"})
	//commonMethod([]string{"userData:"})
}

/*func commonMethod(userId interface{}) {

	typeOf := reflect.TypeOf(userId)
	if typeOf.Kind() == reflect.String {
	//	userId is string
		redisClient := newRedisClient()
		hGet := redisClient.HGetAll(userId)
		if hGet == nil {
			grpcConnection, e := newGrpcConnection()
			if e != nil {
				//	throw error
			}
			// Call method to fetch particular user info
			helloResponse, err := grpcClient.SayHello(nil, nil)
			if err != nil {
				// convert grpc Response to generalUser response & return
				return
			}
			// throw internal server err
		}
		//	convert redis Response to generalUser response & return
		return

	} else if typeOf.Kind() == reflect.Array {
	//	fetch for multiple keys

	}
	else {
	//	 type must be array
	}
}*/

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6397",
		Password: "",
		DB:       0,
	})
	return client
}

func newGrpcConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type RedisGRPCInterface interface {
	getData(keys []interface{}) ([]UserData, error)
	ping() error
	getDataFromRedis(keys []interface{}) ([]UserData, error)
	getDataFromGRPC(keys []interface{}) ([]UserData, error)
}

type UserData struct {
	values interface{}
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

func (rgi *RedisGRPCStruct) getData(keys []interface{}) ([]UserData, error) {
	err := rgi.ping()
	if err != nil {
		return rgi.getDataFromGRPC(keys)
	}
	return rgi.getDataFromRedis(keys)
}

func (rgi *RedisGRPCStruct) ping() error {
	return rgi.RedisClient.Ping().Err()
}

type User struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	MobileNo string `json:"mobileNo"`
}

func (rgi *RedisGRPCStruct) getDataFromRedis(keys []interface{}) ([]UserData, error) {
	var userDatas []UserData
	for _, i := range keys {
		userData := User{}
		var data UserData
		m, err := rgi.RedisClient.HGetAll(i.(string)).Result()
		if err != nil {
			return userDatas, err
		}
		for key, value := range m {
			switch key {
			case "ID":
				userData.ID = value
			case "UserName":
				userData.UserName = value
			case "Email":
				userData.Email = value
			case "MobileNo":
				userData.MobileNo = value
			}
		}
		data.values = userData
		userDatas = append(userDatas, data)
	}
	return userDatas, nil
}

func (rgi *RedisGRPCStruct) getDataFromGRPC(userNames []interface{}) ([]UserData, error) {
	var datas []UserData
	for _, i := range userNames {
		var data UserData
		r, err := rgi.GrpcClient.SayHello(rgi.Ctx, &helloWorld.HelloRequest{UserName: i.(string)})
		if err != nil {
			return datas, err
		}
		data.values = r.UserName
		datas = append(datas, data)
	}
	return datas, nil
}
