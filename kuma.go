package kuma

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/log"
	"github.com/devproje/kuma-engine/utils/mode"
)

const KUMA_ENGINE_VERSION = "v1.0.0"

var (
	act           []*discordgo.Activity
	delay         int  = 10
	engineStarted bool = false
	infoEnabled   bool = true
	infoEphemeral bool = false
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

	if mode.GetMode() == mode.DebugMode {
		log.Logger.Warnln(`Running in "debug" mode. Switch to "release" mode in production.
 - using env:  export ENGINE_MODE=release
 - using code: mode.SetMode(mode.ReleaseMode)`)
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
		for len(act) != 0 {
			for i := 0; i < len(act); i++ {
				k.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
					Status:     string(discordgo.StatusOnline),
					Activities: []*discordgo.Activity{act[i]},
				})

				time.Sleep(time.Second * time.Duration(delay))
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

// Events
func (k *KumaEngine) AddEvent(event interface{}) func() {
	return k.Session.AddHandler(event)
}

func (k *KumaEngine) AddEventOnce(event interface{}) func() {
	return k.Session.AddHandlerOnce(event)
}

// Activity
func (k *KumaEngine) AddAct(a *discordgo.Activity) {
	act = append(act, a)
}

func (k *KumaEngine) AddActs(a ...*discordgo.Activity) {
	act = append(act, a...)
}

func (k *KumaEngine) SetAct(a *discordgo.Activity) {
	act = []*discordgo.Activity{a}
}

func (k *KumaEngine) SetActs(a ...*discordgo.Activity) {
	k.InitActs()
	act = append(act, act...)
}

func (k *KumaEngine) InitActs() {
	act = []*discordgo.Activity{}
}

// Activity Options
func (k *KumaEngine) GetActDelay() int {
	return delay
}

func (k *KumaEngine) SetActDelay(second int) {
	delay = second
}

// KumaEngine Utils
func (k *KumaEngine) Version() string {
	return KUMA_ENGINE_VERSION
}

func (k *KumaEngine) CreateInturruptSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}
