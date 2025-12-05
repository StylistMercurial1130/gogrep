package statemachine

import (
	"testing"
)

func TestCreateNfaSimple(t *testing.T) {
	str := "abc"
	postfix := ExprToPostFix(str)

	postfix2Nfa(postfix)
}

func TestExprsDoNotPanic(t *testing.T) {
	patterns := []string{
		"a",
		"ab",
		"a|b",
		"a*",
		"a+",
		"a?",
		"(ab)c",
		"a(b|c)d",
		"(a|b)*c",
	}

	for _, p := range patterns {
		t.Run(p, func(t *testing.T) {
			postfix := ExprToPostFix(p)
			if postfix == "" {
				t.Fatalf("postfix is empty for %s", p)
			}

			nfa := postfix2Nfa(postfix)
			if nfa == nil {
				t.Fatalf("postfix2Nfa returned nil for %s (postfix=%s)", p, postfix)
			}
		})
	}
}
