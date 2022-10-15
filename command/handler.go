package command

import (
	"fmt"
	"github.com/devproje/kuma-engine/utils"
	"github.com/devproje/kuma-engine/utils/emoji"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/log"
	"github.com/devproje/kuma-engine/utils/mode"
)

func Handler(session *discordgo.Session, event *discordgo.InteractionCreate) {
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

				err := i.Execute(session, event)
				if err != nil {
					str := "An error occurred while executing the code"
					embed := utils.ErrorEmbed(event.Member.User, emoji.Warning, str)
					if mode.GetMode() == mode.DebugMode {
						embed.Description = fmt.Sprintf("%s\n**%s**", str, err.Error())
					}

					_ = session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
						Type: 5,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{},
						},
					})
					return
				}
			}
		}
	}
}
