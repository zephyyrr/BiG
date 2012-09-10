package big

type FungeSpace [25][80]byte

type InstructionPointer struct {
	NS, WE uint8
}

func (ip *InstructionPointer) Add(d Delta) {
	ip.NS += d.NS
	ip.WE += d.WE
	if ip.NS > 24 {
		ip.NS = ip.NS % 25
	}
	if ip.WE > 24 {
		ip.WE = ip.WE % 80
	}
}

type Delta struct {
	NS, WE int8
}

type InstructionSet map[byte]func(VM)

type VM stuct {
	FS *FungeSpace
	IP *InstructionPointer
	IS InstructionSet
	Delta *Delta
	SP Stack
	quitting bool
}

func (vm VM) Tick() {
	vm.IS[vm.FS[vm.IP.NS][vm.IP.WE]](vm)
	vm.IP.Add(vm.delta)
}

func (vm VM) Run() {
	for !vm.quitting {
		vm.Tick()
	}
}

func (vm VM) Exit() {
	vm.quitting = true
}

var StdIS := make(map[byte]func(VM))

func init() {
	for i := 0; i < 10 ; i++ {
		StdIS[i + 48] = func(vm VM) {vm.SP.Push(i)} //0 i ASCII är 48
	}
	
	StdIS['+'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() + vm.SP.Pop())}
	StdIS['-'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() - vm.SP.Pop())}
	StdIS['*'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() * vm.SP.Pop())}
	StdIS['/'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() / vm.SP.Pop())}
	StdIS['%'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() % vm.SP.Pop())}
	StdIS['!'] = func(vm VM) {if vm.SP.Pop() > 0 {vm.SP.Push(0)} else {vm.SP.Push(1)}}
	StdIS['@'] = func(vm VM) {vm.Exit()}
}
