package BiG

import "fmt"

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

func NewStack() Stack {
	return &LinkedStack{nil, 0}
}

func (ls *LinkedStack) String() string {
	s := ""
	e := ls.top
	for i := 0; i < ls.Size(); i++ {
		s += fmt.Sprintf(" %d|", e.value)
		e = e.next
	}
	return s
}

type Element struct {
	value int32
	next  *Element
}

func (ls *LinkedStack) Push(value int32) {
	ls.top = &Element{value, ls.top}
	ls.size++
	//fmt.Println("PUSH:", ls)
}

func (ls *LinkedStack) Pop() int32 {
	if ls.Size() == 0 {
		return 0
	}
	defer func() {
		ls.top = ls.top.next
		ls.size--
	}()
	//fmt.Println("POP:", ls)
	return ls.Peek()
}

func (ls *LinkedStack) Peek() int32 { 
	if ls.Size() == 0 {
		return 0
	}
	return ls.top.value
}

func (ls LinkedStack) Size() int { return ls.size }
