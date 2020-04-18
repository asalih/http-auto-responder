package web

import (
	"encoding/json"
	"net/http"

	"github.com/asalih/http-auto-responder/responder"
)

func init() {
	HTTPHandlerFuncs["/save-rule"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		var rule responder.Rule
		err := decoder.Decode(&rule)
		if err != nil {
			panic(err)
		}

		rule.IsActive = true
		ar.AddOrUpdateRule(&rule)

		w.Write([]byte(`{"status": "OK"}`))
	}

	HTTPHandlerFuncs["/get-rules"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		rules := ar.GetRules()

		encoder := json.NewEncoder(w)
		encoder.Encode(rules)
	}

	HTTPHandlerFuncs["/get-rule"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		pattern := r.URL.Query().Get("pattern")
		method := r.URL.Query().Get("method")

		if method == "" || pattern == "" {
			http.NotFound(w, r)
			return
		}

		rule := ar.GetRule(pattern, method)

		encoder := json.NewEncoder(w)
		encoder.Encode(rule)
	}

	HTTPHandlerFuncs["/remove-rule"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		pattern := r.URL.Query().Get("urlPattern")
		method := r.URL.Query().Get("method")

		if method == "" || pattern == "" {
			http.NotFound(w, r)
			return
		}

		ar.RemoveRule(pattern, method)

		w.Write([]byte(`{"status": "OK"}`))
	}
}
