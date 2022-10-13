package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/log"
	"github.com/devproje/kuma-engine/utils/mode"
)

func CommandHandler(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.Interaction.Type == discordgo.InteractionApplicationCommand {
		for _, i := range Commands {
			if event.ApplicationCommandData().Name == i.Data.Name {
				if mode.GetMode() == mode.DebugMode {
					cmd := event.ApplicationCommandData()
					var str = ""
					if len(cmd.Options) > 0 {
						for _, j := range cmd.Options {
							str += fmt.Sprintf("{%s: %v} ", j.Name, j.Value)
						}
					}

					log.Logger.Infof("%s used command: /%s %s\n", event.Member.User.String(), cmd.Name, str)
				}

				i.Execute(session, event)
			}
		}
	}
}
