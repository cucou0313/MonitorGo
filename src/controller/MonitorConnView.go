/*
Project: MonitorGo
Author: Guo Kaikuo
Create time: 2021-04-15 11:03
IDE: GoLand
*/

package controller

import (
	"MonitorGo/src/logic"
	"MonitorGo/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"sort"
)

func ListClientHandler(ctx *gin.Context) {
	ip := ctx.Query("ip")
	if ip != "" {
		fmt.Printf("client %s 申请注册\n", ip)
		if err := logic.IpInList(ip, &models.IpList); err == nil {
			models.IpList = append(models.IpList, ip)
			fmt.Printf("client %s 注册成功\n", ip)
		} else {
			fmt.Printf("client %s 已经注册过\n", ip)
		}
	}
	fmt.Println(models.IpList)
	ctx.Status(200)
}
func GetAllClientHandler(ctx *gin.Context) {
	sort.Strings(models.IpList)
	//服务端ip放第一
	res := append([]string{models.HostIp}, models.IpList...)
	ctx.JSON(200, gin.H{
		"Code": 0,
		"Data": res,
	})
}

func SendStatusHandler(ctx *gin.Context) {
	logic.SendStatus(models.HostIp, models.ServerIp)
}
