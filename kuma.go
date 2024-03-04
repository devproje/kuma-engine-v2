package kuma

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/v2/command"
	"github.com/devproje/kuma-engine/v2/mode"
	"github.com/devproje/kuma-engine/v2/version"
	"github.com/devproje/plog/level"
	"github.com/devproje/plog/log"
)

type Engine struct {
	Token   string
	Session *discordgo.Session

	writers []io.Writer
	mode    mode.EngineMode

	shards     []*discordgo.Session
	shardCount int
	intents    discordgo.Intent

	listeners       []*listener
	commandHandlers []*command.Handler

	innerHandler *command.Handler
	run          bool
	kumaInfo     bool
}

type listener struct {
	listener interface{}
	once     bool
}

func (k *Engine) init() {
	k.writers = append(k.writers, os.Stdout)
	k.setLogger()

	switch os.Getenv("ENGINE_MODE") {
	case "release":
		k.mode = mode.Release
	case "debug":
		k.mode = mode.Debug
	case "remove":
		k.mode = mode.Remove
	default:
		k.mode = mode.Debug
	}

	k.loadMode()
}

func (k *Engine) setLogger() {
	log.SetOutput(io.MultiWriter(k.writers...))
}

func (k *Engine) loadMode() {
	var msg string
	log.Traceln("loading engine mode...")

	switch k.mode {
	case mode.Release:
		log.SetLevel(level.Info)
	case mode.Debug:
		msg += "Kuma Engine running in \"debug\" mode. Switch to \"release\" mode in production."
		msg += "\n - using env:  export ENGINE_MODE=release"
		msg += "\n - using code: engine.SetMode(mode.Release)"
		log.SetLevel(level.Trace)
	case mode.Remove:
		msg += "Kuma Engine running in \"remove\" mode."
		msg += "\nAll commands and events are disabled."
		msg += "\nAnd will be removed bot command data."
		msg += "\n"
		msg += "Are you sure you want to use this mode? [y/N]:"
		log.SetLevel(level.Info)

		log.Warn(msg)

		var cmd string
		_, _ = fmt.Scanf("%s", &cmd)

		if cmd != "y" && cmd != "Y" {
			log.Infoln("Aborted")
			os.Exit(0)
		}
	}
}

func (k *Engine) removeCommandData() {
	go func() {
		_ = k.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "Remove command...",
					Type: discordgo.ActivityTypeListening,
				},
			},
			Status: string(discordgo.StatusDoNotDisturb),
		})
	}()

	log.Warnln("Removing all command data...")
	for _, h := range k.commandHandlers {
		h.UnregisterCommand(k.Session)
	}

	os.Exit(0)
}

// EngineBuilder Create a new Engine instance
func EngineBuilder() *Engine {
	engine := &Engine{
		run:          false,
		kumaInfo:     true,
		innerHandler: &command.Handler{},
		listeners:    make([]*listener, 0),
	}

	engine.init()
	log.Infof("Loading KumaEngine %s\n", version.KumaEngineVersion)

	return engine
}

// AddEventListener Add a new event listener
func (k *Engine) AddEventListener(ev interface{}) {
	l := listener{
		listener: ev,
		once:     false,
	}
	k.listeners = append(k.listeners, &l)
}

// AddEventOnceListener Add a new event listener once
func (k *Engine) AddEventOnceListener(ev interface{}) {
	l := listener{
		listener: ev,
		once:     true,
	}
	k.listeners = append(k.listeners, &l)
}

// AddCommandHandler Add a new command handler
func (k *Engine) AddCommandHandler(handler *command.Handler) {
	k.commandHandlers = append(k.commandHandlers, handler)
}

// RemoveCommandHandler Remove a command handler
func (k *Engine) RemoveCommandHandler(handler *command.Handler) {
	for i, h := range k.commandHandlers {
		if h == handler {
			k.commandHandlers = append(k.commandHandlers[:i], k.commandHandlers[i+1:]...)
		}
	}
}

// AddCommand Add a new command in inner handler
func (k *Engine) AddCommand(command command.Executor) {
	if k.run {
		return
	}

	k.innerHandler.AddCommand(command)
}

// DropCommand Remove a command in inner handler
func (k *Engine) DropCommand(name string) {
	if k.run {
		return
	}

	k.innerHandler.DropCommand(name)
}

