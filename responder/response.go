package responder

//Response ...
type Response struct {
	ID         uint64     `json:"id"`
	Label      string     `json:"label"`
	StatusCode int        `json:"statusCode"`
	Headers    []*Headers `json:"headers"`
	Body       string     `json:"body"`
}

//Headers response headers key value
type Headers struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
