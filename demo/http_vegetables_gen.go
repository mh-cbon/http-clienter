package main

// file generated by
// github.com/mh-cbon/httper
// do not edit

import (
	httper "github.com/mh-cbon/httper/lib"
	"io"
	"net/http"
	"strconv"
)

var xxStrconvAtoi = strconv.Atoi
var xxIoCopy = io.Copy
var xxHTTPOk = http.StatusOK

// HTTPController is an httper of *JSONController.
// JSONController is jsoner of *Controller.
// Controller of some resources.
type HTTPController struct {
	embed     *JSONController
	cookier   httper.CookieProvider
	dataer    httper.DataerProvider
	sessioner httper.SessionProvider
}

// NewHTTPController constructs an httper of *JSONController
func NewHTTPController(embed *JSONController) *HTTPController {
	ret := &HTTPController{
		embed:     embed,
		cookier:   &httper.CookieHelperProvider{},
		dataer:    &httper.GorillaHTTPDataProvider{},
		sessioner: &httper.GorillaSessionProvider{},
	}
	return ret
}

// HandleError returns http 500 and prints the error.
func (t *HTTPController) HandleError(err error, w http.ResponseWriter, r *http.Request) bool {
	if err == nil {
		return false
	}
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, err.Error())
	return true
}

// HandleSuccess calls for embed.HandleSuccess method.
func (t *HTTPController) HandleSuccess(w http.ResponseWriter, r io.Reader) error {
	return t.embed.HandleSuccess(w, r)
}

// GetByID invoke *JSONController.GetByID using the request body as a json payload.
// GetByID Decodes reqBody as json to invoke *Controller.GetByID.
// Other parameters are passed straight
// GetByID ...
// @route /{id}
// @methods GET
func (t *HTTPController) GetByID(w http.ResponseWriter, r *http.Request) {
	var urlID int
	tempurlID, err := strconv.Atoi(t.dataer.Make(w, r).Get("url", "id"))
	if t.HandleError(err, w, r) {
		return
	}
	urlID = tempurlID

	res, err := t.embed.GetByID(urlID)
	if t.HandleError(err, w, r) {
		return
	}

	t.HandleSuccess(w, res)

}

// UpdateByID invoke *JSONController.UpdateByID using the request body as a json payload.
// UpdateByID Decodes reqBody as json to invoke *Controller.UpdateByID.
// Other parameters are passed straight
// UpdateByID ...
// @route /{id}
// @methods PUT,POST
func (t *HTTPController) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var urlID int
	tempurlID, err := strconv.Atoi(t.dataer.Make(w, r).Get("url", "id"))
	if t.HandleError(err, w, r) {
		return
	}
	urlID = tempurlID
	reqBody := r.Body

	res, err := t.embed.UpdateByID(urlID, reqBody)
	if t.HandleError(err, w, r) {
		return
	}

	t.HandleSuccess(w, res)

}

// DeleteByID invoke *JSONController.DeleteByID using the request body as a json payload.
// DeleteByID Decodes reqBody as json to invoke *Controller.DeleteByID.
// Other parameters are passed straight
// DeleteByID ...
// @route /{id}
// @methods DELETE
func (t *HTTPController) DeleteByID(w http.ResponseWriter, r *http.Request) {
	var REQid int
	tempREQid, err := strconv.Atoi(t.dataer.Make(w, r).Get("req", "id"))
	if t.HandleError(err, w, r) {
		return
	}
	REQid = tempREQid

	res, err := t.embed.DeleteByID(REQid)
	if t.HandleError(err, w, r) {
		return
	}

	t.HandleSuccess(w, res)

}

// TestVars1 invoke *JSONController.TestVars1 using the request body as a json payload.
// TestVars1 Decodes reqBody as json to invoke *Controller.TestVars1.
// Other parameters are passed straight
// TestVars1 ...
func (t *HTTPController) TestVars1(w http.ResponseWriter, r *http.Request) {

	res, err := t.embed.TestVars1(w, r)
	if t.HandleError(err, w, r) {
		return
	}

	t.HandleSuccess(w, res)

}

// TestCookier invoke *JSONController.TestCookier using the request body as a json payload.
// TestCookier Decodes reqBody as json to invoke *Controller.TestCookier.
// Other parameters are passed straight
// TestCookier ...
func (t *HTTPController) TestCookier(w http.ResponseWriter, r *http.Request) {
	var c httper.Cookier
	c = t.cookier.Make(w, r)

	res, err := t.embed.TestCookier(c)
	if t.HandleError(err, w, r) {
		return
	}

	t.HandleSuccess(w, res)

}

// TestSessionner invoke *JSONController.TestSessionner using the request body as a json payload.
// TestSessionner Decodes reqBody as json to invoke *Controller.TestSessionner.
// Other parameters are passed straight
// TestSessionner ...
func (t *HTTPController) TestSessionner(w http.ResponseWriter, r *http.Request) {
	var s httper.Sessionner
	s = t.sessioner.Make(w, r)

	res, err := t.embed.TestSessionner(s)
	if t.HandleError(err, w, r) {
		return
	}

	t.HandleSuccess(w, res)

}

// TestRPCer invoke *JSONController.TestRPCer using the request body as a json payload.
// TestRPCer Decodes r as json to invoke *Controller.TestRPCer.
// TestRPCer ...
func (t *HTTPController) TestRPCer(w http.ResponseWriter, r *http.Request) {

	res, err := t.embed.TestRPCer(r)
	if t.HandleError(err, w, r) {
		return
	}

	t.HandleSuccess(w, res)

}
