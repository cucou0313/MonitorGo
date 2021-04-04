/*
Project: Monitor
Author: Guo Kaikuo
Create time: 2021-04-01 14:41
IDE: GoLand
*/

package main

import (
	"MonitorGo/src/utils"
	"fmt"
)

var mylog = utils.LogInit("main")

func main() {
	test := make(map[string]interface{})

	test["a"] = map[string]interface{}{
		"b": map[string]interface{}{"c": 1},
	}
	fmt.Println(test)
	//router := views.InitRouter()
	//defer models.CloseAllFile(models.MyMonitorTask)
	//go logic.RunMonitorTasks(models.MyMonitorTask)
	//
	//mylog.Info("main start")
	//mylog.Info("start listen port 10000")
	//
	//router.Run(":10000")
}
