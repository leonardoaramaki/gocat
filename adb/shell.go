package adb

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

// Run executes in a device shell all the commands passed to it, and returns the result in provided callback
func Run(device string, callback func(string), cmds ...string) {
	commands := make([]string, 0)
	if device != "" {
		commands = append(commands, device)
	}
	commands = append(commands, "shell")
	commands = append(commands, cmds...)
	cmd := exec.Command("adb", commands...)
	output, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal("Something went wrong...")
	}

	reader := bufio.NewReader(output)
	// Start executing the external command
	err = cmd.Start()

	if err != nil {
		log.Fatal("Something gone wrong...")
	}

	for {
		buffer, err := reader.ReadBytes('\n')
		if err == io.EOF {
			if len(buffer) == 0 {
				break
			}
		} else {
			if err != nil {
				log.Fatal(err)
			}
		}
		callback(string(buffer))
	}
	cmd.Wait()
}
