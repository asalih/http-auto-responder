package web

import (
	"net/http"

	"github.com/asalih/http-auto-responder/responder"
)

//HTTPHandlerFuncs global variable for http handlers
var HTTPHandlerFuncs map[string]func(http.ResponseWriter, *http.Request, *responder.AutoResponder) = make(map[string]func(http.ResponseWriter, *http.Request, *responder.AutoResponder))

//AutoResponderHTTPHandler handler for global http requests
type AutoResponderHTTPHandler struct {
	AutoResponder *responder.AutoResponder
}

//NewAutoResponderHTTPHandler Creates an http wrapper for AutoResponder
func NewAutoResponderHTTPHandler(refAutoResponder *responder.AutoResponder) *AutoResponderHTTPHandler {
	return &AutoResponderHTTPHandler{refAutoResponder}
}

func (handler *AutoResponderHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	hfn := HTTPHandlerFuncs[r.URL.Path]

	if hfn != nil {
		hfn(w, r, handler.AutoResponder)
		return
	}

	rule := handler.AutoResponder.GetRule(r.URL.String(), r.Method)

	if rule == nil || !rule.IsActive {
		http.NotFound(w, r)
		return
	}

	rule.Write(w)
}
