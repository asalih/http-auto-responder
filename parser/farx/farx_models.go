package parser

//AutoResponder Farx response auto responder xml definition
type AutoResponder struct {
	States []*State `xml:"State"`
}

//State Farx response state xml definition
type State struct {
	ResponseRules []*ResponseRule `xml:"ResponseRule"`
}

//ResponseRule Farx response rule xml definition
type ResponseRule struct {
	Match   string `xml:"Match,attr"`
	Headers string `xml:"Headers,attr"`
	Body    string `xml:"Body,attr"`
}
