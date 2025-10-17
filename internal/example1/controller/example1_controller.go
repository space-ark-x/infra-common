package controller

import (
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
)

type Example1Controller struct {
	Ctx context.Context
}

// Get handles GET /
func (c *Example1Controller) Get() mvc.Result {
	return mvc.Response{
		Text: "Hello from Example1 Controller!",
	}
}

// GetBy handles GET /{id}
func (c *Example1Controller) GetBy(id string) mvc.Result {
	return mvc.Response{
		Text: "Get item with id: " + id,
	}
}

// RegisterRoutes 注册路由
func (c *Example1Controller) RegisterRoutes(party router.Party) {
	// 可以在这里自定义路由注册逻辑
	// 但通常使用 Iris MVC 的自动路由功能
}