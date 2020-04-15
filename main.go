package main

import (
	"net/http"
	"strconv"

	"github.com/asalih/http-auto-responder/responder"
)

func main() {

	autoResponder := responder.NewAutoResponder("./rr.db", 80)

	http.ListenAndServe(":"+strconv.Itoa(autoResponder.ListeningPort), &HTTPHandler{autoResponder})
}
