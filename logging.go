package kuma

import (
	"fmt"
	"io"
	"os"

	"github.com/devproje/plog"
)

var writers []io.Writer

func init() {
	writers = append(writers, os.Stdout)
	setLogger()
}

func setLogger() {
	plog.SetOutput(io.MultiWriter(writers...))
}

func AddLoggingFile(name string) {
	f, err := os.OpenFile(fmt.Sprintf("%s.txt", name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		plog.Fatalf("Failed to create '%s.txt' file\n%v", name, err)
		return
	}

	writers = append(writers, f)
	setLogger()
}
