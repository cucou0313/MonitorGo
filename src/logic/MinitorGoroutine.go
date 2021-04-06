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
	"time"
)

// 全局扫描间隔
var ScanInterval uint = 1

func RunMonitorTasks(mt *models.MonitorTask) {
	for {
		for _, task := range mt.Tasks {
			pid_change := false
			var current_pid uint32 = 1
			// 先在用户进程中检索,再在服务进程检索
			if proc, err := utils.GetPidByWmiProcess(task.TaskName); err == nil {
				current_pid = proc.ProcessId
			} else {
				if svc, err := utils.GetPidByWmiService(task.TaskName); err == nil && svc.State == "Running" {
					current_pid = svc.ProcessId
				}
			}

			if task.PId != current_pid && task.PId != 0 {
				pid_change = true
			}
			task.PId = current_pid

			task_json, _ := json.Marshal(task)
			fmt.Println(string(task_json))
		}
		time.Sleep(time.Minute * time.Duration(ScanInterval))
	}
}
