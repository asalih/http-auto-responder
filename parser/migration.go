package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asalih/http-auto-responder/responder"
)

const migrationPath = "./migration"
const uploadPath = "./tmp"

// Migrate is helper func for importing migration folder.
func Migrate(refAutoResponder responder.AutoResponder) {
	fmt.Println("Migration started.")

	filepath.Walk("./migration", func(path string, f os.FileInfo, err error) error {

		if f.IsDir() {
			return nil
		}

		cfileExt := ""
		cfileName := f.Name()
		cfilePath := migrationPath + "/" + f.Name()

		if strings.HasSuffix(cfileName, ".saz") {
			cfileExt = ".saz"
		} else if strings.HasSuffix(cfileName, ".farx") {
			cfileExt = ".farx"
		} else {
			return nil
		}

		var parseErr error
		if cfileExt == ".saz" {
			sazParser := &SazParser{SazFilePath: cfilePath, AutoResponder: refAutoResponder, OrigFileName: cfileName, UploadPath: uploadPath, SazFileName: cfileName}
			parseErr = sazParser.Handle()

			defer os.RemoveAll(cfilePath)
			defer os.RemoveAll(uploadPath + "/" + cfileName)
		} else if cfileExt == ".farx" {
			farxParser := &FarxParser{FarxFilePath: cfilePath, AutoResponder: refAutoResponder, OrigFileName: f.Name()}
			parseErr = farxParser.Handle()

			defer os.RemoveAll(cfilePath)
		}

		if parseErr != nil {
			fmt.Println(parseErr)
		}

		fmt.Printf("Visited: %s\n", path)
		return nil
	})
}
