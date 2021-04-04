/*
   Project: MonitorGo
   Author: Guo Kaikuo
   Create time: 2021-04-04 21:15
   IDE: GoLand
*/

package logic

import (
	"MonitorGo/src/models"
	"encoding/json"
	"fmt"
	"time"
)

func RunMonitorTasks(mt *models.MonitorTask) {
	for {
		for _, task := range mt.Tasks {
			task_json, _ := json.Marshal(task)
			fmt.Println(string(task_json))
		}
		time.Sleep(time.Second * 30)
	}
}
