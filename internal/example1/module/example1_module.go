package module

import (
	"github.com/space-ark-z/infra-common/pkg/app/core"
	"github.com/space-ark-z/infra-common/internal/example1/controller"
	"github.com/space-ark-z/infra-common/internal/example1/service"
)

type Example1Module struct{}

// NewExample1Module 创建一个新的示例模块实例
func NewExample1Module() *Example1Module {
	return &Example1Module{}
}

// RegisterControllers 注册控制器
func (m *Example1Module) RegisterControllers() []core.Controller {
	return []core.Controller{
		&controller.Example1Controller{},
	}
}

// RegisterServices 注册服务
func (m *Example1Module) RegisterServices() []interface{} {
	return []interface{}{
		service.NewExample1Service(),
	}
}

// RegisterDependencies 注册依赖
func (m *Example1Module) RegisterDependencies() []interface{} {
	return []interface{}{}
}