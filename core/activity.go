package core

import "github.com/bwmarrin/discordgo"

// Activity
func (k *KumaEngine) SetAct(a *discordgo.Activity) {
	act = append(act, a)
}

func (k *KumaEngine) SetActs(a ...*discordgo.Activity) {
	act = append(act, a...)
}

// Activity Options
func (k *KumaEngine) GetActDelay() int {
	return delay
}

func (k *KumaEngine) SetActDelay(second int) {
	delay = second
}
