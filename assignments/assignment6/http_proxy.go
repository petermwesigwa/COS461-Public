/*****************************************************************************
 * http_proxy.go                                                                 
 * Names: 
 * NetIds:
 *****************************************************************************/

 // TODO: implement an HTTP proxy

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