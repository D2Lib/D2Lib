package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
)

var cookieHandler = securecookie.New( // generate cookie key
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

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
func LoginHandler(keys []string, enableLogin bool) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// receive login form
		name := request.FormValue("name")
		pass := request.FormValue("pass")
		redirectTarget := "/login"
		hash := sha256.Sum256([]byte(name + " " + pass))                // generate hash
		if contains(keys, hex.EncodeToString(hash[:])) && enableLogin { // compare with stored hash IDs
			// exist hash ID (correct psw)
			setSession(name, response) // set cookie session to remember login status
			log.Printf("[%s] > logged in as %s\n", request.RemoteAddr, name)
			redirectTarget = "/" // set redirect target to root
		}
		http.Redirect(response, request, redirectTarget, 302) // do redirect
	}
}

// LoginPageHandler handle login page
func LoginPageHandler(enableLogin bool, loginPage string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if enableLogin { // check if login is enabled
			// handle login page
			log.Printf("[%s] > request for login page\n", request.RemoteAddr)
			userName := getUserName(request) // get username
			if userName != "" {              // if username exists for current session (already logged in)
				http.Redirect(response, request, "/", 302) // redirect to root
			} else { // if no username (not logged in)
				_, _ = fmt.Fprintf(response, loginPage) // print login page
			}
		} else { // if not
			http.Redirect(response, request, "/", 302) // redirect to root
		}
	}
}

// LogoutHandler handle logout
func LogoutHandler(enableLogin bool) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if enableLogin { // check if login is enabled
			userName := getUserName(request) // get username
			log.Printf("[%s] > logout %s\n", request.RemoteAddr, userName)
			clearSession(response)                     // delete cookie for current session (clear login status)
			http.Redirect(response, request, "/", 302) // redirect to root
		} else { // if not
			http.Redirect(response, request, "/", 302) // redirect to root
		}
	}
}
