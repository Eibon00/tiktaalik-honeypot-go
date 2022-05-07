package commands_fliter

import (
	"bytes"
	"os/exec"
	"strings"
)

// RunCmdReally 允许被直接执行的命令，执行时附带被过滤器允许的参数
func RunCmdReally(cmdStr string) string {
	list := strings.Split(cmdStr, " ")
	cmd := exec.Command(list[0], list[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String()
	} else {
		return out.String()
	}
}

// RunCmdVirtual 假装执行了命令,返回事先准备好的内容
func RunCmdVirtual(cmdStr string) string {
	return cmdStr
}
