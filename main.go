package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	commands_list := make([]string, 10)
	commands_list = append(commands_list, "echo")
	commands_list = append(commands_list, "exit")
	commands_list = append(commands_list, "type")
	path := os.Getenv("PATH")
	if path == "" {
		fmt.Println("PATH variable is empty")
	}
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("error reading command")
		}
		command = strings.TrimSpace(command)
		statements := strings.Split(command, " ")
		if command == "exit 0" {
			os.Exit(0)
		} else if statements[0] == "echo" {
			for _, statement := range statements[1:] {
				fmt.Printf("%v ", statement)
			}
			fmt.Println()
		} else if statements[0] == "type" {
			directories := strings.Split(path, ":")
			for _, dirs := range directories {
				fmt.Printf("%v\n", dirs)
				os.Chdir(dirs)
			}
			if slices.Contains(commands_list, statements[1]) {
				fmt.Printf("%v is a shell builtin\n", statements[1])
			} else {
				fmt.Printf("%v: not found\n", statements[1])
			}
		} else {
			fmt.Println(command[:len(command)] + ": command not found")
		}
	}
}
