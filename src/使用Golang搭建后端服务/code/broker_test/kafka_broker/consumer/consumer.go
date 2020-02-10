package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var csend chan os.Signal
var cexit chan os.Signal
var cget chan os.Signal
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

func get_result(c *kafka.Consumer, p *kafka.Producer, sourceQ string, resultQ string) error {
	fmt.Printf("get_result running\n")
	err := c.SubscribeTopics([]string{sourceQ}, nil)
	if err != nil {
		fmt.Println("collector | err", err)
		return err
	}
Loop:
	for {
		select {
		case s := <-cget:
			fmt.Println()
			fmt.Println("get_result | get exit signal", s)
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
					fmt.Println("get_result err:", err)
				}
				result := d * d
				p.ProduceChannel() <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &resultQ,
						Partition: kafka.PartitionAny,
					},
					Value: []byte(strconv.Itoa(int(result)))}
				fmt.Println("get_result send msg ", result)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}
	fmt.Println("get_result |  exit")
	wg.Done()
	return nil
}
func listen_exit(c *kafka.Consumer, exitCh string) error {
	fmt.Printf("listen_exit running\n")
	err := c.SubscribeTopics([]string{exitCh}, nil)
	if err != nil {
		fmt.Println("listen_exit | err", err)
		return err
	}
Loop:
	for {
		select {
		case s := <-cexit:
			fmt.Println()
			fmt.Println("listen_exit | get exit signal", s)
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
				fmt.Println("listen_exit | get exit signal from remote----------------")
				value := string(e.Value)
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition)
				if value == "Exit" {
					fmt.Println("get exit")
					cget <- os.Interrupt
					csend <- os.Interrupt
					break Loop
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}
	fmt.Println("listen_exit |  exit")
	wg.Done()
	return nil
}
func main() {
	broker := "localhost:9092"
	sourceQ := "sourceQ"
	resultQ := "resultQ"
	exitCh := "exitCh"
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
	kafkaExitConsumer, err2 := kafka.NewConsumer(&consumerConf)
	if err2 != nil {
		panic(err2)
	}
	csend = make(chan os.Signal, 1)
	cget = make(chan os.Signal, 1)
	cexit = make(chan os.Signal, 1)
	signal.Notify(csend, os.Interrupt, os.Kill)
	signal.Notify(cget, os.Interrupt, os.Kill)
	signal.Notify(cexit, os.Interrupt, os.Kill)
	wg.Add(1)
	go get_result(kafkaConsumer, kafkaProducer, sourceQ, resultQ)
	wg.Add(1)
	go secondRef(kafkaProducer)
	wg.Add(1)
	go listen_exit(kafkaExitConsumer, exitCh)
	wg.Wait()
}
