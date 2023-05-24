package emoji

import "fmt"

type Emoji string

type Data struct {
	Name      string
	Id        string
	Animation bool
}

func (d Data) EmojiBuilder() Emoji {
	var str = "<"
	if d.Animation {
		str += "a"
	}

	if d.Id == "" {
		str = fmt.Sprintf(":%s:", d.Name)
		return Emoji(str)
	}

	str += fmt.Sprintf(":%s:%s>", d.Name, d.Id)
	return Emoji(str)
}
