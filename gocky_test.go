package gocky

import (
	"regexp"
	"testing"
	"fmt"
)

func panda() Grammar {
	determiner := Production{key: "DT", nominals: []string{"the"}}
	noun := Production{key: "N", nominals: []string{"panda", "shoots", "leaves"}}
	verb := Production{key: "V", nominals: []string{"eats", "shoots", "leaves"}}
	conjunction := Production{key: "CC", nominals: []string{"and"}}

	CCN := Production{key: "CCN", left: "CC", right: "N"}  // and leaves
	NP0 := Production{key: "NP0", left: "N", right: "CCN"} // shoots and leaves

	CCV := Production{key: "CCV", left: "CC", right: "V"}  // and leaves
	VP0 := Production{key: "VP0", left: "V", right: "CCV"} // shoots and leaves
	VP1 := Production{key: "VP1", left: "V", right: "VP0"} // eats shoots and leaves
	VP2 := Production{key: "VP2", left: "V", right: "NP0"} // eats shoots and leaves

	DN0 := Production{key: "DN0", left: "DT", right: "N"}   // the panda
	DN1 := Production{key: "DN1", left: "DT", right: "NP0"} // the shoots and leaves

	S0 := Production{key: "S0", left: "DN0", right: "V"}    // the panda eats
	S1 := Production{key: "S1", left: "DN0", right: "VP0"}  // the panda shoots and leaves
	S2 := Production{key: "S2", left: "DN0", right: "VP1"}  // the panda eats shoots and leaves (all verbs)
	S3 := Production{key: "S3", left: "DN0", right: "VP2"}  // the panda eats shoots and leaves (verb eats, shoots and leaves are nouns)

	S4 := Production{key: "S4", left: "DN1", right: "V"}    // the shoots and leaves eat
	S5 := Production{key: "S5", left: "DN1", right: "VP0"}  // the shoots and leaves shoots and leaves
	S6 := Production{key: "S6", left: "DN1", right: "VP1"}  // the shoots and leaves eat shoots and leaves (all verbs)
	S7 := Production{key: "S7", left: "DN1", right: "VP2"}  // the shoots and leaves eat shoots and leaves (verb eats, shoots and leaves are nouns)

	return Grammar{
		determiner, noun, verb, conjunction,
		CCN, NP0,
		CCV, VP0, VP1, VP2,
		DN0, DN1,
		S0, S1, S2, S3,
		S4, S5, S6, S7,
	}
}

func bookFlight() Grammar {
	determiner := Production{key: "DT", nominals: []string{"the", "that", "a"}}
	noun := Production{key: "N", nominals: []string{"book", "flight"}}
	verb := Production{key: "V", nominals: []string{"book"}}
	adj  := Production{key: "JJ", nominals: []string{}}

	NP := Production{key: "NP", left: "DT", right: "N"}
	VP := Production{key: "VP", left: "V", right: "NP"}
	return Grammar{determiner, noun, verb, adj, NP, VP}
}

func bigDog() Grammar {
	return Grammar{
		Production{key: "DT", nominals: []string{"the"}},
		Production{key: "N", nominals: []string{"dog"}},
		Production{key: "J", nominals: []string{"big", "gray", "furry"}},
		Production{key: "NP", left: "DT", right: "N"},
		Production{key: "NP", left: "J", right: "N"},
		Production{key: "NP", left: "J", right: "NP"},
		Production{key: "NP", left: "DT", right: "NP"},
	}
}

