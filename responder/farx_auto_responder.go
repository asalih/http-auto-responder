package responder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/asalih/http-auto-responder/utils"
)

//FarxAutoResponder File system json auto responder
type FarxAutoResponder struct {
	FolderPath string
	Rules      []*ruleFile
}

//NewFarxAutoResponder Inits a DB Auto Responder
func NewFarxAutoResponder() FarxAutoResponder {
	return FarxAutoResponder{"./" + utils.Configuration.FarxFilesFolderPath, []*ruleFile{}}
}

//Init auto responder
func (ar *FarxAutoResponder) Init() {
	stat, err := os.Stat(utils.Configuration.FarxFilesFolderPath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(utils.Configuration.FarxFilesFolderPath, 0755)
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

	_, serr := ioutil.ReadDir(utils.Configuration.FarxFilesFolderPath)

	if serr != nil {
		fmt.Println("Can't itarate the given directory")
	}

	ferr := filepath.Walk(utils.Configuration.FarxFilesFolderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			n := info.Name()
			if strings.HasSuffix(n, ".farx") {
				rules, err := ar.loadFarxRule(path)

				if err != nil {
					return err
				}

				if rules != nil {
					ar.Rules = append(ar.Rules, rules...)
				}

			}

			return nil
		})

	if ferr != nil {
		fmt.Println(ferr)
	}
}

//AddOrUpdateRule adds or updates given rule
func (ar *FarxAutoResponder) AddOrUpdateRule(rule *Rule) {
	//NOOP FARX Auto responder no ability to rule.
}

//FindMatchingRule gets the rule with given url pattern and http method
func (ar *FarxAutoResponder) FindMatchingRule(urlPattern string, method string) *Rule {
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

		return rf.rule
	}

	return nil

}

//GetRule gets the rule with given id
func (ar *FarxAutoResponder) GetRule(id uint64) *Rule {
	//NOOP FARX Auto responder no ability to rule.
	return nil
}

//GetRules gets the rules with given url pattern and http method
func (ar *FarxAutoResponder) GetRules() []*Rule {
	//NOOP FARX Auto responder no ability to rule.
	return nil
}

//RemoveRule removes the rule
func (ar *FarxAutoResponder) RemoveRule(id uint64) {
	//NOOP FARX Auto responder no ability to rule.
}

//AddOrUpdateResponse adds or updates given rule
func (ar *FarxAutoResponder) AddOrUpdateResponse(response *Response) {
	//NOOP FARX Auto responder no ability to rule.
}

//GetResponse gets the response with given id
func (ar *FarxAutoResponder) GetResponse(id uint64) *Response {
	//NOOP FARX Auto responder no ability to rule.
	return nil
}

//GetResponses gets the response slice
func (ar *FarxAutoResponder) GetResponses() []*Response {
	//NOOP FARX Auto responder no ability to rule.
	return nil
}

//RemoveResponse removes the response with given id
func (ar *FarxAutoResponder) RemoveResponse(id uint64) {
	//NOOP FARX Auto responder no ability to rule.
}

//Reload updates given farx
func (ar *FarxAutoResponder) Reload(farxPath string) (bool, error) {
	farxPath = strings.ReplaceAll(farxPath, "/", "\\")
	xmlRule, err := ar.loadFarxRule(farxPath)

	if xmlRule == nil || err != nil {
		return false, err
	}

	newSlice := []*ruleFile{}
	for _, rf := range ar.Rules {

		if strings.Index(farxPath, rf.path) > -1 {
			continue
		}
		newSlice = append(newSlice, rf)
	}

	ar.Rules = append(newSlice, xmlRule...)

	return true, nil
}

//LoadFarxRule Loads a farx file and returns rules
func (ar *FarxAutoResponder) loadFarxRule(farxPath string) ([]*ruleFile, error) {
	farx, err := ReadFarxFile(farxPath)
	if err != nil {
		return nil, err
	}

	result := []*ruleFile{}
	for _, s := range farx.States {
		if !s.Enabled {
			continue
		}

		for _, r := range s.ResponseRules {
			if r.Headers == "" && r.Action == "" {
				continue
			}

			response := r.MapToResponse()

			rule := r.MapToRule()
			rule.Response = response

			result = append(result, &ruleFile{farxPath, rule})
		}
	}

	return result, nil
}
