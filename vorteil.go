package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/alankm/vorteil/vorteil"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: vorteil config_file\n")
		return
	}

	vorteil, err := vorteil.New(os.Args[1])
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
