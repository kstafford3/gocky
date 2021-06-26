# Gocky
## Gocky is a [CKY parsing](https://en.wikipedia.org/wiki/CYK_algorithm) library written in Go.

## Concepts

Gocky works with a simple binary tree format to describe grammars.

We call the nodes in a grammar tree "Productions" because the two children "produce" the parent node.

For example, a determiner (`DT`) and a noun (`N`) produce a noun phrase (`NP`). We would describe this relationship as `NP -> DT, N`. As a tree, the `NP` node has a left and a right child `DT` and `N` respectively.

The leaves of this tree are a special case. Instead of left and right children, they can contain some number of "nominals". Nominals are real words for that part of speech. The verb leaf node might contain "walk", "run", "swim", etc. as nominals. We call these leaf nodes "Terminal Productions."

## Representing a Grammar
Gocky encodes grammars in tree where each of the leaf nodes encodes "nominals" representing real words, and each of the branches encodes a potential branch in the parse tree.

The nodes in the tree are linked by production "keys". An individual key can indicate many different productions.

For example, we could describe a simple grammar like so:
```
S -> NP, V
S -> N, V

NP -> DT, N

DT -> the, a
N -> dog
V -> barks
```

This grammar would parse "the dog barks", "a dog barks", or just "dog barks".

There are two methods provided to create these productions.
`NonterminalProduction()` creates a branch in the tree like `NP -> DT, N`.
`TerminalProduction()` creates a leaf node, like `N -> "dog"`.

```go
determiner := TerminalProduction("DT", []string{"the", "a"})
noun := TerminalProduction("N", []string{"dog"})
nounPhrase := NonterminalProduction("NP", "DT", "N")
```
Note that we use the key, not the instance, to link a child node to the production.
We can create as many "N" productions as we want to build out this grammar further.

A `Grammar` is just a collection of `Production`s. We can assemble them as such:
```go
grammar := Grammar{ determiner, noun, nounPhrase }
```
There is no required order of productions in a grammar, though it may affect the order of results when parsing.


## Parsing a Sentence
Parsing an array of nominals against a grammar will give us a list of valid `Parse`s.

```go
parses := Parses([]string{ "the", "dog", "barks" }, grammar)
```

The results will contain any parse that explains all of provided nominals based on the rules provided in your grammar. An empty list indicates that the grammar could not parse the sentence.

To retrieve the individual branches from the parse tree, we can use the `ProductionTerminals` method.
```go
nounPhrase := parse.ProductionTerminals("NP")
// nounPhrase = []string{"the", "dog"}

noun := parse.ProductionTerminals("N")
// noun = []string{"dog"}
```

For a more complex example, see gocky_test.go, where we parse ambiguous sentences.

Copyright 2021 Kyle Stafford
