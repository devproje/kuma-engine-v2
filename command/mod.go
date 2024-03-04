package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/v2/utils"
	"github.com/devproje/plog/log"
)

type Handler struct {
	GuildId  string // Guild Command
	Commands []*Executor
}

type Executor struct {
	Data    *discordgo.ApplicationCommand
	Execute func(event *Event) error
}

type Event struct {
	Session     *discordgo.Session
	Ev          *discordgo.InteractionCreate
	Member      *discordgo.Member
	User        *discordgo.User
	Interaction *discordgo.Interaction
}

// GetCommand get command by name
func (c *Handler) GetCommand(name string) *Executor {
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

// AddCommand add command to the command handler
func (c *Handler) AddCommand(command Executor) {
	log.Debugf("adding command: /%s", command.Data.Name)
	c.Commands = append(c.Commands, &command)
}

// DropCommand drop command from the command handler
func (c *Handler) DropCommand(name string) {
	for i, command := range c.Commands {
		if command.Data.Name == name {
			log.Debugf("dropping command: /%s", command.Data.Name)
			c.Commands = append(c.Commands[:i], c.Commands[i+1:]...)
		}
	}
}

// Build building command handler
func (c *Handler) Build(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.Type != discordgo.InteractionApplicationCommand {
		return
	}

	command := c.GetCommand(event.ApplicationCommandData().Name)
	if command == nil {
		return
	}

	var debug string
	if session.ShardCount > 0 {
		if c.GuildId != "" {
			debug = fmt.Sprintf(
				"using \"%s\" guild's command by <@%s> with #%d shard: /%s",
				event.GuildID, event.Member.User.ID,
				session.ShardID,
				event.ApplicationCommandData().Name,
			)
		} else {
			debug = fmt.Sprintf("using global command by <@%s> with #%d shard: /%s", event.Member.User.ID, session.ShardID, event.ApplicationCommandData().Name)
		}
	} else {
		if c.GuildId != "" {
			debug = fmt.Sprintf("using \"%s\" guild's command by <@%s>: /%s", event.GuildID, event.Member.User.ID, event.ApplicationCommandData().Name)
		} else {
			debug = fmt.Sprintf("using global command by <@%s>: /%s", event.Member.User.ID, event.ApplicationCommandData().Name)
		}
	}
	log.Debugln(debug)

	err := command.Execute(&Event{
		Session:     session,
		Ev:          event,
		Member:      event.Member,
		User:        event.Member.User,
		Interaction: event.Interaction,
	})
	if err != nil {
		err = session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: 4,
			Data: &discordgo.InteractionResponseData{
				Content: "An error occurred while executing the command.",
				Flags:   1 << 6,
			},
		})
		if err != nil {
			log.Errorln(err)
			return
		}
		log.Errorln(err)
	}
}

// RegisterCommand register command to the discord
func (c *Handler) RegisterCommand(session *discordgo.Session) {
	for _, command := range c.Commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, c.GuildId, command.Data)
		if err != nil {
			continue
		}
	}
}

// UnregisterCommand unregister command from the discord
func (c *Handler) UnregisterCommand(session *discordgo.Session) {
	cmds, _ := session.ApplicationCommands(session.State.User.ID, c.GuildId)
	for _, command := range cmds {
		err := session.ApplicationCommandDelete(session.State.User.ID, command.GuildID, command.ID)
		if err != nil {
			continue
		}
	}
}

// Reply send string message to the command
func (ev *Event) Reply(content string) error {
	return send(ev, content, false)
}

// ReplyEphemeral send string message to the command with ephemeral
func (ev *Event) ReplyEphemeral(content string) error {
	return send(ev, content, true)
}

// ReplyEmbed send embed message to the command
func (ev *Event) ReplyEmbed(embed *utils.Embed, ephemeral bool) error {
	return ev.ReplyEmbeds([]*utils.Embed{embed}, false)
}

// ReplyEmbedEphemeral send embed message to the command with ephemeral
func (ev *Event) ReplyEmbedEphemeral(embed *utils.Embed) error {
	return ev.ReplyEmbedsEphemeral([]*utils.Embed{embed})
}

// ReplyEmbeds send embed messages to the command
func (ev *Event) ReplyEmbeds(embeds []*utils.Embed, ephemeral bool) error {
	return sendEmbeds(ev, embeds, false)
}

// ReplyEmbedsEphemeral send embed messages to the command with ephemeral
func (ev *Event) ReplyEmbedsEphemeral(embeds []*utils.Embed) error {
	return sendEmbeds(ev, embeds, true)
}

func send(ev *Event, content string, ephemeral bool) error {
	data := &discordgo.InteractionResponseData{
		Content: content,
	}

	if ephemeral {
		data.Flags = 1 << 6
	}

	return ev.Session.InteractionRespond(ev.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: data,
	})
}

func sendEmbeds(ev *Event, embeds []*utils.Embed, ephemeral bool) error {
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

	return ev.Session.InteractionRespond(ev.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: data,
	})
}
