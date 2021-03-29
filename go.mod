module bitbucket.org/latonaio/kafka-producer

go 1.15

require (
	bitbucket.org/latonaio/aion-core v0.9.4
	github.com/Shopify/sarama v1.26.2
	github.com/frankban/quicktest v1.10.2 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.11.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/robteix/protoconv v1.0.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/goleak v1.1.10
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

//replace github.com/protocolbuffers/protobuf-go => google.golang.org/protobuf v1.25.0
