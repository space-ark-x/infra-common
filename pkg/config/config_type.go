package config

type Record map[string]any

type Type struct {
	AppName     string `yaml:"app_name"`
	AppPort     int    `yaml:"app_port"`
	HealthCheck string `yaml:"health_check"` // 健康检查地址

	Mysql       bool   `yaml:"mysql"`
	MysqlDsn    string `yaml:"mysql_dsn"`
	AutoMigrate bool   `yaml:"auto_migrate"`

	Mongo    bool   `yaml:"mongo"`
	MongoDsn string `yaml:"mongo_dsn"`
	Record
}
