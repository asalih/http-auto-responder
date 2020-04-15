package responder

import "net/http"

//Rule ...
type Rule struct {
	URLPattern string
	Method     string
	ResponseID int
	Response   *Response
}

func (r *Rule) Write(w http.ResponseWriter) {

	for k, v := range r.Response.Headers {
		w.Header().Set(k, v)
	}

	w.Write([]byte(r.Response.Body))
}
