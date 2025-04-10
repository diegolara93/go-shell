package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
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
		if statements[0] == "exit" {
			exitCode, _ := strconv.Atoi(statements[1])
			os.Exit(exitCode)
		} else if statements[0] == "echo" {
			for _, statement := range statements[1:] {
				fmt.Printf("%v ", statement)
			}
			fmt.Println()
		} else if statements[0] == "type" {
			directories := strings.Split(path, ":")
			found := false
			dir := ""
			for _, dirs := range directories {
				// fmt.Printf("%v\n", dirs)
				os.Chdir(dirs)
				if _, err := os.ReadFile(statements[1]); err != nil {
					// fmt.Printf("%v: not found\n", statements[1])
					found = false
				} else {
					// fmt.Printf("%v is %v/%v\n", statements[1], dirs, statements[1])
					found = true
					dir = dirs
					break
				}
			}
			if found {
				if slices.Contains(commands_list, statements[1]) {
					fmt.Printf("%v is a shell builtin\n", statements[1])
				} else {
					fmt.Printf("%v is %v/%v\n", statements[1], dir, statements[1])
				}
			} else if slices.Contains(commands_list, statements[1]) {
				fmt.Printf("%v is a shell builtin\n", statements[1])
			} else {
				fmt.Printf("%v: not found\n", statements[1])
			}
			// if slices.Contains(commands_list, statements[1]) {
			// 	fmt.Printf("%v is a shell builtin\n", statements[1])
			// } else {
			// 	fmt.Printf("%v: not found\n", statements[1])
			// }
		} else {
			directories := strings.Split(path, ":")
			found := false
			// dir := ""
			for _, dirs := range directories {
				if !isFileInDirectory(dirs, statements[0]) {
					// fmt.Printf("%v: not found\n", statements[0])
					found = false
				} else {
					// fmt.Printf("%v is %v/%v\n", statements[1], dirs, statements[1])
					found = true
					// dir = dirs
					break
				}
			}
			if found {
				exe := fmt.Sprintf("%v", statements[0])
				cmd := exec.Command(exe, statements[1])
				out, err := cmd.Output()
				if err != nil {
					// fmt.Printf("error running command: %v\n", cmd.String())
				}
				fmt.Printf(string(out))
			} else {
				fmt.Printf("%v: not found\n", statements[0])
			}
		}
	}
}
func isFileInDirectory(dirPath, fileName string) bool {
	filePath := dirPath + "/" + fileName
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
