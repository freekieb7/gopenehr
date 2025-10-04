// Code generated from AQL.g4 by ANTLR 4.13.2. DO NOT EDIT.

package gen // AQL
import (
	"fmt"
	"strconv"
  	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}


type AQLParser struct {
	*antlr.BaseParser
}

var AQLParserStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  LiteralNames           []string
  SymbolicNames          []string
  RuleNames              []string
  PredictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func aqlParserInit() {
  staticData := &AQLParserStaticData
  staticData.LiteralNames = []string{
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "':'", "';'", "'<'", "'>'", "'<='", 
    "'>='", "'!='", "'='", "'('", "')'", "','", "'/'", "'*'", "'+'", "'-'", 
    "'['", "']'", "'{'", "'}'", "'--'",
  }
  staticData.SymbolicNames = []string{
    "", "WS", "UNICODE_BOM", "SELECT", "AS", "FROM", "WHERE", "ORDER_BY", 
    "GROUP_BY", "DESC", "ASC", "LIMIT", "OFFSET", "DISTINCT", "NULL", "TRUE", 
    "FALSE", "BOOLEAN", "ALL_VERSIONS", "CONTAINS", "AND", "OR", "NOT", 
    "EXISTS", "COMPARISON_OPERATOR", "LIKE", "IN", "ON", "AT", "JOIN", "LEFT", 
    "UNION", "ALL", "LENGTH", "POSITION", "SUBSTRING", "CONCAT", "CONCAT_WS", 
    "ABS", "MOD", "CEIL", "FLOOR", "ROUND", "CURRENT_DATE", "CURRENT_TIME", 
    "CURRENT_DATE_TIME", "CURRENT_TIMEZONE", "COUNT", "MIN", "MAX", "SUM", 
    "AVG", "PARAMETER", "CAST", "IDENTIFIER", "STRING", "INTEGER", "FLOAT", 
    "SCI_INTEGER", "SCI_FLOAT", "SYM_COLON", "SYM_SEMICOLON", "SYM_LT", 
    "SYM_GT", "SYM_LE", "SYM_GE", "SYM_NE", "SYM_EQ", "SYM_LEFT_PAREN", 
    "SYM_RIGHT_PAREN", "SYM_COMMA", "SYM_SLASH", "SYM_ASTERISK", "SYM_PLUS", 
    "SYM_MINUS", "SYM_LEFT_BRACKET", "SYM_RIGHT_BRACKET", "SYM_LEFT_CURLY", 
    "SYM_RIGHT_CURLY", "SYM_DOUBLE_DASH",
  }
  staticData.RuleNames = []string{
    "query", "selectQuery", "selectClause", "fromClause", "joinClause", 
    "whereClause", "groupByClause", "orderByClause", "limitOffsetClause", 
    "selectExpr", "fromExpr", "joinExpr", "whereExpr", "orderByExpr", "columnExpr", 
    "identifiedPath", "objectPath", "pathPart", "pathCondition", "booleanCondition", 
    "pathConditionOperand", "comparisonOperand", "inOperand", "inOperandValue", 
    "primitive", "intPrimitive", "floatPrimitive", "stringOperand", "intOperand", 
    "numbericOperand", "functionCall", "aggregateFunctionCall", "limitOperand",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 79, 475, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 
	4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 
	10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15, 
	2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2, 
	21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26, 
	7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7, 
	31, 2, 32, 7, 32, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 5, 1, 73, 8, 1, 10, 
	1, 12, 1, 76, 9, 1, 1, 1, 3, 1, 79, 8, 1, 1, 1, 3, 1, 82, 8, 1, 1, 1, 1, 
	1, 3, 1, 86, 8, 1, 1, 1, 1, 1, 3, 1, 90, 8, 1, 1, 1, 3, 1, 93, 8, 1, 3, 
	1, 95, 8, 1, 1, 2, 1, 2, 3, 2, 99, 8, 2, 1, 2, 1, 2, 1, 2, 5, 2, 104, 8, 
	2, 10, 2, 12, 2, 107, 9, 2, 1, 3, 1, 3, 1, 3, 1, 4, 3, 4, 113, 8, 4, 1, 
	4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 5, 6, 125, 8, 
	6, 10, 6, 12, 6, 128, 9, 6, 1, 7, 1, 7, 1, 7, 1, 7, 5, 7, 134, 8, 7, 10, 
	7, 12, 7, 137, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 143, 8, 8, 1, 8, 1, 
	8, 1, 8, 1, 8, 3, 8, 149, 8, 8, 3, 8, 151, 8, 8, 1, 9, 1, 9, 1, 9, 3, 9, 
	156, 8, 9, 1, 9, 3, 9, 159, 8, 9, 3, 9, 161, 8, 9, 1, 10, 1, 10, 1, 10, 
	1, 10, 3, 10, 167, 8, 10, 1, 10, 3, 10, 170, 8, 10, 1, 11, 1, 11, 1, 11, 
	1, 11, 3, 11, 176, 8, 11, 1, 11, 3, 11, 179, 8, 11, 1, 11, 1, 11, 1, 11, 
	1, 11, 3, 11, 185, 8, 11, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 191, 8, 11, 
	1, 11, 1, 11, 3, 11, 195, 8, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 
	12, 1, 12, 1, 12, 3, 12, 205, 8, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 
	1, 12, 5, 12, 213, 8, 12, 10, 12, 12, 12, 216, 9, 12, 1, 13, 1, 13, 3, 
	13, 220, 8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 3, 14, 227, 8, 14, 1, 
	15, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 234, 8, 15, 1, 15, 1, 15, 3, 15, 
	238, 8, 15, 1, 15, 3, 15, 241, 8, 15, 1, 16, 1, 16, 1, 16, 5, 16, 246, 
	8, 16, 10, 16, 12, 16, 249, 9, 16, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 3, 
	17, 256, 8, 17, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 
	1, 18, 3, 18, 267, 8, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 5, 
	18, 275, 8, 18, 10, 18, 12, 18, 278, 9, 18, 1, 19, 1, 19, 1, 19, 1, 19, 
	1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 
	19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 302, 8, 19, 
	1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 308, 8, 19, 1, 20, 1, 20, 1, 20, 3, 
	20, 313, 8, 20, 1, 21, 1, 21, 1, 21, 3, 21, 318, 8, 21, 1, 22, 1, 22, 1, 
	22, 1, 22, 5, 22, 324, 8, 22, 10, 22, 12, 22, 327, 9, 22, 3, 22, 329, 8, 
	22, 1, 23, 1, 23, 3, 23, 333, 8, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 
	3, 24, 340, 8, 24, 1, 25, 3, 25, 343, 8, 25, 1, 25, 1, 25, 3, 25, 347, 
	8, 25, 1, 25, 3, 25, 350, 8, 25, 1, 26, 3, 26, 353, 8, 26, 1, 26, 1, 26, 
	3, 26, 357, 8, 26, 1, 26, 3, 26, 360, 8, 26, 1, 27, 1, 27, 1, 28, 1, 28, 
	1, 28, 3, 28, 367, 8, 28, 1, 29, 1, 29, 1, 29, 1, 29, 3, 29, 373, 8, 29, 
	1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 
	30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 
	1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 5, 30, 401, 8, 30, 10, 30, 12, 30, 404, 
	9, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 5, 
	30, 415, 8, 30, 10, 30, 12, 30, 418, 9, 30, 1, 30, 1, 30, 1, 30, 1, 30, 
	1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 
	30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 
	1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 3, 30, 454, 
	8, 30, 1, 31, 1, 31, 1, 31, 3, 31, 459, 8, 31, 1, 31, 1, 31, 3, 31, 463, 
	8, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 3, 31, 471, 8, 31, 1, 
	32, 1, 32, 1, 32, 0, 2, 24, 36, 33, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 
	20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 
	56, 58, 60, 62, 64, 0, 4, 1, 0, 9, 10, 2, 0, 52, 52, 54, 55, 1, 0, 48, 
	51, 2, 0, 52, 52, 56, 56, 529, 0, 66, 1, 0, 0, 0, 2, 69, 1, 0, 0, 0, 4, 
	96, 1, 0, 0, 0, 6, 108, 1, 0, 0, 0, 8, 112, 1, 0, 0, 0, 10, 117, 1, 0, 
	0, 0, 12, 120, 1, 0, 0, 0, 14, 129, 1, 0, 0, 0, 16, 150, 1, 0, 0, 0, 18, 
	160, 1, 0, 0, 0, 20, 162, 1, 0, 0, 0, 22, 194, 1, 0, 0, 0, 24, 204, 1, 
	0, 0, 0, 26, 217, 1, 0, 0, 0, 28, 226, 1, 0, 0, 0, 30, 228, 1, 0, 0, 0, 
	32, 242, 1, 0, 0, 0, 34, 250, 1, 0, 0, 0, 36, 266, 1, 0, 0, 0, 38, 307, 
	1, 0, 0, 0, 40, 312, 1, 0, 0, 0, 42, 317, 1, 0, 0, 0, 44, 328, 1, 0, 0, 
	0, 46, 332, 1, 0, 0, 0, 48, 339, 1, 0, 0, 0, 50, 349, 1, 0, 0, 0, 52, 359, 
	1, 0, 0, 0, 54, 361, 1, 0, 0, 0, 56, 366, 1, 0, 0, 0, 58, 372, 1, 0, 0, 
	0, 60, 453, 1, 0, 0, 0, 62, 470, 1, 0, 0, 0, 64, 472, 1, 0, 0, 0, 66, 67, 
	3, 2, 1, 0, 67, 68, 5, 0, 0, 1, 68, 1, 1, 0, 0, 0, 69, 70, 3, 4, 2, 0, 
	70, 74, 3, 6, 3, 0, 71, 73, 3, 8, 4, 0, 72, 71, 1, 0, 0, 0, 73, 76, 1, 
	0, 0, 0, 74, 72, 1, 0, 0, 0, 74, 75, 1, 0, 0, 0, 75, 78, 1, 0, 0, 0, 76, 
	74, 1, 0, 0, 0, 77, 79, 3, 10, 5, 0, 78, 77, 1, 0, 0, 0, 78, 79, 1, 0, 
	0, 0, 79, 81, 1, 0, 0, 0, 80, 82, 3, 12, 6, 0, 81, 80, 1, 0, 0, 0, 81, 
	82, 1, 0, 0, 0, 82, 94, 1, 0, 0, 0, 83, 85, 5, 31, 0, 0, 84, 86, 5, 32, 
	0, 0, 85, 84, 1, 0, 0, 0, 85, 86, 1, 0, 0, 0, 86, 87, 1, 0, 0, 0, 87, 95, 
	3, 2, 1, 0, 88, 90, 3, 14, 7, 0, 89, 88, 1, 0, 0, 0, 89, 90, 1, 0, 0, 0, 
	90, 92, 1, 0, 0, 0, 91, 93, 3, 16, 8, 0, 92, 91, 1, 0, 0, 0, 92, 93, 1, 
	0, 0, 0, 93, 95, 1, 0, 0, 0, 94, 83, 1, 0, 0, 0, 94, 89, 1, 0, 0, 0, 95, 
	3, 1, 0, 0, 0, 96, 98, 5, 3, 0, 0, 97, 99, 5, 13, 0, 0, 98, 97, 1, 0, 0, 
	0, 98, 99, 1, 0, 0, 0, 99, 100, 1, 0, 0, 0, 100, 105, 3, 18, 9, 0, 101, 
	102, 5, 70, 0, 0, 102, 104, 3, 18, 9, 0, 103, 101, 1, 0, 0, 0, 104, 107, 
	1, 0, 0, 0, 105, 103, 1, 0, 0, 0, 105, 106, 1, 0, 0, 0, 106, 5, 1, 0, 0, 
	0, 107, 105, 1, 0, 0, 0, 108, 109, 5, 5, 0, 0, 109, 110, 3, 20, 10, 0, 
	110, 7, 1, 0, 0, 0, 111, 113, 5, 30, 0, 0, 112, 111, 1, 0, 0, 0, 112, 113, 
	1, 0, 0, 0, 113, 114, 1, 0, 0, 0, 114, 115, 5, 29, 0, 0, 115, 116, 3, 22, 
	11, 0, 116, 9, 1, 0, 0, 0, 117, 118, 5, 6, 0, 0, 118, 119, 3, 24, 12, 0, 
	119, 11, 1, 0, 0, 0, 120, 121, 5, 8, 0, 0, 121, 126, 3, 30, 15, 0, 122, 
	123, 5, 70, 0, 0, 123, 125, 3, 30, 15, 0, 124, 122, 1, 0, 0, 0, 125, 128, 
	1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 126, 127, 1, 0, 0, 0, 127, 13, 1, 0, 
	0, 0, 128, 126, 1, 0, 0, 0, 129, 130, 5, 7, 0, 0, 130, 135, 3, 26, 13, 
	0, 131, 132, 5, 70, 0, 0, 132, 134, 3, 26, 13, 0, 133, 131, 1, 0, 0, 0, 
	134, 137, 1, 0, 0, 0, 135, 133, 1, 0, 0, 0, 135, 136, 1, 0, 0, 0, 136, 
	15, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 138, 139, 5, 11, 0, 0, 139, 142, 
	3, 64, 32, 0, 140, 141, 5, 12, 0, 0, 141, 143, 3, 64, 32, 0, 142, 140, 
	1, 0, 0, 0, 142, 143, 1, 0, 0, 0, 143, 151, 1, 0, 0, 0, 144, 145, 5, 12, 
	0, 0, 145, 148, 3, 64, 32, 0, 146, 147, 5, 11, 0, 0, 147, 149, 3, 64, 32, 
	0, 148, 146, 1, 0, 0, 0, 148, 149, 1, 0, 0, 0, 149, 151, 1, 0, 0, 0, 150, 
	138, 1, 0, 0, 0, 150, 144, 1, 0, 0, 0, 151, 17, 1, 0, 0, 0, 152, 161, 5, 
	72, 0, 0, 153, 158, 3, 28, 14, 0, 154, 156, 5, 4, 0, 0, 155, 154, 1, 0, 
	0, 0, 155, 156, 1, 0, 0, 0, 156, 157, 1, 0, 0, 0, 157, 159, 5, 54, 0, 0, 
	158, 155, 1, 0, 0, 0, 158, 159, 1, 0, 0, 0, 159, 161, 1, 0, 0, 0, 160, 
	152, 1, 0, 0, 0, 160, 153, 1, 0, 0, 0, 161, 19, 1, 0, 0, 0, 162, 166, 5, 
	54, 0, 0, 163, 164, 5, 75, 0, 0, 164, 165, 5, 18, 0, 0, 165, 167, 5, 76, 
	0, 0, 166, 163, 1, 0, 0, 0, 166, 167, 1, 0, 0, 0, 167, 169, 1, 0, 0, 0, 
	168, 170, 5, 54, 0, 0, 169, 168, 1, 0, 0, 0, 169, 170, 1, 0, 0, 0, 170, 
	21, 1, 0, 0, 0, 171, 175, 5, 54, 0, 0, 172, 173, 5, 75, 0, 0, 173, 174, 
	5, 18, 0, 0, 174, 176, 5, 76, 0, 0, 175, 172, 1, 0, 0, 0, 175, 176, 1, 
	0, 0, 0, 176, 178, 1, 0, 0, 0, 177, 179, 5, 54, 0, 0, 178, 177, 1, 0, 0, 
	0, 178, 179, 1, 0, 0, 0, 179, 180, 1, 0, 0, 0, 180, 181, 5, 27, 0, 0, 181, 
	195, 5, 54, 0, 0, 182, 184, 5, 54, 0, 0, 183, 185, 5, 54, 0, 0, 184, 183, 
	1, 0, 0, 0, 184, 185, 1, 0, 0, 0, 185, 186, 1, 0, 0, 0, 186, 187, 5, 26, 
	0, 0, 187, 195, 5, 54, 0, 0, 188, 190, 5, 54, 0, 0, 189, 191, 5, 54, 0, 
	0, 190, 189, 1, 0, 0, 0, 190, 191, 1, 0, 0, 0, 191, 192, 1, 0, 0, 0, 192, 
	193, 5, 28, 0, 0, 193, 195, 3, 30, 15, 0, 194, 171, 1, 0, 0, 0, 194, 182, 
	1, 0, 0, 0, 194, 188, 1, 0, 0, 0, 195, 23, 1, 0, 0, 0, 196, 197, 6, 12, 
	-1, 0, 197, 205, 3, 38, 19, 0, 198, 199, 5, 22, 0, 0, 199, 205, 3, 24, 
	12, 4, 200, 201, 5, 68, 0, 0, 201, 202, 3, 24, 12, 0, 202, 203, 5, 69, 
	0, 0, 203, 205, 1, 0, 0, 0, 204, 196, 1, 0, 0, 0, 204, 198, 1, 0, 0, 0, 
	204, 200, 1, 0, 0, 0, 205, 214, 1, 0, 0, 0, 206, 207, 10, 3, 0, 0, 207, 
	208, 5, 20, 0, 0, 208, 213, 3, 24, 12, 4, 209, 210, 10, 2, 0, 0, 210, 211, 
	5, 21, 0, 0, 211, 213, 3, 24, 12, 3, 212, 206, 1, 0, 0, 0, 212, 209, 1, 
	0, 0, 0, 213, 216, 1, 0, 0, 0, 214, 212, 1, 0, 0, 0, 214, 215, 1, 0, 0, 
	0, 215, 25, 1, 0, 0, 0, 216, 214, 1, 0, 0, 0, 217, 219, 5, 54, 0, 0, 218, 
	220, 7, 0, 0, 0, 219, 218, 1, 0, 0, 0, 219, 220, 1, 0, 0, 0, 220, 27, 1, 
	0, 0, 0, 221, 227, 3, 48, 24, 0, 222, 227, 5, 52, 0, 0, 223, 227, 3, 30, 
	15, 0, 224, 227, 3, 62, 31, 0, 225, 227, 3, 60, 30, 0, 226, 221, 1, 0, 
	0, 0, 226, 222, 1, 0, 0, 0, 226, 223, 1, 0, 0, 0, 226, 224, 1, 0, 0, 0, 
	226, 225, 1, 0, 0, 0, 227, 29, 1, 0, 0, 0, 228, 233, 5, 54, 0, 0, 229, 
	230, 5, 75, 0, 0, 230, 231, 3, 36, 18, 0, 231, 232, 5, 76, 0, 0, 232, 234, 
	1, 0, 0, 0, 233, 229, 1, 0, 0, 0, 233, 234, 1, 0, 0, 0, 234, 237, 1, 0, 
	0, 0, 235, 236, 5, 71, 0, 0, 236, 238, 3, 32, 16, 0, 237, 235, 1, 0, 0, 
	0, 237, 238, 1, 0, 0, 0, 238, 240, 1, 0, 0, 0, 239, 241, 5, 53, 0, 0, 240, 
	239, 1, 0, 0, 0, 240, 241, 1, 0, 0, 0, 241, 31, 1, 0, 0, 0, 242, 247, 3, 
	34, 17, 0, 243, 244, 5, 71, 0, 0, 244, 246, 3, 34, 17, 0, 245, 243, 1, 
	0, 0, 0, 246, 249, 1, 0, 0, 0, 247, 245, 1, 0, 0, 0, 247, 248, 1, 0, 0, 
	0, 248, 33, 1, 0, 0, 0, 249, 247, 1, 0, 0, 0, 250, 255, 5, 54, 0, 0, 251, 
	252, 5, 75, 0, 0, 252, 253, 3, 36, 18, 0, 253, 254, 5, 76, 0, 0, 254, 256, 
	1, 0, 0, 0, 255, 251, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 35, 1, 0, 
	0, 0, 257, 258, 6, 18, -1, 0, 258, 259, 3, 40, 20, 0, 259, 260, 5, 24, 
	0, 0, 260, 261, 3, 40, 20, 0, 261, 267, 1, 0, 0, 0, 262, 263, 5, 68, 0, 
	0, 263, 264, 3, 36, 18, 0, 264, 265, 5, 69, 0, 0, 265, 267, 1, 0, 0, 0, 
	266, 257, 1, 0, 0, 0, 266, 262, 1, 0, 0, 0, 267, 276, 1, 0, 0, 0, 268, 
	269, 10, 3, 0, 0, 269, 270, 5, 20, 0, 0, 270, 275, 3, 36, 18, 4, 271, 272, 
	10, 2, 0, 0, 272, 273, 5, 21, 0, 0, 273, 275, 3, 36, 18, 3, 274, 268, 1, 
	0, 0, 0, 274, 271, 1, 0, 0, 0, 275, 278, 1, 0, 0, 0, 276, 274, 1, 0, 0, 
	0, 276, 277, 1, 0, 0, 0, 277, 37, 1, 0, 0, 0, 278, 276, 1, 0, 0, 0, 279, 
	280, 3, 30, 15, 0, 280, 281, 5, 23, 0, 0, 281, 308, 1, 0, 0, 0, 282, 283, 
	3, 42, 21, 0, 283, 284, 5, 24, 0, 0, 284, 285, 3, 42, 21, 0, 285, 308, 
	1, 0, 0, 0, 286, 287, 3, 42, 21, 0, 287, 288, 5, 25, 0, 0, 288, 289, 3, 
	54, 27, 0, 289, 308, 1, 0, 0, 0, 290, 291, 3, 42, 21, 0, 291, 292, 5, 26, 
	0, 0, 292, 293, 5, 68, 0, 0, 293, 294, 3, 44, 22, 0, 294, 295, 5, 69, 0, 
	0, 295, 308, 1, 0, 0, 0, 296, 297, 5, 54, 0, 0, 297, 298, 5, 19, 0, 0, 
	298, 301, 5, 54, 0, 0, 299, 300, 5, 28, 0, 0, 300, 302, 3, 32, 16, 0, 301, 
	299, 1, 0, 0, 0, 301, 302, 1, 0, 0, 0, 302, 308, 1, 0, 0, 0, 303, 304, 
	5, 68, 0, 0, 304, 305, 3, 38, 19, 0, 305, 306, 5, 69, 0, 0, 306, 308, 1, 
	0, 0, 0, 307, 279, 1, 0, 0, 0, 307, 282, 1, 0, 0, 0, 307, 286, 1, 0, 0, 
	0, 307, 290, 1, 0, 0, 0, 307, 296, 1, 0, 0, 0, 307, 303, 1, 0, 0, 0, 308, 
	39, 1, 0, 0, 0, 309, 313, 3, 48, 24, 0, 310, 313, 3, 32, 16, 0, 311, 313, 
	5, 52, 0, 0, 312, 309, 1, 0, 0, 0, 312, 310, 1, 0, 0, 0, 312, 311, 1, 0, 
	0, 0, 313, 41, 1, 0, 0, 0, 314, 318, 3, 48, 24, 0, 315, 318, 3, 30, 15, 
	0, 316, 318, 5, 52, 0, 0, 317, 314, 1, 0, 0, 0, 317, 315, 1, 0, 0, 0, 317, 
	316, 1, 0, 0, 0, 318, 43, 1, 0, 0, 0, 319, 329, 3, 2, 1, 0, 320, 325, 3, 
	46, 23, 0, 321, 322, 5, 70, 0, 0, 322, 324, 3, 46, 23, 0, 323, 321, 1, 
	0, 0, 0, 324, 327, 1, 0, 0, 0, 325, 323, 1, 0, 0, 0, 325, 326, 1, 0, 0, 
	0, 326, 329, 1, 0, 0, 0, 327, 325, 1, 0, 0, 0, 328, 319, 1, 0, 0, 0, 328, 
	320, 1, 0, 0, 0, 329, 45, 1, 0, 0, 0, 330, 333, 3, 48, 24, 0, 331, 333, 
	5, 52, 0, 0, 332, 330, 1, 0, 0, 0, 332, 331, 1, 0, 0, 0, 333, 47, 1, 0, 
	0, 0, 334, 340, 5, 55, 0, 0, 335, 340, 3, 50, 25, 0, 336, 340, 3, 52, 26, 
	0, 337, 340, 5, 17, 0, 0, 338, 340, 5, 14, 0, 0, 339, 334, 1, 0, 0, 0, 
	339, 335, 1, 0, 0, 0, 339, 336, 1, 0, 0, 0, 339, 337, 1, 0, 0, 0, 339, 
	338, 1, 0, 0, 0, 340, 49, 1, 0, 0, 0, 341, 343, 5, 74, 0, 0, 342, 341, 
	1, 0, 0, 0, 342, 343, 1, 0, 0, 0, 343, 344, 1, 0, 0, 0, 344, 350, 5, 56, 
	0, 0, 345, 347, 5, 74, 0, 0, 346, 345, 1, 0, 0, 0, 346, 347, 1, 0, 0, 0, 
	347, 348, 1, 0, 0, 0, 348, 350, 5, 58, 0, 0, 349, 342, 1, 0, 0, 0, 349, 
	346, 1, 0, 0, 0, 350, 51, 1, 0, 0, 0, 351, 353, 5, 74, 0, 0, 352, 351, 
	1, 0, 0, 0, 352, 353, 1, 0, 0, 0, 353, 354, 1, 0, 0, 0, 354, 360, 5, 57, 
	0, 0, 355, 357, 5, 74, 0, 0, 356, 355, 1, 0, 0, 0, 356, 357, 1, 0, 0, 0, 
	357, 358, 1, 0, 0, 0, 358, 360, 5, 59, 0, 0, 359, 352, 1, 0, 0, 0, 359, 
	356, 1, 0, 0, 0, 360, 53, 1, 0, 0, 0, 361, 362, 7, 1, 0, 0, 362, 55, 1, 
	0, 0, 0, 363, 367, 3, 50, 25, 0, 364, 367, 5, 52, 0, 0, 365, 367, 5, 54, 
	0, 0, 366, 363, 1, 0, 0, 0, 366, 364, 1, 0, 0, 0, 366, 365, 1, 0, 0, 0, 
	367, 57, 1, 0, 0, 0, 368, 373, 3, 50, 25, 0, 369, 373, 3, 52, 26, 0, 370, 
	373, 5, 52, 0, 0, 371, 373, 5, 54, 0, 0, 372, 368, 1, 0, 0, 0, 372, 369, 
	1, 0, 0, 0, 372, 370, 1, 0, 0, 0, 372, 371, 1, 0, 0, 0, 373, 59, 1, 0, 
	0, 0, 374, 375, 5, 33, 0, 0, 375, 376, 5, 68, 0, 0, 376, 377, 3, 54, 27, 
	0, 377, 378, 5, 69, 0, 0, 378, 454, 1, 0, 0, 0, 379, 380, 5, 34, 0, 0, 
	380, 381, 5, 68, 0, 0, 381, 382, 3, 54, 27, 0, 382, 383, 5, 70, 0, 0, 383, 
	384, 3, 56, 28, 0, 384, 385, 5, 69, 0, 0, 385, 454, 1, 0, 0, 0, 386, 387, 
	5, 35, 0, 0, 387, 388, 5, 68, 0, 0, 388, 389, 3, 54, 27, 0, 389, 390, 5, 
	70, 0, 0, 390, 391, 3, 56, 28, 0, 391, 392, 5, 70, 0, 0, 392, 393, 3, 56, 
	28, 0, 393, 394, 5, 69, 0, 0, 394, 454, 1, 0, 0, 0, 395, 396, 5, 36, 0, 
	0, 396, 397, 5, 68, 0, 0, 397, 402, 3, 54, 27, 0, 398, 399, 5, 70, 0, 0, 
	399, 401, 3, 54, 27, 0, 400, 398, 1, 0, 0, 0, 401, 404, 1, 0, 0, 0, 402, 
	400, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0, 403, 405, 1, 0, 0, 0, 404, 402, 
	1, 0, 0, 0, 405, 406, 5, 69, 0, 0, 406, 454, 1, 0, 0, 0, 407, 408, 5, 37, 
	0, 0, 408, 409, 5, 68, 0, 0, 409, 410, 3, 54, 27, 0, 410, 411, 5, 70, 0, 
	0, 411, 416, 3, 54, 27, 0, 412, 413, 5, 70, 0, 0, 413, 415, 3, 54, 27, 
	0, 414, 412, 1, 0, 0, 0, 415, 418, 1, 0, 0, 0, 416, 414, 1, 0, 0, 0, 416, 
	417, 1, 0, 0, 0, 417, 419, 1, 0, 0, 0, 418, 416, 1, 0, 0, 0, 419, 420, 
	5, 69, 0, 0, 420, 454, 1, 0, 0, 0, 421, 422, 5, 38, 0, 0, 422, 423, 5, 
	68, 0, 0, 423, 424, 3, 58, 29, 0, 424, 425, 5, 69, 0, 0, 425, 454, 1, 0, 
	0, 0, 426, 427, 5, 39, 0, 0, 427, 428, 5, 68, 0, 0, 428, 429, 3, 58, 29, 
	0, 429, 430, 5, 70, 0, 0, 430, 431, 3, 58, 29, 0, 431, 432, 5, 69, 0, 0, 
	432, 454, 1, 0, 0, 0, 433, 434, 5, 40, 0, 0, 434, 435, 5, 68, 0, 0, 435, 
	436, 3, 58, 29, 0, 436, 437, 5, 69, 0, 0, 437, 454, 1, 0, 0, 0, 438, 439, 
	5, 41, 0, 0, 439, 440, 5, 68, 0, 0, 440, 441, 3, 58, 29, 0, 441, 442, 5, 
	69, 0, 0, 442, 454, 1, 0, 0, 0, 443, 444, 5, 42, 0, 0, 444, 445, 5, 68, 
	0, 0, 445, 446, 3, 58, 29, 0, 446, 447, 5, 70, 0, 0, 447, 448, 3, 56, 28, 
	0, 448, 449, 5, 69, 0, 0, 449, 454, 1, 0, 0, 0, 450, 454, 5, 43, 0, 0, 
	451, 454, 5, 44, 0, 0, 452, 454, 5, 45, 0, 0, 453, 374, 1, 0, 0, 0, 453, 
	379, 1, 0, 0, 0, 453, 386, 1, 0, 0, 0, 453, 395, 1, 0, 0, 0, 453, 407, 
	1, 0, 0, 0, 453, 421, 1, 0, 0, 0, 453, 426, 1, 0, 0, 0, 453, 433, 1, 0, 
	0, 0, 453, 438, 1, 0, 0, 0, 453, 443, 1, 0, 0, 0, 453, 450, 1, 0, 0, 0, 
	453, 451, 1, 0, 0, 0, 453, 452, 1, 0, 0, 0, 454, 61, 1, 0, 0, 0, 455, 456, 
	5, 47, 0, 0, 456, 462, 5, 68, 0, 0, 457, 459, 5, 13, 0, 0, 458, 457, 1, 
	0, 0, 0, 458, 459, 1, 0, 0, 0, 459, 460, 1, 0, 0, 0, 460, 463, 5, 72, 0, 
	0, 461, 463, 3, 30, 15, 0, 462, 458, 1, 0, 0, 0, 462, 461, 1, 0, 0, 0, 
	463, 464, 1, 0, 0, 0, 464, 471, 5, 69, 0, 0, 465, 466, 7, 2, 0, 0, 466, 
	467, 5, 68, 0, 0, 467, 468, 3, 30, 15, 0, 468, 469, 5, 69, 0, 0, 469, 471, 
	1, 0, 0, 0, 470, 455, 1, 0, 0, 0, 470, 465, 1, 0, 0, 0, 471, 63, 1, 0, 
	0, 0, 472, 473, 7, 3, 0, 0, 473, 65, 1, 0, 0, 0, 60, 74, 78, 81, 85, 89, 
	92, 94, 98, 105, 112, 126, 135, 142, 148, 150, 155, 158, 160, 166, 169, 
	175, 178, 184, 190, 194, 204, 212, 214, 219, 226, 233, 237, 240, 247, 255, 
	266, 274, 276, 301, 307, 312, 317, 325, 328, 332, 339, 342, 346, 349, 352, 
	356, 359, 366, 372, 402, 416, 453, 458, 462, 470,
}
  deserializer := antlr.NewATNDeserializer(nil)
  staticData.atn = deserializer.Deserialize(staticData.serializedATN)
  atn := staticData.atn
  staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
  decisionToDFA := staticData.decisionToDFA
  for index, state := range atn.DecisionToState {
    decisionToDFA[index] = antlr.NewDFA(state, index)
  }
}

