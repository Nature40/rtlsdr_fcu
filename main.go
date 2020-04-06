package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fcu := NewFlowCtrlUnit(os.Stdin, os.Stdout, time.Second)
	fmt.Fprintf(os.Stderr, "err: %v\n", fcu.Wait())
}
