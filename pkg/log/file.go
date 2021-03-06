package log

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LogrusFileHook hook for logrus to write log to file
type LogrusFileHook struct {
	file      *os.File
	flag      int
	chmod     os.FileMode
	formatter logrus.Formatter
}

// NewLogrusFileHook returns new file hook object for logrus
func NewLogrusFileHook(file string, flag int, chmod os.FileMode) (*LogrusFileHook, error) {
	plainFormatter := getFormatter(true)
	logFile, err := os.OpenFile(file, flag, chmod)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err)
		return nil, err
	}

	return &LogrusFileHook{logFile, flag, chmod, plainFormatter}, err
}

// Fire func used by logrus to write the log into a log file
func (hook *LogrusFileHook) Fire(entry *logrus.Entry) error {
	plainFormat, err := hook.formatter.Format(entry)
	if err != nil {
		return errors.WithStack(err)
	}
	line := string(plainFormat)
	_, err = hook.file.WriteString(line)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to write file on filehook (entry.String): %v", err)
		return err
	}

	return nil
}

// Levels defines in which log levels the file hooks works
func (hook *LogrusFileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
