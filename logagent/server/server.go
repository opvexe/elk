package server

import (
	"encoding/json"
	"fmt"
	"shumin-project/elk-server/logagent/kafka"
	"shumin-project/elk-server/logagent/taillog"
	"time"
)

// 启动logagent服务
func Run() (err error) {
	fmt.Println("服务启动...")
	for true{
		//获取每一行日志
		msg := taillog.GetOneLine()
		//发送一行日志到kafka
		err = sendToKafka(msg.Message,msg.Topic)
		if err!=nil{
			fmt.Printf("send to kafka msg:[%v] topic:[%v] failed, err:[%v]", msg.Message, msg.Topic, err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}


func sendToKafka(msg taillog.KafkaMessage,topic string) (err error) {
	smsg,err := json.Marshal(&msg)
	if err!=nil{
		fmt.Printf("send to kafka marshal failed --> msg: [%v], topic:[%s], error: %s", msg, topic, err)
		return
	}
	err = kafka.SendMessageToKafka(string(smsg),topic)
	return
}
