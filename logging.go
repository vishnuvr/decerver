package deCerver

import ()

/*
type LogSystem struct {
	Modules map[string]core.Logger
	// TODO Implement
	logLevel core.LogLevel
	logFile string
	logReader io.Reader
	logWriter io.Writer

	subs      map[string]core.LogSub

	mutex *sync.Mutex
}

func (dc *DeCerver) initLogSystem() {
	logSys := &LogSystem{}
	logSys.logFile = dc.config.LogFile
	dc.logSys = logSys
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

func NewLogSub() core.LogSub {
	ls := &LogSub{
		Channel:  make(chan string),
		SubId:    0,
		LogLevel: 5,
		Enabled:  true,
	}
	return ls
}

func NewEthLogger() *EthLogger {

	el := &EthLogger{}
	el.mutex = &sync.Mutex{};
	el.logLevel = 5
	el.logReader, el.logWriter = io.Pipe()

	ethlog.AddLogSystem(ethlog.NewStdLogSystem(el.logWriter, log.LstdFlags, el.logLevel))

	go func(el *EthLogger) {
		scanner := bufio.NewScanner(el.logReader)
		for scanner.Scan() {
			text := scanner.Text()
			el.mutex.Lock()
			for _, sub := range el.subs {
				sub.Channel <- text
			}
			el.mutex.Unlock()
		}
	}(el)
	return el
	return nil
}

func (el *EthLogger) AddSub(sub *LogSub) {
	el.mutex.Lock()
	el.subs = append(el.subs, sub)
	el.mutex.Unlock()
}

func (el *EthLogger) RemoveSub(sub *LogSub) {
	el.mutex.Lock()
	theIdx := -1
	for idx, s := range el.subs {
		if sub.SubId == s.SubId {
			theIdx = idx
			break
		}
	}
	if theIdx >= 0 {
		el.subs = append(el.subs[:theIdx], el.subs[theIdx+1:]...)
	}
	el.mutex.Unlock()
}*/
