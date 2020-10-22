package loadbalancer

import (
	//	"fmt"
	"github.com/longjoy/micro-go-course/section28/goods/pkg/discovery"
	"testing"
)

func TestHash(t *testing.T) {

	instances := make([]*discovery.InstanceInfo, 3)
	instances[0] = &discovery.InstanceInfo{}
	instances[0].Address = "192.168.1.1"
	instances[1] = &discovery.InstanceInfo{}
	instances[1].Address = "192.168.1.2"
	instances[2] = &discovery.InstanceInfo{}
	instances[2].Address = "192.168.1.3"

	nodeWeight := make(map[string]int)
	instanceMap := make(map[string]*discovery.InstanceInfo)
	for i := 0; i < len(instances); i++ {
		instance := instances[i]
		nodeWeight[instance.Address] = i + 1
		instanceMap[instance.Address] = instance
		println("info:" + instance.Address)
	}
	// 建立Hash环
	hash := NewHashRing()
	// 添加各个服务实例到环上
	hash.AddNodes(nodeWeight)
	// 根据请求的key来获取对应的服务实例
	host := hash.GetNode("2")
	println("ddd" + host)

}
