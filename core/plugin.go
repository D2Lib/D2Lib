package core

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Manifest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	D2LibVersion []int  `json:"d2lib_version"`
	Url          string `json:"url"`
	Command      string `json:"command"`
	Port         int    `json:"port"`
	Method       string `json:"method"`
}

func ScanPlugin(router *mux.Router) {
	log := GetLogger()
	VER_NUM, _ := strconv.Atoi(os.Getenv("D2LIB_vernum"))
	var manifest Manifest
	if _, err := os.Stat(os.Getenv("D2LIB_root") + "/plugins"); os.IsNotExist(err) {
		log.Warn("Plugin folder does not exist. Now creating one...")
		_ = os.Mkdir(os.Getenv("D2LIB_root")+"/plugins", 0755)
	}
	files, _ := ioutil.ReadDir(os.Getenv("D2LIB_root") + "/plugins")
	for _, f := range files {
		if f.IsDir() {
			manifestStr, _ := os.ReadFile(os.Getenv("D2LIB_root") + "/plugins/" + f.Name() + "/manifest.json")
			_ = json.Unmarshal(manifestStr, &manifest)
			log.Debugf("Loading plugin: %s  version: %s", manifest.Name, manifest.Version)
			if manifest.D2LibVersion[0] < VER_NUM && manifest.D2LibVersion[1] > VER_NUM {
				cmd := strings.Split(manifest.Command, " ")
				cmdExec := exec.Command(cmd[0], strings.Join(cmd[1:len(cmd)], " "))
				out, err := cmdExec.Output()
				if err != nil {
					log.Errorf("Failed to load plugin: %s  reason: %v", manifest.Name, err)
				}
				log.Debugf("Plugin `%s` return `%s` while executing system command `%s`", manifest.Name, out, manifest.Command)

				proxy, err := NewProxy("http://localhost:" + strconv.Itoa(manifest.Port))
				if err != nil {
					panic(err)
				}
				if manifest.Method != "" {
					router.HandleFunc(manifest.Url, ProxyRequestHandler(proxy)).Methods(manifest.Method)
				} else {
					router.HandleFunc(manifest.Url, ProxyRequestHandler(proxy))
				}
			} else {
				log.Warnf("Plugin `%s` require D2Lib version `%d`-`%d`  current is `%d`", manifest.Name, manifest.D2LibVersion[0], manifest.D2LibVersion[1], VER_NUM)
			}
		}
	}
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
