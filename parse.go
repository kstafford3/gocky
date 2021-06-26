package gocky

// Parse captures the generated productions or terminal from a generating Production
// A parsed node can be traced through each production back to all generated terminals
//
// A Production describes the structure of a grammar.
// A Parse describes an actual generation from a grammar.
type Parse struct {
	production *Production
	left       *Parse
	right      *Parse
	terminal   string
}

// ProductionTerminals returns each substring representing the provided production key.
// For example, given a production key of "VP", ProductionTerminals will return all terminal combinations that represent a "VP" in this parse.
func (p *Parse) ProductionTerminals(productionKey string) [][]string {
	productionNodes := traverseToKey(p, productionKey)
	terminalRepresentations := [][]string{}
	for _, productionNode := range productionNodes {
		terminalRepresentations = append(terminalRepresentations, nodeTerminals(productionNode)...)
	}
	return terminalRepresentations
}

// traverseToKey traverses the Parse tree to find component Parses that match the given production key
func traverseToKey(node *Parse, productionKey string) []*Parse {
	if node == nil {
		return []*Parse{}
	}
	leftMatches := traverseToKey(node.left, productionKey)
	rightMatches := traverseToKey(node.right, productionKey)
	matches := append(leftMatches, rightMatches...)
	if node.production.key == productionKey {
		matches = append(matches, node)
	}
	return matches
}

// nodeTerminals collects the terminals for a given Parse
func nodeTerminals(node *Parse) [][]string {
	if node == nil {
		return [][]string{}
	}
	if len(node.terminal) > 0 {
		return [][]string{
			[]string{node.terminal},
		}
	}
	leftTerminals := nodeTerminals(node.left)
	rightTerminals := nodeTerminals(node.right)
	terminalCombinations := [][]string{}
	for _, leftTerminal := range leftTerminals {
		for _, rightTerminal := range rightTerminals {
			combination := append(leftTerminal, rightTerminal...)
			terminalCombinations = append(terminalCombinations, combination)
		}
	}
	return terminalCombinations
}
