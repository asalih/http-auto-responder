package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asalih/http-auto-responder/responder"
)

func init() {
	HTTPHandlerFuncs["/http-auto-responder/save-response"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		var response responder.Response
		err := decoder.Decode(&response)
		if err != nil {
			panic(err)
		}

		ar.AddOrUpdateResponse(&response)

		w.Write([]byte(`{"status": "OK"}`))
	}

	HTTPHandlerFuncs["/http-auto-responder/get-responses"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		responses := ar.GetResponses()

		encoder := json.NewEncoder(w)
		encoder.Encode(responses)
	}

	HTTPHandlerFuncs["/http-auto-responder/get-response"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		idq := r.URL.Query().Get("id")
		id, err := strconv.ParseUint(idq, 0, 64)

		if err != nil || idq == "" {
			http.NotFound(w, r)
			return
		}

		response := ar.GetResponse(id)

		encoder := json.NewEncoder(w)
		encoder.Encode(response)
	}

	HTTPHandlerFuncs["/http-auto-responder/remove-response"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
		w.Header().Set("Content-Type", "application/json")

		idq := r.URL.Query().Get("id")
		id, err := strconv.ParseUint(idq, 0, 64)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		ar.RemoveResponse(id)

		w.Write([]byte(`{"status": "OK"}`))
	}
}
