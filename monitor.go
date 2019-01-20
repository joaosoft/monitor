package monitor

import (
	"sync"

	"github.com/labstack/gommon/log"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Monitor struct {
	logger        logger.ILogger
	config        *MonitorConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
}

// NewMonitor ...
func NewMonitor(options ...MonitorOption) (*Monitor, error) {
	config, simpleConfig, err := NewConfig()
	service := &Monitor{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		logger: logger.NewLogDefault("service", logger.InfoLevel),
		config: &config.Monitor,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Monitor.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.Reconfigure(options...)

	if service.config.Host == "" {
		service.config.Host = DefaultURL
	}

	simpleDB := manager.NewSimpleDB(&config.Monitor.Db)
	if err := service.pm.AddDB("db_postgres", simpleDB); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	web := manager.NewSimpleWebServer(service.config.Host)
	controller := service.NewController(service.NewInteractor(service.NewStoragePostgres(simpleDB)))
	controller.RegisterRoutes(web)

	service.pm.AddWeb("api_web", web)

	return service, nil
}

// Start ...
func (m *Monitor) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *Monitor) Stop() error {
	return m.pm.Stop()
}
