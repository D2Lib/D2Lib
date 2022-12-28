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
	switch {                                 // execute commands
	case splitCmd[0] == "quit":
		status, reason = fQuit()
	case splitCmd[0] == "account" && len(splitCmd) == 4:
		status, reason = fAccount(splitCmd)
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

func fAccount(splitCmd []string) (bool, string) {
	return EditAccount(splitCmd)
}
