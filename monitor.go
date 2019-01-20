package monitor

import (
	"sync"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Monitor struct {
	config        *MonitorConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
}

// NewMonitor ...
func NewMonitor(options ...MonitorOption) (*Monitor, error) {
	config, simpleConfig, err := NewConfig()
	monitor := &Monitor{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		config: &config.Monitor,
	}

	if monitor.isLogExternal {
		monitor.pm.Reconfigure(manager.WithLogger(log))
	}

	if err != nil {
		log.Error(err.Error())
	} else {
		monitor.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Monitor.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	monitor.Reconfigure(options...)

	if monitor.config.Host == "" {
		monitor.config.Host = DefaultURL
	}

	simpleDB := manager.NewSimpleDB(&config.Monitor.Db)
	if err := monitor.pm.AddDB("db_postgres", simpleDB); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	web := manager.NewSimpleWebServer(monitor.config.Host)
	controller := NewController(NewInteractor(NewStoragePostgres(simpleDB)))
	controller.RegisterRoutes(web)

	monitor.pm.AddWeb("api_web", web)

	return monitor, nil
}

// Start ...
func (m *Monitor) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *Monitor) Stop() error {
	return m.pm.Stop()
}
