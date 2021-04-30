/*
   Project: MonitorGo
   Author: Guo Kaikuo
   Create time: 2021-04-04 22:32
   IDE: GoLand
*/

package controller

import (
	"MonitorGo/src/logic"
	"MonitorGo/src/models"
	"MonitorGo/src/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func GetAllResHandler(ctx *gin.Context) {
	//resAll := make(map[string]interface{})
	utils.Mylog.Info(ctx.ClientIP(), " 获取所有运行任务结果")
	resAll := gin.H{
		"data1": new(logic.AppResJson),
		"data2": new(logic.AppResJson),
		"data3": new(logic.AppResJson),
		"data4": new(logic.AppResJson),
	}
	var chartNum = 1
	for _, task := range models.MyMonitorTask.Tasks {
		if task.Status == true {
			if res := ReadResFromFile(task.File.Name()); res != nil {
				res.ChartName = fmt.Sprintf("%d.%s(%s)", chartNum, task.TaskName, models.HostIp)
				resAll[fmt.Sprintf("data%d", chartNum)] = res
				chartNum++
			}
		}
	}
	//js, _ := json.Marshal(resAll)
	//fmt.Println("for task", string(js))
	ctx.JSON(http.StatusOK, gin.H{
		"Code": 0,
		"Data": resAll,
	})
}

func ReadResFromFile(file_path string) *logic.AppResJson {
	f, err := os.Open(file_path)
	if err != nil {
		utils.Mylog.Error(file_path, " 文件打开错误:", err.Error())
		return nil
	}
	defer f.Close()
	fs := bufio.NewScanner(f)
	var res = &logic.AppResJson{}
	for fs.Scan() {
		line := fs.Text()
		//fmt.Println(line)
		// process txt
		if each_res := ResParser(line); each_res != nil {
			res.SystemCpu = append(res.SystemCpu, each_res.SystemCpu)
			res.SystemMem = append(res.SystemMem, each_res.SystemMem)
			res.ProcessCpu = append(res.ProcessCpu, each_res.ProcessCpu)
			res.ProcessMem = append(res.ProcessMem, each_res.ProcessMem)
			res.DateTime = append(res.DateTime, each_res.DateTime)
		}
	}
	//fmt.Println("resAll", res)
	return res
}

func ResParser(line string) *logic.MonitorResJson {
	res_json := new(logic.MonitorResJson)
	err := json.Unmarshal([]byte(line), res_json)
	if err != nil {
		utils.Mylog.Error("log解析错误:", err.Error())
		return nil
	} else {
		//fmt.Println("res_json", res_json)
		return res_json
	}
}

func GetOneResHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	id := ctx.Query("id")
	utils.Mylog.Info(ctx.ClientIP(), " 获取指定任务结果", name, id)
	if id == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  "查询数据错误",
		})
	}
	filepath := models.Task_dir + string(os.PathSeparator) + name + "_" + id + ".log"
	utils.Mylog.Info("任务文件地址:", filepath)
	resAll := gin.H{
		"data1": new(logic.AppResJson),
		"data2": new(logic.AppResJson),
		"data3": new(logic.AppResJson),
		"data4": new(logic.AppResJson),
	}
	if res := ReadResFromFile(filepath); res != nil {
		res.ChartName = fmt.Sprintf("%d.%s(%s)", 1, name, models.HostIp)
		resAll["data1"] = res
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Code": 0,
		"Data": resAll,
	})
}

func DownResHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	id := ctx.Query("id")
	utils.Mylog.Info(ctx.ClientIP(), " 申请下载任务日志", name, id)
	if id == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": 1,
			"Msg":  "数据错误",
		})
	}
	filepath := models.Task_dir + string(os.PathSeparator) + name + "_" + id + ".log"
	filename := name + "_" + time.Now().Format("2006-01-02 15:04:05") + ".log"
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(filepath)
}

func SysInfoHandler(ctx *gin.Context) {
	utils.Mylog.Info(ctx.ClientIP(), " 获取系统信息")
	ctx.JSON(200, models.GetHostInfo())
}
