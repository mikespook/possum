package log

import (
	"github.com/rs/zerolog"
)

type Config struct {
	Level    string `mapstructure:"level,omitempty"`
	Filename string `mapstructure:"filename,omitempty"`
}

var (
	defaultConfig = &Config{
		Level: zerolog.LevelTraceValue,
	}
)
