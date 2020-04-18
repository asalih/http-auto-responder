package responder

//Response ...
type Response struct {
	ID      int               `json:"id"`
	Label   string            `json:"label"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
