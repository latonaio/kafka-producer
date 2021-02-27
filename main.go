package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"bitbucket.org/latonaio/aion-core/pkg/log"
	"github.com/Shopify/sarama"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

var msName string = "kafka-producer"

type KafkaConfig struct {
	Addrs string `envconfig:"KAFKA_SERVER" default:"localhost:9092"`
}

type kafkaMsg struct {
	topic   string
	key     string
	content map[string]interface{}
}

func kanbanToKafkaMsg(kanban *msclient.WrapKanban) (*kafkaMsg, error) {
	metadata, err := kanban.GetMetadataByMap()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get metadata.")
	}
	topic, ok := metadata["topic"].(string)
	if !ok {
		return nil, errors.New("not found topic")
	}
	key, ok := metadata["key"].(string)
	if !ok {
		return nil, errors.New("not found key")
	}
	content, ok := metadata["content"].(map[string]interface{})
	if !ok {
		return nil, errors.New("not found content")
	}
	return &kafkaMsg{
		topic:   topic,
		key:     key,
		content: content,
	}, nil
}

func produce(ctx context.Context, dataCh <-chan *msclient.WrapKanban, producer sarama.AsyncProducer) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-producer.Successes():
			log.Printf("success send")
		case <-producer.Errors():
			log.Printf("failed send")
		case data := <-dataCh:
			log.Printf("received data: %v\n", data)
			kmsg, err := kanbanToKafkaMsg(data)
			if err != nil {
				log.Printf("%v\n", err)
			}
			c, err := json.Marshal(kmsg.content)
			if err != nil {
				log.Printf("cannot convert json")
			}
			msg := &sarama.ProducerMessage{
				Topic: kmsg.topic,
				Key:   sarama.StringEncoder(kmsg.key),
				Value: sarama.StringEncoder(c),
			}
			producer.Input() <- msg
		}
	}
}

func main() {
	kConfig := &KafkaConfig{}
	if err := envconfig.Process("", kConfig); err != nil {
		log.Fatal("cannot load environment")
	}
	kafkaAddrs := strings.Split(kConfig.Addrs, ",")
	log.Printf("kafka_server_addresses: %v", kafkaAddrs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := msclient.NewKanbanClient(ctx)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	kanbanCh, err := c.GetKanbanCh(msName, c.GetProcessNumber())
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_2
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewAsyncProducer(kafkaAddrs, config)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer producer.AsyncClose()

	kafkaCh := make(chan *msclient.WrapKanban)
	go produce(ctx, kafkaCh, producer)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)
	for {
		select {
		case s := <-signalCh:
			fmt.Printf("received signal: %s", s.String())
			goto END
		case k := <-kanbanCh:
			if k != nil {
				kafkaCh <- k
			}
		}
	}
END:
}
