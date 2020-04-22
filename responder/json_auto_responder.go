package responder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	config "github.com/asalih/http-auto-responder/c"
	"github.com/minio/minio/pkg/wildcard"
)

var rulePrefix string = "rule"
var responsePrefix string = "responses"

//JSONAutoResponder File system json auto responder
type JSONAutoResponder struct {
	FolderPath string
	Rules      map[uint64]*ruleFile
	Responses  map[uint64]*responseFile
	conf       *config.Config
}

type ruleFile struct {
	path string
	rule *Rule
}

type responseFile struct {
	path     string
	response *Response
}

//NewJSONAutoResponder Inits a DB Auto Responder
func NewJSONAutoResponder(conf *config.Config) JSONAutoResponder {
	return JSONAutoResponder{"./" + conf.JSONsFolderPath, make(map[uint64]*ruleFile), make(map[uint64]*responseFile), conf}
}

//Init auto responder
func (ar *JSONAutoResponder) Init() {
	stat, err := os.Stat(ar.conf.JSONsFolderPath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(ar.conf.JSONsFolderPath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if !stat.IsDir() {
		fmt.Println("Need a directory!")
	}

	_, serr := ioutil.ReadDir(ar.conf.JSONsFolderPath)

	if serr != nil {
		fmt.Println("Can't itarate the given directory")
	}

	ferr := filepath.Walk(ar.conf.JSONsFolderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			n := info.Name()
			if strings.HasSuffix(n, ".json") {
				if strings.HasPrefix(n, rulePrefix) {
					var rule Rule
					content, ferr := ioutil.ReadFile(path)
					if ferr != nil {
						return ferr
					}

					err := json.Unmarshal(content, &rule)
					if err != nil {
						return err
					}

					ar.Rules[rule.ID] = &ruleFile{path: path, rule: &rule}
				} else if strings.HasPrefix(n, responsePrefix) {
					var response Response
					content, ferr := ioutil.ReadFile(path)
					if ferr != nil {
						return ferr
					}

					err := json.Unmarshal(content, &response)
					if err != nil {
						return err
					}
					ar.Responses[response.ID] = &responseFile{path: path, response: &response}
				}
			}

			return nil
		})

	if ferr != nil {
		fmt.Println(ferr)
	}
}

//AddOrUpdateRule adds or updates given rule
func (ar *JSONAutoResponder) AddOrUpdateRule(rule *Rule) {
	if rule.ID == 0 {
		rule.ID = uint64(time.Now().Unix())
	}

	buf, err := json.Marshal(rule)
	if err != nil {
		return
	}

	path := path.Join(ar.FolderPath, rulePrefix+"_"+strconv.FormatUint(rule.ID, 10)+".json")
	ioutil.WriteFile(path, buf, os.ModePerm)

	ar.Rules[rule.ID] = &ruleFile{path: path, rule: rule}
}

//FindMatchingRule gets the rule with given url pattern and http method
func (ar *JSONAutoResponder) FindMatchingRule(urlPattern string, method string) *Rule {
	for _, rf := range ar.Rules {

		if !rf.rule.IsActive {
			continue
		}

		if (!strings.Contains(urlPattern, rf.rule.URLPattern) &&
			!wildcard.Match(rf.rule.URLPattern, urlPattern)) ||
			!wildcard.Match(rf.rule.Method, method) {
			continue
		}

		rf.rule.Response = ar.GetResponse(rf.rule.ResponseID)

		return rf.rule
	}

	return nil

}

//GetRule gets the rule with given id
func (ar *JSONAutoResponder) GetRule(id uint64) *Rule {
	rf := ar.Rules[id]

	if rf != nil {
		return rf.rule
	}

	return nil
}

//GetRules gets the rules with given url pattern and http method
func (ar *JSONAutoResponder) GetRules() []*Rule {
	values := []*Rule{}
	for _, value := range ar.Rules {
		values = append(values, value.rule)
	}

	return values
}

//RemoveRule removes the rule
func (ar *JSONAutoResponder) RemoveRule(id uint64) {
	rf := ar.Rules[id]
	os.RemoveAll(rf.path)
	delete(ar.Rules, id)
}

//AddOrUpdateResponse adds or updates given rule
func (ar *JSONAutoResponder) AddOrUpdateResponse(response *Response) {
	if response.ID == 0 {
		response.ID = uint64(time.Now().Unix())
	}

	buf, err := json.Marshal(response)
	if err != nil {
		return
	}

	path := path.Join(ar.FolderPath, responsePrefix+"_"+strconv.FormatUint(response.ID, 10)+".json")
	ioutil.WriteFile(path, buf, os.ModePerm)

	ar.Responses[response.ID] = &responseFile{path: path, response: response}
}

//GetResponse gets the response with given id
func (ar *JSONAutoResponder) GetResponse(id uint64) *Response {
	rf := ar.Responses[id]

	if rf != nil {
		return rf.response
	}

	return nil
}

//GetResponses gets the response slice
func (ar *JSONAutoResponder) GetResponses() []*Response {
	values := []*Response{}
	for _, value := range ar.Responses {
		values = append(values, value.response)
	}

	return values
}

//RemoveResponse removes the response with given id
func (ar *JSONAutoResponder) RemoveResponse(id uint64) {
	rf := ar.Responses[id]
	os.RemoveAll(rf.path)
	delete(ar.Responses, id)
}
