package utils

import (
	"net/http"
	"net/url"
)

//HTTPRequestModel proxy model for http.Request while executing response template
type HTTPRequestModel struct {
	URL     string
	Method  string
	Headers map[string]string
	Cookies []*http.Cookie
	Query   map[string]string
	Form    map[string]string
}

func mapHeaders(header http.Header) map[string]string {
	headers := make(map[string]string)

	for name, values := range header {
		hval := ""
		for _, value := range values {
			hval += value
		}
		headers[name] = hval
	}
	return headers
}

func mapQueryOrForm(query url.Values) map[string]string {
	q := make(map[string]string)

	for name, values := range query {
		hval := ""
		for _, value := range values {
			hval += value
		}
		q[name] = hval
	}
	return q
}

//MapToHTTPRequestModel Creates a HTTPRequestModel
func MapToHTTPRequestModel(req *http.Request) *HTTPRequestModel {
	req.ParseForm()

	return &HTTPRequestModel{req.URL.String(), req.Method, mapHeaders(req.Header), req.Cookies(), mapQueryOrForm(req.URL.Query()), mapQueryOrForm(req.Form)}

}
