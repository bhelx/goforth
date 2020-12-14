package main

import (
	"flag"
	"fmt"
	"github.com/bhelx/goforth"
	"io/ioutil"
)

func main() {
	repl := flag.Bool("repl", false, "Load a REPL")
	filename := flag.String("file", "", "Load this fth file")
	flag.Parse()

	machine := goforth.NewMachine(1000)

	if *filename != "" {
		b, err := ioutil.ReadFile(*filename)
		if err != nil {
			fmt.Print(err)
		}
		machine.EvaluateString(string(b))
	}

	if *repl {
		i := &goforth.Interpreter{Machine: machine}
		i.Start()
	}
}
