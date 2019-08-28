package main

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-redis/redis"
)

/*const (
	port = ":4001"
)
*/

func main() {
	// create new connection
	client := newClient()
	// set value
	err := hSetValue(client)
	if err != nil {
		fmt.Println("Error at set value", err)
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

func hSetValue(client *redis.Client) error {
	user := []User{{ID: "123456789", UserName: "AneriShah", Email: "anerishah36@gmail.com", MobileNo: "0123456789"}, {ID: "2345678901", UserName: "AbcXyz", Email: "abc.xyz@gmail.com", MobileNo: "1234567890"}}
	for _, i := range user {
		userMap := structs.Map(i)
		err := client.HMSet("userSampleKey:"+i.UserName, userMap).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
