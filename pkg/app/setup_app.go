package app

import (
	"github.com/space-ark-z/infra-common/pkg/config"
)

var appVer *Ver

type SetupAppArgs struct {
}

func SetupApp() {
	appVer = newVer()
	appVer.Config = config.LoadConfig()
}

// GetAppVer 获取应用版本信息
func GetAppVer() *Ver {
	return appVer
}
