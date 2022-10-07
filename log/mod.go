package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		FullTimestamp: true,
	},
	Hooks: make(logrus.LevelHooks),
	Level: logrus.InfoLevel,
}
