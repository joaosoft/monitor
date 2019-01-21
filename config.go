package monitor

import (
	"fmt"

	manager "github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Monitor MonitorConfig `json:"monitor"`
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
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	if appConfig.Monitor.Host == "" {
		appConfig.Monitor.Host = DefaultURL
	}

	return appConfig, simpleConfig, err
}
