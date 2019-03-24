package monitor

import (
	"fmt"

	manager "github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
)

// AppConfig ...
type AppConfig struct {
	Monitor *MonitorConfig `json:"monitor"`
}

// MonitorConfig ...
type MonitorConfig struct {
	Host string           `json:"host"`
	Db   manager.DBConfig `json:"db"`
	Migration         *migration.MigrationConfig `json:"migration"`
	Log  struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}
