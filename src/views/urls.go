/*
Project: Monitor
Author: Guo Kaikuo
Create time: 2021-04-01 15:44
IDE: GoLand
*/

package views

import (
	"MonitorGo/src/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	task := router.Group("/task")
	{
		//测试任务信息
		task.GET("/test", controller.TestTaskHandler)
		//开启系统信息
		task.GET("/open_system", controller.OpenSysHandler)
		//监控任务信息
		task.GET("/get", controller.GetTaskHandler)
		//添加监控任务
		task.GET("/add", controller.AddTaskHandler)
		//删除监控任务
		task.GET("/del", controller.DelTaskHandler)
		//启动监控任务
		task.GET("/start", controller.StartTaskHandler)
		//停止监控任务
		task.GET("/stop", controller.StopTaskHandler)
	}

	res := router.Group("/res")
	{
		//读取监控结果
		res.GET("/get", controller.GetResHandler)
	}

	return router
}
