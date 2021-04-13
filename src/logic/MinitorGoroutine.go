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

// ScanInterval 全局扫描间隔 Minute
var ScanInterval uint = 1

func RunMonitorTasks(mt *models.MonitorTask) {
	for {
		for _, task := range mt.Tasks {
			//task.Logger.Printf("%s test\n",task.TaskName)
			//task_json, _ := json.Marshal(task)
			//fmt.Println(string(task_json))
			if !task.Status {
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
			if task.PId > 0 {
				proc, err := process.NewProcess(int32(task.PId))
				sys_cpu, _ := cpu.Percent(time.Second*2, false)
				sys_mem, _ := mem.VirtualMemory()
				var pro_cpu float64
				var pro_mem float32
				var pro_rss uint64
				if err != nil {
					fmt.Println(err.Error())
					task.File.WriteString(fmt.Sprintf("Process does not exist,PID=%d,err=%s\n", task.PId, err.Error()))
				} else {
					pro_cpu, _ = proc.Percent(time.Second * 2)
					pro_cpu = pro_cpu / float64(task.CoreCount)
					pro_mem, _ = proc.MemoryPercent()
					mem_info, _ := proc.MemoryInfo()
					pro_rss = mem_info.RSS
				}
				resJson := &MonitorResJson{
					SystemCPU: LineChartJson{
						Value: Decimal(sys_cpu[0]),
						More:  fmt.Sprintf("(%d核%d线程)", task.CoreCount, task.LogicalCoreCount),
					},
					SystemMem: LineChartJson{
						Value: sys_mem.UsedPercent,
						More:  fmt.Sprintf("(%.2fGB)", GetPrettyMem(sys_mem.Used, "GB")),
					},
					ProcessCPU: LineChartJson{
						Value: Decimal(pro_cpu),
						More:  "",
					},
					ProcessMem: LineChartJson{
						Value: Decimal(float64(pro_mem)),
						More:  fmt.Sprintf("(%.2fMB)", GetPrettyMem(pro_rss, "MB")),
					},
					DataTime: time.Now().Format("2006-01-02 15:04:05"),
				}
				jsonBytes, _ := json.Marshal(resJson)
				//fmt.Println(string(jsonBytes))
				task.File.WriteString(string(jsonBytes) + "\n")
			}
		}
		time.Sleep(time.Second * 10 * time.Duration(ScanInterval))
	}
}

type MonitorResJson struct {
	SystemCPU  LineChartJson `json:"SystemCPU"`
	SystemMem  LineChartJson `json:"SystemMem"`
	ProcessCPU LineChartJson `json:"ProcessCPU"`
	ProcessMem LineChartJson `json:"ProcessMem"`
	DataTime   string        `json:"DataTime"`
}

// AppResJson 发给前端的数据结构
type AppResJson struct {
	ChartName  string          `json:"ChartName"`
	SystemCPU  []LineChartJson `json:"SystemCPU"`
	SystemMem  []LineChartJson `json:"SystemMem"`
	ProcessCPU []LineChartJson `json:"ProcessCPU"`
	ProcessMem []LineChartJson `json:"ProcessMem"`
	DataTime   []string        `json:"DataTime"`
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
