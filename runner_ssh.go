package hotomata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	ECHO          = 53
	TTY_OP_ISPEED = 128
	TTY_OP_OSPEED = 129
)

type SSHRunner struct {
}

func (r *SSHRunner) Run(machine Machine, command string) *TaskResponse {
	var response = &TaskResponse{
		Log:    &bytes.Buffer{},
		Action: TaskActionContinue,
		Status: TaskStatusSuccess,
	}
	var cmdErr error

	if machine.Hostname == "127.0.0.1" && machine.Port == 0 {
		// Local execution
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdout = response.Log
		cmd.Stderr = response.Log
		cmdErr = cmd.Run()
	} else {
		// Remote execution
		hostname := machine.Hostname + ":" + strconv.Itoa(machine.Port)
		client, err := ssh.Dial("tcp", hostname, machine.SSHConfig)
		if err != nil {
			fmt.Printf("Failed to dial: %s: %s\n", hostname, err.Error())
			os.Exit(1)
		}
		session, err := client.NewSession()
		if err != nil {
			fmt.Printf("Unable to connect: %s: %s\n", hostname, err.Error())
			os.Exit(1)
		}
		defer session.Close()
		defer client.Close()

		modes := ssh.TerminalModes{
			ECHO:          0,
			TTY_OP_ISPEED: 14400,
			TTY_OP_OSPEED: 14400,
		}
		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			fmt.Printf("Request for terminal failed: %s: %s\n", hostname, err.Error())
			os.Exit(1)
		}
		session.Stdout = response.Log
		session.Stderr = response.Log
		cmdErr = session.Run(command)
	}

	if cmdErr != nil {
		response.Log.WriteString(cmdErr.Error() + "\n")
		response.Action = TaskActionAbort
		response.Status = TaskStatusError
		return response
	}

	// if the command outputed json that unmarchals to a TaskResponse
	// treat that as the answer, this is what power our assertions and
	// plugin system
	splittedLog := strings.Split(response.Log.String(), "\n")
	lastSignificantLine := splittedLog[len(splittedLog)-2]
	var taskResponseJson TaskResponse
	if err := json.Unmarshal([]byte(lastSignificantLine), &taskResponseJson); err == nil && taskResponseJson.Action != "" {
		taskResponseJson.Log = response.Log
		return &taskResponseJson
	}

	return response
}
