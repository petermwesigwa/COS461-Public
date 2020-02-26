/*****************************************************************************
 * client-go.go
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

const SEND_BUFFER_SIZE = 2048

/* TODO: client()
 * Open socket and send message from stdin.
 */
func client(server_ip string, server_port string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", server_ip, server_port))
	defer conn.Close()

	if err != nil {
		log.Panicf("Error connecting: %s\n", err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	message := make([]byte, SEND_BUFFER_SIZE)
	readSoFar := SEND_BUFFER_SIZE
	for readSoFar == SEND_BUFFER_SIZE {
		readSoFar, _ = reader.Read(message)
		conn.Write(message[0:readSoFar])
	}
}

// Main parses command-line arguments and calls client function
func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: ./client-go [server IP] [server port] < [message file]")
	}
	server_ip := os.Args[1]
	server_port := os.Args[2]
	client(server_ip, server_port)
}
