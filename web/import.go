package web

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/asalih/http-auto-responder/parser"
	"github.com/asalih/http-auto-responder/responder"
)

const maxUploadSize = 25 * 1024 * 1024 // 25 MB
const uploadPath = "./tmp"

func init() {
	HTTPHandlerFuncs["/http-auto-responder/import"] = func(w http.ResponseWriter, r *http.Request, ar responder.AutoResponder) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		redirectURI := "/http-auto-responder"

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			http.Redirect(w, r, redirectURI+"?err=err", 302)
		}

		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			http.Redirect(w, r, redirectURI+"?err=err", 302)
			fmt.Println(err)
			return
		}

		defer file.Close()

		if fileHeader.Size > maxUploadSize {
			http.Redirect(w, r, redirectURI+"?err=file_size", 302)
			return
		}

		fileExt := ""

		if strings.HasSuffix(fileHeader.Filename, ".saz") {
			fileExt = ".saz"
		} else if strings.HasSuffix(fileHeader.Filename, ".farx") {
			fileExt = ".farx"
		} else {
			http.Redirect(w, r, redirectURI+"?err=file_type", 302)
			return
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Redirect(w, r, redirectURI+"?err=invalid_file", 302)
			return
		}

		fileType := http.DetectContentType(fileBytes)
		if fileType != "application/zip" && fileType != "text/plain; charset=utf-8" {
			http.Redirect(w, r, redirectURI+"?err=invalid_file_type", 302)
			return
		}

		fileName := randToken(12)
		savedFilePath := filepath.Join(uploadPath, fileName+fileExt)

		_, ferr := os.Stat(uploadPath)

		if os.IsNotExist(ferr) {
			errDir := os.MkdirAll(uploadPath, 0755)
			if errDir != nil {
				log.Fatal(err)
			}
		}

		sazFile, err := os.Create(savedFilePath)

		if err != nil {
			http.Redirect(w, r, redirectURI+"?err=cant_save", 302)
			return
		}

		if _, err := sazFile.Write(fileBytes); err != nil {
			http.Redirect(w, r, redirectURI+"?err=cant_write", 302)
			return
		}

		sazFile.Close()

		var parseErr error
		if fileExt == ".saz" {
			sazParser := &parser.SazParser{SazFilePath: savedFilePath, AutoResponder: ar, OrigFileName: fileHeader.Filename, UploadPath: uploadPath, SazFileName: fileName}
			parseErr = sazParser.Handle()
		} else if fileExt == ".farx" {
			farxParser := &parser.FarxParser{FarxFilePath: savedFilePath}
			parseErr = farxParser.Handle()
		}

		if parseErr != nil {
			http.Redirect(w, r, redirectURI+"?err=cant_parse_file", 302)
			return
		}

		http.Redirect(w, r, redirectURI+"?success=1", 302)
	}
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
