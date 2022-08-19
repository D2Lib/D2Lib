package main

/*
D2Lib-Go
Version 0.2.0a
By ArthurZhou
Follows GPL-2.0 License

GitHub repo: https://github.com/D2Lib/D2Lib
*/

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const VER = "0.2.0a"
const AUTHOR = "ArthurZhou"
const ProjRepo = "https://github.com/D2Lib/D2Lib"

// global configurations variables
var addr string
var handleURL string
var storageLocation string
var homePage string
var enableLogin bool

var rootPath, _ = os.Getwd()          // get working dir path
var cookieHandler = securecookie.New( // generate cookie key
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
var router = mux.NewRouter()
var keys []string
var loginPage string
var indexPage string

func configure() {
	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		// config fie does not exist
		log.Println("\033[93m> Config file does not exists. Now creating one...\033[0m")
		newFile, err := os.Create("config.ini") // create a new one
		if err != nil {
			log.Fatalf("\033[91m> Failed to create file: %v", err)
			return
		}
		_, _ = newFile.WriteString("[Network]\naddr=\"127.0.0.1:8080\"\n\n[Handler]\nhandleURL=\"/\"\n\n" +
			"[Storage]\nstorageLocation=\"storage\"\nhomePage=\"home.md\"\n\n[Security]\nenableLogin=false\n")
		_ = newFile.Close()
	}
	cfg, err := ini.Load("config.ini") // read config file
	if err != nil {
		log.Fatalf("\033[91m> FATAL: Failed to read file: %v\n", err)
	}
	// read configurations
	addr = cfg.Section("Network").Key("addr").String()
	handleURL = cfg.Section("Handler").Key("handleURL").String()
	storageLocation = cfg.Section("Storage").Key("storageLocation").String()
	homePage = cfg.Section("Storage").Key("homePage").String()
	enableLogin = cfg.Section("Security").Key("enableLogin").MustBool()

	// load templates
	loginPath := rootPath + "/templates/login.html"
	loFileByte, _ := os.ReadFile(loginPath)
	loginPage = string(loFileByte)

	indexPath := rootPath + "/templates/index.html"
	inFileByte, _ := os.ReadFile(indexPath)
	indexPage = string(inFileByte)
}

