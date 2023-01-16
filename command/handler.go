package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/utils/mode"
	"github.com/devproje/plog/log"
)

func debug(event *discordgo.InteractionCreate) {
	if mode.GetMode() != mode.DebugMode {
		return
	}
	cmd := event.ApplicationCommandData()
	var str = ""
	if len(cmd.Options) > 0 {
		for _, j := range cmd.Options {
			str += fmt.Sprintf("{%s: %v} ", j.Name, j.Value)
		}
	}

	log.Debugf("%s used command: /%s %s\n", event.Member.User.String(), cmd.Name, str)
}

func Handler(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.Interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	for _, i := range Commands {
		if event.ApplicationCommandData().Name == i.Data.Name {
			debug(event)
			err := i.Execute(session, event)
			if err != nil {
				str := "An error occurred while executing command"
				if mode.GetMode() == mode.DebugMode {
					str = fmt.Sprintf("%s\n%s", str, err.Error())
				}

				_ = session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
					Type: 4,
					Data: &discordgo.InteractionResponseData{
						Content: str,
					},
				})

				log.Errorln(err)
				return
			}
		}
	}
}
