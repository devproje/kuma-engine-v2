package kuma

import (
	"fmt"
	"io"
	"os"

	"github.com/devproje/plog/log"
)

func (k *Engine) init() {
	k.writers = append(k.writers, os.Stdout)
	k.setLogger()
}

func (k *Engine) setLogger() {
	log.SetOutput(io.MultiWriter(k.writers...))
}

func (k *Engine) AddLoggingFile(name string) {
	f, err := os.OpenFile(fmt.Sprintf("%s.txt", name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		log.Fatalf("Failed to create '%s.txt' file\n%v", name, err)
		return
	}

	k.writers = append(k.writers, f)
	k.setLogger()
}
