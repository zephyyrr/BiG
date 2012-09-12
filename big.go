package BiG

type FungeSpace [25][80]byte

type InstructionPointer struct {
	WE, NS int8
}

func (ip *InstructionPointer) Add(d Delta) {
	ip.NS += d[1]
	ip.WE += d[0]
	if ip.NS < 0 {
		ip.NS = 25+ip.NS
	}
	if ip.NS > 24 {
		ip.NS = ip.NS % 25
	}
	
	if ip.WE < 0 {
		ip.WE = 25+ip.WE
	}
	if ip.WE > 24 {
		ip.WE = ip.WE % 80
	}
}

type Delta [2]int8

var (
	LEFT = Delta{-1, 0}
	UP = Delta{0, -1}
	RIGHT = Delta{1, 0}
	DOWN = Delta{0, 1}
)

type InstructionSet map[byte]func(VM)

type VM struct {
	FS *FungeSpace
	IP *InstructionPointer
	IS InstructionSet
	Delta *Delta
	SP Stack
	quitting bool
}

func (vm VM) Tick() {
	vm.IS[vm.FS[vm.IP.NS][vm.IP.WE]](vm)
	vm.IP.Add(*vm.Delta)
}

func (vm VM) Run(ticker <-chan bool) {
	for _ = range ticker {
		if vm.quitting {
			break
		}
		vm.Tick()
	}
}

func (vm VM) Exit() {
	vm.quitting = true
}

var StdIS = make(map[byte]func(VM))

func init() {
	for i := byte(0); i < 10 ; i++ {
		StdIS[i + 48] = func(vm VM) {vm.SP.Push(int32(i))} //0 i ASCII är 48
	}
	
	StdIS['<'] = func(vm VM) {vm.Delta = &LEFT}
	StdIS['^'] = func(vm VM) {vm.Delta = &UP}
	StdIS['>'] = func(vm VM) {vm.Delta = &RIGHT}
	StdIS['v'] = func(vm VM) {vm.Delta = &DOWN}
	
	StdIS['+'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() + vm.SP.Pop())}
	StdIS['-'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() - vm.SP.Pop())}
	StdIS['*'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() * vm.SP.Pop())}
	StdIS['/'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() / vm.SP.Pop())}
	StdIS['%'] = func(vm VM) {vm.SP.Push(vm.SP.Pop() % vm.SP.Pop())}
	StdIS['!'] = func(vm VM) { if vm.SP.Pop() == 0 {vm.SP.Push(1)} else {vm.SP.Push(0)} }
	StdIS['`'] = func(vm VM) {
		a, b := vm.SP.Pop(), vm.SP.Pop()
		if b > a {vm.SP.Push(1)} else {vm.SP.Push(0)}
	}
	StdIS['_'] = func(vm VM) {if vm.SP.Pop() == 0 {vm.Delta = &RIGHT} else {vm.Delta = &LEFT}}
	StdIS['|'] = func(vm VM) {if vm.SP.Pop() == 0 {vm.Delta = &DOWN} else {vm.Delta = &UP}}
	
	
	StdIS['#'] = func(vm VM) {vm.IP.Add(*vm.Delta)}
	StdIS['@'] = func(vm VM) {vm.Exit()}
}
