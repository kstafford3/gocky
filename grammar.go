package gocky

// Grammar holds the productions for a context free grammar in chomsky normal form
type Grammar []Production

// Production describes a production in chomsky normal form
// Each production has a key that names the production
//
// Productions have one of the following:
//
// Nominals:
// A list of literal strings that the production represents
// The Production "verb" might have nominals "reach", "throw", "row", "go", etc.
//
// Left and Right Keys:
// References to component productions
// For example, he Production "noun phrase" might have a left "article" and a right "noun"
type Production struct {
	key      string
	left     string
	right    string
	nominals []string
}

// NonterminalProduction creates a non-terminal production in the chomsky normal form
// These Productions describe a Production as a left component and a right component
// The components are also Productions.
// This allows us to create a tree structure of Nonterminal branches reaching Terminal leaves
func NonterminalProduction(key string, left string, right string) Production {
	return Production{
		key:      key,
		left:     left,
		right:    right,
		nominals: []string{},
	}
}

// TerminalProduction creates a terminal production in the chomsky normal form
// These productions describe a set of sting literals
func TerminalProduction(key string, nominals []string) Production {
	return Production{
		key:      key,
		nominals: nominals,
	}
}

// terminalLookup returns a list of Parses for a given nominal
func terminalLookup(nominal string, grammar Grammar) []Parse {
	matchingParses := []Parse{}
	for productionIndex := range grammar {
		production := &grammar[productionIndex]
		if contains(production.nominals, nominal) {
			node := Parse{production: production, terminal: nominal}
			matchingParses = append(matchingParses, node)
		}
	}
	return matchingParses
}

// nonterminalLookup returns a list of matching productions for a pair of child productions
func nonterminalLookup(left *Parse, right *Parse, grammar Grammar) []Parse {
	matchingParses := []Parse{}
	for productionIndex := range grammar {
		production := &grammar[productionIndex]
		if production.left == left.production.key && production.right == right.production.key {
			node := Parse{production: production, left: left, right: right}
			matchingParses = append(matchingParses, node)
		}
	}
	return matchingParses
}
