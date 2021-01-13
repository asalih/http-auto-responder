package responder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asalih/http-auto-responder/utils"
)

var rulePrefix string = "rule"
var responsePrefix string = "responses"

//JSONAutoResponder File system json auto responder
type JSONAutoResponder struct {
	FolderPath string
	Rules      map[uint64]*ruleFile
	Responses  map[uint64]*responseFile
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
func NewJSONAutoResponder() JSONAutoResponder {
	return JSONAutoResponder{"./" + utils.Configuration.JSONsFolderPath, make(map[uint64]*ruleFile), make(map[uint64]*responseFile)}
}

//Init auto responder
func (ar *JSONAutoResponder) Init() {
	stat, err := os.Stat(utils.Configuration.JSONsFolderPath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(utils.Configuration.JSONsFolderPath, 0755)
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

	_, serr := ioutil.ReadDir(utils.Configuration.JSONsFolderPath)

	if serr != nil {
		fmt.Println("Can't itarate the given directory")
	}

	ferr := filepath.Walk(utils.Configuration.JSONsFolderPath,
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

		if !rf.rule.IsActive || !utils.WildcardMatch(rf.rule.Method, method) {
			continue
		}

		mType := utils.GetMatchType(rf.rule.MatchType)
		if (mType == utils.EXACT && urlPattern != rf.rule.URLPattern) ||
			(mType == utils.WILDCARD && !utils.WildcardMatch(rf.rule.URLPattern, urlPattern)) ||
			(mType == utils.CONTAINS && !strings.Contains(urlPattern, rf.rule.URLPattern)) ||
			(mType == utils.NOT && strings.Contains(urlPattern, rf.rule.URLPattern)) {
			continue
		} else if mType == utils.REGEX {
			m, err := regexp.MatchString(rf.rule.URLPattern, urlPattern)

			if !m || err != nil {
				continue
			}
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
func (ar *JSONAutoResponder) GetRules(skip int) []*Rule {
	values := []*Rule{}

	if len(ar.Rules) <= skip {
		return values
	}

	i := 0
	for _, value := range ar.Rules {
		if i > skip && i <= skip+takeSize {
			values = append(values, value.rule)
		}
		i++
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
func (ar *JSONAutoResponder) GetResponses(skip int) []*Response {
	values := []*Response{}

	if len(ar.Responses) <= skip {
		return values
	}

	i := 0
	for _, value := range ar.Responses {
		if i > skip && i <= skip+takeSize {
			values = append(values, value.response)
		}
		i++
	}

	return values
}

//RemoveResponse removes the response with given id
func (ar *JSONAutoResponder) RemoveResponse(id uint64) {
	rf := ar.Responses[id]
	os.RemoveAll(rf.path)
	delete(ar.Responses, id)
}
