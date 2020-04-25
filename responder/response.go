package responder

import (
	"strconv"
	"strings"
)

//Response ...
type Response struct {
	ID         uint64     `json:"id"`
	Label      string     `json:"label"`
	StatusCode int        `json:"statusCode"`
	Headers    []*Headers `json:"headers"`
	Body       string     `json:"body"`
}

//Headers response headers key value
type Headers struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//NewResponseFromString parse raw string content
func NewResponseFromString(content string) *Response {
	lines := strings.Split(content, "\n")

	status, _ := strconv.Atoi(strings.Split(lines[0], " ")[1])
	var headers []*Headers
	body := ""

	bodyStart := false
	for i := 1; i < len(lines); i++ {
		l := lines[i]
		if !bodyStart {
			if len(l) == 1 {
				bodyStart = true
				continue
			}
			idx := strings.Index(l, ":")

			if idx > -1 {
				rs := []rune(l)

				headers = append(headers, &Headers{Key: string(rs[0:idx]), Value: strings.Trim(string(rs[idx+1:]), " ")})
			}
		} else {
			body += l
		}
	}

	return &Response{Headers: headers, Body: body, StatusCode: status}
}
