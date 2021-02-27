package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/kelseyhightower/envconfig"
)

type KafkaConfig struct {
	Addrs     string `envconfig:"KAFKA_SERVER" default:"localhost:9092"`
	Topic     string `envconfig:"KAFKA_TOPIC" default:"Test.A"`
	Partition int32  `envconfig:"KAFKA_PARTITION" default:"0"`
}

func consume(ctx context.Context, kafkaConf *KafkaConfig) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	kafkaAddrs := strings.Split(kafkaConf.Addrs, ",")
	log.Printf("kafka_server_addresses: %v", kafkaAddrs)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(kafkaAddrs, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	partition, err := consumer.ConsumePartition(kafkaConf.Topic, kafkaConf.Partition, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-childCtx.Done():
				return
			case kmsg := <-partition.Messages():
				var content map[string]interface{}
				if err := json.Unmarshal(kmsg.Value, &content); err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%v\n", content)
			}
		}
	}()

	<-ctx.Done()
}

func main() {
	config := &KafkaConfig{}
	if err := envconfig.Process("", config); err != nil {
		fmt.Println("cannot load environment")
		os.Exit(-1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM)

	go consume(ctx, config)

	<-signalCh
}
