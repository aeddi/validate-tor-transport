package main

import (
	"fmt"
	"os"
)

// +build linux darwin

func main() {
	fmt.Println(do(os.Stderr))
}
