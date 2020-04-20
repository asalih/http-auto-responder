package main

import (
	"net/http"
	"strconv"

	config "github.com/asalih/http-auto-responder/c"
	"github.com/asalih/http-auto-responder/responder"
	"github.com/asalih/http-auto-responder/web"
)

func main() {

	config.InitConfig()

	autoResponder := responder.NewAutoResponder("./"+config.Configuration.DatabaseName, config.Configuration.Port)

	http.ListenAndServe(":"+strconv.Itoa(autoResponder.ListeningPort), web.NewAutoResponderHTTPHandler(autoResponder))
}
