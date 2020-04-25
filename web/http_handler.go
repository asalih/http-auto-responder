package web

import (
	"net/http"
	"time"

	"github.com/asalih/http-auto-responder/responder"
	"github.com/asalih/http-auto-responder/utils"
)

//HTTPHandlerFuncs global variable for http handlers
var HTTPHandlerFuncs map[string]func(http.ResponseWriter, *http.Request, responder.AutoResponder) = make(map[string]func(http.ResponseWriter, *http.Request, responder.AutoResponder))

//AutoResponderHTTPHandler handler for global http requests
type AutoResponderHTTPHandler struct {
	AutoResponder responder.AutoResponder
}

//NewAutoResponderHTTPHandler Creates an http wrapper for AutoResponder
func NewAutoResponderHTTPHandler(refAutoResponder responder.AutoResponder) *AutoResponderHTTPHandler {
	return &AutoResponderHTTPHandler{refAutoResponder}
}

//ServeHTTP HTTP Serve handler
func (handler *AutoResponderHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	p := r.URL.Path
	hfn := HTTPHandlerFuncs[p]

	if hfn != nil {
		if utils.Configuration.FarxFilesFolderPath != "" && p != "/http-auto-responder/reload" {
			w.Write([]byte("FARX File System Auto Responder activated. NO UI management available for FARX files."))
			return
		}

		hfn(w, r, handler.AutoResponder)
		return
	}

	url := r.Host

	if r.TLS != nil {
		url = "https://" + url
	} else {
		url = "http://" + url
	}

	url += r.URL.String()

	rule := handler.AutoResponder.FindMatchingRule(url, r.Method)

	if rule == nil {
		http.NotFound(w, r)
		return
	}

	if rule.Latency > 0 {
		time.Sleep(time.Millisecond * time.Duration(rule.Latency))
	}

	rule.Write(w, r)
}