// AQLParserInit initializes any static state used to implement AQLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewAQLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func AQLParserInit() {
  staticData := &AQLParserStaticData
  staticData.once.Do(aqlParserInit)
}

// NewAQLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewAQLParser(input antlr.TokenStream) *AQLParser {
	AQLParserInit()
	this := new(AQLParser)
	this.BaseParser = antlr.NewBaseParser(input)
  staticData := &AQLParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "AQL.g4"

	return this
}


// AQLParser tokens.
const (
	AQLParserEOF = antlr.TokenEOF
	AQLParserWS = 1
	AQLParserUNICODE_BOM = 2
	AQLParserSELECT = 3
	AQLParserAS = 4
	AQLParserFROM = 5
	AQLParserWHERE = 6
	AQLParserORDER_BY = 7
	AQLParserGROUP_BY = 8
	AQLParserDESC = 9
	AQLParserASC = 10
	AQLParserLIMIT = 11
	AQLParserOFFSET = 12
	AQLParserDISTINCT = 13
	AQLParserNULL = 14
	AQLParserTRUE = 15
	AQLParserFALSE = 16
	AQLParserBOOLEAN = 17
	AQLParserALL_VERSIONS = 18
	AQLParserCONTAINS = 19
	AQLParserAND = 20
	AQLParserOR = 21
	AQLParserNOT = 22
	AQLParserEXISTS = 23
	AQLParserCOMPARISON_OPERATOR = 24
	AQLParserLIKE = 25
	AQLParserIN = 26
	AQLParserON = 27
	AQLParserAT = 28
	AQLParserJOIN = 29
	AQLParserLEFT = 30
	AQLParserUNION = 31
	AQLParserALL = 32
	AQLParserLENGTH = 33
	AQLParserPOSITION = 34
	AQLParserSUBSTRING = 35
	AQLParserCONCAT = 36
	AQLParserCONCAT_WS = 37
	AQLParserABS = 38
	AQLParserMOD = 39
	AQLParserCEIL = 40
	AQLParserFLOOR = 41
	AQLParserROUND = 42
	AQLParserCURRENT_DATE = 43
	AQLParserCURRENT_TIME = 44
	AQLParserCURRENT_DATE_TIME = 45
	AQLParserCURRENT_TIMEZONE = 46
	AQLParserCOUNT = 47
	AQLParserMIN = 48
	AQLParserMAX = 49
	AQLParserSUM = 50
	AQLParserAVG = 51
	AQLParserPARAMETER = 52
	AQLParserCAST = 53
	AQLParserIDENTIFIER = 54
	AQLParserSTRING = 55
	AQLParserINTEGER = 56
	AQLParserFLOAT = 57
	AQLParserSCI_INTEGER = 58
	AQLParserSCI_FLOAT = 59
	AQLParserSYM_COLON = 60
	AQLParserSYM_SEMICOLON = 61
	AQLParserSYM_LT = 62
	AQLParserSYM_GT = 63
	AQLParserSYM_LE = 64
	AQLParserSYM_GE = 65
	AQLParserSYM_NE = 66
	AQLParserSYM_EQ = 67
	AQLParserSYM_LEFT_PAREN = 68
	AQLParserSYM_RIGHT_PAREN = 69
	AQLParserSYM_COMMA = 70
	AQLParserSYM_SLASH = 71
	AQLParserSYM_ASTERISK = 72
	AQLParserSYM_PLUS = 73
	AQLParserSYM_MINUS = 74
	AQLParserSYM_LEFT_BRACKET = 75
	AQLParserSYM_RIGHT_BRACKET = 76
	AQLParserSYM_LEFT_CURLY = 77
	AQLParserSYM_RIGHT_CURLY = 78
	AQLParserSYM_DOUBLE_DASH = 79
)

