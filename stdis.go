package BiG

import (
	"fmt"
	"math/rand"
	"time"
)

//StdIS is an implementation of Befunge -93.
var StdIS = make(map[int8]func(*VM))

func init() {

	for i := int8(0); i < 10; i++ {
		StdIS[i+48] = Pusher(i)  //0 i ASCII är 48
	}

	StdIS['+'] = Addition
	StdIS['-'] = Subtraction
	StdIS['*'] = Multiplication
	StdIS['/'] = IntDivision
	StdIS['%'] = Modulo
	StdIS['!'] = Not
	StdIS['`'] = Gt

	StdIS['>'] = func(vm *VM) { vm.Delta = &RIGHT }
	StdIS['<'] = func(vm *VM) { vm.Delta = &LEFT }
	StdIS['^'] = func(vm *VM) { vm.Delta = &UP }
	StdIS['v'] = func(vm *VM) { vm.Delta = &DOWN }
	StdIS['?'] = func(vm *VM) { vm.Delta = RandomDir([]*Delta{&LEFT, &UP, &RIGHT, &DOWN}) }

	StdIS['_'] = IfLeft
	StdIS['|'] = IfUp

	StdIS['"'] = QuickStringMode

	StdIS[':'] = Duplicate
	StdIS['\\'] = Swap
	StdIS['$'] = func(vm *VM) { vm.SP.Pop() }

	StdIS['.'] = PrintInt
	StdIS[','] = PrintChar

	StdIS['#'] = func(vm *VM) { vm.IP.Add(*vm.Delta) }

	StdIS['p'] = Put
	StdIS['g'] = Get

	StdIS['&'] = AskInt
	StdIS['~'] = AskChar

	StdIS['@'] = func(vm *VM) { vm.Exit() }
	StdIS[' '] = func(vm *VM) {}
}

func Pusher(i int8) func(*VM) {
	return func(vm *VM) {vm.SP.Push(int32(i)) }
	}

func Addition(vm *VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(a + b)
}

func Subtraction(vm *VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(b - a)
}

func Multiplication(vm *VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(a * b)
}

func IntDivision(vm *VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(b / a)
}

func Modulo(vm *VM) {
	a := vm.SP.Pop()
	b := vm.SP.Pop()
	vm.SP.Push(b % a)
}

func Not(vm *VM) {
	a := vm.SP.Pop()
	if a == 0 {
		vm.SP.Push(1)
	} else {
		vm.SP.Push(0)
	}
}

func Gt(vm *VM) {
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

func IfLeft(vm *VM) {
	a := vm.SP.Pop()
	if a == 0 {
		vm.Delta = &RIGHT
	} else {
		vm.Delta = &LEFT
	}
}

func IfUp(vm *VM) {
	a := vm.SP.Pop()
	if a == 0 {
		vm.Delta = &DOWN
	} else {
		vm.Delta = &UP
	}
}

func StartStringMode(vm *VM) {
	backup := vm.IS.Clone()
	for i := int8(127); i > 0; i-- {
		vm.IS[i] = func(vm *VM) {
			//fmt.Println("PUSHING:", vm.FS[vm.IP.NS][vm.IP.WE])
			vm.SP.Push(int32(vm.FS[vm.IP.NS][vm.IP.WE]))
		}
	}
	vm.IS['"'] = func(vm *VM) {
		vm.IS = backup
	}
}

func QuickStringMode(vm *VM) {
	vm.IP.Add(*vm.Delta)
	for next := vm.FS[vm.IP.NS][vm.IP.WE];next != '"'; next = vm.FS[vm.IP.NS][vm.IP.WE] {
		vm.SP.Push(int32(next))
		vm.IP.Add(*vm.Delta)
	}
}

func Duplicate(vm *VM) {
	a := vm.SP.Peek()
	vm.SP.Push(a)
}

func Swap(vm *VM) {
	a, b := vm.SP.Pop(), vm.SP.Pop()
	vm.SP.Push(a)
	vm.SP.Push(b)
}

func PrintInt(vm *VM) {
	fmt.Fprintf(vm.Stdout, "%d ", vm.SP.Pop())
}

func PrintChar(vm *VM) {
	vm.Stdout.Write([]byte{byte(vm.SP.Pop())})
}

func Put(vm *VM) {
	y, x, v := vm.SP.Pop(), vm.SP.Pop(), vm.SP.Pop()
	vm.FS[y][x] = int8(v)
}

func Get(vm *VM) {
	y, x := vm.SP.Pop(), vm.SP.Pop()
	if y < PAGEHEIGHT && y >= 0 && x < PAGEWIDTH && x >= 0 {
		vm.SP.Push(int32(vm.FS[y][x]))
	} else {
		vm.SP.Push(0)
	}
}

func AskInt(vm *VM) {
	var i int32
	fmt.Fscanf(vm.Stdin, "%d", &i)
	vm.SP.Push(i)
}

func AskChar(vm *VM) {
	var i byte
	fmt.Fscanf(vm.Stdin, "%c", &i)
	vm.SP.Push(int32(i))
}
