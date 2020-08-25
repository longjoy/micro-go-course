package service

import (
	"errors"
)

const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")

	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

type Service interface {
	// Concat a and b
	Concat(req StringRequest, ret *string) error
}

type StringService struct {
}

func (s StringService) Concat(req StringRequest, ret *string) error {
	// test for length overflow
	if len(req.A)+len(req.B) > StrMaxSize {
		*ret = ""
		return ErrMaxSize
	}
	*ret = req.A + req.B
	return nil
}
