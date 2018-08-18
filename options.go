package monitor

import (
	logger "github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// MonitorOption ...
type MonitorOption func(client *Monitor)

// Reconfigure ...
func (monitor *Monitor) Reconfigure(options ...MonitorOption) {
	for _, option := range options {
		option(monitor)
	}
}

// WithConfiguration ...
func WithConfiguration(config *MonitorConfig) MonitorOption {
	return func(client *Monitor) {
		client.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) MonitorOption {
	return func(monitor *Monitor) {
		log = logger
		monitor.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) MonitorOption {
	return func(monitor *Monitor) {
		log.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) MonitorOption {
	return func(dropbox *Monitor) {
		dropbox.pm = mgr
	}
}
