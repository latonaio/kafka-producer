module bitbucket.org/latonaio/kafka-producer

go 1.15

require (
	bitbucket.org/latonaio/aion-core v0.9.3
	github.com/Shopify/sarama v1.27.2
	github.com/golang/mock v1.4.0
	github.com/golang/protobuf v1.4.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.8.1
	github.com/robteix/protoconv v1.0.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/goleak v1.1.10
	google.golang.org/protobuf v1.24.0
)

//replace github.com/protocolbuffers/protobuf-go => google.golang.org/protobuf v1.25.0
