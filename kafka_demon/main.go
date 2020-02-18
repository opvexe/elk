package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //leader 和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个 分区
	config.Producer.Return.Successes = true                   // 成功交付的消息将在sucess channel 返回

	//构造一个消息
	msg := &sarama.ProducerMessage{
		Topic:     "web_log",
		Value:     sarama.StringEncoder("this is a test log"),
		Timestamp: time.Time{},
	}
	//连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("produce closed ,err:", err)
		return
	}
	defer client.Close()
	//发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed,err:", err)
		return
	}
	fmt.Println("pid:", pid, "offset:", offset)
}
