package gocky

import "testing"

func TestProductionTerminals(t *testing.T) {
	type test struct {
		name                      string
		parse                     *Parse
		targetProduction          string
		expectedTerminalSequences [][]string
	}

	verb := &Production{key: "V", nominals: []string{"eats", "shoots", "leaves"}}
	conjunction := &Production{key: "CC", nominals: []string{"and"}}

	CCV := &Production{key: "CCV", left: "CC", right: "V"}   // and leaves
	VP0 := &Production{key: "VP0", left: "V", right: "CCV"}  // shoots and leaves
	VP1 := &Production{key: "VP1", left: "V", right: "VP1b"} // eats shoots and leaves

	eats := &Parse{terminal: "eats", production: verb}
	shoots := &Parse{terminal: "shoots", production: verb}
	and := &Parse{terminal: "and", production: conjunction}
	leaves := &Parse{terminal: "leaves", production: verb}

	andLeaves := &Parse{production: CCV, left: and, right: leaves}
	shootsAndLeaves := &Parse{production: VP0, left: shoots, right: andLeaves}
	eatsShootsAndLeaves := &Parse{production: VP1, left: eats, right: shootsAndLeaves}

	// Note that for these tests, "eats", "shoots", and "leaves" are all being classified as verbs, regardless of panda habits.
	testCases := []test{
		{
			name:                      "missing",
			parse:                     eatsShootsAndLeaves,
			targetProduction:          "NP",
			expectedTerminalSequences: [][]string{},
		},
		{
			name:             "verbs",
			parse:            eatsShootsAndLeaves,
			targetProduction: "V",
			expectedTerminalSequences: [][]string{
				[]string{"eats"},
				[]string{"shoots"},
				[]string{"leaves"},
			},
		},
		{
			name:             "conjunctions",
			parse:            eatsShootsAndLeaves,
			targetProduction: "CC",
			expectedTerminalSequences: [][]string{
				[]string{"and"},
			},
		},
		{
			name:             "conj verb phrase",
			parse:            eatsShootsAndLeaves,
			targetProduction: "CCV",
			expectedTerminalSequences: [][]string{
				[]string{"and", "leaves"},
			},
		},
		{
			name:             "verb phrase",
			parse:            eatsShootsAndLeaves,
			targetProduction: "VP0",
			expectedTerminalSequences: [][]string{
				[]string{"shoots", "and", "leaves"},
			},
		},
		{
			name:             "extended verb phrase",
			parse:            eatsShootsAndLeaves,
			targetProduction: "VP1",
			expectedTerminalSequences: [][]string{
				[]string{"eats", "shoots", "and", "leaves"},
			},
		},
	}

	for _, testCase := range testCases {
		actualTerminalSequences := testCase.parse.ProductionTerminals(testCase.targetProduction)
		if len(actualTerminalSequences) != len(testCase.expectedTerminalSequences) {
			t.Fatalf("(Test \"%s\"), num terminal sequences expected %d, got %d", testCase.name, len(testCase.expectedTerminalSequences), len(actualTerminalSequences))
		}

		for actualTerminalSequenceIndex, actualTerminalSequence := range actualTerminalSequences {
			expectedTerminalSequence := testCase.expectedTerminalSequences[actualTerminalSequenceIndex]
			if len(actualTerminalSequence) != len(expectedTerminalSequence) {
				t.Fatalf("(Test \"%s\"), terminal sequence size expected %d, got %d", testCase.name, len(expectedTerminalSequence), len(actualTerminalSequence))
			}

			for actualTerminalIndex, actualTerminal := range actualTerminalSequence {
				expectedTerminal := expectedTerminalSequence[actualTerminalIndex]
				if actualTerminal != expectedTerminal {
					t.Errorf("(Test \"%s\"), expected terminal %s, got %s", testCase.name, expectedTerminal, actualTerminal)
				}
			}
		}
	}
}
