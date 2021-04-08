/*
   Project: MonitorGo
   Author: Guo Kaikuo
   Create time: 2021-04-04 12:38
   IDE: GoLand
*/

package models

import (
	"MonitorGo/src/utils"
	"errors"
	"github.com/go-basic/uuid"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"os"
	"path/filepath"
	"time"
)

var MyMonitorTask = new(MonitorTask)
var task_dir = filepath.Join(utils.GetExecPath(), "TaskLog")
var MyCoreCount = GetCoreCount()

// 创建task日志路径
func init() {
	if !utils.FileExist(task_dir) {
		_ = utils.CreateDir(task_dir)
	}
}

type TaskInterface interface {
	AddTask(name string, ip string, scanInt int) error
	DelTask(id string) error
}

type TaskInfo struct {
	TaskId   string `json:"TaskId"`
	TaskName string `json:"TaskName"`
	TaskTime int64  `json:"TaskTime"`
	HostIp   string `json:"HostIp"`
	// 初始pid=0
	PId       uint32 `json:"PId"`
	CoreCount int    `json:"CoreCount"`
	// 该任务的日志采集文件指针
	File *os.File `json:"File,omitempty"`
	// 该任务的log.Logger指针
	Logger *log.Logger `json:"Logger,omitempty"`
}

type MonitorTask struct {
	Tasks []TaskInfo
}

/**
 * @Description: 生成Monitor任务
 * @receiver mt
 * @param name
 * @param ip
 * @return error
 */
func (mt *MonitorTask) AddTask(name string, ip string) error {
	if name == "" || ip == "" {
		return errors.New("task name or ip is empty")
	}
	for _, task := range mt.Tasks {
		if task.TaskName == name && task.HostIp == ip {
			return errors.New("this monitor task has been existed")
		}
	}
	ip_path := filepath.Join(task_dir, ip)
	if !utils.FileExist(ip_path) {
		_ = utils.CreateDir(ip_path)
	}
	id := uuid.New()
	file_path := ip_path + string(os.PathSeparator) + name + "_" + id + ".log"
	f, l := utils.CreateFile(file_path)
	//fmt.Println(f, l)
	newTask := &TaskInfo{
		TaskId:    id,
		TaskName:  name,
		TaskTime:  time.Now().Unix(),
		HostIp:    ip,
		PId:       0,
		CoreCount: MyCoreCount,
		File:      f,
		Logger:    l,
	}
	//task_json, _ := json.Marshal(newTask)
	//l.Println(string(task_json))
	mt.Tasks = append(mt.Tasks, *newTask)
	return nil
}

func (mt *MonitorTask) DelTask(id string) error {
	if id == "" {
		return errors.New("task id is empty")
	}
	del_falg := false
	for index, task := range mt.Tasks {
		if task.TaskId == id {
			mt.Tasks = append(mt.Tasks[:index], mt.Tasks[index+1:]...)
			del_falg = true
			task.File.Close()
			break
		}
	}
	if del_falg {
		return nil
	} else {
		return errors.New("this monitor id does not exist")
	}
}

func CloseAllFile(mt *MonitorTask) {
	for _, task := range mt.Tasks {
		for {
			err := task.File.Close()
			if err != nil {
				continue
			}
			break
		}
	}
}

/**
 * @Description: 操作系统内核数,非逻辑内核
 * @return int
 */
func GetCoreCount() int {
	if core_count, err := cpu.Counts(false); err != nil {
		return 1
	} else {
		return core_count
	}
}
