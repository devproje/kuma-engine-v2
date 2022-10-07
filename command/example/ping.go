package example

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine/core"
	"github.com/devproje/kuma-engine/emoji"
	"github.com/devproje/kuma-engine/utils"
)

func Ping(engine *core.KumaEngine, session *discordgo.Session, event *discordgo.InteractionCreate) {
	ping := session.HeartbeatLatency()
	embed := utils.Embed{
		Title:       fmt.Sprintf("%s **Measuring...**", emoji.HourGlassFlowingSand),
		Description: "Just hold on sec...",
		Color:       0x0D0D0D,
	}.Build(event.Member.User)
	before := time.Now()
	_ = session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	after := time.Since(before)
	embed = utils.Embed{
		Title: fmt.Sprintf("%s **Pong!**", emoji.PingPong),
		Description: fmt.Sprintf(
			"**BOT**: %s**ms**\n**API**: %s**ms**",
			strconv.FormatInt(after.Milliseconds(), 10),
			strconv.FormatInt(ping.Milliseconds(), 10)),
		Color: engine.Color,
	}.Build(event.Member.User)
	_, _ = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}
