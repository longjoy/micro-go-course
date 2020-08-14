package discovery

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"os"
	"strconv"
)


type DiscoveryClient struct {
	client consul.Client
	register sd.Registrar
	config *api.Config
	registration *api.AgentServiceRegistration
}

func NewAgentServiceRegistration(serviceName, instanceId, healthCheckUrl, serviceAddr string, servicePort int, meta map[string]string) *api.AgentServiceRegistration {
	return &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: serviceAddr,
		Port:    servicePort,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + serviceAddr + ":" + strconv.Itoa(servicePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

}

func NewDiscoveryClient(host string, port int, registration *api.AgentServiceRegistration) (*DiscoveryClient, error) {
	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)

	client, err := api.NewClient(config)
	if err != nil{
		return nil, err
	}

	sdClient := consul.NewClient(client)


	return &DiscoveryClient{
		client: sdClient,
		config: config,
		registration: registration,
		register: consul.NewRegistrar(sdClient, registration, log.NewLogfmtLogger(os.Stderr)),
	}, nil
}

func (consulClient *DiscoveryClient) Register(ctx context.Context)  {
	consulClient.register.Register()
}

func (consulClient *DiscoveryClient) Deregister(ctx context.Context)  {
	consulClient.register.Deregister()
}


func (consulClient *DiscoveryClient) DiscoverServices(ctx context.Context, serviceName string) ([] *api.AgentService, error) {

	result, _, err := consulClient.client.Service(serviceName, "", false, nil)

	if err != nil{
		return nil, err
	}

	rsp := make([]*api.AgentService, 0, len(result))

	for _, v := range result{
		rsp = append(rsp, v.Service)
	}

	return rsp, err

}
