package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lawtech0902/go_gin_blog/backend/pkg/setting"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type Level int

const (
	INFO Level = iota
	WARNING
	ERROR
	FATAL
)

var (
	file *os.File
	err  error
)

func CustomLogger() gin.HandlerFunc {
	gin.DisableConsoleColor()
	
	file, err = os.OpenFile(GetLogFileFullPath("gin"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	
	g := gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		levelFlags := []string{"INFO", "WARN", "ERROR", "FATAL"}
		var level string
		status := params.StatusCode
		
		switch {
		case status > 499:
			level = levelFlags[FATAL]
		case status > 399:
			level = levelFlags[ERROR]
		case status > 299:
			level = levelFlags[WARNING]
		default:
			level = levelFlags[INFO]
		}
		
		return fmt.Sprintf("[%s] - %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			level,
			params.ClientIP,
			params.TimeStamp.Format(setting.AppInfo.TimeFormat),
			params.Method,
			params.Path,
			params.Request.Proto,
			status,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)
	})
	
	return g
}

func CloseLogFile() {
	if err = file.Close(); err != nil {
		return
	}
}

func GetLogFileFullPath(prefix string) string {
	return path.Join(setting.AppInfo.RootBasePath, setting.AppInfo.LogBasePath, GetLogFileName(prefix))
}

func GetLogFileName(prefix string) string {
	return fmt.Sprintf("%s_%s.log",
		prefix,
		time.Now().Format(strings.Split(setting.AppInfo.TimeFormat, " ")[0]))
}
