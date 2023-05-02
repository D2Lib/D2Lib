package core

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Manifest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Version      string   `json:"version"`
	D2LibVersion []int    `json:"d2lib_version"`
	Urls         []string `json:"urls"`
	Command      string   `json:"command"`
	Port         int      `json:"port"`
	Method       string   `json:"method"`
}

type logrusErrorWriter struct{}

func (w logrusErrorWriter) Write(p []byte) (n int, err error) {
	GetLogger().Errorf("%s", strings.ReplaceAll(string(p), "\n", ""))
	return len(p), nil
}

func ScanPlugin(router *mux.Router) {
	logger := GetLogger()
	VerNum, _ := strconv.Atoi(os.Getenv("D2LIB_vernum"))
	var manifest Manifest
	if _, err := os.Stat(os.Getenv("D2LIB_root") + "/plugins"); os.IsNotExist(err) {
		logger.Warn("Plugin folder does not exist. Now creating one...")
		_ = os.Mkdir(os.Getenv("D2LIB_root")+"/plugins", 0755)
	}
	files, _ := ioutil.ReadDir(os.Getenv("D2LIB_root") + "/plugins")
	for _, f := range files {
		if f.IsDir() {
			manifestStr, _ := os.ReadFile(os.Getenv("D2LIB_root") + "/plugins/" + f.Name() + "/manifest.json")
			_ = json.Unmarshal(manifestStr, &manifest)
			logger.Debugf("Loading plugin: %s  version: %s", manifest.Name, manifest.Version)
			if manifest.D2LibVersion[0] <= VerNum && manifest.D2LibVersion[1] >= VerNum {
				cmd := strings.Split(manifest.Command, " ")
				cmdExec := exec.Command(cmd[0], strings.Join(cmd[1:len(cmd)], " "))
				out, err := cmdExec.Output()
				if err != nil {
					logger.Errorf("Failed to load plugin: %s  reason: %v", manifest.Name, err)
				}
				logger.Debugf("Plugin `%s` return `%s` while executing system command `%s`", manifest.Name, out, manifest.Command)

				proxy, err := NewProxy("http://localhost:"+strconv.Itoa(manifest.Port), manifest.Urls[0])
				errorLogger := log.New(logrusErrorWriter{}, "", 0)
				proxy.ErrorLog = errorLogger
				if err != nil {
					panic(err)
				}
				router.HandleFunc("/test/{path}", ProxyRequestHandler(proxy, "/test"))
				if manifest.Method != "" {
					for _, url := range manifest.Urls {
						router.HandleFunc(url, ProxyRequestHandler(proxy, url)).Methods(manifest.Method)
					}
				} else {
					for _, url := range manifest.Urls {
						router.HandleFunc(url, ProxyRequestHandler(proxy, url))
					}
				}
			} else {
				logger.Warnf("Plugin `%s` require D2Lib version `%d`-`%d`  current is `%d`", manifest.Name, manifest.D2LibVersion[0], manifest.D2LibVersion[1], VerNum)
			}
		}
	}
}

func NewProxy(targetHost string, rootUrl string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return NewSingleHostReverseProxy(url, rootUrl), nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, url string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func NewSingleHostReverseProxy(target *url.URL, rootUrl string) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		originUrl := strings.Replace(req.URL.Path, rootUrl, "", 1) + "?" + req.URL.RawQuery
		if originUrl == "" {
			originUrl = "/"
		}
		fmt.Println(originUrl)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Header.Set("Origin-Url", originUrl)
	}
	return &httputil.ReverseProxy{Director: director}
}
