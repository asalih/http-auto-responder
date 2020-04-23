package parser

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/asalih/http-auto-responder/responder"
)

//SazParser Implements a parsing struct for saz parsing
type SazParser struct {
	SazFileName   string
	SazFilePath   string
	UploadPath    string
	OrigFileName  string
	AutoResponder responder.AutoResponder
}

//Handle Starts file import process
func (parser *SazParser) Handle() error {
	zf, err := zip.OpenReader(parser.SazFilePath)

	if err != nil {
		return err
	}

	sazFolderName := parser.UploadPath + "/" + parser.SazFileName
	for _, file := range zf.File {
		parser.readZipFile(file, sazFolderName)
	}

	zf.Close()

	parser.parseAllFilesAndSave(sazFolderName, parser.OrigFileName)

	return nil
}

func (parser *SazParser) readZipFile(file *zip.File, folderPath string) {
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

func (parser *SazParser) parseAllFilesAndSave(folderPath string, orgFileName string) {
	i := 1
	for {
		fileC, err := ioutil.ReadFile(folderPath + "/" + strconv.Itoa(i) + "_c.txt")
		if err != nil {
			break
		}
		fileS, _ := ioutil.ReadFile(folderPath + "/" + strconv.Itoa(i) + "_s.txt")

		response := parser.parseResponse(string(fileS))
		response.Label = orgFileName + "_" + strconv.Itoa(i)

		parser.AutoResponder.AddOrUpdateResponse(response)

		rule := parser.parseRule(string(fileC))
		rule.ResponseID = response.ID

		parser.AutoResponder.AddOrUpdateRule(rule)

		i++
	}
}

func (parser *SazParser) parseRule(content string) *responder.Rule {
	firstLine := strings.Split(strings.Split(content, "\n")[0], " ")

	return &responder.Rule{IsActive: true, URLPattern: firstLine[1], Method: firstLine[0], MatchType: "CONTAINS"}
}

func (parser *SazParser) parseResponse(content string) *responder.Response {
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

			if idx > -1 {
				rs := []rune(l)

				headers = append(headers, &responder.Headers{Key: string(rs[0:idx]), Value: strings.Trim(string(rs[idx+1:]), " ")})
			}
		} else {
			body += l
		}
	}

	return &responder.Response{Headers: headers, Body: body, StatusCode: status}
}
