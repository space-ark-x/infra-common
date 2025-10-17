package app

import "github.com/space-ark-z/infra-common/pkg/config"

func newVer() *ver {
	return &ver{}
}

type ver struct {
	Config *config.Type
}
