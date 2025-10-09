package aql

import (
	"errors"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/freekieb7/gopenehr/internal/aql/gen"
)

type TreeShapeListener struct {
	*gen.BaseAQLListener
	Query *gen.QueryContext
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (t *TreeShapeListener) EnterQuery(ctx *gen.QueryContext) {
	t.Query = ctx
}

type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []error
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		Errors: make([]error, 0),
	}
}

func (e *ErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, charPositionInLine int, msg string, _ antlr.RecognitionException) {
	e.Errors = append(e.Errors, fmt.Errorf("error at %d:%d %s", line, charPositionInLine, msg))
}

func QueryContext(aql string) (gen.IQueryContext, error) {
	listener := NewTreeShapeListener()
	errorListener := NewErrorListener()

	input := antlr.NewInputStream(aql)
	lexer := gen.NewAQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	p := gen.NewAQLParser(stream)
	p.AddErrorListener(errorListener)

	antlr.ParseTreeWalkerDefault.Walk(listener, p.Query())

	if len(errorListener.Errors) > 0 {
		return nil, errors.Join(errorListener.Errors...)
	}

	return listener.Query, nil
}
