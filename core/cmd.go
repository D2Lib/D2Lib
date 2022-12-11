package core

import (
	"bufio"
	"os"
	"strings"
)

func Cmd() {
	log.Debug("Command Line Tool started")
	scanner := bufio.NewScanner(os.Stdin) // check input
	for scanner.Scan() {
		cmdInput := scanner.Text() // scan input
		Executor(cmdInput)
	}
}

func Executor(cmdInput string) {
	splitCmd := strings.Split(cmdInput, " ") // split args
	switch {                                 // execute commands
	case splitCmd[0] == "quit":
		fQuit()
	case splitCmd[0] == "account" && len(splitCmd) == 4:
		fAccount(splitCmd)
	default:
		log.Errorf("Unknown command: %s", cmdInput)
	}
}

func fQuit() {
	log.Warn("Due to some issues on windows systems, we`ve removed this function permanently! Please use Ctrl+C instead!")
}

func fAccount(splitCmd []string) {
	EditAccount(splitCmd)
}
