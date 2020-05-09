/*****************************************************************************
 * http_proxy.go                                                                 
 * Names: 
 * NetIds:
 *****************************************************************************/

 // TODO: implement an HTTP proxy

 package main

 import (
 	"os"
 	"net"
 	"net/http"
 	"bufio"
 	"log"
 	//"io"
 	//"fmt",
 	//"bytes"
 )

 func main() {
 	var port string
 	if len(os.Args) < 2 {
    	port = "80"
  	} else {
  		port = os.Args[1]
  	}

  	// listen for connection
	ln, err := net.Listen("tcp", ":" + port)
  	if err != nil {
    	log.Fatal(err)
  	}

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

 	BUFFER_SIZE := 100000

 	// read data into buffer
  rb := make([]byte, BUFFER_SIZE)
  wb := make([]byte, BUFFER_SIZE)
 	_, err := conn.Read(rb)

  // write to file
  fw, err := os.Create("/http")
  fw.Write(rb)

  // create bufio reader
  fr, err := os.Open("/http")
 	r := bufio.NewReader(fr)

 	// read data into request data type
 	req, err := http.ReadRequest(r)
 	if err != nil || req.Method != "GET" {
 		resp := []byte("HTTP/1.1 500 Internal Server Error\r\n")
 		conn.Write(resp)
 		conn.Close()
 		return
 	}

 	// add connection = close
 	req.Header.Del("Connection")
 	req.Header.Add("Connection", "close")

 	// send request
 	client := &http.Client{}
	resp, err := client.Do(req)

	// write to buffer and then to connection
	err = resp.Write(fw)
  fr.Read(wb)
	conn.Write(wb)

	// close connection
	conn.Close()

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

