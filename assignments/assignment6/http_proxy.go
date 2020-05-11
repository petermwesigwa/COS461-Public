/*****************************************************************************
 * http_proxy.go
 * Names:
 * NetIds:
 *****************************************************************************/

// TODO: implement an HTTP proxy

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

func main() {
	var port string
	if len(os.Args) < 2 {
		port = "80"
	} else {
		port = os.Args[1]
	}

	// listen for connection
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on port %s\n", port)

	for {
		// accept a connection
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		// spawn go routine to handle connection
		go handle_http(conn)
	}

}

func handle_http(conn net.Conn) {

	// close the connection once done
	defer conn.Close()

	// create bufio reader
	r := bufio.NewReader(conn)

	// read data into request data type
	req, err := http.ReadRequest(r)
	if err != nil {
		resp := []byte("HTTP/1.1 500 Internal Server Error\r\n")
		conn.Write(resp)
		return
	}

	if req.Method != "GET" {
		resp := []byte("HTTP/1.1 500 Internal Server Error\r\n")
		conn.Write(resp)
		return
	}

	// add connection = close
	req.Header.Del("Connection")
	req.Header.Add("Connection", "close")

	// have to set req.URL not req.RequestURI for http.Client to work.
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		resp := []byte("HTTP/1.1 500 Internal Server Error\r\n")
		conn.Write(resp)
		return
	}
	req.URL = u
	req.RequestURI = ""

	// send request
	client := &http.Client{
		// Avoid following redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("net/http: use last response")
		},
	}
	var RedirectAttemptedError = errors.New("net/http: use last response")
	resp, err := client.Do(req)
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		err = nil
	}
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		resp := []byte(fmt.Sprintf("HTTP/1.1 500 Internal Server Error: %s\r\n", err.Error()))
		conn.Write(resp)
		return
	}

	// write the response form the remote server to the client.
	err = resp.Write(conn)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		resp := []byte("HTTP/1.1 500 Internal Server Error\r\n")
		conn.Write(resp)
		return
	}
	// close connection
}

/*
 func client(server_ip string, server_port string) {

  b := make([]byte, SEND_BUFFER_SIZE)
  server := server_ip + ":" + server_port

  conn, err := net.Dial("tcp", server)
  if err != nil {
    log.Fatal(err)
  }

  r := bufio.NewReader(os.Stdin)

  n, err := r.Read(b)
  for n == SEND_BUFFER_SIZE {
    conn.Write(b[:n])
    n, err = r.Read(b)
  }
  conn.Write(b[:n])

}


func server(server_port string) {

  b := make([]byte, RECV_BUFFER_SIZE)

  ln, err := net.Listen("tcp", ":" + server_port)
  if err != nil {
    log.Fatal(err)
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      continue
    }

    n, err := conn.Read(b)
    for n == RECV_BUFFER_SIZE {
      //fmt.Printf(string(b[:n]))
      os.Stdout.Write(b[:n])
      n, err = conn.Read(b)
    }
    //fmt.Printf(string(b[:n]))
    os.Stdout.Write(b[:n])


  }

}

*/
