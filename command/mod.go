package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/devproje/plog/log"
)

type CommandHandler struct {
	GuildId  string // Guild Command
	Commands []*CommandExecutor
}

type CommandExecutor struct {
	Data     *discordgo.ApplicationCommand
	Executor func(event *CommandEvent) error
}

type CommandEvent struct {
	Session           *discordgo.Session
	InteractionCreate *discordgo.InteractionCreate
	Member            *discordgo.Member
	User              *discordgo.User
}

func (c *CommandHandler) GetCommand(name string) *CommandExecutor {
	for _, command := range c.Commands {
		if command.Data.Name == name {
			return command
		}
	}
	return nil
}

func (c *CommandHandler) AddCommand(command CommandExecutor) {
	c.Commands = append(c.Commands, &command)
}

func (c *CommandHandler) DropCommand(name string) {
	for i, command := range c.Commands {
		if command.Data.Name == name {
			c.Commands = append(c.Commands[:i], c.Commands[i+1:]...)
		}
	}
}

func (c *CommandHandler) BuildHandler(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.Type != discordgo.InteractionApplicationCommand {
		return
	}

	command := c.GetCommand(event.ApplicationCommandData().Name)
	if command == nil {
		return
	}

	if c.GuildId != "" {
		log.Debugf("using \"%s\" guild's command by <@%s>: /%s", event.GuildID, event.Member.User.ID, event.ApplicationCommandData().Name)
	} else {
		log.Debugf("using global command by <@%s>: /%s", event.Member.User.ID, event.ApplicationCommandData().Name)
	}

	err := command.Executor(&CommandEvent{
		Session:           session,
		InteractionCreate: event,
		Member:            event.Member,
		User:              event.Member.User,
	})
	if err != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: 4,
			Data: &discordgo.InteractionResponseData{
				Content: "An error occurred while executing the command.",
				Flags:   1 << 6,
			},
		})
		log.Errorln(err)
	}
}

func (c *CommandHandler) RegisterCommand(session *discordgo.Session, guildId string) {
	for _, command := range c.Commands {
		log.Infof("Registering command: /%s\n", command.Data.Name)
		session.ApplicationCommandCreate("", c.GuildId, command.Data)
	}
}

func (c *CommandHandler) UnregisterCommand(session *discordgo.Session, guildId string) {
	cmds, _ := session.ApplicationCommands("", c.GuildId)
	for i, command := range cmds {
		log.Infof("Unregistering command (%d/%d)\n", i+1, len(cmds))
		session.ApplicationCommandDelete("", c.GuildId, command.ID)
	}
}
