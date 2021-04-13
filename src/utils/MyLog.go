/*
   Project: MonitorGo
   Author: Guo Kaikuo
   Create time: 2021-04-03 15:21
   IDE: GoLand
*/

package utils

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var projectName = "Monitor"
var log_dir = filepath.Join(GetExecPath(), "GoLog")

func LogInit(file_name string) *logrus.Logger {
	if !FileExist(log_dir) {
		_ = CreateDir(log_dir)
	}
	file_path := log_dir + string(os.PathSeparator) + projectName + "_" + file_name
	/* 日志轮转相关函数
	WithLinkName 为最新的日志建立软连接
	WithRotationTime 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	WithMaxAge 设置文件清理前的最长保存时间
	WithRotationCount 设置文件清理前最多保存的个数
	*/
	writer, _ := rotatelogs.New(
		file_path+".log.%Y-%m-%d",
		rotatelogs.WithLinkName(log_dir+string(os.PathSeparator)+file_name+"_latest.log"),
		rotatelogs.WithRotationCount(7),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(new(MyFormatter))
	logger.SetOutput(writer)
	return logger
}

/*
自定义logrus的输出格式,官方说法它提供了接口
type Formatter interface {
	Format(*Entry) ([]byte, error)
}
需要用户自己实现这个接口即可
*/
type MyFormatter struct {
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timeStamp := time.Now().Format("2006/01/02 15:04:05")
	logLevel := strings.ToUpper(entry.Level.String())
	funcName := entry.Caller.Function
	funcLine := entry.Caller.Line
	msg := fmt.Sprintf("%s [%s] %s:%d: %s\n", timeStamp, logLevel, filepath.Base(funcName), funcLine, entry.Message)
	return []byte(msg), nil
}

//CreateDir  文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	_ = os.Chmod(path, os.ModePerm)
	return nil
}

//FileExist  判断文件夹/文件是否存在  存在返回 true
func FileExist(file_path string) bool {
	_, err := os.Stat(file_path)
	return err == nil || os.IsExist(err)
}

/**
 * @Description: 获取程序运行目录,解决log和config的相对路径问题
 * @return string 程序运行目录
 */
func GetExecPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return path
}

/**
 * @Description: Monitor任务采集日志保存
 * @param file_path  日志文件路径
 * @return *os.File 读写日志文件的指针,用于主程序退出时close
 */
func CreateFile(file_path string) *os.File {
	f, err := os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	return f
}
