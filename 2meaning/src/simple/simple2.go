package main

import (
	"fmt"
)

type Reducible interface {
	Reduce() Reducible
	ReduceOnce() Reducible
	IsReducible() bool
}


/* ---[ Compute Types ]--- */

type Number struct {
	value int
}

type Add struct {
	left Reducible
	right Reducible
}

type Multiply struct {
	left Reducible
	right Reducible
}


func (n Number) Reduce() Reducible {
	return n
}

func (a Add) Reduce() Reducible {
	numLeft := a.left.Reduce().(Number)
	numRight := a.right.Reduce().(Number)
	return Number{numLeft.value + numRight.value}
}

func (m Multiply) Reduce() Reducible {
	numLeft := m.left.Reduce().(Number)
	numRight := m.right.Reduce().(Number)
	return Number{numLeft.value * numRight.value}
}


func (n Number) ReduceOnce() Reducible {
	return n
}

func (a Add) ReduceOnce() Reducible {
	if a.left.IsReducible() {
		return Add{a.left.Reduce(), a.right}
	} else if a.right.IsReducible() {
		return Add{a.left, a.right.Reduce()}
	} else {
		numLeft := a.left.(Number)
		numRight := a.left.(Number)
		return Number{numLeft.value + numRight.value}
	}
}

func (m Multiply) ReduceOnce() Reducible {
	if m.left.IsReducible() {
		return Multiply{m.left.Reduce(), m.right}
	} else if m.right.IsReducible() {
		return Multiply{m.left, m.right.Reduce()}
	} else {
		numLeft := m.left.(Number)
		numRight := m.left.(Number)
		return Number{numLeft.value * numRight.value}
	}
}




func (n Number) IsReducible() bool {
	return false
}

func (a Add) IsReducible() bool {
	return true
}

func (m Multiply) IsReducible() bool {
	return true
}


func (n Number) String() string {
	return fmt.Sprintf("%d", n.value)
}

func (a Add) String() string {
	return fmt.Sprintf("«%v + %v»", a.left, a.right)
}

func (m Multiply) String() string {
	return fmt.Sprintf("«%v * %v»", m.left, m.right)
}





func main() {
	n1 := Number{2}
	n2 := Number{3}

	fmt.Printf("%v\n", n1)

	a := Add{n1, n2}
	fmt.Printf("%v\n", a)

	m := Multiply{n1, a}
	fmt.Printf("%v\n", m)

	println("------- reduce ---------")
	fmt.Printf("%v\n", n1.Reduce())
	fmt.Printf("%v\n", a.Reduce())
	fmt.Printf("%v\n", m.Reduce())

	println("------- reducible? ---------")
	fmt.Printf("%v\n", n1.IsReducible())
	fmt.Printf("%v\n", a.IsReducible())
	fmt.Printf("%v\n", m.IsReducible())

	println("------- reduceOnce ---------")
	r := n1.ReduceOnce()
	fmt.Printf("%v\n", r)
	r = a.ReduceOnce()
	fmt.Printf("%v\n", r)
	r = m.ReduceOnce()
	fmt.Printf("%v\n", r)

	for r.IsReducible() {
		r = r.ReduceOnce()
		fmt.Printf("%v\n", r)
	}
}

