package goforth

import (
	"fmt"
	"log"
)

const Kernel = `
: true -1 ;
: false 0 ;

variable  temp
: swap   >r temp ! r> temp @ ;
: over   >r temp ! temp @ r> temp @ ;
: rot    >r swap r> swap ;

: dup    temp ! temp @ temp @ ;
: 2dup   over over ;
: ?dup   temp ! temp @ if temp @ temp @ then ;

: nip    >r temp ! r> ;

: invert   -1 nand ;
: negate   invert 1 + ;
: -        negate + ;

: 1+   1 + ;
: 1-   -1 + ;
: +!   dup >r @ + r> ! ;
: 0=   if 0 else -1 then ;
: =    - 0= ;
: <>   = 0= ;

: or   invert swap invert nand ;
: xor   2dup nand 1+ dup + + + ;
: and   nand invert ;
: 2*    dup + ;

: <   2dup xor 0< if drop 0< else - 0< then ;
: u<   2dup xor 0< if nip 0< else - 0< then ;
: >   swap < ;
: u>   swap u> ;
`

type Stack []int

func (s *Stack) Push(i int) {
	*s = append(*s, i)
}

func (s *Stack) Pop() (int, bool) {
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, true
}

type Machine struct {
	Dictionary map[string]AST
	Stack      Stack
	RStack     Stack
	Memory     []int
	Addr       int
	Variables  map[string]int
}

func NewMachine(memsize int) *Machine {
	machine := &Machine{
		Memory:     make([]int, 1000+memsize),
		Dictionary: make(map[string]AST),
		Variables:  make(map[string]int),
	}
	kernelToks := Tokenize(Kernel)
	kernel := Parse(kernelToks)
	machine.Evaluate(kernel)
	return machine
}

type BinaryOp func(x, y int) int

func (m *Machine) BinaryOp(op BinaryOp) *int {
	stackLen := len(m.Stack)
	if stackLen < 2 {
		fmt.Printf("Stack underflow %+v", m)
		return &stackLen
	} else {
		x, _ := m.Stack.Pop()
		y, _ := m.Stack.Pop()
		m.Stack.Push(op(x, y))
	}
	return nil
}

func (m *Machine) execNative(word string) int {
	switch word {
	case ".":
		elem, ok := m.Stack.Pop()
		if ok {
			fmt.Printf("%+v\n", elem)
		} else {
			fmt.Println("Stack underflow")
		}
	case "drop":
		m.Stack.Pop()
	case ">r":
		x, _ := m.Stack.Pop()
		m.RStack.Push(x)
	case "r>":
		x, _ := m.RStack.Pop()
		m.Stack.Push(x)
	case "+":
		ret := m.BinaryOp(func(x, y int) int { return x + y })
		if ret != nil {
			return *ret
		}
	case "nand":
		ret := m.BinaryOp(func(x, y int) int { return x &^ y })
		if ret != nil {
			return *ret
		}
	case "!":
		addr, _ := m.Stack.Pop()
		val, _ := m.Stack.Pop()
		m.Memory[addr] = val
	case "@":
		last := len(m.Stack) - 1
		m.Stack[last] = m.Memory[m.Stack[last]]
	case "0<":
		last := len(m.Stack) - 1
		if m.Stack[last] < 0 {
			// true
			m.Stack[last] = -1
		} else {
			// false
			m.Stack[last] = 0
		}
	default:
		return 0
	}
	return 1
}

func (m *Machine) eval(expr Expr) {
	//fmt.Printf("Expr: %+v\n", expr)
	//fmt.Printf("Stack: %+v\n", m.Stack)
	switch expr.Type {
	case ExprPush:
		m.Stack.Push(expr.ValueInt)
	case ExprExec:
		// try native op first then lookup in dictionary then variables
		if result := m.execNative(expr.Value); result == 0 {
			if ast, ok := m.Dictionary[expr.Value]; ok {
				m.Evaluate(ast)
			} else {
				if addr, vok := m.Variables[expr.Value]; vok {
					m.eval(
						Expr{
							Type:     ExprPush,
							ValueInt: addr,
						},
					)
				}
			}
		}
	case ExprDefine:
		m.Dictionary[expr.Value] = expr.Body
	case ExprIfThen:
		cond, _ := m.Stack.Pop()
		if cond == 0 {
			m.Evaluate(expr.AltBody)
		} else {
			m.Evaluate(expr.Body)
		}
	case ExprLoop:
		start, _ := m.Stack.Pop()
		end, _ := m.Stack.Pop()
		m.Dictionary["i"] = []Expr{{Type: ExprPush}}
		for start != end {
			m.Dictionary["i"][0].ValueInt = start
			m.Evaluate(expr.Body)
			start++
		}
		// this is destructive but just don't use `i` outside of loops for now
		m.Dictionary["i"] = nil
	case ExprVariable:
		name := expr.Value
		m.Variables[name] = m.Addr
		m.Addr++
	default:
		log.Fatalf("Don't know how to evaluate expression %+v", expr)
	}
}

func (m *Machine) Evaluate(exprs AST) {
	for _, e := range exprs {
		m.eval(e)
	}
}

func (m *Machine) EvaluateString(code string) {
	toks := Tokenize(code)
	exprs := Parse(toks)
	m.Evaluate(exprs)
}
