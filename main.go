package main

// #cgo CFLAGS: -I/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/include -Wall -g
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/lib -L. -lei -lpthread
// #include <stdlib.h>
// #include "gotp.h"
// #include "ei.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("CGO Skeleton Example")

	C.ei_init()

	// Example usage of C function
	result := C.hello_world()
	goResult := C.GoString(result)
	fmt.Printf("C function returned: %s\n", goResult)

	// Example with parameter
	input := C.CString("Go")
	defer C.free(unsafe.Pointer(input))

	greeting := C.greet(input)
	goGreeting := C.GoString(greeting)
	fmt.Printf("Greeting: %s\n", goGreeting)
}
