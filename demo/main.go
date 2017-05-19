package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

	client := &httpClient{base: "http://localhost:8080"}
	api := NewHTTPClientController(router)

	go func() {
		log.Fatal(
			http.ListenAndServe(":8080", nil),
		)
	}()

	<-time.After(time.Second)
	req, err := api.GetByID(0)
	if err != nil {
		log.Fatal(err)
	}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		log.Fatal(reqErr)
	}
	defer res.Body.Close()
	fmt.Println(res)
	io.Copy(os.Stdout, res.Body)
	fmt.Println()
}

type httpClient struct {
	http.Client
	base string
}

func (c httpClient) Do(req *http.Request) (*http.Response, error) {
	newURL, err := url.Parse(fmt.Sprintf("%v%v", c.base, req.URL.String()))
	if err != nil {
		return nil, err
	}
	req.URL = newURL
	return c.Client.Do(req)
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
