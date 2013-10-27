package main

import (
	"testing"
	"fmt"
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
