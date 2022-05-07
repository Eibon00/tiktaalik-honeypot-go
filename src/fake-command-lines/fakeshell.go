package fake_command_lines

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
	"io"
	"log"
	"strings"
	"tiktaalik-honeypot-go/src/dbutil"
	"tiktaalik-honeypot-go/src/dbutil/dbstructs"
	"tiktaalik-honeypot-go/src/dbutil/dbwork"
	fake_command_lines "tiktaalik-honeypot-go/src/fake-command-lines/commands-fliter"
	write_line "tiktaalik-honeypot-go/src/write-line"
	"tiktaalik-honeypot-go/src/write-line/colors"
)

// FakeShell 假的shell，用于展示给入侵者，糊弄他们,如果入侵者被ban，则展示野兽先辈(雾
//当连接断开时返回记录的数据
func FakeShell(s ssh.Session) []string {
	//当允许登录时
	io.WriteString(s, "Linux"+fake_command_lines.RunCmdReally("uname -r")+"\n")

	//创建假的终端，设置颜色样式
	fakeTerm := term.NewTerminal(s, fmt.Sprintf("%s%s%s@%s%s%s:~#%s ",
		colors.Yellow,
		s.User(),
		colors.Green,
		colors.Yellow,
		s.LocalAddr(),
		colors.Green,
		colors.Reset,
	))

	var Record []string
	for {
		cmdline, err := fakeTerm.ReadLine()
		//记录日志,并防止日志太多给数据库挤爆了
		if len(Record) < 50 {
			Record = append(Record, cmdline)
		}

		//接收到ctrl+c或exit时退出,退出前存储日志
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		if cmdline == "exit" {
			break
		}

		//在此处应当有专门的方法先对cmdline进行预处理,暂时先这样
		var cmd dbstructs.Commands
		cmdWithArgs := strings.Split(cmdline, " ")
		command := cmdWithArgs[0]
		cmd = dbutil.SearchingInCommands(command)

		//发现未知的命令,添加到数据库
		if dbwork.IsEmpty(cmd) {
			dbutil.AddUnknownCommand(cmdline)
		}

		//命令不允许执行时,输出警告
		if cmd.IsAllowed != true {
			write_line.ColorWrite(fakeTerm, fmt.Sprintf("unknown command: %s", command), colors.Red)
			write_line.PrintEnd(fakeTerm, 1)
		}
	}

	//关闭连接，异常退出时发送警告
	err := s.Close()
	if err != nil {
		log.Fatal(err)
	}
	return Record
}

//更新对方窗口大小
//func updateSize(s ssh.Pty) {
//	go func() {
//		sigwinchCh := make(chan os.Signal, 1)
//		signal.Notify(sigwinchCh, syscall.SIG)
//	}()
//}
