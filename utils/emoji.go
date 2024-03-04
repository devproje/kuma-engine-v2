package utils

import "fmt"

// Emoji emoji data
type Emoji string

// SimpleEmojiBuilder build simple emoji
func SimpleEmojiBuilder(name string) Emoji {
	return ExternalEmojiBuilder(name, "", false)
}

// ExternalEmojiBuilder build external server emoji
func ExternalEmojiBuilder(name, id string, animate bool) Emoji {
	var str = "<"
	if id == "" {
		str = fmt.Sprintf(":%s:", name)
		return Emoji(str)
	}

	if animate {
		str += "a"
	}

	str += fmt.Sprintf(":%s:%s>", name, id)
	return Emoji(str)
}

func (e Emoji) Build() string {
	return string(e)
}
