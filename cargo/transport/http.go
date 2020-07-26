package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/longjoy/micro-go-course/section08/user/endpoint"
	"net/http"
	"os"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints *endpoint.UserEndpoints) http.Handler {
	r := mux.NewRouter()

	kitLog := log.NewLogfmtLogger(os.Stderr)

	kitLog = log.With(kitLog, "ts", log.DefaultTimestampUTC)
	kitLog = log.With(kitLog, "caller", log.DefaultCaller)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeJSONResponse,
		options...,
	))


	r.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeJSONResponse,
		options...,
	))



	return r
}
func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || password == "" || email == ""{
		return nil, ErrorBadRequest
	}
	return &endpoint.RegisterRequest{
		Username:username,
		Password:password,
		Email:email,
	},nil
}


func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return nil, ErrorBadRequest
	}
	return &endpoint.LoginRequest{
		Email:email,
		Password:password,
	},nil
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

