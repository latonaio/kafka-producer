# kafka-producer
kafka-producerは、AIONのプラットフォーム上で動作するマイクロサービスから、[Kafka](https://kafka.apache.org/)へメッセージを送信するマイクロサービスです。  

## 概要
kafka-producerはその他のマイクロサービスからKanbanを受け取り、Kanbanに記述されたメタデータを元に、Kafkaへメッセージを送信します。  

## 動作環境
kafka-producerは、aion-coreのプラットフォーム上での動作を前提としています。 使用する際は、事前に下記の通りAIONの動作環境を用意してください。

* ARM CPU搭載のデバイス(NVIDIA Jetson シリーズ等)   
* OS: Linux Ubuntu OS   
* CPU: ARM64     
* Kubernetes   
* AION のリソース   

## セットアップ
このリポジトリをクローンし、`make`コマンドを用いてDocker container imageのビルドを行ってください。
```
$ cd /path/to/kafka-producer
$ make docker-build
```

## 起動方法
### 環境変数
|変数名|パラメータ|
|-|-|
|KAFKA_SERVER|"{your kafka service address}:{your kafka service port}"|

### デプロイ on AION
kafka-producerをデプロイするには、project.ymlに以下の構成でサービスを追加した上で、AIONのデプロイを行ってください。
```
  kafka-producer:
    scale: 1
    startup: yes
    always: yes
    env:
      KAFKA_SERVER: {your kafka service port}:{your kafka service port}
```

## I/O
### Input
kafka-producerを利用してKafkaへメッセージを送信するには、Kanbanのmetadataに以下の項目を含めて送信してください。  
`topic`, `key`には任意の文字列、`content`には任意の連想配列を指定できます。

```Python
metadata = {
    "topic": "TopicA",
    "key": "Key1",
    "content": {
         "calc_result": 100,
         "timestamp"" 20210209133650
    }
}

```

### Output
Kafkaには、`Input`で指定した`topic`, `key`に対して、`content`の内容が記録されます。
