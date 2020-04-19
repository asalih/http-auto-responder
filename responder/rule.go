package responder

import "net/http"

//Rule ...
type Rule struct {
	ID         uint64    `json:"id"`
	URLPattern string    `json:"urlPattern"`
	Method     string    `json:"method"`
	ResponseID uint64    `json:"responseID"`
	IsActive   bool      `json:"isActive"`
	Response   *Response `json:"-"`
}

func (r *Rule) Write(w http.ResponseWriter) {

	for _, v := range r.Response.Headers {
		w.Header().Set(v.Key, v.Value)
	}

	w.WriteHeader(r.Response.StatusCode)
	w.Write([]byte(r.Response.Body))
}
