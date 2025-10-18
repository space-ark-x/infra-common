package app

import "github.com/space-ark-z/infra-common/pkg/config"

func newVer() *Ver {
	return &Ver{}
}

type Ver struct {
	Config *config.Type
}
