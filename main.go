package main

// #cgo CFLAGS: -Wall -g
// // #cgo LDFLAGS: -L. -lgotp
// #include <stdlib.h>
// #include "gotp.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("CGO Skeleton Example")

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
