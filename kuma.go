package kuma

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/mode"
	"github.com/devproje/kuma-engine/utils"
	"github.com/devproje/plog/level"
	"github.com/devproje/plog/log"
)

var innerMode mode.EngineMode

type KumaEngine struct {
	Token   string
	Session *discordgo.Session

	writers []io.Writer
	mode    mode.EngineMode

	shardID    int
	shardCount int
	intents    discordgo.Intent

	listeners       []*listener
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
	case "release":
		innerMode = mode.Release
	case "debug":
		innerMode = mode.Debug
	case "test":
		innerMode = mode.Test
	case "remove":
		innerMode = mode.Remove
	default:
		innerMode = mode.Debug
	}
}

func (k *KumaEngine) init() {
	k.writers = append(k.writers, os.Stdout)
	k.setLogger()
}

func (k *KumaEngine) setLogger() {
	log.SetOutput(io.MultiWriter(k.writers...))
}

func (k *KumaEngine) loadMode() {
	var msg string
	log.Traceln("loading engine mode...")

	switch k.mode {
	case mode.Release:
		log.SetLevel(level.Info)
	case mode.Debug:
		msg += "Kuma Engine running in \"debug\" mode. Switch to \"release\" mode in production."
		msg += "\n - using env:  export ENGINE_MODE=release"
		msg += "\n - using code: engine.SetMode(mode.Release)"
		log.SetLevel(level.Debug)
	case mode.Test:
		msg += "KUMA ENGINE IS RUNNING IN \"TEST\" MODE"
		msg += "\nDO NOT USE THIS MODE IN PRODUCTION"
		log.SetLevel(level.Trace)
	case mode.Remove:
		msg += "Kuma Engine running in \"remove\" mode."
		msg += "\nAll commands and events are disabled."
		msg += "\nAnd will be removed bot's command datas."
		msg += "\n"
		msg += "\nAre you sure you want to use this mode?"
		log.SetLevel(level.Info)
	}

	if msg != "" {
		log.Warnln(msg)
		var command string
		if k.mode == mode.Remove {
			fmt.Printf("Press 'y' to continue: ")
			fmt.Scanf("%s", command)

			if command != "y" {
				log.Infoln("Aborted")
				os.Exit(0)
			}
		}
	}
}

func (k *KumaEngine) loadModule() {
	if k.kumaInfo {
		log.Traceln("adding kuma info command...")
		k.innerHandler.AddCommand(command.KumaInfo)
	}

	log.Traceln("loading command handlers...")
	k.commandHandlers = append(k.commandHandlers, k.innerHandler)
	for _, h := range k.commandHandlers {
		k.listeners = append(k.listeners, &listener{
			listener: h.Build,
			once:     false,
		})
	}

	log.Traceln("loading event listeners...")
	for i, l := range k.listeners {
		if l.once {
			log.Debugf("Loading event once (%d/%d)\n", i+1, len(k.listeners))
			k.Session.AddHandlerOnce(l.listener)
			continue
		}

		log.Debugf("Loading event (%d/%d)\n", i+1, len(k.listeners))
		k.Session.AddHandler(l.listener)
	}
}

func (k *KumaEngine) removeData() {
	go func() {
		k.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "Remove command...",
					Type: discordgo.ActivityTypeListening,
				},
			},
			Status: "dnd",
		})
	}()

	log.Traceln("removing all command datas...")
	for _, h := range k.commandHandlers {
		h.UnregisterCommand(k.Session, h.GuildId)
	}

	os.Exit(0)
}

// EngineBuilder: Create a new KumaEngine instance
func EngineBuilder() *KumaEngine {
	return &KumaEngine{
		run:          false,
		kumaInfo:     true,
		innerHandler: &command.CommandHandler{},
		listeners:    make([]*listener, 0),
	}
}