// AQLParser rules.
const (
	AQLParserRULE_query = 0
	AQLParserRULE_selectQuery = 1
	AQLParserRULE_selectClause = 2
	AQLParserRULE_fromClause = 3
	AQLParserRULE_joinClause = 4
	AQLParserRULE_whereClause = 5
	AQLParserRULE_groupByClause = 6
	AQLParserRULE_orderByClause = 7
	AQLParserRULE_limitOffsetClause = 8
	AQLParserRULE_selectExpr = 9
	AQLParserRULE_fromExpr = 10
	AQLParserRULE_joinExpr = 11
	AQLParserRULE_whereExpr = 12
	AQLParserRULE_orderByExpr = 13
	AQLParserRULE_columnExpr = 14
	AQLParserRULE_identifiedPath = 15
	AQLParserRULE_objectPath = 16
	AQLParserRULE_pathPart = 17
	AQLParserRULE_pathCondition = 18
	AQLParserRULE_booleanCondition = 19
	AQLParserRULE_pathConditionOperand = 20
	AQLParserRULE_comparisonOperand = 21
	AQLParserRULE_inOperand = 22
	AQLParserRULE_inOperandValue = 23
	AQLParserRULE_primitive = 24
	AQLParserRULE_intPrimitive = 25
	AQLParserRULE_floatPrimitive = 26
	AQLParserRULE_stringOperand = 27
	AQLParserRULE_intOperand = 28
	AQLParserRULE_numbericOperand = 29
	AQLParserRULE_functionCall = 30
	AQLParserRULE_aggregateFunctionCall = 31
	AQLParserRULE_limitOperand = 32
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SelectQuery() ISelectQueryContext
	EOF() antlr.TerminalNode

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) SelectQuery() ISelectQueryContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectQueryContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectQueryContext)
}

func (s *QueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(AQLParserEOF, 0)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitQuery(s)
	}
}




func (p *AQLParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, AQLParserRULE_query)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(66)
		p.SelectQuery()
	}
	{
		p.SetState(67)
		p.Match(AQLParserEOF)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ISelectQueryContext is an interface to support dynamic dispatch.
type ISelectQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SelectClause() ISelectClauseContext
	FromClause() IFromClauseContext
	AllJoinClause() []IJoinClauseContext
	JoinClause(i int) IJoinClauseContext
	WhereClause() IWhereClauseContext
	GroupByClause() IGroupByClauseContext
	UNION() antlr.TerminalNode
	SelectQuery() ISelectQueryContext
	ALL() antlr.TerminalNode
	OrderByClause() IOrderByClauseContext
	LimitOffsetClause() ILimitOffsetClauseContext

	// IsSelectQueryContext differentiates from other interfaces.
	IsSelectQueryContext()
}

type SelectQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectQueryContext() *SelectQueryContext {
	var p = new(SelectQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_selectQuery
	return p
}

func InitEmptySelectQueryContext(p *SelectQueryContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_selectQuery
}

func (*SelectQueryContext) IsSelectQueryContext() {}

func NewSelectQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectQueryContext {
	var p = new(SelectQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_selectQuery

	return p
}

func (s *SelectQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectQueryContext) SelectClause() ISelectClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectClauseContext)
}

func (s *SelectQueryContext) FromClause() IFromClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *SelectQueryContext) AllJoinClause() []IJoinClauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IJoinClauseContext); ok {
			len++
		}
	}

	tst := make([]IJoinClauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IJoinClauseContext); ok {
			tst[i] = t.(IJoinClauseContext)
			i++
		}
	}

	return tst
}

func (s *SelectQueryContext) JoinClause(i int) IJoinClauseContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoinClauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoinClauseContext)
}

func (s *SelectQueryContext) WhereClause() IWhereClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereClauseContext)
}

func (s *SelectQueryContext) GroupByClause() IGroupByClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByClauseContext)
}

func (s *SelectQueryContext) UNION() antlr.TerminalNode {
	return s.GetToken(AQLParserUNION, 0)
}

func (s *SelectQueryContext) SelectQuery() ISelectQueryContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectQueryContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectQueryContext)
}

func (s *SelectQueryContext) ALL() antlr.TerminalNode {
	return s.GetToken(AQLParserALL, 0)
}

func (s *SelectQueryContext) OrderByClause() IOrderByClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrderByClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrderByClauseContext)
}

func (s *SelectQueryContext) LimitOffsetClause() ILimitOffsetClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitOffsetClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitOffsetClauseContext)
}

func (s *SelectQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterSelectQuery(s)
	}
}

func (s *SelectQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitSelectQuery(s)
	}
}




func (p *AQLParser) SelectQuery() (localctx ISelectQueryContext) {
	localctx = NewSelectQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, AQLParserRULE_selectQuery)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.SelectClause()
	}
	{
		p.SetState(70)
		p.FromClause()
	}
	p.SetState(74)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AQLParserJOIN || _la == AQLParserLEFT {
		{
			p.SetState(71)
			p.JoinClause()
		}


		p.SetState(76)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserWHERE {
		{
			p.SetState(77)
			p.WhereClause()
		}

	}
	p.SetState(81)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserGROUP_BY {
		{
			p.SetState(80)
			p.GroupByClause()
		}

	}
	p.SetState(94)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserUNION:
		{
			p.SetState(83)
			p.Match(AQLParserUNION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(85)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserALL {
			{
				p.SetState(84)
				p.Match(AQLParserALL)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(87)
			p.SelectQuery()
		}



	case AQLParserEOF, AQLParserORDER_BY, AQLParserLIMIT, AQLParserOFFSET, AQLParserSYM_RIGHT_PAREN:
		p.SetState(89)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserORDER_BY {
			{
				p.SetState(88)
				p.OrderByClause()
			}

		}
		p.SetState(92)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserLIMIT || _la == AQLParserOFFSET {
			{
				p.SetState(91)
				p.LimitOffsetClause()
			}

		}




	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ISelectClauseContext is an interface to support dynamic dispatch.
type ISelectClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELECT() antlr.TerminalNode
	AllSelectExpr() []ISelectExprContext
	SelectExpr(i int) ISelectExprContext
	DISTINCT() antlr.TerminalNode
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode

	// IsSelectClauseContext differentiates from other interfaces.
	IsSelectClauseContext()
}

type SelectClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectClauseContext() *SelectClauseContext {
	var p = new(SelectClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_selectClause
	return p
}

func InitEmptySelectClauseContext(p *SelectClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_selectClause
}

func (*SelectClauseContext) IsSelectClauseContext() {}

func NewSelectClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectClauseContext {
	var p = new(SelectClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_selectClause

	return p
}

func (s *SelectClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectClauseContext) SELECT() antlr.TerminalNode {
	return s.GetToken(AQLParserSELECT, 0)
}

func (s *SelectClauseContext) AllSelectExpr() []ISelectExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISelectExprContext); ok {
			len++
		}
	}

	tst := make([]ISelectExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISelectExprContext); ok {
			tst[i] = t.(ISelectExprContext)
			i++
		}
	}

	return tst
}

func (s *SelectClauseContext) SelectExpr(i int) ISelectExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectExprContext)
}

func (s *SelectClauseContext) DISTINCT() antlr.TerminalNode {
	return s.GetToken(AQLParserDISTINCT, 0)
}