// SetMode Set the engine mode
func (k *Engine) SetMode(m mode.EngineMode) {
	if k.run {
		return
	}

	k.mode = m
}

// SetToken Set the bot token
func (k *Engine) SetToken(token string) {
	if k.run {
		return
	}

	k.Token = token
}

// SetIntent Set the bot intents
func (k *Engine) SetIntent(intent discordgo.Intent) {
	if k.run {
		return
	}

	k.intents = intent
}

// IsKumaInfo Check if the kuma info command is enabled
func (k *Engine) IsKumaInfo() bool {
	return k.kumaInfo
}

// SetKumaInfo Set the kuma info command
func (k *Engine) SetKumaInfo(value bool) {
	if k.run {
		return
	}

	k.kumaInfo = value
}

func (k *Engine) GetShard(shardId int) (*discordgo.Session, error) {
	if len(k.shards) == 0 {
		return nil, fmt.Errorf("you are not using sharding")
	}

	if shardId >= len(k.shards) && shardId < 0 {
		return nil, fmt.Errorf("shard not found")
	}

	return k.shards[shardId], nil
}

// GetShardCount Get the bot shard count
func (k *Engine) GetShardCount() int {
	return k.shardCount
}

// SetShardCount Set the bot shard count
func (k *Engine) SetShardCount(shardCount int) {
	if k.run {
		return
	}

	k.shardCount = shardCount
}

// AddLoggingFile Add a new logging file
func (k *Engine) AddLoggingFile(name string) {
	f, err := os.OpenFile(fmt.Sprintf("%s.txt", name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		log.Fatalf("Failed to create '%s.txt' file\n%v", name, err)
		return
	}

	k.writers = append(k.writers, f)
	k.setLogger()
}

func (k *Engine) createSession(shardId int) error {
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", k.Token))
	if err != nil {
		log.Errorln("failed to create bot session, please check token and try again.")
		return err
	}

	bot.Identify.Intents = k.intents
	if k.shardCount > 0 {
		bot.ShardID = shardId
		bot.ShardCount = k.shardCount
	} else {
		k.Session = bot
	}

	for i, l := range k.listeners {
		if l.once {
			log.Debugf("Loading event once (%d/%d)\n", i+1, len(k.listeners))
			bot.AddHandlerOnce(l.listener)
			continue
		}

		log.Debugf("Loading event (%d/%d)\n", i+1, len(k.listeners))
		bot.AddHandler(l.listener)
	}

	err = bot.Open()
	if err != nil {
		return err
	}

	return nil
}

// Build Building the bot session
func (k *Engine) Build() error {
	if k.run {
		return nil
	}

	if k.kumaInfo {
		k.innerHandler.AddCommand(command.KumaInfo)
	}

	k.commandHandlers = append(k.commandHandlers, k.innerHandler)
	for _, h := range k.commandHandlers {
		k.listeners = append(k.listeners, &listener{
			listener: h.Build,
			once:     false,
		})
	}

	log.Traceln("creating bot session...")
	if k.mode == mode.Remove {
		err := k.createSession(0)
		if err != nil {
			return err
		}

		k.removeCommandData()
	}

	if k.shardCount > 0 {
		for i := 0; i < k.shardCount; i++ {
			log.Debugf("creating bot session for shard: %d\n", i)
			go func(shardId int) {
				err := k.createSession(shardId)
				if err != nil {
					log.Errorln(err)
				}

				log.Infof("Loading shard for id: #%d\n", shardId)
			}(i)
		}

	} else {
		err := k.createSession(0)
		if err != nil {
			return err
		}
	}

	k.run = true
	log.Traceln("registering commands...")
	go func() {
		for _, h := range k.commandHandlers {
			if k.shardCount > 0 {
				for _, s := range k.shards {
					h.RegisterCommand(s)
				}
			} else {
				h.RegisterCommand(k.Session)
			}
		}
	}()

	log.Traceln("bot session created successfully")
	return nil
}

// CreateInterruptSignal Create an interrupt signal
func (k *Engine) CreateInterruptSignal() {
	if !k.run {
		return
	}

	log.Infoln("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}

// Close Closing bot session
func (k *Engine) Close() error {
	if !k.run {
		return fmt.Errorf("bot session not running")
	}

	if k.shardCount > 0 {
		for _, s := range k.shards {
			err := s.Close()
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := k.Session.Close()
	if err != nil {
		return err
	}

	return nil
}
