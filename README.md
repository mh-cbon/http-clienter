# http-clienter

[![travis Status](https://travis-ci.org/mh-cbon/http-clienter.svg?branch=master)](https://travis-ci.org/mh-cbon/http-clienter) [![Appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/http-clienter?branch=master&svg=true)](https://ci.appveyor.com/projects/mh-cbon/http-clienter) [![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/http-clienter)](https://goreportcard.com/report/github.com/mh-cbon/http-clienter) [![GoDoc](https://godoc.org/github.com/mh-cbon/http-clienter?status.svg)](http://godoc.org/github.com/mh-cbon/http-clienter) [![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Package http-clienter generates http client of a type


# TOC
- [Install](#install)
  - [Usage](#usage)
    - [$ http-clienter -help](#-http-clienter--help)
  - [Cli examples](#cli-examples)
- [API example](#api-example)
  - [Annotations](#annotations)
  - [> demo/main.go](#-demomaingo)
  - [> demo/httpclientcontroller.go](#-demohttpclientcontrollergo)
- [Recipes](#recipes)
  - [Release the project](#release-the-project)
- [History](#history)

# Install
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/http-clienter
cd $GOPATH/src/github.com/mh-cbon/http-clienter
git clone https://github.com/mh-cbon/http-clienter.git .
glide install
go install
```

## Usage

#### $ http-clienter -help
```sh
http-clienter 0.0.0

Usage

	http-clienter [-p name] [-mode name] [...types]

  types:  A list of types such as src:dst.
          A type is defined by its package path and its type name,
          [pkgpath/]name
          If the Package path is empty, it is set to the package name being generated.
          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName 
  -p:     The name of the package output.
  -mode:  The generation mode gorilla|std.
```

## Cli examples

```sh
# Create an http client of Tomate to MyTomate
http-clienter tomate_gen.go Tomate:MyTomate
```

# API example

Following example demonstates a program using it to generate an http cleint of a type.

#### Annotations

`http-clienter` reads and interprets annotations on `struct` and `methods`.

The `struct` annotations are used as default for the `methods` annotations.

| Name | Description |
| --- | --- |
| @route | The route path such as `/{param}` |
| @name | The route name `name` |
| @host | The route name `host` |
| @methods | The route methods `GET,POST,PUT` |
| @schemes | The route methods `http, https` |

#### > demo/main.go
```go
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
```

Following code is the generated implementation of the goriller binder.

#### > demo/httpclientcontroller.go
```go
package main

// file generated by
// github.com/mh-cbon/http-clienter
// do not edit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClientController is an http-clienter of *Controller.
// Controller of some resources.
type HTTPClientController struct {
	router *mux.Router
	Base   string
}

// NewHTTPClientController constructs an http-clienter of *Controller
func NewHTTPClientController(router *mux.Router) *HTTPClientController {
	ret := &HTTPClientController{
		router: router,
	}
	return ret
}

// GetByID constructs a request to /{id}
func (t HTTPClientController) GetByID(urlID int) (*http.Request, error) {
	var ret *http.Request
	var body io.Reader
	// var err error

	surl := "/{id}"
	surl = strings.Replace(surl, "{id}", fmt.Sprintf("%v", urlID), 1)
	url, URLerr := url.ParseRequestURI(surl)
	if URLerr != nil {
		return nil, URLerr
	}
	finalURL := url.String()
	finalURL = fmt.Sprintf("%v%v", t.Base, finalURL)

	req, reqErr := http.NewRequest("GET", finalURL, body)
	if reqErr != nil {
		return nil, reqErr
	}
	ret = req

	return ret, nil
}

// UpdateByID constructs a request to /{id}
func (t HTTPClientController) UpdateByID(urlID int, reqBody *Tomate) (*http.Request, error) {
	var ret *http.Request
	var body io.Reader
	// var err error

	data, reqBodyErr := json.Marshal(reqBody)
	if reqBodyErr != nil {
		return nil, reqBodyErr
	}

	body = bytes.NewBuffer(data)
	surl := "/{id}"
	surl = strings.Replace(surl, "{id}", fmt.Sprintf("%v", urlID), 1)
	url, URLerr := url.ParseRequestURI(surl)
	if URLerr != nil {
		return nil, URLerr
	}
	finalURL := url.String()
	finalURL = fmt.Sprintf("%v%v", t.Base, finalURL)

	req, reqErr := http.NewRequest("GET", finalURL, body)
	if reqErr != nil {
		return nil, reqErr
	}
	ret = req

	return ret, nil
}

// DeleteByID constructs a request to /{id}
func (t HTTPClientController) DeleteByID(REQid int) (*http.Request, error) {
	var ret *http.Request
	var body io.Reader
	// var err error

	surl := "/{id}"
	surl = strings.Replace(surl, "{id}", fmt.Sprintf("%v", REQid), 1)
	url, URLerr := url.ParseRequestURI(surl)
	if URLerr != nil {
		return nil, URLerr
	}
	finalURL := url.String()
	finalURL = fmt.Sprintf("%v%v", t.Base, finalURL)

	req, reqErr := http.NewRequest("GET", finalURL, body)
	if reqErr != nil {
		return nil, reqErr
	}
	ret = req

	return ret, nil
}
```


# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)
