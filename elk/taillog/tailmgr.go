package taillog

import (
	"fmt"
	"shumin-project/elk-server/elk/etcd"
	"time"
)

type TaskMgr struct {
	logEntry    []*etcd.LogEntry      // 最开始的配置
	TaskMap     map[string]*LogTask   // 存储当前配置的map
	NewConfChan chan []*etcd.LogEntry // watch监控得到的新配置
}

var taskMgr TaskMgr // 全局任务管理者

func Init(LogConf []*etcd.LogEntry) {
	taskMgr = TaskMgr{
		logEntry:    LogConf,
		TaskMap:     make(map[string]*LogTask, etcd.EtcdConfig.TaskNum),
		NewConfChan: make(chan []*etcd.LogEntry),
	}
	for _, logConfig := range LogConf { // 遍历获取每一个task
		fmt.Println("当前配置------------>  Path:  " + logConfig.Path + "    Topic:  " + logConfig.Topic)
		taskObj := NewLogTask(logConfig.Path, logConfig.Topic)
		key := fmt.Sprintf(logConfig.Path + logConfig.Topic)
		taskMgr.TaskMap[key] = taskObj // 将每个任务放到map中做记录，方便监控收到的新配置进行判断
	}
	go taskMgr.NewRun()
}

func PushConfToChan() chan<- []*etcd.LogEntry {
	return taskMgr.NewConfChan
}

// watch传来新的配置时，开启新的任务
func (t *TaskMgr) NewRun() {
	for {
		select {
		case newLogConf := <-taskMgr.NewConfChan:
			// 将新任务与添加到配置中
			for _, logConfig := range newLogConf {
				key := fmt.Sprintf(logConfig.Path + logConfig.Topic)
				// 判断key是否在原来的map中
				if _, ok := t.TaskMap[key]; ok {
					continue
				} else {
					fmt.Println("新的配置------------>  Path:  " + logConfig.Path + "    Topic:  " + logConfig.Topic)
					// 开启新的任务
					taskObj := NewLogTask(logConfig.Path, logConfig.Topic)
					// 并在map中做记录
					t.TaskMap[key] = taskObj
					// 放到当前的LogEntry
					newLogEntry := etcd.LogEntry{
						Path:  logConfig.Path,
						Topic: logConfig.Topic,
					}
					t.logEntry = append(t.logEntry, &newLogEntry)
				}
			}

			// 移除之前的配置
			for index, oldConf := range t.logEntry {
				isExist := true
				for _, newConf := range newLogConf {
					if oldConf.Path == newConf.Path && oldConf.Topic == newConf.Topic {
						isExist = false
						continue
					}
				}
				if isExist {
					// 停止原来的配置
					key := fmt.Sprintf(oldConf.Path + oldConf.Topic)
					t.TaskMap[key].cancelFunc()
					// 从map中删除该task
					delete(t.TaskMap, key)
					// 从logEntry移除该task
					t.logEntry = append(t.logEntry[:index], t.logEntry[index+1:]...)
				}
			}
		default:
			time.Sleep(time.Microsecond)
		}
	}
}
