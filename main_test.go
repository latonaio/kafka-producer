package main

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama/mocks"
	"github.com/robteix/protoconv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"google.golang.org/protobuf/types/known/structpb"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
)

func genTestData(topic string, key string, content map[string]interface{}) *msclient.WrapKanban {
	kanbanData := &msclient.WrapKanban{}
	innerStruct := protoconv.NewStruct()
	for k, v := range content {
		innerStruct.Set(k, protoconv.StringVal(v.(string)))
	}
	strct := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"topic":   protoconv.StringVal(topic),
			"key":     protoconv.StringVal(key),
			"content": innerStruct.Value(),
		},
	}

	kanbanData.Metadata = strct

	return kanbanData
}

func TestKanbanToKafkaMsg(t *testing.T) {
	topic := "Test.A"
	key := "service-a:001"
	content := map[string]interface{}{
		"data": "hello",
	}

	kanbanData := genTestData(topic, key, content)

	kafkaMsg, err := kanbanToKafkaMsg(kanbanData)
	if err != nil {
		t.Errorf("fialed test %#v\n", err)
	}

	assert.Equal(t, topic, kafkaMsg.topic)
	assert.Equal(t, key, kafkaMsg.key)
	assert.Equal(t, content, kafkaMsg.content)

}

func TestProduce(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataCh := make(chan *msclient.WrapKanban)

	config := mocks.NewTestConfig()
	m := mocks.NewAsyncProducer(t, config)
	defer m.AsyncClose()

	m.ExpectInputAndSucceed()

	go produce(ctx, dataCh, m)

	topic := "Test.A"
	key := "service-a:001"
	content := map[string]interface{}{
		"data": "hello",
	}
	kanbanData := genTestData(topic, key, content)

	dataCh <- kanbanData
	time.Sleep(2 * time.Second)

}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
