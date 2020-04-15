package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/asalih/http-auto-responder/responder"
)

//HTTPHandler handler for global http requests
type HTTPHandler struct {
	AutoResponder *responder.AutoResponder
}

func (handler *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/http-auto-responder" {
		handler.handleUI(w, r)
		return
	} else if r.URL.Path == "/save-rule" {
		handler.saveRule(w, r)
	}

	rule := handler.AutoResponder.GetRule(r.URL.String(), r.Method)

	if rule == nil {
		return
	}

	rule.Write(w)
}

func (handler *HTTPHandler) handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}

	t.Execute(w, nil)
}

func (handler *HTTPHandler) saveRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var rule responder.Rule
	err := decoder.Decode(&rule)
	if err != nil {
		panic(err)
	}

	handler.AutoResponder.AddOrUpdateRule(&rule)

	w.Write([]byte("OK"))
}
