package utils

import (
	"github.com/bwmarrin/discordgo"
)

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
