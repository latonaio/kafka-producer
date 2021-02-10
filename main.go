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
	"github.com/pkg/errors"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"bitbucket.org/latonaio/aion-core/pkg/log"
)

var msName string = "kafka-producer"

type kafkaConfig struct {
	addr string `envconfig:"KAFKA_SERVER" default:"localhost:9092"`
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

func produce(ctx context.Context, dataCh <-chan *msclient.WrapKanban, kafkaConf *kafkaConfig) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewAsyncProducer([]string{kafkaConf.addr}, config)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer producer.AsyncClose()

	go func() {
		for {
			select {
			case <-childCtx.Done():
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
	}()
	<-ctx.Done()
}

func main() {
	config := &kafkaConfig{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal("cannot load environment")
	}
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
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM)

	kafkaCh := make(chan *msclient.WrapKanban)
	go produce(ctx, kafkaCh, config)

	for {
		select {
		case s := <-signalCh:
			fmt.Printf("received signal: %s", s.String())
			goto END
		case k := <-kanbanCh:
			kafkaCh <- k
		}
	}
END:
}
