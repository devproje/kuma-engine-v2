package command

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/v2/utils"
	"github.com/devproje/kuma-engine/v2/utils/emoji"
	"github.com/devproje/kuma-engine/v2/version"
)

const logo = "https://github.com/devproje/kuma-engine/raw/master/assets/kuma-engine-logo.png"

// Kumainfo: framework info command
var KumaInfo = CommandExecutor{
	Data: &discordgo.ApplicationCommand{
		Name:        "kumainfo",
		Description: "KumaEngine system information",
	},
	Executor: func(event *CommandEvent) error {
		embed := utils.EmbedBuilder(fmt.Sprintf("%s **KumaInfo**", emoji.SimpleBuilder("dart")), "KumaEngine system information")
		embed.SetThumbnail(logo, 512, 512)
		embed.AddField(fmt.Sprintf("%s **ENGINE VERSION**", emoji.SimpleBuilder("electric_plug")), fmt.Sprintf("`%s`", version.KUMA_ENGINE_VERSION), true)
		embed.AddField(fmt.Sprintf("%s **GO VERSION**", emoji.SimpleBuilder("page_facing_up")), fmt.Sprintf("`%s`", runtime.Version()), true)
		embed.AddField(fmt.Sprintf("%s **API LATENCY**", emoji.SimpleBuilder("ping_pong")), fmt.Sprintf("`%dms`", event.Session.HeartbeatLatency().Milliseconds()), true)
		embed.AddField(fmt.Sprintf("%s **OS**", emoji.SimpleBuilder("desktop")), fmt.Sprintf("`%s/%s`", runtime.GOOS, runtime.GOARCH), true)
		embed.AddField(fmt.Sprintf("%s **BOT SERVERS**", emoji.SimpleBuilder("satellite")), fmt.Sprintf("`%d`", len(event.Session.State.Guilds)), true)
		embed.AddField(fmt.Sprintf("%s **SYSTEM PID**", emoji.SimpleBuilder("file_folder")), fmt.Sprintf("`%d`", os.Getpid()), true)
		embed.SetColor(rand.Intn(0xFFFFFF))
		embed.SetFooter(event.Member.User.String(), event.Member.User.AvatarURL("512x512"))

		err := event.ReplyEmbed(embed, true)
		if err != nil {
			return err
		}

		return nil
	},
}
