package web

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/asalih/http-auto-responder/responder"
)

const maxUploadSize = 25 * 1024 * 1024 // 25 MB
const uploadPath = "./tmp"

func init() {
	HTTPHandlerFuncs["/http-auto-responder/import"] = func(w http.ResponseWriter, r *http.Request, ar *responder.AutoResponder) {
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

		if fileHeader.Size > maxUploadSize || !strings.HasSuffix(fileHeader.Filename, ".saz") {
			http.Redirect(w, r, redirectURI+"?err=file_size", 302)
			return
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Redirect(w, r, redirectURI+"?err=invalid_file", 302)
			return
		}

		fileType := http.DetectContentType(fileBytes)
		if fileType != "application/zip" {
			http.Redirect(w, r, redirectURI+"?err=invalid_file_type", 302)
			return
		}

		fileName := randToken(12)
		newPath := filepath.Join(uploadPath, fileName+".saz")

		_, ferr := os.Stat(uploadPath)

		if os.IsNotExist(ferr) {
			errDir := os.MkdirAll(uploadPath, 0755)
			if errDir != nil {
				log.Fatal(err)
			}
		}

		newFile, err := os.Create(newPath)

		if err != nil {
			http.Redirect(w, r, redirectURI+"?err=cant_save", 302)
			return
		}

		if _, err := newFile.Write(fileBytes); err != nil {
			http.Redirect(w, r, redirectURI+"?err=cant_write", 302)
			return
		}

		newFile.Close()

		zf, err := zip.OpenReader(newPath)

		if err != nil {
			http.Redirect(w, r, redirectURI+"?err=saved_zip_read", 302)
			return
		}

		sazFolderName := uploadPath + "/" + fileName
		for _, file := range zf.File {
			readZipFile(file, sazFolderName)
		}

		zf.Close()

		parseAllFilesAndSave(sazFolderName, fileHeader.Filename, ar)

		defer os.RemoveAll(sazFolderName + ".saz")
		defer os.RemoveAll(sazFolderName)

		http.Redirect(w, r, redirectURI+"?success=1", 302)
	}
}

func readZipFile(file *zip.File, folderPath string) {
	fc, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer fc.Close()

	_, ferr := os.Stat(folderPath)

	if os.IsNotExist(ferr) {
		errDir := os.MkdirAll(folderPath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	newFile, ferr := os.Create(folderPath + "/" + strings.ReplaceAll(file.Name, "raw/", ""))
	defer newFile.Close()

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		panic(err)
	}

	newFile.Write(content)
}

func parseAllFilesAndSave(folderPath string, orgFileName string, ar *responder.AutoResponder) {
	i := 1
	for {
		fileC, err := ioutil.ReadFile(folderPath + "/" + strconv.Itoa(i) + "_c.txt")
		if err != nil {
			break
		}
		fileS, _ := ioutil.ReadFile(folderPath + "/" + strconv.Itoa(i) + "_s.txt")

		response := parseResponse(string(fileS))
		response.Label = orgFileName + "_" + strconv.Itoa(i)

		ar.AddOrUpdateResponse(response)

		rule := parseRule(string(fileC))
		rule.ResponseID = response.ID

		ar.AddOrUpdateRule(rule)

		i++
	}
}

func parseRule(content string) *responder.Rule {
	firstLine := strings.Split(strings.Split(content, "\n")[0], " ")

	return &responder.Rule{IsActive: true, URLPattern: firstLine[1], Method: firstLine[0]}
}

func parseResponse(content string) *responder.Response {
	lines := strings.Split(content, "\n")

	status, _ := strconv.Atoi(strings.Split(lines[0], " ")[1])
	var headers []*responder.Headers
	body := ""

	bodyStart := false
	for i := 1; i < len(lines); i++ {
		l := lines[i]
		if !bodyStart {
			if len(l) == 1 {
				bodyStart = true
				continue
			}
			idx := strings.Index(l, ":")
			rs := []rune(l)

			headers = append(headers, &responder.Headers{Key: string(rs[0:idx]), Value: strings.Trim(string(rs[idx+1:]), " ")})
		} else {
			body += l
		}
	}

	return &responder.Response{Headers: headers, Body: body, StatusCode: status}
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
