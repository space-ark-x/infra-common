package main

import (
	"github.com/space-ark-z/infra-common/pkg/app/application"
	example1 "github.com/space-ark-z/infra-common/internal/example1/module"
)

func main() {
	// 创建应用实例
	app := application.New()
	
	// 注册模块
	app.RegisterModule(example1.NewExample1Module())
	
	// 设置应用
	app.Setup()
	
	// 启动应用
	app.Run(":8080")
}