package core

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
rmexec.go
Handle remote commands and send it to core.cmd.Executor
*/

func RemoteExecutor() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		log := GetLogger()                         // get logger
		authKey := request.URL.Query().Get("auth") // get auth key from url param: auth
		command := request.URL.Query().Get("cmd")  // get command from url param: cmd
		log.Infof("Incoming remote task: `%s` from client: %s", command, request.RemoteAddr)
		if authKey == os.Getenv("D2LIB_rmkey") { // authorize client by auth key
			if command == "" { // check if param `cmd` is empty or not
				log.Infof("Blank command from client: %s", request.RemoteAddr)
				_, _ = fmt.Fprint(response, "NONE_TASK==|==Cannot execute blank task!")
			} else {
				status, reason := Executor(command) // send task to core.cmd.Executor
				_, _ = fmt.Fprint(response, strings.ToUpper(strconv.FormatBool(status))+"==|=="+reason)
			}
		} else { // wrong key
			log.Warnf("Failed to authorize client: %s", request.RemoteAddr)
			_, _ = fmt.Fprint(response, "AUTH_FAIL==|==Failed to authorize!")
		}
	}
}
