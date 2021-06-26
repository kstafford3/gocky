package gocky

// MatchingParses produces a list of parses based on a list of words, a grammar, and a target production key.
// Only parses that can be generated from the target production keys will be returned.
func MatchingParses(words []string, grammar Grammar, targetProductionKeys []string) []Parse {
	parses := Parses(words, grammar)
	matchingParses := []Parse{}
	for _, parse := range parses {
		if contains(targetProductionKeys, parse.production.key) {
			matchingParses = append(matchingParses, parse)
		}
	}
	return matchingParses
}

// Parses produces a list of parses based on a list of words and a grammar.
// Each parse will describe a different parse tree for the words based on the grammar.
func Parses(words []string, grammar Grammar) []Parse {
	return ckyParse(words, grammar)
}

// ckyParse performs a parse based on the CKY algorithm.
// https://en.wikipedia.org/wiki/CYK_algorithm
func ckyParse(words []string, grammar Grammar) []Parse {
	table := make([][][]Parse, len(words)+1)
	for endIndex := 1; endIndex <= len(words); endIndex++ {
		table[endIndex-1] = make([][]Parse, len(words)+1)
		table[endIndex-1][endIndex] = terminalLookup(words[endIndex-1], grammar)
		for startIndex := endIndex - 2; startIndex >= 0; startIndex-- {
			table[startIndex][endIndex] = []Parse{}
			for splitIndex := startIndex + 1; splitIndex < endIndex; splitIndex++ {
				splitProductions := getGeneratingProductions(table[startIndex][splitIndex], table[splitIndex][endIndex], grammar)
				table[startIndex][endIndex] = append(table[startIndex][endIndex], splitProductions...)
			}
		}
	}
	return table[0][len(words)]
}

// getGeneratingProductions takes two parses and generates a list of parses that that could explain them as left and right components of a Production
//
func getGeneratingProductions(leftParses []Parse, rightParses []Parse, grammar Grammar) []Parse {
	allProductions := []Parse{}
	for leftParseIndex := range leftParses {
		leftParse := &leftParses[leftParseIndex]
		for rightParseIndex := range rightParses {
			rightParse := &rightParses[rightParseIndex]
			localProductions := nonterminalLookup(leftParse, rightParse, grammar)
			allProductions = append(allProductions, localProductions...)
		}
	}
	return allProductions
}
