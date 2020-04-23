package parser

import (
	"strconv"

	"github.com/asalih/http-auto-responder/responder"
)

//FarxParser Implements a parsing struct for farx parsing
type FarxParser struct {
	FarxFilePath  string
	OrigFileName  string
	AutoResponder responder.AutoResponder
}

//Handle Starts file import process
func (parser *FarxParser) Handle() error {
	farx, err := responder.ReadFarxFile(parser.FarxFilePath)

	if farx == nil || err != nil {
		return err
	}

	i := 0
	for _, s := range farx.States {
		if !s.Enabled {
			continue
		}

		for _, r := range s.ResponseRules {
			i++

			if r.Headers == "" {
				continue
			}

			response := r.MapToResponse()
			response.Label = parser.OrigFileName + "_" + strconv.Itoa(i)

			parser.AutoResponder.AddOrUpdateResponse(response)

			rule := r.MapToRule()
			rule.ResponseID = response.ID

			parser.AutoResponder.AddOrUpdateRule(rule)
		}
	}

	return nil
}
