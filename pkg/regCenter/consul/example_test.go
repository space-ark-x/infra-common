package consul_test

import (
	"fmt"
	"log"
	"time"

	"space.ark-z.common/pkg/regCenter"
	"space.ark-z.common/pkg/regCenter/consul"
)

const SERVER_NAME = "example-user-provider"

func ExampleRegistry() {
	// 创建Consul配置
	config := &consul.Config{
		Address: "127.0.0.1:8500", // Consul服务器地址
		Token:   "your-acl-token", // ACL Token
	}

	// 创建注册中心实例
	registry, err := consul.NewRegistry(config)
	if err != nil {
		log.Fatal("Failed to create registry:", err)
	}
	defer registry.Close()

	// 获取本机IP
	ip, err := consul.GetLocalIP()
	if err != nil {
		log.Fatal("Failed to get local IP:", err)
	}

	// 创建服务实例
	instance := &regCenter.ServiceInstance{
		ID:      "service-1",
		Name:    SERVER_NAME,
		Address: ip,
		Port:    8080,
		Tags:    []string{"primary"},
		Meta: map[string]string{
			"version": "v1.0.0",
			"env":     "prod",
		},
		Weight:  10,
		Enable:  true,
		Healthy: true,
	}

	// 注册服务
	err = registry.Register(instance)
	if err != nil {
		log.Fatal("Failed to register service:", err)
	}

	fmt.Println("Service registered successfully")

	// 发现服务
	instances, err := registry.Discovery(SERVER_NAME)
	if err != nil {
		log.Fatal("Failed to discover service:", err)
	}

	fmt.Printf("Found %d service instances\n", len(instances))

	// 随机获取一个服务实例
	randomInstance, err := registry.RandomOne(SERVER_NAME)
	if err != nil {
		log.Fatal("Failed to get random service instance:", err)
	}

	fmt.Printf("Random instance: %s:%d\n", randomInstance.Address, randomInstance.Port)

	// 等待一段时间以观察心跳
	time.Sleep(30 * time.Second)

	// 输出:
	// Service registered successfully
	// Found 1 service instances
	// Random instance: 192.168.1.100:8080
}