func (s *SelectClauseContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *SelectClauseContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
}

func (s *SelectClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterSelectClause(s)
	}
}

func (s *SelectClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitSelectClause(s)
	}
}




func (p *AQLParser) SelectClause() (localctx ISelectClauseContext) {
	localctx = NewSelectClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, AQLParserRULE_selectClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(96)
		p.Match(AQLParserSELECT)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(98)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserDISTINCT {
		{
			p.SetState(97)
			p.Match(AQLParserDISTINCT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	}
	{
		p.SetState(100)
		p.SelectExpr()
	}
	p.SetState(105)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AQLParserSYM_COMMA {
		{
			p.SetState(101)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(102)
			p.SelectExpr()
		}


		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IFromClauseContext is an interface to support dynamic dispatch.
type IFromClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FROM() antlr.TerminalNode
	FromExpr() IFromExprContext

	// IsFromClauseContext differentiates from other interfaces.
	IsFromClauseContext()
}

type FromClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFromClauseContext() *FromClauseContext {
	var p = new(FromClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_fromClause
	return p
}

func InitEmptyFromClauseContext(p *FromClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_fromClause
}

func (*FromClauseContext) IsFromClauseContext() {}

func NewFromClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromClauseContext {
	var p = new(FromClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_fromClause

	return p
}

func (s *FromClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FromClauseContext) FROM() antlr.TerminalNode {
	return s.GetToken(AQLParserFROM, 0)
}

func (s *FromClauseContext) FromExpr() IFromExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromExprContext)
}

func (s *FromClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FromClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterFromClause(s)
	}
}

func (s *FromClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitFromClause(s)
	}
}




func (p *AQLParser) FromClause() (localctx IFromClauseContext) {
	localctx = NewFromClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, AQLParserRULE_fromClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(108)
		p.Match(AQLParserFROM)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(109)
		p.FromExpr()
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IJoinClauseContext is an interface to support dynamic dispatch.
type IJoinClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	JOIN() antlr.TerminalNode
	JoinExpr() IJoinExprContext
	LEFT() antlr.TerminalNode

	// IsJoinClauseContext differentiates from other interfaces.
	IsJoinClauseContext()
}

type JoinClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJoinClauseContext() *JoinClauseContext {
	var p = new(JoinClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_joinClause
	return p
}

func InitEmptyJoinClauseContext(p *JoinClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_joinClause
}

func (*JoinClauseContext) IsJoinClauseContext() {}

func NewJoinClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JoinClauseContext {
	var p = new(JoinClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_joinClause

	return p
}

func (s *JoinClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *JoinClauseContext) JOIN() antlr.TerminalNode {
	return s.GetToken(AQLParserJOIN, 0)
}

func (s *JoinClauseContext) JoinExpr() IJoinExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoinExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoinExprContext)
}

func (s *JoinClauseContext) LEFT() antlr.TerminalNode {
	return s.GetToken(AQLParserLEFT, 0)
}

func (s *JoinClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *JoinClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterJoinClause(s)
	}
}

func (s *JoinClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitJoinClause(s)
	}
}




func (p *AQLParser) JoinClause() (localctx IJoinClauseContext) {
	localctx = NewJoinClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, AQLParserRULE_joinClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserLEFT {
		{
			p.SetState(111)
			p.Match(AQLParserLEFT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	}
	{
		p.SetState(114)
		p.Match(AQLParserJOIN)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(115)
		p.JoinExpr()
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IWhereClauseContext is an interface to support dynamic dispatch.
type IWhereClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	WhereExpr() IWhereExprContext

	// IsWhereClauseContext differentiates from other interfaces.
	IsWhereClauseContext()
}

type WhereClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereClauseContext() *WhereClauseContext {
	var p = new(WhereClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_whereClause
	return p
}

func InitEmptyWhereClauseContext(p *WhereClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_whereClause
}

func (*WhereClauseContext) IsWhereClauseContext() {}

func NewWhereClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereClauseContext {
	var p = new(WhereClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_whereClause

	return p
}

func (s *WhereClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereClauseContext) WHERE() antlr.TerminalNode {
	return s.GetToken(AQLParserWHERE, 0)
}

func (s *WhereClauseContext) WhereExpr() IWhereExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereExprContext)
}

func (s *WhereClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WhereClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterWhereClause(s)
	}
}

func (s *WhereClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitWhereClause(s)
	}
}




func (p *AQLParser) WhereClause() (localctx IWhereClauseContext) {
	localctx = NewWhereClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, AQLParserRULE_whereClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		p.Match(AQLParserWHERE)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(118)
		p.whereExpr(0)
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IGroupByClauseContext is an interface to support dynamic dispatch.
type IGroupByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GROUP_BY() antlr.TerminalNode
	AllIdentifiedPath() []IIdentifiedPathContext
	IdentifiedPath(i int) IIdentifiedPathContext
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode

	// IsGroupByClauseContext differentiates from other interfaces.
	IsGroupByClauseContext()
}

type GroupByClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByClauseContext() *GroupByClauseContext {
	var p = new(GroupByClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_groupByClause
	return p
}

func InitEmptyGroupByClauseContext(p *GroupByClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_groupByClause
}

func (*GroupByClauseContext) IsGroupByClauseContext() {}

func NewGroupByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByClauseContext {
	var p = new(GroupByClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_groupByClause

	return p
}

func (s *GroupByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByClauseContext) GROUP_BY() antlr.TerminalNode {
	return s.GetToken(AQLParserGROUP_BY, 0)
}

func (s *GroupByClauseContext) AllIdentifiedPath() []IIdentifiedPathContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIdentifiedPathContext); ok {
			len++
		}
	}

	tst := make([]IIdentifiedPathContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIdentifiedPathContext); ok {
			tst[i] = t.(IIdentifiedPathContext)
			i++
		}
	}

	return tst
}

func (s *GroupByClauseContext) IdentifiedPath(i int) IIdentifiedPathContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedPathContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedPathContext)
}

func (s *GroupByClauseContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *GroupByClauseContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
}

func (s *GroupByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *GroupByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterGroupByClause(s)
	}
}

func (s *GroupByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitGroupByClause(s)
	}
}




func (p *AQLParser) GroupByClause() (localctx IGroupByClauseContext) {
	localctx = NewGroupByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, AQLParserRULE_groupByClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(120)
		p.Match(AQLParserGROUP_BY)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(121)
		p.IdentifiedPath()
	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AQLParserSYM_COMMA {
		{
			p.SetState(122)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(123)
			p.IdentifiedPath()
		}


		p.SetState(128)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IOrderByClauseContext is an interface to support dynamic dispatch.
type IOrderByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ORDER_BY() antlr.TerminalNode
	AllOrderByExpr() []IOrderByExprContext
	OrderByExpr(i int) IOrderByExprContext
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode

	// IsOrderByClauseContext differentiates from other interfaces.
	IsOrderByClauseContext()
}

type OrderByClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrderByClauseContext() *OrderByClauseContext {
	var p = new(OrderByClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_orderByClause
	return p
}

func InitEmptyOrderByClauseContext(p *OrderByClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_orderByClause
}

func (*OrderByClauseContext) IsOrderByClauseContext() {}

func NewOrderByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByClauseContext {
	var p = new(OrderByClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_orderByClause

	return p
}

func (s *OrderByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByClauseContext) ORDER_BY() antlr.TerminalNode {
	return s.GetToken(AQLParserORDER_BY, 0)
}

func (s *OrderByClauseContext) AllOrderByExpr() []IOrderByExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOrderByExprContext); ok {
			len++
		}
	}

	tst := make([]IOrderByExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOrderByExprContext); ok {
			tst[i] = t.(IOrderByExprContext)
			i++
		}
	}

	return tst
}

func (s *OrderByClauseContext) OrderByExpr(i int) IOrderByExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrderByExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrderByExprContext)
}

func (s *OrderByClauseContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *OrderByClauseContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
}

func (s *OrderByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *OrderByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterOrderByClause(s)
	}
}

func (s *OrderByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitOrderByClause(s)
	}
}




func (p *AQLParser) OrderByClause() (localctx IOrderByClauseContext) {
	localctx = NewOrderByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, AQLParserRULE_orderByClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(129)
		p.Match(AQLParserORDER_BY)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(130)
		p.OrderByExpr()
	}
	p.SetState(135)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AQLParserSYM_COMMA {
		{
			p.SetState(131)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(132)
			p.OrderByExpr()
		}


		p.SetState(137)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ILimitOffsetClauseContext is an interface to support dynamic dispatch.
type ILimitOffsetClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeftLimit returns the leftLimit rule contexts.
	GetLeftLimit() ILimitOperandContext

	// GetRightOffset returns the rightOffset rule contexts.
	GetRightOffset() ILimitOperandContext

	// GetLeftOffset returns the leftOffset rule contexts.
	GetLeftOffset() ILimitOperandContext

	// GetRightLimit returns the rightLimit rule contexts.
	GetRightLimit() ILimitOperandContext


	// SetLeftLimit sets the leftLimit rule contexts.
	SetLeftLimit(ILimitOperandContext)

	// SetRightOffset sets the rightOffset rule contexts.
	SetRightOffset(ILimitOperandContext)

	// SetLeftOffset sets the leftOffset rule contexts.
	SetLeftOffset(ILimitOperandContext)

	// SetRightLimit sets the rightLimit rule contexts.
	SetRightLimit(ILimitOperandContext)


	// Getter signatures
	LIMIT() antlr.TerminalNode
	AllLimitOperand() []ILimitOperandContext
	LimitOperand(i int) ILimitOperandContext
	OFFSET() antlr.TerminalNode

	// IsLimitOffsetClauseContext differentiates from other interfaces.
	IsLimitOffsetClauseContext()
}

type LimitOffsetClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	leftLimit ILimitOperandContext 
	rightOffset ILimitOperandContext 
	leftOffset ILimitOperandContext 
	rightLimit ILimitOperandContext 
}

func NewEmptyLimitOffsetClauseContext() *LimitOffsetClauseContext {
	var p = new(LimitOffsetClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_limitOffsetClause
	return p
}

func InitEmptyLimitOffsetClauseContext(p *LimitOffsetClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_limitOffsetClause
}

func (*LimitOffsetClauseContext) IsLimitOffsetClauseContext() {}

func NewLimitOffsetClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitOffsetClauseContext {
	var p = new(LimitOffsetClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_limitOffsetClause

	return p
}

func (s *LimitOffsetClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitOffsetClauseContext) GetLeftLimit() ILimitOperandContext { return s.leftLimit }

func (s *LimitOffsetClauseContext) GetRightOffset() ILimitOperandContext { return s.rightOffset }

func (s *LimitOffsetClauseContext) GetLeftOffset() ILimitOperandContext { return s.leftOffset }

func (s *LimitOffsetClauseContext) GetRightLimit() ILimitOperandContext { return s.rightLimit }


func (s *LimitOffsetClauseContext) SetLeftLimit(v ILimitOperandContext) { s.leftLimit = v }

func (s *LimitOffsetClauseContext) SetRightOffset(v ILimitOperandContext) { s.rightOffset = v }

func (s *LimitOffsetClauseContext) SetLeftOffset(v ILimitOperandContext) { s.leftOffset = v }

func (s *LimitOffsetClauseContext) SetRightLimit(v ILimitOperandContext) { s.rightLimit = v }


func (s *LimitOffsetClauseContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(AQLParserLIMIT, 0)
}

func (s *LimitOffsetClauseContext) AllLimitOperand() []ILimitOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILimitOperandContext); ok {
			len++
		}
	}

	tst := make([]ILimitOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILimitOperandContext); ok {
			tst[i] = t.(ILimitOperandContext)
			i++
		}
	}

	return tst
}

func (s *LimitOffsetClauseContext) LimitOperand(i int) ILimitOperandContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitOperandContext)
}

func (s *LimitOffsetClauseContext) OFFSET() antlr.TerminalNode {
	return s.GetToken(AQLParserOFFSET, 0)
}

func (s *LimitOffsetClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitOffsetClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LimitOffsetClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterLimitOffsetClause(s)
	}
}

func (s *LimitOffsetClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitLimitOffsetClause(s)
	}
}




func (p *AQLParser) LimitOffsetClause() (localctx ILimitOffsetClauseContext) {
	localctx = NewLimitOffsetClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, AQLParserRULE_limitOffsetClause)
	var _la int

	p.SetState(150)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserLIMIT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(138)
			p.Match(AQLParserLIMIT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(139)

			var _x = p.LimitOperand()


			localctx.(*LimitOffsetClauseContext).leftLimit = _x
		}
		p.SetState(142)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserOFFSET {
			{
				p.SetState(140)
				p.Match(AQLParserOFFSET)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(141)

				var _x = p.LimitOperand()


				localctx.(*LimitOffsetClauseContext).rightOffset = _x
			}

		}


	case AQLParserOFFSET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(144)
			p.Match(AQLParserOFFSET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(145)

			var _x = p.LimitOperand()


			localctx.(*LimitOffsetClauseContext).leftOffset = _x
		}
		p.SetState(148)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserLIMIT {
			{
				p.SetState(146)
				p.Match(AQLParserLIMIT)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(147)

				var _x = p.LimitOperand()


				localctx.(*LimitOffsetClauseContext).rightLimit = _x
			}

		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ISelectExprContext is an interface to support dynamic dispatch.
type ISelectExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SYM_ASTERISK() antlr.TerminalNode
	ColumnExpr() IColumnExprContext
	IDENTIFIER() antlr.TerminalNode
	AS() antlr.TerminalNode

	// IsSelectExprContext differentiates from other interfaces.
	IsSelectExprContext()
}

type SelectExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectExprContext() *SelectExprContext {
	var p = new(SelectExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_selectExpr
	return p
}

func InitEmptySelectExprContext(p *SelectExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_selectExpr
}

func (*SelectExprContext) IsSelectExprContext() {}

func NewSelectExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectExprContext {
	var p = new(SelectExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_selectExpr

	return p
}

func (s *SelectExprContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectExprContext) SYM_ASTERISK() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_ASTERISK, 0)
}

func (s *SelectExprContext) ColumnExpr() IColumnExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnExprContext)
}

func (s *SelectExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
}

