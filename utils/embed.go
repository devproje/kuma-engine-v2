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

// EmbedBuilder general embed builder
func EmbedBuilder(title, description string) *Embed {
	return &Embed{Title: title, Description: description}
}

// SetTitle set title
func (e *Embed) SetTitle(title string) *Embed {
	e.Title = title
	return e
}

// SetDescription set description
func (e *Embed) SetDescription(description string) *Embed {
	e.Description = description
	return e
}

// SetColor set color
func (e *Embed) SetColor(color int) *Embed {
	e.Color = color
	return e
}

// SetURL set url
func (e *Embed) SetURL(url string) *Embed {
	e.URL = url
	return e
}

// SetType set type
func (e *Embed) SetType(embedType discordgo.EmbedType) *Embed {
	e.Type = embedType
	return e
}

// SetTimestamp set timestamp
func (e *Embed) SetTimestamp(timestamp string) *Embed {
	e.Timestamp = timestamp
	return e
}

// SetImage set image
func (e *Embed) SetImage(url string, w, h int) *Embed {
	e.Image = &discordgo.MessageEmbedImage{URL: url, Width: w, Height: h}
	return e
}

// SetAuthor set author
func (e *Embed) SetAuthor(name, url, icon string) *Embed {
	e.Author = &discordgo.MessageEmbedAuthor{Name: name, URL: url, IconURL: icon}
	return e
}

// SetProvider set provider
func (e *Embed) SetProvider(name, url string) *Embed {
	e.Provider = &discordgo.MessageEmbedProvider{Name: name, URL: url}
	return e
}

// SetFooter set footer
func (e *Embed) SetFooter(text, icon string) *Embed {
	e.Footer = &discordgo.MessageEmbedFooter{Text: text, IconURL: icon}
	return e
}

// AddField add field
func (e *Embed) AddField(name, value string, inline bool) *Embed {
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})

	return e
}

// SetFields set fields
func (e *Embed) SetFields(fields []*discordgo.MessageEmbedField) *Embed {
	e.Fields = fields
	return e
}

// SetThumbnail set thumbnail
func (e *Embed) SetThumbnail(url string, w, h int) *Embed {
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: url, Width: w, Height: h}
	return e
}

// SetVideo set video
func (e *Embed) SetVideo(url string, w, h int) *Embed {
	e.Video = &discordgo.MessageEmbedVideo{
		URL:    url,
		Width:  w,
		Height: h,
	}

	return e
}

// Build building embed
func (e *Embed) Build() *discordgo.MessageEmbed {
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
