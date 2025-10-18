package intf

import "github.com/kataras/iris/v12"

type IController interface {
	Build(app *iris.Application)
	GetName() string
	Init(mo *IModule)
}

var _ IController = (*CoreController)(nil)

type CoreController struct {
	Mo   *IModule
	Name string
}

func (c *CoreController) Build(app *iris.Application) {

}

func (c *CoreController) GetName() string {
	return c.Name
}

func (c *CoreController) Init(mo *IModule) {
	c.Mo = mo
}
