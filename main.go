package main

import (
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"tiktaalik-honeypot-go/src/configurator"
	"tiktaalik-honeypot-go/src/dbutil"
	fakeshell "tiktaalik-honeypot-go/src/fake-command-lines"
	"tiktaalik-honeypot-go/src/logger"
	private_host_key "tiktaalik-honeypot-go/src/private-host-key"
	"time"
)

func init() {
	flag.BoolVar(&enableWebUI, "webui", false, "Enable Web UI")
}

var enableWebUI bool
var attempts = 0
var config configurator.Config

func sessionHandler(s ssh.Session) {
	Result := fakeshell.FakeShell(s)
	logger.LogClient(s.RemoteAddr().String(), strconv.FormatInt(time.Now().UnixNano(), 10), Result)
}

func authHandler(ctx ssh.Context, passwd string) bool {
	attempts++
	var clientIP string
	if strings.ContainsRune(ctx.RemoteAddr().String(), ':') {
		clientIP, _, _ = net.SplitHostPort(ctx.RemoteAddr().String())
	} else {
		clientIP = ctx.RemoteAddr().String()
	}
	body := fmt.Sprintf("User: %s,Password: %s, Address: %s, Status: ", ctx.User(), passwd, clientIP)

	if ctx.User() != config.Auth.User || passwd != config.Auth.Password {
		status := "failed"
		log.Printf(fmt.Sprintf("[%d]%s%s", attempts, body, status))
		return false
	}
	status := "connected"
	logger.LogBlacklist(clientIP, strconv.FormatInt(time.Now().UnixNano(), 10))
	log.Printf(fmt.Sprintf("[%d]%s%s", attempts, body, status))
	return true
}

func main() {
	//fmt.Println("[+] wdnmd,is running!")
	config = configurator.ParseConfigFile()
	dbutil.WriteConfigFile(config)
	s := &ssh.Server{
		Addr:            fmt.Sprintf("0.0.0.0:%s", string(rune(config.Auth.Port))),
		Handler:         sessionHandler,
		PasswordHandler: authHandler,
		IdleTimeout:     30 * time.Second,
	}

	var keyFilePath string
	if os.Getenv("HONEYPOT_HOSTKEYFILE") == "" {
		keyFilePath = fmt.Sprintf("%s/config/hostkey_rsa", configurator.GetRootPath())
	} else {
		keyFilePath = os.Getenv("HONEYPOT_HOSTKEYFILE")
	}
	keyFileExists := configurator.FileExists(keyFilePath)
	if keyFileExists {
		key, err := private_host_key.ReadHostKeyFile(keyFilePath)
		if err != nil {
			log.Fatal(err)
		}
		s.AddHostKey(key)
	}
	log.Printf("[+]Starting Honeypot Server on Address: %v\n", s.Addr)
	if keyFileExists {
		log.Printf("[+]Honeypot HostKey Mode: user-input-file")
	} else {
		log.Print("[+]Honeypot HostKey Mode: auto-generated")
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
