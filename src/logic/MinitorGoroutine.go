/*
   Project: MonitorGo
   Author: Guo Kaikuo
   Create time: 2021-04-04 21:15
   IDE: GoLand
*/

package logic

import (
	"MonitorGo/src/models"
	"MonitorGo/src/utils"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"time"
)

// 全局扫描间隔 Minute
var ScanInterval uint = 30

func RunMonitorTasks(mt *models.MonitorTask) {
	for {
		for _, task := range mt.Tasks {
			//task.Logger.Printf("%s test\n",task.TaskName)
			//task_json, _ := json.Marshal(task)
			//fmt.Println(string(task_json))
			var current_pid uint32 = 0
			// 先在用户进程中检索,再在服务进程检索
			if proc, err := utils.GetPidByWmiProcess(task.TaskName); err == nil {
				current_pid = proc.ProcessId
			} else {
				if svc, err := utils.GetPidByWmiService(task.TaskName); err == nil && svc.State == "Running" {
					current_pid = svc.ProcessId
				}
			}
			task.PId = current_pid
			if task.PId > 0 {
				proc, err := process.NewProcess(int32(task.PId))
				if err != nil {
					//fmt.Println(err.Error())
					task.Logger.Printf("Process does not exist,PID=%d,err=%s\n", task.PId, err.Error())
				} else {
					cpu_rate, _ := proc.Percent(time.Second * 2)
					mem_rate, _ := proc.MemoryPercent()
					mem_info, _ := proc.MemoryInfo()
					mem_rss := mem_info.RSS
					mem_pretty := mem_rss / (1024 * 1024)
					procJson := &ProcessResJson{
						CpuRate:   fmt.Sprintf("%.2f", cpu_rate/float64(task.CoreCount)),
						MemRate:   fmt.Sprintf("%.2f", mem_rate),
						MemRss:    fmt.Sprintf("%dMB", mem_pretty),
						Timestamp: time.Now().Unix(),
					}
					jsonBytes, _ := json.Marshal(procJson)
					//fmt.Println(cpu_rate, mem_rate, mem_info.String())
					task.Logger.Println(string(jsonBytes))
				}
			}
		}
		time.Sleep(time.Second * time.Duration(ScanInterval))
	}
}

type ProcessResJson struct {
	CpuRate   string `json:"CpuRate"`
	MemRate   string `json:"MemRate"`
	MemRss    string `json:"MemRss"`
	Timestamp int64  `json:"Timestamp"`
}
