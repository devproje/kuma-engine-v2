package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/devproje/plog/log"
)

var Commands []Command

type Command struct {
	Data    *discordgo.ApplicationCommand
	Usage   string
	Execute func(session *discordgo.Session, event *discordgo.InteractionCreate) error
}

// GetCommandData getting target command data
func GetCommandData(name string) *discordgo.ApplicationCommand {
	for _, i := range Commands {
		if i.Data.Name == name {
			return i.Data
		}
	}

	return nil
}

// QueryCommandList getting all command list
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

// AddCommand add target application command handler
func AddCommand(cmd Command) {
	Commands = append(Commands, cmd)
}

// AddCommands add many application command handlers
func AddCommands(cmds ...Command) {
	Commands = append(Commands, cmds...)
}

// DropCommand delete target application command handler
func DropCommand(cmd Command) {
	for i, j := range Commands {
		if j.Data.Name == cmd.Data.Name {
			Commands = append(Commands[:i], Commands[i+1:]...)
		}
	}
}

// IsCommandNil checking Commands array is nil
func IsCommandNil() bool {
	return len(Commands) == 0
}

// AddData add all application commands data
func AddData(session *discordgo.Session) error {
	for i, j := range Commands {
		log.Infof("Register command %s data (%d/%d)", j.Data.Name, i+1, len(Commands))
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", j.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// DropData delete all application commands data
func DropData(session *discordgo.Session) error {
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, i := range commands {
		log.Infof("Remove command %s data", i.Name)
		err = session.ApplicationCommandDelete(session.State.User.ID, "", i.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// DropDataManual delete target application command data
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
