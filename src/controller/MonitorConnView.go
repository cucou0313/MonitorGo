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
	"MonitorGo/src/utils"
	"github.com/gin-gonic/gin"
	"sort"
)

func ClientAckHandler(ctx *gin.Context) {
	ip := ctx.Query("ip")
	if ip != "" {
		utils.Mylog.Info(ip, " 发送心跳")
		if err := logic.IpInList(ip, &models.IpList); err == nil {
			models.IpList = append(models.IpList, ip)
			utils.Mylog.Info(ip, " 新客户端注册")
		} else {
			utils.Mylog.Info(ip, " 通讯正常")
		}
	}
	ctx.Status(200)
}

func GetAllClientHandler(ctx *gin.Context) {
	utils.Mylog.Info(ctx.ClientIP(), " 获取终端ip列表")

	sort.Strings(models.IpList)
	//服务端ip放第一
	res := append([]string{models.HostIp}, models.IpList...)
	ctx.JSON(200, gin.H{
		"Code": 0,
		"Data": res,
	})
}
