package pkg

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
)

// GET the running proceses
func GetRunningProcesses() ([]RunningProcess, error) {
	res := make([]RunningProcess, 0)
	args := "-i -P -n"
	cmd := exec.Command("lsof", strings.Split(args, " ")...)

	stdoutput, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdoutput)
	scanner.Split(bufio.ScanLines)
	//skip headers
	scanner.Scan()
	for scanner.Scan() {
		m := strings.Fields(scanner.Text())
		device, _ := strconv.Atoi(m[5])
		port, _ := strconv.Atoi((strings.Split(strings.Split(m[8], " ")[0], ":"))[1])
		process := RunningProcess{
			COMMAND: m[0],
			PID:     m[1],
			USER:    m[2],
			FD:      m[3],
			TYPE:    m[4],
			DEVICE:  device,
			SIZE:    m[6],
			NODE:    m[7],
			PORT:    port,
			NAME:    m[8],
		}
		res = append(res, process)
	}
	cmd.Wait()
	return res, nil
}
