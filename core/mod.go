package core

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/log"
)

const KUMA_ENGINE_VERSION = "v0.4.0"

var (
	act           []*discordgo.Activity
	delay         int  = 10
	engineStarted bool = false
	infoEnabled   bool = true
)

type KumaEngine struct {
	Token   string
	Color   int
	Session *discordgo.Session
}

// Create Engine
func (k KumaEngine) Create() (*KumaEngine, error) {
	log.Logger.Infof("KumaEngine %s\n", KUMA_ENGINE_VERSION)
	var err error
	k.Session, err = discordgo.New(fmt.Sprintf("Bot %s", k.Token))
	if err != nil {
		return nil, err
	}

	k.Session.AddHandler(command.CommandHandler)
	if infoEnabled {
		command.AddCommand(kumaInfo)
	}

	return &k, nil
}

func (k KumaEngine) CreateIntents(intent discordgo.Intent) (*KumaEngine, error) {
	engine, err := k.Create()
	if err != nil {
		return nil, err
	}

	k.Session.Identify.Intents = intent
	return engine, nil
}

// Engine Options
func (k *KumaEngine) Start() error {
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
		err := command.AddData(k.Session)
		if err != nil {
			log.Logger.Errorln(err)
		}
	}

	engineStarted = true
	return nil
}

func (k *KumaEngine) Stop() error {
	err := k.Session.Close()
	if err != nil {
		return err
	}

	return nil
}

// KumaEngine Version
func (k *KumaEngine) Version() string {
	return KUMA_ENGINE_VERSION
}
