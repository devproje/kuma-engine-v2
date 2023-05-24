package kuma

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/utils/mode"
	"github.com/devproje/plog/log"
)

const KUMA_ENGINE_VERSION = "v1.5.0"

type Engine struct {
	act         []*discordgo.Activity
	writers     []io.Writer
	delay       int
	started     bool
	infoEnabled bool

	Token   string
	Session *discordgo.Session
}

// Create default engine
func (k *Engine) Create() (*Engine, error) {
	var err error
	k.init()

	k.delay = 10
	k.started = false
	k.infoEnabled = true

	log.Infof("KumaEngine %s\n", KUMA_ENGINE_VERSION)
	k.Session, err = discordgo.New(fmt.Sprintf("Bot %s", k.Token))
	if err != nil {
		return nil, err
	}

	if mode.GetMode() == mode.DebugMode {
		log.Warnln(`Running in "debug" mode. Switch to "release" mode in production.
 - using env:  export ENGINE_MODE=release
 - using code: mode.SetMode(mode.ReleaseMode)`)
	}

	k.Session.AddHandler(command.Handler)
	if k.infoEnabled {
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
		for len(k.act) != 0 {
			for i := 0; i < len(k.act); i++ {
				_ = k.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
					Status:     string(discordgo.StatusOnline),
					Activities: []*discordgo.Activity{k.act[i]},
				})

				time.Sleep(time.Second * time.Duration(delay))
			}
		}
	}(k.delay)

	if !command.IsCommandNil() {
		err = command.AddData(k.Session)
		if err != nil {
			log.Errorln(err)
		}
	}

	k.started = true
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
	k.act = append(k.act, a)
}

// AddActs add many activities
func (k *Engine) AddActs(a ...*discordgo.Activity) {
	k.act = append(k.act, a...)
}

// SetAct set activity
func (k *Engine) SetAct(a *discordgo.Activity) {
	k.act = []*discordgo.Activity{a}
}

// SetActs set activities
func (k *Engine) SetActs(a ...*discordgo.Activity) {
	k.InitActs()
	k.act = append(k.act, a...)
}

// InitActs Initializing activity array
func (k *Engine) InitActs() {
	k.act = []*discordgo.Activity{}
}

// GetActDelay Get activity change time
func (k *Engine) GetActDelay() int {
	return k.delay
}

// SetActDelay Set activity change time
func (k *Engine) SetActDelay(second int) {
	k.delay = second
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
