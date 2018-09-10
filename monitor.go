package monitor

import (
	"fmt"

	"sync"

	logger "github.com/joaosoft/logger"
	manager "github.com/joaosoft/manager"
)

type Monitor struct {
	config        *MonitorConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
}

// NewMonitor ...
func NewMonitor(options ...MonitorOption) (*Monitor, error) {
	monitor := &Monitor{
		pm: manager.NewManager(manager.WithRunInBackground(false)),
	}

	if monitor.isLogExternal {
		monitor.pm.Reconfigure(manager.WithLogger(log))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		monitor.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Monitor.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	monitor.config = &appConfig.Monitor

	monitor.Reconfigure(options...)

	if monitor.config.Host == "" {
		monitor.config.Host = DefaultURL
	}

	simpleDB := manager.NewSimpleDB(&appConfig.Monitor.Db)
	if err := monitor.pm.AddDB("db_postgres", simpleDB); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	web := manager.NewSimpleWebServer(monitor.config.Host)
	controller := NewDbMigration(NewInteractor(NewStoragePostgres(simpleDB)))
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
