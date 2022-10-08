package core

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/log"
)

var (
	act   []*discordgo.Activity
	delay int
)

type KumaEngine struct {
	Token   string
	Color   int
	Session *discordgo.Session
}

func (k KumaEngine) Create() (*KumaEngine, error) {
	log.Logger.Infof("KumaEngine %s\n", KUMA_ENGINE_VERSION)
	var err error
	k.Session, err = discordgo.New("Bot " + k.Token)
	if err != nil {
		return nil, err
	}

	k.Session.AddHandler(command.CommandHandler)
	command.RegisterCommand(KumaInfo)

	return &k, nil
}

func (k KumaEngine) CreateIntents() (*KumaEngine, error) {
	engine, err := k.Create()
	if err != nil {
		return nil, err
	}

	k.Session.Identify.Intents = discordgo.IntentsAll

	return engine, nil
}

func (k *KumaEngine) RegisterEvent(event interface{}) func() {
	return k.Session.AddHandler(event)
}

func (k *KumaEngine) RegisterEventOnce(event interface{}) func() {
	return k.Session.AddHandlerOnce(event)
}

func (k *KumaEngine) SetActivity(a *discordgo.Activity) {
	act = append(act, a)
}

func (k *KumaEngine) SetActivities(a ...*discordgo.Activity) {
	act = append(act, a...)
}

func (k *KumaEngine) Version() string {
	return KUMA_ENGINE_VERSION
}

func (k *KumaEngine) Build() error {
	err := k.Session.Open()
	if err != nil {
		return err
	}

	go func(delay int) {
		for {
			if len(act) != 0 {
				for i := 0; i < len(act); i++ {
					k.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
						Status:     string(discordgo.StatusOnline),
						Activities: []*discordgo.Activity{act[i]},
					})

					time.Sleep(time.Second * time.Duration(delay))
				}
			}
		}
	}(delay)

	if !command.IsCommandNull() {
		command.RegisterData(k.Session)
	}

	return nil
}

func (k *KumaEngine) Close() error {
	err := k.Session.Close()
	if err != nil {
		return err
	}

	return nil
}

func (k *KumaEngine) CreateInturruptSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}
