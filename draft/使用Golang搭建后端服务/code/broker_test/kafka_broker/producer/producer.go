package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var csend chan os.Signal
var cprod chan os.Signal
var ccoll chan os.Signal
var wg sync.WaitGroup

func secondRef(p *kafka.Producer) {
	fmt.Println("secondRef | start")
Loop:
	for {
		select {
		case s := <-csend:
			fmt.Println()
			fmt.Println("secondRef | get exit signal", s)
			break Loop
		case e := <-p.Events():
			fmt.Println("secondRef | get event")
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		default:
		}
	}
	fmt.Println("secondRef |  exit")
	wg.Done()
}

func producer(p *kafka.Producer, sourceQ string, exitCh string) {
Loop:
	for {
		select {
		case s := <-cprod:
			fmt.Println()
			fmt.Println("Producer | get exit signal", s)
			p.ProduceChannel() <- &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &exitCh,
					Partition: kafka.PartitionAny,
				},
				Value: []byte("Exit")}
			time.Sleep(1 * time.Second)
			fmt.Println("Producer send msg exit to exitch")
			break Loop
		default:
		}
		data := rand.Int31n(400)
		p.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &sourceQ,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(strconv.Itoa(int(data)))}
		fmt.Println("Producer send msg ", data, " to sourceQ")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Producer |  exit")
	wg.Done()
}

func collector(c *kafka.Consumer, resultQ string) error {
	err := c.SubscribeTopics([]string{resultQ}, nil)
	if err != nil {
		fmt.Println("collector | err", err)
		return err
	}
	var sum int64 = 0
Loop:
	for {
		select {
		case s := <-ccoll:
			fmt.Println()
			fmt.Println("collector | get exit signal", s)
			break Loop

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				value := string(e.Value)
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition)
				d, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					fmt.Println("collector err:", err)
				}
				sum += d
				fmt.Println("collector get sum ", sum)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}
	fmt.Println("collector | exit")
	wg.Done()
	return nil
}

func main() {
	sourceQ := "sourceQ"
	resultQ := "resultQ"
	exitCh := "exitCh"
	broker := "localhost:9092"
	producerConf := kafka.ConfigMap{
		"bootstrap.servers": broker,
	}
	consumerConf := kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        resultQ,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
		"auto.offset.reset":               "latest"}
	kafkaProducer, err1 := kafka.NewProducer(&producerConf)
	if err1 != nil {
		panic(err1)
	}
	kafkaConsumer, err2 := kafka.NewConsumer(&consumerConf)
	if err2 != nil {
		panic(err2)
	}
	csend = make(chan os.Signal, 1)
	cprod = make(chan os.Signal, 1)
	ccoll = make(chan os.Signal, 1)
	signal.Notify(ccoll, os.Interrupt, os.Kill)
	signal.Notify(cprod, os.Interrupt, os.Kill)
	signal.Notify(csend, os.Interrupt, os.Kill)
	wg.Add(1)
	go collector(kafkaConsumer, resultQ)
	wg.Add(1)
	go producer(kafkaProducer, sourceQ, exitCh)
	wg.Add(1)
	go secondRef(kafkaProducer)
	wg.Wait()
}
