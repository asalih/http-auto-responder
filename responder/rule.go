package responder

import "net/http"

//Rule ...
type Rule struct {
	URLPattern string    `json:"urlPattern"`
	Method     string    `json:"method"`
	ResponseID int       `json:"responseID"`
	IsActive   bool      `json:"isActive"`
	Response   *Response `json:"-"`
}

func (r *Rule) Write(w http.ResponseWriter) {

	for k, v := range r.Response.Headers {
		w.Header().Set(k, v)
	}

	w.Write([]byte(r.Response.Body))
}
