package core

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"net/http"
	"os"
	"strings"
)

/*
render.go
Handle and render normal requests
*/

func RequestHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		log := GetLogger()
		// check if the user has logged in
		userName := getUserName(request)
		if userName == "" && os.Getenv("D2LIB_elogn") == "true" { // not logged in
			log.Tracef("[%s] > redirect because not logged in", request.RemoteAddr)
			http.Redirect(response, request, "/login", 302)
		} else { // logged in
			reqURL := "/" + request.URL.Query().Get("path")
			if len(reqURL) > 1 {
				log.Tracef("[%s] > request for doc: %s", request.RemoteAddr, reqURL)
				if _, err := os.Stat(os.Getenv("D2LIB_root") + "/" + os.Getenv("D2LIB_sloc") + reqURL); !os.IsNotExist(err) {
					// url exists
					filePath := os.Getenv("D2LIB_root") + "/" + os.Getenv("D2LIB_sloc") + reqURL
					fileByte, _ := os.ReadFile(filePath) // read file as byte array
					fileText := string(fileByte)         // convert byte array to string
					if reqURL[len(reqURL)-3:] == ".md" { // is this a markdown file?
						splPath := strings.Split(reqURL, "/")
						fileName := strings.Join(splPath[len(splPath)-1:], "")
						// render markdown to html and put it into the template
						fileText = strings.ReplaceAll(os.Getenv("D2LIB_ipage"), "{{ TITLE }}", fileName)
						fileText = strings.ReplaceAll(fileText, "{{ CONTENT }}", string(markdown.ToHTML(fileByte, nil, nil)))
						fileText = strings.ReplaceAll(fileText, "{{ MENU }}", os.Getenv("D2LIB_menu"))
					} else if reqURL[len(reqURL)-5:] == ".html" { // is this a markdown file?
						splPath := strings.Split(reqURL, "/")
						fileName := strings.Join(splPath[len(splPath)-1:], "")
						// replace hooks
						fileText = strings.ReplaceAll(string(fileByte), "{{ TITLE }}", fileName)
						fileText = strings.ReplaceAll(fileText, "{{ MENU }}", os.Getenv("D2LIB_menu"))
						fileText = strings.ReplaceAll(fileText, "{{ STYLE }}", "<style>"+os.Getenv("D2LIB_")+os.Getenv("D2LIB_istyle")+"</style>")
					}
					_, _ = fmt.Fprint(response, fileText) // output to http.ResponseWriter
				} else {
					// url does not exist
					fnfHandler(request, response, reqURL)
				}
			} else {
				log.Tracef("[%s] > blank url", request.RemoteAddr)
				fnfHandler(request, response, reqURL)
			}
		}
	}
}

func RedirectHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// if request for root, redirect to home page
		http.Redirect(response, request, "/docs?path="+os.Getenv("D2LIB_hpage"), 302)
	}
}

func FaviconHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) { // handle favicon
		log := GetLogger()
		log.Tracef("[%s] > request for favicon", request.RemoteAddr)
		http.ServeFile(response, request, os.Getenv("D2LIB_root")+"/templates/favicon.ico")
	}
}

func fnfHandler(request *http.Request, response http.ResponseWriter, reqURL string) { // handle 404 page
	log := GetLogger()
	log.Tracef("[%s] > url does not exist: %s", request.RemoteAddr, reqURL)
	fileText := strings.ReplaceAll(os.Getenv("D2LIB_ipage"), "{{ TITLE }}", "404 Page Not Found")
	fileText = strings.ReplaceAll(fileText, "{{ CONTENT }}", os.Getenv("D2LIB_fpage"))
	fileText = strings.ReplaceAll(fileText, "{{ MENU }}", os.Getenv("D2LIB_menu"))
	_, _ = fmt.Fprint(response, fileText) // output to http.ResponseWriter
}
