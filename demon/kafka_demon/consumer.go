package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

// kafka消费者
type KafkaClient struct {
	client sarama.Consumer
	topic []string // topic列表
}

var (
	kafkaClient *KafkaClient
	runTopic []string	//正在运行的topic任务
)

func NewKafkaConsumer(hosts []string,topics []string) (err error) {
	consumer,err := sarama.NewConsumer(hosts,nil)
	if err!=nil{
		return
	}
	kafkaClient = &KafkaClient{
		client: consumer,
		topic:  topics,
	}
	//创建topic任务
	for _, topic := range topics{
		go createTopicTask(topic)
	}
	return
}

// 判断topic是否在运行
func isTopicExists(topic string) (ok bool) {
	ok = false
	for _, t := range runTopic{
		if t == topic {
			ok = true
			break
		}
	}
	return
}

//创建kafka topic消费者任务
func createTopicTask(topic string)  {
	time.Sleep(time.Second)
	if isTopicExists(topic) == true{
		return
	}
	var wg sync.WaitGroup
	partionList,err := kafkaClient.client.Partitions(topic)
	if err!=nil{
		logs.Error("get topic: [%s] partitions failed, err: %s", topic, err)
		return
	}
	for partion := range partionList{
		pc, err := kafkaClient.client.ConsumePartition(topic, int32(partion), sarama.OffsetNewest)
		if err!=nil{
			logs.Warn("topic: [%s] start cousumer partition failed, err: %s", topic, err)
			continue
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages(){
				fmt.Println(msg)
			}
			wg.Done()
		}(pc)
	}
	// 启动成功的topic任务添加到列表中去
	runTopic = append(runTopic,topic)
	wg.Wait()
}

//更新topic任务
func UpdataTopicTask(topics []string) (err error) {
	for _, newTopic := range topics {
		var topicStauts = false
		for _, oldTopic := range kafkaClient.topic {
			if newTopic == oldTopic {
				topicStauts = true
				break
			}
		}
		if topicStauts == false {
			go createTopicTask(newTopic)
		}
	}
	return
}