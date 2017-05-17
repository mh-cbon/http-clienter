package main

// file generated by
// github.com/mh-cbon/http-clienter
// do not edit

import (
	"errors"
	httper "github.com/mh-cbon/httper/lib"
	"net/http"
)

// HTTPClientControllerRPC is an http-clienter of *Controller.
// Controller of some resources.
type HTTPClientControllerRPC struct {
	Base string
}

// NewHTTPClientControllerRPC constructs an http-clienter of *Controller
func NewHTTPClientControllerRPC() *HTTPClientControllerRPC {
	ret := &HTTPClientControllerRPC{}
	return ret
}

// GetByID constructs a request to GetByID
func (t HTTPClientControllerRPC) GetByID(urlID int) (*http.Request, error) {
	return nil, errors.New("todo")
}

// UpdateByID constructs a request to UpdateByID
func (t HTTPClientControllerRPC) UpdateByID(urlID int, reqBody *Tomate) (*http.Request, error) {
	return nil, errors.New("todo")
}

// DeleteByID constructs a request to DeleteByID
func (t HTTPClientControllerRPC) DeleteByID(REQid int) (*http.Request, error) {
	return nil, errors.New("todo")
}

// TestVars1 constructs a request to TestVars1
func (t HTTPClientControllerRPC) TestVars1(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	return nil, errors.New("todo")
}

// TestCookier constructs a request to TestCookier
func (t HTTPClientControllerRPC) TestCookier(c httper.Cookier) (*http.Request, error) {
	return nil, errors.New("todo")
}

// TestSessionner constructs a request to TestSessionner
func (t HTTPClientControllerRPC) TestSessionner(s httper.Sessionner) (*http.Request, error) {
	return nil, errors.New("todo")
}

// TestRPCer constructs a request to TestRPCer
func (t HTTPClientControllerRPC) TestRPCer(id int) (*http.Request, error) {
	return nil, errors.New("todo")
}