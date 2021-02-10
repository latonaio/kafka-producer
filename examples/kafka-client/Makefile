GO_SRCS := $(shell find . -type f -name '*.go')

docker-build: $(GO_SRCS)
	bash ./scripts/build.sh

go-build: $(GO_SRCS)
	go build ./
