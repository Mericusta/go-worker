package logger

import (
	"fmt"
	"os"

	"github.com/go-worker/config"
	"github.com/go-worker/global"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

var LogChan chan []byte
var LogControl chan bool
var LogFile *os.File

func init() {
	LogChan = make(chan []byte)
	LogControl = make(chan bool)
	initLogFile()
}

func initLogFile() {
	var initLogFileError error
	var initLogFileErrorInfo string
	if utility.IsExist(config.WorkerConfig.LogFilePath) {
		LogFile, initLogFileError = os.OpenFile(config.WorkerConfig.LogFilePath, os.O_APPEND, os.ModeAppend)
		initLogFileErrorInfo = ui.CommonError5
		if initLogFileError == nil {
			return
		}
	} else {
		LogFile, initLogFileError = utility.CreateFile(config.WorkerConfig.LogFilePath)
		if initLogFileError == nil {
			return
		}
		initLogFileErrorInfo = ui.CommonError4
	}
	ui.OutputErrorInfo(initLogFileErrorInfo, config.WorkerConfig.LogFilePath, initLogFileError)
}

// OutputNoteInfoWithOutFormat 输出常规信息
func OutputNoteInfoWithOutFormat(content ...interface{}) {
	sendLogMsg(global.LogMarkNote, "%v", content...)
}

// OutputNoteInfo 输出带有格式的常规信息
func OutputNoteInfo(format string, content ...interface{}) {
	sendLogMsg(global.LogMarkNote, format, content...)
}

// OutputWarnInfo 输出警告信息
func OutputWarnInfo(format string, content ...interface{}) {
	sendLogMsg(global.LogMarkWarn, format, content...)
}

// OutputErrorInfo 输出错误信息
func OutputErrorInfo(format string, content ...interface{}) {
	sendLogMsg(global.LogMarkError, format, content...)
}

// OutputEmptyLine 输出空白行
func OutputEmptyLine() {
	LogChan <- []byte("\n")
}

func sendLogMsg(logMark, format string, content ...interface{}) {
	formatContent := fmt.Sprintf(format, content...)
	LogChan <- []byte(fmt.Sprintf("[%v] %v\n", logMark, formatContent))
}

func Run() {
CHAN:
	for {
		select {
		case msg := <-LogChan:
			LogFile.Write(msg)
		case <-LogControl:
			if len(LogChan) == 0 {
				break CHAN
			}
		}
	}
	Close()
}

func Close() {
	LogFile.Close()
	close(LogChan)
	close(LogControl)
}
