package responder

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/asalih/http-auto-responder/utils"
	bolt "github.com/etcd-io/bbolt"
	"github.com/minio/minio/pkg/wildcard"
)

var ruleBucketName []byte = []byte("rules")
var responseBucketName []byte = []byte("responses")

//DBAutoResponder DB auto responder
type DBAutoResponder struct {
	DBPath string
	db     *bolt.DB
}

//NewDBAutoResponder Inits a DB Auto Responder
func NewDBAutoResponder() DBAutoResponder {
	return DBAutoResponder{"./" + utils.Configuration.DatabaseName, nil}
}

//Init auto responder
func (ar *DBAutoResponder) Init() {
	db, err := bolt.Open(ar.DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	ar.db = db
}

//AddOrUpdateRule adds or updates given rule
func (ar *DBAutoResponder) AddOrUpdateRule(rule *Rule) {
	ar.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(ruleBucketName)
		b := tx.Bucket(ruleBucketName)

		if rule.ID == 0 {
			id, _ := b.NextSequence()
			rule.ID = id
		}

		buf, err := json.Marshal(rule)
		if err != nil {
			return err
		}

		return b.Put(itob(rule.ID), buf)
	})
}

//FindMatchingRule gets the rule with given url pattern and http method
func (ar *DBAutoResponder) FindMatchingRule(urlPattern string, method string) *Rule {
	var extractedRule *Rule
	ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ruleBucketName)

		if b == nil {
			return nil
		}

		c := b.Cursor()

		if c == nil {
			return nil
		}

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rule Rule
			err := json.Unmarshal(v, &rule)

			if err != nil {
				continue
			}

			if !rule.IsActive || !wildcard.Match(rule.Method, method) {
				continue
			}

			mType := utils.GetMatchType(rule.MatchType)
			if (mType == utils.EXACT && urlPattern != strings.ToLower(rule.URLPattern)) ||
				(mType == utils.WILDCARD && !wildcard.Match(rule.URLPattern, urlPattern)) ||
				(mType == utils.CONTAINS && !strings.Contains(urlPattern, rule.URLPattern)) ||
				(mType == utils.NOT && strings.Contains(urlPattern, rule.URLPattern)) {
				continue
			} else if mType == utils.REGEX {
				m, err := regexp.MatchString(rule.URLPattern, urlPattern)

				if !m || err != nil {
					continue
				}
			}

			rule.Response = ar.GetResponse(rule.ResponseID)

			extractedRule = &rule
			return nil
		}

		return nil
	})

	return extractedRule
}

//GetRule gets the rule with given id
func (ar *DBAutoResponder) GetRule(id uint64) *Rule {
	var rule Rule
	ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ruleBucketName)
		if b == nil {
			return nil
		}

		v := b.Get(itob(id))

		err := json.Unmarshal(v, &rule)

		if err != nil {
			return err
		}

		return nil
	})

	return &rule
}

//GetRules gets the rules with given url pattern and http method
func (ar *DBAutoResponder) GetRules() []*Rule {
	var rules []*Rule
	ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ruleBucketName)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		if c == nil {
			return nil
		}

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var rule Rule
			err := json.Unmarshal(v, &rule)
			if err != nil {
				continue
			}
			rules = append(rules, &rule)
		}

		return nil
	})

	return rules
}

//RemoveRule removes the rule
func (ar *DBAutoResponder) RemoveRule(id uint64) {
	ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ruleBucketName)

		if b == nil {
			return nil
		}

		return b.Delete(itob(id))
	})
}

//AddOrUpdateResponse adds or updates given rule
func (ar *DBAutoResponder) AddOrUpdateResponse(response *Response) {
	ar.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(responseBucketName)
		b := tx.Bucket(responseBucketName)

		if response.ID == 0 {
			id, _ := b.NextSequence()
			response.ID = id
		}

		buf, err := json.Marshal(response)
		if err != nil {
			return err
		}

		return b.Put(itob(response.ID), buf)
	})
}

//GetResponse gets the response with given id
func (ar *DBAutoResponder) GetResponse(id uint64) *Response {
	var response Response
	ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(responseBucketName)
		if b == nil {
			return nil
		}

		v := b.Get(itob(id))

		err := json.Unmarshal(v, &response)

		if err != nil {
			return err
		}

		return nil
	})

	return &response
}

//GetResponses gets the response slice
func (ar *DBAutoResponder) GetResponses() []*Response {
	var responses []*Response
	ar.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(responseBucketName)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		if c == nil {
			return nil
		}

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var response Response
			err := json.Unmarshal(v, &response)

			if err != nil {
				continue
			}
			responses = append(responses, &response)
		}

		return nil
	})

	return responses
}

//RemoveResponse removes the response with given id
func (ar *DBAutoResponder) RemoveResponse(id uint64) {
	ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(responseBucketName)
		if b == nil {
			return nil
		}

		return b.Delete(itob(id))
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
