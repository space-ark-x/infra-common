package consul

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"space.ark-z.common/pkg/regCenter"

	"github.com/hashicorp/consul/api"
)

// Config Consul配置
type Config struct {
	// Consul服务器地址
	Address string `json:"address" yaml:"address"`
	
	// ACL Token用于认证
	Token string `json:"token" yaml:"token"`
	
	// 连接超时时间
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	
	// 心跳间隔
	HeartbeatInterval time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	
	// 健康检查失败后的注销时间
	DeregisterCriticalServiceAfter time.Duration `json:"deregister_critical_service_after" yaml:"deregister_critical_service_after"`
}

// Registry Consul注册中心
type Registry struct {
	cli    *api.Client
	config *Config
	
	// 用于存储当前服务的注册信息，以便在关闭时注销
	services map[string]*api.AgentServiceRegistration
	
	// 心跳检测的上下文和取消函数
	heartbeatCtx    context.Context
	heartbeatCancel context.CancelFunc
}

// NewRegistry 创建一个新的Consul注册中心实例
func NewRegistry(config *Config) (*Registry, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	
	// 设置默认值
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}
	
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 10 * time.Second
	}
	
	if config.DeregisterCriticalServiceAfter == 0 {
		config.DeregisterCriticalServiceAfter = 30 * time.Second
	}
	
	// 创建Consul客户端配置
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Address
	consulConfig.Token = config.Token
	consulConfig.HttpClient.Timeout = config.Timeout
	
	// 创建Consul客户端
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}
	
	// 测试连接
	_, err = client.Agent().Self()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to consul: %w", err)
	}
	
	// 初始化上下文
	ctx, cancel := context.WithCancel(context.Background())
	
	registry := &Registry{
		cli:             client,
		config:          config,
		services:        make(map[string]*api.AgentServiceRegistration),
		heartbeatCtx:    ctx,
		heartbeatCancel: cancel,
	}
	
	return registry, nil
}

// Register 注册服务
func (r *Registry) Register(instance *regCenter.ServiceInstance) error {
	if instance == nil {
		return errors.New("service instance is nil")
	}
	
	// 构建服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      instance.ID,
		Name:    instance.Name,
		Address: instance.Address,
		Port:    instance.Port,
		Tags:    instance.Tags,
		Meta:    instance.Meta,
		Weights: &api.AgentWeights{
			Passing: instance.Weight,
			Warning: instance.Weight,
		},
		Check: &api.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("service:%s", instance.ID),
			TTL:                            (r.config.HeartbeatInterval * 2).String(),
			DeregisterCriticalServiceAfter: r.config.DeregisterCriticalServiceAfter.String(),
			Status:                         api.HealthPassing,
		},
	}
	
	// 注册服务
	err := r.cli.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	
	// 存储注册信息
	r.services[instance.ID] = registration
	
	// 启动心跳检测
	go r.heartbeat(instance.ID)
	
	return nil
}

// Deregister 注销服务
func (r *Registry) Deregister(serviceID string) error {
	if serviceID == "" {
		return errors.New("service id is empty")
	}
	
	// 从Consul注销服务
	err := r.cli.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	
	// 从本地存储中移除
	delete(r.services, serviceID)
	
	return nil
}

// Discovery 根据服务名发现服务实例
func (r *Registry) Discovery(serviceName string) ([]*regCenter.ServiceInstance, error) {
	if serviceName == "" {
		return nil, errors.New("service name is empty")
	}
	
	// 查询服务实例
	services, _, err := r.cli.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service: %w", err)
	}
	
	// 转换为统一格式
	instances := make([]*regCenter.ServiceInstance, 0, len(services))
	for _, service := range services {
		instances = append(instances, &regCenter.ServiceInstance{
			ID:      service.Service.ID,
			Name:    service.Service.Service,
			Address: service.Service.Address,
			Port:    service.Service.Port,
			Tags:    service.Service.Tags,
			Meta:    service.Service.Meta,
			Weight:  service.Service.Weights.Passing,
			Enable:  true,
			Healthy: true,
		})
	}
	
	return instances, nil
}

// GetService 根据服务ID获取单个服务实例
func (r *Registry) GetService(serviceID string) (*regCenter.ServiceInstance, error) {
	if serviceID == "" {
		return nil, errors.New("service id is empty")
	}
	
	// 查询服务实例
	service, _, err := r.cli.Agent().Service(serviceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	
	// 转换为统一格式
	instance := &regCenter.ServiceInstance{
		ID:      service.ID,
		Name:    service.Service,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Meta:    service.Meta,
		Weight:  service.Weights.Passing,
		Enable:  true,
		Healthy: true,
	}
	
	return instance, nil
}

// RandomOne 随机获取一个健康的服务实例
func (r *Registry) RandomOne(serviceName string) (*regCenter.ServiceInstance, error) {
	instances, err := r.Discovery(serviceName)
	if err != nil {
		return nil, err
	}
	
	if len(instances) == 0 {
		return nil, errors.New("no available service instances")
	}
	
	// 随机返回一个实例
	return instances[rand.Intn(len(instances))], nil
}

// Close 关闭注册中心，注销所有已注册的服务
func (r *Registry) Close() error {
	// 取消心跳检测
	if r.heartbeatCancel != nil {
		r.heartbeatCancel()
	}
	
	// 注销所有服务
	for serviceID := range r.services {
		_ = r.Deregister(serviceID)
	}
	
	return nil
}

// heartbeat 心跳检测，定期更新服务状态
func (r *Registry) heartbeat(serviceID string) {
	ticker := time.NewTicker(r.config.HeartbeatInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-r.heartbeatCtx.Done():
			return
		case <-ticker.C:
			// 更新TTL
			_ = r.cli.Agent().UpdateTTL(
				fmt.Sprintf("service:%s", serviceID),
				"",
				api.HealthPassing,
			)
		}
	}
}

// GetLocalIP 获取本机IP地址
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}