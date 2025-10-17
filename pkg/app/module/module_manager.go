package module

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/space-ark-z/infra-common/pkg/app/core"
)

type ModuleManager struct {
	app     *iris.Application
	modules []core.Module
}

func NewModuleManager(app *iris.Application) *ModuleManager {
	return &ModuleManager{
		app:     app,
		modules: make([]core.Module, 0),
	}
}

// RegisterModule 注册模块
func (mm *ModuleManager) RegisterModule(module core.Module) {
	mm.modules = append(mm.modules, module)
}

// Init 初始化所有模块
func (mm *ModuleManager) Init() {
	for _, module := range mm.modules {
		// 创建一个新的MVC应用
		mvcApp := mvc.New(mm.app.Party("/"))
		
		// 注册依赖
		for _, dep := range module.RegisterDependencies() {
			mvcApp.Register(dep)
		}
		
		// 注册服务
		for _, service := range module.RegisterServices() {
			mvcApp.Register(service)
		}
		
		// 注册控制器
		for _, controller := range module.RegisterControllers() {
			mvcApp.Handle(controller)
		}
	}
}