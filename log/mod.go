package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	writers []io.Writer
	Logger  = &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.InfoLevel,
	}
)

func init() {
	writers = append(writers, os.Stdout)
	setLogger()
}

func setLogger() {
	Logger.SetOutput(io.MultiWriter(writers...))
}

func AddLoggingFile(name string) {
	f, err := os.OpenFile(fmt.Sprintf("%s.txt", name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		Logger.Fatalf("Failed to create '%s.txt' file\n%v", name, err)
		return
	}

	writers = append(writers, f)
	setLogger()
}
