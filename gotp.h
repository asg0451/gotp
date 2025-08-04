#ifndef GOTP_H
#define GOTP_H

#ifdef __cplusplus
extern "C" {
#endif

#include "ei.h"



// Simple inline wrappers for variadic functions (CGO doesn't support variadic)
static inline int gotp_wrapped_ei_x_format_wo_ver_0(ei_x_buff* x, const char *fmt) {
    return ei_x_format_wo_ver(x, fmt);
}

static inline int gotp_wrapped_ei_x_format_wo_ver_1(ei_x_buff* x, const char *fmt, const char *arg1) {
    return ei_x_format_wo_ver(x, fmt, arg1);
}

static inline int gotp_wrapped_ei_x_format_wo_ver_2(ei_x_buff* x, const char *fmt, const char *arg1, const char *arg2) {
    return ei_x_format_wo_ver(x, fmt, arg1, arg2);
}

static inline int gotp_wrapped_ei_x_format_wo_ver_3(ei_x_buff* x, const char *fmt, const char *arg1, const char *arg2, const char *arg3) {
    return ei_x_format_wo_ver(x, fmt, arg1, arg2, arg3);
}


#ifdef __cplusplus
}
#endif

#endif // GOTP_H
