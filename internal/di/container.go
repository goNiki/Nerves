package di

import (
	"fmt"

	"github.com/goNiki/Nerves/internal/config"
	incidentDataSource "github.com/goNiki/Nerves/internal/datasource/incident"
	locationDataSource "github.com/goNiki/Nerves/internal/datasource/location"
	webhookDataSource "github.com/goNiki/Nerves/internal/datasource/webhook"
	healthHandler "github.com/goNiki/Nerves/internal/handlers/health"
	incidentHandler "github.com/goNiki/Nerves/internal/handlers/incident"
	locationHandler "github.com/goNiki/Nerves/internal/handlers/location"
	statsHandler "github.com/goNiki/Nerves/internal/handlers/stats"
	"github.com/goNiki/Nerves/internal/infrastructure/database"
	"github.com/goNiki/Nerves/internal/infrastructure/logger"
	infraredis "github.com/goNiki/Nerves/internal/infrastructure/redis"
	cacheRepo "github.com/goNiki/Nerves/internal/repository/cache"
	incidentRepo "github.com/goNiki/Nerves/internal/repository/incident"
	locationRepo "github.com/goNiki/Nerves/internal/repository/location"
	queueRepo "github.com/goNiki/Nerves/internal/repository/queue"
	healthService "github.com/goNiki/Nerves/internal/service/health"
	incidentService "github.com/goNiki/Nerves/internal/service/incident"
	locationService "github.com/goNiki/Nerves/internal/service/location"
	statsService "github.com/goNiki/Nerves/internal/service/stats"
	webhookService "github.com/goNiki/Nerves/internal/service/webhook"
	cacheWorker "github.com/goNiki/Nerves/internal/worker/cache"
	webhookWorker "github.com/goNiki/Nerves/internal/worker/webhook"
)

type Container struct {
	Config *config.Config

	DB     *database.Postgres
	Redis  *infraredis.Client
	Logger *logger.Logger

	IncidentRepo      *incidentRepo.Repository
	LocationCheckRepo *locationRepo.Repository
	CacheRepo         *cacheRepo.Repository
	QueueRepo         *queueRepo.Repository

	IncidentDataSource *incidentDataSource.DataSource
	LocationDataSource *locationDataSource.DataSource
	WebhookDataSource  *webhookDataSource.DataSource

	IncidentService *incidentService.Service
	LocationService *locationService.Service
	StatsService    *statsService.Service
	HealthService   *healthService.Service
	WebhookService  *webhookService.Service

	IncidentHandler *incidentHandler.Handler
	LocationHandler *locationHandler.Handler
	StatsHandler    *statsHandler.Handler
	HealthHandler   *healthHandler.Handler

	WebhookWorker *webhookWorker.Worker
	CacheWorker   *cacheWorker.Worker
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{Config: cfg}

	loggerInstance := logger.InitLogger(cfg.Log)
	c.Logger = loggerInstance
	log := loggerInstance.Log

	db, err := database.InitDatabase(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	c.DB = db
	log.Info("postgres connected")

	redis, err := infraredis.New(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	c.Redis = redis
	log.Info("redis connected")

	c.IncidentRepo = incidentRepo.NewIncidentRepository(db)
	c.LocationCheckRepo = locationRepo.NewLocationRepository(db)
	c.CacheRepo = cacheRepo.NewIncidentCacheRepo(redis)
	c.QueueRepo = queueRepo.NewQueueRepository(redis)

	log.Info("repositories initialized")

	c.IncidentDataSource = incidentDataSource.NewDataSource(
		c.IncidentRepo,
		c.CacheRepo,
		log,
	)

	c.LocationDataSource = locationDataSource.New(
		c.CacheRepo,
		c.IncidentRepo,
		log,
	)

	c.WebhookDataSource = webhookDataSource.New(
		c.QueueRepo,
		log,
	)

	log.Info("datasources initialized")

	c.IncidentService = incidentService.NewIncidentService(c.IncidentDataSource, log)

	c.WebhookService = webhookService.New(c.WebhookDataSource, log)

	c.LocationService = locationService.New(
		c.LocationDataSource,
		c.LocationCheckRepo,
		c.WebhookService,
		log,
	)

	c.StatsService = statsService.New(c.LocationCheckRepo, cfg.Stats, log)

	c.HealthService = healthService.New(db, redis, log)

	log.Info("services initialized")

	c.IncidentHandler = incidentHandler.NewHandler(c.IncidentService)
	c.LocationHandler = locationHandler.NewHandler(c.LocationService)
	c.StatsHandler = statsHandler.NewHandler(c.StatsService)
	c.HealthHandler = healthHandler.NewHandler(c.HealthService)

	log.Info("handlers initialized")

	c.WebhookWorker = webhookWorker.New(
		c.WebhookDataSource,
		cfg.Webhook.Address(),
		cfg.Webhook.Timeout(),
		cfg.Webhook.MaxRetries(),
		log,
	)

	c.CacheWorker = cacheWorker.New(
		c.IncidentDataSource,
		cfg.Cache.RefreshInterval(),
		log,
	)

	log.Info("workers initialized")

	return c, nil
}

func (c *Container) Close() {
	if c.DB != nil {
		c.Logger.Log.Info("postgres connection closed")
	}

	if c.Redis != nil {
		c.Logger.Log.Info("redis connection closed")
	}
}
