package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"redis_test/push"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var c chan os.Signal
var c1 chan os.Signal
var wg sync.WaitGroup

func producer(client *redis.Client, sourceQ string, exitCh string) {
Loop:
	for {
		select {
		case s := <-c:
			fmt.Println()
			fmt.Println("Producer | get exit signal", s)
			client.Publish(exitCh, "Exit")
			break Loop
		default:
		}
		data := rand.Int31n(400)
		err := push.Push(client, sourceQ, strconv.Itoa(int(data)))
		if err != nil {
			fmt.Println("err:", err)
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Producer |  exit")
	wg.Done()
}

func collector(client *redis.Client, resultQ string) error {
	var sum int64 = 0
Loop:
	for {
		select {
		case s := <-c1:
			fmt.Println()
			fmt.Println("collector | get exit signal", s)
			break Loop
		default:
		}
		isExists, err1 := client.Exists(resultQ).Result()
		if err1 != nil {
			if err1 == redis.Nil {
				fmt.Printf("key %v 不存在\n", resultQ)
			} else {
				fmt.Println("collector error:", err1)
			}
			return err1
		}
		type_, err2 := client.Type(resultQ).Result()
		if err2 != nil {
			if err2 == redis.Nil {
				fmt.Printf("key %v 不存在\n", resultQ)
			} else {
				fmt.Println("collector error:", err2)
			}
			return err2
		}
		if isExists > 0 && type_ == "list" {
			data, err3 := client.RPop(resultQ).Result()
			if err3 != nil {
				if err3 == redis.Nil {
					fmt.Printf("key %v 不存在\n", resultQ)
					time.Sleep(1 * time.Second)
				} else {
					fmt.Println("collector error:", err3)
				}
				return err3
			}
			fmt.Println("collector received data: ", data)
			d, err := strconv.ParseInt(data, 10, 64)
			if err != nil {
				fmt.Println("collector err:", err)
			}
			sum += d
			fmt.Println("collector get sum ", sum)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("collector | exit")
	wg.Done()
	return nil
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

	c = make(chan os.Signal, 1)
	c1 = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	signal.Notify(c1, os.Interrupt, os.Kill)
	wg.Add(1)
	go collector(client, resultQ)
	wg.Add(1)
	go producer(client, sourceQ, exitCh)
	wg.Wait()
}
