package parser

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

	matchFiles, err := filepath.Glob(folderPath + "/*_s.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range matchFiles {
		is := strings.ReplaceAll(filepath.Base(f), "_s.txt", "")

		fileC, err := ioutil.ReadFile(folderPath + "/" + is + "_c.txt")
		if err != nil {
			fmt.Println(err)
			break
		}
		fileS, _ := ioutil.ReadFile(folderPath + "/" + is + "_s.txt")

		response := ParseStringContent(string(fileS))
		response.Label = orgFileName + "_" + is

		parser.AutoResponder.AddOrUpdateResponse(response)

		rule := parser.parseRule(string(fileC))
		rule.ResponseID = response.ID

		parser.AutoResponder.AddOrUpdateRule(rule)
	}
}

var hostRegex = regexp.MustCompile(`(?i)host: `)

func (parser *SazParser) parseRule(content string) *responder.Rule {
	linez := strings.Split(content, "\n")
	uriLine := strings.Split(linez[0], " ")
	hostLine := strings.TrimSpace(hostRegex.ReplaceAllString(linez[1], ""))

	//hostline is a backup, first line might have full uri
	uriPattern := uriLine[1]
	if !strings.Contains(uriLine[1], hostLine) {
		uriPattern = "http://" + hostLine + uriPattern
	}

	return &responder.Rule{IsActive: true, URLPattern: uriPattern, Method: uriLine[0], MatchType: "EXACT"}
}

//ParseStringContent parses given http content
func ParseStringContent(content string) *responder.Response {
	return responder.NewResponseFromString(content)
}
