package application

import (
	"github.com/kataras/iris/v12"
	"github.com/space-ark-z/infra-common/pkg/app/module"
	"github.com/space-ark-z/infra-common/pkg/app/core"
)

type Application struct {
	irisApp       *iris.Application
	moduleManager *module.ModuleManager
}

// New 创建一个新的应用实例
func New() *Application {
	app := iris.New()
	
	// 创建应用实例
	application := &Application{
		irisApp: app,
	}
	
	// 初始化模块管理器
	application.moduleManager = module.NewModuleManager(app)
	
	return application
}

// RegisterModule 注册模块
func (app *Application) RegisterModule(module core.Module) {
	app.moduleManager.RegisterModule(module)
}

// Setup 初始化应用
func (app *Application) Setup() {
	app.moduleManager.Init()
}

// Run 启动应用
func (app *Application) Run(addr string) error {
	return app.irisApp.Listen(addr)
}

// GetIrisApp 获取iris应用实例
func (app *Application) GetIrisApp() *iris.Application {
	return app.irisApp
}