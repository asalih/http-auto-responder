package main

import (
	"fmt"
	"net/http"
	"strconv"

	config "github.com/asalih/http-auto-responder/c"
	"github.com/asalih/http-auto-responder/responder"
	"github.com/asalih/http-auto-responder/web"
)

func main() {

	config.InitConfig()

	autoResponder := responder.NewAutoResponder(&config.Configuration)

	if autoResponder == nil {
		fmt.Println("No Auto Responder provided. Provide DB Name or Folder path for auto responder data handling!")
		return
	}

	http.ListenAndServe(":"+strconv.Itoa(config.Configuration.Port), web.NewAutoResponderHTTPHandler(autoResponder))
}
