package erlang

// #include <stdlib.h>
// #include <stdio.h>
// #include <string.h>
// #include "gotp.h"
// #include "ei.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// Connection represents an Erlang node connection
type Connection struct {
	ec C.ei_cnode
}

// Config holds connection configuration
type Config struct {
	Cookie        string
	MyNodeName    string
	RemoteNodeName string
	Creation      uint
	Port          int
}

// NewConnection creates a new Erlang connection
func NewConnection(config Config) (*Connection, error) {
	C.ei_init()

	conn := &Connection{}
	
	cookie := C.CString(config.Cookie)
	defer C.free(unsafe.Pointer(cookie))
	
	myNodeName := C.CString(config.MyNodeName)
	defer C.free(unsafe.Pointer(myNodeName))
	
	creation := C.uint(config.Creation)

	if C.ei_connect_init(&conn.ec, myNodeName, cookie, creation) != 0 {
		return nil, fmt.Errorf("ei_connect_init failed: %s", getErlError())
	}

	return conn, nil
}

// Publish publishes a port to epmd
func (c *Connection) Publish(port int) (int, error) {
	pubFd := C.ei_publish(&c.ec, C.int(port))
	if pubFd < 0 {
		return 0, fmt.Errorf("ei_publish failed: %s", getErlError())
	}
	return int(pubFd), nil
}

// Connect connects to a remote Erlang node
func (c *Connection) Connect(remoteNodeName string) (int, error) {
	remoteNodeNameC := C.CString(remoteNodeName)
	defer C.free(unsafe.Pointer(remoteNodeNameC))

	remoteNodeSockFd := C.ei_connect(&c.ec, remoteNodeNameC)
	if remoteNodeSockFd < 0 {
		return 0, fmt.Errorf("ei_connect failed: %s", getErlError())
	}
	return int(remoteNodeSockFd), nil
}

// SendRPC sends an RPC call to a remote process
func (c *Connection) SendRPC(sockFd int, module, function, format string, args ...string) ([]string, error) {
	var request C.ei_x_buff
	C.ei_x_new(&request)
	defer C.ei_x_free(&request)

	formatC := C.CString(format)
	defer C.free(unsafe.Pointer(formatC))

	if len(args) > 0 {
		arg1C := C.CString(args[0])
		defer C.free(unsafe.Pointer(arg1C))
		C.gotp_wrapped_ei_x_format_wo_ver_1(&request, formatC, arg1C)
	} else {
		C.gotp_wrapped_ei_x_format_wo_ver_0(&request, formatC)
	}

	var response C.ei_x_buff
	C.ei_x_new(&response)
	defer C.ei_x_free(&response)

	moduleC := C.CString(module)
	defer C.free(unsafe.Pointer(moduleC))
	functionC := C.CString(function)
	defer C.free(unsafe.Pointer(functionC))

	if C.ei_rpc(&c.ec, C.int(sockFd), moduleC, functionC, request.buff, request.index, &response) != 0 {
		return nil, fmt.Errorf("ei_rpc failed: %s", getErlError())
	}

	var idx C.int
	vals, err := decodeTuple(response.buff, &idx)
	if err != nil {
		return nil, fmt.Errorf("decodeTuple failed: %s", err)
	}

	return vals, nil
}

// getErlError returns the current Erlang error message
func getErlError() string {
	return C.GoString(C.strerror(C.erl_errno))
}

// decodeTuple decodes an Erlang tuple from a buffer
func decodeTuple(buff *C.char, idx *C.int) ([]string, error) {
	var size C.int
	if ret := C.ei_decode_tuple_header(buff, idx, &size); ret < 0 {
		return nil, fmt.Errorf("ei_decode_tuple_header failed (ret: %d): %s", ret, getErlError())
	}

	var vals []string
	for i := 0; i < int(size); i++ {
		var val C.ei_term
		if ret := C.ei_decode_ei_term(buff, idx, &val); ret < 0 {
			return nil, fmt.Errorf("ei_decode_ei_term failed (ret: %d): %s", ret, getErlError())
		}
		switch val.ei_type {
		case C.ERL_ATOM_EXT:
			decodedVal := nullTerminatedBytesToString(val.value[:])
			vals = append(vals, decodedVal)
		case C.ERL_SMALL_TUPLE_EXT, C.ERL_LARGE_TUPLE_EXT:
			subVals, err := decodeTuple(buff, idx)
			if err != nil {
				return nil, err
			}
			vals = append(vals, subVals...)
		default:
			return nil, fmt.Errorf("unknown type: %d", val.ei_type)
		}
	}
	return vals, nil
}

// nullTerminatedBytesToString converts a null-terminated byte array to a Go string
func nullTerminatedBytesToString(b []byte) string {
	return C.GoString((*C.char)(unsafe.Pointer(&b[0])))
}