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
	resAll := make(map[string]interface{})
	for i, task := range models.MyMonitorTask.Tasks {
		if res := ReadResFromFile(task.File.Name()); res != nil {
			res.ChartName = fmt.Sprintf("%s(%s)", task.TaskName, task.HostIp)
			resAll[fmt.Sprintf("data%d", i)] = res
		}
	}
	js, _ := json.Marshal(resAll)
	fmt.Println("for task", string(js))
	ctx.JSON(http.StatusOK, gin.H{
		"errCode": 0,
		"data":    resAll,
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
		fmt.Println(line)
		// process txt
		if each_res := ResParser(line); each_res != nil {
			res.SystemCPU = append(res.ProcessCPU, each_res.SystemCPU)
			res.SystemMem = append(res.SystemMem, each_res.SystemMem)
			res.ProcessCPU = append(res.ProcessCPU, each_res.ProcessCPU)
			res.ProcessMem = append(res.ProcessMem, each_res.ProcessMem)
			res.DataTime = append(res.DataTime, each_res.DataTime)
		}
	}
	fmt.Println("resAll", res)
	return res
}

func ResParser(line string) *logic.MonitorResJson {
	res_json := new(logic.MonitorResJson)
	err := json.Unmarshal([]byte(line), res_json)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		fmt.Println("res_json", res_json)
		return res_json
	}
}
