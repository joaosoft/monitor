package monitor

import (
	"fmt"

	manager "github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Monitor *MonitorConfig `json:"monitor"`
}

// MonitorConfig ...
type MonitorConfig struct {
	Host string           `json:"host"`
	Db   manager.DBConfig `json:"db"`
	Log  struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig(host string, db manager.DBConfig) *MonitorConfig {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	}

	appConfig.Monitor.Host = host

	return appConfig.Monitor
}
