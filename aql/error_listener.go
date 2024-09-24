package aql

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"strconv"
)

type ErrorListener struct {
	Errors []string
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		Errors: make([]string, 0),
	}
}

func (d *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	d.Errors = append(d.Errors, fmt.Sprintf("line "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg))
}

func (d *ErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (d *ErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (d *ErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

func (d *ErrorListener) Count() int {
	return len(d.Errors)
}
