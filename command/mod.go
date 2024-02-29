package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/v2/utils"
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

// GetCommand: get command by name
func (c *CommandHandler) GetCommand(name string) *CommandExecutor {
	if c.Commands == nil {
		return nil
	}

	for _, command := range c.Commands {
		if command.Data.Name == name {
			return command
		}
	}
	return nil
}

// AddCommand: add command to the command handler
func (c *CommandHandler) AddCommand(command CommandExecutor) {
	c.Commands = append(c.Commands, &command)
}

// DropCommand: drop command from the command handler
func (c *CommandHandler) DropCommand(name string) {
	for i, command := range c.Commands {
		if command.Data.Name == name {
			c.Commands = append(c.Commands[:i], c.Commands[i+1:]...)
		}
	}
}

// BuildHandler: build command handler
func (c *CommandHandler) Build(session *discordgo.Session, event *discordgo.InteractionCreate) {
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

// RegisterCommand: register command to the discord
func (c *CommandHandler) RegisterCommand(session *discordgo.Session, guildId string) {
	for _, command := range c.Commands {
		log.Infof("Registering command: /%s\n", command.Data.Name)
		_, err := session.ApplicationCommandCreate("", c.GuildId, command.Data)
		if err != nil {
			log.Errorf("An error occurred while registering the command: /%s\n", command.Data.Name)
			continue
		}
	}
}

// UnregisterCommand: unregister command from the discord
func (c *CommandHandler) UnregisterCommand(session *discordgo.Session, guildId string) {
	cmds, _ := session.ApplicationCommands("", c.GuildId)
	for i, command := range cmds {
		var msg = "Unregistering global command"
		if command.GuildID != "" {
			msg = fmt.Sprintf("Unregistering %s guild's command", command.GuildID)
		}

		log.Infof("%s (%d/%d)\n", msg, i+1, len(cmds))
		err := session.ApplicationCommandDelete("", c.GuildId, command.ID)
		if err != nil {
			log.Errorf("An error occurred while unregistering the command: /%s\n", command.Name)
			continue
		}
	}
}

// Reply: send string message to the command
func (ev *CommandEvent) Reply(content string, ephemeral bool) error {
	data := &discordgo.InteractionResponseData{
		Content: content,
	}

	if ephemeral {
		data.Flags = 1 << 6
	}

	return ev.Session.InteractionRespond(ev.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: data,
	})
}

// ReplyEmbed: send embed message to the command
func (ev *CommandEvent) ReplyEmbed(embed *utils.Embed, ephemeral bool) error {
	return ev.ReplyEmbeds([]*utils.Embed{embed}, ephemeral)
}

// ReplyEmbeds: send embed messages to the command
func (ev *CommandEvent) ReplyEmbeds(embeds []*utils.Embed, ephemeral bool) error {
	var me []*discordgo.MessageEmbed
	for _, embed := range embeds {
		me = append(me, embed.Build())
	}

	data := &discordgo.InteractionResponseData{
		Embeds: me,
	}

	if ephemeral {
		data.Flags = 1 << 6
	}

	return ev.Session.InteractionRespond(ev.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: data,
	})
}
