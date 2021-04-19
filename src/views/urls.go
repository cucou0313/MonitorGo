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
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())

	task := router.Group("/task")
	{
		//测试任务信息
		task.GET("/check", controller.TestTaskHandler)
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

	conn := router.Group("/conn")
	{
		//Listen Client
		conn.GET("/listen_client", controller.ListClientHandler)
		//获取在线终端
		conn.GET("/get_all_ip", controller.GetAllClientHandler)
	}

	res := router.Group("/res")
	{
		//读取所有运行监控结果
		res.GET("/get_all", controller.GetAllResHandler)
		//读取指定任务结果
		res.GET("/get_one", controller.GetOneResHandler)
		//下载指定任务结果
		res.GET("/down_file", controller.DownResHandler)
	}

	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//接收客户端发送的origin （重要！）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")
		//允许类型校验
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
