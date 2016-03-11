package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/alankm/vorteil/core"
)

func main() {
	if l := len(os.Args); l != 2 {
		fmt.Fprintf(os.Stderr, "usage: vorteil config\n")
		return
	}
	vorteil, err := core.New(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	err = vorteil.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	vorteil.Stop()
}
