package statemachine

import "testing"

func TestEvaluatorCreation(t *testing.T) {
	expr := "a|b"
	str := "a"

	status := Evaluate(expr, str)

	if !status {
		t.Error("test failed !")
	}
}
