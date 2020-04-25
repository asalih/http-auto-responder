package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/asalih/http-auto-responder/responder"
	"github.com/asalih/http-auto-responder/utils"
)

func init() {
	HTTPHandlerFuncs["/http-auto-responder/save-rule"] = func(w http.ResponseWriter, r *http.Request, ar responder.AutoResponder) {
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

	HTTPHandlerFuncs["/http-auto-responder/get-rules"] = func(w http.ResponseWriter, r *http.Request, ar responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		rules := ar.GetRules()

		encoder := json.NewEncoder(w)
		encoder.Encode(rules)
	}

	HTTPHandlerFuncs["/http-auto-responder/get-rule"] = func(w http.ResponseWriter, r *http.Request, ar responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		idq := r.URL.Query().Get("id")
		id, err := strconv.ParseUint(idq, 0, 64)

		if err != nil || idq == "" {
			http.NotFound(w, r)
			return
		}

		rule := ar.GetRule(id)

		encoder := json.NewEncoder(w)
		encoder.Encode(rule)
	}

	HTTPHandlerFuncs["/http-auto-responder/remove-rule"] = func(w http.ResponseWriter, r *http.Request, ar responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		idq := r.URL.Query().Get("id")
		id, err := strconv.ParseUint(idq, 0, 64)

		if err != nil || idq == "" {
			http.NotFound(w, r)
			return
		}

		ar.RemoveRule(id)

		w.Write([]byte(`{"status": "OK"}`))
	}

	HTTPHandlerFuncs["/http-auto-responder/reload"] = func(w http.ResponseWriter, r *http.Request, ar responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		if utils.Configuration.FarxFilesFolderPath == "" {
			w.Write([]byte(`{"status": "NA"}`))
			return
		}

		pq := r.URL.Query().Get("path")

		if path.Ext(pq) != ".farx" {
			w.Write([]byte(`{"status": "Unexpected extension"}`))
			return
		}

		data := ar.(*responder.FarxAutoResponder)
		success, err := data.Reload(utils.Configuration.FarxFilesFolderPath + "/" + pq)

		if err != nil {
			fmt.Println(err)
		}

		if success {
			w.Write([]byte(`{"status": "OK"}`))
		} else {
			w.Write([]byte(`{"status": "Err"}`))
		}
	}
}
