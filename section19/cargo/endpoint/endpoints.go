package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	shipping "github.com/longjoy/micro-go-course/section19/cargo/model"
	"github.com/longjoy/micro-go-course/section19/cargo/service/booking"
	"github.com/longjoy/micro-go-course/section19/cargo/service/handling"
	"time"
)

type CargoEndpoints struct {
	BookCargoEndpoint             endpoint.Endpoint
	LoadCargoEndpoint             endpoint.Endpoint
	AssignCargoToRouteEndpoint    endpoint.Endpoint
	ChangeDestinationEndpoint     endpoint.Endpoint
	CargosEndpoint                endpoint.Endpoint
	LocationsEndpoint             endpoint.Endpoint
	RegisterHandlingEventEndpoint endpoint.Endpoint
	//TrackEndpoint                 endpoint.Endpoint
}

type RegisterHandlingEventRequest struct {
	Completed    time.Time
	Id           shipping.TrackingID
	VoyageNumber shipping.VoyageNumber
	UnLocode     shipping.UNLocode
	EventType    shipping.HandlingEventType
}

type RegisterHandlingEventResponse struct {
	Res bool `json:"res"`
}

func RegisterHandlingEventEndpoint(handlingService handling.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*RegisterHandlingEventRequest)
		res, err := handlingService.RegisterHandlingEvent(req.Completed, req.Id, req.VoyageNumber, req.UnLocode, req.EventType)
		return &RegisterHandlingEventResponse{Res: res}, err
	}
}

type LocationsResponse struct {
	Locations []booking.Location `json:"locations"`
}

func LocationsEndpoint(bookService booking.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res := bookService.Locations()
		return &LocationsResponse{Locations: res}, err
	}
}

type CargosResponse struct {
	Cargos []booking.Cargo `json:"cargos"`
}

func CargosEndpoint(bookService booking.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res := bookService.Cargos()
		return &CargosResponse{Cargos: res}, err
	}
}

type ChangeDestinationRequest struct {
	Id          shipping.TrackingID
	Destination shipping.UNLocode
}

type ChangeDestinationResponse struct {
	Res bool `json:"res"`
}

func ChangeDestinationEndpoint(bookService booking.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*ChangeDestinationRequest)
		res, err := bookService.ChangeDestination(req.Id, req.Destination)
		return &ChangeDestinationResponse{Res: res}, err
	}
}

type AssignCargoToRouteRequest struct {
	Id        shipping.TrackingID
	Itinerary shipping.Itinerary
}

type AssignCargoToRouteResponse struct {
	Res bool `json:"res"`
}

func AssignCargoToRouteEndpoint(bookService booking.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*AssignCargoToRouteRequest)
		res, err := bookService.AssignCargoToRoute(req.Id, req.Itinerary)
		return &AssignCargoToRouteResponse{Res: res}, err
	}
}

type LoadCargoRequest struct {
	Id shipping.TrackingID
}

type LoadCargoResponse struct {
	Cargo booking.Cargo `json:"cargo"`
}

func MakeLoadCargoEndpoint(bookService booking.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*LoadCargoRequest)
		cargo, err := bookService.LoadCargo(req.Id)
		return &LoadCargoResponse{Cargo: cargo}, err

	}
}

type BookCargoRequest struct {
	Origin      shipping.UNLocode
	Destination shipping.UNLocode
	Deadline    time.Time
}

type BookCargoResponse struct {
	TrackingID shipping.TrackingID `json:"tracking_id"`
}

type LocationsRequest struct {
}

func MakeBookCargoEndpoint(bookService booking.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*BookCargoRequest)
		trackingID, err := bookService.BookNewCargo(req.Origin, req.Destination, req.Deadline)
		return &BookCargoResponse{TrackingID: trackingID}, err

	}
}

/*// UserEndpoint define endpoint
func (ce *CargoEndpoints) BookNewCargo(ctx context.Context, origin shipping.UNLocode, destination shipping.UNLocode,
	deadline time.Time) (shipping.TrackingID, error) {
	if origin == "" || destination == "" || deadline.IsZero() {
		return "", booking.ErrInvalidArgument
	}

	id := shipping.NextTrackingID()
	rs := shipping.RouteSpecification{
		Origin:          origin,
		Destination:     destination,
		ArrivalDeadline: deadline,
	}

	c := shipping.NewCargo(id, rs)

	if _, err := s.cargos.Store(c); err != nil {
		return "", err
	}

	return c.TrackingID, nil
}*/
