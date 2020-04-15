package responder

//Response ...
type Response struct {
	ID      int
	Label   string
	Headers map[string]string
	Body    string
}
