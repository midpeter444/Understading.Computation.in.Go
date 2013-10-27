package main

import (
	"fmt"
)

type Reducible interface {
	Evaluate(env Env) (Reducible, Env)
	Reduce(env Env) (Reducible, Env)
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

type Variable struct {
	name string
}

type DoNothing struct {}

type Assign struct {
	name string
	expression Reducible
}

type Env map[string]Reducible


// merges the key-val pairs in env2
// into env1 and returns a new Env object.
// When the same key is in both env1 and env2
// the new Env object will have the value of env2
func EnvMerge(env1 Env, env2 Env) Env {
	nenv := Env{}
	for k,v := range env1 {
		nenv[k] = v
	}
	for k,v := range env2 {
		nenv[k] = v
	}
	return nenv
}


/* ---[ Implementations of Reducible ]--- */

func (n Number) Evaluate(env Env) (Reducible, Env) {
	return n, env
}

func (b Boolean) Evaluate(env Env) (Reducible, Env) {
	return b, env
}

func (a Add) Evaluate(env Env) (Reducible, Env) {
	var rLeft, rRight Reducible
	rLeft, env = a.left.Evaluate(env)
	rRight, env = a.right.Evaluate(env)
	numLeft := rLeft.(Number)
	numRight := rRight.(Number)
	return Number{numLeft.value + numRight.value}, env
}

func (m Multiply) Evaluate(env Env) (Reducible, Env) {
	var rLeft, rRight Reducible
	rLeft, env = m.left.Evaluate(env)
	rRight, env = m.right.Evaluate(env)
	numLeft := rLeft.(Number)
	numRight := rRight.(Number)
	return Number{numLeft.value * numRight.value}, env
}

func (lt LessThan) Evaluate(env Env) (Reducible, Env) {
	var rLeft, rRight Reducible
	rLeft, env = lt.left.Evaluate(env)
	rRight, env = lt.right.Evaluate(env)
	numLeft := rLeft.(Number)
	numRight := rRight.(Number)
	return Boolean{numLeft.value < numRight.value}, env
}

func (v Variable) Evaluate(env Env) (Reducible, Env) {
	// TODO: shouldn't this check to reduce until points
	// to a terminal?
	return env[v.name], env
}

func (d DoNothing) Evaluate(env Env) (Reducible, Env) {
	return d, env
}

func (a Assign) Evaluate(env Env) (Reducible, Env) {
	var r Reducible
	r = a
	for r.IsReducible() {
		r, env = r.Reduce(env)
	}
	return r, env
}


func (n Number) Reduce(env Env) (Reducible, Env) {
	return n, env
}

func (b Boolean) Reduce(env Env) (Reducible, Env) {
	return b, env
}

func (a Add) Reduce(env Env) (Reducible, Env) {
	if a.left.IsReducible() {
		var newleft Reducible
		newleft, env = a.left.Evaluate(env)
		return Add{newleft, a.right}, env

	} else if a.right.IsReducible() {
		var newright Reducible
		newright, env = a.right.Evaluate(env)
		return Add{a.left, newright}, env

	} else {
		return a.Evaluate(env)
	}
}

func (m Multiply) Reduce(env Env) (Reducible, Env) {
	if m.left.IsReducible() {
		var newleft Reducible
		newleft, env = m.left.Evaluate(env)
		return Multiply{newleft, m.right}, env

	} else if m.right.IsReducible() {
		var newright Reducible
		newright, env = m.right.Evaluate(env)
		return Multiply{m.left, newright}, env

	} else {
		return m.Evaluate(env)
	}
}

func (lt LessThan) Reduce(env Env) (Reducible, Env) {
	if lt.left.IsReducible() {
		var newleft Reducible
		newleft, env = lt.left.Evaluate(env)
		return LessThan{newleft, lt.right}, env

	} else if lt.right.IsReducible() {
		var newright Reducible
		newright, env = lt.right.Evaluate(env)
		return LessThan{lt.left, newright}, env

	} else {
		var rLeft, rRight Reducible
		rLeft, env = lt.left.Evaluate(env)
		rRight, env = lt.right.Evaluate(env)
		numLeft := rLeft.(Number)
		numRight := rRight.(Number)
		return Boolean{numLeft.value < numRight.value}, env
	}
}

func (v Variable) Reduce(env Env) (Reducible, Env) {
	// TODO: shouldn't this reduce once if pointing to a non-terminal?
	return env[v.name], env
}

func (d DoNothing) Reduce(env Env) (Reducible, Env) {
	return d, env
}

func (a Assign) Reduce(env Env) (Reducible, Env) {
	if a.expression.IsReducible() {
		var r Reducible
		r, env = a.expression.Reduce(env)
		return Assign{a.name, r}, env
	} else {
		return DoNothing{}, EnvMerge(env, Env{a.name: a.expression})
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

func (d DoNothing) IsReducible() bool {
	return false
}

func (v Variable) IsReducible() bool {
	// TODO: shouldn't this depend on whether the
	// assignment points to a terminal or not?
	return true
}

func (a Assign) IsReducible() bool {
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

func (v Variable) String() string {
	return fmt.Sprintf("«%v»", v.name)
	// return v.name
}

func (d DoNothing) String() string {
	return "do-nothing"
}

func (a Assign) String() string {
	return fmt.Sprintf("«%v = %v»", a.name, a.expression)
}


/* ---[ Machine ]--- */

type Machine struct {
	expression Reducible
}

func (m *Machine) step(env Env) Env {
	m.expression, env = m.expression.Reduce(env)
	return env
}

func (m *Machine) run(env Env) {
	for m.expression.IsReducible() {
		fmt.Printf("%v\n", m.expression)
		env = m.step(env)
	}
	fmt.Printf("%v\n", m.expression)
}


/* ---[ Main ]--- */

func firstR(r Reducible, env Env) Reducible {
	return r
}

func main() {
	env := Env{}
	n1 := Number{2}
	n2 := Number{3}

	fmt.Printf("%v\n", n1)

	a := Add{n1, n2}
	fmt.Printf("%v\n", a)

	m := Multiply{n1, a}
	fmt.Printf("%v\n", m)

	println("------- reduce ---------")
	var r Reducible
	r, env = n1.Evaluate(env)
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", firstR(a.Evaluate(env)))
	fmt.Printf("%v\n", firstR(m.Evaluate(env)))

	println("------- reducible? ---------")
	fmt.Printf("%v\n", n1.IsReducible())
	fmt.Printf("%v\n", a.IsReducible())
	fmt.Printf("%v\n", m.IsReducible())

	println("------- reduceOnce ---------")
	r, env = n1.Reduce(env)
	fmt.Printf("%v\n", r)
	r, env = a.Reduce(env)
	fmt.Printf("%v\n", r)
	r, env = m.Reduce(env)
	fmt.Printf("%v\n", r)

	for r.IsReducible() {
		r, env = r.Reduce(env)
		fmt.Printf("%v\n", r)
	}

	println("------- Machine Test 1 ---------")
	machine := Machine{m}
	machine.run(env)


	println("------ Boolean and LessThan -------")
	bb := Boolean{true}
	fmt.Printf("Boolean{true} = %v\n", bb)

	lt := LessThan{Number{77}, Number{14}}
	fmt.Printf("%v\n", lt)
	fmt.Printf("LessThan Reduce: %v\n", firstR(lt.Evaluate(env)))
	fmt.Printf("LessThan Evaluate: %v\n", firstR(lt.Evaluate(env)))

	lt2 := LessThan{n1, a}
	fmt.Printf("lt2: %v\n", lt2)
	fmt.Printf("lt2: LessThan Reduce: %v\n", firstR(lt2.Reduce(env)))
	fmt.Printf("lt2: LessThan Evaluate: %v\n", firstR(lt2.Evaluate(env)))

	println("------ Variable -------")
	x := Variable{"x"}
	y := Variable{"y"}
	machine = Machine{ Add{x,y} } 
	machine.run(Env{"x": Number{3}, "y": Number{4}})

	println("------ Assign -------")
	assign := Assign{"x", Add{Variable{"x"}, Number{1}}}
	fmt.Printf("Assign is reducible?: %v\n", assign.IsReducible())
	env = Env{"x": Number{2}}

	var stmt Reducible
	stmt, env = assign.Reduce(env)
	fmt.Printf("stmt: %T :: %v\n", stmt, stmt)
	fmt.Printf("env: %v\n", env)
	stmt, env = stmt.Reduce(env)
	fmt.Printf("stmt: %T :: %v\n", stmt, stmt)
	fmt.Printf("env: %v\n", env)
	stmt, env = stmt.Reduce(env)
	fmt.Printf("stmt: %T :: %v\n", stmt, stmt)
	fmt.Printf("env: %v\n", env)
}
