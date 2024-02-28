package command

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/utils"
	"github.com/devproje/kuma-engine/utils/emoji"
)

const logo = "https://github.com/devproje/kuma-engine/raw/master/assets/kuma-engine-logo.png"

var KumaInfo = CommandExecutor{
	Data: &discordgo.ApplicationCommand{
		Name:        "kumainfo",
		Description: "KumaEngine system information",
	},
	Executor: func(event *CommandEvent) error {
		embed := utils.Embed{
			Title:       fmt.Sprintf("%s **KumaInfo**", emoji.Dart),
			Description: "KumaEngine system information",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL:    logo,
				Width:  512,
				Height: 512,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   fmt.Sprintf("%s **ENGINE VERSION**", emoji.ElectricPlug),
					Value:  fmt.Sprintf("`%s`", utils.KUMA_ENGINE_VERSION),
					Inline: true,
				},
				{
					Name:   fmt.Sprintf("%s **GO VERSION**", emoji.PageFacingUp),
					Value:  fmt.Sprintf("`%s`", runtime.Version()),
					Inline: true,
				},
				{
					Name:   fmt.Sprintf("%s **API LATENCY**", emoji.PingPong),
					Value:  fmt.Sprintf("`%dms`", event.Session.HeartbeatLatency().Milliseconds()),
					Inline: true,
				},
				{
					Name:   fmt.Sprintf("%s **OS**", emoji.Desktop),
					Value:  fmt.Sprintf("`%s/%s`", runtime.GOOS, runtime.GOARCH),
					Inline: true,
				},
				{
					Name:   fmt.Sprintf("%s **BOT SERVERS**", emoji.Satellite),
					Value:  fmt.Sprintf("`%d`", len(event.Session.State.Guilds)),
					Inline: true,
				},
				{
					Name:   fmt.Sprintf("%s **SYSTEM PID**", emoji.FileFolder),
					Value:  fmt.Sprintf("`%d`", os.Getpid()),
					Inline: true,
				},
			},
			Color: rand.Intn(0xFFFFFF),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    event.Member.User.String(),
				IconURL: event.Member.User.AvatarURL("512x512"),
			},
		}.Build()

		data := &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		}

		err := event.Session.InteractionRespond(event.InteractionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: data,
		})
		if err != nil {
			return err
		}

		return nil
	},
}
