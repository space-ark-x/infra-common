package core

import "github.com/kataras/iris/v12/core/router"

// Controller 控制器接口
type Controller interface {
	RegisterRoutes(party router.Party)
}

// Service 服务接口
type Service interface{}

// Module 模块接口
type Module interface {
	RegisterControllers() []Controller
	RegisterServices() []interface{}
	RegisterDependencies() []interface{}
}
