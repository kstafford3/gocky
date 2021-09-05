package gocky

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
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

	S0 := Production{key: "S0", left: "DN0", right: "V"}   // the panda eats
	S1 := Production{key: "S1", left: "DN0", right: "VP0"} // the panda shoots and leaves
	S2 := Production{key: "S2", left: "DN0", right: "VP1"} // the panda eats shoots and leaves (all verbs)
	S3 := Production{key: "S3", left: "DN0", right: "VP2"} // the panda eats shoots and leaves (verb eats, shoots and leaves are nouns)

	S4 := Production{key: "S4", left: "DN1", right: "V"}   // the shoots and leaves eat
	S5 := Production{key: "S5", left: "DN1", right: "VP0"} // the shoots and leaves shoots and leaves
	S6 := Production{key: "S6", left: "DN1", right: "VP1"} // the shoots and leaves eat shoots and leaves (all verbs)
	S7 := Production{key: "S7", left: "DN1", right: "VP2"} // the shoots and leaves eat shoots and leaves (verb eats, shoots and leaves are nouns)

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
	adj := Production{key: "JJ", nominals: []string{}}

	NP := Production{key: "NP", left: "DT", right: "N"}
	VP := Production{key: "VP", left: "V", right: "NP"}
	return Grammar{determiner, noun, verb, adj, NP, VP}
}

func bigDog() Grammar {
	return Grammar{
		Production{key: "N", nominals: []string{"dog"}},
		Production{key: "DT", nominals: []string{"the"}},
		Production{key: "J", nominals: []string{"big", "gray", "furry"}},
		Production{key: "N", left: "DT", right: "N"},
		Production{key: "N", left: "J", right: "N"},
		Production{key: "N", left: "J", right: "NP"},
		Production{key: "N", left: "DT", right: "NP"},
	}
}

// compareNestedStringArray tests whether two nested string arrays are equivalent
func compareNestedStringArray(t *testing.T, name string, expected [][]string, actual [][]string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s expected %v but got %v", name, expected, actual)
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
		expectedProductionKeys       [][]string
		expectedParseRepresentations []map[string][][]string
	}

	testCases := []test{
		{
			grammar:                bookFlight(),
			sentence:               "book that flight",
			expectedProductionKeys: [][]string{{"VP", "V", "NP", "DT", "N"}},
			expectedParseRepresentations: []map[string][][]string{
				{
					"V":  {{"book"}},
					"N":  {{"flight"}},
					"NP": {{"that", "flight"}},
					"VP": {{"book", "that", "flight"}},
					"JJ": {},
				},
			},
		},
		{
			grammar:  panda(),
			sentence: "the panda eats shoots and leaves",
			expectedProductionKeys: [][]string{
				{"S3", "DN0", "DT", "N", "VP2", "V", "NP0", "N", "CCN", "CC", "N"},
				{"S2", "DN0", "DT", "N", "VP1", "V", "VP0", "V", "CCV", "CC", "V"},
			},
			expectedParseRepresentations: []map[string][][]string{
				{
					"V": [][]string{
						{"eats"},
					},
					"N": [][]string{
						{"panda"},
						{"shoots"},
						{"leaves"},
					},
					"DN0": [][]string{
						{"the", "panda"},
					},
					"NP0": [][]string{
						{"shoots", "and", "leaves"},
					},
					"VP2": [][]string{
						{"eats", "shoots", "and", "leaves"},
					},
					"S3": [][]string{
						{"the", "panda", "eats", "shoots", "and", "leaves"},
					},
				},
				{
					"V": [][]string{
						{"eats"},
						{"shoots"},
						{"leaves"},
					},
					"N": [][]string{
						{"panda"},
					},
					"VP1": [][]string{
						{"eats", "shoots", "and", "leaves"},
					},
					"S2": [][]string{
						{"the", "panda", "eats", "shoots", "and", "leaves"},
					},
				},
			},
		},
		{
			grammar:                bigDog(),
			sentence:               "the big gray furry dog",
			expectedProductionKeys: [][]string{{"N", "DT", "N", "J", "N", "J", "N", "J", "N"}},
			expectedParseRepresentations: []map[string][][]string{
				{
					"N": [][]string{
						{"dog"},
						{"furry", "dog"},
						{"gray", "furry", "dog"},
						{"big", "gray", "furry", "dog"},
						{"the", "big", "gray", "furry", "dog"},
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

			expectedProductionKeys := testCase.expectedProductionKeys[parseIndex]
			actualProductionKeys := actualParse.ProductionKeys()
			if !reflect.DeepEqual(expectedProductionKeys, actualProductionKeys) {
				t.Errorf("(Test \"%s\"), expected production keys %v, but got %v", testCase.sentence, expectedProductionKeys, actualProductionKeys)
			}

			for targetProductionKey, expectedProductionTerminalSequences := range parseRepresentations {
				actualProductionTerminalSequences := actualParse.ProductionTerminals(targetProductionKey)
				testName := fmt.Sprintf("(Test \"%s\") production key \"%s\"", testCase.sentence, targetProductionKey)
				compareNestedStringArray(t, testName, expectedProductionTerminalSequences, actualProductionTerminalSequences)
			}
		}
	}
}
