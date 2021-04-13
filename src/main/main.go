/*
Project: Monitor
Author: Guo Kaikuo
Create time: 2021-04-01 14:41
IDE: GoLand
*/

package main

import (
	"MonitorGo/src/logic"
	"MonitorGo/src/models"
	"MonitorGo/src/utils"
	"MonitorGo/src/views"
)

var mylog = utils.LogInit("main")

func main() {
	router := views.InitRouter()
	defer models.CloseAllFile(models.MyMonitorTask)
	go logic.RunMonitorTasks(models.MyMonitorTask)

	mylog.Info("main start")
	mylog.Info("start listen port 10000")

	router.Run(":10000")
}
