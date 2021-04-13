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

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"SystemCPU":  []int{8, 9, 7, 10, 8, 20},
			"SystemMem":  []int{33, 36, 40, 44, 55, 66},
			"ProcessCPU": []int{1, 2, 1, 0, 0, 1},
			"ProcessMem": []int{1, 2, 2, 1, 1, 0},
			"ProcessRss": []int{40, 60, 400, 100, 50, 60},
			"DateTime":   []int{1, 2, 3, 4, 5, 6},
		})
	})

	task := router.Group("/task")
	{
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
