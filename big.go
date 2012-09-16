package BiG

import (
	"io"
	"os"
	"time"
	"bytes"
	"fmt"
)

const (
	PAGEWIDTH = 80
	PAGEHEIGHT = 25
)

type FungeSpace [PAGEHEIGHT][PAGEWIDTH]byte

func (fs FungeSpace) String() string {
	buff := new(bytes.Buffer)
	for i := 0; i < PAGEHEIGHT; i++ {
		for j := 0; j < PAGEWIDTH; j++ {
			fmt.Fprint(buff, string(fs[i][j]))
		}
		fmt.Fprintln(buff)
	}
	return buff.String()
}

type InstructionPointer struct {
	WE, NS int8
}

func (ip *InstructionPointer) Add(d Delta) {
	ip.NS += d[1]
	ip.WE += d[0]
	if ip.NS < 0 {
		ip.NS = PAGEHEIGHT + ip.NS
	}
	if ip.NS > 24 {
		ip.NS = ip.NS % PAGEHEIGHT
	}

	if ip.WE < 0 {
		ip.WE = PAGEWIDTH + ip.WE
	}
	if ip.WE > 24 {
		ip.WE = ip.WE % PAGEWIDTH
	}
}

type Delta [2]int8

var (
	LEFT  = Delta{-1, 0}
	UP    = Delta{0, -1}
	RIGHT = Delta{1, 0}
	DOWN  = Delta{0, 1}
)

type InstructionSet map[byte]func(*VM)

func (old InstructionSet) Clone() (newIS InstructionSet) {
	newIS = make(InstructionSet)
	for k, v := range old {
		newIS[k] = v
	}
	return
}

type VM struct {
	FS       *FungeSpace
	IP       *InstructionPointer
	IS       InstructionSet
	Delta    *Delta
	SP       Stack
	quitting bool
	Stdin    io.Reader
	Stdout   io.Writer
}

func NewVM(fs *FungeSpace) (vm *VM) {
	vm = new(VM)
	vm.Delta = &RIGHT
	vm.quitting = false
	vm.SP = NewStack()
	vm.IS = StdIS
	vm.IP = &InstructionPointer{0, 0}
	vm.FS = fs
	vm.Stdin = os.Stdin
	vm.Stdout = os.Stdout
	return
}

func (vm *VM) Tick() {
	//fmt.Println("Tick!")
	f := vm.IS[vm.FS[vm.IP.NS][vm.IP.WE]]
	if f != nil {
		f(vm)
	}
	vm.IP.Add(*vm.Delta)
}

func (vm *VM) Run(ticker *time.Ticker) {
	for _ = range ticker.C {
		if vm.quitting {
			break
		}
		vm.Tick()
	}
}

func (vm VM) Done() bool {
	return vm.quitting
}

func (vm *VM) Exit() {
	vm.quitting = true
}
