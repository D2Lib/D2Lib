package main

/*
D2Lib-Go
By ArthurZhou
Follows GPL-2.0 License

GitHub repo: https://github.com/D2Lib/D2Lib
*/

import (
	"D2Lib/core"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/ini.v1"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const VER = "0.2.2-s20221211"
const AUTHOR = "ArthurZhou"
const ProjRepo = "https://github.com/D2Lib/D2Lib"

// global configurations variables

var rootPath, _ = os.Getwd() // get working dir path
var router = mux.NewRouter()

var log = core.Log

func configure() {
	_ = os.Setenv("D2LIB_root", rootPath)

	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		// config fie does not exist
		log.Warn("Config file does not exists. Now creating one...")
		newFile, err := os.Create("config.ini") // create a new one
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
			return
		}
		_, _ = newFile.WriteString("[Network]\naddr=\"127.0.0.1:8080\"\n\n[Storage]\nstorageLocation=\"storage\"\nhomePage=\"home.md\"\nfnfPage=\"<h1>404</h1><br><center><p>Page Not Found</p></center>\"\n\n[Security]\nenableLogin=false\n")
		_ = newFile.Close()
	}
	cfg, err := ini.Load("config.ini") // read config file
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	// read configurations
	_ = os.Setenv("D2LIB_addr", cfg.Section("Network").Key("addr").String())
	_ = os.Setenv("D2LIB_sloc", cfg.Section("Storage").Key("storageLocation").String())
	_ = os.Setenv("D2LIB_hpage", cfg.Section("Storage").Key("homePage").String())
	_ = os.Setenv("D2LIB_elogn", cfg.Section("Security").Key("enableLogin").String())
	_ = os.Setenv("D2LIB_fpage", cfg.Section("Storage").Key("fnfPage").String())

	// load templates
	loginPath := rootPath + "/templates/login.html"
	loFileByte, _ := os.ReadFile(loginPath)
	_ = os.Setenv("D2LIB_lpage", string(loFileByte))

	indexStylePath := rootPath + "/templates/index.css"
	insFileByte, _ := os.ReadFile(indexStylePath)
	indexStyle := string(insFileByte)
	_ = os.Setenv("D2LIB_istyle", indexStyle)
	indexPath := rootPath + "/templates/index.html"
	inFileByte, _ := os.ReadFile(indexPath)
	_ = os.Setenv("D2LIB_ipage", strings.ReplaceAll(string(inFileByte), "{{ STYLE }}", "<style>"+indexStyle+"</style>"))
}

