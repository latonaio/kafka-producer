# syntax = docker/dockerfile:1.0.0-experimental
# Build Container
FROM golang:1.15.2 as builder

ENV GO111MODULE on
ENV GOPRIVATE "bitbucket.org/latonaio"
WORKDIR /go/src/bitbucket.org/latonaio

COPY go.mod .

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
RUN mkdir /root/.ssh/ && touch /root/.ssh/known_hosts && ssh-keyscan -t rsa bitbucket.org >> /root/.ssh/known_hosts
RUN --mount=type=secret,id=ssh,target=/root/.ssh/id_rsa go mod download

COPY . .

RUN go build -o kafka-producer .


# Runtime Container
FROM alpine:3.12

RUN apk add --no-cache libc6-compat tzdata

COPY --from=builder /go/src/bitbucket.org/latonaio/kafka-producer .

CMD ["./kafka-producer"]
