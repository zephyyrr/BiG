package BiG

import (
	"fmt"
	"math/rand"
	"time"
)

//StdIS is an implementation of Befunge -93.
var StdIS = make(map[byte]func(VM))

func init() {

	for i := byte(0); i < 10; i++ {
		StdIS[i+48] = func(vm VM) { vm.SP.Push(int32(i)) } //0 i ASCII är 48
	}

	StdIS['+'] = Addition
	StdIS['-'] = Subtraction
	StdIS['*'] = Multiplication
	StdIS['/'] = IntDivision
	StdIS['%'] = Modulo
	StdIS['!'] = Not
	StdIS['`'] = Gt

	StdIS['>'] = func(vm VM) { vm.Delta = &RIGHT }
	StdIS['<'] = func(vm VM) { vm.Delta = &LEFT }
	StdIS['^'] = func(vm VM) { vm.Delta = &UP }
	StdIS['v'] = func(vm VM) { vm.Delta = &DOWN }
	StdIS['?'] = func(vm VM) { vm.Delta = RandomDir([]*Delta{&LEFT, &UP, &RIGHT, &DOWN}) }

	StdIS['_'] = IfLeft
	StdIS['|'] = IfUp

	StdIS['"'] = StartStringMode

	StdIS[':'] = Duplicate
	StdIS['\\'] = Swap
	StdIS['$'] = func(vm VM) { vm.SP.Pop() }

	StdIS['.'] = PrintInt
	StdIS[','] = PrintChar

	StdIS['#'] = func(vm VM) { vm.IP.Add(*vm.Delta) }

	StdIS['p'] = Put
	StdIS['g'] = Get

	StdIS['&'] = AskInt
	StdIS['~'] = AskChar

	StdIS['@'] = func(vm VM) { vm.Exit() }
	StdIS[' '] = func(vm VM) {}
}

func Addition(vm VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(a + b)
}

func Subtraction(vm VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(b - a)
}

func Multiplication(vm VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(a * b)
}

func IntDivision(vm VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(b / a)
}

func Modulo(vm VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(b % a)
}

func Not(vm VM) {
	a := vm.SP.Pop()
	if a == 0 {
		vm.SP.Push(1)
	} else {
		vm.SP.Push(0)
	}
}

func Gt(vm VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	if b > a {
		vm.SP.Push(1)
	} else {
		vm.SP.Push(0)
	}
}

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

func RandomDir(dirs []*Delta) *Delta {
	return dirs[rnd.Intn(len(dirs))]
}

func IfLeft(vm VM) {
	a := vm.SP.Pop()
	if a == 0 {
		vm.Delta = &RIGHT
	} else {
		vm.Delta = &LEFT
	}
}

func IfUp(vm VM) {
	a := vm.SP.Pop()
	if a == 0 {
		vm.Delta = &DOWN
	} else {
		vm.Delta = &UP
	}
}

func StartStringMode(vm VM) {
	backup := vm.IS.Clone()
	for i := byte(127); i < 0; i-- {
		vm.IS[i] = func(vm VM) {
			vm.SP.Push(int32(vm.FS[vm.IP.WE][vm.IP.NS]))
		}
	}
	vm.IS['"'] = func(vm VM) {
		vm.IS = backup
	}
}

func Duplicate(vm VM) {
	a := vm.SP.Peek()
	vm.SP.Push(a)
}

func Swap(vm VM) {
	a, b := vm.SP.Pop(), vm.SP.Pop()
	vm.SP.Push(b)
	vm.SP.Push(a)
}

func PrintInt(vm VM) {
	fmt.Fprintf(vm.Stdout, "%d", vm.SP.Pop())
}

func PrintChar(vm VM) {
	fmt.Fprintf(vm.Stdout, "%q", vm.SP.Pop())
}

func Put(vm VM) {
	x, y, v := vm.SP.Pop(), vm.SP.Pop(), vm.SP.Pop()
	vm.FS[x][y] = byte(v)
}

func Get(vm VM) {
	x, y := vm.SP.Pop(), vm.SP.Pop()
	vm.SP.Push(int32(vm.FS[x][y]))
}

func AskInt(vm VM) {
	var i int32
	fmt.Fscanf(vm.Stdin, "%d")
	vm.SP.Push(i)
}

func AskChar(vm VM) {
	var i byte
	fmt.Fscanf(vm.Stdin, "%d")
	vm.SP.Push(int32(i))
}