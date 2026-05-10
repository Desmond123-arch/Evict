package pkg

import (
	"fmt"
	"os/exec"
	"strings"
)

type RunningProcess struct {
	COMMAND string
	PID     string
	USER    string
	FD      string
	TYPE    string
	DEVICE  int
	SIZE    string
	NODE    string
	PORT    int
	NAME    string
}

func (process RunningProcess) KillProcess() error {
	args := fmt.Sprintf("-9 %s", process.PID)

	cmd := exec.Command("kill", strings.Split(args, " ")...)

	_, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	fmt.Printf(cmd.String())
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
