package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type LogLevel uint8

var printf = fmt.Printf

// logFunc represents a log function

type LoggerColor struct {
	level       int
	trace       string
	debug       string
	info        string
	warn        string
	err         string
	fatal       string
	ExitOnFatal bool
}

type Logger struct {
	LogEnable   bool
	FormatLog   string
	Prefix      string
	LogFlag     int
	LogAvalable []int
	logColor    *LoggerColor
	Logger      *log.Logger
}

var Log *Logger

func New(out string, prefix string, logAvalable []int, logEnable bool) (logger *Logger, er error) {
	colorLogger := &LoggerColor{
		level:       INFO,
		trace:       "[TRACE]",
		debug:       "[DEBUG]",
		info:        "[INFO]",
		warn:        "[WARN]",
		err:         "[ERROR]",
		fatal:       "[FATAL]",
		ExitOnFatal: false,
	}
	logger = &Logger{
		LogFlag:     log.Ldate | log.Ltime | log.Lshortfile,
		Prefix:      prefix,
		LogAvalable: logAvalable,
		logColor:    colorLogger,
		LogEnable:   logEnable,
		FormatLog:   "%s",
	}
	Log = logger
	if logEnable {
		Writer, err := initWrite(out)
		if err != nil {
			er = err
			return
		}
		logger.Logger = log.New(Writer, logger.Prefix, logger.LogFlag)
	}
	return
}

func initWrite(out string) (file *os.File, err error) {
	if out != "" {
		file, err = os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return
		}
		log.SetOutput(file)
	}
	return
}

func (l *Logger) log(level string, a ...interface{}) {
	if l.LogEnable {
		l.Logger.Printf("%s %s %s\n", level, "=>", fmt.Sprint(a...))
		return
	}
	printf("%s %s %s %s\n", l.now(),
		level,
		"=>",
		fmt.Sprint(a...))
	return
}

func (l *Logger) logf(level string, format string, a ...interface{}) {
	if l.LogEnable {
		l.Logger.Printf("%s %s %s\n", level, "=>", fmt.Sprintf(format, a...))
		return
	}
	printf("%s %s %s %s\n", l.now(),
		level,
		"=>",
		fmt.Sprintf(format, a...))
}

func (l *Logger) Print(v ...interface{}) (e error) {
	l.log(l.logColor.info, v...)
	return
}

func (l *Logger) Error(v ...interface{}) (e error) {
	if !l.checkAvalableLogging(ERROR) {
		return
	}
	l.log(l.logColor.err, v...)
	return
}

func (l *Logger) Warning(v ...interface{}) (e error) {
	if !l.checkAvalableLogging(WARN) {
		return
	}
	l.log(l.logColor.warn, v...)
	return
}

func (l *Logger) Info(v ...interface{}) (err error) {
	if !l.checkAvalableLogging(INFO) {
		return
	}
	l.log(l.logColor.info, v...)
	return
}

func (l *Logger) Debug(v ...interface{}) (err error) {
	if !l.checkAvalableLogging(DEBUG) {
		return
	}
	l.log(l.logColor.debug, v...)
	return
}

func (l *Logger) Errorf(format string, a ...interface{}) (e error) {
	if !l.checkAvalableLogging(ERROR) {
		return
	}
	l.logf(l.logColor.err, format, a...)
	return
}

func (l *Logger) Warningf(format string, a ...interface{}) (e error) {
	if !l.checkAvalableLogging(WARN) {
		return
	}
	l.logf(l.logColor.warn, format, a...)
	return
}

func (l *Logger) Infof(format string, a ...interface{}) (e error) {
	if !l.checkAvalableLogging(INFO) {
		return
	}
	l.logf(l.logColor.info, format, a...)
	return
}

func (l *Logger) Debugf(format string, a ...interface{}) (err error) {
	if !l.checkAvalableLogging(DEBUG) {
		return
	}
	l.logf(l.logColor.debug, format, a...)
	return
}

func (l *Logger) SetAvalableLogging(logAvalable int) {
	l.LogAvalable = append(l.LogAvalable, logAvalable)
}

func (l *Logger) now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l *Logger) checkAvalableLogging(avalableLog int) bool {
	if len(l.LogAvalable) == 0 {
		return false
	}
	for _, enableLog := range l.LogAvalable {
		if enableLog == avalableLog {
			return true
		}
	}
	return false
}
