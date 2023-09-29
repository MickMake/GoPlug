package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MickMake/GoUnify/Only"
	"github.com/hashicorp/go-hclog"

	"github.com/MickMake/GoPlug/utils/Return"
)

type Logger struct {
	Name string
	out  io.Writer
	file *os.File
	log  hclog.Logger
}

func NewLogger(name string, filename string) (Logger, Return.Error) {
	var ret Logger
	var err Return.Error

	for range Only.Once {
		ret.Name = name

		ret.SetLogFile(filename)
		ret.log = hclog.New(&hclog.LoggerOptions{
			Name:                     name,
			Level:                    hclog.Debug,
			Output:                   ret.out,
			Mutex:                    nil,
			JSONFormat:               false,
			IncludeLocation:          false,
			AdditionalLocationOffset: 0,
			TimeFormat:               "",
			TimeFn:                   nil,
			DisableTime:              true,
			Color:                    hclog.ForceColor,
			ColorHeaderOnly:          false,
			ColorHeaderAndFields:     false,
			Exclude:                  nil,
			IndependentLevels:        false,
			SubloggerHook:            nil,
		})
		ret.Info("[%s] Logger started", name)
	}

	return ret, err
}

func (l *Logger) IsValid() Return.Error {
	if l == nil {
		var err Return.Error
		err.SetError("logger is not defined")
		return err
	}
	return Return.Ok
}

func (l *Logger) Gethclog() hclog.Logger {
	if l == nil {
		return nil
	}
	return l.log
}

func (l *Logger) GetLevel() hclog.Level {
	return l.log.GetLevel()
}

func (l *Logger) SetLevel(level hclog.Level) {
	l.log.SetLevel(level)
}

func (l *Logger) SetNamed(name string) {
	l.log.Named(name)
}

func (l *Logger) Close() {
	if l == nil {
		return
	}
	if l.file == nil {
		return
	}
	//goland:noinspection GoUnhandledErrorResult
	l.file.Close()
}

func (l *Logger) Info(msg string, args ...any) {
	if l.log != nil {
		l.log.Info(l.Name + " => " + fmt.Sprintf(msg, args...))
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	if l.log != nil {
		l.log.Debug(l.Name + " => " + fmt.Sprintf(msg, args...))
	}
}

func (l *Logger) SetName(name string) {
	if name != "" {
		l.Name = name
	}
}

func (l *Logger) SetLogFile(filename string) Return.Error {
	var err Return.Error

	for range Only.Once {
		if filename == "" {
			if l.file != nil {
				//goland:noinspection GoUnhandledErrorResult
				l.file.Close()
			}
			l.file = nil
			l.out = os.Stderr
			// l.out = l.log.StandardLogger(&hclog.StandardLoggerOptions{
			// 	InferLevels:              true,
			// 	InferLevelsWithTimestamp: true,
			// 	ForceLevel:               l.log.GetLevel(),
			// })
			log.SetOutput(l.out)
			log.SetPrefix("")
			log.SetFlags(log.Lshortfile)
			break
		}

		// @TODO - Maybe use l.log.StandardWriter() instead, then can use log.SetOutput()
		var e error
		l.out, e = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if e != nil {
			err.SetError("error opening file: %v", e)
			break
		}
		// l.out = l.file

		// l.out = l.log.StandardWriter(&hclog.StandardLoggerOptions{
		// 	InferLevels:              true,
		// 	InferLevelsWithTimestamp: true,
		// 	ForceLevel:               l.log.GetLevel(),
		// })
		log.SetOutput(l.out)
		log.SetPrefix("")
		log.SetFlags(log.Lshortfile)
		// mw := io.MultiWriter(os.Stdout, logFile)

	}

	return err
}
