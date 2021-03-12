package log4go

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
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
	mu         sync.Mutex // ensures atomic writes; protects the following fields
	buf        []byte     // 存放日志内容
	out        io.Writer  // 输出日志到目标流
	level      Level
	dateFormat string
	flag       bool
}

var std = New(nil, "", DEBUG, true)

func New(out io.Writer, dateFormat string, level Level, flag bool) *Logger {
	if level > ERROR {
		fmt.Println("WARN: log level error, need 0-3. given: " + strconv.Itoa(int(level)))
		level = ERROR
	}
	logger := &Logger{out: out, level: level, dateFormat: "2006-01-02 15:04:05", flag: flag}
	if dateFormat != "" {
		logger.dateFormat = dateFormat
	}
	return logger
}

// 生成输出语句
func (l *Logger) OutPut(logType Level, s string) error {
	now := time.Now()
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
	information := fmt.Sprintf("[%s] [%s] %s:%d: %s", now.Format(l.dateFormat), levelStr, file, line, s)
	l.buf = l.buf[:0]
	l.buf = append(l.buf, information...)
	if len(information) == 0 || information[len(information)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	// 判断日志等级输出到目标文件
	if l.out != nil {
		if logType >= l.level {
			l.mu.Lock()
			_, err := l.out.Write(l.buf)
			l.mu.Unlock()
			if err != nil {
				return err
			}
		}
	}
	if l.flag {
		_, _ = os.Stderr.Write(l.buf)
	}
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

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// Writer returns the output destination for the logger.
func (l *Logger) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

func (l *Logger) Level() Level {
	return l.level
}

// TODO deal error
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) LevelString() string {
	return l.getLevel(l.level)
}

func (l *Logger) DateFormat() string {
	return l.dateFormat
}

func (l *Logger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *Logger) Flag() bool {
	return l.flag
}

func (l *Logger) SetFlag(flag bool) {
	l.flag = flag
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = w
}

// Writer returns the output destination for the logger.
func Writer() io.Writer {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.out
}

func Debug(format string, v ...interface{}) {
	std.OutPut(DEBUG, fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	std.OutPut(INFO, fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	std.OutPut(WARN, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	std.OutPut(ERROR, fmt.Sprintf(format, v...))
}

func GetLevel() Level {
	return std.level
}

// TODO deal error
func SetLevel(level Level) {
	std.level = level
}

func LevelString() string {
	return std.getLevel(std.level)
}

func DateFormat() string {
	return std.dateFormat
}

func SetDateFormat(format string) {
	std.dateFormat = format
}

func Flag() bool {
	return std.flag
}

func SetFlag(flag bool) {
	std.flag = flag
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
