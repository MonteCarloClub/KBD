package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func openLogFile(datadir string, filename string) *os.File {
	filePath := filename
	if !path.IsAbs(filename) {
		filePath = path.Join(datadir, filename)
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening log file '%s': %v", filename, err))
	}
	return file
}

func New(datadir string, logFile string, logLevel int) LogSystem {
	var writer io.Writer
	if logFile == "" {
		writer = os.Stdout
	} else {
		writer = openLogFile(datadir, logFile)
	}

	var sys LogSystem
	sys = NewStdLogSystem(writer, log.LstdFlags, LogLevel(logLevel))
	AddLogSystem(sys)

	return sys
}

func NewJSONsystem(datadir string, logFile string) LogSystem {
	var writer io.Writer
	if logFile == "-" {
		writer = os.Stdout
	} else {
		writer = openLogFile(datadir, logFile)
	}

	var sys LogSystem
	sys = NewJsonLogSystem(writer)
	AddLogSystem(sys)

	return sys
}
