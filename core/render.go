package core

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"log"
	"net/http"
	"os"
	"strings"
)

func RequestHandler(enableLogin bool, rootPath string, storageLocation string, indexPage string, menuRender string, fnfPage string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// check if the user has logged in
		userName := getUserName(request)
		if userName == "" && enableLogin { // not logged in
			log.Printf("[%s] > redirect because not logged in\n", request.RemoteAddr)
			http.Redirect(response, request, "/login", 302)
		} else { // logged in
			reqURL := "/" + request.URL.Query().Get("path")
			log.Printf("[%s] > request for doc: %s\n", request.RemoteAddr, reqURL)
			if _, err := os.Stat(rootPath + "/" + storageLocation + reqURL); !os.IsNotExist(err) {
				// url exists
				filePath := rootPath + "/" + storageLocation + reqURL
				fileByte, _ := os.ReadFile(filePath) // read file as byte array
				fileText := string(fileByte)         // convert byte array to string
				strLength := len(reqURL) - 3
				if reqURL[strLength:] == ".md" { // is this a markdown file?
					splPath := strings.Split(reqURL, "/")
					fileName := strings.Join(splPath[len(splPath)-1:], "")
					// render markdown as html
					fileText = strings.ReplaceAll(indexPage, "{{ TITLE }}", fileName)
					fileText = strings.ReplaceAll(fileText, "{{ CONTENT }}", string(markdown.ToHTML(fileByte, nil, nil)))
					fileText = strings.ReplaceAll(fileText, "{{ MENU }}", menuRender)
				}
				_, _ = fmt.Fprint(response, fileText) // output to http.ResponseWriter
			} else {
				// url does not exist
				log.Printf("[%s] > url does not exist: %s\n", request.RemoteAddr, reqURL)
				fileText := strings.ReplaceAll(indexPage, "{{ TITLE }}", "404 Page Not Found")
				fileText = strings.ReplaceAll(fileText, "{{ CONTENT }}", fnfPage)
				fileText = strings.ReplaceAll(fileText, "{{ MENU }}", menuRender)
				_, _ = fmt.Fprint(response, fileText) // output to http.ResponseWriter
			}
		}
	}
}

func RedirectHandler(homePage string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// if request for root, redirect to home page
		http.Redirect(response, request, "/docs?path="+homePage, 302)
	}
}

func FaviconHandler(rootPath string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		log.Println("> request for favicon")
		http.ServeFile(response, request, rootPath+"/templates/favicon.ico")
	}
}
