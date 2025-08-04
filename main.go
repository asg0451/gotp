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
	myNodeName := C.CString("itest")
	remoteNodeName := C.CString("itestapp@localhost")
	creation := C.uint(1)

	var ec C.ei_cnode
	if C.ei_connect_init(&ec, myNodeName, cookie, creation) != 0 {
		return fmt.Errorf("ei_connect_init failed: %s", getErlError())
	}

	var pubFd C.int
	if pubFd = C.ei_publish(&ec, 9999); pubFd < 0 {
		return fmt.Errorf("ei_publish failed: %s", getErlError())
	}

	fmt.Println("Published port 9999 to epmd")

	var remoteNodeSockFd C.int
	if remoteNodeSockFd = C.ei_connect(&ec, remoteNodeName); remoteNodeSockFd < 0 {
		return fmt.Errorf("ei_connect failed: %s", getErlError())
	}

	fmt.Println("Connected to remote Erlang node")

	// send message to remote node
	// 1. create tuple to send
	var req C.ei_x_buff
	C.ei_x_new_with_version(&req)
	C.ei_x_encode_tuple_header(&req, 2)
	C.ei_x_encode_pid(&req, C.ei_self(&ec))
	C.ei_x_encode_atom(&req, C.CString("Hello world"))
	defer C.ei_x_free(&req)
	// 2. send message
	if C.ei_reg_send(&ec, remoteNodeSockFd, C.CString("ItestElixirApp.Worker"), req.buff, req.index) != 0 {
		return fmt.Errorf("ei_reg_send failed: %s", getErlError())
	}
	fmt.Println("Sent message to remote Erlang node")

	// ei_x_buff buf;
	// ei_x_new_with_version(&buf);
	// ei_x_encode_tuple_header(&buf, 2);
	// ei_x_encode_pid(&buf, ei_self(ec));
	// ei_x_encode_atom(&buf, "Hello world");
	// ei_reg_send(&ec, fd, "my_server", buf.buff, buf.index);

	return nil
}

func getErlError() string {
	return C.GoString(C.strerror(C.erl_errno))
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
