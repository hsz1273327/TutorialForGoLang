package main

import (
	"fmt"
	"redis_test/push"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func get_result(client *redis.Client, sourceQ string, resultQ string) error {
	fmt.Printf("get_result running\n")

	for {
		isExists, err1 := client.Exists(sourceQ).Result()
		if err1 != nil {
			if err1 == redis.Nil {
				fmt.Printf("key %v 不存在\n", sourceQ)
			} else {
				fmt.Println("error:", err1)
			}
			return err1
		}
		type_, err2 := client.Type(sourceQ).Result()
		if err2 != nil {
			if err2 == redis.Nil {
				fmt.Printf("key %v 不存在\n", sourceQ)
			} else {
				fmt.Println("error:", err2)
			}
			return err2
		}
		if isExists > 0 && type_ == "list" {
			data, err3 := client.RPop(sourceQ).Result()
			if err3 != nil {
				if err3 == redis.Nil {
					fmt.Printf("key %v 不存在\n", resultQ)
					time.Sleep(1 * time.Second)
				} else {
					fmt.Println("error:", err3)
				}
				return err3
			}
			fmt.Println("received data: ", data)
			d, err := strconv.ParseInt(data, 10, 32)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				result := d * d
				push.Push(client, resultQ, strconv.Itoa(int(result)))
			}
		}
	}
}

func main() {
	opt, err := redis.ParseURL("redis://localhost:6379/1")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	sourceQ := "sourceQ"
	resultQ := "resultQ"
	exitCh := "exitCh"
	go get_result(client, sourceQ, resultQ)
	pubsub := client.Subscribe(exitCh)
	defer pubsub.Close()
	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		if msg.Payload == "Exit" {
			return
		}
	}
}
