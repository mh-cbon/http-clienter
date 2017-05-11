package main

import (
	"net/http"

	httper "github.com/mh-cbon/httper/lib"
)

//go:generate lister vegetables_gen.go *Tomate:Tomates
//go:generate channeler tomate_chan_gen.go *Tomates:ChanTomates

//go:generate jsoner -mode gorilla json_controller_gen.go *Controller:JSONController
//go:generate httper -mode gorilla http_vegetables_gen.go *JSONController:HTTPController
//go:generate goriller goriller_vegetables_gen.go *HTTPController:GorillerTomate
//go:generate http-clienter -mode gorilla http_client_gen.go *Controller:HTTPClientController

func main() {

	backend := NewChanTomates()
	backend.Push(&Tomate{Name: "red"})

}

// Tomate is about red vegetables to make famous italian food.
type Tomate struct {
	ID   int
	Name string
}

// GetID return the ID of the Tomate.
func (t *Tomate) GetID() int {
	return t.ID
}

// TomateBackend ...
type TomateBackend interface {
	Filter(...func(*Tomate) bool) *ChanTomates // i want to return interface here, like TomateBackend.
	First() *Tomate
	Remove(*Tomate) bool
}

// Controller of some resources.
type Controller struct {
	backend TomateBackend
}

// NewController ...
func NewController(backend TomateBackend) *Controller {
	return &Controller{
		backend: backend,
	}
}

// GetByID ...
// @route /{id}
// @methods GET
func (t *Controller) GetByID(urlID int) *Tomate {
	return t.backend.Filter(FilterTomates.ByID(urlID)).First()
}

// UpdateByID ...
// @route /{id}
// @methods PUT,POST
func (t *Controller) UpdateByID(urlID int, reqBody *Tomate) *Tomate {
	var ret *Tomate
	t.backend.Filter(func(v *Tomate) bool {
		if v.ID == urlID {
			v.Name = reqBody.Name
			ret = v
		}
		return true
	})
	return ret
}

// DeleteByID ...
// @route /{id}
// @methods DELETE
func (t *Controller) DeleteByID(REQid int) bool {
	return t.backend.Remove(&Tomate{ID: REQid})
}

// TestVars1 ...
func (t *Controller) TestVars1(w http.ResponseWriter, r *http.Request) {
}

// TestCookier ...
func (t *Controller) TestCookier(c httper.Cookier) {
}

// TestSessionner ...
func (t *Controller) TestSessionner(s httper.Sessionner) {
}

// TestRPCer ...
func (t *Controller) TestRPCer(id int) bool {
	return false
}
