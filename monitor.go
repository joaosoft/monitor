package monitor

import (
	"sync"

	"github.com/labstack/gommon/log"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
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
		logger: logger.NewLogDefault("monitor", logger.WarnLevel),
		config: config.Monitor,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Monitor != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Monitor.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	} else {
		config.Monitor = &MonitorConfig{
			Host: DefaultURL,
		}
	}

	service.Reconfigure(options...)

	// execute migrations
	migrationService, err := migration.NewCmdService(migration.WithCmdConfiguration(service.config.Migration))
	if err != nil {
		return nil, err
	}

	if _, err := migrationService.Execute(migration.OptionUp, 0, migration.ExecutorModeDatabase); err != nil {
		return nil, err
	}

	simpleDB := service.pm.NewSimpleDB(&config.Monitor.Db)
	if err := service.pm.AddDB("db_postgres", simpleDB); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	web := service.pm.NewSimpleWebServer(service.config.Host)
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
