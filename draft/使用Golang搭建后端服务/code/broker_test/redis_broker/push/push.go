package push

import (
	"fmt"

	"github.com/go-redis/redis"
)

func Push(client *redis.Client, Q string, value string) error {
	isExists, err1 := client.Exists(Q).Result()
	if err1 != nil {
		if err1 == redis.Nil {
			fmt.Printf("key %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err1)
		}
		return err1
	}

	type_, err2 := client.Type(Q).Result()
	if err2 != nil {
		if err2 == redis.Nil {
			fmt.Printf("queue %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err2)
		}
		return err2
	}
	if isExists > 0 && type_ == "list" {
		_, err := client.LPushX(Q, value).Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Printf("queue %v 不存在\n", Q)
			} else {
				fmt.Println("error:", err)
			}
			return err
		}
		fmt.Printf("send %v to %v\n", value, Q)
		return nil
	}
	_, err := client.Del(Q).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err)
		}
		return err
	}
	_, err = client.LPush(Q, value).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err)
		}
		return err
	}
	fmt.Printf("send %v to %v\n", value, Q)
	return nil
}
