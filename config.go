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
func NewConfig() (*MonitorConfig, error) {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())

		return &MonitorConfig{}, err
	}

	return appConfig.Monitor, nil
}
