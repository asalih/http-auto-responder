package responder

import "github.com/asalih/http-auto-responder/utils"

const takeSize = 50

//AutoResponder implements an interface for autoresponder
type AutoResponder interface {
	//Init auto responder
	Init()
	//AddOrUpdateRule adds or updates given rule
	AddOrUpdateRule(rule *Rule)
	//FindMatchingRule gets the rule with given url pattern and http method
	FindMatchingRule(urlPattern string, method string) *Rule
	//GetRule gets the rule with given id
	GetRule(id uint64) *Rule
	//GetRules gets the rules with given url pattern and http method
	GetRules(skip int) []*Rule
	//RemoveRule removes the rule
	RemoveRule(id uint64)
	//AddOrUpdateResponse adds or updates given rule
	AddOrUpdateResponse(response *Response)
	//GetResponse gets the response with given id
	GetResponse(id uint64) *Response
	//GetResponses gets the response slice
	GetResponses(skip int) []*Response
	//RemoveResponse removes the response with given id
	RemoveResponse(id uint64)
}

//NewAutoResponder Inits an Auto Responder
func NewAutoResponder(conf *utils.Config) AutoResponder {
	if conf.DatabaseName != "" {
		dbResponder := NewDBAutoResponder()
		dbResponder.Init()

		return &dbResponder
	} else if conf.JSONsFolderPath != "" {
		jsonResponder := NewJSONAutoResponder()
		jsonResponder.Init()

		return &jsonResponder
	} else if conf.FarxFilesFolderPath != "" {
		farxResponder := NewFarxAutoResponder()
		farxResponder.Init()

		return &farxResponder
	}

	return nil
}
