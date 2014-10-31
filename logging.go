package deCerver

import (
	"fmt"
	"log"
	"os"
)

const LOGGER_PREFIX = "[DECERVER] "
const LOG_MINLEVEL = 0
const LOG_MAXLEVEL = 5

// For convenience.
var logger *log.Logger

var defaultLoggerFile *os.File = os.Stdout

type LogSystem struct {
	// Decerver Core
	DCLogger *log.Logger
	
	Modules map[string]*log.Logger
	// TODO Implement
	logLevel int
	logFile string
}

func (dc *DeCerver) initLogSystem() {
	dc.logSys = &LogSystem{}
	dc.logSys.logFile = dc.config.LogFile
	logger = log.New(defaultLoggerFile,LOGGER_PREFIX, dc.config.LogLevel)
	dc.logSys.DCLogger = logger
	dc.logSys.Modules = make(map[string]*log.Logger)
	dc.logSys.logLevel = dc.config.LogLevel
}

func (ls *LogSystem) openLogFile() *os.File {
	file, err := os.OpenFile(ls.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening log file: %v", err))
	}
	return file
}

func (ls *LogSystem) AddLogger(id string, logger *log.Logger){
	logger.SetPrefix(LOGGER_PREFIX + "(" + logger.Prefix() + ") ")
	ls.Modules[id] = logger
}

func (ls *LogSystem) RemoveLogger(id string) bool {
	if _ , ok := ls.Modules[id]; !ok {
		return false
	}
	delete(ls.Modules,id)
	return true
}

func (ls *LogSystem) SetLogLevel(newLevel int){
	if LOG_MINLEVEL > newLevel || newLevel > LOG_MAXLEVEL {
		ls.DCLogger.Printf("The log level must be between %d and %d. New value: %d", LOG_MINLEVEL, LOG_MAXLEVEL,newLevel)
		return
	}
	if ls.logLevel == newLevel {
		ls.DCLogger.Printf("The log level is already set to '%d'", newLevel)
		return
	}
	ls.logLevel = newLevel
}

func (ls *LogSystem) LogLevel() int {
	return ls.logLevel
}