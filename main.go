package main

// #include <stdlib.h>
// #include <stdio.h>
// #include <string.h>
// #include "gotp.h"
// #include "ei.h"
import "C"

import (
	"bytes"
	"fmt"
	"os"
	"unsafe"

	"github.com/goerlang/etf"
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

	// Use goerlang/etf to decode the response
	responseBytes := C.GoBytes(unsafe.Pointer(response.buff), response.index)
	decodedTerm, err := decodeETFResponse(responseBytes)
	if err != nil {
		return fmt.Errorf("decodeETFResponse failed: %s", err)
	}
	fmt.Println("Response: ", formatETFTerm(decodedTerm))

	fmt.Println("Sent message to remote Erlang node")

	return nil
}

func getErlError() string {
	return C.GoString(C.strerror(C.erl_errno))
}

// decodeETFResponse decodes Erlang External Term Format using goerlang/etf
func decodeETFResponse(data []byte) (etf.Term, error) {
	// Skip the version byte (131) if present
	if len(data) > 0 && data[0] == 131 {
		data = data[1:]
	}
	
	reader := bytes.NewReader(data)
	context := &etf.Context{}
	term, err := context.Read(reader)
	if err != nil {
		return nil, fmt.Errorf("ETF decode failed: %v", err)
	}
	return term, nil
}

// Helper function to convert ETF terms to readable strings
func formatETFTerm(term etf.Term) string {
	switch v := term.(type) {
	case etf.Atom:
		return string(v)
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case etf.Tuple:
		result := "("
		for i, elem := range v {
			if i > 0 {
				result += ", "
			}
			result += formatETFTerm(elem)
		}
		result += ")"
		return result
	case etf.List:
		result := "["
		for i, elem := range v {
			if i > 0 {
				result += ", "
			}
			result += formatETFTerm(elem)
		}
		result += "]"
		return result
	default:
		return fmt.Sprintf("%v", v)
	}
}
