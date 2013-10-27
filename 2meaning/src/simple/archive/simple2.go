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

type Boolean struct {
	value bool
}

type LessThan struct {
	left Reducible
	right Reducible
}


/* ---[ Implementations of Reducible ]--- */

func (n Number) Reduce() Reducible {
	return n
}

func (b Boolean) Reduce() Reducible {
	return b
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

func (lt LessThan) Reduce() Reducible {
	numLeft := lt.left.Reduce().(Number)
	numRight := lt.right.Reduce().(Number)
	return Boolean{numLeft.value < numRight.value}	
}



func (n Number) ReduceOnce() Reducible {
	return n
}

func (b Boolean) ReduceOnce() Reducible {
	return b
}

func (a Add) ReduceOnce() Reducible {
	if a.left.IsReducible() {
		return Add{a.left.Reduce(), a.right}
	} else if a.right.IsReducible() {
		return Add{a.left, a.right.Reduce()}
	} else {
		return a.Reduce()
	}
}

func (m Multiply) ReduceOnce() Reducible {
	if m.left.IsReducible() {
		return Multiply{m.left.Reduce(), m.right}
	} else if m.right.IsReducible() {
		return Multiply{m.left, m.right.Reduce()}
	} else {
		return m.Reduce()
	}
}

func (lt LessThan) ReduceOnce() Reducible {
	if lt.left.IsReducible() {
		return LessThan{lt.left.Reduce(), lt.right}
	} else if lt.right.IsReducible() {
		return LessThan{lt.left, lt.right.Reduce()}
	} else {
		numLeft := lt.left.Reduce().(Number)
		numRight := lt.right.Reduce().(Number)
		return Boolean{numLeft.value < numRight.value}
	}
}



func (n Number) IsReducible() bool {
	return false
}

func (b Boolean) IsReducible() bool {
	return false
}

func (a Add) IsReducible() bool {
	return true
}

func (m Multiply) IsReducible() bool {
	return true
}

func (lt LessThan) IsReducible() bool {
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

func (b Boolean) String() string {
	return fmt.Sprintf("%v", b.value)
}

func (lt LessThan) String() string {
	return fmt.Sprintf("%v < %v", lt.left, lt.right)
}



/* ---[ Machine ]--- */

type Machine struct {
	expression Reducible
}

func (m *Machine) step() {
	m.expression = m.expression.ReduceOnce()
}

func (m *Machine) run() {
	for m.expression.IsReducible() {
		fmt.Printf("%v\n", m.expression)
		m.step()
	}
	fmt.Printf("%v\n", m.expression)
}


/* ---[ Main ]--- */

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

	println("------- Machine Test 1 ---------")
	machine := Machine{m}
	machine.run()


	println("------ Boolean and LessThan -------")
	bb := Boolean{true}
	fmt.Printf("Boolean{true} = %v\n", bb)

	lt := LessThan{Number{77}, Number{14}}
	fmt.Printf("%v\n", lt)
	fmt.Printf("LessThan ReduceOnce: %v\n", lt.Reduce())
	fmt.Printf("LessThan Reduce: %v\n", lt.Reduce())

	lt2 := LessThan{n1, a}
	fmt.Printf("lt2: %v\n", lt2)
	fmt.Printf("lt2: LessThan ReduceOnce: %v\n", lt2.ReduceOnce())
	fmt.Printf("lt2: LessThan Reduce: %v\n", lt2.Reduce())
	
}

