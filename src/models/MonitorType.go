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
	"github.com/shirou/gopsutil/v3/host"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var MyMonitorTask = &MonitorTask{
	CoreCount:        GetCoreCount(),
	LogicalCoreCount: GetLogicalCoreCount(),
	CollectSysInfo:   false,
}
var IpList []string
var HostIp = utils.GetHostIp()
var ServerIp = HostIp
var Task_dir = filepath.Join(utils.GetExecPath(), "TaskLog")

// 创建task日志路径
func init() {
	if !utils.FileExist(Task_dir) {
		_ = utils.CreateDir(Task_dir)
	}
}

type TaskInterface interface {
	AddTask(name string) error
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
func (mt *MonitorTask) AddTask(name string) error {
	if name == "" {
		return errors.New("task name is empty")
	}
	name_lower := strings.ToLower(name)
	for _, task := range mt.Tasks {
		if strings.ToLower(task.TaskName) == name_lower {
			return errors.New("this monitor task has been existed")
		}
	}
	id := uuid.New()
	file_path := Task_dir + string(os.PathSeparator) + name + "_" + id + ".log"
	f := utils.CreateFile(file_path)
	//fmt.Println(f, l)
	newTask := &TaskInfo{
		TaskId:   id,
		TaskName: name,
		TaskTime: time.Now(),
		Status:   false,
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

func GetHostInfo() HostInfo {
	info, err := host.Info()
	res := &HostInfo{}
	if err != nil {
		return *res
	}
	res.Hostname = info.Hostname
	res.BootTime = time.Unix(int64(info.BootTime), 0).Format("2006-01-02 15:04:05")
	res.OS = info.OS
	res.Procs = info.Procs
	res.Platform = info.Platform
	res.PlatformFamily = info.PlatformFamily
	res.PlatformVersion = info.PlatformVersion
	res.KernelArch = info.KernelArch
	return *res
}

type HostInfo struct {
	Hostname        string `json:"hostname"`
	BootTime        string `json:"bootTime"`
	OS              string `json:"os"`              // ex: freebsd, linux
	Procs           uint64 `json:"procs"`           // number of processes
	Platform        string `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily  string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion string `json:"platformVersion"` // version of the complete OS
	KernelArch      string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
}
