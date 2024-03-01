package emoji

import "fmt"

// Emoji: emoji data
type Emoji string

// SimpleBuilder: build simple emoji
func SimpleBuilder(name string) Emoji {
	return EmojiBuilder(name, "", false)
}

// EmojiBuilder: build emoji
func EmojiBuilder(name, id string, animate bool) Emoji {
	var str = "<"
	if animate {
		str += "a"
	}

	if id == "" {
		str = fmt.Sprintf(":%s:", name)
		return Emoji(str)
	}

	str += fmt.Sprintf(":%s:%s>", name, id)
	return Emoji(str)
}

func (e Emoji) Build() string {
	return string(e)
}
