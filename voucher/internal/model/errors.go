package model

import (
	"errors"
)

var (
	ErrDuplicate error = errors.New("duplicate")
)

type ResponseError struct {
	err        error
	StatusCode int
	StatusDesc string
}

func (r ResponseError) Error() string {
	return r.err.Error()
}

func (r ResponseError) DescError() string {
	if len(r.StatusDesc) == 0 {
		return r.Error()
	}
	return r.StatusDesc
}

type ResponseErrorOpt func(*ResponseError)

func WithStatusCode(code int) ResponseErrorOpt {
	return func(s *ResponseError) {
		s.StatusCode = code
	}
}

func WithDescription(desc string) ResponseErrorOpt {
	return func(s *ResponseError) {
		s.StatusDesc = desc
	}
}

func NewRespError(err error, opt ...ResponseErrorOpt) *ResponseError {
	var respErr = &ResponseError{err: err}
	for _, o := range opt {
		o(respErr)
	}
	return respErr
}

var (
	ErrDBNotFound error = errors.New("db not found")
)
