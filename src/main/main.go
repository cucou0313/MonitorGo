/*
Project: Monitor
Author: Guo Kaikuo
Create time: 2021-04-01 14:41
IDE: GoLand
*/

package main

import (
	"MonitorGo/src/utils"
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

var mylog = utils.LogInit("main")

type ulong int32
type ulong_ptr uintptr

type PROCESSENTRY32 struct {
	dwSize              ulong
	cntUsage            ulong
	th32ProcessID       ulong
	th32DefaultHeapID   ulong_ptr
	th32ModuleID        ulong
	cntThreads          ulong
	th32ParentProcessID ulong
	pcPriClassBase      ulong
	dwFlags             ulong
	szExeFile           [260]byte
}

func main() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pHandle, _, _ := CreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) == -1 {
		return
	}
	Process32Next := kernel32.NewProc("Process32Next")
	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := Process32Next.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			fmt.Println("ProcessName : " + string(proc.szExeFile[0:]))
			fmt.Println("th32ModuleID : " + strconv.Itoa(int(proc.th32ModuleID)))
			fmt.Println("ProcessID : " + strconv.Itoa(int(proc.th32ProcessID)))
		} else {
			break
		}
	}
	CloseHandle := kernel32.NewProc("CloseHandle")
	_, _, _ = CloseHandle.Call(pHandle)
	//router := views.InitRouter()
	//defer models.CloseAllFile(models.MyMonitorTask)
	//go logic.RunMonitorTasks(models.MyMonitorTask)
	//
	//mylog.Info("main start")
	//mylog.Info("start listen port 10000")
	//
	//router.Run(":10000")
}
