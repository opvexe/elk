package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

// tail的用法示例
func main() {
	fileName := "./my.log"
	config := tail.Config{
		ReOpen: true, // 重新打开 【文件大小500m，切日志时，重新打开日志】
		Location: &tail.SeekInfo{ //从文件的哪个地方开始读【从上次记录的位置开始读取】
			Offset: 0,
			Whence: 2,
		},
		MustExist: false, // 文件不存在不报错
		Poll:      false, // 轮询文件
		Follow:    true,  //实时跟踪 【可能有没读完的】
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Println("tail file closed", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("msg:", line.Text)
	}
}