func (s *SelectExprContext) AS() antlr.TerminalNode {
	return s.GetToken(AQLParserAS, 0)
}

func (s *SelectExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterSelectExpr(s)
	}
}

func (s *SelectExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitSelectExpr(s)
	}
}




func (p *AQLParser) SelectExpr() (localctx ISelectExprContext) {
	localctx = NewSelectExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, AQLParserRULE_selectExpr)
	var _la int

	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserSYM_ASTERISK:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(152)
			p.Match(AQLParserSYM_ASTERISK)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserNULL, AQLParserBOOLEAN, AQLParserLENGTH, AQLParserPOSITION, AQLParserSUBSTRING, AQLParserCONCAT, AQLParserCONCAT_WS, AQLParserABS, AQLParserMOD, AQLParserCEIL, AQLParserFLOOR, AQLParserROUND, AQLParserCURRENT_DATE, AQLParserCURRENT_TIME, AQLParserCURRENT_DATE_TIME, AQLParserCOUNT, AQLParserMIN, AQLParserMAX, AQLParserSUM, AQLParserAVG, AQLParserPARAMETER, AQLParserIDENTIFIER, AQLParserSTRING, AQLParserINTEGER, AQLParserFLOAT, AQLParserSCI_INTEGER, AQLParserSCI_FLOAT, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(153)
			p.ColumnExpr()
		}
		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserAS || _la == AQLParserIDENTIFIER {
			p.SetState(155)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			if _la == AQLParserAS {
				{
					p.SetState(154)
					p.Match(AQLParserAS)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}

			}
			{
				p.SetState(157)
				p.Match(AQLParserIDENTIFIER)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IFromExprContext is an interface to support dynamic dispatch.
type IFromExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAlias returns the alias token.
	GetAlias() antlr.Token 


	// SetAlias sets the alias token.
	SetAlias(antlr.Token) 


	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	SYM_LEFT_BRACKET() antlr.TerminalNode
	ALL_VERSIONS() antlr.TerminalNode
	SYM_RIGHT_BRACKET() antlr.TerminalNode

	// IsFromExprContext differentiates from other interfaces.
	IsFromExprContext()
}

type FromExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	alias antlr.Token
}

func NewEmptyFromExprContext() *FromExprContext {
	var p = new(FromExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_fromExpr
	return p
}

func InitEmptyFromExprContext(p *FromExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_fromExpr
}

func (*FromExprContext) IsFromExprContext() {}

func NewFromExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromExprContext {
	var p = new(FromExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_fromExpr

	return p
}

func (s *FromExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FromExprContext) GetAlias() antlr.Token { return s.alias }


func (s *FromExprContext) SetAlias(v antlr.Token) { s.alias = v }


func (s *FromExprContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(AQLParserIDENTIFIER)
}

func (s *FromExprContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, i)
}

func (s *FromExprContext) SYM_LEFT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_BRACKET, 0)
}

func (s *FromExprContext) ALL_VERSIONS() antlr.TerminalNode {
	return s.GetToken(AQLParserALL_VERSIONS, 0)
}

func (s *FromExprContext) SYM_RIGHT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_BRACKET, 0)
}

func (s *FromExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FromExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterFromExpr(s)
	}
}

func (s *FromExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitFromExpr(s)
	}
}




func (p *AQLParser) FromExpr() (localctx IFromExprContext) {
	localctx = NewFromExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, AQLParserRULE_fromExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(162)
		p.Match(AQLParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(166)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserSYM_LEFT_BRACKET {
		{
			p.SetState(163)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(164)
			p.Match(AQLParserALL_VERSIONS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(165)
			p.Match(AQLParserSYM_RIGHT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	}
	p.SetState(169)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserIDENTIFIER {
		{
			p.SetState(168)

			var _m = p.Match(AQLParserIDENTIFIER)

			localctx.(*FromExprContext).alias = _m
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IJoinExprContext is an interface to support dynamic dispatch.
type IJoinExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAlias returns the alias token.
	GetAlias() antlr.Token 

	// GetSource returns the source token.
	GetSource() antlr.Token 


	// SetAlias sets the alias token.
	SetAlias(antlr.Token) 

	// SetSource sets the source token.
	SetSource(antlr.Token) 


	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	ON() antlr.TerminalNode
	SYM_LEFT_BRACKET() antlr.TerminalNode
	ALL_VERSIONS() antlr.TerminalNode
	SYM_RIGHT_BRACKET() antlr.TerminalNode
	IN() antlr.TerminalNode
	AT() antlr.TerminalNode
	IdentifiedPath() IIdentifiedPathContext

	// IsJoinExprContext differentiates from other interfaces.
	IsJoinExprContext()
}

type JoinExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	alias antlr.Token
	source antlr.Token
}

func NewEmptyJoinExprContext() *JoinExprContext {
	var p = new(JoinExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_joinExpr
	return p
}

func InitEmptyJoinExprContext(p *JoinExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_joinExpr
}

func (*JoinExprContext) IsJoinExprContext() {}

func NewJoinExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JoinExprContext {
	var p = new(JoinExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_joinExpr

	return p
}

func (s *JoinExprContext) GetParser() antlr.Parser { return s.parser }

func (s *JoinExprContext) GetAlias() antlr.Token { return s.alias }

func (s *JoinExprContext) GetSource() antlr.Token { return s.source }


func (s *JoinExprContext) SetAlias(v antlr.Token) { s.alias = v }

func (s *JoinExprContext) SetSource(v antlr.Token) { s.source = v }


func (s *JoinExprContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(AQLParserIDENTIFIER)
}

func (s *JoinExprContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, i)
}

func (s *JoinExprContext) ON() antlr.TerminalNode {
	return s.GetToken(AQLParserON, 0)
}

func (s *JoinExprContext) SYM_LEFT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_BRACKET, 0)
}

func (s *JoinExprContext) ALL_VERSIONS() antlr.TerminalNode {
	return s.GetToken(AQLParserALL_VERSIONS, 0)
}

func (s *JoinExprContext) SYM_RIGHT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_BRACKET, 0)
}

func (s *JoinExprContext) IN() antlr.TerminalNode {
	return s.GetToken(AQLParserIN, 0)
}

func (s *JoinExprContext) AT() antlr.TerminalNode {
	return s.GetToken(AQLParserAT, 0)
}

func (s *JoinExprContext) IdentifiedPath() IIdentifiedPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedPathContext)
}

func (s *JoinExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *JoinExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterJoinExpr(s)
	}
}

func (s *JoinExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitJoinExpr(s)
	}
}




func (p *AQLParser) JoinExpr() (localctx IJoinExprContext) {
	localctx = NewJoinExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, AQLParserRULE_joinExpr)
	var _la int

	p.SetState(194)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(171)
			p.Match(AQLParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(175)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserSYM_LEFT_BRACKET {
			{
				p.SetState(172)
				p.Match(AQLParserSYM_LEFT_BRACKET)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(173)
				p.Match(AQLParserALL_VERSIONS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(174)
				p.Match(AQLParserSYM_RIGHT_BRACKET)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		p.SetState(178)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserIDENTIFIER {
			{
				p.SetState(177)

				var _m = p.Match(AQLParserIDENTIFIER)

				localctx.(*JoinExprContext).alias = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(180)
			p.Match(AQLParserON)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(181)

			var _m = p.Match(AQLParserIDENTIFIER)

			localctx.(*JoinExprContext).source = _m
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(182)
			p.Match(AQLParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(184)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserIDENTIFIER {
			{
				p.SetState(183)

				var _m = p.Match(AQLParserIDENTIFIER)

				localctx.(*JoinExprContext).alias = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(186)
			p.Match(AQLParserIN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(187)

			var _m = p.Match(AQLParserIDENTIFIER)

			localctx.(*JoinExprContext).source = _m
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(188)
			p.Match(AQLParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserIDENTIFIER {
			{
				p.SetState(189)

				var _m = p.Match(AQLParserIDENTIFIER)

				localctx.(*JoinExprContext).alias = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(192)
			p.Match(AQLParserAT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(193)
			p.IdentifiedPath()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IWhereExprContext is an interface to support dynamic dispatch.
type IWhereExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BooleanCondition() IBooleanConditionContext
	NOT() antlr.TerminalNode
	AllWhereExpr() []IWhereExprContext
	WhereExpr(i int) IWhereExprContext
	SYM_LEFT_PAREN() antlr.TerminalNode
	SYM_RIGHT_PAREN() antlr.TerminalNode
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode

	// IsWhereExprContext differentiates from other interfaces.
	IsWhereExprContext()
}

type WhereExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereExprContext() *WhereExprContext {
	var p = new(WhereExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_whereExpr
	return p
}

func InitEmptyWhereExprContext(p *WhereExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_whereExpr
}

func (*WhereExprContext) IsWhereExprContext() {}

func NewWhereExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereExprContext {
	var p = new(WhereExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_whereExpr

	return p
}

func (s *WhereExprContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereExprContext) BooleanCondition() IBooleanConditionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanConditionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanConditionContext)
}

func (s *WhereExprContext) NOT() antlr.TerminalNode {
	return s.GetToken(AQLParserNOT, 0)
}

func (s *WhereExprContext) AllWhereExpr() []IWhereExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IWhereExprContext); ok {
			len++
		}
	}

	tst := make([]IWhereExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IWhereExprContext); ok {
			tst[i] = t.(IWhereExprContext)
			i++
		}
	}

	return tst
}

func (s *WhereExprContext) WhereExpr(i int) IWhereExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereExprContext)
}

func (s *WhereExprContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *WhereExprContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *WhereExprContext) AND() antlr.TerminalNode {
	return s.GetToken(AQLParserAND, 0)
}

func (s *WhereExprContext) OR() antlr.TerminalNode {
	return s.GetToken(AQLParserOR, 0)
}

func (s *WhereExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WhereExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterWhereExpr(s)
	}
}

func (s *WhereExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitWhereExpr(s)
	}
}





func (p *AQLParser) WhereExpr() (localctx IWhereExprContext) {
	return p.whereExpr(0)
}

func (p *AQLParser) whereExpr(_p int) (localctx IWhereExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewWhereExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IWhereExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 24
	p.EnterRecursionRule(localctx, 24, AQLParserRULE_whereExpr, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(197)
			p.BooleanCondition()
		}


	case 2:
		{
			p.SetState(198)
			p.Match(AQLParserNOT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(199)
			p.whereExpr(4)
		}


	case 3:
		{
			p.SetState(200)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(201)
			p.whereExpr(0)
		}
		{
			p.SetState(202)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(214)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(212)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) {
			case 1:
				localctx = NewWhereExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_whereExpr)
				p.SetState(206)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(207)
					p.Match(AQLParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(208)
					p.whereExpr(4)
				}


			case 2:
				localctx = NewWhereExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_whereExpr)
				p.SetState(209)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(210)
					p.Match(AQLParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(211)
					p.whereExpr(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(216)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}



	errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IOrderByExprContext is an interface to support dynamic dispatch.
type IOrderByExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOrder returns the order token.
	GetOrder() antlr.Token 


	// SetOrder sets the order token.
	SetOrder(antlr.Token) 


	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	DESC() antlr.TerminalNode
	ASC() antlr.TerminalNode

	// IsOrderByExprContext differentiates from other interfaces.
	IsOrderByExprContext()
}

type OrderByExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	order antlr.Token
}

func NewEmptyOrderByExprContext() *OrderByExprContext {
	var p = new(OrderByExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_orderByExpr
	return p
}

func InitEmptyOrderByExprContext(p *OrderByExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_orderByExpr
}

func (*OrderByExprContext) IsOrderByExprContext() {}

func NewOrderByExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByExprContext {
	var p = new(OrderByExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_orderByExpr

	return p
}

func (s *OrderByExprContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByExprContext) GetOrder() antlr.Token { return s.order }


func (s *OrderByExprContext) SetOrder(v antlr.Token) { s.order = v }


func (s *OrderByExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
}

func (s *OrderByExprContext) DESC() antlr.TerminalNode {
	return s.GetToken(AQLParserDESC, 0)
}

func (s *OrderByExprContext) ASC() antlr.TerminalNode {
	return s.GetToken(AQLParserASC, 0)
}

func (s *OrderByExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *OrderByExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterOrderByExpr(s)
	}
}

func (s *OrderByExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitOrderByExpr(s)
	}
}




func (p *AQLParser) OrderByExpr() (localctx IOrderByExprContext) {
	localctx = NewOrderByExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, AQLParserRULE_orderByExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(217)
		p.Match(AQLParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(219)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserDESC || _la == AQLParserASC {
		{
			p.SetState(218)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*OrderByExprContext).order = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == AQLParserDESC || _la == AQLParserASC) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*OrderByExprContext).order = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IColumnExprContext is an interface to support dynamic dispatch.
type IColumnExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primitive() IPrimitiveContext
	PARAMETER() antlr.TerminalNode
	IdentifiedPath() IIdentifiedPathContext
	AggregateFunctionCall() IAggregateFunctionCallContext
	FunctionCall() IFunctionCallContext

	// IsColumnExprContext differentiates from other interfaces.
	IsColumnExprContext()
}

type ColumnExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnExprContext() *ColumnExprContext {
	var p = new(ColumnExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_columnExpr
	return p
}

func InitEmptyColumnExprContext(p *ColumnExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_columnExpr
}

func (*ColumnExprContext) IsColumnExprContext() {}

func NewColumnExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnExprContext {
	var p = new(ColumnExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_columnExpr

	return p
}

func (s *ColumnExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnExprContext) Primitive() IPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimitiveContext)
}

func (s *ColumnExprContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *ColumnExprContext) IdentifiedPath() IIdentifiedPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedPathContext)
}

func (s *ColumnExprContext) AggregateFunctionCall() IAggregateFunctionCallContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregateFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregateFunctionCallContext)
}

func (s *ColumnExprContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *ColumnExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ColumnExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterColumnExpr(s)
	}
}

func (s *ColumnExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitColumnExpr(s)
	}
}




