package app

import (
	"github.com/space-ark-z/infra-common/pkg/config"
)

var app_ver *ver

type SetupAppArgs struct {
}

func SetupApp() {
	app_ver = newVer()
	app_ver.Config = config.LoadConfig()
}

// GetAppVer 获取应用版本信息
func GetAppVer() *ver {
	return app_ver
}