package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	prevDir      = getExpandedCurrentDir()
	prevExitCode = 0
	PID          = os.Getppid()
)

func initShell() {
	os.Setenv("OLDPWD", getExpandedCurrentDir())
	os.Setenv("PWD", getExpandedCurrentDir())
	os.Setenv("SHELL", "/usr/bin/ohmygosh")
}

func precmd() {
	title := getCurrentUsername() + "@" + getCurrentHostname() + ":" + getCurrentDir()
	command := []string{"echo", "-ne", "\033]0;" + title + "\a"}
	RunCommand(command)
}

func getprompt() {
	// TODO: use PS1 to customize prompt
	_, ok := os.LookupEnv("PS1")

	if !ok {
		fmt.Printf("%s@%s:%s$ ", getCurrentUsername(), getCurrentHostname(), getCurrentDir())
	} else {
		fmt.Printf("> ")
	}
}

func RunShell(command string) {
	r := regexp.MustCompile(`[^\s"]+|"'([^"]*)"`)
	args := r.FindAllString(command, -1)

	switch args[0] {
	case "exit":
		os.Exit(0)
	case "cd":
		if len(args[1:]) > 1 {
			panic("wrong number of arguments to cd")
		} else if args[1] == "-" {
			os.Chdir(prevDir)
		} else {
			prevDir = getExpandedCurrentDir()
			os.Setenv("OLDPWD", prevDir)
			os.Chdir(args[1])
			os.Setenv("PWD", getExpandedCurrentDir())
		}
	case "echo":
		if len(args[1:]) >= 1 {
			for i := 1; i <= len(args[1:]); i++ {
				if strings.HasPrefix(args[i], "$") {
					fmt.Printf(os.Getenv(strings.Replace(args[i], "$", "", 1)) + " ")
				}
			}
			fmt.Println()
		}
	case "$?":
		// print exitcode of previous command
		fmt.Printf("%d\n", prevExitCode)
	default:
		RunCommand(args[0:])
	}
}

func RunCommand(command []string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			prevExitCode = exitError.ExitCode()
		} else {
			prevExitCode = 127
			fmt.Println(err)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	initShell()
	for {
		getprompt()
		precmd()
		input, _ := reader.ReadString('\n')

		// handle ctrl + d
		rune := []rune(input)
		if len(rune) == 0 {
			os.Exit(0)
		}
		commandString := strings.TrimSuffix(input, "\n")
		commandString = strings.Trim(commandString, " ")
		if commandString == "" {
			continue
		}
		RunShell(commandString)
	}
}
