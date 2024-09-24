package aql

import "github.com/acontrolfreak/openehr/aql/gen"

type Node struct {
	Type      string
	ClassName string
	Alias     string
	Number    int
	Predicate *gen.PathPredicateContext
	Children  []*Node
}