func (p *AQLParser) ColumnExpr() (localctx IColumnExprContext) {
	localctx = NewColumnExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, AQLParserRULE_columnExpr)
	p.SetState(226)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserSTRING, AQLParserINTEGER, AQLParserFLOAT, AQLParserSCI_INTEGER, AQLParserSCI_FLOAT, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(221)
			p.Primitive()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(222)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(223)
			p.IdentifiedPath()
		}


	case AQLParserCOUNT, AQLParserMIN, AQLParserMAX, AQLParserSUM, AQLParserAVG:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(224)
			p.AggregateFunctionCall()
		}


	case AQLParserLENGTH, AQLParserPOSITION, AQLParserSUBSTRING, AQLParserCONCAT, AQLParserCONCAT_WS, AQLParserABS, AQLParserMOD, AQLParserCEIL, AQLParserFLOOR, AQLParserROUND, AQLParserCURRENT_DATE, AQLParserCURRENT_TIME, AQLParserCURRENT_DATE_TIME:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(225)
			p.FunctionCall()
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IIdentifiedPathContext is an interface to support dynamic dispatch.
type IIdentifiedPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	SYM_LEFT_BRACKET() antlr.TerminalNode
	PathCondition() IPathConditionContext
	SYM_RIGHT_BRACKET() antlr.TerminalNode
	SYM_SLASH() antlr.TerminalNode
	ObjectPath() IObjectPathContext
	CAST() antlr.TerminalNode

	// IsIdentifiedPathContext differentiates from other interfaces.
	IsIdentifiedPathContext()
}

type IdentifiedPathContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifiedPathContext() *IdentifiedPathContext {
	var p = new(IdentifiedPathContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_identifiedPath
	return p
}

func InitEmptyIdentifiedPathContext(p *IdentifiedPathContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_identifiedPath
}

func (*IdentifiedPathContext) IsIdentifiedPathContext() {}

func NewIdentifiedPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifiedPathContext {
	var p = new(IdentifiedPathContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_identifiedPath

	return p
}

func (s *IdentifiedPathContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifiedPathContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
}

func (s *IdentifiedPathContext) SYM_LEFT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_BRACKET, 0)
}

func (s *IdentifiedPathContext) PathCondition() IPathConditionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathConditionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathConditionContext)
}

func (s *IdentifiedPathContext) SYM_RIGHT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_BRACKET, 0)
}

func (s *IdentifiedPathContext) SYM_SLASH() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_SLASH, 0)
}

func (s *IdentifiedPathContext) ObjectPath() IObjectPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectPathContext)
}

func (s *IdentifiedPathContext) CAST() antlr.TerminalNode {
	return s.GetToken(AQLParserCAST, 0)
}

func (s *IdentifiedPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifiedPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IdentifiedPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterIdentifiedPath(s)
	}
}

func (s *IdentifiedPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitIdentifiedPath(s)
	}
}




func (p *AQLParser) IdentifiedPath() (localctx IIdentifiedPathContext) {
	localctx = NewIdentifiedPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, AQLParserRULE_identifiedPath)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(228)
		p.Match(AQLParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(233)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 30, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(229)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(230)
			p.pathCondition(0)
		}
		{
			p.SetState(231)
			p.Match(AQLParserSYM_RIGHT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

		} else if p.HasError() { // JIM
			goto errorExit
	}
	p.SetState(237)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(235)
			p.Match(AQLParserSYM_SLASH)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(236)
			p.ObjectPath()
		}

		} else if p.HasError() { // JIM
			goto errorExit
	}
	p.SetState(240)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(239)
			p.Match(AQLParserCAST)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

		} else if p.HasError() { // JIM
			goto errorExit
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IObjectPathContext is an interface to support dynamic dispatch.
type IObjectPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPathPart() []IPathPartContext
	PathPart(i int) IPathPartContext
	AllSYM_SLASH() []antlr.TerminalNode
	SYM_SLASH(i int) antlr.TerminalNode

	// IsObjectPathContext differentiates from other interfaces.
	IsObjectPathContext()
}

type ObjectPathContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjectPathContext() *ObjectPathContext {
	var p = new(ObjectPathContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_objectPath
	return p
}

func InitEmptyObjectPathContext(p *ObjectPathContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_objectPath
}

func (*ObjectPathContext) IsObjectPathContext() {}

func NewObjectPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectPathContext {
	var p = new(ObjectPathContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_objectPath

	return p
}

func (s *ObjectPathContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectPathContext) AllPathPart() []IPathPartContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPathPartContext); ok {
			len++
		}
	}

	tst := make([]IPathPartContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPathPartContext); ok {
			tst[i] = t.(IPathPartContext)
			i++
		}
	}

	return tst
}

func (s *ObjectPathContext) PathPart(i int) IPathPartContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathPartContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathPartContext)
}

func (s *ObjectPathContext) AllSYM_SLASH() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_SLASH)
}

func (s *ObjectPathContext) SYM_SLASH(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_SLASH, i)
}

func (s *ObjectPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ObjectPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterObjectPath(s)
	}
}

func (s *ObjectPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitObjectPath(s)
	}
}




func (p *AQLParser) ObjectPath() (localctx IObjectPathContext) {
	localctx = NewObjectPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, AQLParserRULE_objectPath)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(242)
		p.PathPart()
	}
	p.SetState(247)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(243)
				p.Match(AQLParserSYM_SLASH)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(244)
				p.PathPart()
			}


		}
		p.SetState(249)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IPathPartContext is an interface to support dynamic dispatch.
type IPathPartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	SYM_LEFT_BRACKET() antlr.TerminalNode
	PathCondition() IPathConditionContext
	SYM_RIGHT_BRACKET() antlr.TerminalNode

	// IsPathPartContext differentiates from other interfaces.
	IsPathPartContext()
}

type PathPartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathPartContext() *PathPartContext {
	var p = new(PathPartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathPart
	return p
}

func InitEmptyPathPartContext(p *PathPartContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathPart
}

func (*PathPartContext) IsPathPartContext() {}

func NewPathPartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathPartContext {
	var p = new(PathPartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_pathPart

	return p
}

func (s *PathPartContext) GetParser() antlr.Parser { return s.parser }

func (s *PathPartContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
}

func (s *PathPartContext) SYM_LEFT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_BRACKET, 0)
}

func (s *PathPartContext) PathCondition() IPathConditionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathConditionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathConditionContext)
}

func (s *PathPartContext) SYM_RIGHT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_BRACKET, 0)
}

func (s *PathPartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathPartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathPartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterPathPart(s)
	}
}

func (s *PathPartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitPathPart(s)
	}
}




func (p *AQLParser) PathPart() (localctx IPathPartContext) {
	localctx = NewPathPartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, AQLParserRULE_pathPart)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(250)
		p.Match(AQLParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(255)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(251)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(252)
			p.pathCondition(0)
		}
		{
			p.SetState(253)
			p.Match(AQLParserSYM_RIGHT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

		} else if p.HasError() { // JIM
			goto errorExit
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IPathConditionContext is an interface to support dynamic dispatch.
type IPathConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPathConditionOperand() []IPathConditionOperandContext
	PathConditionOperand(i int) IPathConditionOperandContext
	COMPARISON_OPERATOR() antlr.TerminalNode
	SYM_LEFT_PAREN() antlr.TerminalNode
	AllPathCondition() []IPathConditionContext
	PathCondition(i int) IPathConditionContext
	SYM_RIGHT_PAREN() antlr.TerminalNode
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode

	// IsPathConditionContext differentiates from other interfaces.
	IsPathConditionContext()
}

type PathConditionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathConditionContext() *PathConditionContext {
	var p = new(PathConditionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathCondition
	return p
}

func InitEmptyPathConditionContext(p *PathConditionContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathCondition
}

func (*PathConditionContext) IsPathConditionContext() {}

func NewPathConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathConditionContext {
	var p = new(PathConditionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_pathCondition

	return p
}

func (s *PathConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *PathConditionContext) AllPathConditionOperand() []IPathConditionOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPathConditionOperandContext); ok {
			len++
		}
	}

	tst := make([]IPathConditionOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPathConditionOperandContext); ok {
			tst[i] = t.(IPathConditionOperandContext)
			i++
		}
	}

	return tst
}

func (s *PathConditionContext) PathConditionOperand(i int) IPathConditionOperandContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathConditionOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathConditionOperandContext)
}

func (s *PathConditionContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(AQLParserCOMPARISON_OPERATOR, 0)
}

func (s *PathConditionContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *PathConditionContext) AllPathCondition() []IPathConditionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPathConditionContext); ok {
			len++
		}
	}

	tst := make([]IPathConditionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPathConditionContext); ok {
			tst[i] = t.(IPathConditionContext)
			i++
		}
	}

	return tst
}

func (s *PathConditionContext) PathCondition(i int) IPathConditionContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathConditionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathConditionContext)
}

func (s *PathConditionContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *PathConditionContext) AND() antlr.TerminalNode {
	return s.GetToken(AQLParserAND, 0)
}

func (s *PathConditionContext) OR() antlr.TerminalNode {
	return s.GetToken(AQLParserOR, 0)
}

func (s *PathConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterPathCondition(s)
	}
}

func (s *PathConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitPathCondition(s)
	}
}





func (p *AQLParser) PathCondition() (localctx IPathConditionContext) {
	return p.pathCondition(0)
}

func (p *AQLParser) pathCondition(_p int) (localctx IPathConditionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewPathConditionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IPathConditionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 36
	p.EnterRecursionRule(localctx, 36, AQLParserRULE_pathCondition, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(266)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserPARAMETER, AQLParserIDENTIFIER, AQLParserSTRING, AQLParserINTEGER, AQLParserFLOAT, AQLParserSCI_INTEGER, AQLParserSCI_FLOAT, AQLParserSYM_MINUS:
		{
			p.SetState(258)
			p.PathConditionOperand()
		}
		{
			p.SetState(259)
			p.Match(AQLParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(260)
			p.PathConditionOperand()
		}


	case AQLParserSYM_LEFT_PAREN:
		{
			p.SetState(262)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(263)
			p.pathCondition(0)
		}
		{
			p.SetState(264)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(276)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(274)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext()) {
			case 1:
				localctx = NewPathConditionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_pathCondition)
				p.SetState(268)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(269)
					p.Match(AQLParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(270)
					p.pathCondition(4)
				}


			case 2:
				localctx = NewPathConditionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_pathCondition)
				p.SetState(271)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(272)
					p.Match(AQLParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(273)
					p.pathCondition(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(278)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}



	errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IBooleanConditionContext is an interface to support dynamic dispatch.
type IBooleanConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IdentifiedPath() IIdentifiedPathContext
	EXISTS() antlr.TerminalNode
	AllComparisonOperand() []IComparisonOperandContext
	ComparisonOperand(i int) IComparisonOperandContext
	COMPARISON_OPERATOR() antlr.TerminalNode
	LIKE() antlr.TerminalNode
	StringOperand() IStringOperandContext
	IN() antlr.TerminalNode
	SYM_LEFT_PAREN() antlr.TerminalNode
	InOperand() IInOperandContext
	SYM_RIGHT_PAREN() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	CONTAINS() antlr.TerminalNode
	AT() antlr.TerminalNode
	ObjectPath() IObjectPathContext
	BooleanCondition() IBooleanConditionContext

	// IsBooleanConditionContext differentiates from other interfaces.
	IsBooleanConditionContext()
}

type BooleanConditionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanConditionContext() *BooleanConditionContext {
	var p = new(BooleanConditionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_booleanCondition
	return p
}

func InitEmptyBooleanConditionContext(p *BooleanConditionContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_booleanCondition
}

func (*BooleanConditionContext) IsBooleanConditionContext() {}

func NewBooleanConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanConditionContext {
	var p = new(BooleanConditionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_booleanCondition

	return p
}

func (s *BooleanConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanConditionContext) IdentifiedPath() IIdentifiedPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedPathContext)
}

func (s *BooleanConditionContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(AQLParserEXISTS, 0)
}

func (s *BooleanConditionContext) AllComparisonOperand() []IComparisonOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IComparisonOperandContext); ok {
			len++
		}
	}

	tst := make([]IComparisonOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IComparisonOperandContext); ok {
			tst[i] = t.(IComparisonOperandContext)
			i++
		}
	}

	return tst
}

func (s *BooleanConditionContext) ComparisonOperand(i int) IComparisonOperandContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonOperandContext)
}

func (s *BooleanConditionContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(AQLParserCOMPARISON_OPERATOR, 0)
}

func (s *BooleanConditionContext) LIKE() antlr.TerminalNode {
	return s.GetToken(AQLParserLIKE, 0)
}

func (s *BooleanConditionContext) StringOperand() IStringOperandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringOperandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringOperandContext)
}

func (s *BooleanConditionContext) IN() antlr.TerminalNode {
	return s.GetToken(AQLParserIN, 0)
}

func (s *BooleanConditionContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *BooleanConditionContext) InOperand() IInOperandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInOperandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInOperandContext)
}

func (s *BooleanConditionContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *BooleanConditionContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(AQLParserIDENTIFIER)
}

func (s *BooleanConditionContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, i)
}

func (s *BooleanConditionContext) CONTAINS() antlr.TerminalNode {
	return s.GetToken(AQLParserCONTAINS, 0)
}

func (s *BooleanConditionContext) AT() antlr.TerminalNode {
	return s.GetToken(AQLParserAT, 0)
}

func (s *BooleanConditionContext) ObjectPath() IObjectPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectPathContext)
}

func (s *BooleanConditionContext) BooleanCondition() IBooleanConditionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanConditionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanConditionContext)
}

func (s *BooleanConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BooleanConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterBooleanCondition(s)
	}
}

func (s *BooleanConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitBooleanCondition(s)
	}
}