func dirScan() {
	fixTimes := 0

	if _, err := os.Stat(rootPath + "/" + os.Getenv("D2LIB_sloc")); os.IsNotExist(err) {
		// storage folder does not exist
		log.Warn("Storage folder does not exist. Now creating one...")
		_ = os.Mkdir(rootPath+"/"+os.Getenv("D2LIB_sloc"), 0755)
		fixTimes += 1
	}
	if _, err := os.Stat(rootPath + "/" + os.Getenv("D2LIB_sloc") + "/" + os.Getenv("D2LIB_hpage")); os.IsNotExist(err) {
		// home page does not exist
		log.Warn("Home page does not exist. Now creating one...")
		newFile, _ := os.Create(rootPath + "/" + os.Getenv("D2LIB_sloc") + "/" + os.Getenv("D2LIB_hpage"))
		_, _ = newFile.WriteString("# Home Page")
		_ = newFile.Close()
		fixTimes += 1
	}
	if _, err := os.Stat(rootPath + "/keypool.lock"); os.IsNotExist(err) {
		// keypool does not exist
		log.Warn("Key pool does not exist. Now creating one...")
		newFile, _ := os.Create(rootPath + "/keypool.lock")
		_ = newFile.Close()
		fixTimes += 1
	}

	if _, err := os.Stat(rootPath + "/templates"); os.IsNotExist(err) {
		// templates folder does not exist
		log.Warn("Templates folder does not exist. Now creating one...")
		_ = os.Mkdir(rootPath+"/templates", 0755)
		fixTimes += 1
	}

	if _, err := os.Stat(rootPath + "/templates/login.html"); os.IsNotExist(err) {
		// login template does not exist
		log.Warn("Login template does not exist. Now creating one...")
		newFile, _ := os.Create(rootPath + "/templates/login.html")
		_, _ = newFile.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Login</title>\n    <style>\n        body {\n            background-color: #292929;\n        }\n\n        div {\n            margin: 20px;\n            padding: 10px;\n        }\n\n        hr {\n            border-top: 5px solid #c3c3c3;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            border-radius: 3px;\n        }\n\n        h1 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 250%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h2 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 220%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h3 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 190%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h4 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 170%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h5 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 150%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h6 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 130%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        code {\n            color: #c8c8c8;\n            font-family: Courier New, serif;\n        }\n\n        a {\n            text-decoration: None;\n            color: #58748d;\n            font-family: sans-serif;\n            letter-spacing: 1px;\n        }\n\n        a:link, a:visited {\n            color: #58748d;\n        }\n\n        a:hover {\n            color: #539899;\n            text-decoration: none;\n        }\n\n        a:active {\n            color: #c3c3c3;\n            background: #101010;\n        }\n\n        p {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        ul {\n            list-style-type: square;\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        ol {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        table {\n            border: 2px solid #101010;\n            font-family: Helvetica, serif;\n        }\n\n        th {\n            border: 1px solid #101010;\n            font-family: Helvetica, serif;\n            color: #c3c0c3;\n            font-weight: bold;\n            text-align: center;\n            padding: 10px;\n        }\n\n        td {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n            text-align: center;\n            padding: 2px;\n        }\n\n        input {\n            color: #c3c3c3;\n            font-family: Courier, serif;\n            background: #101010;\n            border-top-width: 0;\n            border-bottom-width: 2px;\n            border-left-width: 0;\n            border-right-width: 0;\n            height: 30px;\n            width: 500px;\n            font-size: 15px;\n        }\n\n        ::placeholder {\n            text-align: center;\n        }\n    </style>\n</head>\n<body>\n<center>\n    <h1>Login</h1>\n    <form method=\"post\" action=\"/login\">\n        <label for=\"name\">\n        <input type=\"text\" id=\"name\" name=\"name\" placeholder=\"> Username <\"></label>\n        <br>\n        <label for=\"pass\">\n        <input type=\"password\" id=\"pass\" name=\"pass\" placeholder=\"> Password <\"></label>\n        <br>\n        <input type=\"submit\" name=\"Login\">\n    </form>\n</center>\n</body>\n</html>")
		_ = newFile.Close()
		fixTimes += 1
	}

	if _, err := os.Stat(rootPath + "/templates/index.html"); os.IsNotExist(err) {
		// index template does not exist
		log.Warn("Index template does not exist. Now creating one...")
		newFile, _ := os.Create(rootPath + "/templates/index.html")
		_, _ = newFile.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>\n{{ TITLE }}\n    </title>\n    <style>\n        body {\n            background-color: #292929;\n        }\n\n        @keyframes fadeInAnimation {\n            0% {\n                opacity: 0;\n            }\n            100% {\n                opacity: 1;\n            }\n        }\n\n        div {\n            margin: 20px;\n            padding: 10px;\n        }\n\n        hr {\n            border-top: 5px solid #c3c3c3;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            border-radius: 3px;\n        }\n\n        h1 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 250%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h2 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 220%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h3 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 190%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h4 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 170%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h5 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 150%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        h6 {\n            color: #c3c3c3;\n            font-family: Arial, serif;\n            font-size: 130%;\n            text-align: center;\n            letter-spacing: 3px;\n        }\n\n        code {\n            color: #c8c8c8;\n            font-family: Courier New, serif;\n        }\n\n        a {\n            text-decoration: None;\n            color: #58748d;\n            font-family: sans-serif;\n            letter-spacing: 1px;\n        }\n\n        a:link, a:visited {\n            color: #58748d;\n        }\n\n        a:hover {\n            color: #539899;\n            text-decoration: none;\n        }\n\n        a:active {\n            color: #c3c3c3;\n            background: #101010;\n        }\n\n        p {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        p.warn {\n            color: #e33a3a;\n            font-family: Helvetica, serif;\n            font-size: 100%;\n            display: inline;\n            text-indent: 100px;\n            letter-spacing: 1px;\n            line-height: 120%;\n        }\n\n        ul {\n            list-style-type: square;\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        ol {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n        }\n\n        table {\n            border: 2px solid #101010;\n            font-family: Helvetica, serif;\n        }\n\n        th {\n            border: 1px solid #101010;\n            font-family: Helvetica, serif;\n            color: #c3c0c3;\n            font-weight: bold;\n            text-align: center;\n            padding: 10px;\n        }\n\n        td {\n            font-family: Helvetica, serif;\n            color: #c3c3c3;\n            text-align: center;\n            padding: 2px;\n        }\n\n        input {\n            color: #c3c3c3;\n            font-family: Helvetica, serif;\n            background: #101010;\n            border-top-width: 0;\n            border-bottom-width: 0;\n            border-left-width: 0;\n            border-right-width: 0;\n            height: 20px;\n            width: 200px;\n        }\n\n        div.fade {\n            animation: fadeInAnimation ease 0.3s;\n            animation-iteration-count: 1;\n            animation-fill-mode: forwards;\n        }\n\n        ul.menu {\n            list-style-type: none;\n            margin: 0;\n            padding: 0;\n            overflow: hidden;\n            background-color: #333;\n        }\n\n        li.menu {\n            float: left;\n        }\n\n        li.logout {\n            float: right;\n        }\n\n        li.menu a {\n            display: block;\n            color: white;\n            text-align: center;\n            padding: 14px 16px;\n            text-decoration: none;\n        }\n\n        li.menu a:hover {\n            background-color: #111;\n        }\n\n        li.logout a {\n            display: block;\n            color: #958a4b;\n            text-align: center;\n            padding: 14px 16px;\n            text-decoration: none;\n        }\n\n        li.logout a:hover {\n            background-color: #111;\n        }\n    </style>\n</head>\n<body>\n<div>\n    <ul class=\"menu\">{{ MENU }}</ul>\n</div>\n<div class=\"fade\">\n{{ CONTENT }}\n</div>\n<div>\n    <br><hr><p>Powered by D2Lib</p>\n</div>\n</body>\n</html>")
		_ = newFile.Close()
		fixTimes += 1
	}

	if fixTimes != 0 {
		log.Info("Setup finished! Restart is required!")
		os.Exit(0)
	}
}

func main() {
	// add deferred functions to prevent uncompleted shutdowns
	defer os.Exit(0)
	defer log.Trace("Process ended by deferred auto shutdown")
	fmt.Print(`
  ____ ____  _     _ _     
 |  _ \___ \| |   (_) |__  
 | | | |__) | |   | | '_ \ 
 | |_| / __/| |___| | |_) |
 |____/_____|_____|_|_.__/ 
                           `)
	fmt.Print("\n")

	signalChannel := make(chan os.Signal, 2) // bind for signals
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() { // handle Ctrl+C signal and force kill signal
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			// handle SIGINT
			print("\n")
			log.Trace("SIGINT(Interrupt Signal) received. Shutting down server...")
			log.Info("Server stopped!")
			os.Exit(0)
		case syscall.SIGTERM:
			// handle SIGTERM
			log.Fatal("Process terminated")
		}
	}()

	log.Infof("D2Lib-Go Version %s by %s  GitHub repo %s", VER, AUTHOR, ProjRepo)
	log.Info("Press Ctrl+C to stop.")
	log.Debug("Loading configurations")
	configure() // load config
	log.Debugf("Working dir: %s", rootPath)
	log.Debug("Scanning working directory...")
	dirScan() // check dir
	log.Debug("Rendering menu bar...")
	menuRender := "<div><ul class=\"menu\">"
	if os.Getenv("D2LIB_elogn") == "true" { // add "logout" button to menubar
		menuRender += "<li class=\"logout\"><a class=\"logout\" href=\"/logout\">Log out</a></li>"
	}
	menuRender += "<li class=\"menu\"><a class=\"menu\" href=\"/\">Home</a></li>" // add "home" button to menubar
	files, _ := os.ReadDir(rootPath + "/" + os.Getenv("D2LIB_sloc"))              // search for folders in current dir
	for _, f := range files {                                                     // render menubar
		if f.IsDir() {
			menuRender += "<li class=\"menu\"><a class=\"menu\" href=\"/docs?path=" + f.Name() + "/" + os.Getenv("D2LIB_hpage") + "\">" + f.Name() + "</a></li>"
		}
	}
	menuRender += "</ul></div>"
	_ = os.Setenv("D2LIB_menu", menuRender)
	log.Info("Done!")
	go core.Cmd() // start cmd

	// set handlers
	if os.Getenv("D2LIB_elogn") == "true" { // set auth functions
		router.HandleFunc("/login", core.LoginPageHandler()).Methods("GET")
		router.HandleFunc("/login", core.LoginHandler()).Methods("POST")
		router.HandleFunc("/logout", core.LogoutHandler()).Methods("GET")
	}
	router.HandleFunc("/favicon.ico", core.FaviconHandler()).Methods("GET")
	router.HandleFunc("/docs", core.RequestHandler()).Methods("GET")
	router.HandleFunc("/", core.RedirectHandler()).Methods("GET")
	log.Infof("Server opened on %s", os.Getenv("D2LIB_addr"))
	http.Handle("/", router) // handle requests to mux router

	err := http.ListenAndServe(os.Getenv("D2LIB_addr"), nil) // start http server
	if err != nil {
		log.Panicf("%v", err)
		return
	}
}
