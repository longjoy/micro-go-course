package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/longjoy/micro-go-course/section19/cargo/endpoint"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints *endpoint.CargoEndpoints) http.Handler {
	r := mux.NewRouter()

	kitLog := log.NewLogfmtLogger(os.Stderr)

	kitLog = log.With(kitLog, "ts", log.DefaultTimestampUTC)
	kitLog = log.With(kitLog, "caller", log.DefaultCaller)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/locations").Handler(kithttp.NewServer(
		endpoints.LocationsEndpoint,
		decodeLocationsRequest,
		encodeJSONResponse,
		options...,
	))
	r.Methods("POST").Path("/incidents").Handler(kithttp.NewServer(
		endpoints.RegisterHandlingEventEndpoint,
		decodeRegisterHandlingEventRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("POST").PathPrefix("/cargos").Path("/").Handler(kithttp.NewServer(
		endpoints.LoadCargoEndpoint,
		decodeLoadCargoRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("POST").PathPrefix("/cargos").Path("/change_destination").Handler(kithttp.NewServer(
		endpoints.ChangeDestinationEndpoint,
		decodeChangeDestinationRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("POST").PathPrefix("/cargos").Path("/").Handler(kithttp.NewServer(
		endpoints.AssignCargoToRouteEndpoint,
		decodeAssignCargoToRouteRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("POST").PathPrefix("/cargos").Path("/assign_to_route").Handler(kithttp.NewServer(
		endpoints.BookCargoEndpoint,
		decodeBookCargoRequest,
		encodeJSONResponse,
		options...,
	))

	return r
}

// decodeHealthCheckRequest decode request
func decodeLocationsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.LocationsRequest{}, nil
}

func decodeBookCargoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return nil, err
	}
	println("json:", string(body))

	var bcr endpoint.BookCargoRequest
	if err = json.Unmarshal(body, &bcr); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return nil, err
	}
	fmt.Printf("%+v", bcr)

	return &endpoint.BookCargoRequest{
		Origin:      bcr.Origin,
		Destination: bcr.Destination,
		Deadline:    bcr.Deadline,
	}, nil
}

func decodeAssignCargoToRouteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return nil, err
	}
	println("json:", string(body))

	var acr endpoint.AssignCargoToRouteRequest
	if err = json.Unmarshal(body, &acr); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return nil, err
	}
	fmt.Printf("%+v", acr)
	return &endpoint.AssignCargoToRouteRequest{
		Id:        acr.Id,
		Itinerary: acr.Itinerary,
	}, nil
}
func decodeLoadCargoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return nil, err
	}
	println("json:", string(body))

	var lcr endpoint.LoadCargoRequest
	if err = json.Unmarshal(body, &lcr); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return nil, err
	}
	fmt.Printf("%+v", lcr)
	return &endpoint.LoadCargoRequest{
		Id: lcr.Id,
	}, nil
}

func decodeChangeDestinationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return nil, err
	}
	println("json:", string(body))

	var cdr endpoint.ChangeDestinationRequest
	if err = json.Unmarshal(body, &cdr); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return nil, err
	}
	fmt.Printf("%+v", cdr)
	return &endpoint.ChangeDestinationRequest{
		Id:          cdr.Id,
		Destination: cdr.Destination,
	}, nil
}

func decodeRegisterHandlingEventRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return nil, err
	}
	println("json:", string(body))

	var rhe endpoint.RegisterHandlingEventRequest
	if err = json.Unmarshal(body, &rhe); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return nil, err
	}
	fmt.Printf("%+v", rhe)
	return &endpoint.RegisterHandlingEventRequest{
		Id:           rhe.Id,
		Completed:    rhe.Completed,
		VoyageNumber: rhe.VoyageNumber,
		UnLocode:     rhe.UnLocode,
		EventType:    rhe.EventType,
	}, nil
}

func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
