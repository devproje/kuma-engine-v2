package utils

import (
	"fmt"
	"github.com/devproje/kuma-engine/utils/emoji"

	"github.com/bwmarrin/discordgo"
)

const ERROR_COLOR = 0xDD0000

type Embed struct {
	Title       string
	Description string
	Fields      []*discordgo.MessageEmbedField
	Footer      *discordgo.MessageEmbedFooter
	Color       int
	URL         string
	Type        discordgo.EmbedType
	Timestamp   string
	Image       *discordgo.MessageEmbedImage
	Thumbnail   *discordgo.MessageEmbedThumbnail
	Video       *discordgo.MessageEmbedVideo
	Provider    *discordgo.MessageEmbedProvider
	Author      *discordgo.MessageEmbedAuthor
}

func (e Embed) Build() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		URL:         e.URL,
		Type:        e.Type,
		Title:       e.Title,
		Description: e.Description,
		Timestamp:   e.Timestamp,
		Color:       e.Color,
		Footer:      e.Footer,
		Image:       e.Image,
		Thumbnail:   e.Thumbnail,
		Video:       e.Video,
		Provider:    e.Provider,
		Author:      e.Author,
		Fields:      e.Fields,
	}
}

func ErrorEmbed(executor *discordgo.User, emoji emoji.Emoji, message string) *discordgo.MessageEmbed {
	return Embed{
		Title:       fmt.Sprintf("%s **Error!**", emoji),
		Description: message,
		Color:       ERROR_COLOR,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    executor.String(),
			IconURL: executor.AvatarURL("512x512"),
		},
	}.Build()
}
