package kuma

import (
	"fmt"
	"io"
	"os"

	"github.com/devproje/plog/log"
)

func (k *KumaEngine) init() {
	k.writers = append(k.writers, os.Stdout)
	k.setLogger()
}

func (k *KumaEngine) setLogger() {
	log.SetOutput(io.MultiWriter(k.writers...))
}

func (k *KumaEngine) AddLoggingFile(name string) {
	f, err := os.OpenFile(fmt.Sprintf("%s.txt", name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		log.Fatalf("Failed to create '%s.txt' file\n%v", name, err)
		return
	}

	k.writers = append(k.writers, f)
	k.setLogger()
}
