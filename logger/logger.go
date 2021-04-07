package logger

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func InitLog (level string) {
	formatter := runtime.Formatter{ChildFormatter: &log.JSONFormatter{}}
	formatter.File = true
	formatter.Line = true
	log.SetFormatter(&formatter)

	lvl := log.WarnLevel
	switch strings.ToUpper(level) {
	case "TRACE":
		lvl = log.TraceLevel
		break
	case "INFO":
		lvl = log.InfoLevel
		break
	case "DEBUG":
		lvl = log.DebugLevel
		break
	case "WARN":
	case "WARNING":
		lvl = log.WarnLevel
		break
	case "ERROR":
		lvl = log.ErrorLevel
		break
	}
	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
}