package errorapp

import "errors"

//Ошибки иннициализации
var (
	ErrLoadEnvFile       = errors.New("error load env file")
	ErrParseDBConfig     = errors.New("failed wher parse dbconfig")
	ErrParseRedisConfig  = errors.New("failed wher parse redisconfig")
	ErrParseLoggerConfig = errors.New("failed wher parse loggerconfig")
	ErrServerConfig      = errors.New("failed wher parse serverconfig")
	ErrAPIConfig         = errors.New("failed wher parse apikey")
	ErrListed            = errors.New("failed to listed server")
	ErrServe             = errors.New("failed to serve server")
	ErrInitPostgres      = errors.New("failed init postgres")
	ErrInitRedis         = errors.New("failed init redis")
	ErrWebHookConfig     = errors.New("failed wher parse webhook config")
	ErrStatsConfig       = errors.New("failed wher stats config")
	ErrCacheConfig       = errors.New("failed when parse cache config")
)