func (p *AQLParser) BooleanCondition() (localctx IBooleanConditionContext) {
	localctx = NewBooleanConditionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, AQLParserRULE_booleanCondition)
	p.SetState(307)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 39, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(279)
			p.IdentifiedPath()
		}
		{
			p.SetState(280)
			p.Match(AQLParserEXISTS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(282)
			p.ComparisonOperand()
		}
		{
			p.SetState(283)
			p.Match(AQLParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(284)
			p.ComparisonOperand()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(286)
			p.ComparisonOperand()
		}
		{
			p.SetState(287)
			p.Match(AQLParserLIKE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(288)
			p.StringOperand()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(290)
			p.ComparisonOperand()
		}
		{
			p.SetState(291)
			p.Match(AQLParserIN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(292)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(293)
			p.InOperand()
		}
		{
			p.SetState(294)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(296)
			p.Match(AQLParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(297)
			p.Match(AQLParserCONTAINS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(298)
			p.Match(AQLParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(301)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(299)
				p.Match(AQLParserAT)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(300)
				p.ObjectPath()
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(303)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(304)
			p.BooleanCondition()
		}
		{
			p.SetState(305)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IPathConditionOperandContext is an interface to support dynamic dispatch.
type IPathConditionOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primitive() IPrimitiveContext
	ObjectPath() IObjectPathContext
	PARAMETER() antlr.TerminalNode

	// IsPathConditionOperandContext differentiates from other interfaces.
	IsPathConditionOperandContext()
}

type PathConditionOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathConditionOperandContext() *PathConditionOperandContext {
	var p = new(PathConditionOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathConditionOperand
	return p
}

func InitEmptyPathConditionOperandContext(p *PathConditionOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathConditionOperand
}

func (*PathConditionOperandContext) IsPathConditionOperandContext() {}

func NewPathConditionOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathConditionOperandContext {
	var p = new(PathConditionOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_pathConditionOperand

	return p
}

func (s *PathConditionOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *PathConditionOperandContext) Primitive() IPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimitiveContext)
}

func (s *PathConditionOperandContext) ObjectPath() IObjectPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectPathContext)
}

func (s *PathConditionOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *PathConditionOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathConditionOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathConditionOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterPathConditionOperand(s)
	}
}

func (s *PathConditionOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitPathConditionOperand(s)
	}
}




func (p *AQLParser) PathConditionOperand() (localctx IPathConditionOperandContext) {
	localctx = NewPathConditionOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, AQLParserRULE_pathConditionOperand)
	p.SetState(312)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserSTRING, AQLParserINTEGER, AQLParserFLOAT, AQLParserSCI_INTEGER, AQLParserSCI_FLOAT, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(309)
			p.Primitive()
		}


	case AQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(310)
			p.ObjectPath()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(311)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IComparisonOperandContext is an interface to support dynamic dispatch.
type IComparisonOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primitive() IPrimitiveContext
	IdentifiedPath() IIdentifiedPathContext
	PARAMETER() antlr.TerminalNode

	// IsComparisonOperandContext differentiates from other interfaces.
	IsComparisonOperandContext()
}

type ComparisonOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonOperandContext() *ComparisonOperandContext {
	var p = new(ComparisonOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_comparisonOperand
	return p
}

func InitEmptyComparisonOperandContext(p *ComparisonOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_comparisonOperand
}

func (*ComparisonOperandContext) IsComparisonOperandContext() {}

func NewComparisonOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonOperandContext {
	var p = new(ComparisonOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_comparisonOperand

	return p
}

func (s *ComparisonOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonOperandContext) Primitive() IPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimitiveContext)
}

func (s *ComparisonOperandContext) IdentifiedPath() IIdentifiedPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedPathContext)
}

func (s *ComparisonOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *ComparisonOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ComparisonOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterComparisonOperand(s)
	}
}

func (s *ComparisonOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitComparisonOperand(s)
	}
}




func (p *AQLParser) ComparisonOperand() (localctx IComparisonOperandContext) {
	localctx = NewComparisonOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, AQLParserRULE_comparisonOperand)
	p.SetState(317)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserSTRING, AQLParserINTEGER, AQLParserFLOAT, AQLParserSCI_INTEGER, AQLParserSCI_FLOAT, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(314)
			p.Primitive()
		}


	case AQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(315)
			p.IdentifiedPath()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(316)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IInOperandContext is an interface to support dynamic dispatch.
type IInOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SelectQuery() ISelectQueryContext
	AllInOperandValue() []IInOperandValueContext
	InOperandValue(i int) IInOperandValueContext
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode

	// IsInOperandContext differentiates from other interfaces.
	IsInOperandContext()
}

type InOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInOperandContext() *InOperandContext {
	var p = new(InOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_inOperand
	return p
}

func InitEmptyInOperandContext(p *InOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_inOperand
}

func (*InOperandContext) IsInOperandContext() {}

func NewInOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InOperandContext {
	var p = new(InOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_inOperand

	return p
}

func (s *InOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *InOperandContext) SelectQuery() ISelectQueryContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectQueryContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectQueryContext)
}

func (s *InOperandContext) AllInOperandValue() []IInOperandValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IInOperandValueContext); ok {
			len++
		}
	}

	tst := make([]IInOperandValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IInOperandValueContext); ok {
			tst[i] = t.(IInOperandValueContext)
			i++
		}
	}

	return tst
}

func (s *InOperandContext) InOperandValue(i int) IInOperandValueContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInOperandValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInOperandValueContext)
}

func (s *InOperandContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *InOperandContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
}

func (s *InOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *InOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterInOperand(s)
	}
}

func (s *InOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitInOperand(s)
	}
}




func (p *AQLParser) InOperand() (localctx IInOperandContext) {
	localctx = NewInOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, AQLParserRULE_inOperand)
	var _la int

	p.SetState(328)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserSELECT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(319)
			p.SelectQuery()
		}


	case AQLParserNULL, AQLParserBOOLEAN, AQLParserPARAMETER, AQLParserSTRING, AQLParserINTEGER, AQLParserFLOAT, AQLParserSCI_INTEGER, AQLParserSCI_FLOAT, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(320)
			p.InOperandValue()
		}
		p.SetState(325)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		for _la == AQLParserSYM_COMMA {
			{
				p.SetState(321)
				p.Match(AQLParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(322)
				p.InOperandValue()
			}


			p.SetState(327)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
		    	goto errorExit
		    }
			_la = p.GetTokenStream().LA(1)
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IInOperandValueContext is an interface to support dynamic dispatch.
type IInOperandValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primitive() IPrimitiveContext
	PARAMETER() antlr.TerminalNode

	// IsInOperandValueContext differentiates from other interfaces.
	IsInOperandValueContext()
}

type InOperandValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInOperandValueContext() *InOperandValueContext {
	var p = new(InOperandValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_inOperandValue
	return p
}

func InitEmptyInOperandValueContext(p *InOperandValueContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_inOperandValue
}

func (*InOperandValueContext) IsInOperandValueContext() {}

func NewInOperandValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InOperandValueContext {
	var p = new(InOperandValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_inOperandValue

	return p
}

func (s *InOperandValueContext) GetParser() antlr.Parser { return s.parser }

func (s *InOperandValueContext) Primitive() IPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimitiveContext)
}

func (s *InOperandValueContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *InOperandValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InOperandValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *InOperandValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterInOperandValue(s)
	}
}

func (s *InOperandValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitInOperandValue(s)
	}
}




func (p *AQLParser) InOperandValue() (localctx IInOperandValueContext) {
	localctx = NewInOperandValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, AQLParserRULE_inOperandValue)
	p.SetState(332)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserSTRING, AQLParserINTEGER, AQLParserFLOAT, AQLParserSCI_INTEGER, AQLParserSCI_FLOAT, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(330)
			p.Primitive()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(331)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IPrimitiveContext is an interface to support dynamic dispatch.
type IPrimitiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING() antlr.TerminalNode
	IntPrimitive() IIntPrimitiveContext
	FloatPrimitive() IFloatPrimitiveContext
	BOOLEAN() antlr.TerminalNode
	NULL() antlr.TerminalNode

	// IsPrimitiveContext differentiates from other interfaces.
	IsPrimitiveContext()
}

type PrimitiveContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimitiveContext() *PrimitiveContext {
	var p = new(PrimitiveContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_primitive
	return p
}

func InitEmptyPrimitiveContext(p *PrimitiveContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_primitive
}

func (*PrimitiveContext) IsPrimitiveContext() {}

func NewPrimitiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimitiveContext {
	var p = new(PrimitiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_primitive

	return p
}

func (s *PrimitiveContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimitiveContext) STRING() antlr.TerminalNode {
	return s.GetToken(AQLParserSTRING, 0)
}

func (s *PrimitiveContext) IntPrimitive() IIntPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntPrimitiveContext)
}

func (s *PrimitiveContext) FloatPrimitive() IFloatPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFloatPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFloatPrimitiveContext)
}

func (s *PrimitiveContext) BOOLEAN() antlr.TerminalNode {
	return s.GetToken(AQLParserBOOLEAN, 0)
}

func (s *PrimitiveContext) NULL() antlr.TerminalNode {
	return s.GetToken(AQLParserNULL, 0)
}

func (s *PrimitiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimitiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PrimitiveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterPrimitive(s)
	}
}

func (s *PrimitiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitPrimitive(s)
	}
}




func (p *AQLParser) Primitive() (localctx IPrimitiveContext) {
	localctx = NewPrimitiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, AQLParserRULE_primitive)
	p.SetState(339)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 45, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(334)
			p.Match(AQLParserSTRING)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(335)
			p.IntPrimitive()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(336)
			p.FloatPrimitive()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(337)
			p.Match(AQLParserBOOLEAN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(338)
			p.Match(AQLParserNULL)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IIntPrimitiveContext is an interface to support dynamic dispatch.
type IIntPrimitiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER() antlr.TerminalNode
	SYM_MINUS() antlr.TerminalNode
	SCI_INTEGER() antlr.TerminalNode

	// IsIntPrimitiveContext differentiates from other interfaces.
	IsIntPrimitiveContext()
}

type IntPrimitiveContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntPrimitiveContext() *IntPrimitiveContext {
	var p = new(IntPrimitiveContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_intPrimitive
	return p
}

func InitEmptyIntPrimitiveContext(p *IntPrimitiveContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_intPrimitive
}

func (*IntPrimitiveContext) IsIntPrimitiveContext() {}

func NewIntPrimitiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntPrimitiveContext {
	var p = new(IntPrimitiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_intPrimitive

	return p
}

func (s *IntPrimitiveContext) GetParser() antlr.Parser { return s.parser }

func (s *IntPrimitiveContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(AQLParserINTEGER, 0)
}

func (s *IntPrimitiveContext) SYM_MINUS() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_MINUS, 0)
}

func (s *IntPrimitiveContext) SCI_INTEGER() antlr.TerminalNode {
	return s.GetToken(AQLParserSCI_INTEGER, 0)
}

func (s *IntPrimitiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntPrimitiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IntPrimitiveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterIntPrimitive(s)
	}
}

func (s *IntPrimitiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitIntPrimitive(s)
	}
}




