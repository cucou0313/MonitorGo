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
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetResHandler(ctx *gin.Context) {
	//resAll := make(map[string]interface{})
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
				res.ChartName = fmt.Sprintf("%d.%s(%s)", chartNum, task.TaskName, task.HostIp)
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
		fmt.Println("open error")
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
		fmt.Println(err.Error())
		return nil
	} else {
		//fmt.Println("res_json", res_json)
		return res_json
	}
}
