#ifndef GOTP_H
#define GOTP_H

#ifdef __cplusplus
extern "C" {
#endif

char* hello_world(void);
char* greet(const char* name);
void cleanup_string(char* str);

#ifdef __cplusplus
}
#endif

#endif // GOTP_H