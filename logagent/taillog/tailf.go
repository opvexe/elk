package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"sync"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

// 日志收集任务定义
type TailCollection struct {
	Topic   string `json:"topic"`    //主题
	LogPath string `json:"log_path"` //日志路径
}

// tail 任务对象
type TailObject struct {
	taiObject    *tail.Tail
	collecction  TailCollection
	status       int
	closeChannel chan int
}

// 发送到kafka消息
type KafkaMessage struct {
	Message string `json:"message"`
	Ip      string `json:"ip"`
}

// 消息体
type TextMessage struct {
	Message KafkaMessage
	Topic   string
}

// tail 任务对象管理
type TailObjectManager struct {
	tailObjects    []*TailObject
	messageChannel chan *TextMessage
	lock           sync.Mutex
}

var (
	hostIP           string
	taiObjectManager *TailObjectManager
)

// 初始化
func NewTail(collecctions []TailCollection, size int, ip string) (err error) {
	taiObjectManager = &TailObjectManager{
		messageChannel: make(chan *TextMessage, size),
	}
	if len(collecctions) == 0 {
		fmt.Println("collection task is nill")
	}
	hostIP = ip
	for _, v := range collecctions {
		createTask(v)
	}
	return
}

// 创建tail task
func createTask(collection TailCollection) {
	obj, err := tail.TailFile(collection.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	})
	if err != nil {
		fmt.Println("tail create err:", err)
	}
	tailObject := &TailObject{
		taiObject:    obj,
		collecction:  collection,
		closeChannel: make(chan int, 1),
	}
	// 开启go程去监听
	go readFromTail(tailObject, collection.Topic)
	taiObjectManager.tailObjects = append(taiObjectManager.tailObjects, tailObject)
}

// 读取监听日志文件内容
func readFromTail(tailObject *TailObject, topic string) {
	for true {
		select {
		//任务正常运行
		case lineMessage, ok := <-tailObject.taiObject.Lines:
			if !ok {
				fmt.Println("read objecct err:", topic)
				continue
			}
			if len(lineMessage.Text) == 0 {
				continue
			}
			kafkaMessage := KafkaMessage{
				Message: lineMessage.Text,
				Ip:      hostIP,
			}
			kafkaObject := &TextMessage{
				Message: kafkaMessage,
				Topic:   topic,
			}
			taiObjectManager.messageChannel <- kafkaObject
		case <-tailObject.closeChannel:
			fmt.Println("tail object will exited")
			return
		}
	}
}

// 从channel中获取一行的数据
func GetOneLine() (msg *TextMessage) {
	msg = <-taiObjectManager.messageChannel
	return
}

// 更新tail任务
func UpdataTailTask(collections []TailCollection) (err error) {
	taiObjectManager.lock.Lock()
	defer taiObjectManager.lock.Unlock()
	for _, newColect := range collections {
		//判断tail是否运行状态
		var isRun bool = false
		for _, oldColect := range taiObjectManager.tailObjects {
			if newColect.LogPath == oldColect.collecction.LogPath {
				isRun = true
				break
			}
		}
		// 如果tail任务不存在,创建新任务
		if isRun == false {
			createTask(newColect)
		}
	}
	//更新tail任务管理内容
	var tailObjects []*TailObject
	for _, oldObject := range taiObjectManager.tailObjects {
		oldObject.status = StatusDelete
		for _, newColl := range collections {
			if newColl.LogPath == oldObject.collecction.LogPath {
				oldObject.status = StatusNormal
				break
			}
		}
		if oldObject.status == StatusDelete {
			oldObject.closeChannel <- 1
			continue
		}
		tailObjects = append(tailObjects, oldObject)
	}
	taiObjectManager.tailObjects = tailObjects
	return
}
