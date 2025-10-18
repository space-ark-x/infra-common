package intf

type IService interface {
	Init(*IModule)
	GetName() string
	Module() *IModule
}

type CoreService struct {
	Name string
	Mo   *IModule
}

func (s *CoreService) Init(mo *IModule) {
	s.Mo = mo
}

func (s *CoreService) GetName() string {
	return s.Name
}

func (s *CoreService) Module() *IModule {
	return s.Mo
}