// compareNestedStringArray tests whether two nested string arrays are equivalent
//
// The two nested arrays are equivalent if
// the top-level arrays are the same length
//     len(expected) == len(actual)
// the nested arrays at the same index are the same length
//     len(expected[i]) == len(actual[i])
// and the nested values are thes same for the same index pair
//     expected[i][j] == actual[i][j]
func compareNestedStringArray(t *testing.T, name string, expected [][]string, actual [][]string) {
	if len(expected) != len(actual) {
		t.Errorf("%s expected length of %d, but actual length is %d", name, len(expected), len(actual))
	}

	for i := range expected {
		if len(expected[i]) != len(actual[i]) {
			t.Errorf("%s expected length of %d at index %d, but actual length is %d", name, len(expected[i]), i, len(actual[i]))
		}
	}

	for i := range expected {
		for j := range expected[i] {
			if expected[i][j] != actual[i][j] {
				t.Errorf("%s expected \"%s\" at index [%d, %d], but actual value is \"%s\"", name, expected[i][j], i, j, actual[i][j])
			}
		}
	}
}


// TestParses parses a sentence against the provided grammar, then verifies the expected terminal productions are present in the parses.
// We first split the sentence into component words.
// Then we parse the word array against the grammar.
// For each parse, we have a map of productions to terminals that we want to test.
// For each production, we pull out all matching terminal productions out of the phrase and match them.
func TestParses(t *testing.T) {
	type test struct {
		sentence                     string
		grammar                      Grammar
		expectedParseRepresentations []map[string][][]string
	}

	testCases := []test{
		{
			grammar:  bookFlight(),
			sentence: "book that flight",
			expectedParseRepresentations: []map[string][][]string{
				map[string][][]string{
					"V":  [][]string{[]string{"book"}},
					"N":  [][]string{[]string{"flight"}},
					"NP": [][]string{[]string{"that", "flight"}},
					"VP": [][]string{[]string{"book", "that", "flight"}},
					"JJ": [][]string{},
				},
			},
		},
		{
			grammar:  panda(),
			sentence: "the panda eats shoots and leaves",
			expectedParseRepresentations: []map[string][][]string{
				map[string][][]string{
					"V": [][]string{
						[]string{"eats"},
					},
					"N": [][]string{
						[]string{"panda"},
						[]string{"shoots"},
						[]string{"leaves"},
					},
					"DN0": [][]string{
						[]string{"the", "panda"},
					},
					"NP0": [][]string{
						[]string{"shoots", "and", "leaves"},
					},
					"VP2": [][]string{
						[]string{"eats", "shoots", "and", "leaves"},
					},
					"S3": [][]string{
						[]string{"the", "panda", "eats", "shoots", "and", "leaves"},
					},
				},
				map[string][][]string{
					"V": [][]string{
						[]string{"eats"},
						[]string{"shoots"},
						[]string{"leaves"},
					},
					"N": [][]string{
						[]string{"panda"},
					},
					"VP1": [][]string{
						[]string{"eats", "shoots", "and", "leaves"},
					},
					"S2": [][]string{
						[]string{"the", "panda", "eats", "shoots", "and", "leaves"},
					},
				},
			},
		},
		{
			grammar: bigDog(),
			sentence: "the big gray furry dog",
			expectedParseRepresentations: []map[string][][]string{
				map[string][][]string{
					"NP": [][]string{
						[]string{"furry", "dog"},
						[]string{"gray", "furry", "dog"},
						[]string{"big", "gray", "furry", "dog"},
						[]string{"the", "big", "gray", "furry", "dog"},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		words := regexp.MustCompile("\\s+").Split(testCase.sentence, -1)
		actualParses := Parses(words, testCase.grammar)
		if len(actualParses) != len(testCase.expectedParseRepresentations) {
			t.Fatalf("(Test \"%s\"), num parses expected %d, got %d", testCase.sentence, len(testCase.expectedParseRepresentations), len(actualParses))
		}

		for parseIndex, parseRepresentations := range testCase.expectedParseRepresentations {
			actualParse := actualParses[parseIndex]

			for targetProductionKey, expectedProductionTerminalSequences := range parseRepresentations {
				actualProductionTerminalSequences := actualParse.ProductionTerminals(targetProductionKey)
				testName := fmt.Sprintf("(Test \"%s\") production key \"%s\"", testCase.sentence, targetProductionKey)
				compareNestedStringArray(t, testName, expectedProductionTerminalSequences, actualProductionTerminalSequences)
			}
		}
	}
}