package main

// #cgo CFLAGS: -I/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/include -Wall -g
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/lib -L. -lei -lpthread
// #include <stdlib.h>
// #include <stdio.h>
// #include <string.h>
// #include "gotp.h"
// #include "ei.h"
import "C"

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("CGO Skeleton Example")

	C.ei_init()

	cookie := C.CString("super_secret")
	nodeName := C.CString("itestapp@localhost")
	creation := C.uint(1)

	var ec C.ei_cnode
	if C.ei_connect_init(&ec, nodeName, cookie, creation) != 0 {
		// TODO: check errors: https://www.erlang.org/docs/20/man/erl_error
		el := C.GoString(C.strerror(C.erl_errno))
		return fmt.Errorf("ei_connect_init failed: %s", el)
	}

	var sockfd C.int
	if sockfd = C.ei_connect(&ec, nodeName); sockfd < 0 {
		el := C.GoString(C.strerror(C.erl_errno))
		return fmt.Errorf("ei_connect failed: %s", el)
	}

	fmt.Println("Connected to Erlang node")

	return nil
}

// result := C.hello_world()
// goResult := C.GoString(result)
// fmt.Printf("C function returned: %s\n", goResult)

// // Example with parameter
// input := C.CString("Go")
// defer C.free(unsafe.Pointer(input))

// greeting := C.greet(input)
// goGreeting := C.GoString(greeting)
// fmt.Printf("Greeting: %s\n", goGreeting)
