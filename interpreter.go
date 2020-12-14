package goforth

import (
	"bufio"
	"fmt"
	"os"
)

type Interpreter struct {
	Machine *Machine
}

func (i *Interpreter) Start() {
	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			s := string(sentence)
			if s == "exit\n" {
				break
			}
			toks := Tokenize(string(sentence))
			exprs := Parse(toks)
			i.Machine.Evaluate(exprs)
			fmt.Printf("Stack: %+v\n", i.Machine.Stack)
		}
	}
}
