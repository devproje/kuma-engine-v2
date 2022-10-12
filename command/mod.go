package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/log"
)

var Commands []Command

type Command struct {
	Data    *discordgo.ApplicationCommand
	Usage   string
	Execute func(session *discordgo.Session, event *discordgo.InteractionCreate)
}

func GetCommandData(name string) *discordgo.ApplicationCommand {
	for _, i := range Commands {
		if i.Data.Name == name {
			return i.Data
		}
	}

	return nil
}

func QueryCommandList() []*discordgo.ApplicationCommandOptionChoice {
	var list []*discordgo.ApplicationCommandOptionChoice
	for _, i := range Commands {
		list = append(list, &discordgo.ApplicationCommandOptionChoice{
			Name:  i.Data.Name,
			Value: i.Data.Name,
		})
	}

	return list
}

func AddCommand(cmd Command) {
	Commands = append(Commands, cmd)
}

func AddCommands(cmds ...Command) {
	Commands = append(Commands, cmds...)
}

func DropCommand(cmd Command) {
	for i, j := range Commands {
		if j.Data.Name == cmd.Data.Name {
			Commands = append(Commands[:i], Commands[i+1:]...)
		}
	}
}

func IsCommandNull() bool {
	return len(Commands) == 0
}

func AddData(session *discordgo.Session) error {
	for i, j := range Commands {
		log.Logger.Infof("Register command %s data (%d/%d)", j.Data.Name, i+1, len(Commands))
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", j.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func DropData(session *discordgo.Session) error {
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, i := range commands {
		log.Logger.Infof("Remove command %s data", i.Name)
		err = session.ApplicationCommandDelete(session.State.User.ID, "", i.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func DropDataManual(session *discordgo.Session, command Command) error {
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, i := range commands {
		if i.Name == command.Data.Name {
			err = session.ApplicationCommandDelete(session.State.User.ID, "", i.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Deprecated: Use AddCommand(cmd Command) instead
func RegisterCommand(cmd Command) {
	AddCommand(cmd)
}

// Deprecated: Use AddCommands(cmds ...Command) instead
func RegisterCommands(cmds ...Command) {
	AddCommands(cmds...)
}

// Deprecated: Use AddData(session *discordgo.Session) instead
func RegisterData(session *discordgo.Session) error {
	err := AddData(session)
	if err != nil {
		return err
	}

	return nil
}

// Deprecated: Use DropData(session *discordgo.Session) instead
func UnregisterData(session *discordgo.Session) error {
	err := DropData(session)
	if err != nil {
		return err
	}

	return nil
}
