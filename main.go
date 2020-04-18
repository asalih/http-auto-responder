package main

import (
	"net/http"
	"strconv"

	"github.com/asalih/http-auto-responder/responder"
	"github.com/asalih/http-auto-responder/web"
)

func main() {

	autoResponder := responder.NewAutoResponder("./ar.db", 80)

	http.ListenAndServe(":"+strconv.Itoa(autoResponder.ListeningPort), web.NewAutoResponderHTTPHandler(autoResponder))
}
