package log

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	},
	Hooks: make(logrus.LevelHooks),
	Level: logrus.InfoLevel,
}
