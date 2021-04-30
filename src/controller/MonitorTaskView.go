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
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func TestTaskHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	var Msg string
	if name == "" {
		Msg = "输入的进程名为空,则表示只监听系统信息"
	} else {
		utils.Mylog.Info(ctx.ClientIP(), " 申请进程检查,", name)
		var current_pid uint32 = 0
		// 先在用户进程中检索,再在服务进程检索
		if proc, err := utils.GetPidByWmiProcess(name); err == nil {
			current_pid = proc.ProcessId
		} else {
			if svc, err := utils.GetPidByWmiService(name); err == nil && svc.State == "Running" {
				current_pid = svc.ProcessId
			}
		}
		utils.Mylog.Info("检查结果pid=", current_pid)

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
	utils.Mylog.Info(ctx.ClientIP(), " 切换系统收集开关, flag=", sys_flag)

	ctx.JSON(http.StatusOK, gin.H{
		"Code": 0,
		"Msg":  "Switch system info collection success",
	})
}

func GetTaskHandler(ctx *gin.Context) {
	utils.Mylog.Info(ctx.ClientIP(), " 获取所有任务信息")

	var data []models.TaskInfo
	var runCount = 0
	var totalCount = len(models.MyMonitorTask.Tasks)
	var taskStatus = gin.H{}
	for _, task := range models.MyMonitorTask.Tasks {
		if task.Status {
			runCount++
		}
		taskStatus[task.TaskId] = task.Status
		data = append(data, *task)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Code":             0,
		"Data":             data,
		"CoreCount":        models.MyMonitorTask.CoreCount,
		"LogicalCoreCount": models.MyMonitorTask.LogicalCoreCount,
		"CollectSysInfo":   models.MyMonitorTask.CollectSysInfo,
		"SysInfo":          models.GetHostInfo(),
		"runCount":         runCount,
		"totalCount":       totalCount,
		"taskStatus":       taskStatus,
	})
}

func AddTaskHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	utils.Mylog.Info(ctx.ClientIP(), " 添加新任务, name=", name)

	err := models.MyMonitorTask.AddTask(name)
	if err != nil {
		utils.Mylog.Info("添加任务失败, err=", err.Error())

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		utils.Mylog.Info("添加任务成功")

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Add new task successfully,name=%s ", name),
		})
	}
}

func StartTaskHandler(ctx *gin.Context) {

	if models.MyMonitorTask.GetRunningNum() >= 4 {
		utils.Mylog.Info("超出可同时监听数")

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  "超出可同时监听数！",
		})
		return
	}
	id := ctx.Query("id")
	utils.Mylog.Info(ctx.ClientIP(), " 启动任务, taskid=", id)

	err := models.MyMonitorTask.StartTask(id)
	if err != nil {
		utils.Mylog.Info("启动任务失败, err=", err.Error())

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		utils.Mylog.Info("启动任务成功")

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Start this task successfully,id=%s", id),
		})
	}
}

func StopTaskHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	utils.Mylog.Info(ctx.ClientIP(), " 停止任务, taskid=", id)

	err := models.MyMonitorTask.StopTask(id)
	if err != nil {
		utils.Mylog.Info("停止任务失败, err=", err.Error())

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		utils.Mylog.Info("停止任务成功")

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Stop this task successfully,id=%s", id),
		})
	}
}

func DelTaskHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	utils.Mylog.Info(ctx.ClientIP(), " 删除任务, taskid=", id)

	err := models.MyMonitorTask.DelTask(id)
	if err != nil {
		utils.Mylog.Info("删除任务失败, err=", err.Error())

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  err.Error(),
		})
	} else {
		utils.Mylog.Info("删除任务成功")

		ctx.JSON(http.StatusOK, gin.H{
			"Code": 0,
			"Msg":  fmt.Sprintf("Del this task successfully,id=%s", id),
		})
	}
}
