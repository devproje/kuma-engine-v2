package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const ErrorColor = 0xDD0000

type Embed struct {
	Title       string
	Description string
	Image       *discordgo.MessageEmbedImage
	Thumbnail   *discordgo.MessageEmbedThumbnail
	Fields      []*discordgo.MessageEmbedField
	Color       int
}

func writeFooter(executor *discordgo.User) *discordgo.MessageEmbedFooter {
	return &discordgo.MessageEmbedFooter{
		Text:    executor.String(),
		IconURL: executor.AvatarURL("512x512"),
	}
}

func (e Embed) Build(executor *discordgo.User) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       e.Title,
		Description: e.Description,
		Image:       e.Image,
		Fields:      e.Fields,
		Color:       e.Color,
		Footer:      writeFooter(executor),
	}
}

func ErrorEmbed(executor *discordgo.User, emoji string) *discordgo.MessageEmbed {
	return Embed{
		Title:       fmt.Sprintf("%s **Error!**", emoji),
		Description: "명령어를 실행하는 도중에 오류가 발생 했어요",
		Color:       ErrorColor,
	}.Build(executor)
}

func CustomErrorEmbed(executor *discordgo.User, emoji string, message string) *discordgo.MessageEmbed {
	return Embed{
		Title:       fmt.Sprintf("%s **Error!**", emoji),
		Description: message,
		Color:       ErrorColor,
	}.Build(executor)
}

func PermissionDenied(executor *discordgo.User, emoji string) *discordgo.MessageEmbed {
	return CustomErrorEmbed(executor, emoji, "이 명령어를 수행할 권한이 없어요")
}
