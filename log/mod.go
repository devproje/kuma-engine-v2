package log

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger main logging system
var Logger = &logrus.Logger{
	Out: os.Stdout,
	Formatter: &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	},
	Hooks: make(logrus.LevelHooks),
	Level: logrus.InfoLevel,
}
