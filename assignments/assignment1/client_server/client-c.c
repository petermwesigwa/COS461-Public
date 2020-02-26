/*****************************************************************************
 * client-c.c                                                                 
 * Name:
 * NetId:
 *****************************************************************************/

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netdb.h>
#include <netinet/in.h>
#include <errno.h>

#define SEND_BUFFER_SIZE 2048


/* TODO: client()
 * Open socket and send message from stdin.
 * Return 0 on success, non-zero on failure
*/
int client(char *server_ip, char *server_port) {
  int yes = 1;
  int bytes_read = SEND_BUFFER_SIZE;
  int len, status, sockfd, i;
  struct addrinfo hints, *res;
  char buf[SEND_BUFFER_SIZE];
  char c;

  // set IP andtransport protocol
  memset(&hints, 0, sizeof hints);
  hints.ai_family = AF_UNSPEC;
  hints.ai_socktype = SOCK_STREAM;

  // store address information for server
  if ((status = getaddrinfo(server_ip, server_port, &hints, &res)) != 0) {
    fprintf(stderr, "getaddrinfo error: %s\n", gai_strerror(status));
    exit(1);
  }

  // create the socket 
  sockfd = socket(res->ai_family, res->ai_socktype, res->ai_protocol);
  if (sockfd == -1) {
    fprintf(stderr, "Error creating socket: %s\n", strerror(errno));
    exit(1);
  }

  // connect to the socket
  if (connect(sockfd, res->ai_addr, res->ai_addrlen) == -1) {
    close(sockfd);
    fprintf(stderr, "Error Connecting: %s\n", strerror(errno));
    exit(EXIT_FAILURE);
  }

  // read data from stdin and transmit it to the socket
  while (bytes_read == SEND_BUFFER_SIZE) {
    bytes_read = read(0, buf, SEND_BUFFER_SIZE);
    if (send(sockfd, buf, bytes_read, 0) == -1) {
      fprintf(stderr, "Error writing to socket: %s\n", strerror(errno));
    } else {
      for (i=0; i < bytes_read; i++)
      {
        printf("%c", buf[i]);
      }
    }
  }

  // cleanup and exit
  freeaddrinfo(res);
  close(sockfd);
  return 0;
}

/*
 * main()
 * Parse command-line arguments and call client function
*/
int main(int argc, char **argv) {
  char *server_ip;
  char *server_port;

  if (argc != 3) {
    fprintf(stderr, "Usage: ./client-c [server IP] [server port] < [message]\n");
    exit(EXIT_FAILURE);
  }

  server_ip = argv[1];
  server_port = argv[2];
  return client(server_ip, server_port);
}
