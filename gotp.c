#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "gotp.h"

char* hello_world(void) {
    char* message = malloc(32);
    if (message == NULL) {
        return NULL;
    }
    strcpy(message, "Hello from C!");
    return message;
}

char* greet(const char* name) {
    if (name == NULL) {
        return NULL;
    }

    size_t len = strlen("Hello, ") + strlen(name) + strlen("!") + 1;
    char* greeting = malloc(len);
    if (greeting == NULL) {
        return NULL;
    }

    snprintf(greeting, len, "Hello, %s!", name);
    return greeting;
}

void cleanup_string(char* str) {
    if (str != NULL) {
        free(str);
    }
}
