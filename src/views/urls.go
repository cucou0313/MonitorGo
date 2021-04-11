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
		task.GET("/get", controller.GetTaskHandler)
		task.GET("/add", controller.AddTaskHandler)
		task.GET("/del", controller.DelTaskHandler)
		task.GET("/start", controller.StartTaskHandler)
		task.GET("/stop", controller.StopTaskHandler)
	}

	res := router.Group("/res")
	{
		res.GET("/get", controller.GetResHandler)
	}

	return router
}
