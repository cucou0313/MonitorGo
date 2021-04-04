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

func GetTaskHandler(ctx *gin.Context) {
	var data []map[string]interface{}
	for _, task := range models.MyMonitorTask.Tasks {
		var tmp map[string]interface{}
		if task_json, err := json.Marshal(task); err == nil {
			//data=append(data,string(task_json))
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
	scanInt := ctx.DefaultQuery("scanInt", "60")
	si, err := strconv.Atoi(scanInt)
	if err != nil || (si < 30) {
		fmt.Println("Scan Interval time is invalid")
	}
	err = models.MyMonitorTask.AddTask(name, ip, si)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"errCode": 1,
			"errMsg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"errCode": 0,
			"errMsg":  fmt.Sprintf("add new task success,name=%s,ip=%s", name, ip),
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
			"errMsg":  fmt.Sprintf("del this task success,id=%s", id),
		})
	}
}
