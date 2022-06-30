package base

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"github.com/tietang/props/kvs"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var lfh *utils.LineNumLogrusHook

func init() {
	// 日志格式
	formatter := &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
		DisableColors:   false,
		ForceFormatting: true,
	}
	// 控制台高亮
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		TimestampStyle:  "blue",
	})
	logrus.SetFormatter(formatter)

	// 日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetLevel(logrus.DebugLevel)

	// 开启调用函数、文件、代码行信息的输出
	logrus.SetReportCaller(true)

	// 日志文件和滚动配置
	SetLineNumLogrusHook()

}

func SetLineNumLogrusHook() {
	lfh = utils.NewLineNumLogrusHook()
	lfh.EnableFileNameLog = true
	lfh.EnableFuncNameLog = true
	logrus.AddHook(lfh)
}

//将滚动日志writer共享给iris glog output
var log_writer io.Writer

//初始化log配置，配置logrus日志文件滚动生成和
func InitLog(conf kvs.ConfigSource) {
	//设置日志输出级别
	level, err := logrus.ParseLevel(conf.GetDefault("log.level", "info"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	if conf.GetBoolDefault("log.enableLineLog", true) {
		lfh.EnableFileNameLog = true
		lfh.EnableFuncNameLog = true
	} else {
		lfh.EnableFileNameLog = false
		lfh.EnableFuncNameLog = false
	}

	//配置日志输出目录
	logDir := conf.GetDefault("log.dir", "./logs")
	logTestDir, err := conf.Get("log.test.dir")
	if err == nil {
		logDir = logTestDir
	}
	logPath := logDir //+ "/logs"
	logFilePath, _ := filepath.Abs(logPath)
	logrus.Infof("log dir: %s", logFilePath)
	logFileName := conf.GetDefault("log.file.name", "red-envelop")
	maxAge := conf.GetDurationDefault("log.max.age", time.Hour*24)
	rotationTime := conf.GetDurationDefault("log.rotation.time", time.Hour*1)
	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	//设置滚动日志输出writer
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", err)
	}
	//设置日志文件输出的日志格式
	formatter := &logrus.TextFormatter{}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		return function, fmt.Sprintf("%s/%s:%d", f, filename, frame.Line)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, formatter)

	logrus.AddHook(lfHook)
	//
	log_writer = writer
}
