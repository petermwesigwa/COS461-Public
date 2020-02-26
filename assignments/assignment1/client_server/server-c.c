/*****************************************************************************
 * server-c.c                                                                 
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

#define QUEUE_LENGTH 10
#define RECV_BUFFER_SIZE 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 * Return 0 on success, non-zero on failure
*/
int server(char *server_port) {
  int status, sin_size, sockfd, new_socket, i;
  int yes = 1;
  int bytes_received = RECV_BUFFER_SIZE;
  char buf[RECV_BUFFER_SIZE];
  struct sockaddr_storage their_addr;
  struct addrinfo hints, *res;
  char *server_ip = "127.0.0.1";

  memset(&hints, 0, sizeof hints);
  hints.ai_family = AF_UNSPEC;
  hints.ai_socktype = SOCK_STREAM;

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

  // clear up "Address already in use" error
  if (setsockopt(sockfd, SOL_SOCKET, SO_REUSEADDR, &yes, sizeof yes) == -1) {
    perror("setsockopt");
    exit(1);
  } 

  // bind the socket
  if(bind(sockfd, res->ai_addr, res->ai_addrlen) == -1) {
    fprintf(stderr, "Error binding to socket: %s\n", strerror(errno));
    exit(1);
  }

  if (listen(sockfd, QUEUE_LENGTH) == -1) {
    fprintf(stderr, "Error listening %s\n", strerror(errno));
    exit(1);
  }

  while (1) {
    sin_size = sizeof their_addr;
    new_socket = accept(sockfd, (struct sockaddr *)&their_addr, &sin_size);
    if (new_socket == -1) {
        fprintf(stderr, "Error accepting connection: %s\n", strerror(errno));
        continue;
    }

    while (bytes_received == RECV_BUFFER_SIZE) {
      bytes_received = recv(new_socket, buf, RECV_BUFFER_SIZE, 0);
      for(i=0; i < bytes_received; i++) {
        printf("%c", buf[i]);
      }
      fflush(stdout);
    }

    // clean up
    close(new_socket);
    bytes_received = RECV_BUFFER_SIZE; // gotta reset this so the while loop works again
  } 
  close(sockfd);
  freeaddrinfo(res);
  return 0;
}

/*
 * main():
 * Parse command-line arguments and call server function
*/
int main(int argc, char **argv) {
  char *server_port;

  if (argc != 2) {
    fprintf(stderr, "Usage: ./server-c [server port]\n");
    exit(EXIT_FAILURE);
  }

  server_port = argv[1];
  return server(server_port);
}
