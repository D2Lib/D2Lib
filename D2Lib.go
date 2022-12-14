package main

/*
D2Lib-Go
By ArthurZhou
Follows GPL-2.0 License

GitHub repo: https://github.com/D2Lib/D2Lib
*/

import (
	"D2Lib/core"
	"bytes"
	"github.com/gorilla/mux"
	"gopkg.in/ini.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const VER = "0.2.2-s20221112-2-hotfix"
const AUTHOR = "ArthurZhou"
const ProjRepo = "https://github.com/D2Lib/D2Lib"

// global configurations variables
var addr string
var storageLocation string
var homePage string
var enableLogin bool
var fnfPage string

var rootPath, _ = os.Getwd() // get working dir path
var router = mux.NewRouter()
var loginPage string
var indexPage string
var menuRender string

func configure() {
	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		// config fie does not exist
		log.Println("\033[93m> Config file does not exists. Now creating one...\033[0m")
		newFile, err := os.Create("config.ini") // create a new one
		if err != nil {
			log.Fatalf("\033[91m> Failed to create file: %v", err)
			return
		}
		_, _ = newFile.WriteString("[Network]\naddr=\"127.0.0.1:8080\"\n\n[Storage]\nstorageLocation=\"storage\"\nhomePage=\"home.md\"\nfnfPage=\"<h1>404</h1><br><center><p>Page Not Found</p></center>\"\n\n[Security]\nenableLogin=false\n")
		_ = newFile.Close()
	}
	cfg, err := ini.Load("config.ini") // read config file
	if err != nil {
		log.Fatalf("\033[91m> FATAL: Failed to read file: %v\n", err)
	}
	// read configurations
	addr = cfg.Section("Network").Key("addr").String()
	storageLocation = cfg.Section("Storage").Key("storageLocation").String()
	homePage = cfg.Section("Storage").Key("homePage").String()
	enableLogin = cfg.Section("Security").Key("enableLogin").MustBool()
	fnfPage = cfg.Section("Storage").Key("fnfPage").String()

	// load templates
	loginPath := rootPath + "/templates/login.html"
	loFileByte, _ := os.ReadFile(loginPath)
	loginPage = string(loFileByte)

	indexPath := rootPath + "/templates/index.html"
	inFileByte, _ := os.ReadFile(indexPath)
	indexPage = string(inFileByte)
}

