package responder

import (
	"encoding/binary"
	"encoding/json"
	"log"

	bolt "github.com/etcd-io/bbolt"
	"github.com/minio/minio/pkg/wildcard"
)

var ruleBucketName []byte = []byte("rules")
var responseBucketName []byte = []byte("responses")

//AutoResponder ...
type AutoResponder struct {
	DBPath        string
	ListeningPort int
	db            *bolt.DB
}

//NewAutoResponder Inits an Auto Responder
func NewAutoResponder(dbPath string, listeningPort int) *AutoResponder {
	autoResponder := &AutoResponder{dbPath, listeningPort, nil}

	autoResponder.init()
	return autoResponder
}

func (ar *AutoResponder) init() {
	db, err := bolt.Open(ar.DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	ar.db = db
}

//AddOrUpdateRule adds or updates given rule
func (ar *AutoResponder) AddOrUpdateRule(rule *Rule) {
	ar.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(ruleBucketName)
		b := tx.Bucket(ruleBucketName)

		buf, err := json.Marshal(rule)
		if err != nil {
			return err
		}

		return b.Put([]byte(rule.URLPattern+rule.Method), buf)
	})
}

//GetRule gets the rule with given url pattern and http method
func (ar *AutoResponder) GetRule(urlPattern string, method string) *Rule {
	var extractedRule *Rule
	ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ruleBucketName)

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

			if !wildcard.Match(rule.URLPattern, urlPattern) &&
				!wildcard.Match(rule.Method, method) {
				continue
			}

			rule.Response = ar.GetResponse(rule.ResponseID)

			extractedRule = &rule
			return nil
		}

		return nil
	})

	return extractedRule
}

//GetRules gets the rules with given url pattern and http method
func (ar *AutoResponder) GetRules() []*Rule {
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
func (ar *AutoResponder) RemoveRule(urlPattern string, method string) {
	ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ruleBucketName)

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

			if !wildcard.Match(rule.URLPattern, urlPattern) &&
				!wildcard.Match(rule.Method, method) {
				continue
			}

			return c.Delete()
		}

		return nil
	})
}

//AddOrUpdateResponse adds or updates given rule
func (ar *AutoResponder) AddOrUpdateResponse(response *Response) {
	ar.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(responseBucketName)
		b := tx.Bucket(responseBucketName)

		buf, err := json.Marshal(response)
		if err != nil {
			return err
		}

		return b.Put(itob(response.ID), buf)
	})
}

//GetResponse gets the response with given id
func (ar *AutoResponder) GetResponse(id int) *Response {
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
func (ar *AutoResponder) GetResponses() []*Response {
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
func (ar *AutoResponder) RemoveResponse(id int) {
	ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(responseBucketName)
		if b == nil {
			return nil
		}

		return b.Delete(itob(id))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
