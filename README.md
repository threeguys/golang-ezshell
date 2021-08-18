# golang-ezshell
##### A Golang library for creating simple interactive shells

An extremely lightweight and simple library for creating interactive shells, golang-ezshell is my attempt
to create a library that does what I need it to do, with little extra fluff. I often build tooling based
on an interactive shell as part of testing out new technology. Some languages have great support or libraries
to do this but my (very brief) search turned up libraries that were more geared to glitz with color
support and advanced completion schemes. I was building an epub book reader library and learning about
the format and a shell seemed the best way to go about it (since they're just zip files) and thus
golang-ezshell was born.

For now, the comments are sparse and documentation minimal as I have built it to do what I need and am
not spending a lot of time on it. I will circle back around.

# Building
This library does use one of my other golang libs, [golang-toolkit](https://github.com/threeguys/golang-toolkit)
for assertion support in testing, however it has no runtime dependencies. You can install the module by
running:

```
go get github.com/threeguys/golang-ezshell
```

# Examples
See the [examples](examples) directory for a couple of simple examples. The [hello-ezshell](examples/hello-ezshell/main.go)
is a very short, quick start on implementing a couple of commands to get you started. The [ezbash](examples/ezbash/main.go)
example is a bit more in-depth, illustrating the ability to use different modes in order to have varying
sets of commands available, depending on the mode.

The most basic example is a null shell, with only the exit command available:

```go
package main

import (
	"github.com/threeguys/golang-ezshell/shell"
	"log"
	"os"
)

func main() {
	sh := shell.NewCommandShell("$ ", []*shell.Command{
		{
			Name:        "exit",
			Description: "quits the shell",
			Flags:       0,
			Handler:     func(_ []string) error { os.Exit(0); return nil },
		},
	})
	log.Fatal(sh.Run())
}
```

# Contributing
Please feel free to post PRs and bugs, I will be responsive and fix what seems to be a problem (because it
will affect me too).

# License
Released under [Apache 2.0](LICENSE)
 