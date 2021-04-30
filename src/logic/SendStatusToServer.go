/*
Project: MonitorGo
Author: Guo Kaikuo
Create time: 2021-04-15 11:26
IDE: GoLand
*/

package logic

import (
	"MonitorGo/src/utils"
	"errors"
	"fmt"
	"net/http"
	"time"
)

/**
 * @Description: 向Server发送自己的ip,维护在线终端表
	方便起见,直接放在 RunMonitorTasks 进行定时
 * @param myip
 * @param serip
*/
func SendStatus(myip, serip string) {
	if myip != serip {
		utils.Mylog.Info("向服务端发送心跳")
		//fmt.Printf("client %s 向 server %s 发送心跳\n", myip, serip)
		request, _ := http.NewRequest(
			"GET",
			fmt.Sprintf("http://%s:10000/conn/listen_client?ip=%s", serip, myip),
			nil,
		)
		do, err := (&http.Client{
			//2秒超时,防止等待耗时
			Timeout: 2 * time.Second,
		}).Do(request)
		if err != nil {
			utils.Mylog.Info("心跳失败,err=", err.Error())
			//fmt.Println(err.Error())
			return
		} else {
			utils.Mylog.Info("通讯正常")
			//fmt.Println("通讯正常")
		}
		do.Body.Close()
	}
}

func IpInList(ip string, iplist *[]string) error {
	for _, each_ip := range *iplist {
		if each_ip == ip {
			return errors.New("this ip has existed")
		}
	}
	return nil
}
