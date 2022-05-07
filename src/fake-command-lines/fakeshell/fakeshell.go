package fakeshell

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
	"io"
	"log"
	"tiktaalik-honeypot-go/src/write-line/colors"
)

// FakeShell 假的shell，用于展示给入侵者，糊弄他们,如果入侵者被ban，则展示野兽先辈字符画
func FakeShell(s ssh.Session) {

	//创建假的终端，设置颜色样式
	fakeTerm := term.NewTerminal(s, fmt.Sprintf("%s%s%s@%s%s%s:~#%s ",
		colors.Yellow,
		s.User(),
		colors.Green,
		colors.Blue,
		s.LocalAddr(),
		colors.Green,
		colors.Reset,
	))

	//接收到ctrl+c或exit时退出
	for true {
		cmdline, err := fakeTerm.ReadLine()

		//在此处应当先对cmdline进行预处理
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
	}

	//关闭连接，异常退出时发送警告
	err := s.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}

//func updateSize(s ssh.Pty) {
//	go func() {
//		sigwinchCh := make(chan os.Signal, 1)
//		signal.Notify(sigwinchCh, syscall.SIG)
//	}()
//}
