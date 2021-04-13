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
	"strings"
)

func GetResHandler(ctx *gin.Context) {
	resAll := make(map[string]interface{})
	for _, task := range models.MyMonitorTask.Tasks {
		if res := ReadResFromFile(task.File.Name()); res != nil {
			name := fmt.Sprintf("%s(%s)", task.TaskName, task.HostIp)
			resAll[name] = res
		}
		fmt.Println("for task", resAll)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"errCode": 0,
		"data":    resAll,
	})
}

func ReadResFromFile(file_path string) []logic.MonitorResJson {
	f, err := os.Open(file_path)
	if err != nil {
		fmt.Println("open error")
		return nil
	}
	defer f.Close()
	fs := bufio.NewScanner(f)
	var resAll []logic.MonitorResJson
	for fs.Scan() {
		line := fs.Text()
		fmt.Println(line)
		// process txt
		if each_res := ResParser(line); each_res != nil {
			resAll = append(resAll, *each_res)
		}
	}
	fmt.Println("resAll", resAll)
	return resAll
}

func ResParser(line string) *logic.MonitorResJson {
	space_pos := strings.LastIndex(line, " ")
	res_string := line[space_pos+1:]
	res_json := new(logic.MonitorResJson)
	err := json.Unmarshal([]byte(res_string), res_json)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		fmt.Println("res_json", res_json)
		return res_json
	}
}
