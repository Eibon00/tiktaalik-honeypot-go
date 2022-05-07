package notifier

import "os/exec"

// SendNotify 通知发送工具，通过命令发送Linux桌面通知，待改进
func SendNotify(appName, title, body string) {
	var args []string
	args = append(args, "-a", appName)
	args = append(args, title)
	args = append(args, body)

	cmd := exec.Command("notify-send", args...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
