package command

import "github.com/bwmarrin/discordgo"

func CommandHandler(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.Interaction.Type == discordgo.InteractionApplicationCommand {
		for _, i := range Commands {
			if event.ApplicationCommandData().Name == i.Data.Name {
				i.Execute(session, event)
			}
		}
	}
}