func (p *AQLParser) IntPrimitive() (localctx IIntPrimitiveContext) {
	localctx = NewIntPrimitiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, AQLParserRULE_intPrimitive)
	var _la int

	p.SetState(349)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(342)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserSYM_MINUS {
			{
				p.SetState(341)
				p.Match(AQLParserSYM_MINUS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(344)
			p.Match(AQLParserINTEGER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(346)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserSYM_MINUS {
			{
				p.SetState(345)
				p.Match(AQLParserSYM_MINUS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(348)
			p.Match(AQLParserSCI_INTEGER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IFloatPrimitiveContext is an interface to support dynamic dispatch.
type IFloatPrimitiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FLOAT() antlr.TerminalNode
	SYM_MINUS() antlr.TerminalNode
	SCI_FLOAT() antlr.TerminalNode

	// IsFloatPrimitiveContext differentiates from other interfaces.
	IsFloatPrimitiveContext()
}

type FloatPrimitiveContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFloatPrimitiveContext() *FloatPrimitiveContext {
	var p = new(FloatPrimitiveContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_floatPrimitive
	return p
}

func InitEmptyFloatPrimitiveContext(p *FloatPrimitiveContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_floatPrimitive
}

func (*FloatPrimitiveContext) IsFloatPrimitiveContext() {}

func NewFloatPrimitiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FloatPrimitiveContext {
	var p = new(FloatPrimitiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_floatPrimitive

	return p
}

func (s *FloatPrimitiveContext) GetParser() antlr.Parser { return s.parser }

func (s *FloatPrimitiveContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(AQLParserFLOAT, 0)
}

func (s *FloatPrimitiveContext) SYM_MINUS() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_MINUS, 0)
}

func (s *FloatPrimitiveContext) SCI_FLOAT() antlr.TerminalNode {
	return s.GetToken(AQLParserSCI_FLOAT, 0)
}

func (s *FloatPrimitiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FloatPrimitiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FloatPrimitiveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterFloatPrimitive(s)
	}
}

func (s *FloatPrimitiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitFloatPrimitive(s)
	}
}




func (p *AQLParser) FloatPrimitive() (localctx IFloatPrimitiveContext) {
	localctx = NewFloatPrimitiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, AQLParserRULE_floatPrimitive)
	var _la int

	p.SetState(359)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(352)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserSYM_MINUS {
			{
				p.SetState(351)
				p.Match(AQLParserSYM_MINUS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(354)
			p.Match(AQLParserFLOAT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(356)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserSYM_MINUS {
			{
				p.SetState(355)
				p.Match(AQLParserSYM_MINUS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

		}
		{
			p.SetState(358)
			p.Match(AQLParserSCI_FLOAT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IStringOperandContext is an interface to support dynamic dispatch.
type IStringOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING() antlr.TerminalNode
	PARAMETER() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsStringOperandContext differentiates from other interfaces.
	IsStringOperandContext()
}

type StringOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringOperandContext() *StringOperandContext {
	var p = new(StringOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_stringOperand
	return p
}

func InitEmptyStringOperandContext(p *StringOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_stringOperand
}

func (*StringOperandContext) IsStringOperandContext() {}

func NewStringOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringOperandContext {
	var p = new(StringOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_stringOperand

	return p
}

func (s *StringOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *StringOperandContext) STRING() antlr.TerminalNode {
	return s.GetToken(AQLParserSTRING, 0)
}

func (s *StringOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *StringOperandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
}

func (s *StringOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StringOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterStringOperand(s)
	}
}

func (s *StringOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitStringOperand(s)
	}
}




func (p *AQLParser) StringOperand() (localctx IStringOperandContext) {
	localctx = NewStringOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, AQLParserRULE_stringOperand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(361)
		_la = p.GetTokenStream().LA(1)

		if !(((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 58546795155816448) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IIntOperandContext is an interface to support dynamic dispatch.
type IIntOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IntPrimitive() IIntPrimitiveContext
	PARAMETER() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsIntOperandContext differentiates from other interfaces.
	IsIntOperandContext()
}

type IntOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntOperandContext() *IntOperandContext {
	var p = new(IntOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_intOperand
	return p
}

func InitEmptyIntOperandContext(p *IntOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_intOperand
}

func (*IntOperandContext) IsIntOperandContext() {}

func NewIntOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntOperandContext {
	var p = new(IntOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_intOperand

	return p
}

func (s *IntOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *IntOperandContext) IntPrimitive() IIntPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntPrimitiveContext)
}

func (s *IntOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *IntOperandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
}

func (s *IntOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IntOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterIntOperand(s)
	}
}

func (s *IntOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitIntOperand(s)
	}
}




func (p *AQLParser) IntOperand() (localctx IIntOperandContext) {
	localctx = NewIntOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, AQLParserRULE_intOperand)
	p.SetState(366)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserINTEGER, AQLParserSCI_INTEGER, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(363)
			p.IntPrimitive()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(364)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(365)
			p.Match(AQLParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// INumbericOperandContext is an interface to support dynamic dispatch.
type INumbericOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IntPrimitive() IIntPrimitiveContext
	FloatPrimitive() IFloatPrimitiveContext
	PARAMETER() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsNumbericOperandContext differentiates from other interfaces.
	IsNumbericOperandContext()
}

type NumbericOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumbericOperandContext() *NumbericOperandContext {
	var p = new(NumbericOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_numbericOperand
	return p
}

func InitEmptyNumbericOperandContext(p *NumbericOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_numbericOperand
}

func (*NumbericOperandContext) IsNumbericOperandContext() {}

func NewNumbericOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumbericOperandContext {
	var p = new(NumbericOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_numbericOperand

	return p
}

func (s *NumbericOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *NumbericOperandContext) IntPrimitive() IIntPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntPrimitiveContext)
}

func (s *NumbericOperandContext) FloatPrimitive() IFloatPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFloatPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFloatPrimitiveContext)
}

func (s *NumbericOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *NumbericOperandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
}

func (s *NumbericOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumbericOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NumbericOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterNumbericOperand(s)
	}
}

func (s *NumbericOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitNumbericOperand(s)
	}
}




func (p *AQLParser) NumbericOperand() (localctx INumbericOperandContext) {
	localctx = NewNumbericOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, AQLParserRULE_numbericOperand)
	p.SetState(372)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 53, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(368)
			p.IntPrimitive()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(369)
			p.FloatPrimitive()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(370)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(371)
			p.Match(AQLParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LENGTH() antlr.TerminalNode
	SYM_LEFT_PAREN() antlr.TerminalNode
	AllStringOperand() []IStringOperandContext
	StringOperand(i int) IStringOperandContext
	SYM_RIGHT_PAREN() antlr.TerminalNode
	POSITION() antlr.TerminalNode
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode
	AllIntOperand() []IIntOperandContext
	IntOperand(i int) IIntOperandContext
	SUBSTRING() antlr.TerminalNode
	CONCAT() antlr.TerminalNode
	CONCAT_WS() antlr.TerminalNode
	ABS() antlr.TerminalNode
	AllNumbericOperand() []INumbericOperandContext
	NumbericOperand(i int) INumbericOperandContext
	MOD() antlr.TerminalNode
	CEIL() antlr.TerminalNode
	FLOOR() antlr.TerminalNode
	ROUND() antlr.TerminalNode
	CURRENT_DATE() antlr.TerminalNode
	CURRENT_TIME() antlr.TerminalNode
	CURRENT_DATE_TIME() antlr.TerminalNode

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) LENGTH() antlr.TerminalNode {
	return s.GetToken(AQLParserLENGTH, 0)
}

func (s *FunctionCallContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *FunctionCallContext) AllStringOperand() []IStringOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStringOperandContext); ok {
			len++
		}
	}

	tst := make([]IStringOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStringOperandContext); ok {
			tst[i] = t.(IStringOperandContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) StringOperand(i int) IStringOperandContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringOperandContext)
}

func (s *FunctionCallContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *FunctionCallContext) POSITION() antlr.TerminalNode {
	return s.GetToken(AQLParserPOSITION, 0)
}

func (s *FunctionCallContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *FunctionCallContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
}

func (s *FunctionCallContext) AllIntOperand() []IIntOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIntOperandContext); ok {
			len++
		}
	}

	tst := make([]IIntOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIntOperandContext); ok {
			tst[i] = t.(IIntOperandContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) IntOperand(i int) IIntOperandContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntOperandContext)
}

func (s *FunctionCallContext) SUBSTRING() antlr.TerminalNode {
	return s.GetToken(AQLParserSUBSTRING, 0)
}

func (s *FunctionCallContext) CONCAT() antlr.TerminalNode {
	return s.GetToken(AQLParserCONCAT, 0)
}

func (s *FunctionCallContext) CONCAT_WS() antlr.TerminalNode {
	return s.GetToken(AQLParserCONCAT_WS, 0)
}

func (s *FunctionCallContext) ABS() antlr.TerminalNode {
	return s.GetToken(AQLParserABS, 0)
}

func (s *FunctionCallContext) AllNumbericOperand() []INumbericOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INumbericOperandContext); ok {
			len++
		}
	}

	tst := make([]INumbericOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INumbericOperandContext); ok {
			tst[i] = t.(INumbericOperandContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) NumbericOperand(i int) INumbericOperandContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumbericOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumbericOperandContext)
}

func (s *FunctionCallContext) MOD() antlr.TerminalNode {
	return s.GetToken(AQLParserMOD, 0)
}

func (s *FunctionCallContext) CEIL() antlr.TerminalNode {
	return s.GetToken(AQLParserCEIL, 0)
}

func (s *FunctionCallContext) FLOOR() antlr.TerminalNode {
	return s.GetToken(AQLParserFLOOR, 0)
}

func (s *FunctionCallContext) ROUND() antlr.TerminalNode {
	return s.GetToken(AQLParserROUND, 0)
}

func (s *FunctionCallContext) CURRENT_DATE() antlr.TerminalNode {
	return s.GetToken(AQLParserCURRENT_DATE, 0)
}

func (s *FunctionCallContext) CURRENT_TIME() antlr.TerminalNode {
	return s.GetToken(AQLParserCURRENT_TIME, 0)
}

func (s *FunctionCallContext) CURRENT_DATE_TIME() antlr.TerminalNode {
	return s.GetToken(AQLParserCURRENT_DATE_TIME, 0)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}




func (p *AQLParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, AQLParserRULE_functionCall)
	var _la int

	p.SetState(453)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserLENGTH:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(374)
			p.Match(AQLParserLENGTH)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(375)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(376)
			p.StringOperand()
		}
		{
			p.SetState(377)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserPOSITION:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(379)
			p.Match(AQLParserPOSITION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(380)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(381)
			p.StringOperand()
		}
		{
			p.SetState(382)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(383)
			p.IntOperand()
		}
		{
			p.SetState(384)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserSUBSTRING:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(386)
			p.Match(AQLParserSUBSTRING)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(387)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(388)
			p.StringOperand()
		}
		{
			p.SetState(389)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(390)
			p.IntOperand()
		}
		{
			p.SetState(391)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(392)
			p.IntOperand()
		}
		{
			p.SetState(393)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserCONCAT:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(395)
			p.Match(AQLParserCONCAT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(396)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(397)
			p.StringOperand()
		}
		p.SetState(402)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		for _la == AQLParserSYM_COMMA {
			{
				p.SetState(398)
				p.Match(AQLParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(399)
				p.StringOperand()
			}


			p.SetState(404)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
		    	goto errorExit
		    }
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(405)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserCONCAT_WS:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(407)
			p.Match(AQLParserCONCAT_WS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(408)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(409)
			p.StringOperand()
		}
		{
			p.SetState(410)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(411)
			p.StringOperand()
		}
		p.SetState(416)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		for _la == AQLParserSYM_COMMA {
			{
				p.SetState(412)
				p.Match(AQLParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(413)
				p.StringOperand()
			}


			p.SetState(418)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
		    	goto errorExit
		    }
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(419)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserABS:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(421)
			p.Match(AQLParserABS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(422)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(423)
			p.NumbericOperand()
		}
		{
			p.SetState(424)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserMOD:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(426)
			p.Match(AQLParserMOD)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(427)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(428)
			p.NumbericOperand()
		}
		{
			p.SetState(429)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(430)
			p.NumbericOperand()
		}
		{
			p.SetState(431)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserCEIL:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(433)
			p.Match(AQLParserCEIL)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(434)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(435)
			p.NumbericOperand()
		}
		{
			p.SetState(436)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserFLOOR:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(438)
			p.Match(AQLParserFLOOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(439)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(440)
			p.NumbericOperand()
		}
		{
			p.SetState(441)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserROUND:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(443)
			p.Match(AQLParserROUND)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(444)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(445)
			p.NumbericOperand()
		}
		{
			p.SetState(446)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(447)
			p.IntOperand()
		}
		{
			p.SetState(448)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserCURRENT_DATE:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(450)
			p.Match(AQLParserCURRENT_DATE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserCURRENT_TIME:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(451)
			p.Match(AQLParserCURRENT_TIME)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserCURRENT_DATE_TIME:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(452)
			p.Match(AQLParserCURRENT_DATE_TIME)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IAggregateFunctionCallContext is an interface to support dynamic dispatch.
type IAggregateFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name token.
	GetName() antlr.Token 


	// SetName sets the name token.
	SetName(antlr.Token) 


	// Getter signatures
	SYM_LEFT_PAREN() antlr.TerminalNode
	SYM_RIGHT_PAREN() antlr.TerminalNode
	COUNT() antlr.TerminalNode
	SYM_ASTERISK() antlr.TerminalNode
	IdentifiedPath() IIdentifiedPathContext
	DISTINCT() antlr.TerminalNode
	MIN() antlr.TerminalNode
	MAX() antlr.TerminalNode
	SUM() antlr.TerminalNode
	AVG() antlr.TerminalNode

	// IsAggregateFunctionCallContext differentiates from other interfaces.
	IsAggregateFunctionCallContext()
}

type AggregateFunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	name antlr.Token
}

func NewEmptyAggregateFunctionCallContext() *AggregateFunctionCallContext {
	var p = new(AggregateFunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_aggregateFunctionCall
	return p
}

func InitEmptyAggregateFunctionCallContext(p *AggregateFunctionCallContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_aggregateFunctionCall
}

func (*AggregateFunctionCallContext) IsAggregateFunctionCallContext() {}

func NewAggregateFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregateFunctionCallContext {
	var p = new(AggregateFunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_aggregateFunctionCall

	return p
}

func (s *AggregateFunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregateFunctionCallContext) GetName() antlr.Token { return s.name }


func (s *AggregateFunctionCallContext) SetName(v antlr.Token) { s.name = v }


func (s *AggregateFunctionCallContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *AggregateFunctionCallContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *AggregateFunctionCallContext) COUNT() antlr.TerminalNode {
	return s.GetToken(AQLParserCOUNT, 0)
}

func (s *AggregateFunctionCallContext) SYM_ASTERISK() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_ASTERISK, 0)
}

func (s *AggregateFunctionCallContext) IdentifiedPath() IIdentifiedPathContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedPathContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedPathContext)
}

func (s *AggregateFunctionCallContext) DISTINCT() antlr.TerminalNode {
	return s.GetToken(AQLParserDISTINCT, 0)
}

func (s *AggregateFunctionCallContext) MIN() antlr.TerminalNode {
	return s.GetToken(AQLParserMIN, 0)
}

func (s *AggregateFunctionCallContext) MAX() antlr.TerminalNode {
	return s.GetToken(AQLParserMAX, 0)
}

func (s *AggregateFunctionCallContext) SUM() antlr.TerminalNode {
	return s.GetToken(AQLParserSUM, 0)
}

func (s *AggregateFunctionCallContext) AVG() antlr.TerminalNode {
	return s.GetToken(AQLParserAVG, 0)
}

func (s *AggregateFunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregateFunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *AggregateFunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterAggregateFunctionCall(s)
	}
}

func (s *AggregateFunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitAggregateFunctionCall(s)
	}
}




func (p *AQLParser) AggregateFunctionCall() (localctx IAggregateFunctionCallContext) {
	localctx = NewAggregateFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, AQLParserRULE_aggregateFunctionCall)
	var _la int

	p.SetState(470)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserCOUNT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(455)

			var _m = p.Match(AQLParserCOUNT)

			localctx.(*AggregateFunctionCallContext).name = _m
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(456)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(462)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case AQLParserDISTINCT, AQLParserSYM_ASTERISK:
			p.SetState(458)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			if _la == AQLParserDISTINCT {
				{
					p.SetState(457)
					p.Match(AQLParserDISTINCT)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}

			}
			{
				p.SetState(460)
				p.Match(AQLParserSYM_ASTERISK)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}


		case AQLParserIDENTIFIER:
			{
				p.SetState(461)
				p.IdentifiedPath()
			}



		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(464)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserMIN, AQLParserMAX, AQLParserSUM, AQLParserAVG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(465)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*AggregateFunctionCallContext).name = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 4222124650659840) != 0)) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*AggregateFunctionCallContext).name = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(466)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(467)
			p.IdentifiedPath()
		}
		{
			p.SetState(468)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ILimitOperandContext is an interface to support dynamic dispatch.
type ILimitOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER() antlr.TerminalNode
	PARAMETER() antlr.TerminalNode

	// IsLimitOperandContext differentiates from other interfaces.
	IsLimitOperandContext()
}

type LimitOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitOperandContext() *LimitOperandContext {
	var p = new(LimitOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_limitOperand
	return p
}

func InitEmptyLimitOperandContext(p *LimitOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_limitOperand
}

func (*LimitOperandContext) IsLimitOperandContext() {}

func NewLimitOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitOperandContext {
	var p = new(LimitOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_limitOperand

	return p
}

func (s *LimitOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitOperandContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(AQLParserINTEGER, 0)
}

func (s *LimitOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *LimitOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LimitOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterLimitOperand(s)
	}
}

func (s *LimitOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitLimitOperand(s)
	}
}




func (p *AQLParser) LimitOperand() (localctx ILimitOperandContext) {
	localctx = NewLimitOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, AQLParserRULE_limitOperand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(472)
		_la = p.GetTokenStream().LA(1)

		if !(_la == AQLParserPARAMETER || _la == AQLParserINTEGER) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


func (p *AQLParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 12:
			var t *WhereExprContext = nil
			if localctx != nil { t = localctx.(*WhereExprContext) }
			return p.WhereExpr_Sempred(t, predIndex)

	case 18:
			var t *PathConditionContext = nil
			if localctx != nil { t = localctx.(*PathConditionContext) }
			return p.PathCondition_Sempred(t, predIndex)


	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *AQLParser) WhereExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
			return p.Precpred(p.GetParserRuleContext(), 3)

	case 1:
			return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *AQLParser) PathCondition_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
			return p.Precpred(p.GetParserRuleContext(), 3)

	case 3:
			return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

