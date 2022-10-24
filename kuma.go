package kuma

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/utils/mode"
	"github.com/devproje/plog"
)

const KUMA_ENGINE_VERSION = "v1.3.2"

var (
	act           []*discordgo.Activity
	delay         = 10
	engineStarted = false
	infoEnabled   = true
	infoEphemeral = false
)

type Engine struct {
	Token   string
	Color   int
	Session *discordgo.Session
}

// Create default engine
func (k *Engine) Create() (*Engine, error) {
	plog.Infof("KumaEngine %s\n", KUMA_ENGINE_VERSION)
	var err error
	k.Session, err = discordgo.New(fmt.Sprintf("Bot %s", k.Token))
	if err != nil {
		return nil, err
	}

	if mode.GetMode() == mode.DebugMode {
		plog.Warnln(`Running in "debug" mode. Switch to "release" mode in production.
 - using env:  export ENGINE_MODE=release
 - using code: mode.SetMode(mode.ReleaseMode)`)
	}

	k.Session.AddHandler(command.Handler)
	if infoEnabled {
		command.AddCommand(kumaInfo)
	}

	return k, nil
}

// CreateIntents create engine with discord intents
func (k *Engine) CreateIntents(intent discordgo.Intent) (*Engine, error) {
	engine, err := k.Create()
	if err != nil {
		return nil, err
	}

	k.Session.Identify.Intents = intent
	return engine, nil
}

// Start starting engine
func (k *Engine) Start() error {
	err := k.Session.Open()
	if err != nil {
		return err
	}

	go func(delay int) {
		for len(act) != 0 {
			for i := 0; i < len(act); i++ {
				_ = k.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
					Status:     string(discordgo.StatusOnline),
					Activities: []*discordgo.Activity{act[i]},
				})

				time.Sleep(time.Second * time.Duration(delay))
			}
		}
	}(delay)

	if !command.IsCommandNil() {
		err = command.AddData(k.Session)
		if err != nil {
			plog.Errorln(err)
		}
	}

	engineStarted = true
	return nil
}

// Stop stopped engine
func (k *Engine) Stop() error {
	err := k.Session.Close()
	if err != nil {
		return err
	}

	return nil
}

// AddEvent Add discord event handler
func (k *Engine) AddEvent(event interface{}) func() {
	return k.Session.AddHandler(event)
}

// AddEventOnce Add discord event handler once
func (k *Engine) AddEventOnce(event interface{}) func() {
	return k.Session.AddHandlerOnce(event)
}

// AddAct add one activity
func (k *Engine) AddAct(a *discordgo.Activity) {
	act = append(act, a)
}

// AddActs add many activities
func (k *Engine) AddActs(a ...*discordgo.Activity) {
	act = append(act, a...)
}

// SetAct set activity
func (k *Engine) SetAct(a *discordgo.Activity) {
	act = []*discordgo.Activity{a}
}

// SetActs set activities
func (k *Engine) SetActs(a ...*discordgo.Activity) {
	k.InitActs()
	act = append(act, a...)
}

// InitActs Initializing activity array
func (k *Engine) InitActs() {
	act = []*discordgo.Activity{}
}

// GetActDelay Get activity change time
func (k *Engine) GetActDelay() int {
	return delay
}

// SetActDelay Set activity change time
func (k *Engine) SetActDelay(second int) {
	delay = second
}

// Version Checking engine version
func (k *Engine) Version() string {
	return KUMA_ENGINE_VERSION
}

// CreateInterruptSignal Creating Ctrl+C interrupt signal
func (k *Engine) CreateInterruptSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}
