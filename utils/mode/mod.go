package mode

import (
	"fmt"
	"os"

	"github.com/devproje/kuma-engine/log"
)

type EngineMode string

const (
	ReleaseMode EngineMode = "release"
	DebugMode   EngineMode = "debug"
)

var mode EngineMode = DebugMode

func init() {
	m := os.Getenv("ENGINE_MODE")
	modeEnv(m)
}

func modeEnv(t string) {
	if t == "" {
		return
	}
	m := EngineMode(t)
	switch m {
	case ReleaseMode:
		break
	case DebugMode:
		break
	default:
		log.Logger.Panicf(fmt.Sprintf("unknown mode: %s (avaliable mode: release, debug)", t))
	}

	SetMode(m)
}

func GetMode() EngineMode {
	return mode
}

func SetMode(t EngineMode) {
	mode = t
}
