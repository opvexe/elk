package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

// 生产者对象
var kafkaClient sarama.SyncProducer

// 初始化生产者
func NewKafka(address []string) (err error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Version = sarama.V0_10_0_1  //注意设置版本

	kafkaClient, err = sarama.NewSyncProducer(address, kafkaConfig)
	if err != nil {
		fmt.Println("kafka producer err:", err)
		return
	}
	return
}

// 发送消息
func SendMessageToKafka(msg, topic string) (err error) {
	msgObject := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(msg),
		Timestamp: time.Time{},
	}
	pid, offset, err := kafkaClient.SendMessage(msgObject)
	if err != nil {
		fmt.Printf("send msg to kafka topic:[%v] msg:[%v] faile, %v", msg, topic, err)
		return
	}
	fmt.Printf("topic: [%v] pid: [%v], offset: [%v]", topic, pid, offset)
	return
}
