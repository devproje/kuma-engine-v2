package command

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/v2/utils"
	"github.com/devproje/kuma-engine/v2/version"
)

const logo = "https://github.com/devproje/kuma-engine/raw/master/assets/kuma-engine-logo.png"

// KumaInfo framework info command
var KumaInfo = Executor{Data: data, Execute: execute}

var data = &discordgo.ApplicationCommand{
	Name:        "kumainfo",
	Description: "KumaEngine system information",
}

func execute(event *Event) error {
	embed := &utils.Embed{
		Title:       fmt.Sprintf("%s **KumaInfo**", utils.SimpleEmojiBuilder("dart")),
		Description: "KumaEngine system information",
		Color:       rand.Intn(0xFFFFFF),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    logo,
			Width:  512,
			Height: 512,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    event.Member.User.String(),
			IconURL: event.Member.User.AvatarURL("1024"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("%s **ENGINE VERSION**", utils.SimpleEmojiBuilder("electric_plug")),
				Value:  fmt.Sprintf("`%s`", version.KumaEngineVersion),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s **GO VERSION**", utils.SimpleEmojiBuilder("page_facing_up")),
				Value:  fmt.Sprintf("`%s`", runtime.Version()),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s **API LATENCY**", utils.SimpleEmojiBuilder("ping_pong")),
				Value:  fmt.Sprintf("`%dms`", event.Session.HeartbeatLatency().Milliseconds()),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s **OS**", utils.SimpleEmojiBuilder("desktop")),
				Value:  fmt.Sprintf("`%s/%s`", runtime.GOOS, runtime.GOARCH),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s **BOT SERVERS**", utils.SimpleEmojiBuilder("satellite")),
				Value:  fmt.Sprintf("`%d`", len(event.Session.State.Guilds)),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("%s **SYSTEM PID**", utils.SimpleEmojiBuilder("file_folder")),
				Value:  fmt.Sprintf("`%d`", os.Getpid()),
				Inline: true,
			},
		},
	}

	err := event.ReplyEmbedEphemeral(embed)
	if err != nil {
		return err
	}

	return nil
}
