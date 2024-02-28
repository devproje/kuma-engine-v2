package kuma

import (
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/mode"
	"github.com/devproje/kuma-engine/utils"
	"github.com/devproje/plog/level"
	"github.com/devproje/plog/log"
)

var innerMode *mode.EngineMode

type KumaEngine struct {
	Token string

	writers []io.Writer
	mode    mode.EngineMode

	shardID    int
	shardCount int
	intents    discordgo.Intent

	listeners       []listener
	commandHandlers []*command.CommandHandler

	innerHandler *command.CommandHandler
	run          bool
	kumaInfo     bool
}

type listener struct {
	listener interface{}
	once     bool
}

func init() {
	switch os.Getenv("ENGINE_MODE") {
	case "RELEASE_MODE":
		innerMode = &mode.Release
	case "DEBUG_MODE":
		innerMode = &mode.Debug
	case "TEST_MODE":
		innerMode = &mode.Test
	}
}

func EngineBuilder() *KumaEngine {
	return &KumaEngine{run: false}
}

func EngineBuilderWithShard(shardID, shards int) *KumaEngine {
	return &KumaEngine{
		shardID:    shardID,
		shardCount: shards,
		run:        false,
	}
}

func (k *KumaEngine) AddEventListener(ev interface{}) {
	l := listener{
		listener: ev,
		once:     false,
	}
	k.listeners = append(k.listeners, l)
}

func (k *KumaEngine) AddEventOnceListener(ev interface{}) {
	l := listener{
		listener: ev,
		once:     true,
	}
	k.listeners = append(k.listeners, l)
}

func (k *KumaEngine) AddCommandHandler(handler *command.CommandHandler) {
	k.commandHandlers = append(k.commandHandlers, handler)
}

func (k *KumaEngine) RemoveCommandHandler(handler *command.CommandHandler) {
	for i, h := range k.commandHandlers {
		if h == handler {
			k.commandHandlers = append(k.commandHandlers[:i], k.commandHandlers[i+1:]...)
		}
	}
}

func (k *KumaEngine) AddCommand(command command.CommandExecutor) {
	if k.run {
		return
	}

	k.innerHandler.AddCommand(command)
}

func (k *KumaEngine) DropCommand(name string) {
	k.innerHandler.DropCommand(name)
}

func (k *KumaEngine) SetMode(m mode.EngineMode) {
	if k.run {
		return
	}

	k.mode = m
}

func (k *KumaEngine) SetToken(token string) {
	if k.run {
		return
	}

	k.Token = token
}

func (k *KumaEngine) SetIntent(intent discordgo.Intent) {
	if k.run {
		return
	}

	k.intents = intent
}

func (k *KumaEngine) IsKumaInfo() bool {
	return k.kumaInfo
}

func (k *KumaEngine) SetKumaInfo(value bool) {
	if k.run {
		return
	}

	k.kumaInfo = value
}

func (k *KumaEngine) Build() (*discordgo.Session, error) {
	log.Infof("Loading KumaEngine %s\n", utils.KUMA_ENGINE_VERSION)
	k.settingMode()

	log.Traceln("creating bot session...")
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", k.Token))
	if err != nil {
		log.Errorln("failed to create bot session, please check token and try again.")
		return nil, err
	}

	log.Traceln("setting bot intents...")
	bot.Identify.Intents = k.intents

	if (k.shardID != 0) && (k.shardCount != 0) {
		log.Traceln("sharding bot session...")
		bot.ShardID = k.shardID
		bot.ShardCount = k.shardCount
	}

	if k.kumaInfo {
		log.Traceln("adding kuma info command...")
		k.innerHandler.AddCommand(command.KumaInfo)
	}

	log.Traceln("loading command handlers...")
	k.commandHandlers = append(k.commandHandlers, k.innerHandler)
	for _, h := range k.commandHandlers {
		k.listeners = append(k.listeners, listener{
			listener: h.BuildHandler,
			once:     false,
		})
	}

	log.Traceln("loading event listeners...")
	for i, l := range k.listeners {
		if l.once {
			log.Debugf("Loading event once (%d/%d)\n", i+1, len(k.listeners))
			bot.AddHandlerOnce(l.listener)
			continue
		}

		log.Debugf("Loading event (%d/%d)\n", i+1, len(k.listeners))
		bot.AddHandler(l.listener)
	}

	log.Traceln("opening bot session...")
	err = bot.Open()
	if err != nil {
		return nil, err
	}

	k.run = true

	log.Traceln("bot session created successfully")

	go func() {
		for _, h := range k.commandHandlers {
			h.RegisterCommand(bot, h.GuildId)
		}
	}()

	return bot, nil
}

func (k *KumaEngine) settingMode() {
	if innerMode != nil {
		k.mode = *innerMode
	}

	var msg string

	switch k.mode {
	case mode.Release:
		log.SetLevel(level.Info)
	case mode.Debug:
		msg += "Running in \"debug\" mode. Switch to \"release\" mode in production."
		msg += "\n - using env:  export ENGINE_MODE=release"
		msg += "\n - using code: engine.SetMode(mode.Release)"
		log.SetLevel(level.Debug)
	case mode.Test:
		msg += "KUMA ENGINE IS RUNNING IN \"TEST\" MODE"
		msg += "DO NOT USE THIS MODE IN PRODUCTION"
		log.SetLevel(level.Trace)
	}

	if msg != "" {
		log.Warnln(msg)
	}
}
