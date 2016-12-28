package replayer

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Arguments struct {
	File      FileArgument
	Http      HttpArgument
	Parse     ParseArgument
	Regexp    RegexpArgument
	ThreadMax int
}

type FileArgument struct {
	FilePath string
	FileType string
}

type HttpArgument struct {
	BaseUri string
	Headers string
	Cookies string
}

type ParseArgument struct {
	Delimiter       string
	UriStemColumn   int
	UriQueryColumn  int
	BodyColumn      int
	BeginLine       int
	UserAgentColumn int
	VerbColumn      int
	BodyTypeColumn  int
}

type RegexpArgument struct {
	Regexp  string
	Pattern string
	Replace string
}

// Replay Test
func Replay(args Arguments) {
	file, err := os.Open(args.File.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	delimiter := args.Parse.Delimiter
	lineCounter := 0
	beginLine := 1
	if args.Parse.BeginLine != 0 {
		beginLine = args.Parse.BeginLine
	}

	for scanner.Scan() {
		lineCounter = lineCounter + 1
		if lineCounter >= beginLine {

			//TODO : CHECK FOR REGEXP

			lineArray := delimit(scanner.Text(), delimiter)
			//todo call to httpCall
			fmt.Println(lineArray)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func delimit(line string, delimiter string) []string {

	//Prepare string to split.
	line = replaceCharBetweenQuotes(line, delimiter, true)
	result := strings.Split(line, delimiter)

	for i, value := range result {
		result[i] = replaceCharBetweenQuotes(value, delimiter, false)
		result[i] = strings.Replace(result[i], "\"", "", -1)
	}

	return result
}

func replaceCharBetweenQuotes(line string, delimiter string, replaceDelimiter bool) string {

	betweenQuotesRegexp := regexp.MustCompile("\".*?\"")
	replaceBy := "§"

	matches := betweenQuotesRegexp.FindAllStringSubmatch(line, -1)
	if matches != nil && len(matches) > 0 {
		for _, match := range matches {
			sentence := match[0]

			if replaceDelimiter == true {
				line = strings.Replace(line, sentence, strings.Replace(sentence, delimiter, replaceBy, -1), -1)
			} else {
				line = strings.Replace(line, sentence, strings.Replace(sentence, replaceBy, delimiter, -1), -1)
			}
		}
	}

	return line
}

func httpCall(info []string, args Arguments) {
	var client http.Client
	var req *http.Request
	var url string

	if info[args.Parse.VerbColumn-1] == "GET" {
		url = args.Http.BaseUri + info[args.Parse.UriStemColumn-1] + info[args.Parse.UriQueryColumn]
		req, _ = http.NewRequest("GET", url, nil)

	} else if info[args.Parse.VerbColumn-1] == "POST" && args.Parse.BodyTypeColumn > 0 && args.Parse.BodyColumn > 0 {
		url = args.Http.BaseUri + info[args.Parse.UriStemColumn-1]
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer([]byte(info[args.Parse.BodyColumn])))
	}

	if args.Http.Headers != "" {
		//TODO
	}

	if args.Http.Cookies != "" {
		//TODO
	}

	client.Do(req)

}
