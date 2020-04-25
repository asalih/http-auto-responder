package responder

import (
	"net/http"
	"text/template"

	"github.com/asalih/http-auto-responder/utils"
)

//Rule ...
type Rule struct {
	ID         uint64    `json:"id"`
	URLPattern string    `json:"urlPattern"`
	MatchType  string    `json:"matchType"`
	Method     string    `json:"method"`
	ResponseID uint64    `json:"responseID"`
	Latency    int       `json:"latency"`
	IsActive   bool      `json:"isActive"`
	Response   *Response `json:"-"`
}

func (r *Rule) Write(w http.ResponseWriter, req *http.Request) {

	for _, v := range r.Response.Headers {
		if v.Key == "" || v.Value == "" {
			continue
		}

		if v.Key == "Content-Length" {
			continue
		}
		w.Header().Set(v.Key, v.Value)
	}

	w.WriteHeader(r.Response.StatusCode)

	templ, err := template.New(r.URLPattern).Parse(r.Response.Body)
	if err != nil {
		w.Write([]byte("Err on templating"))
	}

	model := utils.MapToHTTPRequestModel(req)

	templ.Execute(w, model)
}
