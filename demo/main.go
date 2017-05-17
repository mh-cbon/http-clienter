package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httper "github.com/mh-cbon/httper/lib"
)

//go:generate lister *Tomate:TomatesGen
//go:generate channeler TomatesGen:TomatesSyncGen

//go:generate jsoner -mode gorilla *Controller:ControllerJSONGen
//go:generate httper -mode gorilla *ControllerJSONGen:ControllerHTTPGen
//go:generate goriller *ControllerHTTPGen:ControllerGoriller
//go:generate goriller -mode rpc *ControllerHTTPGen:ControllerGorillerRPC

//go:generate http-clienter -mode gorilla *Controller:HTTPClientController
//go:generate http-clienter -mode std *Controller:HTTPClientControllerRPC

func main() {

	backend := NewTomatesSyncGen()
	backend.Push(&Tomate{Name: "red"})

	router := mux.NewRouter()

	controller := NewController(backend)
	jsoner := NewControllerJSONGen(controller, nil)
	httper := NewControllerHTTPGen(jsoner, nil)
	goriller := NewControllerGoriller(httper)

	goriller.Bind(router)

	http.Handle("/", router)

	client := NewHTTPClientController(router)
	client.Base = "http://localhost:8080"

	go func() {
		<-time.After(time.Second)
		req, err := client.GetByID(0)
		fmt.Println(err)
		fmt.Println(http.DefaultClient.Do(req))
	}()

	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
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
	// what if i want to return interface here, like TomateBackend.
	Filter(...func(*Tomate) bool) *TomatesGen
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
	res := t.backend.Filter(FilterTomatesGen.ByID(urlID))
	fmt.Println("res", res)
	return res.First()
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