func dirScan() {
	fixTimes := 0

	if _, err := os.Stat(rootPath + "/" + storageLocation); os.IsNotExist(err) {
		// storage folder does not exist
		log.Println("\033[93m> Storage folder does not exist. Now creating one...\033[0m")
		_ = os.Mkdir(rootPath+"/"+storageLocation, 0755)
		fixTimes += 1
	}
	if _, err := os.Stat(rootPath + "/" + storageLocation + "/" + homePage); os.IsNotExist(err) {
		// home page does not exist
		log.Println("\033[93m> Home page does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/" + storageLocation + "/" + homePage)
		_, _ = newFile.WriteString("# Home Page")
		_ = newFile.Close()
		fixTimes += 1
	}
	if _, err := os.Stat(rootPath + "/keypool.lock"); os.IsNotExist(err) {
		// keypool does not exist
		log.Println("\033[93m> Key pool does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/keypool.lock")
		_ = newFile.Close()
		fixTimes += 1
	}

	if _, err := os.Stat(rootPath + "/templates"); os.IsNotExist(err) {
		// templates folder does not exist
		log.Println("\033[93m> Templates folder does not exist. Now creating one...\033[0m")
		_ = os.Mkdir(rootPath+"/templates", 0755)
		fixTimes += 1
	}

	if _, err := os.Stat(rootPath + "/templates/login.html"); os.IsNotExist(err) {
		// login template does not exist
		log.Println("\033[93m> Login template does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/templates/login.html")
		_, _ = newFile.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Login</title>\n    <style>\n        body {\n            background-color: #292929;\n        }\n\n        div {\n            margin: 20px;\n            padding: 10px;\n        }\n\n        hr {\n            border-top: 5px solid #c3c3c3;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            border-radius: 3px;\n        }\n\n        h1 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 250%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h2 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 220%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h3 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 190%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h4 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 170%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h5 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 150%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h6 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 130%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        code {\n            color: #c8c8c8;\n            font-family: Courier New, serif;\n        }\n\n        a {\n            text-decoration: None;\n            color: #58748d;\n            font-family: sans-serif;\n            letter-spacing: 1px;\n        }\n\n        a:link, a:visited {\n            color: #58748d;\n        }\n\n        a:hover {\n            color: #539899;\n            text-decoration: none;\n        }\n\n        a:active {\n            color: #c3c3c3;\n            background: #101010;\n        }\n\n        p {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        ul {\n            list-style-type: square;\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        ol {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        table {\n            border: 2px solid #101010;\n            font-family: Helvetica, serif;\n        }\n\n        th {\n            border: 1px solid #101010;\n            font-family: Helvetica, serif;\n            color: #c3c0c3;\n            font-weight: bold;\n            text-align: center;\n            padding: 10px;\n        }\n\n        td {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n            text-align: center;\n            padding: 2px;\n        }\n\n        input {\n            color: #c3c3c3;\n            font-family: Courier, serif;\n            background: #101010;\n            border-top-width: 0;\n            border-bottom-width: 2px;\n            border-left-width: 0;\n            border-right-width: 0;\n            height: 30px;\n            width: 500px;\n            font-size: 15px;\n        }\n\n        ::placeholder {\n            text-align: center;\n        }\n    </style>\n</head>\n<body>\n<center>\n    <h1>Login</h1>\n    <form method=\"post\" action=\"/login\">\n        <label for=\"name\">\n        <input type=\"text\" id=\"name\" name=\"name\" placeholder=\"> Username <\"></label>\n        <br>\n        <label for=\"pass\">\n        <input type=\"password\" id=\"pass\" name=\"pass\" placeholder=\"> Password <\"></label>\n        <br>\n        <input type=\"submit\" name=\"Login\">\n    </form>\n</center>\n</body>\n</html>")
		_ = newFile.Close()
		fixTimes += 1
	}

	if _, err := os.Stat(rootPath + "/templates/index.html"); os.IsNotExist(err) {
		// index template does not exist
		log.Println("\033[93m> Index template does not exist. Now creating one...\033[0m")
		newFile, _ := os.Create(rootPath + "/templates/index.html")
		_, _ = newFile.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>\n{{ TITLE }}\n    </title>\n    <style>\n        body {\n            background-color: #292929;\n        }\n\n        @keyframes fadeInAnimation {\n            0% {\n                opacity: 0;\n            }\n            100% {\n                opacity: 1;\n            }\n        }\n\n        div {\n            margin: 20px;\n            padding: 10px;\n        }\n\n        hr {\n            border-top: 5px solid #c3c3c3;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            border-radius: 3px;\n        }\n\n        h1 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 250%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h2 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 220%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h3 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 190%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h4 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 170%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h5 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 150%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h6 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 130%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        code {\n            color: #c8c8c8;\n            font-family: Courier New, serif;\n        }\n\n        a {\n            text-decoration: None;\n            color: #58748d;\n            font-family: sans-serif;\n            letter-spacing: 1px;\n        }\n\n        a:link, a:visited {\n            color: #58748d;\n        }\n\n        a:hover {\n            color: #539899;\n            text-decoration: none;\n        }\n\n        a:active {\n            color: #c3c3c3;\n            background: #101010;\n        }\n\n        p {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        p.warn {\n            color: #e33a3a;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        ul {\n            list-style-type: square;\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        ol {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        table {\n            border: 2px solid #101010;\n            font-family: Helvetica, serif;\n        }\n\n        th {\n            border: 1px solid #101010;\n            font-family: Helvetica, serif;\n            color: #c3c0c3;\n            font-weight: bold;\n            text-align: center;\n            padding: 10px;\n        }\n\n        td {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n            text-align: center;\n            padding: 2px;\n        }\n\n        input {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            background: #101010;\n            border-top-width: 0;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            height: 20px;\n            width: 200px;\n        }\n\n        div.fade {\n            animation: fadeInAnimation ease 0.3s;\n            animation-iteration-count: 1;\n            animation-fill-mode: forwards;\n        }\n\n        ul.menu {\n            list-style-type: none;\n            margin: 0;\n            padding: 0;\n            overflow: hidden;\n            background-color: #333;\n        }\n\n        li.menu {\n            float: left;\n        }\n\n        li.logout {\n            float: right;\n        }\n\n        li.menu a {\n            display: block;\n            color: white;\n            text-align: center;\n            padding: 14px 16px;\n            text-decoration: none;\n        }\n\n        li.menu a:hover {\n            background-color: #111;\n        }\n\n        li.logout a {\n            display: block;\n            color: #958a4b;\n            text-align: center;\n            padding: 14px 16px;\n            text-decoration: none;\n        }\n\n        li.logout a:hover {\n            background-color: #111;\n        }\n    </style>\n</head>\n<body>\n<div>\n    <ul class=\"menu\">{{ MENU }}</ul>\n</div>\n<div class=\"fade\">\n{{ CONTENT }}\n</div>\n<div>\n    <br><hr><p>Powered by D2Lib</p>\n</div>\n</body>\n</html>")
		_ = newFile.Close()
		fixTimes += 1
	}

	if fixTimes != 0 {
		log.Println("> Setup finished! Restart is required!")
		os.Exit(0)
	}
}

func main() {
	// add deferred functions to prevent uncompleted shutdowns
	defer os.Exit(0)
	defer log.Fatalln("\033[91m> Process ended by deferred auto shutdown")

	buf := bytes.Buffer{}                          // set a new buffer to store logs
	log.SetOutput(io.MultiWriter(os.Stdout, &buf)) // set logger output

	log.SetPrefix("STARTUP > ")
	signalChannel := make(chan os.Signal, 2) // bind for signals
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() { // handle Ctrl+C signal and force kill signal
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			// handle SIGINT
			print("\n")
			log.Println("\033[1;30m> SIGINT(Interrupt Signal) received. Shutting down server...\033[0m")
			log.Println("\033[1;30m> Server stopped!\033[0m")
			os.Exit(0)
		case syscall.SIGTERM:
			// handle SIGTERM
			log.Fatalln("\033[91m> FATAL: Process terminated")
		}
	}()

	log.Printf("> D2Lib-Go Version %s by %s  GitHub repo %s\n", VER, AUTHOR, ProjRepo)
	log.Println("> Press Ctrl+C to stop.")
	log.Println("> Loading configurations")
	configure() // load config
	log.Printf("\033[95m> Working dir: %s\033[0m\n", rootPath)

	log.Println("> Scanning working directory...")
	dirScan() // check dir
	log.Println("> Rendering menu bar...")
	if enableLogin { // add "logout" button to menubar
		menuRender += "<li class=\"logout\"><a class=\"logout\" href=\"/logout\">Log out</a></li>"
	}
	menuRender += "<li class=\"menu\"><a class=\"menu\" href=\"/\">Home</a></li>" // add "home" button to menubar
	files, _ := ioutil.ReadDir(rootPath + "/" + storageLocation)                  // search for folders in current dir
	for _, f := range files {                                                     // render menubar
		if f.IsDir() {
			menuRender += "<li class=\"menu\"><a class=\"menu\" href=\"/docs?path=" + f.Name() + "/" + homePage + "\">" + f.Name() + "</a></li>"
		}
	}
	log.Println("> Done!")
	go core.Cmd(rootPath) // start cmd

	log.SetPrefix("MAIN > ")                                                            // set a prefix
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC | log.Lmicroseconds) // ste logger flags

	// set handlers
	if enableLogin { // set auth functions
		router.HandleFunc("/login", core.LoginPageHandler(enableLogin, loginPage)).Methods("GET")
		router.HandleFunc("/login", core.LoginHandler(enableLogin)).Methods("POST")
		router.HandleFunc("/logout", core.LogoutHandler(enableLogin)).Methods("GET")
	}
	router.HandleFunc("/favicon.ico", core.FaviconHandler(rootPath)).Methods("GET")
	router.HandleFunc("/docs", core.RequestHandler(enableLogin, rootPath, storageLocation, indexPage, menuRender, fnfPage)).Methods("GET")
	router.HandleFunc("/", core.RedirectHandler(homePage)).Methods("GET")
	log.Printf("\033[95m> Server opened on %s\033[0m\n", addr)
	http.Handle("/", router) // handle requests to mux router

	err := http.ListenAndServe(addr, nil) // start http server
	if err != nil {
		log.Fatalf("\033[91m> FATAL: %v\n", err)
		return
	}
}
