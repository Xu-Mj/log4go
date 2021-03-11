package log

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"time"
)

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	buf        []byte    // 存放日志内容
	out        io.Writer // 输出日志到目标流
	level      Level
	dateFormat string
}

func New(out io.Writer, dateFormat string, level Level) *Logger {
	if level > ERROR {
		fmt.Println("WARN: log level error, need 0-3. given: " + strconv.Itoa(int(level)))
		level = ERROR
	}
	logger := &Logger{out: out, level: level, dateFormat: "2006-01-02 15:04:05"}
	if dateFormat != "" {
		logger.dateFormat = dateFormat
	}
	return logger
}

// 生成输出语句
func (l *Logger) OutPut(logType Level, s string) error {
	//timeStr := fmt.Sprintf("[%s]", time.Now().Format(l.dateFormat))
	levelStr := l.getLevel(logType)
	// get called position
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown position"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	information := fmt.Sprintf("[%s] [%s] %s:%d: %s", time.Now().Format(l.dateFormat), levelStr, file, line, s)
	l.buf = l.buf[:0]
	l.buf = append(l.buf, information...)
	if len(information) == 0 || information[len(information)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	// 判断日志等级输出到目标文件
	if l.out != nil {
		if logType >= l.level {
			_, err := l.out.Write(l.buf)
			if err != nil {
				return err
			}
		}
	}
	fmt.Print(string(l.buf))
	return nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.OutPut(DEBUG, fmt.Sprintf(format, v...))
}
func (l *Logger) Info(format string, v ...interface{}) {
	l.OutPut(INFO, fmt.Sprintf(format, v...))
}
func (l *Logger) Warn(format string, v ...interface{}) {
	l.OutPut(WARN, fmt.Sprintf(format, v...))
}
func (l *Logger) Error(format string, v ...interface{}) {
	l.OutPut(ERROR, fmt.Sprintf(format, v...))
}

func (l *Logger) getLevel(logType Level) string {
	switch logType {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	default:
		return "ERROR"
	}
}
