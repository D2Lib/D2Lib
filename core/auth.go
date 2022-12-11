package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/securecookie"
	"net/http"
	"os"
	"strings"
)

var cookieHandler = securecookie.New( // generate cookie key
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
var keys = getKeys()

func getKeys() []string {
	fileByte, _ := os.ReadFile("./keypool.lock")
	return strings.Split(string(fileByte), "\n")
}

// get username for current session
func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil { // if the cookie of this session exists in local cookie pool
		cookieValue := make(map[string]string)                                             // get encoded cookie value
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil { // decode cookie
			userName = cookieValue["name"] // get username
		}
	}
	return userName // return username
}

// set login status
func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{ // define a new cookie value structure
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{ // define a new cookie structure for current username
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

// clear login status
func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{ // define an empty cookie structure
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie) // set for current session
}

// check if key exists in key pool
func contains(s []string, str string) bool {
	for _, v := range s { // walk in keypool
		if v == str { // if current key exists in key pool
			return true // stop loop and return true
		}
	}

	return false // if not exists, return false
}

// LoginHandler handle login
func LoginHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// receive login form
		name := request.FormValue("name")
		pass := request.FormValue("pass")
		redirectTarget := "/login?status=err"
		hash := sha256.Sum256([]byte(name + " " + pass))                                       // generate hash
		if contains(keys, hex.EncodeToString(hash[:])) && os.Getenv("D2LIB_elogn") == "true" { // compare with stored hash IDs
			// exist hash ID (correct psw)
			setSession(name, response) // set cookie session to remember login status
			log.Tracef("[%s] > logged in as %s", request.RemoteAddr, name)
			redirectTarget = "/" // set redirect target to root
		}
		http.Redirect(response, request, redirectTarget, 302) // do redirect
	}
}

// LoginPageHandler handle login page
func LoginPageHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if os.Getenv("D2LIB_elogn") == "true" { // check if login is enabled
			// handle login page
			log.Tracef("[%s] > request for login page", request.RemoteAddr)
			userName := getUserName(request) // get username
			if userName != "" {              // if username exists for current session (already logged in)
				http.Redirect(response, request, "/", 302) // redirect to root
			} else { // if no username (not logged in)
				loginStatus := ""
				if request.URL.Query().Get("status") == "err" {
					loginStatus = "Error username or password"
				}
				_, _ = fmt.Fprintf(response, strings.ReplaceAll(os.Getenv("D2LIB_lpage"), "{{ ERR }}", loginStatus)) // print login page
			}
		} else { // if not
			http.Redirect(response, request, "/", 302) // redirect to root
		}
	}
}

// LogoutHandler handle logout
func LogoutHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if os.Getenv("D2LIB_elogn") == "true" { // check if login is enabled
			userName := getUserName(request) // get username
			log.Tracef("[%s] > logout %s", request.RemoteAddr, userName)
			clearSession(response)                     // delete cookie for current session (clear login status)
			http.Redirect(response, request, "/", 302) // redirect to root
		} else { // if not
			http.Redirect(response, request, "/", 302) // redirect to root
		}
	}
}

func EditAccount(splitCmd []string) {
	if splitCmd[1] == "add" {
		hash := sha256.Sum256([]byte(splitCmd[2] + " " + splitCmd[3]))
		openFile, _ := os.OpenFile(os.Getenv("D2LIB_root")+"/keypool.lock", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		_, _ = openFile.Write([]byte(hex.EncodeToString(hash[:]) + "\n"))
		_ = openFile.Close()
		keys = append(keys, hex.EncodeToString(hash[:]))
		log.Warnf("Successfully added account: %s %s", splitCmd[2], splitCmd[3])
	} else if splitCmd[1] == "del" {
		hash := sha256.Sum256([]byte(splitCmd[2] + " " + splitCmd[3]))
		poolByte, _ := os.ReadFile(os.Getenv("D2LIB_root") + "/keypool.lock")
		openFile, _ := os.OpenFile(os.Getenv("D2LIB_root")+"/keypool.lock", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		_, _ = openFile.Write([]byte(strings.ReplaceAll(string(poolByte), hex.EncodeToString(hash[:])+"\n", "")))
		_ = openFile.Close()
		for i, v := range keys {
			if v == hex.EncodeToString(hash[:]) {
				keys = append(keys[:i], keys[i+1:]...)
				break
			}
		}
		log.Warnf("Successfully deleted account: %s %s", splitCmd[2], splitCmd[3])
	}
}
