package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/kelseyhightower/envconfig"
)

type kafkaConfig struct {
	addr      string `envconfig:"KAFKA_SERVER" default:"localhost:9092"`
	topic     string `envconfig:"KAFKA_TOPIC" default:"Test.A"`
	partition int32  `envconfig:"KAFKA_PARTITION" default:"0"`
}

func consume(ctx context.Context, kafkaConf *kafkaConfig) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{kafkaConf.addr}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	partition, err := consumer.ConsumePartition(kafkaConf.topic, kafkaConf.partition, sarama.OffsetNewest)
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
	config := &kafkaConfig{}
	if err := envconfig.Process("", &config); err != nil {
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
