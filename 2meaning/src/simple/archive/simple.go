package main

import (
	"fmt"
)

type Reducible interface {
	Reduce() interface{}  // return type might have to change to interface{}
	// ReduceOnce()
	IsReducible() bool
}

/* ---[ Compute Types ]--- */

type number struct {
	value int
}

type add struct {
	left *number
	right *number
}

type multiply struct {
	left *number
	right *number
}

/* ---[ Reduce methods ]--- */

func (n number) Reduce() interface{} {
	return n.value
}

func (a add) Reduce() interface{} {
	numLeft := a.left.Reduce().(number)
	numRight := a.right.Reduce().(number)
	return numLeft.value + numRight.value
}

func (m multiply) Reduce() interface{} {
	numLeft := m.left.Reduce().(number)
	numRight := m.right.Reduce().(number)
	return numLeft.value * numRight.value
}

// func (n number) ReduceOnce() {
// 	// no op
// }

// func (a add) ReduceOnce() {
// 	if a.left.IsReducible() {
// 		a.left = a.left.Reduce()
// 	} else if (a.right.IsReducible()) {
// 		a.right = a.right.Reduce()
// 	} else {
// 		a.Reduce()
// 	}
// }

// func (m multiply) ReduceOnce() {
// 	if m.left.IsReducible() {
// 		m.left = number( m.left.Reduce() )
// 	} else if (m.right.IsReducible()) {
// 		m.right = number( m.right.Reduce() )
// 	} else {
// 		m.Reduce()
// 	}
// }


// Q1: should this be a fn or a method?
// Q2: shoudl this return a number of a *number ?
func NewNumber(value int) *number {
	return &number{value}
}

func NewAdd(left *number, right *number) *add {
	return &add{left, right}
}

func NewMultiply(left *number, right *number) *multiply {
	return &multiply{left, right}
}


func (n *number) String() string {
	return fmt.Sprintf("%d", n.value)
}

func (a add) String() string {
	return fmt.Sprintf("«%v + %v»", a.left, a.right)
}

func (m multiply) String() string {
	return fmt.Sprintf("«%v * %v»", m.left, m.right)
}

func main() {
	n1 := NewNumber(22)
	n2 := NewNumber(33)
	fmt.Println(n1)

	a := NewAdd(n1, n2)
	fmt.Println(a)

	n3 := NewNumber(44)

	m := multiply(n3, a)
	fmt.Println(m)

	/// reduce
	println("--------------------------")
	// fmt.Printf("reduceOnce add : %v\n", a.ReduceOnce())
	// fmt.Printf("reduce add : %v\n", a.Reduce())
	// fmt.Printf("reduce mult: %v\n", m.Reduce())
}
