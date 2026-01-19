package config

import "time"

type Server interface {
	Address() string
	Port() string
	ReadTimeOut() time.Duration
	WriteTimeOut() time.Duration
	IdleTimeOut() time.Duration
	ShutDownTimeOut() time.Duration
}

type Postgres interface {
	Address() string
	MaxConns() int32
	MinConns() int32
	MaxConnLifeTime() time.Duration
	MaxConnIdleTime() time.Duration
	HealthCheckPeriod() time.Duration
}

type Log interface {
	Env() string
	Level() string
	Format() string
	FilePath() string
	Output() string
}

type Api interface {
	APIKey() string
}

type Redis interface {
	Address() string
	Password() string
	DB() int
	MaxRetries() int
	PoolSize() int
	MinIdleConns() int
	DialTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	PoolTimeout() time.Duration
	IdleTimeout() time.Duration
	MaxConnLifetime() time.Duration
}

type WebHook interface {
	Address() string
	Host() string
	Port() string
	Timeout() time.Duration
	MaxRetries() int
}

type Stats interface {
	TimeWindiwMinutes() int
}

type Cache interface {
	RefreshInterval() time.Duration
}
