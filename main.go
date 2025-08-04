package main

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
	vals, err := decodeTuple(response.buff, &idx)
	if err != nil {
		return fmt.Errorf("decodeTuple failed: %s", err)
	}
	fmt.Println("Response: ", vals)

	fmt.Println("Sent message to remote Erlang node")

	return nil
}

func getErlError() string {
	return C.GoString(C.strerror(C.erl_errno))
}

func decodeTuple(buff *C.char, idx *C.int) ([]string, error) {
	var size C.int
	if ret := C.ei_decode_tuple_header(buff, idx, &size); ret < 0 {
		return nil, fmt.Errorf("ei_decode_tuple_header failed (ret: %d): %s", ret, getErlError())
	}

	var vals []string
	for i := 0; i < int(size); i++ {
		var val C.ei_term
		// note: ret = 1 if the decoded data is in val.value, 0 if it didnt fit / isnt
		if ret := C.ei_decode_ei_term(buff, idx, &val); ret < 0 {
			return nil, fmt.Errorf("ei_decode_ei_term failed (ret: %d): %s", ret, getErlError())
		}
		switch val.ei_type {
		case C.ERL_ATOM_EXT:
			decodedVal := nullTerminatedBytesToString(val.value[:])
			fmt.Println("decoded atom: ", decodedVal)
			vals = append(vals, decodedVal)
		case C.ERL_SMALL_TUPLE_EXT, C.ERL_LARGE_TUPLE_EXT:
			fmt.Println("saw tuple, recursing")
			// or something..?
			subVals, err := decodeTuple(buff, idx)
			if err != nil {
				return nil, err
			}
			// todo not flatten
			vals = append(vals, subVals...)
		default:
			return nil, fmt.Errorf("unknown type: %d", val.ei_type)
		}
	}
	return vals, nil
}

func nullTerminatedBytesToString(b []byte) string {
	return C.GoString((*C.char)(unsafe.Pointer(&b[0])))
}
