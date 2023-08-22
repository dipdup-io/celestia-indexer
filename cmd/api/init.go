package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})
}

func initConfig() (*Config, error) {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("command line execute")
		return nil, err
	}

	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		log.Panic().Err(err).Msg("config command line arg is required")
		return nil, err
	}

	var cfg Config
	if err := config.Parse(*configPath, &cfg); err != nil {
		log.Panic().Err(err).Msg("parsing config file")
		return nil, err
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelInfoValue
	}

	return &cfg, nil
}

func initLogger(level string) error {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Panic().Err(err).Msg("parsing log level")
		return err
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.Logger.With().Caller().Logger()

	return nil
}

func initEcho(cfg ApiConfig) *echo.Echo {
	e := echo.New()
	e.Validator = handler.NewCelestiaApiValidator()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogMethod:    true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("uri", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Str("method", v.Method).
				Str("user-agent", v.UserAgent).
				Str("ip", c.RealIP()).
				Msg("request")

			return nil
		},
	}))
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CSRF())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Pre(middleware.RemoveTrailingSlash())

	timeout := 30 * time.Second
	if cfg.RequestTimeout > 0 {
		timeout = time.Duration(cfg.RequestTimeout) * time.Second
	}
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		Timeout: timeout,
	}))

	if cfg.Prometheus {
		e.Use(echoprometheus.NewMiddleware("celestia_api"))
	}
	if cfg.RateLimit > 0 {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit))))
	}
	return e
}

func initDatabase(cfg config.Database) postgres.Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := postgres.Create(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return db
}

func initHandlers(e *echo.Echo, cfg ApiConfig, db postgres.Storage) {
	v1 := e.Group("v1")

	stateHandlers := handler.NewStateHandler(db.State)
	v1.GET("/head", stateHandlers.Head)

	addessHandlers := handler.NewAddressHandler(db.Address)
	addressGroup := v1.Group("/address")
	{
		addressGroup.GET("", addessHandlers.List)
		addressGroup.GET("/:hash", addessHandlers.Get)
	}

	blockHandlers := handler.NewBlockHandler(db.Blocks, db.Event)
	blockGroup := v1.Group("/block")
	{
		blockGroup.GET("", blockHandlers.List)
		blockGroup.GET("/:height", blockHandlers.Get)
		blockGroup.GET("/:height/events", blockHandlers.GetEvents)
	}

	txHandlers := handler.NewTxHandler(db.Tx, db.Event, db.Message)
	txGroup := v1.Group("/tx")
	{
		txGroup.GET("", txHandlers.List)
		txGroup.GET("/:hash", txHandlers.Get)
		txGroup.GET("/:hash/events", txHandlers.GetEvents)
		txGroup.GET("/:hash/messages", txHandlers.GetMessages)
	}

	namespaceHandlers := handler.NewNamespaceHandler(db.Namespace)
	namespaceGroup := v1.Group("/namespace")
	{
		namespaceGroup.GET("", namespaceHandlers.List)
		namespaceGroup.GET("/:id", namespaceHandlers.Get)
		namespaceGroup.GET("/:id/:version", namespaceHandlers.GetWithVersion)
	}
	namespaceByHashGroup := v1.Group("/namespace_by_hash")
	{
		namespaceByHashGroup.GET("/:hash", namespaceHandlers.GetByHash)
	}

	if cfg.Prometheus {
		v1.GET("/metrics", echoprometheus.NewHandler())
	}

	log.Info().Msg("API routes:")
	for _, route := range e.Routes() {
		log.Info().Msgf("[%s] %s -> %s", route.Method, route.Path, route.Name)
	}
}
