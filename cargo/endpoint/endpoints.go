package endpoint

import (
	"github.com/go-kit/kit/endpoint"
)

type CargoEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}
