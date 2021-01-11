package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asalih/http-auto-responder/parser"
	"github.com/asalih/http-auto-responder/responder"
	"github.com/asalih/http-auto-responder/utils"
	"github.com/asalih/http-auto-responder/web"
)

func main() {

	utils.InitConfig()

	autoResponder := responder.NewAutoResponder(&utils.Configuration)

	if autoResponder == nil {
		fmt.Println("No Auto Responder provided. Provide DB Name or Folder path for auto responder data handling!")
		return
	}

	parser.Migrate(autoResponder)

	http.ListenAndServe(":"+strconv.Itoa(utils.Configuration.Port), web.NewAutoResponderHTTPHandler(autoResponder))
}
