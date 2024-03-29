package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/qiniu/x/log.v7"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := []string{"127.0.0.1:9092"}
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()
	_, e := master.Partitions("topic-python")
	if e != nil {
		log.Println(e)
	}
	consumer, err := master.ConsumePartition("topic-python", 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				fmt.Println("Received messages", string(msg.Key), string(msg.Value), msg.Topic)
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
}
