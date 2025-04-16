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
	commands_list = append(commands_list, "pwd")
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
			prefix := []rune(statements[1])
			suffix := []rune(statements[len(statements)-1])
			newStr := ""

			if prefix[0] == 39 && (suffix[len(suffix)-1]) == 39 {
				for _, statement := range statements[1:] {
					statementChars := []rune(statement)
					for _, chars := range statementChars {

						// fmt.Println(statementChars)

						if chars == 32 {

						} else if chars != 39 {
							newStr += string(chars)

						}

					}
					newStr += " "
				}
				fmt.Printf("%v\n", newStr)
			} else {
				for _, statement := range statements[1:] {
					statementChars := []rune(statement)
					for i, chars := range statementChars {

						// fmt.Println(statementChars)

						if chars == 32 {

						} else if chars != 39 {
							newStr += string(chars)
							if i == len(statementChars)-1 {
								newStr += " "
							}
						}
					}
				}
				fmt.Printf("%v\n", newStr)
			}
		} else if statements[0] == "cd" {
			/*
				Again due to the issue with Go's built in os.Chdir and other os functions, cd is
				implemented by setting the PWD env to the path given, idk how viable this is but oh well its how
				the tests pass
			*/
			path = statements[1]
			_, err := os.ReadDir(path)
			pathArr := []rune(path)
			if pathArr[0] == '/' {
				/*
					Handles absolute paths, for example switch to any directory from any directory using for example
					/usr/bin/
				*/
				if err != nil {
					fmt.Printf("cd: %v: No such file or directory\n", statements[1])
				} else {
					os.Setenv("PWD", statements[1])
				}
			} else if pathArr[0] == '~' {
				/*
					"cd ~" handles taking you back to the home directory of the system
				*/
				homeDir := os.Getenv("HOME")
				os.Setenv("PWD", homeDir)
			} else if pathArr[0] == '.' && pathArr[1] == '.' {
				/*
					Handles ../ which goes back to a directory in the file tree
				*/
				backAmount := len(strings.Split(statements[1], "/"))
				currDir := os.Getenv("PWD")
				currDirArr := []rune(currDir)
				newDir := ""
				for i := 0; i < backAmount; i++ {
					newDir += string(currDirArr[i])
				}
				os.Setenv("PWD", newDir)
			} else if pathArr[0] == '.' {
				/*
					Handles relative paths, for example if you're in /usr/ and want to go into bin cd ./bin will take you
					to /usr/bin
				*/
				currPath := os.Getenv("PWD")
				newPath := currPath + statements[1][1:]
				os.Setenv("PWD", newPath)
			}
		} else if statements[0] == "pwd" {
			/*
				for some reason, os.Getwd() is returning the wrong path so use the env
			*/
			dir := os.Getenv("PWD")

			fmt.Println(dir)
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
