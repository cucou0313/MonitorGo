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
	"os"
	"path/filepath"
	"time"
)

var MyMonitorTask = &MonitorTask{
	CoreCount:        GetCoreCount(),
	LogicalCoreCount: GetLogicalCoreCount(),
	CollectSysInfo:   false,
}
var IpList []string
var HostIp = utils.GetHostIp()
var ServerIp = "10.65.10.187"
var task_dir = filepath.Join(utils.GetExecPath(), "TaskLog")

// 创建task日志路径
func init() {
	if !utils.FileExist(task_dir) {
		_ = utils.CreateDir(task_dir)
	}
}

type TaskInterface interface {
	AddTask(name string, ip string) error
	StartTask(id string) error
	StopTask(id string) error
	DelTask(id string) error
	GetRunningNum() int
}

type TaskInfo struct {
	TaskId   string `json:"TaskId"`
	TaskName string `json:"TaskName"`
	// 设置持续写日志24小时
	TaskTime time.Time `json:"TaskTime"`
	HostIp   string    `json:"HostIp"`
	// Status任务是否启动
	Status bool `json:"Status"`
	// 初始pid=0
	PId uint32 `json:"PId"`
	// 该任务的日志采集文件指针
	File *os.File `json:"File,omitempty"`
}

type MonitorTask struct {
	Tasks []*TaskInfo
	// 系统内核数
	CoreCount        int
	LogicalCoreCount int
	// 是否同步采集系统信息
	CollectSysInfo bool
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
	f := utils.CreateFile(file_path)
	//fmt.Println(f, l)
	newTask := &TaskInfo{
		TaskId:   id,
		TaskName: name,
		TaskTime: time.Now(),
		Status:   false,
		HostIp:   ip,
		PId:      0,
		File:     f,
	}
	//task_json, _ := json.Marshal(newTask)
	//l.Println(string(task_json))
	mt.Tasks = append(mt.Tasks, newTask)
	return nil
}

func (mt *MonitorTask) StartTask(id string) error {
	if id == "" {
		return errors.New("task id is empty")
	}
	for _, task := range mt.Tasks {
		if task.TaskId == id {
			task.Status = true
			task.TaskTime = time.Now().Local()
			return nil
		}
	}
	return errors.New("this monitor task id does not exist")
}

func (mt *MonitorTask) StopTask(id string) error {
	if id == "" {
		return errors.New("task id is empty")
	}
	for _, task := range mt.Tasks {
		if task.TaskId == id {
			task.Status = false
			return nil
		}
	}
	return errors.New("this monitor task id does not exist")
}

func (mt *MonitorTask) DelTask(id string) error {
	if id == "" {
		return errors.New("task id is empty")
	}
	for index, task := range mt.Tasks {
		if task.TaskId == id {
			mt.Tasks = append(mt.Tasks[:index], mt.Tasks[index+1:]...)
			task.File.Close()
			return nil
		}
	}
	return errors.New("this monitor task id does not exist")
}

func (mt *MonitorTask) GetRunningNum() int {
	var runningNum int
	for _, task := range mt.Tasks {
		if task.Status == true {
			runningNum++
		}
	}
	return runningNum
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

func GetLogicalCoreCount() int {
	if core_count, err := cpu.Counts(true); err != nil {
		return 1
	} else {
		return core_count
	}
}
