package main

import (
	"testing"
	"fmt"
	"reflect"
)

func assertEqual(exp, actual interface{}, t *testing.T) {
	if exp != actual {
		t.Errorf("exp != actual: >>%v<< != >>%v<<", exp, actual)
	}
}

func assertTrue(assertion bool, t *testing.T) {
	if !assertion {
		t.Errorf("assertion failed")
	}
}

func TestMerge(t *testing.T) {
	env1 := Env{}
	env2 := Env{}

	env1["aa"] = Number{22}
	env1["bb"] = Number{33}
	env2["cc"] = Number{44}

	env3 := EnvMerge(env1, env2)

	var actual, expected string

	actual = fmt.Sprintf("%v", env3)
	expected = "map[aa:22 bb:33 cc:44]"
	assertEqual(expected, actual, t)

	env4 := Env{"aa": Number{555}}
	env5 := EnvMerge(env3, env4)
	actual = fmt.Sprintf("%v", env5)
	expected = "map[aa:555 bb:33 cc:44]"
	assertEqual(expected, actual, t)

	env6 := EnvMerge(env4, env3)
	actual = fmt.Sprintf("%v", env6)
	expected = "map[aa:22 bb:33 cc:44]"
	assertEqual(expected, actual, t)
}


func TestReduceAssign(t *testing.T) {
	assign := Assign{"x", Add{Variable{"x"}, Number{1}}}
	assertTrue(assign.IsReducible(), t)
	env := Env{"x": Number{2}}

	var stmt Reducible
	stmt, env = assign.Reduce(env)
	stmt, env = stmt.Reduce(env)
	stmt, env = stmt.Reduce(env)

	var expected, actual string
	expected = "map[x:3]"
	actual = fmt.Sprintf("%v", env)
	assertEqual(expected, actual, t)
	assertEqual(fmt.Sprintf("%T", stmt), "main.DoNothing", t)
}

func TestAssignWithMachine(t *testing.T) {
	assign := Assign{"x", Add{Variable{"x"}, Number{1}}}
	env := Env{"x": Number{22}}
	machine := Machine{assign}
	env = machine.Run(env)
	
	var expected, actual string
	expected = "map[x:23]"
	actual = fmt.Sprintf("%v", env)
	assertEqual(expected, actual, t)
}

func TestEvaluateAssign(t *testing.T) {
	assign := Assign{"x", Add{Variable{"x"}, Number{1}}}
	env := Env{"x": Number{2}}

	var stmt Reducible
	stmt, env = assign.Evaluate(env)
	
	var expected, actual string
	expected = "map[x:3]"
	actual = fmt.Sprintf("%v", env)
	assertEqual(expected, actual, t)
	assertEqual(fmt.Sprintf("%T", stmt), "main.DoNothing", t)
}


func TestReduceIf(t *testing.T) {
	x := Variable{"x"}
	consequence := Assign{"y", Number{1}}
	alternative := Assign{"y", Number{2}}
	ifStmt := If{x, consequence, alternative}
	env := Env{"x": Boolean{true}}

	mach := Machine{ifStmt}
	env = mach.Run(env)

	one := Number{1}
	if env["y"] != one { t.Errorf("y was not set to one: %v", env) }
}

func TestEvaluateIf(t *testing.T) {
	x := Variable{"x"}
	consequence := Assign{"y", Number{1}}
	alternative := Assign{"y", Number{2}}
	ifStmt := If{x, consequence, alternative}
	env := Env{"x": Boolean{true}}

	var r Reducible
	r, env = ifStmt.Evaluate(env)
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", env)
}

func TestEvaluateSequence(t *testing.T) {
	asgn1 := Assign{"x", Add{Number{1}, Number{1}}}
	asgn2 := Assign{"y", Add{Variable{"x"}, Number{3}}}
	seq := Sequence{asgn1, asgn2}
	env := Env{}

	var r, doNothing Reducible
	r, env = seq.Evaluate(env)

	rKind := reflect.TypeOf(r).Kind()
	doNothing = DoNothing{}
	
	doNothingKind := reflect.TypeOf(doNothing).Kind()
	if rKind != doNothingKind { t.Errorf("%v", rKind) }

	two := Number{2}
	five := Number{5}
	if env["x"] != two { t.Errorf("x is wrong: %v", env) }
	if env["y"] != five { t.Errorf("y is wrong: %v", env) }
}


func TestReduceSequence(t *testing.T) {
	cond := LessThan{Variable{"x"}, Number{5}}
	body := Assign{"x", Multiply{Variable{"x"}, Number{3}}}
	while := While{cond, body}
	env := Env{"x": Number{1}}
	
	mach := Machine{while}
	env = mach.Run(env)

	nine := Number{9}
	if env["x"] != nine { t.Errorf("x is wrong: %v", env) }
}
