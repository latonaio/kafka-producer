# syntax = docker/dockerfile:1.0.0-experimental
# Build Container
FROM golang:1.15.2 as builder

ENV GO111MODULE on
ENV GOPRIVATE "bitbucket.org/latonaio"
WORKDIR /go/src/bitbucket.org/latonaio

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o kafka-consumer


# Runtime Container
FROM alpine:3.12

RUN apk add --no-cache libc6-compat tzdata

COPY --from=builder /go/src/bitbucket.org/latonaio/kafka-consumer .

CMD ["./kafka-consumer"]
