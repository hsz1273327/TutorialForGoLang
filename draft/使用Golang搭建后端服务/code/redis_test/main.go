package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func get_key(client *redis.Client, key string) {
	fmt.Println("get key")
	val, err := client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key %v 不存在\n", key)
		} else {
			fmt.Println("error:", err)
		}
	} else {
		fmt.Println("key", key, val)
	}
}

func do_get_key(client *redis.Client, key string) {
	fmt.Println("do get key")
	val, err := client.Do("get", key).String()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key %v 不存在\n", key)
		} else {
			fmt.Println("error:", err)
		}
	} else {
		fmt.Println("key", key, val)
	}
}

func simple_set(client *redis.Client) {
	key := "testkey"
	get_key(client, key)
	val, err := client.Set(key, 105, 0).Result()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("set:", val)
		do_get_key(client, key)
		val, err := client.Del(key).Result()
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("delete suc:", val)
		}
	}
}

func incr_pipeline(client *redis.Client) {
	pipe := client.Pipeline()

	incr := pipe.Incr("pipeline_counter")
	pipe.Expire("pipeline_counter", time.Hour)
	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)

}

func exit_check(client *redis.Client) {
	val, err := client.Exists("asd").Result()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("exists:", val)
	}
	client.Set("asd", 105, 0).Result()
	val, err = client.Exists("asd").Result()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("exists:", val)
	}
}

func rpop(client *redis.Client) {
	data, err := client.RPop("list1").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key 不存在\n")
		} else {
			fmt.Println("error:", err)
		}
	} else {
		fmt.Println("exists:", data)
	}

}

func main() {
	opt, err := redis.ParseURL("redis://localhost:6379/1")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	simple_set(client)
	incr_pipeline(client)
	exit_check(client)
	rpop(client)
}
