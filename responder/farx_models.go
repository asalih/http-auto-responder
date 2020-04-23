package responder

import (
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/asalih/http-auto-responder/utils"
)

//FarxFileXML Farx response auto responder xml definition
type FarxFileXML struct {
	XMLName xml.Name `xml:"AutoResponder"`
	States  []*State `xml:"State"`
}

//State Farx response state xml definition
type State struct {
	ResponseRules []*ResponseRule `xml:"ResponseRule"`
	Enabled       bool            `xml:"Enabled,attr"`
}

//ResponseRule Farx response rule xml definition
type ResponseRule struct {
	Match   string `xml:"Match,attr"`
	Headers string `xml:"Headers,attr"`
	Body    string `xml:"Body,attr"`
}

//ReadFarxFile Farx file reader
func ReadFarxFile(path string) (*FarxFileXML, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var farxAutoResponder FarxFileXML
	xerr := xml.Unmarshal(file, &farxAutoResponder)

	if xerr != nil {
		return nil, xerr
	}

	if farxAutoResponder.States == nil || len(farxAutoResponder.States) == 0 {
		return nil, nil
	}

	return &farxAutoResponder, nil
}

//MapToRule Transforms ResponseRule to Rule
func (xmlRule *ResponseRule) MapToRule() *Rule {

	idx := strings.Index(xmlRule.Match, ":")

	matchType := utils.CONTAINS
	match := xmlRule.Match

	if idx > -1 {
		rs := []rune(xmlRule.Match)
		matchType = utils.GetMatchType(string(rs[0:idx]))
		match = string(rs[idx+1:])
	}

	return &Rule{IsActive: true, URLPattern: match, MatchType: matchType, Method: "*"}
}

//MapToResponse Transforms ResponseRule to Response
func (xmlRule *ResponseRule) MapToResponse() *Response {
	headersRaw, _ := base64.StdEncoding.DecodeString(xmlRule.Headers)
	lines := strings.Split(string(headersRaw), "\n")

	status, _ := strconv.Atoi(strings.Split(lines[0], " ")[1])
	var headers []*Headers

	for i := 1; i < len(lines); i++ {
		l := lines[i]
		idx := strings.Index(l, ":")
		if idx > -1 {
			rs := []rune(l)

			headers = append(headers, &Headers{Key: string(rs[0:idx]), Value: strings.Trim(string(rs[idx+1:]), " ")})
		}
	}

	bodyRaw, _ := base64.StdEncoding.DecodeString(xmlRule.Body)

	return &Response{Headers: headers, Body: string(bodyRaw), StatusCode: status}
}
