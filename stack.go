package BiG

type Stack interface {
	Pop() int32
	Push(value int32)
	Peek() int32
	Size() int
}

type LinkedStack struct {
	top  *Element
	size int
}

type Element struct {
	value int32
	next  *Element
}

func (ls LinkedStack) Push(value int32) {
	ls.top = &Element{value, ls.top}
	ls.size++
}

func (ls LinkedStack) Pop() int32 {
	defer func() {
		ls.top = ls.top.next
		ls.size--
	}()
	return ls.Peek()
}

func (ls LinkedStack) Peek() int32 { return ls.top.value }

func (ls LinkedStack) Size() int { return ls.size }
