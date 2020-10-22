package loadbalancer

import (
	"errors"
	"github.com/longjoy/micro-go-course/section28/goods/pkg/discovery"
	"math/rand"
	"sort"
)

// 负载均衡器
type LoadBalancer interface {
	SelectService(service []*discovery.InstanceInfo) (*discovery.InstanceInfo, error)
	SelectServiceByKey(service []*discovery.InstanceInfo, key string) (*discovery.InstanceInfo, error)
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{}
}

func NewWeightRoundRobinLoadBalancer() *WeightRoundRobinLoadBalancer {
	return &WeightRoundRobinLoadBalancer{}
}

type RandomLoadBalancer struct {
}

// 随机负载均衡
func (loadBalance *RandomLoadBalancer) SelectService(services []*discovery.InstanceInfo) (*discovery.InstanceInfo, error) {

	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}

	return services[rand.Intn(len(services))], nil
}

func (loadBalance *RandomLoadBalancer) SelectServiceByKey(services []*discovery.InstanceInfo, key string) (*discovery.InstanceInfo, error) {
	return loadBalance.SelectService(services)
}

type WeightRoundRobinLoadBalancer struct {
}

// 权重平滑负载均衡
func (loadBalance *WeightRoundRobinLoadBalancer) SelectService(services []*discovery.InstanceInfo) (best *discovery.InstanceInfo, err error) {

	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}

	total := 0
	for i := 0; i < len(services); i++ {
		w := services[i]
		if w == nil {
			continue
		}

		w.CurWeight += w.Weights.Passing

		total += w.Weights.Passing

		if best == nil || w.CurWeight > best.CurWeight {
			best = w
		}
	}

	if best == nil {
		return nil, nil
	}

	best.CurWeight -= total
	return best, nil
}

func (loadBalance *WeightRoundRobinLoadBalancer) SelectServiceByKey(services []*discovery.InstanceInfo, key string) (*discovery.InstanceInfo, error) {
	return loadBalance.SelectService(services)
}

func NewHashLoadBalancer() *HashLoadBalancer {
	return &HashLoadBalancer{}
}

func (loadBalance *HashLoadBalancer) SelectService(services []*discovery.InstanceInfo) (*discovery.InstanceInfo, error) {

	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}
	return services[rand.Intn(len(services))], nil
}

type HashLoadBalancer struct {
}

func (loadBalance *HashLoadBalancer) SelectServiceByKey(services []*discovery.InstanceInfo, key string) (*discovery.InstanceInfo, error) {
	lens := len(services)
	if services == nil || lens == 0 {
		return nil, errors.New("service instances are not exist")
	}

	nodeWeight := make(map[string]int)
	instanceMap := make(map[string]*discovery.InstanceInfo)
	for i := 0; i < len(services); i++ {
		instance := services[i]
		nodeWeight[instance.Address] = i
		instanceMap[instance.Address] = instance
	}
	sort.Sort()
	// 建立Hash环
	hash := NewHashRing()
	// 添加各个服务实例到环上
	hash.AddNodes(nodeWeight)
	// 根据请求的key来获取对应的服务实例
	host := hash.GetNode(key)
	return instanceMap[host], nil
}
