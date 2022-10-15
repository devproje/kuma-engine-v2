<br/>

![License](https://img.shields.io/github/license/devproje/kuma-engine)
[![GoDoc](https://godoc.org/github.com/devproje/kuma-engine?status.svg)](https://godoc.org/github.com/devproje/kuma-engine)
<img width="200" height="200" align="right" src="https://github.com/devproje/kuma-engine/raw/master/assets/kuma-engine-logo.png" alt=""/>

# KumaEngine
Personal discordgo extend library

## How to use

### 1. Installation
```shell
go get -u github.com/devproje/kuma-engine
```

### 2. Example code

```go
package main

import (
	"flag"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/devproje/kuma-engine"
	"github.com/devproje/kuma-engine/command"
	"github.com/devproje/kuma-engine/log"
	"github.com/devproje/kuma-engine/utils"
)

var (
	token = flag.String("token", "", "Type Discord Token")
	e     *kuma.Engine
)

func ready(session *discordgo.Session, ready *discordgo.Ready) {
	log.Logger.Infof("Logged in as %s", ready.User.String())
}

func main() {
	flag.Parse()
	var err error
	engine := &kuma.Engine{
		Token: token,
		Color: 0xFF0000,
	}
	engine, err = engine.Create()
	if err != nil {
		log.Logger.Fatalln(err)
	}
	
	engine.AddEventOnce(ready)

	err = engine.Start()
	if err != nil {
		log.Logger.Fatalln(err)
	}

	e = engine
	log.Logger.Infoln("Bot is now running. Press CTRL-C to exit.")
	engine.CreateInterruptSignal()

	err = engine.Stop()
	if err != nil {
		log.Logger.Fatalln(err)
	}
}
```
