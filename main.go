package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "logging to stderr").Short('v').Bool()
	cmd     = kingpin.Flag("cmd", "command to be executed afterwards").Short('e').String()
	shWrap  = kingpin.Flag("sh", "execute the command within a \"sh\"").Bool()
	bufSize = kingpin.Flag("buf", "buffer size in bytes").Short('b').Default("65536").Int()
	timeout = kingpin.Arg("timeout", "duration until termination, e.g., 5s").Required().Duration()
)

func init() {
	kingpin.CommandLine.Help = "Pipe stdin to stderr, but terminates if nothing has been sent for some time."
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
}

func main() {
	fcu := NewFlowCtrlUnit(os.Stdin, os.Stdout, *timeout, *bufSize)
	fcuErr := fcu.Wait()

	if *verbose {
		fmt.Fprintf(os.Stderr, "quitting, error: %v\n", fcuErr)
	}

	if *cmd != "" {
		var c *exec.Cmd

		if *shWrap {
			c = exec.Command("sh", "-c", *cmd)
		} else {
			cmdParts := strings.Fields(*cmd)
			c = exec.Command(cmdParts[0], cmdParts[1:]...)
		}

		if *verbose {
			fmt.Fprintf(os.Stderr, "exec: %v\n", c)
		}

		if cmdErr := c.Start(); cmdErr != nil && *verbose {
			fmt.Fprintf(os.Stderr, "exec failed: %v\n", cmdErr)
			os.Exit(1)
		}
	}
}
