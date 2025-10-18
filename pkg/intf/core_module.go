package intf

import (
	"github.com/kataras/iris/v12"
)

type IModule interface {
	Name() string
	Build(app *iris.Application, mo IModule)
	GetService(s IService) (IService, bool)
	AddSubModule(IModule)
	AddController(con IController)
	AddService(s IService)
}

type CoreModule struct {
	mList map[string]IModule
	sList map[string]IService
	cList map[string]IController
	name  string
}

func NewModule(name string) IModule {
	return &CoreModule{
		name:  name,
		mList: make(map[string]IModule),
		sList: make(map[string]IService),
		cList: make(map[string]IController),
	}
}

func (m *CoreModule) Build(app *iris.Application, mo IModule) {
	for _, module := range m.mList {
		module.Build(app, m)
	}
	for _, controller := range m.cList {
		controller.Build(app)
	}
}

func (m *CoreModule) Name() string {
	return m.name
}

func (m *CoreModule) GetService(s IService) (IService, bool) {
	service, ok := m.sList[s.GetName()]
	return service, ok
}

func (m *CoreModule) AddController(con IController) {
	_, ok := m.cList[(con).GetName()]
	if !ok {
		(con).Init(m)
		m.cList[(con).GetName()] = con
		return
	}
	panic("controller already exists")
}

func (m *CoreModule) AddService(s IService) {
	_, ok := m.sList[s.GetName()]
	if !ok {
		s.Init(m)
		m.sList[s.GetName()] = s
		return
	}
	panic("service already exists")
}

func (m *CoreModule) AddSubModule(sub IModule) {
	_, ok := m.mList[sub.Name()]
	if !ok {
		m.mList[sub.Name()] = sub
		return
	}
	panic("sub module already exists")
}

func (m *CoreModule) SetInterceptor() {

}