// AddEventListener: Add a new event listener
func (k *KumaEngine) AddEventListener(ev interface{}) {
	l := listener{
		listener: ev,
		once:     false,
	}
	k.listeners = append(k.listeners, &l)
}

// AddEventOnceListener: Add a new event listener once
func (k *KumaEngine) AddEventOnceListener(ev interface{}) {
	l := listener{
		listener: ev,
		once:     true,
	}
	k.listeners = append(k.listeners, &l)
}

// AddCommandHandler: Add a new command handler
func (k *KumaEngine) AddCommandHandler(handler *command.CommandHandler) {
	k.commandHandlers = append(k.commandHandlers, handler)
}

// RemoveCommandHandler: Remove a command handler
func (k *KumaEngine) RemoveCommandHandler(handler *command.CommandHandler) {
	for i, h := range k.commandHandlers {
		if h == handler {
			k.commandHandlers = append(k.commandHandlers[:i], k.commandHandlers[i+1:]...)
		}
	}
}

// AddCommand: Add a new command in inner handler
func (k *KumaEngine) AddCommand(command command.CommandExecutor) {
	if k.run {
		return
	}

	k.innerHandler.AddCommand(command)
}

// DropCommand: Remove a command in inner handler
func (k *KumaEngine) DropCommand(name string) {
	k.innerHandler.DropCommand(name)
}

// SetMode: Set the engine mode
func (k *KumaEngine) SetMode(m mode.EngineMode) {
	if k.run {
		return
	}

	k.mode = m
}

// SetToken: Set the bot token
func (k *KumaEngine) SetToken(token string) {
	if k.run {
		return
	}

	k.Token = token
}

// SetIntent: Set the bot intents
func (k *KumaEngine) SetIntent(intent discordgo.Intent) {
	if k.run {
		return
	}

	k.intents = intent
}

// IsKumaInfo: Check if the kuma info command is enabled
func (k *KumaEngine) IsKumaInfo() bool {
	return k.kumaInfo
}

// SetKumaInfo: Set the kuma info command
func (k *KumaEngine) SetKumaInfo(value bool) {
	if k.run {
		return
	}

	k.kumaInfo = value
}

func (k *KumaEngine) GetShard() (int, int) {
	return k.shardID, k.shardCount
}

// SetShard: Set the bot shard
func (k *KumaEngine) SetShard(shardId, shardCount int) {
	if k.run {
		return
	}

	k.shardID = shardId
	k.shardCount = shardCount
}

// AddLoggingFile: Add a new logging file
func (k *KumaEngine) AddLoggingFile(name string) {
	f, err := os.OpenFile(fmt.Sprintf("%s.txt", name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		log.Fatalf("Failed to create '%s.txt' file\n%v", name, err)
		return
	}

	k.writers = append(k.writers, f)
	k.setLogger()
}

// Build: Build the bot session
func (k *KumaEngine) Build() error {
	k.init()
	k.loadMode()

	log.Infof("Loading KumaEngine %s\n", utils.KUMA_ENGINE_VERSION)

	log.Traceln("creating bot session...")
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", k.Token))
	if err != nil {
		log.Errorln("failed to create bot session, please check token and try again.")
		return err
	}

	k.Session = bot
	if k.mode == mode.Remove {
		log.Infoln("Kuma Engine is starting...")
		_ = bot.Open()
		k.run = true

		k.removeData()
	}

	k.loadModule()

	log.Traceln("setting bot intents...")
	bot.Identify.Intents = k.intents

	if (k.shardID != 0) && (k.shardCount != 0) {
		log.Traceln("sharding bot session...")
		bot.ShardID = k.shardID
		bot.ShardCount = k.shardCount
	}

	log.Traceln("opening bot session...")
	err = bot.Open()
	if err != nil {
		return err
	}

	k.run = true
	log.Traceln("bot session created successfully")

	log.Infoln("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	go func() {
		for _, h := range k.commandHandlers {
			h.RegisterCommand(bot, h.GuildId)
		}
	}()

	return nil
}
