/*****************************************************************************
 * server-go.go
 * Name:
 * NetId:
 *****************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const RECV_BUFFER_SIZE = 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 */
func server(server_port string) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", server_port))
	if err != nil {
		log.Printf("Error listening: %s\n", err.Error())
		os.Exit(1)
	}
	for {
		conn, err2 := ln.Accept()
		if err2 != nil {
			log.Printf("Error Connecting: %s\n", err.Error())
		}
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	writer := bufio.NewWriter(os.Stdout)
	message := make([]byte, RECV_BUFFER_SIZE)
	readSofar := RECV_BUFFER_SIZE
	for readSofar == RECV_BUFFER_SIZE {
		readSoFar, _ := conn.Read(message)
		writer.Write(message[0:readSoFar])
		writer.Flush()
	}
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