func dirScan() {
	if _, err := os.Stat(rootPath + "/" + storageLocation); os.IsNotExist(err) {
		// storage folder does not exist
		log.Println("\033[93m> Storage folder does not exist. Now creating one...\033[0m")
		_ = os.Mkdir(rootPath+"/"+storageLocation, 0755)
	}
	if _, err := os.Stat(rootPath + "/" + storageLocation + "/" + homePage); os.IsNotExist(err) {
		// home page does not exist
		log.Println("\033[93m> Home page does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/" + storageLocation + "/" + homePage)
		_, _ = newFile.WriteString("# Home Page")
		_ = newFile.Close()
	}
	if _, err := os.Stat(rootPath + "/keypool.lock"); os.IsNotExist(err) {
		// keypool does not exist
		log.Println("\033[93m> Key pool does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/keypool.lock")
		_ = newFile.Close()
	}

	if _, err := os.Stat(rootPath + "/templates"); os.IsNotExist(err) {
		// templates folder does not exist
		log.Println("\033[93m> Templates folder does not exist. Now creating one...\033[0m")
		_ = os.Mkdir(rootPath+"/templates", 0755)
	}

	if _, err := os.Stat(rootPath + "/templates/login.html"); os.IsNotExist(err) {
		// login template does not exist
		log.Println("\033[93m> Login template does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/templates/login.html")
		_, _ = newFile.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Login</title>\n    <style>\n        body {\n            background-color: #292929;\n        }\n\n        div {\n            margin: 20px;\n            padding: 10px;\n        }\n\n        hr {\n            border-top: 5px solid #c3c3c3;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            border-radius: 3px;\n        }\n\n        h1 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 250%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h2 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 220%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h3 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 190%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h4 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 170%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h5 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 150%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h6 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 130%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        code {\n            color: #c8c8c8;\n            font-family: Courier New, serif;\n        }\n\n        a {\n            text-decoration: None;\n            color: #58748d;\n            font-family: sans-serif;\n            letter-spacing: 1px;\n        }\n\n        a:link, a:visited {\n            color: #58748d;\n        }\n\n        a:hover {\n            color: #539899;\n            text-decoration: none;\n        }\n\n        a:active {\n            color: #c3c3c3;\n            background: #101010;\n        }\n\n        p {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        ul {\n            list-style-type: square;\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        ol {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        table {\n            border: 2px solid #101010;\n            font-family: Helvetica, serif;\n        }\n\n        th {\n            border: 1px solid #101010;\n            font-family: Helvetica, serif;\n            color: #c3c0c3;\n            font-weight: bold;\n            text-align: center;\n            padding: 10px;\n        }\n\n        td {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n            text-align: center;\n            padding: 2px;\n        }\n\n        input {\n            color: #c3c3c3;\n            font-family: Courier, serif;\n            background: #101010;\n            border-top-width: 0;\n            border-bottom-width: 2px;\n            border-left-width: 0;\n            border-right-width: 0;\n            height: 30px;\n            width: 500px;\n            font-size: 15px;\n        }\n\n        ::placeholder {\n            text-align: center;\n        }\n    </style>\n</head>\n<body>\n<center>\n    <h1>Login</h1>\n    <form method=\"post\" action=\"/login\">\n        <label for=\"name\">\n        <input type=\"text\" id=\"name\" name=\"name\" placeholder=\"> Username <\"></label>\n        <br>\n        <label for=\"pass\">\n        <input type=\"password\" id=\"pass\" name=\"pass\" placeholder=\"> Password <\"></label>\n        <br>\n        <input type=\"submit\" name=\"Login\">\n    </form>\n</center>\n</body>\n</html>")
		_ = newFile.Close()
	}

	if _, err := os.Stat(rootPath + "/templates/index.html"); os.IsNotExist(err) {
		// index template does not exist
		log.Println("\033[93m> Index template does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/templates/index.html")
		_, _ = newFile.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>\n{{ TITLE }}\n    </title>\n    <style>\n        body {\n            background-color: #292929;\n        }\n\n        @keyframes fadeInAnimation {\n            0% {\n                opacity: 0;\n            }\n            100% {\n                opacity: 1;\n            }\n        }\n\n        div {\n            margin: 20px;\n            padding: 10px;\n        }\n\n        hr {\n            border-top: 5px solid #c3c3c3;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            border-radius: 3px;\n        }\n\n        h1 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 250%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h2 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 220%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h3 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 190%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h4 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 170%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h5 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 150%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h6 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 130%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        code {\n            color: #c8c8c8;\n            font-family: Courier New, serif;\n        }\n\n        a {\n            text-decoration: None;\n            color: #58748d;\n            font-family: sans-serif;\n            letter-spacing: 1px;\n        }\n\n        a:link, a:visited {\n            color: #58748d;\n        }\n\n        a:hover {\n            color: #539899;\n            text-decoration: none;\n        }\n\n        a:active {\n            color: #c3c3c3;\n            background: #101010;\n        }\n\n        p {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        p.warn {\n            color: #e33a3a;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        ul {\n            list-style-type: square;\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        ol {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        table {\n            border: 2px solid #101010;\n            font-family: Helvetica, serif;\n        }\n\n        th {\n            border: 1px solid #101010;\n            font-family: Helvetica, serif;\n            color: #c3c0c3;\n            font-weight: bold;\n            text-align: center;\n            padding: 10px;\n        }\n\n        td {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n            text-align: center;\n            padding: 2px;\n        }\n\n        input {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            background: #101010;\n            border-top-width: 0;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            height: 20px;\n            width: 200px;\n        }\n\n        div.fade {\n            animation: fadeInAnimation ease 0.3s;\n            animation-iteration-count: 1;\n            animation-fill-mode: forwards;\n        }\n\n        li.logout {\n            float: right;\n        }\n\n        li.menu a {\n            display: block;\n            color: white;\n            text-align: center;\n            padding: 14px 16px;\n            text-decoration: none;\n        }\n\n        li.menu a:hover {\n            background-color: #111;\n        }\n\n        li.logout a {\n            display: block;\n            color: #958a4b;\n            text-align: center;\n            padding: 14px 16px;\n            text-decoration: none;\n        }\n\n        li.logout a:hover {\n            background-color: #111;\n        }\n    </style>\n</head>\n<body>\n<div><p class=\"warn\">{{ ACCOUNT }}</p></div>\n<div class=\"fade\">\n{{ CONTENT }}\n</div>\n<div>\n    <br><hr><p>Powered by D2Lib</p>\n</div>\n</body>\n</html>")
		_ = newFile.Close()
	}

	log.Println("> Setup finished! Restart is required!")
	os.Exit(0)
}

func getKeys() {
	fileByte, _ := os.ReadFile(rootPath + "/keypool.lock")
	keys = strings.Split(string(fileByte), "\n")
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	// check if the user has logged in
	userName := getUserName(request)
	if userName == "" && enableLogin { // not logged in
		log.Printf("[%s] > redirect because not logged in\n", request.RemoteAddr)
		http.Redirect(response, request, "/login", 302)
	} else { // logged in
		params := mux.Vars(request)
		reqURL := "/" + params["path"] // get request url
		log.Printf("[%s] > request for url: %s\n", request.RemoteAddr, reqURL)

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
				if enableLogin {
					fileText = strings.ReplaceAll(fileText, "{{ ACCOUNT }}", "<a href=\"/logout\">Logout</a>")
				} else {
					fileText = strings.ReplaceAll(fileText, "{{ ACCOUNT }}", "")
				}
			}
			_, _ = fmt.Fprint(response, fileText) // output to http.ResponseWriter
		} else {
			// url does not exist
			log.Printf("[%s] > url does not exist: %s\n", request.RemoteAddr, reqURL)
			_, _ = fmt.Fprint(response, http.ErrMissingFile) // output to http.ResponseWriter
		}
	}
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// login handler
func loginHandler(response http.ResponseWriter, request *http.Request) {

	// handle login form
	name := request.FormValue("name")
	pass := request.FormValue("pass")
	redirectTarget := "/login"
	hash := sha256.Sum256([]byte(name + " " + pass))
	if contains(keys, hex.EncodeToString(hash[:])) && enableLogin {
		setSession(name, response)
		log.Printf("[%s] > logged in as %s\n", request.RemoteAddr, name)
		redirectTarget = "/"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func loginPageHandler(response http.ResponseWriter, request *http.Request) {
	if enableLogin {
		// handle login page
		log.Printf("[%s] > request for login page\n", request.RemoteAddr)
		userName := getUserName(request)
		if userName != "" {
			http.Redirect(response, request, "/", 302)
		} else {
			_, _ = fmt.Fprintf(response, loginPage)
		}
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// logout handler
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	if enableLogin {
		userName := getUserName(request)
		log.Printf("[%s] > logout %s\n", request.RemoteAddr, userName)
		clearSession(response)
		http.Redirect(response, request, "/", 302)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func redirectHandler(response http.ResponseWriter, request *http.Request) {
	// if request for root, redirect to home page
	http.Redirect(response, request, "/"+homePage, 302)
}

func cmd() {
	log.Println("> Command Line Tool started")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmdInput := scanner.Text()
		splitCmd := strings.Split(cmdInput, " ")
		switch {
		case splitCmd[0] == "quit":
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		case splitCmd[0] == "account" && len(splitCmd) == 4:
			if splitCmd[1] == "add" {
				hash := sha256.Sum256([]byte(splitCmd[2] + " " + splitCmd[3]))
				openFile, _ := os.OpenFile(rootPath+"/keypool.lock", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				_, _ = openFile.Write([]byte(hex.EncodeToString(hash[:]) + "\n"))
				_ = openFile.Close()
				keys = append(keys, hex.EncodeToString(hash[:]))
				log.Printf("Successfully added account: %s %s\n", splitCmd[2], splitCmd[3])
			} else if splitCmd[1] == "del" {
				// pass
			}
		default:
			log.Printf("> Unknown command: %s\n", cmdInput)
		}
	}
}

func main() {
	defer os.Exit(0)
	defer log.Fatalln("\033[91m> Process ended by deferred auto shutdown")
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			// handle SIGINT
			print("\n")
			log.Println("\033[1;30;47m> SIGINT(Interrupt Signal) received. Shutting down server...\033[0m")
			log.Println("\033[1;30;47m> Server stopped!\033[0m")
			os.Exit(0)
		case syscall.SIGTERM:
			// handle SIGTERM
			log.Fatalln("\033[91m> FATAL: Process terminated")
		}
	}()

	log.Printf("> D2Lib-Go Version %s by %s  GitHub repo %s\n", VER, AUTHOR, ProjRepo)
	log.Println("> Press Ctrl+C to stop.")
	log.Println("> Loading configurations")
	configure()
	log.Printf("\033[95m> Handle URL set to \"%s\". Working dir: %s\033[0m\n", handleURL, rootPath)

	log.Println("> Scanning working directory...")
	dirScan()
	log.Println("> Loading key pool...")
	getKeys()
	log.Println("> Done!")
	log.Printf("\033[95m> Server opened on %s\033[0m\n", addr)
	go cmd() // start cmd

	// set handlers
	if enableLogin {
		router.HandleFunc("/login", loginPageHandler).Methods("GET")
		router.HandleFunc("/login", loginHandler).Methods("POST")
		router.HandleFunc("/logout", logoutHandler).Methods("GET")
	}
	router.HandleFunc("/{path}", requestHandler).Methods("GET")
	router.HandleFunc("/", redirectHandler).Methods("GET")
	http.Handle(handleURL, router) // handle requests to requestHandler

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("\033[91m> FATAL: %v\n", err)
		return
	}
}
