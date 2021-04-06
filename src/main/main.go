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
	"github.com/shirou/gopsutil/v3/winservices"
)

var mylog = utils.LogInit("main")

func main() {
	//var pid uint32
	//if res, err := utils.GetPidByWmiProcess("goland64.exe"); err == nil {
	//	fmt.Println("ok")
	//	jsonBytes, _ := json.MarshalIndent(res, "", "\t")
	//	fmt.Println(string(jsonBytes))
	//	pid = res.ProcessId
	//	fmt.Println(pid)
	//} else {
	//	fmt.Println("err")
	//	fmt.Println(err)
	//}
	svcs, err := winservices.ListServices()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(svcs)
	}
	//router := views.InitRouter()
	//defer models.CloseAllFile(models.MyMonitorTask)
	//go logic.RunMonitorTasks(models.MyMonitorTask)
	//
	//mylog.Info("main start")
	//mylog.Info("start listen port 10000")
	//
	//router.Run(":10000")
}
