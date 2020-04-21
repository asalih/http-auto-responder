package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/asalih/http-auto-responder/responder"
)

func init() {
	HTTPHandlerFuncs["/http-auto-responder"] = func(w http.ResponseWriter, r *http.Request, ar responder.AutoResponder) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "Unable to load template")
		}

		t.Execute(w, nil)
	}
}
