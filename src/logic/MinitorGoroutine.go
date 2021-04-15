/*
   Project: MonitorGo
   Author: Guo Kaikuo
   Create time: 2021-04-04 21:15
   IDE: GoLand
*/

package logic

import (
	"MonitorGo/src/models"
	"MonitorGo/src/utils"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"strconv"
	"time"
)

// ScanInterval 全局扫描间隔 Second
var ScanInterval uint = 30

func RunMonitorTasks(mt *models.MonitorTask) {
	for {
		timeStart := time.Now().Local()
		var SysCpuRate float64
		var SysMemRate float64
		var SysMemValue uint64
		if mt.CollectSysInfo {
			c, _ := cpu.Percent(time.Second*2, false)
			m, _ := mem.VirtualMemory()
			SysCpuRate = c[0]
			SysMemRate = m.UsedPercent
			SysMemValue = m.Used
		}
		for _, task := range mt.Tasks {
			if !task.Status {
				continue
			}
			//日志默认只写24小时, 防止写入过多无法导出
			timeNow := time.Now().Local()
			subDay := timeNow.Sub(task.TaskTime).Hours()
			if subDay > 24 {
				task.Status = false
				continue
			}
			var current_pid uint32 = 0
			// 先在用户进程中检索,再在服务进程检索
			if proc, err := utils.GetPidByWmiProcess(task.TaskName); err == nil {
				current_pid = proc.ProcessId
			} else {
				if svc, err := utils.GetPidByWmiService(task.TaskName); err == nil && svc.State == "Running" {
					current_pid = svc.ProcessId
				}
			}
			task.PId = current_pid
			resJson := new(MonitorResJson)
			if task.PId > 0 {
				proc, err := process.NewProcess(int32(task.PId))
				var pro_cpu float64
				var pro_mem float32
				var pro_rss uint64
				if err != nil {
					fmt.Println(err.Error())
					task.File.WriteString(fmt.Sprintf("Process does not exist,PID=%d,err=%s\n", task.PId, err.Error()))
				} else {
					pro_cpu, _ = proc.Percent(time.Second * 2)
					pro_cpu = pro_cpu / float64(mt.CoreCount)
					pro_mem, _ = proc.MemoryPercent()
					mem_info, _ := proc.MemoryInfo()
					pro_rss = mem_info.RSS
				}
				resJson.ProcessCpu = LineChartJson{
					Value: Decimal(pro_cpu),
					More:  fmt.Sprintf("(PID:%d)", task.PId),
				}
				resJson.ProcessMem = LineChartJson{
					Value: Decimal(float64(pro_mem)),
					More:  fmt.Sprintf("(%.2fMB)", GetPrettyMem(pro_rss, "MB")),
				}
				resJson.DateTime = time.Now().Format("2006-01-02 15:04:05")
			}

			if mt.CollectSysInfo {
				resJson.SystemCpu = LineChartJson{
					Value: Decimal(SysCpuRate),
					More:  fmt.Sprintf("(%d核%d线程)", mt.CoreCount, mt.LogicalCoreCount),
				}
				resJson.SystemMem = LineChartJson{
					Value: SysMemRate,
					More:  fmt.Sprintf("(%.2fGB)", GetPrettyMem(SysMemValue, "GB")),
				}
			}
			jsonBytes, _ := json.Marshal(resJson)
			//fmt.Println(string(jsonBytes))
			task.File.WriteString(string(jsonBytes) + "\n")
		}
		timeEnd := time.Now().Local()
		subSecond := timeEnd.Sub(timeStart).Milliseconds()
		trueSleep := time.Millisecond * time.Duration(int64(ScanInterval*1000)-subSecond)
		SendStatus(models.HostIp, models.ServerIp)
		time.Sleep(trueSleep)
	}
}

type MonitorResJson struct {
	SystemCpu  LineChartJson `json:"SystemCpu"`
	SystemMem  LineChartJson `json:"SystemMem"`
	ProcessCpu LineChartJson `json:"ProcessCpu"`
	ProcessMem LineChartJson `json:"ProcessMem"`
	DateTime   string        `json:"DateTime"`
}

// AppResJson 发给前端的数据结构
type AppResJson struct {
	ChartName  string          `json:"ChartName"`
	SystemCpu  []LineChartJson `json:"SystemCpu"`
	SystemMem  []LineChartJson `json:"SystemMem"`
	ProcessCpu []LineChartJson `json:"ProcessCpu"`
	ProcessMem []LineChartJson `json:"ProcessMem"`
	DateTime   []string        `json:"DateTime"`
}

type LineChartJson struct {
	Value float64 `json:"value"`
	More  string  `json:"more"`
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func GetPrettyMem(value uint64, format string) float64 {
	switch format {
	case "KB":
		return float64(value) / 1024
	case "MB":
		return float64(value) / (1024 * 1024)
	case "GB":
		return float64(value) / (1024 * 1024 * 1024)
	default:
		return float64(value)
	}
}
