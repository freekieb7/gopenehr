package aql

import (
	"github.com/acontrolfreak/openehr/aql/gen"
	"log"
)

type ParserListener struct {
	*gen.BaseAqlParserListener
	nodes            []*Node
	sources          map[string]*Node
	selects          []*gen.SelectExprContext
	where            *gen.WhereExprContext
	sourceNodeNumber int
}

func NewAqlParserListener() *ParserListener {
	return &ParserListener{
		sources: make(map[string]*Node),
	}
}

func (s *ParserListener) Selects() []*gen.SelectExprContext {
	return s.selects
}

func (s *ParserListener) EnterSelectExpr(ctx *gen.SelectExprContext) {
	s.selects = append(s.selects, ctx)
}

func (s *ParserListener) EnterWhereExpr(ctx *gen.WhereExprContext) {
	s.where = ctx
}

func (s *ParserListener) ExitContainsExpr(ctx *gen.ContainsExprContext) {
	if ctx.CONTAINS() != nil {
		if ctx.NOT() != nil {
			// TODO not supported
			return
		}

		l := len(s.nodes)

		right := s.nodes[l-1]
		left := s.nodes[l-2]

		s.nodes = s.nodes[:l-2]

		var node Node
		node.Type = "CONTAINS"
		node.Children = append(node.Children, left, right)

		s.nodes = append(s.nodes, &node)
		return
	}

	if ctx.AND() != nil {
		return
	}

	if ctx.OR() != nil {
		return
	}
}

func (s *ParserListener) EnterClassExpression(ctx *gen.ClassExpressionContext) {
	var node Node

	node.Type = "CLASS"

	node.Number = s.sourceNodeNumber
	s.sourceNodeNumber++

	if i := ctx.IDENTIFIER(0); i != nil {
		node.ClassName = i.GetText()
	} else {
		log.Fatal("no class identifier available for node")
	}

	if i := ctx.IDENTIFIER(1); i != nil {
		node.Alias = i.GetText()

		if _, exists := s.sources[node.ClassName]; !exists {
			// TODO error dup alias
		}

		s.sources[node.Alias] = &node
	}

	if ctx.PathPredicate() != nil {
		predicate, ok := ctx.PathPredicate().(*gen.PathPredicateContext)

		if !ok {
			log.Fatal("could not assert path predicate")
		}

		node.Predicate = predicate
	}

	s.nodes = append(s.nodes, &node)
}
