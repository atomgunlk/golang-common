package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func isLogfileOpen() bool {
	if logFile == nil || mw == nil {
		return false
	}
	return true
}
func logFileOpen() {
	//create your file with desired read/write permissions
	var err error
	logFilePath := ""
	if logFileNameDateSuffix {
		logFilePath = fmt.Sprintf("%s/%s_%s.log", logPath, logFileNamePrefix, time.Now().Format("20060102-15"))
	} else {
		logFilePath = fmt.Sprintf("%s/%s.log", logPath, logFileNamePrefix)
	}
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		log.Println(err)
		mw = io.MultiWriter(logStdOut)
	} else {
		//defer to close when you're done with it, not because you think it's idiomatic!
		// defer logFile.Close()
		if logToStdout {
			mw = io.MultiWriter(logStdOut, logFile)
		} else {
			mw = io.MultiWriter(logFile)
		}
	}

	// SetOutput Can be any io.Writer
	logger.SetOutput(mw)
	lastLogTime = time.Now()
	lastLogFileOpenTime = time.Now()
	fmt.Printf("Logger LOGFILE OPENED Path %s\r\n", logFilePath)
}

// logFileClose ... Clese log file fd
func logFileClose() {
	if isLogfileOpen() {
		// close
		if err := logFile.Close(); err != nil {
			log.Println(err)
		}
		logFile = nil
		mw = nil
		if logToStdout {
			logger.SetOutput(logStdOut)
		} else {
			logger.SetOutput(nil)
		}

	}
}

// Go routine for close file after no log for few second
func logFileCloseTimer() {
	forceCheckCloseFileTimer := time.Now()
	for {
		time.Sleep(time.Second * 10)
		if isLogfileOpen() {
			elapsed := int(time.Since(lastLogTime) / time.Second)
			if elapsed > logfileCloseTimeout { // check timeout sec
				fmt.Printf("Logger LOGFILE TIMEOUT %v sec, CLOSE FILE\r\n", elapsed)
				logFileClose()
			}
		}

		if isLogfileOpen() {
			forceCheckElapsed := int(time.Since(forceCheckCloseFileTimer) / time.Minute)
			if forceCheckElapsed >= 10 { // min
				forceCheckCloseFileTimer = time.Now()
				if lastLogFileOpenTime.YearDay() != time.Now().YearDay() {
					fmt.Printf("Logger LOGFILE Force close 10 min\r\n")
					logFileClose()
				}
			}
		}
	}
}
