package config

import "os"

type Record map[string]string

type Type struct {
	AppName     string `yaml:"APP_NAME"`
	AppPort     int    `yaml:"APP_PORT"`
	HealthCheck string `yaml:"HEALTH_CHECK"` // 健康检查地址

	Mysql       bool   `yaml:"MYSQL"`
	MysqlDsn    string `yaml:"MYSQL_DSN"`
	AutoMigrate bool   `yaml:"AUTO_MIGRATE"`

	Mongo    bool   `yaml:"MONGO"`
	MongoDsn string `yaml:"MONGO_DSN"`
	Record   Record
	env      string
}

func (t *Type) Get(record string, defaultValue string) string {
	envValue := os.Getenv(record)
	if envValue != "" {
		return envValue
	}
	configValue := t.Record[record]
	if configValue != "" {
		return configValue
	}
	return defaultValue
}
