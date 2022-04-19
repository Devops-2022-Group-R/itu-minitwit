package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	goLog "log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type LogLevel int

const (
	Silent LogLevel = iota + 1
	Error
	Warn
	Info
	Fatal
	Panic
)

type Log struct {
	buffer   *bufio.Writer
	logger   *goLog.Logger
	logLevel LogLevel
}

func New() *Log {
	buf := bufio.NewWriter(os.Stdout)

	return NewUsingBuffer(buf)
}

// Create a new logger using a buffer, such as to a file
func NewUsingBuffer(buffer *bufio.Writer) *Log {
	return &Log{
		buffer: buffer,
		logger: goLog.New(buffer, "", 0),
	}
}

func getLogLevel(level LogLevel) string {
	switch level {
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	case Panic:
		return "PANIC"
	case Silent:
		return "SILENT"
	default:
		return "INFO"
	}
}

// Used internally to print in the desired format
func (l *Log) formatString(level LogLevel, v ...interface{}) string {
	jsonMap := map[string]interface{}{
		"time":  time.Now().UTC().Format("02-01-2006 15:04:05 -0700"),
		"level": getLogLevel(level),
		"log":   fmt.Sprint(v...),
	}
	jsonStr, _ := json.Marshal(jsonMap)

	return string(jsonStr)
}

func (l *Log) formatStringF(level LogLevel, format string, args ...interface{}) string {
	return l.formatString(level, fmt.Sprintf(format, args...))
}

func (l *Log) formatStringLn(level LogLevel, v ...interface{}) string {
	return l.formatString(level, fmt.Sprintln(v...))
}

func (l *Log) print(message string) {
	l.logger.Print(message)
	l.buffer.Flush()
}

func (l *Log) Printf(format string, args ...interface{}) {
	l.print(l.formatStringF(Info, format, args...))
}

func (l *Log) Print(v ...interface{}) {
	l.print(l.formatString(Info, v...))
}

func (l *Log) Println(message string) {
	l.print(l.formatStringLn(Info, message))
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	l.print(l.formatStringF(Fatal, format, args...))
	os.Exit(1)
}

func (l *Log) Fatal(v ...interface{}) {
	l.print(l.formatString(Fatal, v...))
	os.Exit(1)
}

func (l *Log) Fatalln(message string) {
	l.print(l.formatStringLn(Fatal, message))
	os.Exit(1)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.print(l.formatStringF(Warn, format, args...))
}

func (l *Log) Warnln(message string) {
	l.print(l.formatStringLn(Warn, message))
}

func (l *Log) LogMode(level logger.LogLevel) logger.Interface {
	switch level {
	case logger.Silent:
		l.logLevel = Silent
	case logger.Error:
		l.logLevel = Error
	case logger.Warn:
		l.logLevel = Warn
	case logger.Info:
		l.logLevel = Info
	}
	return l
}

func (l *Log) Info(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel >= Info {
		l.print("[GORM]" + l.formatStringF(Info, msg, data...))
	}
}

func (l *Log) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel >= Warn {
		l.print("[GORM]" + l.formatStringF(Warn, msg, data...))
	}
}

func (l *Log) Error(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel >= Error {
		l.print("[GORM]" + l.formatStringF(Error, msg, data...))
	}
}

// Trace print sql message
func (l *Log) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= Silent {
		return
	}

	elapsed := time.Since(begin)

	traceStr := "[GORM] %s\n[%.3fms] [rows:%v] %s"
	traceErrStr := "[GORM] %s %s\n[%.3fms] [rows:%v] %s"
	traceWarnStr := "[GORM] %s %s\n[%.3fms] [rows:%v] %s"

	SlowThreshold := 200 * time.Millisecond
	switch {
	case err != nil && l.logLevel >= Error && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > SlowThreshold && l.logLevel >= Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", SlowThreshold)
		if rows == -1 {
			l.Printf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.logLevel == Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func LogJson(level LogLevel, jsonMap map[string]interface{}) {
	jsonMap["level"] = getLogLevel(level)
	jsonMap["time"] = time.Now().UTC().Format("02-01-2006 15:04:05 -0700")
	jsonStr, err := json.Marshal(jsonMap)

	if err != nil {
		Logger.Warnf("log failed to marshal into json: %v", jsonMap)
		return
	}

	goLog.Printf("%s", jsonStr)
}

var (
	Logger = New()
)
