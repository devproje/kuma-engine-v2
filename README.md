<br/>

![License](https://img.shields.io/github/license/devproje/kuma-engine-v2)
[![GoDoc](https://godoc.org/github.com/devproje/kuma-engine-v2?status.svg)](https://godoc.org/github.com/devproje/kuma-engine)
<img width="200" height="200" align="right" src="https://github.com/devproje/kuma-engine-v2/raw/master/assets/kuma-engine-logo.png" alt=""/>

# KumaEngine
Personal discordgo extend library

## How to use

### 1. Installation
```shell
go get -u github.com/devproje/kuma-engine-v2
```

### 2. Example code

```go
package main

import (
	"flag"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine"
	"github.com/devproje/kuma-engine/v2/command"
	"github.com/devproje/kuma-engine/v2/utils"
	"github.com/devproje/plog/log"
)

var (
	token = *flag.String("token", "", "Type Discord Token")
	test  = command.Executor{
		Data: &discordgo.ApplicationCommand{
			Name:        "test",
			Description: "Test Command",
		},
        Execute: func(ev *command.Event) error {
            err := ev.Reply("Test command")
			if err != nil {
				return err
            }
			
			return nil
        },
    }
)

func ready(session *discordgo.Session, ready *discordgo.Ready) {
	log.Infof("Logged in as %s", ready.User.String())
}

func main() {
	builder := kuma.EngineBuilder()
	builder.SetToken(token)
	
	builder.AddCommand(test)
	builder.AddEventOnceListener(ready)
	builder.SetKumaInfo(true)

	err := builder.Build()
	if err != nil {
		log.Fatalln(err)
	}
}
```
