package gocky

import "testing"

func TestTerminalLookup(t *testing.T) {
	type test struct {
		name                string
		grammar             Grammar
		lookup              string
		expectedProductions []Production
	}

	target1 := TerminalProduction("target1", []string{"test1", "test2"})
	target2 := TerminalProduction("target2", []string{"test2", "test3"})
	redHerring1 := TerminalProduction("red_herring1", []string{"redHerring"})
	redHerring2 := TerminalProduction("red_herring2", []string{"redHerring"})
	nonterminal1 := NonterminalProduction("nonterminal1", "red_herring1", "red_herring2")
	nonterminal2 := NonterminalProduction("nonterminal2", "red_herring1", "target")
	sentenceProduction := NonterminalProduction("sentence", "nonterminal1", "nonterminal2")
	grammar := Grammar{redHerring1, redHerring2, target1, target2, nonterminal1, nonterminal2, sentenceProduction}

	tests := []test{
		{name: "missing", grammar: grammar, lookup: "missing", expectedProductions: []Production{}},
		{name: "single", grammar: grammar, lookup: "test1", expectedProductions: []Production{target1}},
		{name: "double", grammar: grammar, lookup: "test2", expectedProductions: []Production{target1, target2}},
	}

	for _, testCase := range tests {
		actualMatches := terminalLookup(testCase.lookup, testCase.grammar)

		if len(actualMatches) != len(testCase.expectedProductions) {
			t.Fatalf("(Test \"%s\"), num matches expected %d, got %d", testCase.name, len(testCase.expectedProductions), len(actualMatches))
		}

		for matchIndex, expectedProduction := range testCase.expectedProductions {
			actualMatch := actualMatches[matchIndex]

			if expectedProduction.key != actualMatch.production.key {
				t.Errorf("(Test \"%s\"), match key expected %s, got %s", testCase.name, expectedProduction.key, actualMatch.production.key)
			}

			if testCase.lookup != actualMatch.terminal {
				t.Errorf("(Test \"%s\"), match terminal expected %s, got %s", testCase.name, testCase.lookup, actualMatch.terminal)
			}
		}
	}
}

func TestNonterminalLookup(t *testing.T) {
	type test struct {
		name                string
		grammar             Grammar
		left                Parse
		right               Parse
		expectedProductions []Production
	}

	leftTerminal := TerminalProduction("leftTerminal", []string{"leftNominal"})
	rightTerminal := TerminalProduction("rightTerminal", []string{"rightNominal"})
	redHerringTerminal := TerminalProduction("redHerringTerminal", []string{"redHerring"})
	targetNonterminal := NonterminalProduction("target", "leftTerminal", "rightTerminal")
	redHerringNonterminalLeft := NonterminalProduction("redHerringNonterminalLeft", "redHerringTerminal", "rightTerminal")
	redHerringNonterminalRight := NonterminalProduction("redHerringNonterminalRight", "leftTerminal", "redHerringTerminal")
	unambiguousGrammar := Grammar{leftTerminal, rightTerminal, redHerringTerminal, targetNonterminal, redHerringNonterminalLeft, redHerringNonterminalRight}

	targetNonterminal1 := NonterminalProduction("target1", "leftTerminal", "rightTerminal")
	targetNonterminal2 := NonterminalProduction("target2", "leftTerminal", "rightTerminal")
	ambiguousGrammar := Grammar{leftTerminal, rightTerminal, redHerringTerminal, targetNonterminal1, targetNonterminal2, redHerringNonterminalLeft, redHerringNonterminalRight}

	missingTerminal := TerminalProduction("missing", []string{"missing"})

	tests := []test{
		{
			name:                "missingLeft",
			grammar:             unambiguousGrammar,
			left:                Parse{production: &missingTerminal, terminal: "missing"},
			right:               Parse{production: &rightTerminal, terminal: "rightNominal"},
			expectedProductions: []Production{},
		},
		{
			name:                "missingRight",
			grammar:             unambiguousGrammar,
			left:                Parse{production: &leftTerminal, terminal: "leftNominal"},
			right:               Parse{production: &missingTerminal, terminal: "missing"},
			expectedProductions: []Production{},
		},
		{
			name:                "single",
			grammar:             unambiguousGrammar,
			left:                Parse{production: &leftTerminal, terminal: "leftNominal"},
			right:               Parse{production: &rightTerminal, terminal: "rightNominal"},
			expectedProductions: []Production{targetNonterminal},
		},
		{
			name:                "double",
			grammar:             ambiguousGrammar,
			left:                Parse{production: &leftTerminal, terminal: "leftNominal"},
			right:               Parse{production: &rightTerminal, terminal: "rightNominal"},
			expectedProductions: []Production{targetNonterminal1, targetNonterminal2},
		},
	}

	for _, testCase := range tests {
		actualMatches := nonterminalLookup(&testCase.left, &testCase.right, testCase.grammar)
		if len(actualMatches) != len(testCase.expectedProductions) {
			t.Fatalf("(Test \"%s\"), num matches expected %d, got %d", testCase.name, len(testCase.expectedProductions), len(actualMatches))
		}

		for matchIndex, expectedProduction := range testCase.expectedProductions {
			actualMatch := actualMatches[matchIndex]

			if expectedProduction.key != actualMatch.production.key {
				t.Errorf("(Test \"%s\"), match key expected %s, got %s", testCase.name, expectedProduction.key, actualMatch.production.key)
			}

			if &testCase.left != actualMatch.left {
				t.Errorf("(Test \"%s\"), left node address does not match provided left node for production %s", testCase.name, actualMatch.production.key)
			}

			if &testCase.right != actualMatch.right {
				t.Errorf("(Test \"%s\"), right node address does not match provided right node for production %s", testCase.name, actualMatch.production.key)
			}
		}
	}
}
