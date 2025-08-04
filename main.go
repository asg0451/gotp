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
	"unsafe"
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

	// TODO: why proper way no work?
	// // send message to remote node
	// // 1. create tuple to send
	// var req C.ei_x_buff
	// C.ei_x_new_with_version(&req)
	// C.ei_x_encode_tuple_header(&req, 2)
	// C.ei_x_encode_pid(&req, C.ei_self(&ec))
	// C.ei_x_encode_atom(&req, C.CString("Hello world"))
	// defer C.ei_x_free(&req)
	// // 2. send message
	// if C.ei_reg_send(&ec, remoteNodeSockFd, C.CString("ItestElixirApp.Worker"), req.buff, req.index) != 0 {
	// 	return fmt.Errorf("ei_reg_send failed: %s", getErlError())
	// }

	// workaround from https://stackoverflow.com/questions/57893503/send-message-to-genserver-from-c
	var request C.ei_x_buff
	C.ei_x_new(&request)
	C.gotp_wrapped_ei_x_format_wo_ver_1(&request, C.CString("[~s]"), C.CString("Hello world"))
	defer C.ei_x_free(&request)
	var response C.ei_x_buff
	C.ei_x_new(&response)
	defer C.ei_x_free(&response)
	if C.ei_rpc(&ec, remoteNodeSockFd, C.CString("ItestElixirApp.Worker"), C.CString("send_to_self"), request.buff, request.index, &response) != 0 {
		return fmt.Errorf("ei_rpc failed: %s", getErlError())
	}

	var idx C.int
	var size C.int
	if C.ei_decode_tuple_header(response.buff, &idx, &size) != 0 {
		return fmt.Errorf("ei_decode_ei_term failed: %s", getErlError())
	}

	var vals []string
	for i := 0; i < int(size); i++ {
		var val C.ei_term
		if C.ei_decode_ei_term(response.buff, &idx, &val) != 1 {
			return fmt.Errorf("ei_decode_ei_term failed: %s", getErlError())
		}
		// TODO: actually value is a union; see https://sunzenshen.github.io/tutorials/2015/05/09/cgotchas-intro.html and https://www.erlang.org/docs/23/man/ei#ei_term
		switch val.ei_type {
		case C.ERL_ATOM_EXT:
			decodedVal := nullTerminatedBytesToString(val.value[:])
			fmt.Println("decoded atom: ", decodedVal)
			vals = append(vals, decodedVal)
		case C.ERL_SMALL_TUPLE_EXT, C.ERL_LARGE_TUPLE_EXT:
			// TODO recurse
		default:
			fmt.Println("Unknown type: ", val.ei_type)
			continue
		}

		vals = append(vals, string(val.value[:val.size]))
	}
	fmt.Println("Response: ", vals)

	// print response
	fmt.Println("Response: ", C.GoString(response.buff))

	fmt.Println("Sent message to remote Erlang node")

	return nil
}

func getErlError() string {
	return C.GoString(C.strerror(C.erl_errno))
}

func nullTerminatedBytesToString(b []byte) string {
	return C.GoString((*C.char)(unsafe.Pointer(&b[0])))
}
