package parser

import (
	"strconv"
	"strings"

	"github.com/asalih/http-auto-responder/responder"
)

//FarxParser Implements a parsing struct for farx parsing
type FarxParser struct {
	FarxFileName  string
	FarxFilePath  string
	UploadPath    string
	OrigFileName  string
	AutoResponder responder.AutoResponder
}

//Handle Starts file import process
func (parser *FarxParser) Handle() error {
	return nil
}

func (parser *FarxParser) parseRule(content string) *responder.Rule {
	firstLine := strings.Split(strings.Split(content, "\n")[0], " ")

	return &responder.Rule{IsActive: true, URLPattern: firstLine[1], Method: firstLine[0]}
}

func (parser *FarxParser) parseResponse(content string) *responder.Response {
	lines := strings.Split(content, "\n")

	status, _ := strconv.Atoi(strings.Split(lines[0], " ")[1])
	var headers []*responder.Headers
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
			rs := []rune(l)

			headers = append(headers, &responder.Headers{Key: string(rs[0:idx]), Value: strings.Trim(string(rs[idx+1:]), " ")})
		} else {
			body += l
		}
	}

	return &responder.Response{Headers: headers, Body: body, StatusCode: status}
}
