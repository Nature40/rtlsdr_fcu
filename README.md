# `fcu`, flow control unit

Simple command line util to pipe `stdin` to `stdout`, but terminates if nothing
has been sent for some time. The exiting is accompanied by the closing of both
devices or also by executing some command.

```
$ ./fcu --help
usage: fcu [<flags>] <timeout>

Pipe stdin to stderr, but terminates if nothing has been sent for some time.

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
  -e, --cmd=CMD  command to be executed afterwards
      --sh       execute the command within a "sh"
  -v, --verbose  logging to stderr

Args:
  <timeout>  duration until termination, e.g., 5s
```


## Build

1. Install [Go](https://golang.org/), version 1.11 or later
2. Clone and change to this repository
3. `go build -o fcu`
