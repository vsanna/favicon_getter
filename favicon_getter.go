package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	urlFilenname := os.Args[1]

	// ioutil.ReadFile: 全データを[]byteとして読み込む。string(data)で表示できる
	// file, err := ioutil.ReadFile(filename)

	// os.Open: File structを返す file_unix.go
	file, err := os.Open(urlFilenname)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		urlString := scanner.Text()
		url, err := url.Parse(urlString)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("start to fetch url: %s\n", url)
		response, err := http.Get(url.String())
		if err != nil {
			fmt.Println(err)
			continue
		}

		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		htmlTokens := html.NewTokenizer(strings.NewReader(string(b)))

	PARSE_TOKEN:
		for {
			tokenType := htmlTokens.Next()

			switch tokenType {

			case html.ErrorToken:
				break PARSE_TOKEN

			case html.TextToken:

			case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
				token := htmlTokens.Token()
				if token.Data != "link" {
					continue
				}

				rel, ok := attrByName(&token, "rel")
				if !ok {
					continue
				}

				if !hasIcon(rel) {
					continue
				}

				href, _ := attrByName(&token, "href")
				if strings.HasPrefix(href, "http") {
					fmt.Println(href)
				} else {
					fmt.Println((*url).Scheme + "://" + (*url).Host + href)
				}
			}
		}

		response.Body.Close()
	}
}

func attrByName(token *html.Token, attr string) (string, bool) {
	for _, attribute := range token.Attr {
		// fmt.Println(attribute.Key, attribute.Val)
		if attribute.Key == attr {
			return attribute.Val, true
		}
	}
	return "none", false
}

type anySlice []interface{}

func (arr anySlice) contains(item interface{}) bool {
	for _, el := range arr {
		if el == item {
			return true
		}
	}
	return false
}

func hasIcon(rel string) bool {
	relItems := strings.Split(rel, " ")
	// https://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces-in-go
	// NOTE string[]を
	var a anySlice = make([]interface{}, len(relItems))
	for i := range relItems {
		a[i] = relItems[i]
	}

	return a.contains("icon")
}
