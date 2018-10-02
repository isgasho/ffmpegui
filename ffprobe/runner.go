package ffprobe

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCmd(args []string, ignExitCode bool) (string, error) {
	// args = append([]string{"-c"}, args...)
	// cmd := exec.Command("/bin/sh", args...)
	var out bytes.Buffer
	log.Info(args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil && !ignExitCode {
		fmt.Println(out.String())
		panic(err)
	}
	return out.String(), nil
}
func runCmdStr(cmd string, ignExitCode bool) string {
	if str, err := runCmd(strings.Split(cmd, " "), ignExitCode); err == nil {
		return str
	}
	return ""
}

func (mp *proberCommon) runCmdPipe(cmdstr []string) (*bufio.Scanner, error) { //TODO ret error
	log.Info(cmdstr)
	cmd := exec.Command(cmdstr[0], cmdstr[1:]...)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	outscanner := bufio.NewScanner(stdout)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	go func() {
		<-mp.done
		log.Info("Received stop signal. Sending quit")
		defer stdin.Close()
		cmd.Process.Signal(os.Interrupt)
		// io.WriteString(stdin, "exit()\r\n")
		// io.WriteString(stdin, "+") TODO send keys for more/less status
		cmd.Wait()
		mp.done <- true //signal process done
	}()

	return outscanner, nil
}