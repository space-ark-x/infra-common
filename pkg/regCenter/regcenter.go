package regCenter

// ServiceInstance 服务实例接口
type ServiceInstance struct {
	ID       string            // 服务ID
	Name     string            // 服务名称
	Address  string            // 服务地址
	Port     int               // 服务端口
	Tags     []string          // 服务标签
	Meta     map[string]string // 元数据
	Weight   int               // 权重
	Enable   bool              // 是否启用
	Healthy  bool              // 是否健康
}

// Registry 注册中心接口
type Registry interface {
	// Register 注册服务
	Register(instance *ServiceInstance) error
	
	// Deregister 注销服务
	Deregister(serviceID string) error
	
	// Discovery 根据服务名发现服务实例
	Discovery(serviceName string) ([]*ServiceInstance, error)
	
	// GetService 根据服务ID获取单个服务实例
	GetService(serviceID string) (*ServiceInstance, error)
	
	// RandomOne 随机获取一个健康的服务实例
	RandomOne(serviceName string) (*ServiceInstance, error)
	
	// Close 关闭注册中心
	Close() error
}