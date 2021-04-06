/*
Project: MonitorGo
Author: Guo Kaikuo
Create time: 2021-04-06 10:58
IDE: GoLand
*/

package utils

import (
	"errors"
	"fmt"
	"github.com/StackExchange/wmi"
	"syscall"
	"unsafe"
)

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
type ProcessRes struct {
	ProcessName string
	ProcessID   uint32
}

/**
 * @Description: 基于进程快照的方式,获取所有Windows进程信息,但无法拿到具体的服务类进程
 * @return []ProcessRes
 */
func GetPidByProNameInWindows() []ProcessRes {
	var res []ProcessRes
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pHandle, _, _ := CreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) == -1 {
		return res
	}
	Process32Next := kernel32.NewProc("Process32Next")
	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := Process32Next.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			res = append(res, ProcessRes{
				ProcessName: string(proc.szExeFile[0:]),
				ProcessID:   uint32(proc.th32ProcessID),
			})
		} else {
			break
		}
	}
	CloseHandle := kernel32.NewProc("CloseHandle")
	_, _, _ = CloseHandle.Call(pHandle)
	return res
}

//https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-service
type Win32_Service struct {
	Name      string
	ProcessId uint32
	State     string
	Status    string
}

/**
 * @Description: 通过Wmi Win32_Service获取服务进程信息
 * @param serviceName 服务名
 * @return Win32_Service
 * @return error
 */
func GetPidByWmiService(serviceName string) (Win32_Service, error) {
	var dst []Win32_Service
	q := wmi.CreateQuery(&dst, fmt.Sprintf("WHERE name='%s'", serviceName))
	errQuery := wmi.Query(q, &dst)
	if errQuery != nil {
		return Win32_Service{}, errQuery
	} else if len(dst) == 0 {
		return Win32_Service{}, errors.New("query res is nil")
	} else {
		return dst[0], errQuery
	}
}

type Win32_Process struct {
	Name            string
	ProcessId       uint32
	ParentProcessId uint32
	Status          string
	CommandLine     string
	CSName          string
	ExecutablePath  string
	OSName          string
	PeakVirtualSize uint64
	ThreadCount     uint32
	VirtualSize     uint64
	WindowsVersion  string
}

/**
 * @Description: 通过Wmi Win32_Process获取进程信息
 * @param processName 进程名
 * @return Win32_Process
 * @return error
 */
func GetPidByWmiProcess(processName string) (Win32_Process, error) {
	var dst []Win32_Process
	q := wmi.CreateQuery(&dst, fmt.Sprintf("WHERE name='%s'", processName))
	errQuery := wmi.Query(q, &dst)
	if errQuery != nil {
		return Win32_Process{}, errQuery
	} else if len(dst) == 0 {
		return Win32_Process{}, errors.New("query res is nil")
	} else {
		return dst[0], errQuery
	}
}
