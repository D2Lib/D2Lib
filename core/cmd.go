package core

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func Cmd(rootPath string) {
	log.Println("> Command Line Tool started")
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
		log.Printf("> Unknown command: %s\n", cmdInput)
	}
}

func fQuit() {
	log.Println("\033[93m> Due to some issues on windows systems, we`ve removed this function permantely! Please use Ctrl+C instead!\033[0m")
}

func fAccount(splitCmd []string, rootPath string) {
	EditAccount(splitCmd, rootPath)
}
