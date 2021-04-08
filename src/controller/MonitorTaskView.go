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
)

var mylog = utils.LogInit("test")

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
		"data": data,
	})
}

func AddTaskHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	ip := ctx.Query("ip")
	err := models.MyMonitorTask.AddTask(name, ip)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"errCode": 1,
			"errMsg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"errCode": 0,
			"Msg":     fmt.Sprintf("Add new task successfully,name=%s,ip=%s", name, ip),
		})
	}
}

func DelTaskHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	err := models.MyMonitorTask.DelTask(id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"errCode": 1,
			"errMsg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"errCode": 0,
			"Msg":     fmt.Sprintf("Del this task successfully,id=%s", id),
		})
	}
}
