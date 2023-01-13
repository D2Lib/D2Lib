package core

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

func LoadConfig() (bool, string) {
	if _, errRead := os.Stat("config.ini"); os.IsNotExist(errRead) {
		// config fie does not exist
		fmt.Println("Config file does not exists. Now creating one...")
		newFile, errCreate := os.Create("config.ini") // create a new one
		if errCreate != nil {
			fmt.Printf("Failed to create config.ini: %v \n", errCreate)
			return false, "Failed to create config.ini"
		}
		_, _ = newFile.WriteString("[Network]\naddr=\"127.0.0.1:8080\"\n\n[Storage]\nstorageLocation=\"storage\"\nhomePage=\"home.md\"\nfnfPage=\"<h1>404</h1><br><center><p>Page Not Found</p></center>\"\n\n[Security]\nenableLogin=false\nremoteExecute=false\nremoteKey=auth\n\n[Logger]\nlogLevel=debug\nlogColor=true\nsocketLogger=false\n")
		_ = newFile.Close()
	}
	cfg, errLoad := ini.Load("config.ini") // read config file
	if errLoad != nil {
		fmt.Printf("Failed to read config.ini: %v \n", errLoad)
		return false, "Failed to read config.ini"
	}
	// read configurations
	_ = os.Setenv("D2LIB_addr", cfg.Section("Network").Key("addr").String())
	_ = os.Setenv("D2LIB_sloc", cfg.Section("Storage").Key("storageLocation").String())
	_ = os.Setenv("D2LIB_hpage", cfg.Section("Storage").Key("homePage").String())
	_ = os.Setenv("D2LIB_fpage", cfg.Section("Storage").Key("fnfPage").String())
	_ = os.Setenv("D2LIB_elogn", cfg.Section("Security").Key("enableLogin").String())
	_ = os.Setenv("D2LIB_rmexe", cfg.Section("Security").Key("remoteExecute").String())
	_ = os.Setenv("D2LIB_rmkey", cfg.Section("Security").Key("remoteKey").String())
	_ = os.Setenv("D2LIB_loglv", cfg.Section("Logger").Key("logLevel").String())
	_ = os.Setenv("D2LIB_logcl", cfg.Section("Logger").Key("logColor").String())
	_ = os.Setenv("D2LIB_sockl", cfg.Section("Logger").Key("socketLogger").String())
	_ = os.Setenv("D2LIB_saddr", cfg.Section("Logger").Key("socketAddress").String())
	_ = os.Setenv("D2LIB_sprot", cfg.Section("Logger").Key("socketProto").String())
	_ = os.Setenv("D2LIB_sapp", cfg.Section("Logger").Key("socketApp").String())
	return true, "Success"
}

func LoadTemplate() (bool, string) {
	// load templates
	loginPath := os.Getenv("D2LIB_root") + "/templates/login.html"
	loFileByte, errLo := os.ReadFile(loginPath)
	if errLo != nil {
		fmt.Printf("Failed to load login.html: %v \n", errLo)
		return false, "Failed to load login.html"
	}
	_ = os.Setenv("D2LIB_lpage", string(loFileByte))

	indexStylePath := os.Getenv("D2LIB_root") + "/templates/index.css"
	insFileByte, errIns := os.ReadFile(indexStylePath)
	if errIns != nil {
		fmt.Printf("Failed to load index.css: %v \n", errIns)
		return false, "Failed to load index.css"
	}
	indexStyle := string(insFileByte)
	_ = os.Setenv("D2LIB_istyle", indexStyle)
	indexPath := os.Getenv("D2LIB_root") + "/templates/index.html"
	inFileByte, errIn := os.ReadFile(indexPath)
	if errIn != nil {
		fmt.Printf("Failed to load index.html: %v \n", errIn)
		return false, "Failed to load index.html"
	}
	_ = os.Setenv("D2LIB_ipage", strings.ReplaceAll(string(inFileByte), "{{ STYLE }}", "<style>"+indexStyle+"</style>"))

	// render menubar
	menuRender := "<div><ul class=\"menu\">"
	if os.Getenv("D2LIB_elogn") == "true" { // add "logout" button to menubar
		menuRender += "<li class=\"logout\"><a class=\"logout\" href=\"/logout\">Log out</a></li>"
	}
	menuRender += "<li class=\"menu\"><a class=\"menu\" href=\"/\">Home</a></li>"   // add "home" button to menubar
	files, _ := os.ReadDir(os.Getenv("D2LIB_root") + "/" + os.Getenv("D2LIB_sloc")) // search for folders in storage
	for _, f := range files {                                                       // render menubar
		if f.IsDir() {
			menuRender += "<li class=\"menu\"><a class=\"menu\" href=\"/docs/" + f.Name() + "/" + os.Getenv("D2LIB_hpage") + "\">" + f.Name() + "</a></li>"
		}
	}
	menuRender += "</ul></div>"
	_ = os.Setenv("D2LIB_menu", menuRender)

	// render 404 page
	fnfText := strings.ReplaceAll(os.Getenv("D2LIB_ipage"), "{{ TITLE }}", "404 Page Not Found")
	fnfText = strings.ReplaceAll(fnfText, "{{ CONTENT }}", os.Getenv("D2LIB_fpage"))
	fnfText = strings.ReplaceAll(fnfText, "{{ MENU }}", os.Getenv("D2LIB_menu"))
	_ = os.Setenv("D2LIB_fnf", fnfText)

	return true, "Success"
}
