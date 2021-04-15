/*
Project: MonitorGo
Author: Guo Kaikuo
Create time: 2021-04-15 11:14
IDE: GoLand
*/

package utils

import (
	"net"
	"strings"
)

func GetHostIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	addr := conn.LocalAddr().String()
	addr_list := strings.SplitN(addr, ":", 2)
	return addr_list[0]
}
