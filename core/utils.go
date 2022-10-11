package core

import (
	"os"
	"os/signal"
	"syscall"
)

func (k *KumaEngine) CreateInturruptSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}
