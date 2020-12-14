# GoForth

A small Forth implemented in Go. Similar in design to my [rust/LLVM forth](https://github.com/bhelx/llvm-forth/) project.
I'm making this to get comfortable writing Go again for work purposes.

## Design

Unlike the LLVM project, this implementation is interpreted and not compiled. The forth machine runs directly
on the AST. There are no intermediate steps or bytecode. To keep things very simple, the machine
only implements 9 primitive operations:

* `drop`
* `@`
* `!`
* `r>`
* `>r`
* `nand`
* `+`
* `<0`
* `.`

All other forth words are dervied from these in the [kernel](https://github.com/bhelx/llvm-forth/#kernelfth) which
is itself written in forth.

## Example

There are 2 flags to the command line program: `file` and `repl`.

Example:
```
go run cmd/goforth/main.go --repl
> 50 8
Stack: [50 8]
> - .
42
Stack: []
> : fib dup 1 > if dup 2 - fib swap 1 - fib + then ;
Stack: []
> : print-fib-numbers 10 0 do i fib . loop ;
Stack: []
> print-fib-numbers
0
1
1
2
3
5
8
13
21
34
Stack: []
> exit
```

fib.fth:
```fth
: fib dup 1 > if dup 2 - fib swap 1 - fib + then ;
: print-fib-numbers 10 0 do i fib . loop ;
print-fib-numbers
```

```
go run cmd/goforth/main.go --file=fib.fth --repl 
0
1
1
2
3
5
8
13
21
34
> 10 fib
Stack: [55]
> exit
```

