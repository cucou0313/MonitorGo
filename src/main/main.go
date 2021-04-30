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

func main() {
	defer models.CloseAllFile(models.MyMonitorTask)
	go logic.RunMonitorTasks(models.MyMonitorTask)

	router := views.InitRouter()
	utils.Mylog.Info("Monitor start")
	utils.Mylog.Info("Monitor server ip=", models.HostIp)
	utils.Mylog.Info("Start listen on port 10000")

	router.Run("0.0.0.0:10000")
}
