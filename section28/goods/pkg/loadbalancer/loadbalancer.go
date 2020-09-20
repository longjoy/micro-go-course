package loadbalancer

import (
	"errors"
	"github.com/longjoy/micro-go-course/section28/goods/pkg/discovery"
	"math/rand"
)

// 负载均衡器
type LoadBalancer interface {
	SelectService(service []*discovery.InstanceInfo) (*discovery.InstanceInfo, error)
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
