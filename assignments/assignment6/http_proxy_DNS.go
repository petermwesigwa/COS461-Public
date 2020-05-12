/*****************************************************************************
 * http_proxy_DNS.go
 * Names: Peter Mwesigwa, Matthew Fastow, Greg Smith
 * NetIds: mwesigwa, mfastow, gpsmith
 *****************************************************************************/

// TODO: implement an HTTP proxy with DNS Prefetching

// Note: it is highly recommended to complete http_proxy.go first, then copy it
// with the name http_proxy_DNS.go, thus overwriting this file, then edit it
// to add DNS prefetching (don't forget to change the filename in the header
// to http_proxy_DNS.go in the copy of http_proxy.go)

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/html"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	resp, err := client.Do(req)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "net/http: use last response") {
			fmt.Printf("Error: %s", err.Error())
			resp := []byte(fmt.Sprintf("HTTP/1.1 500 Internal Server Error: %s\r\n", err.Error()))
			conn.Write(resp)
			return
		}
	}

	// write the response form the remote server to the client.
	err = resp.Write(conn)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		resp := []byte("HTTP/1.1 500 Internal Server Error\r\n")
		conn.Write(resp)
		return
	}
	// DNS prefetching
	resp_dns := resp
	go send_dns(resp_dns.Body)
	// close connection
}

func send_dns(r io.Reader) error {

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return z.Err()
		case html.StartTagToken:

			name, hasToken := z.TagName()

			// check if it is html link tag
			if string(name) == "a" {

				// send dns requests for all links in tag
				for hasToken {
					var key, val []byte
					key, val, hasToken = z.TagAttr()
					if string(key) == "href" && string(val[:4]) == "http" {
						go net.LookupHost(string(val))
					}
				}
			}
		}
	}

}
