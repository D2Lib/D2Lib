package core

import (
	"bufio"
	"os"
	"strings"
)

/*
cmd.go
Bind for inputs, extract inputs and execute functions
*/

func Cmd() {
	log := GetLogger()
	log.Debug("Command Line Tool started")
	scanner := bufio.NewScanner(os.Stdin) // check input
	for scanner.Scan() {
		cmdInput := scanner.Text() // scan input
		Executor(cmdInput)
	}
}

func Executor(cmdInput string) (bool, string) {
	log := GetLogger()
	var status bool
	var reason string
	splitCmd := strings.Split(cmdInput, " ") // split args
	switch splitCmd[0] {                     // execute commands
	case "quit":
		status, reason = fQuit()
	case "version":
		status, reason = fVersion()
	case "account":
		if len(splitCmd) == 4 {
			status, reason = fAccount(splitCmd)
		} else {
			log.Error("No enough params: ", splitCmd)
			status = false
			reason = "No enough params: " + strings.Join(splitCmd, " ")
		}
	case "reload":
		if len(splitCmd) == 2 {
			status, reason = fReload(splitCmd)
			if status {
				log.Warn("Success!")
			}
		} else {
			log.Error("No enough params: ", splitCmd)
			status = false
			reason = "No enough params: " + strings.Join(splitCmd, " ")
		}

	default:
		log.Errorf("Unknown command: %s", cmdInput)
		status = false
		reason = "Unknown command: " + cmdInput
	}
	return status, reason
}

func fQuit() (bool, string) {
	log := GetLogger()
	print("\n")
	log.Warn("Server stopped!")
	os.Exit(0)
	return true, "ok"
}

func fVersion() (bool, string) {
	log := GetLogger()
	log.Warn(os.Getenv("D2LIB_ver"))
	return true, os.Getenv("D2LIB_ver")
}

func fAccount(splitCmd []string) (bool, string) {
	return EditAccount(splitCmd)
}

func fReload(splitCmd []string) (bool, string) {
	log := GetLogger()
	switch splitCmd[1] {
	case "config":
		log.Warn("Reload configurations may cause many issues. For a complete reload, please restart your server.")
		log.Warn("Reloading configurations")
		return LoadConfig()
	case "template":
		log.Warn("Reloading templates")
		return LoadTemplate()
	default:
		log.Error("Unknown param: " + splitCmd[1])
		return false, "Unknown param: " + splitCmd[1]
	}
}
