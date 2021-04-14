/*
   Project: MonitorGo
   Author: Guo Kaikuo
   Create time: 2021-04-04 22:31
   IDE: GoLand
*/

package controller

import (
	"MonitorGo/src/models"
	"MonitorGo/src/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var mylog = utils.LogInit("test")

func TestTaskHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	var Msg string
	if name == "" {
		Msg = "输入的进程名为空,则表示只监听系统信息"
	} else {
		var current_pid uint32 = 0
		// 先在用户进程中检索,再在服务进程检索
		if proc, err := utils.GetPidByWmiProcess(name); err == nil {
			current_pid = proc.ProcessId
		} else {
			if svc, err := utils.GetPidByWmiService(name); err == nil && svc.State == "Running" {
				current_pid = svc.ProcessId
			}
		}
		if current_pid > 0 {
			Msg = fmt.Sprintf("找到了名为:'%s'的运行进程,进程ID:%d", name, current_pid)
		} else {
			Msg = fmt.Sprintf("现在没有找到名为:'%s'的进程,该进程可能没有运行,任务可在运行后自动监听", name)
		}
	}
	ctx.JSON(200, gin.H{
		"Code": 0,
		"Msg":  Msg,
	})
}

func OpenSysHandler(ctx *gin.Context) {
	flag := ctx.Query("flag")
	sys_flag, _ := strconv.ParseBool(flag)
	models.MyMonitorTask.CollectSysInfo = sys_flag
	fmt.Println(models.MyMonitorTask.CollectSysInfo)
	ctx.JSON(http.StatusOK, gin.H{
		"Code": 0,
		"Msg":  "Open system info collection success",
	})
}

func GetTaskHandler(ctx *gin.Context) {
	var data []map[string]interface{}
	for _, task := range models.MyMonitorTask.Tasks {
		var tmp map[string]interface{}
		if task_json, err := json.Marshal(task); err == nil {
			_ = json.Unmarshal(task_json, &tmp)
			data = append(data, tmp)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Code": 0,
		"Data": data,
	})
}

func AddTaskHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	ip := ctx.Query("ip")
	err := models.MyMonitorTask.AddTask(name, ip)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Add new task successfully,name=%s,ip=%s", name, ip),
		})
	}
}

func StartTaskHandler(ctx *gin.Context) {
	if models.MyMonitorTask.GetRunningNum() >= 4 {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  "超出可同时监听数！",
		})
		return
	}
	id := ctx.Query("id")
	err := models.MyMonitorTask.StartTask(id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Start this task successfully,id=%s", id),
		})
	}
}

func StopTaskHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	err := models.MyMonitorTask.StopTask(id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Stop this task successfully,id=%s", id),
		})
	}
}

func DelTaskHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	err := models.MyMonitorTask.DelTask(id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Del this task successfully,id=%s", id),
		})
	}
}
