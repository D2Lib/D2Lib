package core

import (
	"bufio"
	"os"
	"strings"
)

func Cmd(rootPath string) {
	log.Debugln("Command Line Tool started")
	scanner := bufio.NewScanner(os.Stdin) // check input
	for scanner.Scan() {
		cmdInput := scanner.Text() // scan input
		Executor(cmdInput, rootPath)
	}
}

func Executor(cmdInput string, rootPath string) {
	splitCmd := strings.Split(cmdInput, " ") // split args
	switch {                                 // execute commands
	case splitCmd[0] == "quit":
		fQuit()
	case splitCmd[0] == "account" && len(splitCmd) == 4:
		fAccount(splitCmd, rootPath)
	default:
		log.Errorf("Unknown command: %s", cmdInput)
	}
}

func fQuit() {
	log.Warnln("Due to some issues on windows systems, we`ve removed this function permantely! Please use Ctrl+C instead!")
}

func fAccount(splitCmd []string, rootPath string) {
	EditAccount(splitCmd, rootPath)
}
