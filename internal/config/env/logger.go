package env

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type LoggerEnvConfig struct {
	Env      string `env:"LOG_ENV,required"`
	Level    string `env:"LOG_LEVEL,required"`
	Format   string `env:"LOG_FORMAT" envDefault:"json"`
	FilePath string `env:"LOG_FILE_PATH" envDefault:"/var/log/app.log"`
	Output   string `env:"LOG_OUTPUT,required"`
}

type LoggerConfig struct {
	raw LoggerEnvConfig
}

func NewLoggerConfig() (*LoggerConfig, error) {
	var raw LoggerEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrParseLoggerConfig, err)
	}

	return &LoggerConfig{
		raw: raw,
	}, nil

}

func (cfg *LoggerConfig) Env() string {
	return cfg.raw.Env
}

func (cfg *LoggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *LoggerConfig) Format() string {
	return cfg.raw.Format
}

func (cfg *LoggerConfig) FilePath() string {
	return cfg.raw.FilePath
}

func (cfg *LoggerConfig) Output() string {
	return cfg.raw.Output
}
