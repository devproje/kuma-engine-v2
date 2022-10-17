package mode

import (
	"fmt"
	"os"

	"github.com/devproje/plog"
	"github.com/devproje/plog/level"
)

type EngineMode string

const (
	ReleaseMode EngineMode = "release"
	DebugMode   EngineMode = "debug"
)

var mode = DebugMode

func init() {
	m := os.Getenv("ENGINE_MODE")
	modeEnv(m)

	if mode == DebugMode {
		plog.SetLevel(level.Debug)
	}
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
		plog.Panicf(fmt.Sprintf("unknown mode: %s (avaliable mode: release, debug)", t))
	}

	SetMode(m)
}

func GetMode() EngineMode {
	return mode
}

func SetMode(t EngineMode) {
	mode = t
}
