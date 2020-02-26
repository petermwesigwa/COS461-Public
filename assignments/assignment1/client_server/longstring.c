#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv) {
    int length = atoi(argv[1]);
    int i;

    for (i = 0; i < length; i++) {
        // printf("%c", (char) 48 + rand() % 74);
        // printf("x");
        // printf("%c", (char) 48 + i % 74);
        printf("%d ", i);
    }

    printf("xx\nx");
}