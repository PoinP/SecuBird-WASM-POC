//go:build js

package main

import (
	"fmt"
	_ "wasm-test/js"
)

func main() {
	fmt.Println("SDK is ready to be used!")
	<-make(chan struct{})
}
