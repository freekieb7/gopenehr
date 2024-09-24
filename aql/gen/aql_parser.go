// Code generated from AqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package gen // AqlParser
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


type AqlParser struct {
	*antlr.BaseParser
}

var AqlParserParserStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  LiteralNames           []string
  SymbolicNames          []string
  RuleNames              []string
  PredictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func aqlparserParserInit() {
  staticData := &AqlParserParserStaticData
  staticData.LiteralNames = []string{
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "';'", "'<'", "'>'", "'<='", "'>='", "'!='", "'='", 
    "'('", "')'", "','", "'/'", "'*'", "'+'", "'-'", "'['", "']'", "'{'", 
    "'}'", "'--'",
  }
  staticData.SymbolicNames = []string{
    "", "WS", "UNICODE_BOM", "COMMENT", "SELECT", "AS", "FROM", "WHERE", 
    "ORDER", "BY", "DESC", "DESCENDING", "ASC", "ASCENDING", "LIMIT", "OFFSET", 
    "DISTINCT", "VERSION", "LATEST_VERSION", "ALL_VERSIONS", "NULL", "BOOLEAN", 
    "TOP", "FORWARD", "BACKWARD", "CONTAINS", "AND", "OR", "NOT", "EXISTS", 
    "COMPARISON_OPERATOR", "LIKE", "MATCHES", "STRING_FUNCTION_ID", "NUMERIC_FUNCTION_ID", 
    "DATE_TIME_FUNCTION_ID", "LENGTH", "POSITION", "SUBSTRING", "CONCAT", 
    "CONCAT_WS", "ABS", "MOD", "CEIL", "FLOOR", "ROUND", "CURRENT_DATE", 
    "CURRENT_TIME", "CURRENT_DATE_TIME", "NOW", "CURRENT_TIMEZONE", "COUNT", 
    "MIN", "MAX", "SUM", "AVG", "TERMINOLOGY", "PARAMETER", "ID_CODE", "AT_CODE", 
    "CONTAINED_REGEX", "ARCHETYPE_HRID", "IDENTIFIER", "TERM_CODE", "URI", 
    "INTEGER", "REAL", "SCI_INTEGER", "SCI_REAL", "DATE", "TIME", "DATETIME", 
    "STRING", "SYM_SEMICOLON", "SYM_LT", "SYM_GT", "SYM_LE", "SYM_GE", "SYM_NE", 
    "SYM_EQ", "SYM_LEFT_PAREN", "SYM_RIGHT_PAREN", "SYM_COMMA", "SYM_SLASH", 
    "SYM_ASTERISK", "SYM_PLUS", "SYM_MINUS", "SYM_LEFT_BRACKET", "SYM_RIGHT_BRACKET", 
    "SYM_LEFT_CURLY", "SYM_RIGHT_CURLY", "SYM_DOUBLE_DASH",
  }
  staticData.RuleNames = []string{
    "selectQuery", "selectClause", "fromClause", "whereClause", "orderByClause", 
    "limitClause", "selectExpr", "fromExpr", "whereExpr", "orderByExpr", 
    "columnExpr", "containsExpr", "identifiedExpr", "classExprOperand", 
    "terminal", "identifiedPath", "pathPredicate", "standardPredicate", 
    "archetypePredicate", "nodePredicate", "versionPredicate", "pathPredicateOperand", 
    "objectPath", "pathPart", "likeOperand", "matchesOperand", "valueListItem", 
    "primitive", "numericPrimitive", "functionCall", "aggregateFunctionCall", 
    "terminologyFunction", "top",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 91, 400, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 
	4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 
	10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15, 
	2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2, 
	21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26, 
	7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7, 
	31, 2, 32, 7, 32, 1, 0, 1, 0, 1, 0, 3, 0, 70, 8, 0, 1, 0, 3, 0, 73, 8, 
	0, 1, 0, 3, 0, 76, 8, 0, 1, 0, 3, 0, 79, 8, 0, 1, 0, 1, 0, 1, 1, 1, 1, 
	3, 1, 85, 8, 1, 1, 1, 3, 1, 88, 8, 1, 1, 1, 1, 1, 1, 1, 5, 1, 93, 8, 1, 
	10, 1, 12, 1, 96, 9, 1, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 
	1, 4, 1, 4, 1, 4, 5, 4, 109, 8, 4, 10, 4, 12, 4, 112, 9, 4, 1, 5, 1, 5, 
	1, 5, 1, 5, 3, 5, 118, 8, 5, 1, 6, 1, 6, 1, 6, 3, 6, 123, 8, 6, 1, 7, 1, 
	7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 135, 8, 8, 1, 
	8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 5, 8, 143, 8, 8, 10, 8, 12, 8, 146, 9, 
	8, 1, 9, 1, 9, 3, 9, 150, 8, 9, 1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 156, 
	8, 10, 1, 11, 1, 11, 1, 11, 3, 11, 161, 8, 11, 1, 11, 1, 11, 3, 11, 165, 
	8, 11, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 171, 8, 11, 1, 11, 1, 11, 1, 
	11, 1, 11, 1, 11, 1, 11, 5, 11, 179, 8, 11, 10, 11, 12, 11, 182, 9, 11, 
	1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 
	12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 
	1, 12, 3, 12, 206, 8, 12, 1, 13, 1, 13, 3, 13, 210, 8, 13, 1, 13, 3, 13, 
	213, 8, 13, 1, 13, 1, 13, 3, 13, 217, 8, 13, 1, 13, 1, 13, 1, 13, 1, 13, 
	3, 13, 223, 8, 13, 3, 13, 225, 8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 3, 14, 
	231, 8, 14, 1, 15, 1, 15, 3, 15, 235, 8, 15, 1, 15, 1, 15, 3, 15, 239, 
	8, 15, 1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 245, 8, 16, 1, 16, 1, 16, 1, 
	17, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 
	259, 8, 19, 1, 19, 1, 19, 1, 19, 3, 19, 264, 8, 19, 1, 19, 1, 19, 1, 19, 
	1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 275, 8, 19, 1, 19, 1, 
	19, 1, 19, 1, 19, 1, 19, 1, 19, 5, 19, 283, 8, 19, 10, 19, 12, 19, 286, 
	9, 19, 1, 20, 1, 20, 1, 20, 3, 20, 291, 8, 20, 1, 21, 1, 21, 1, 21, 1, 
	21, 1, 21, 3, 21, 298, 8, 21, 1, 22, 1, 22, 1, 22, 5, 22, 303, 8, 22, 10, 
	22, 12, 22, 306, 9, 22, 1, 23, 1, 23, 3, 23, 310, 8, 23, 1, 24, 1, 24, 
	1, 25, 1, 25, 1, 25, 1, 25, 5, 25, 318, 8, 25, 10, 25, 12, 25, 321, 9, 
	25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 329, 8, 25, 1, 26, 
	1, 26, 1, 26, 3, 26, 334, 8, 26, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 
	27, 1, 27, 3, 27, 343, 8, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 
	3, 28, 351, 8, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 5, 29, 359, 
	8, 29, 10, 29, 12, 29, 362, 9, 29, 3, 29, 364, 8, 29, 1, 29, 3, 29, 367, 
	8, 29, 1, 30, 1, 30, 1, 30, 3, 30, 372, 8, 30, 1, 30, 1, 30, 3, 30, 376, 
	8, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 3, 30, 384, 8, 30, 1, 
	31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 32, 1, 32, 
	1, 32, 3, 32, 398, 8, 32, 1, 32, 0, 3, 16, 22, 38, 33, 0, 2, 4, 6, 8, 10, 
	12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 
	48, 50, 52, 54, 56, 58, 60, 62, 64, 0, 8, 1, 0, 10, 13, 2, 0, 57, 57, 61, 
	61, 1, 0, 58, 59, 3, 0, 57, 59, 63, 63, 72, 72, 2, 0, 57, 57, 72, 72, 2, 
	0, 33, 35, 62, 62, 1, 0, 52, 55, 1, 0, 23, 24, 444, 0, 66, 1, 0, 0, 0, 
	2, 82, 1, 0, 0, 0, 4, 97, 1, 0, 0, 0, 6, 100, 1, 0, 0, 0, 8, 103, 1, 0, 
	0, 0, 10, 113, 1, 0, 0, 0, 12, 119, 1, 0, 0, 0, 14, 124, 1, 0, 0, 0, 16, 
	134, 1, 0, 0, 0, 18, 147, 1, 0, 0, 0, 20, 155, 1, 0, 0, 0, 22, 170, 1, 
	0, 0, 0, 24, 205, 1, 0, 0, 0, 26, 224, 1, 0, 0, 0, 28, 230, 1, 0, 0, 0, 
	30, 232, 1, 0, 0, 0, 32, 240, 1, 0, 0, 0, 34, 248, 1, 0, 0, 0, 36, 252, 
	1, 0, 0, 0, 38, 274, 1, 0, 0, 0, 40, 290, 1, 0, 0, 0, 42, 297, 1, 0, 0, 
	0, 44, 299, 1, 0, 0, 0, 46, 307, 1, 0, 0, 0, 48, 311, 1, 0, 0, 0, 50, 328, 
	1, 0, 0, 0, 52, 333, 1, 0, 0, 0, 54, 342, 1, 0, 0, 0, 56, 350, 1, 0, 0, 
	0, 58, 366, 1, 0, 0, 0, 60, 383, 1, 0, 0, 0, 62, 385, 1, 0, 0, 0, 64, 394, 
	1, 0, 0, 0, 66, 67, 3, 2, 1, 0, 67, 69, 3, 4, 2, 0, 68, 70, 3, 6, 3, 0, 
	69, 68, 1, 0, 0, 0, 69, 70, 1, 0, 0, 0, 70, 72, 1, 0, 0, 0, 71, 73, 3, 
	8, 4, 0, 72, 71, 1, 0, 0, 0, 72, 73, 1, 0, 0, 0, 73, 75, 1, 0, 0, 0, 74, 
	76, 3, 10, 5, 0, 75, 74, 1, 0, 0, 0, 75, 76, 1, 0, 0, 0, 76, 78, 1, 0, 
	0, 0, 77, 79, 5, 91, 0, 0, 78, 77, 1, 0, 0, 0, 78, 79, 1, 0, 0, 0, 79, 
	80, 1, 0, 0, 0, 80, 81, 5, 0, 0, 1, 81, 1, 1, 0, 0, 0, 82, 84, 5, 4, 0, 
	0, 83, 85, 5, 16, 0, 0, 84, 83, 1, 0, 0, 0, 84, 85, 1, 0, 0, 0, 85, 87, 
	1, 0, 0, 0, 86, 88, 3, 64, 32, 0, 87, 86, 1, 0, 0, 0, 87, 88, 1, 0, 0, 
	0, 88, 89, 1, 0, 0, 0, 89, 94, 3, 12, 6, 0, 90, 91, 5, 82, 0, 0, 91, 93, 
	3, 12, 6, 0, 92, 90, 1, 0, 0, 0, 93, 96, 1, 0, 0, 0, 94, 92, 1, 0, 0, 0, 
	94, 95, 1, 0, 0, 0, 95, 3, 1, 0, 0, 0, 96, 94, 1, 0, 0, 0, 97, 98, 5, 6, 
	0, 0, 98, 99, 3, 14, 7, 0, 99, 5, 1, 0, 0, 0, 100, 101, 5, 7, 0, 0, 101, 
	102, 3, 16, 8, 0, 102, 7, 1, 0, 0, 0, 103, 104, 5, 8, 0, 0, 104, 105, 5, 
	9, 0, 0, 105, 110, 3, 18, 9, 0, 106, 107, 5, 82, 0, 0, 107, 109, 3, 18, 
	9, 0, 108, 106, 1, 0, 0, 0, 109, 112, 1, 0, 0, 0, 110, 108, 1, 0, 0, 0, 
	110, 111, 1, 0, 0, 0, 111, 9, 1, 0, 0, 0, 112, 110, 1, 0, 0, 0, 113, 114, 
	5, 14, 0, 0, 114, 117, 5, 65, 0, 0, 115, 116, 5, 15, 0, 0, 116, 118, 5, 
	65, 0, 0, 117, 115, 1, 0, 0, 0, 117, 118, 1, 0, 0, 0, 118, 11, 1, 0, 0, 
	0, 119, 122, 3, 20, 10, 0, 120, 121, 5, 5, 0, 0, 121, 123, 5, 62, 0, 0, 
	122, 120, 1, 0, 0, 0, 122, 123, 1, 0, 0, 0, 123, 13, 1, 0, 0, 0, 124, 125, 
	3, 22, 11, 0, 125, 15, 1, 0, 0, 0, 126, 127, 6, 8, -1, 0, 127, 135, 3, 
	24, 12, 0, 128, 129, 5, 28, 0, 0, 129, 135, 3, 16, 8, 4, 130, 131, 5, 80, 
	0, 0, 131, 132, 3, 16, 8, 0, 132, 133, 5, 81, 0, 0, 133, 135, 1, 0, 0, 
	0, 134, 126, 1, 0, 0, 0, 134, 128, 1, 0, 0, 0, 134, 130, 1, 0, 0, 0, 135, 
	144, 1, 0, 0, 0, 136, 137, 10, 3, 0, 0, 137, 138, 5, 26, 0, 0, 138, 143, 
	3, 16, 8, 4, 139, 140, 10, 2, 0, 0, 140, 141, 5, 27, 0, 0, 141, 143, 3, 
	16, 8, 3, 142, 136, 1, 0, 0, 0, 142, 139, 1, 0, 0, 0, 143, 146, 1, 0, 0, 
	0, 144, 142, 1, 0, 0, 0, 144, 145, 1, 0, 0, 0, 145, 17, 1, 0, 0, 0, 146, 
	144, 1, 0, 0, 0, 147, 149, 3, 30, 15, 0, 148, 150, 7, 0, 0, 0, 149, 148, 
	1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 19, 1, 0, 0, 0, 151, 156, 3, 30, 
	15, 0, 152, 156, 3, 54, 27, 0, 153, 156, 3, 60, 30, 0, 154, 156, 3, 58, 
	29, 0, 155, 151, 1, 0, 0, 0, 155, 152, 1, 0, 0, 0, 155, 153, 1, 0, 0, 0, 
	155, 154, 1, 0, 0, 0, 156, 21, 1, 0, 0, 0, 157, 158, 6, 11, -1, 0, 158, 
	164, 3, 26, 13, 0, 159, 161, 5, 28, 0, 0, 160, 159, 1, 0, 0, 0, 160, 161, 
	1, 0, 0, 0, 161, 162, 1, 0, 0, 0, 162, 163, 5, 25, 0, 0, 163, 165, 3, 22, 
	11, 0, 164, 160, 1, 0, 0, 0, 164, 165, 1, 0, 0, 0, 165, 171, 1, 0, 0, 0, 
	166, 167, 5, 80, 0, 0, 167, 168, 3, 22, 11, 0, 168, 169, 5, 81, 0, 0, 169, 
	171, 1, 0, 0, 0, 170, 157, 1, 0, 0, 0, 170, 166, 1, 0, 0, 0, 171, 180, 
	1, 0, 0, 0, 172, 173, 10, 3, 0, 0, 173, 174, 5, 26, 0, 0, 174, 179, 3, 
	22, 11, 4, 175, 176, 10, 2, 0, 0, 176, 177, 5, 27, 0, 0, 177, 179, 3, 22, 
	11, 3, 178, 172, 1, 0, 0, 0, 178, 175, 1, 0, 0, 0, 179, 182, 1, 0, 0, 0, 
	180, 178, 1, 0, 0, 0, 180, 181, 1, 0, 0, 0, 181, 23, 1, 0, 0, 0, 182, 180, 
	1, 0, 0, 0, 183, 184, 5, 29, 0, 0, 184, 206, 3, 30, 15, 0, 185, 186, 3, 
	30, 15, 0, 186, 187, 5, 30, 0, 0, 187, 188, 3, 28, 14, 0, 188, 206, 1, 
	0, 0, 0, 189, 190, 3, 58, 29, 0, 190, 191, 5, 30, 0, 0, 191, 192, 3, 28, 
	14, 0, 192, 206, 1, 0, 0, 0, 193, 194, 3, 30, 15, 0, 194, 195, 5, 31, 0, 
	0, 195, 196, 3, 48, 24, 0, 196, 206, 1, 0, 0, 0, 197, 198, 3, 30, 15, 0, 
	198, 199, 5, 32, 0, 0, 199, 200, 3, 50, 25, 0, 200, 206, 1, 0, 0, 0, 201, 
	202, 5, 80, 0, 0, 202, 203, 3, 24, 12, 0, 203, 204, 5, 81, 0, 0, 204, 206, 
	1, 0, 0, 0, 205, 183, 1, 0, 0, 0, 205, 185, 1, 0, 0, 0, 205, 189, 1, 0, 
	0, 0, 205, 193, 1, 0, 0, 0, 205, 197, 1, 0, 0, 0, 205, 201, 1, 0, 0, 0, 
	206, 25, 1, 0, 0, 0, 207, 209, 5, 62, 0, 0, 208, 210, 5, 62, 0, 0, 209, 
	208, 1, 0, 0, 0, 209, 210, 1, 0, 0, 0, 210, 212, 1, 0, 0, 0, 211, 213, 
	3, 32, 16, 0, 212, 211, 1, 0, 0, 0, 212, 213, 1, 0, 0, 0, 213, 225, 1, 
	0, 0, 0, 214, 216, 5, 17, 0, 0, 215, 217, 5, 62, 0, 0, 216, 215, 1, 0, 
	0, 0, 216, 217, 1, 0, 0, 0, 217, 222, 1, 0, 0, 0, 218, 219, 5, 87, 0, 0, 
	219, 220, 3, 40, 20, 0, 220, 221, 5, 88, 0, 0, 221, 223, 1, 0, 0, 0, 222, 
	218, 1, 0, 0, 0, 222, 223, 1, 0, 0, 0, 223, 225, 1, 0, 0, 0, 224, 207, 
	1, 0, 0, 0, 224, 214, 1, 0, 0, 0, 225, 27, 1, 0, 0, 0, 226, 231, 3, 54, 
	27, 0, 227, 231, 5, 57, 0, 0, 228, 231, 3, 30, 15, 0, 229, 231, 3, 58, 
	29, 0, 230, 226, 1, 0, 0, 0, 230, 227, 1, 0, 0, 0, 230, 228, 1, 0, 0, 0, 
	230, 229, 1, 0, 0, 0, 231, 29, 1, 0, 0, 0, 232, 234, 5, 62, 0, 0, 233, 
	235, 3, 32, 16, 0, 234, 233, 1, 0, 0, 0, 234, 235, 1, 0, 0, 0, 235, 238, 
	1, 0, 0, 0, 236, 237, 5, 83, 0, 0, 237, 239, 3, 44, 22, 0, 238, 236, 1, 
	0, 0, 0, 238, 239, 1, 0, 0, 0, 239, 31, 1, 0, 0, 0, 240, 244, 5, 87, 0, 
	0, 241, 245, 3, 34, 17, 0, 242, 245, 3, 36, 18, 0, 243, 245, 3, 38, 19, 
	0, 244, 241, 1, 0, 0, 0, 244, 242, 1, 0, 0, 0, 244, 243, 1, 0, 0, 0, 245, 
	246, 1, 0, 0, 0, 246, 247, 5, 88, 0, 0, 247, 33, 1, 0, 0, 0, 248, 249, 
	3, 44, 22, 0, 249, 250, 5, 30, 0, 0, 250, 251, 3, 42, 21, 0, 251, 35, 1, 
	0, 0, 0, 252, 253, 7, 1, 0, 0, 253, 37, 1, 0, 0, 0, 254, 255, 6, 19, -1, 
	0, 255, 258, 7, 2, 0, 0, 256, 257, 5, 82, 0, 0, 257, 259, 7, 3, 0, 0, 258, 
	256, 1, 0, 0, 0, 258, 259, 1, 0, 0, 0, 259, 275, 1, 0, 0, 0, 260, 263, 
	5, 61, 0, 0, 261, 262, 5, 82, 0, 0, 262, 264, 7, 3, 0, 0, 263, 261, 1, 
	0, 0, 0, 263, 264, 1, 0, 0, 0, 264, 275, 1, 0, 0, 0, 265, 275, 5, 57, 0, 
	0, 266, 267, 3, 44, 22, 0, 267, 268, 5, 30, 0, 0, 268, 269, 3, 42, 21, 
	0, 269, 275, 1, 0, 0, 0, 270, 271, 3, 44, 22, 0, 271, 272, 5, 32, 0, 0, 
	272, 273, 5, 60, 0, 0, 273, 275, 1, 0, 0, 0, 274, 254, 1, 0, 0, 0, 274, 
	260, 1, 0, 0, 0, 274, 265, 1, 0, 0, 0, 274, 266, 1, 0, 0, 0, 274, 270, 
	1, 0, 0, 0, 275, 284, 1, 0, 0, 0, 276, 277, 10, 2, 0, 0, 277, 278, 5, 26, 
	0, 0, 278, 283, 3, 38, 19, 3, 279, 280, 10, 1, 0, 0, 280, 281, 5, 27, 0, 
	0, 281, 283, 3, 38, 19, 2, 282, 276, 1, 0, 0, 0, 282, 279, 1, 0, 0, 0, 
	283, 286, 1, 0, 0, 0, 284, 282, 1, 0, 0, 0, 284, 285, 1, 0, 0, 0, 285, 
	39, 1, 0, 0, 0, 286, 284, 1, 0, 0, 0, 287, 291, 5, 18, 0, 0, 288, 291, 
	5, 19, 0, 0, 289, 291, 3, 34, 17, 0, 290, 287, 1, 0, 0, 0, 290, 288, 1, 
	0, 0, 0, 290, 289, 1, 0, 0, 0, 291, 41, 1, 0, 0, 0, 292, 298, 3, 54, 27, 
	0, 293, 298, 3, 44, 22, 0, 294, 298, 5, 57, 0, 0, 295, 298, 5, 58, 0, 0, 
	296, 298, 5, 59, 0, 0, 297, 292, 1, 0, 0, 0, 297, 293, 1, 0, 0, 0, 297, 
	294, 1, 0, 0, 0, 297, 295, 1, 0, 0, 0, 297, 296, 1, 0, 0, 0, 298, 43, 1, 
	0, 0, 0, 299, 304, 3, 46, 23, 0, 300, 301, 5, 83, 0, 0, 301, 303, 3, 46, 
	23, 0, 302, 300, 1, 0, 0, 0, 303, 306, 1, 0, 0, 0, 304, 302, 1, 0, 0, 0, 
	304, 305, 1, 0, 0, 0, 305, 45, 1, 0, 0, 0, 306, 304, 1, 0, 0, 0, 307, 309, 
	5, 62, 0, 0, 308, 310, 3, 32, 16, 0, 309, 308, 1, 0, 0, 0, 309, 310, 1, 
	0, 0, 0, 310, 47, 1, 0, 0, 0, 311, 312, 7, 4, 0, 0, 312, 49, 1, 0, 0, 0, 
	313, 314, 5, 89, 0, 0, 314, 319, 3, 52, 26, 0, 315, 316, 5, 82, 0, 0, 316, 
	318, 3, 52, 26, 0, 317, 315, 1, 0, 0, 0, 318, 321, 1, 0, 0, 0, 319, 317, 
	1, 0, 0, 0, 319, 320, 1, 0, 0, 0, 320, 322, 1, 0, 0, 0, 321, 319, 1, 0, 
	0, 0, 322, 323, 5, 90, 0, 0, 323, 329, 1, 0, 0, 0, 324, 329, 3, 62, 31, 
	0, 325, 326, 5, 89, 0, 0, 326, 327, 5, 64, 0, 0, 327, 329, 5, 90, 0, 0, 
	328, 313, 1, 0, 0, 0, 328, 324, 1, 0, 0, 0, 328, 325, 1, 0, 0, 0, 329, 
	51, 1, 0, 0, 0, 330, 334, 3, 54, 27, 0, 331, 334, 5, 57, 0, 0, 332, 334, 
	3, 62, 31, 0, 333, 330, 1, 0, 0, 0, 333, 331, 1, 0, 0, 0, 333, 332, 1, 
	0, 0, 0, 334, 53, 1, 0, 0, 0, 335, 343, 5, 72, 0, 0, 336, 343, 3, 56, 28, 
	0, 337, 343, 5, 69, 0, 0, 338, 343, 5, 70, 0, 0, 339, 343, 5, 71, 0, 0, 
	340, 343, 5, 21, 0, 0, 341, 343, 5, 20, 0, 0, 342, 335, 1, 0, 0, 0, 342, 
	336, 1, 0, 0, 0, 342, 337, 1, 0, 0, 0, 342, 338, 1, 0, 0, 0, 342, 339, 
	1, 0, 0, 0, 342, 340, 1, 0, 0, 0, 342, 341, 1, 0, 0, 0, 343, 55, 1, 0, 
	0, 0, 344, 351, 5, 65, 0, 0, 345, 351, 5, 66, 0, 0, 346, 351, 5, 67, 0, 
	0, 347, 351, 5, 68, 0, 0, 348, 349, 5, 86, 0, 0, 349, 351, 3, 56, 28, 0, 
	350, 344, 1, 0, 0, 0, 350, 345, 1, 0, 0, 0, 350, 346, 1, 0, 0, 0, 350, 
	347, 1, 0, 0, 0, 350, 348, 1, 0, 0, 0, 351, 57, 1, 0, 0, 0, 352, 367, 3, 
	62, 31, 0, 353, 354, 7, 5, 0, 0, 354, 363, 5, 80, 0, 0, 355, 360, 3, 28, 
	14, 0, 356, 357, 5, 82, 0, 0, 357, 359, 3, 28, 14, 0, 358, 356, 1, 0, 0, 
	0, 359, 362, 1, 0, 0, 0, 360, 358, 1, 0, 0, 0, 360, 361, 1, 0, 0, 0, 361, 
	364, 1, 0, 0, 0, 362, 360, 1, 0, 0, 0, 363, 355, 1, 0, 0, 0, 363, 364, 
	1, 0, 0, 0, 364, 365, 1, 0, 0, 0, 365, 367, 5, 81, 0, 0, 366, 352, 1, 0, 
	0, 0, 366, 353, 1, 0, 0, 0, 367, 59, 1, 0, 0, 0, 368, 369, 5, 51, 0, 0, 
	369, 375, 5, 80, 0, 0, 370, 372, 5, 16, 0, 0, 371, 370, 1, 0, 0, 0, 371, 
	372, 1, 0, 0, 0, 372, 373, 1, 0, 0, 0, 373, 376, 3, 30, 15, 0, 374, 376, 
	5, 84, 0, 0, 375, 371, 1, 0, 0, 0, 375, 374, 1, 0, 0, 0, 376, 377, 1, 0, 
	0, 0, 377, 384, 5, 81, 0, 0, 378, 379, 7, 6, 0, 0, 379, 380, 5, 80, 0, 
	0, 380, 381, 3, 30, 15, 0, 381, 382, 5, 81, 0, 0, 382, 384, 1, 0, 0, 0, 
	383, 368, 1, 0, 0, 0, 383, 378, 1, 0, 0, 0, 384, 61, 1, 0, 0, 0, 385, 386, 
	5, 56, 0, 0, 386, 387, 5, 80, 0, 0, 387, 388, 5, 72, 0, 0, 388, 389, 5, 
	82, 0, 0, 389, 390, 5, 72, 0, 0, 390, 391, 5, 82, 0, 0, 391, 392, 5, 72, 
	0, 0, 392, 393, 5, 81, 0, 0, 393, 63, 1, 0, 0, 0, 394, 395, 5, 22, 0, 0, 
	395, 397, 5, 65, 0, 0, 396, 398, 7, 7, 0, 0, 397, 396, 1, 0, 0, 0, 397, 
	398, 1, 0, 0, 0, 398, 65, 1, 0, 0, 0, 51, 69, 72, 75, 78, 84, 87, 94, 110, 
	117, 122, 134, 142, 144, 149, 155, 160, 164, 170, 178, 180, 205, 209, 212, 
	216, 222, 224, 230, 234, 238, 244, 258, 263, 274, 282, 284, 290, 297, 304, 
	309, 319, 328, 333, 342, 350, 360, 363, 366, 371, 375, 383, 397,
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

// AqlParserInit initializes any static state used to implement AqlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewAqlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func AqlParserInit() {
  staticData := &AqlParserParserStaticData
  staticData.once.Do(aqlparserParserInit)
}

// NewAqlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewAqlParser(input antlr.TokenStream) *AqlParser {
	AqlParserInit()
	this := new(AqlParser)
	this.BaseParser = antlr.NewBaseParser(input)
  staticData := &AqlParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "AqlParser.g4"

	return this
}


// AqlParser tokens.
const (
	AqlParserEOF = antlr.TokenEOF
	AqlParserWS = 1
	AqlParserUNICODE_BOM = 2
	AqlParserCOMMENT = 3
	AqlParserSELECT = 4
	AqlParserAS = 5
	AqlParserFROM = 6
	AqlParserWHERE = 7
	AqlParserORDER = 8
	AqlParserBY = 9
	AqlParserDESC = 10
	AqlParserDESCENDING = 11
	AqlParserASC = 12
	AqlParserASCENDING = 13
	AqlParserLIMIT = 14
	AqlParserOFFSET = 15
	AqlParserDISTINCT = 16
	AqlParserVERSION = 17
	AqlParserLATEST_VERSION = 18
	AqlParserALL_VERSIONS = 19
	AqlParserNULL = 20
	AqlParserBOOLEAN = 21
	AqlParserTOP = 22
	AqlParserFORWARD = 23
	AqlParserBACKWARD = 24
	AqlParserCONTAINS = 25
	AqlParserAND = 26
	AqlParserOR = 27
	AqlParserNOT = 28
	AqlParserEXISTS = 29
	AqlParserCOMPARISON_OPERATOR = 30
	AqlParserLIKE = 31
	AqlParserMATCHES = 32
	AqlParserSTRING_FUNCTION_ID = 33
	AqlParserNUMERIC_FUNCTION_ID = 34
	AqlParserDATE_TIME_FUNCTION_ID = 35
	AqlParserLENGTH = 36
	AqlParserPOSITION = 37
	AqlParserSUBSTRING = 38
	AqlParserCONCAT = 39
	AqlParserCONCAT_WS = 40
	AqlParserABS = 41
	AqlParserMOD = 42
	AqlParserCEIL = 43
	AqlParserFLOOR = 44
	AqlParserROUND = 45
	AqlParserCURRENT_DATE = 46
	AqlParserCURRENT_TIME = 47
	AqlParserCURRENT_DATE_TIME = 48
	AqlParserNOW = 49
	AqlParserCURRENT_TIMEZONE = 50
	AqlParserCOUNT = 51
	AqlParserMIN = 52
	AqlParserMAX = 53
	AqlParserSUM = 54
	AqlParserAVG = 55
	AqlParserTERMINOLOGY = 56
	AqlParserPARAMETER = 57
	AqlParserID_CODE = 58
	AqlParserAT_CODE = 59
	AqlParserCONTAINED_REGEX = 60
	AqlParserARCHETYPE_HRID = 61
	AqlParserIDENTIFIER = 62
	AqlParserTERM_CODE = 63
	AqlParserURI = 64
	AqlParserINTEGER = 65
	AqlParserREAL = 66
	AqlParserSCI_INTEGER = 67
	AqlParserSCI_REAL = 68
	AqlParserDATE = 69
	AqlParserTIME = 70
	AqlParserDATETIME = 71
	AqlParserSTRING = 72
	AqlParserSYM_SEMICOLON = 73
	AqlParserSYM_LT = 74
	AqlParserSYM_GT = 75
	AqlParserSYM_LE = 76
	AqlParserSYM_GE = 77
	AqlParserSYM_NE = 78
	AqlParserSYM_EQ = 79
	AqlParserSYM_LEFT_PAREN = 80
	AqlParserSYM_RIGHT_PAREN = 81
	AqlParserSYM_COMMA = 82
	AqlParserSYM_SLASH = 83
	AqlParserSYM_ASTERISK = 84
	AqlParserSYM_PLUS = 85
	AqlParserSYM_MINUS = 86
	AqlParserSYM_LEFT_BRACKET = 87
	AqlParserSYM_RIGHT_BRACKET = 88
	AqlParserSYM_LEFT_CURLY = 89
	AqlParserSYM_RIGHT_CURLY = 90
	AqlParserSYM_DOUBLE_DASH = 91
)

// AqlParser rules.
const (
	AqlParserRULE_selectQuery = 0
	AqlParserRULE_selectClause = 1
	AqlParserRULE_fromClause = 2
	AqlParserRULE_whereClause = 3
	AqlParserRULE_orderByClause = 4
	AqlParserRULE_limitClause = 5
	AqlParserRULE_selectExpr = 6
	AqlParserRULE_fromExpr = 7
	AqlParserRULE_whereExpr = 8
	AqlParserRULE_orderByExpr = 9
	AqlParserRULE_columnExpr = 10
	AqlParserRULE_containsExpr = 11
	AqlParserRULE_identifiedExpr = 12
	AqlParserRULE_classExprOperand = 13
	AqlParserRULE_terminal = 14
	AqlParserRULE_identifiedPath = 15
	AqlParserRULE_pathPredicate = 16
	AqlParserRULE_standardPredicate = 17
	AqlParserRULE_archetypePredicate = 18
	AqlParserRULE_nodePredicate = 19
	AqlParserRULE_versionPredicate = 20
	AqlParserRULE_pathPredicateOperand = 21
	AqlParserRULE_objectPath = 22
	AqlParserRULE_pathPart = 23
	AqlParserRULE_likeOperand = 24
	AqlParserRULE_matchesOperand = 25
	AqlParserRULE_valueListItem = 26
	AqlParserRULE_primitive = 27
	AqlParserRULE_numericPrimitive = 28
	AqlParserRULE_functionCall = 29
	AqlParserRULE_aggregateFunctionCall = 30
	AqlParserRULE_terminologyFunction = 31
	AqlParserRULE_top = 32
)

// ISelectQueryContext is an interface to support dynamic dispatch.
type ISelectQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SelectClause() ISelectClauseContext
	FromClause() IFromClauseContext
	EOF() antlr.TerminalNode
	WhereClause() IWhereClauseContext
	OrderByClause() IOrderByClauseContext
	LimitClause() ILimitClauseContext
	SYM_DOUBLE_DASH() antlr.TerminalNode

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
	p.RuleIndex = AqlParserRULE_selectQuery
	return p
}

func InitEmptySelectQueryContext(p *SelectQueryContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_selectQuery
}

func (*SelectQueryContext) IsSelectQueryContext() {}

func NewSelectQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectQueryContext {
	var p = new(SelectQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_selectQuery

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

func (s *SelectQueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(AqlParserEOF, 0)
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

func (s *SelectQueryContext) LimitClause() ILimitClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *SelectQueryContext) SYM_DOUBLE_DASH() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_DOUBLE_DASH, 0)
}

func (s *SelectQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterSelectQuery(s)
	}
}

func (s *SelectQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitSelectQuery(s)
	}
}




func (p *AqlParser) SelectQuery() (localctx ISelectQueryContext) {
	localctx = NewSelectQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, AqlParserRULE_selectQuery)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(66)
		p.SelectClause()
	}
	{
		p.SetState(67)
		p.FromClause()
	}
	p.SetState(69)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserWHERE {
		{
			p.SetState(68)
			p.WhereClause()
		}

	}
	p.SetState(72)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserORDER {
		{
			p.SetState(71)
			p.OrderByClause()
		}

	}
	p.SetState(75)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserLIMIT {
		{
			p.SetState(74)
			p.LimitClause()
		}

	}
	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserSYM_DOUBLE_DASH {
		{
			p.SetState(77)
			p.Match(AqlParserSYM_DOUBLE_DASH)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	}
	{
		p.SetState(80)
		p.Match(AqlParserEOF)
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
	Top() ITopContext
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
	p.RuleIndex = AqlParserRULE_selectClause
	return p
}

func InitEmptySelectClauseContext(p *SelectClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_selectClause
}

func (*SelectClauseContext) IsSelectClauseContext() {}

func NewSelectClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectClauseContext {
	var p = new(SelectClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_selectClause

	return p
}

func (s *SelectClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectClauseContext) SELECT() antlr.TerminalNode {
	return s.GetToken(AqlParserSELECT, 0)
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
	return s.GetToken(AqlParserDISTINCT, 0)
}

func (s *SelectClauseContext) Top() ITopContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITopContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITopContext)
}

func (s *SelectClauseContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AqlParserSYM_COMMA)
}

func (s *SelectClauseContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_COMMA, i)
}

func (s *SelectClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterSelectClause(s)
	}
}

func (s *SelectClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitSelectClause(s)
	}
}




func (p *AqlParser) SelectClause() (localctx ISelectClauseContext) {
	localctx = NewSelectClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, AqlParserRULE_selectClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(82)
		p.Match(AqlParserSELECT)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserDISTINCT {
		{
			p.SetState(83)
			p.Match(AqlParserDISTINCT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	}
	p.SetState(87)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserTOP {
		{
			p.SetState(86)
			p.Top()
		}

	}
	{
		p.SetState(89)
		p.SelectExpr()
	}
	p.SetState(94)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AqlParserSYM_COMMA {
		{
			p.SetState(90)
			p.Match(AqlParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(91)
			p.SelectExpr()
		}


		p.SetState(96)
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
	p.RuleIndex = AqlParserRULE_fromClause
	return p
}

func InitEmptyFromClauseContext(p *FromClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_fromClause
}

func (*FromClauseContext) IsFromClauseContext() {}

func NewFromClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromClauseContext {
	var p = new(FromClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_fromClause

	return p
}

func (s *FromClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FromClauseContext) FROM() antlr.TerminalNode {
	return s.GetToken(AqlParserFROM, 0)
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
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterFromClause(s)
	}
}

func (s *FromClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitFromClause(s)
	}
}




func (p *AqlParser) FromClause() (localctx IFromClauseContext) {
	localctx = NewFromClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, AqlParserRULE_fromClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(97)
		p.Match(AqlParserFROM)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(98)
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
	p.RuleIndex = AqlParserRULE_whereClause
	return p
}

func InitEmptyWhereClauseContext(p *WhereClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_whereClause
}

func (*WhereClauseContext) IsWhereClauseContext() {}

func NewWhereClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereClauseContext {
	var p = new(WhereClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_whereClause

	return p
}

func (s *WhereClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereClauseContext) WHERE() antlr.TerminalNode {
	return s.GetToken(AqlParserWHERE, 0)
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
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterWhereClause(s)
	}
}

func (s *WhereClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitWhereClause(s)
	}
}




func (p *AqlParser) WhereClause() (localctx IWhereClauseContext) {
	localctx = NewWhereClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, AqlParserRULE_whereClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(100)
		p.Match(AqlParserWHERE)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(101)
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


// IOrderByClauseContext is an interface to support dynamic dispatch.
type IOrderByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ORDER() antlr.TerminalNode
	BY() antlr.TerminalNode
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
	p.RuleIndex = AqlParserRULE_orderByClause
	return p
}

func InitEmptyOrderByClauseContext(p *OrderByClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_orderByClause
}

func (*OrderByClauseContext) IsOrderByClauseContext() {}

func NewOrderByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByClauseContext {
	var p = new(OrderByClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_orderByClause

	return p
}

func (s *OrderByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByClauseContext) ORDER() antlr.TerminalNode {
	return s.GetToken(AqlParserORDER, 0)
}

func (s *OrderByClauseContext) BY() antlr.TerminalNode {
	return s.GetToken(AqlParserBY, 0)
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
	return s.GetTokens(AqlParserSYM_COMMA)
}

func (s *OrderByClauseContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_COMMA, i)
}

func (s *OrderByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *OrderByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterOrderByClause(s)
	}
}

func (s *OrderByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitOrderByClause(s)
	}
}




func (p *AqlParser) OrderByClause() (localctx IOrderByClauseContext) {
	localctx = NewOrderByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, AqlParserRULE_orderByClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Match(AqlParserORDER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(104)
		p.Match(AqlParserBY)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(105)
		p.OrderByExpr()
	}
	p.SetState(110)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AqlParserSYM_COMMA {
		{
			p.SetState(106)
			p.Match(AqlParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(107)
			p.OrderByExpr()
		}


		p.SetState(112)
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


// ILimitClauseContext is an interface to support dynamic dispatch.
type ILimitClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLimit returns the limit token.
	GetLimit() antlr.Token 

	// GetOffset returns the offset token.
	GetOffset() antlr.Token 


	// SetLimit sets the limit token.
	SetLimit(antlr.Token) 

	// SetOffset sets the offset token.
	SetOffset(antlr.Token) 


	// Getter signatures
	LIMIT() antlr.TerminalNode
	AllINTEGER() []antlr.TerminalNode
	INTEGER(i int) antlr.TerminalNode
	OFFSET() antlr.TerminalNode

	// IsLimitClauseContext differentiates from other interfaces.
	IsLimitClauseContext()
}

type LimitClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	limit antlr.Token
	offset antlr.Token
}

func NewEmptyLimitClauseContext() *LimitClauseContext {
	var p = new(LimitClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_limitClause
	return p
}

func InitEmptyLimitClauseContext(p *LimitClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_limitClause
}

func (*LimitClauseContext) IsLimitClauseContext() {}

func NewLimitClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseContext {
	var p = new(LimitClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_limitClause

	return p
}

func (s *LimitClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseContext) GetLimit() antlr.Token { return s.limit }

func (s *LimitClauseContext) GetOffset() antlr.Token { return s.offset }


func (s *LimitClauseContext) SetLimit(v antlr.Token) { s.limit = v }

func (s *LimitClauseContext) SetOffset(v antlr.Token) { s.offset = v }


func (s *LimitClauseContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(AqlParserLIMIT, 0)
}

func (s *LimitClauseContext) AllINTEGER() []antlr.TerminalNode {
	return s.GetTokens(AqlParserINTEGER)
}

func (s *LimitClauseContext) INTEGER(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserINTEGER, i)
}

func (s *LimitClauseContext) OFFSET() antlr.TerminalNode {
	return s.GetToken(AqlParserOFFSET, 0)
}

func (s *LimitClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LimitClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterLimitClause(s)
	}
}

func (s *LimitClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitLimitClause(s)
	}
}




func (p *AqlParser) LimitClause() (localctx ILimitClauseContext) {
	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, AqlParserRULE_limitClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(113)
		p.Match(AqlParserLIMIT)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(114)

		var _m = p.Match(AqlParserINTEGER)

		localctx.(*LimitClauseContext).limit = _m
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserOFFSET {
		{
			p.SetState(115)
			p.Match(AqlParserOFFSET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(116)

			var _m = p.Match(AqlParserINTEGER)

			localctx.(*LimitClauseContext).offset = _m
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


// ISelectExprContext is an interface to support dynamic dispatch.
type ISelectExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAliasName returns the aliasName token.
	GetAliasName() antlr.Token 


	// SetAliasName sets the aliasName token.
	SetAliasName(antlr.Token) 


	// Getter signatures
	ColumnExpr() IColumnExprContext
	AS() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsSelectExprContext differentiates from other interfaces.
	IsSelectExprContext()
}

type SelectExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	aliasName antlr.Token
}

func NewEmptySelectExprContext() *SelectExprContext {
	var p = new(SelectExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_selectExpr
	return p
}

func InitEmptySelectExprContext(p *SelectExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_selectExpr
}

func (*SelectExprContext) IsSelectExprContext() {}

func NewSelectExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectExprContext {
	var p = new(SelectExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_selectExpr

	return p
}

func (s *SelectExprContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectExprContext) GetAliasName() antlr.Token { return s.aliasName }


func (s *SelectExprContext) SetAliasName(v antlr.Token) { s.aliasName = v }


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

func (s *SelectExprContext) AS() antlr.TerminalNode {
	return s.GetToken(AqlParserAS, 0)
}

func (s *SelectExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AqlParserIDENTIFIER, 0)
}

func (s *SelectExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterSelectExpr(s)
	}
}

func (s *SelectExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitSelectExpr(s)
	}
}




func (p *AqlParser) SelectExpr() (localctx ISelectExprContext) {
	localctx = NewSelectExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, AqlParserRULE_selectExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
		p.ColumnExpr()
	}
	p.SetState(122)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserAS {
		{
			p.SetState(120)
			p.Match(AqlParserAS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(121)

			var _m = p.Match(AqlParserIDENTIFIER)

			localctx.(*SelectExprContext).aliasName = _m
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


// IFromExprContext is an interface to support dynamic dispatch.
type IFromExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ContainsExpr() IContainsExprContext

	// IsFromExprContext differentiates from other interfaces.
	IsFromExprContext()
}

type FromExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFromExprContext() *FromExprContext {
	var p = new(FromExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_fromExpr
	return p
}

func InitEmptyFromExprContext(p *FromExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_fromExpr
}

func (*FromExprContext) IsFromExprContext() {}

func NewFromExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromExprContext {
	var p = new(FromExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_fromExpr

	return p
}

func (s *FromExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FromExprContext) ContainsExpr() IContainsExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContainsExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContainsExprContext)
}

func (s *FromExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FromExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterFromExpr(s)
	}
}

func (s *FromExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitFromExpr(s)
	}
}




func (p *AqlParser) FromExpr() (localctx IFromExprContext) {
	localctx = NewFromExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, AqlParserRULE_fromExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(124)
		p.containsExpr(0)
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
	IdentifiedExpr() IIdentifiedExprContext
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
	p.RuleIndex = AqlParserRULE_whereExpr
	return p
}

func InitEmptyWhereExprContext(p *WhereExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_whereExpr
}

func (*WhereExprContext) IsWhereExprContext() {}

func NewWhereExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereExprContext {
	var p = new(WhereExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_whereExpr

	return p
}

func (s *WhereExprContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereExprContext) IdentifiedExpr() IIdentifiedExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedExprContext)
}

func (s *WhereExprContext) NOT() antlr.TerminalNode {
	return s.GetToken(AqlParserNOT, 0)
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
	return s.GetToken(AqlParserSYM_LEFT_PAREN, 0)
}

func (s *WhereExprContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_PAREN, 0)
}

func (s *WhereExprContext) AND() antlr.TerminalNode {
	return s.GetToken(AqlParserAND, 0)
}

func (s *WhereExprContext) OR() antlr.TerminalNode {
	return s.GetToken(AqlParserOR, 0)
}

func (s *WhereExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WhereExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterWhereExpr(s)
	}
}

func (s *WhereExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitWhereExpr(s)
	}
}





func (p *AqlParser) WhereExpr() (localctx IWhereExprContext) {
	return p.whereExpr(0)
}

func (p *AqlParser) whereExpr(_p int) (localctx IWhereExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewWhereExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IWhereExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 16
	p.EnterRecursionRule(localctx, 16, AqlParserRULE_whereExpr, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(134)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(127)
			p.IdentifiedExpr()
		}


	case 2:
		{
			p.SetState(128)
			p.Match(AqlParserNOT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(129)
			p.whereExpr(4)
		}


	case 3:
		{
			p.SetState(130)
			p.Match(AqlParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(131)
			p.whereExpr(0)
		}
		{
			p.SetState(132)
			p.Match(AqlParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(142)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext()) {
			case 1:
				localctx = NewWhereExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AqlParserRULE_whereExpr)
				p.SetState(136)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(137)
					p.Match(AqlParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(138)
					p.whereExpr(4)
				}


			case 2:
				localctx = NewWhereExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AqlParserRULE_whereExpr)
				p.SetState(139)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(140)
					p.Match(AqlParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(141)
					p.whereExpr(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(146)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext())
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
	IdentifiedPath() IIdentifiedPathContext
	DESCENDING() antlr.TerminalNode
	DESC() antlr.TerminalNode
	ASCENDING() antlr.TerminalNode
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
	p.RuleIndex = AqlParserRULE_orderByExpr
	return p
}

func InitEmptyOrderByExprContext(p *OrderByExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_orderByExpr
}

func (*OrderByExprContext) IsOrderByExprContext() {}

func NewOrderByExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByExprContext {
	var p = new(OrderByExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_orderByExpr

	return p
}

func (s *OrderByExprContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByExprContext) GetOrder() antlr.Token { return s.order }


func (s *OrderByExprContext) SetOrder(v antlr.Token) { s.order = v }


func (s *OrderByExprContext) IdentifiedPath() IIdentifiedPathContext {
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

func (s *OrderByExprContext) DESCENDING() antlr.TerminalNode {
	return s.GetToken(AqlParserDESCENDING, 0)
}

func (s *OrderByExprContext) DESC() antlr.TerminalNode {
	return s.GetToken(AqlParserDESC, 0)
}

func (s *OrderByExprContext) ASCENDING() antlr.TerminalNode {
	return s.GetToken(AqlParserASCENDING, 0)
}

func (s *OrderByExprContext) ASC() antlr.TerminalNode {
	return s.GetToken(AqlParserASC, 0)
}

func (s *OrderByExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *OrderByExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterOrderByExpr(s)
	}
}

func (s *OrderByExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitOrderByExpr(s)
	}
}




func (p *AqlParser) OrderByExpr() (localctx IOrderByExprContext) {
	localctx = NewOrderByExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, AqlParserRULE_orderByExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(147)
		p.IdentifiedPath()
	}
	p.SetState(149)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if ((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 15360) != 0) {
		{
			p.SetState(148)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*OrderByExprContext).order = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 15360) != 0)) {
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
	IdentifiedPath() IIdentifiedPathContext
	Primitive() IPrimitiveContext
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
	p.RuleIndex = AqlParserRULE_columnExpr
	return p
}

func InitEmptyColumnExprContext(p *ColumnExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_columnExpr
}

func (*ColumnExprContext) IsColumnExprContext() {}

func NewColumnExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnExprContext {
	var p = new(ColumnExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_columnExpr

	return p
}

func (s *ColumnExprContext) GetParser() antlr.Parser { return s.parser }

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
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterColumnExpr(s)
	}
}

func (s *ColumnExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitColumnExpr(s)
	}
}




func (p *AqlParser) ColumnExpr() (localctx IColumnExprContext) {
	localctx = NewColumnExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, AqlParserRULE_columnExpr)
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(151)
			p.IdentifiedPath()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(152)
			p.Primitive()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(153)
			p.AggregateFunctionCall()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(154)
			p.FunctionCall()
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


// IContainsExprContext is an interface to support dynamic dispatch.
type IContainsExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ClassExprOperand() IClassExprOperandContext
	CONTAINS() antlr.TerminalNode
	AllContainsExpr() []IContainsExprContext
	ContainsExpr(i int) IContainsExprContext
	NOT() antlr.TerminalNode
	SYM_LEFT_PAREN() antlr.TerminalNode
	SYM_RIGHT_PAREN() antlr.TerminalNode
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode

	// IsContainsExprContext differentiates from other interfaces.
	IsContainsExprContext()
}

type ContainsExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyContainsExprContext() *ContainsExprContext {
	var p = new(ContainsExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_containsExpr
	return p
}

func InitEmptyContainsExprContext(p *ContainsExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_containsExpr
}

func (*ContainsExprContext) IsContainsExprContext() {}

func NewContainsExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContainsExprContext {
	var p = new(ContainsExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_containsExpr

	return p
}

func (s *ContainsExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ContainsExprContext) ClassExprOperand() IClassExprOperandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClassExprOperandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClassExprOperandContext)
}

func (s *ContainsExprContext) CONTAINS() antlr.TerminalNode {
	return s.GetToken(AqlParserCONTAINS, 0)
}

func (s *ContainsExprContext) AllContainsExpr() []IContainsExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IContainsExprContext); ok {
			len++
		}
	}

	tst := make([]IContainsExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IContainsExprContext); ok {
			tst[i] = t.(IContainsExprContext)
			i++
		}
	}

	return tst
}

func (s *ContainsExprContext) ContainsExpr(i int) IContainsExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContainsExprContext); ok {
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

	return t.(IContainsExprContext)
}

func (s *ContainsExprContext) NOT() antlr.TerminalNode {
	return s.GetToken(AqlParserNOT, 0)
}

func (s *ContainsExprContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_PAREN, 0)
}

func (s *ContainsExprContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_PAREN, 0)
}

func (s *ContainsExprContext) AND() antlr.TerminalNode {
	return s.GetToken(AqlParserAND, 0)
}

func (s *ContainsExprContext) OR() antlr.TerminalNode {
	return s.GetToken(AqlParserOR, 0)
}

func (s *ContainsExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContainsExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ContainsExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterContainsExpr(s)
	}
}

func (s *ContainsExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitContainsExpr(s)
	}
}





func (p *AqlParser) ContainsExpr() (localctx IContainsExprContext) {
	return p.containsExpr(0)
}

func (p *AqlParser) containsExpr(_p int) (localctx IContainsExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewContainsExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IContainsExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 22
	p.EnterRecursionRule(localctx, 22, AqlParserRULE_containsExpr, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserVERSION, AqlParserIDENTIFIER:
		{
			p.SetState(158)
			p.ClassExprOperand()
		}
		p.SetState(164)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 16, p.GetParserRuleContext()) == 1 {
			p.SetState(160)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			if _la == AqlParserNOT {
				{
					p.SetState(159)
					p.Match(AqlParserNOT)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}

			}
			{
				p.SetState(162)
				p.Match(AqlParserCONTAINS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(163)
				p.containsExpr(0)
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}


	case AqlParserSYM_LEFT_PAREN:
		{
			p.SetState(166)
			p.Match(AqlParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(167)
			p.containsExpr(0)
		}
		{
			p.SetState(168)
			p.Match(AqlParserSYM_RIGHT_PAREN)
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
	p.SetState(180)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(178)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) {
			case 1:
				localctx = NewContainsExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AqlParserRULE_containsExpr)
				p.SetState(172)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(173)
					p.Match(AqlParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(174)
					p.containsExpr(4)
				}


			case 2:
				localctx = NewContainsExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AqlParserRULE_containsExpr)
				p.SetState(175)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(176)
					p.Match(AqlParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(177)
					p.containsExpr(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(182)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext())
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


// IIdentifiedExprContext is an interface to support dynamic dispatch.
type IIdentifiedExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXISTS() antlr.TerminalNode
	IdentifiedPath() IIdentifiedPathContext
	COMPARISON_OPERATOR() antlr.TerminalNode
	Terminal() ITerminalContext
	FunctionCall() IFunctionCallContext
	LIKE() antlr.TerminalNode
	LikeOperand() ILikeOperandContext
	MATCHES() antlr.TerminalNode
	MatchesOperand() IMatchesOperandContext
	SYM_LEFT_PAREN() antlr.TerminalNode
	IdentifiedExpr() IIdentifiedExprContext
	SYM_RIGHT_PAREN() antlr.TerminalNode

	// IsIdentifiedExprContext differentiates from other interfaces.
	IsIdentifiedExprContext()
}

type IdentifiedExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifiedExprContext() *IdentifiedExprContext {
	var p = new(IdentifiedExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_identifiedExpr
	return p
}

func InitEmptyIdentifiedExprContext(p *IdentifiedExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_identifiedExpr
}

func (*IdentifiedExprContext) IsIdentifiedExprContext() {}

func NewIdentifiedExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifiedExprContext {
	var p = new(IdentifiedExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_identifiedExpr

	return p
}

func (s *IdentifiedExprContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifiedExprContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(AqlParserEXISTS, 0)
}

func (s *IdentifiedExprContext) IdentifiedPath() IIdentifiedPathContext {
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

func (s *IdentifiedExprContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(AqlParserCOMPARISON_OPERATOR, 0)
}

func (s *IdentifiedExprContext) Terminal() ITerminalContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITerminalContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITerminalContext)
}

func (s *IdentifiedExprContext) FunctionCall() IFunctionCallContext {
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

func (s *IdentifiedExprContext) LIKE() antlr.TerminalNode {
	return s.GetToken(AqlParserLIKE, 0)
}

func (s *IdentifiedExprContext) LikeOperand() ILikeOperandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILikeOperandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILikeOperandContext)
}

func (s *IdentifiedExprContext) MATCHES() antlr.TerminalNode {
	return s.GetToken(AqlParserMATCHES, 0)
}

func (s *IdentifiedExprContext) MatchesOperand() IMatchesOperandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMatchesOperandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMatchesOperandContext)
}

func (s *IdentifiedExprContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_PAREN, 0)
}

func (s *IdentifiedExprContext) IdentifiedExpr() IIdentifiedExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifiedExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifiedExprContext)
}

func (s *IdentifiedExprContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_PAREN, 0)
}

func (s *IdentifiedExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifiedExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IdentifiedExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterIdentifiedExpr(s)
	}
}

func (s *IdentifiedExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitIdentifiedExpr(s)
	}
}




func (p *AqlParser) IdentifiedExpr() (localctx IIdentifiedExprContext) {
	localctx = NewIdentifiedExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, AqlParserRULE_identifiedExpr)
	p.SetState(205)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(183)
			p.Match(AqlParserEXISTS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(184)
			p.IdentifiedPath()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(185)
			p.IdentifiedPath()
		}
		{
			p.SetState(186)
			p.Match(AqlParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(187)
			p.Terminal()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(189)
			p.FunctionCall()
		}
		{
			p.SetState(190)
			p.Match(AqlParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(191)
			p.Terminal()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(193)
			p.IdentifiedPath()
		}
		{
			p.SetState(194)
			p.Match(AqlParserLIKE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(195)
			p.LikeOperand()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(197)
			p.IdentifiedPath()
		}
		{
			p.SetState(198)
			p.Match(AqlParserMATCHES)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(199)
			p.MatchesOperand()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(201)
			p.Match(AqlParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(202)
			p.IdentifiedExpr()
		}
		{
			p.SetState(203)
			p.Match(AqlParserSYM_RIGHT_PAREN)
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


// IClassExprOperandContext is an interface to support dynamic dispatch.
type IClassExprOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsClassExprOperandContext differentiates from other interfaces.
	IsClassExprOperandContext()
}

type ClassExprOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClassExprOperandContext() *ClassExprOperandContext {
	var p = new(ClassExprOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_classExprOperand
	return p
}

func InitEmptyClassExprOperandContext(p *ClassExprOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_classExprOperand
}

func (*ClassExprOperandContext) IsClassExprOperandContext() {}

func NewClassExprOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassExprOperandContext {
	var p = new(ClassExprOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_classExprOperand

	return p
}

func (s *ClassExprOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *ClassExprOperandContext) CopyAll(ctx *ClassExprOperandContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ClassExprOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClassExprOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}




type ClassExpressionContext struct {
	ClassExprOperandContext
	variable antlr.Token
}

func NewClassExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ClassExpressionContext {
	var p = new(ClassExpressionContext)

	InitEmptyClassExprOperandContext(&p.ClassExprOperandContext)
	p.parser = parser
	p.CopyAll(ctx.(*ClassExprOperandContext))

	return p
}


func (s *ClassExpressionContext) GetVariable() antlr.Token { return s.variable }


func (s *ClassExpressionContext) SetVariable(v antlr.Token) { s.variable = v }

func (s *ClassExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClassExpressionContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(AqlParserIDENTIFIER)
}

func (s *ClassExpressionContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserIDENTIFIER, i)
}

func (s *ClassExpressionContext) PathPredicate() IPathPredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathPredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathPredicateContext)
}


func (s *ClassExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterClassExpression(s)
	}
}

func (s *ClassExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitClassExpression(s)
	}
}


type VersionClassExprContext struct {
	ClassExprOperandContext
	variable antlr.Token
}

func NewVersionClassExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VersionClassExprContext {
	var p = new(VersionClassExprContext)

	InitEmptyClassExprOperandContext(&p.ClassExprOperandContext)
	p.parser = parser
	p.CopyAll(ctx.(*ClassExprOperandContext))

	return p
}


func (s *VersionClassExprContext) GetVariable() antlr.Token { return s.variable }


func (s *VersionClassExprContext) SetVariable(v antlr.Token) { s.variable = v }

func (s *VersionClassExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VersionClassExprContext) VERSION() antlr.TerminalNode {
	return s.GetToken(AqlParserVERSION, 0)
}

func (s *VersionClassExprContext) SYM_LEFT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_BRACKET, 0)
}

func (s *VersionClassExprContext) VersionPredicate() IVersionPredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVersionPredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVersionPredicateContext)
}

func (s *VersionClassExprContext) SYM_RIGHT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_BRACKET, 0)
}

func (s *VersionClassExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AqlParserIDENTIFIER, 0)
}


func (s *VersionClassExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterVersionClassExpr(s)
	}
}

func (s *VersionClassExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitVersionClassExpr(s)
	}
}



func (p *AqlParser) ClassExprOperand() (localctx IClassExprOperandContext) {
	localctx = NewClassExprOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, AqlParserRULE_classExprOperand)
	p.SetState(224)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserIDENTIFIER:
		localctx = NewClassExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(207)
			p.Match(AqlParserIDENTIFIER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(209)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(208)

				var _m = p.Match(AqlParserIDENTIFIER)

				localctx.(*ClassExpressionContext).variable = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}
		p.SetState(212)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 22, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(211)
				p.PathPredicate()
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}


	case AqlParserVERSION:
		localctx = NewVersionClassExprContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(214)
			p.Match(AqlParserVERSION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(216)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(215)

				var _m = p.Match(AqlParserIDENTIFIER)

				localctx.(*VersionClassExprContext).variable = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}
		p.SetState(222)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(218)
				p.Match(AqlParserSYM_LEFT_BRACKET)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(219)
				p.VersionPredicate()
			}
			{
				p.SetState(220)
				p.Match(AqlParserSYM_RIGHT_BRACKET)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}

			} else if p.HasError() { // JIM
				goto errorExit
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


// ITerminalContext is an interface to support dynamic dispatch.
type ITerminalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primitive() IPrimitiveContext
	PARAMETER() antlr.TerminalNode
	IdentifiedPath() IIdentifiedPathContext
	FunctionCall() IFunctionCallContext

	// IsTerminalContext differentiates from other interfaces.
	IsTerminalContext()
}

type TerminalContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTerminalContext() *TerminalContext {
	var p = new(TerminalContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_terminal
	return p
}

func InitEmptyTerminalContext(p *TerminalContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_terminal
}

func (*TerminalContext) IsTerminalContext() {}

func NewTerminalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TerminalContext {
	var p = new(TerminalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_terminal

	return p
}

func (s *TerminalContext) GetParser() antlr.Parser { return s.parser }

func (s *TerminalContext) Primitive() IPrimitiveContext {
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

func (s *TerminalContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AqlParserPARAMETER, 0)
}

func (s *TerminalContext) IdentifiedPath() IIdentifiedPathContext {
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

func (s *TerminalContext) FunctionCall() IFunctionCallContext {
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

func (s *TerminalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TerminalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TerminalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterTerminal(s)
	}
}

func (s *TerminalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitTerminal(s)
	}
}




func (p *AqlParser) Terminal() (localctx ITerminalContext) {
	localctx = NewTerminalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, AqlParserRULE_terminal)
	p.SetState(230)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(226)
			p.Primitive()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(227)
			p.Match(AqlParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(228)
			p.IdentifiedPath()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(229)
			p.FunctionCall()
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


// IIdentifiedPathContext is an interface to support dynamic dispatch.
type IIdentifiedPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	PathPredicate() IPathPredicateContext
	SYM_SLASH() antlr.TerminalNode
	ObjectPath() IObjectPathContext

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
	p.RuleIndex = AqlParserRULE_identifiedPath
	return p
}

func InitEmptyIdentifiedPathContext(p *IdentifiedPathContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_identifiedPath
}

func (*IdentifiedPathContext) IsIdentifiedPathContext() {}

func NewIdentifiedPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifiedPathContext {
	var p = new(IdentifiedPathContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_identifiedPath

	return p
}

func (s *IdentifiedPathContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifiedPathContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AqlParserIDENTIFIER, 0)
}

func (s *IdentifiedPathContext) PathPredicate() IPathPredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathPredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathPredicateContext)
}

func (s *IdentifiedPathContext) SYM_SLASH() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_SLASH, 0)
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

func (s *IdentifiedPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifiedPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IdentifiedPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterIdentifiedPath(s)
	}
}

func (s *IdentifiedPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitIdentifiedPath(s)
	}
}




func (p *AqlParser) IdentifiedPath() (localctx IIdentifiedPathContext) {
	localctx = NewIdentifiedPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, AqlParserRULE_identifiedPath)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(232)
		p.Match(AqlParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(234)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(233)
			p.PathPredicate()
		}

		} else if p.HasError() { // JIM
			goto errorExit
	}
	p.SetState(238)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 28, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(236)
			p.Match(AqlParserSYM_SLASH)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(237)
			p.ObjectPath()
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


// IPathPredicateContext is an interface to support dynamic dispatch.
type IPathPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SYM_LEFT_BRACKET() antlr.TerminalNode
	SYM_RIGHT_BRACKET() antlr.TerminalNode
	StandardPredicate() IStandardPredicateContext
	ArchetypePredicate() IArchetypePredicateContext
	NodePredicate() INodePredicateContext

	// IsPathPredicateContext differentiates from other interfaces.
	IsPathPredicateContext()
}

type PathPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathPredicateContext() *PathPredicateContext {
	var p = new(PathPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_pathPredicate
	return p
}

func InitEmptyPathPredicateContext(p *PathPredicateContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_pathPredicate
}

func (*PathPredicateContext) IsPathPredicateContext() {}

func NewPathPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathPredicateContext {
	var p = new(PathPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_pathPredicate

	return p
}

func (s *PathPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PathPredicateContext) SYM_LEFT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_BRACKET, 0)
}

func (s *PathPredicateContext) SYM_RIGHT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_BRACKET, 0)
}

func (s *PathPredicateContext) StandardPredicate() IStandardPredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStandardPredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStandardPredicateContext)
}

func (s *PathPredicateContext) ArchetypePredicate() IArchetypePredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArchetypePredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArchetypePredicateContext)
}

func (s *PathPredicateContext) NodePredicate() INodePredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INodePredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INodePredicateContext)
}

func (s *PathPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterPathPredicate(s)
	}
}

func (s *PathPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitPathPredicate(s)
	}
}




func (p *AqlParser) PathPredicate() (localctx IPathPredicateContext) {
	localctx = NewPathPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, AqlParserRULE_pathPredicate)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(240)
		p.Match(AqlParserSYM_LEFT_BRACKET)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(244)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(241)
			p.StandardPredicate()
		}


	case 2:
		{
			p.SetState(242)
			p.ArchetypePredicate()
		}


	case 3:
		{
			p.SetState(243)
			p.nodePredicate(0)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	{
		p.SetState(246)
		p.Match(AqlParserSYM_RIGHT_BRACKET)
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


// IStandardPredicateContext is an interface to support dynamic dispatch.
type IStandardPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ObjectPath() IObjectPathContext
	COMPARISON_OPERATOR() antlr.TerminalNode
	PathPredicateOperand() IPathPredicateOperandContext

	// IsStandardPredicateContext differentiates from other interfaces.
	IsStandardPredicateContext()
}

type StandardPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStandardPredicateContext() *StandardPredicateContext {
	var p = new(StandardPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_standardPredicate
	return p
}

func InitEmptyStandardPredicateContext(p *StandardPredicateContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_standardPredicate
}

func (*StandardPredicateContext) IsStandardPredicateContext() {}

func NewStandardPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StandardPredicateContext {
	var p = new(StandardPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_standardPredicate

	return p
}

func (s *StandardPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *StandardPredicateContext) ObjectPath() IObjectPathContext {
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

func (s *StandardPredicateContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(AqlParserCOMPARISON_OPERATOR, 0)
}

func (s *StandardPredicateContext) PathPredicateOperand() IPathPredicateOperandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathPredicateOperandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathPredicateOperandContext)
}

func (s *StandardPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StandardPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StandardPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterStandardPredicate(s)
	}
}

func (s *StandardPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitStandardPredicate(s)
	}
}




func (p *AqlParser) StandardPredicate() (localctx IStandardPredicateContext) {
	localctx = NewStandardPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, AqlParserRULE_standardPredicate)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(248)
		p.ObjectPath()
	}
	{
		p.SetState(249)
		p.Match(AqlParserCOMPARISON_OPERATOR)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(250)
		p.PathPredicateOperand()
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


// IArchetypePredicateContext is an interface to support dynamic dispatch.
type IArchetypePredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ARCHETYPE_HRID() antlr.TerminalNode
	PARAMETER() antlr.TerminalNode

	// IsArchetypePredicateContext differentiates from other interfaces.
	IsArchetypePredicateContext()
}

type ArchetypePredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArchetypePredicateContext() *ArchetypePredicateContext {
	var p = new(ArchetypePredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_archetypePredicate
	return p
}

func InitEmptyArchetypePredicateContext(p *ArchetypePredicateContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_archetypePredicate
}

func (*ArchetypePredicateContext) IsArchetypePredicateContext() {}

func NewArchetypePredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArchetypePredicateContext {
	var p = new(ArchetypePredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_archetypePredicate

	return p
}

func (s *ArchetypePredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *ArchetypePredicateContext) ARCHETYPE_HRID() antlr.TerminalNode {
	return s.GetToken(AqlParserARCHETYPE_HRID, 0)
}

func (s *ArchetypePredicateContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AqlParserPARAMETER, 0)
}

func (s *ArchetypePredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArchetypePredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ArchetypePredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterArchetypePredicate(s)
	}
}

func (s *ArchetypePredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitArchetypePredicate(s)
	}
}




func (p *AqlParser) ArchetypePredicate() (localctx IArchetypePredicateContext) {
	localctx = NewArchetypePredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, AqlParserRULE_archetypePredicate)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(252)
		_la = p.GetTokenStream().LA(1)

		if !(_la == AqlParserPARAMETER || _la == AqlParserARCHETYPE_HRID) {
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


// INodePredicateContext is an interface to support dynamic dispatch.
type INodePredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID_CODE() []antlr.TerminalNode
	ID_CODE(i int) antlr.TerminalNode
	AllAT_CODE() []antlr.TerminalNode
	AT_CODE(i int) antlr.TerminalNode
	SYM_COMMA() antlr.TerminalNode
	STRING() antlr.TerminalNode
	PARAMETER() antlr.TerminalNode
	TERM_CODE() antlr.TerminalNode
	ARCHETYPE_HRID() antlr.TerminalNode
	ObjectPath() IObjectPathContext
	COMPARISON_OPERATOR() antlr.TerminalNode
	PathPredicateOperand() IPathPredicateOperandContext
	MATCHES() antlr.TerminalNode
	CONTAINED_REGEX() antlr.TerminalNode
	AllNodePredicate() []INodePredicateContext
	NodePredicate(i int) INodePredicateContext
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode

	// IsNodePredicateContext differentiates from other interfaces.
	IsNodePredicateContext()
}

type NodePredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNodePredicateContext() *NodePredicateContext {
	var p = new(NodePredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_nodePredicate
	return p
}

func InitEmptyNodePredicateContext(p *NodePredicateContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_nodePredicate
}

func (*NodePredicateContext) IsNodePredicateContext() {}

func NewNodePredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NodePredicateContext {
	var p = new(NodePredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_nodePredicate

	return p
}

func (s *NodePredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *NodePredicateContext) AllID_CODE() []antlr.TerminalNode {
	return s.GetTokens(AqlParserID_CODE)
}

func (s *NodePredicateContext) ID_CODE(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserID_CODE, i)
}

func (s *NodePredicateContext) AllAT_CODE() []antlr.TerminalNode {
	return s.GetTokens(AqlParserAT_CODE)
}

func (s *NodePredicateContext) AT_CODE(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserAT_CODE, i)
}

func (s *NodePredicateContext) SYM_COMMA() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_COMMA, 0)
}

func (s *NodePredicateContext) STRING() antlr.TerminalNode {
	return s.GetToken(AqlParserSTRING, 0)
}

func (s *NodePredicateContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AqlParserPARAMETER, 0)
}

func (s *NodePredicateContext) TERM_CODE() antlr.TerminalNode {
	return s.GetToken(AqlParserTERM_CODE, 0)
}

func (s *NodePredicateContext) ARCHETYPE_HRID() antlr.TerminalNode {
	return s.GetToken(AqlParserARCHETYPE_HRID, 0)
}

func (s *NodePredicateContext) ObjectPath() IObjectPathContext {
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

func (s *NodePredicateContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(AqlParserCOMPARISON_OPERATOR, 0)
}

func (s *NodePredicateContext) PathPredicateOperand() IPathPredicateOperandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathPredicateOperandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathPredicateOperandContext)
}

func (s *NodePredicateContext) MATCHES() antlr.TerminalNode {
	return s.GetToken(AqlParserMATCHES, 0)
}

func (s *NodePredicateContext) CONTAINED_REGEX() antlr.TerminalNode {
	return s.GetToken(AqlParserCONTAINED_REGEX, 0)
}

func (s *NodePredicateContext) AllNodePredicate() []INodePredicateContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INodePredicateContext); ok {
			len++
		}
	}

	tst := make([]INodePredicateContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INodePredicateContext); ok {
			tst[i] = t.(INodePredicateContext)
			i++
		}
	}

	return tst
}

func (s *NodePredicateContext) NodePredicate(i int) INodePredicateContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INodePredicateContext); ok {
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

	return t.(INodePredicateContext)
}

func (s *NodePredicateContext) AND() antlr.TerminalNode {
	return s.GetToken(AqlParserAND, 0)
}

func (s *NodePredicateContext) OR() antlr.TerminalNode {
	return s.GetToken(AqlParserOR, 0)
}

func (s *NodePredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NodePredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NodePredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterNodePredicate(s)
	}
}

func (s *NodePredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitNodePredicate(s)
	}
}





func (p *AqlParser) NodePredicate() (localctx INodePredicateContext) {
	return p.nodePredicate(0)
}

func (p *AqlParser) nodePredicate(_p int) (localctx INodePredicateContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewNodePredicateContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx INodePredicateContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 38
	p.EnterRecursionRule(localctx, 38, AqlParserRULE_nodePredicate, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(274)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(255)
			_la = p.GetTokenStream().LA(1)

			if !(_la == AqlParserID_CODE || _la == AqlParserAT_CODE) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(258)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 30, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(256)
				p.Match(AqlParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(257)
				_la = p.GetTokenStream().LA(1)

				if !(((int64((_la - 57)) & ^0x3f) == 0 && ((int64(1) << (_la - 57)) & 32839) != 0)) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}


	case 2:
		{
			p.SetState(260)
			p.Match(AqlParserARCHETYPE_HRID)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(263)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(261)
				p.Match(AqlParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(262)
				_la = p.GetTokenStream().LA(1)

				if !(((int64((_la - 57)) & ^0x3f) == 0 && ((int64(1) << (_la - 57)) & 32839) != 0)) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}


	case 3:
		{
			p.SetState(265)
			p.Match(AqlParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 4:
		{
			p.SetState(266)
			p.ObjectPath()
		}
		{
			p.SetState(267)
			p.Match(AqlParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(268)
			p.PathPredicateOperand()
		}


	case 5:
		{
			p.SetState(270)
			p.ObjectPath()
		}
		{
			p.SetState(271)
			p.Match(AqlParserMATCHES)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(272)
			p.Match(AqlParserCONTAINED_REGEX)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(284)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(282)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext()) {
			case 1:
				localctx = NewNodePredicateContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AqlParserRULE_nodePredicate)
				p.SetState(276)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(277)
					p.Match(AqlParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(278)
					p.nodePredicate(3)
				}


			case 2:
				localctx = NewNodePredicateContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AqlParserRULE_nodePredicate)
				p.SetState(279)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(280)
					p.Match(AqlParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(281)
					p.nodePredicate(2)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(286)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext())
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


// IVersionPredicateContext is an interface to support dynamic dispatch.
type IVersionPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LATEST_VERSION() antlr.TerminalNode
	ALL_VERSIONS() antlr.TerminalNode
	StandardPredicate() IStandardPredicateContext

	// IsVersionPredicateContext differentiates from other interfaces.
	IsVersionPredicateContext()
}

type VersionPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVersionPredicateContext() *VersionPredicateContext {
	var p = new(VersionPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_versionPredicate
	return p
}

func InitEmptyVersionPredicateContext(p *VersionPredicateContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_versionPredicate
}

func (*VersionPredicateContext) IsVersionPredicateContext() {}

func NewVersionPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VersionPredicateContext {
	var p = new(VersionPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_versionPredicate

	return p
}

func (s *VersionPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *VersionPredicateContext) LATEST_VERSION() antlr.TerminalNode {
	return s.GetToken(AqlParserLATEST_VERSION, 0)
}

func (s *VersionPredicateContext) ALL_VERSIONS() antlr.TerminalNode {
	return s.GetToken(AqlParserALL_VERSIONS, 0)
}

func (s *VersionPredicateContext) StandardPredicate() IStandardPredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStandardPredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStandardPredicateContext)
}

func (s *VersionPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VersionPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *VersionPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterVersionPredicate(s)
	}
}

func (s *VersionPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitVersionPredicate(s)
	}
}




func (p *AqlParser) VersionPredicate() (localctx IVersionPredicateContext) {
	localctx = NewVersionPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, AqlParserRULE_versionPredicate)
	p.SetState(290)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserLATEST_VERSION:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(287)
			p.Match(AqlParserLATEST_VERSION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserALL_VERSIONS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(288)
			p.Match(AqlParserALL_VERSIONS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(289)
			p.StandardPredicate()
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


// IPathPredicateOperandContext is an interface to support dynamic dispatch.
type IPathPredicateOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primitive() IPrimitiveContext
	ObjectPath() IObjectPathContext
	PARAMETER() antlr.TerminalNode
	ID_CODE() antlr.TerminalNode
	AT_CODE() antlr.TerminalNode

	// IsPathPredicateOperandContext differentiates from other interfaces.
	IsPathPredicateOperandContext()
}

type PathPredicateOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathPredicateOperandContext() *PathPredicateOperandContext {
	var p = new(PathPredicateOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_pathPredicateOperand
	return p
}

func InitEmptyPathPredicateOperandContext(p *PathPredicateOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_pathPredicateOperand
}

func (*PathPredicateOperandContext) IsPathPredicateOperandContext() {}

func NewPathPredicateOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathPredicateOperandContext {
	var p = new(PathPredicateOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_pathPredicateOperand

	return p
}

func (s *PathPredicateOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *PathPredicateOperandContext) Primitive() IPrimitiveContext {
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

func (s *PathPredicateOperandContext) ObjectPath() IObjectPathContext {
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

func (s *PathPredicateOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AqlParserPARAMETER, 0)
}

func (s *PathPredicateOperandContext) ID_CODE() antlr.TerminalNode {
	return s.GetToken(AqlParserID_CODE, 0)
}

func (s *PathPredicateOperandContext) AT_CODE() antlr.TerminalNode {
	return s.GetToken(AqlParserAT_CODE, 0)
}

func (s *PathPredicateOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathPredicateOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathPredicateOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterPathPredicateOperand(s)
	}
}

func (s *PathPredicateOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitPathPredicateOperand(s)
	}
}




func (p *AqlParser) PathPredicateOperand() (localctx IPathPredicateOperandContext) {
	localctx = NewPathPredicateOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, AqlParserRULE_pathPredicateOperand)
	p.SetState(297)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserNULL, AqlParserBOOLEAN, AqlParserINTEGER, AqlParserREAL, AqlParserSCI_INTEGER, AqlParserSCI_REAL, AqlParserDATE, AqlParserTIME, AqlParserDATETIME, AqlParserSTRING, AqlParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(292)
			p.Primitive()
		}


	case AqlParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(293)
			p.ObjectPath()
		}


	case AqlParserPARAMETER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(294)
			p.Match(AqlParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserID_CODE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(295)
			p.Match(AqlParserID_CODE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserAT_CODE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(296)
			p.Match(AqlParserAT_CODE)
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
	p.RuleIndex = AqlParserRULE_objectPath
	return p
}

func InitEmptyObjectPathContext(p *ObjectPathContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_objectPath
}

func (*ObjectPathContext) IsObjectPathContext() {}

func NewObjectPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectPathContext {
	var p = new(ObjectPathContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_objectPath

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
	return s.GetTokens(AqlParserSYM_SLASH)
}

func (s *ObjectPathContext) SYM_SLASH(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_SLASH, i)
}

func (s *ObjectPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ObjectPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterObjectPath(s)
	}
}

func (s *ObjectPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitObjectPath(s)
	}
}




func (p *AqlParser) ObjectPath() (localctx IObjectPathContext) {
	localctx = NewObjectPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, AqlParserRULE_objectPath)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(299)
		p.PathPart()
	}
	p.SetState(304)
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
			{
				p.SetState(300)
				p.Match(AqlParserSYM_SLASH)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(301)
				p.PathPart()
			}


		}
		p.SetState(306)
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
	PathPredicate() IPathPredicateContext

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
	p.RuleIndex = AqlParserRULE_pathPart
	return p
}

func InitEmptyPathPartContext(p *PathPartContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_pathPart
}

func (*PathPartContext) IsPathPartContext() {}

func NewPathPartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathPartContext {
	var p = new(PathPartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_pathPart

	return p
}

func (s *PathPartContext) GetParser() antlr.Parser { return s.parser }

func (s *PathPartContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AqlParserIDENTIFIER, 0)
}

func (s *PathPartContext) PathPredicate() IPathPredicateContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathPredicateContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPathPredicateContext)
}

func (s *PathPartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathPartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathPartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterPathPart(s)
	}
}

func (s *PathPartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitPathPart(s)
	}
}




func (p *AqlParser) PathPart() (localctx IPathPartContext) {
	localctx = NewPathPartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, AqlParserRULE_pathPart)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(307)
		p.Match(AqlParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(309)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(308)
			p.PathPredicate()
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


// ILikeOperandContext is an interface to support dynamic dispatch.
type ILikeOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING() antlr.TerminalNode
	PARAMETER() antlr.TerminalNode

	// IsLikeOperandContext differentiates from other interfaces.
	IsLikeOperandContext()
}

type LikeOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLikeOperandContext() *LikeOperandContext {
	var p = new(LikeOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_likeOperand
	return p
}

func InitEmptyLikeOperandContext(p *LikeOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_likeOperand
}

func (*LikeOperandContext) IsLikeOperandContext() {}

func NewLikeOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LikeOperandContext {
	var p = new(LikeOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_likeOperand

	return p
}

func (s *LikeOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *LikeOperandContext) STRING() antlr.TerminalNode {
	return s.GetToken(AqlParserSTRING, 0)
}

func (s *LikeOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AqlParserPARAMETER, 0)
}

func (s *LikeOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LikeOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LikeOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterLikeOperand(s)
	}
}

func (s *LikeOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitLikeOperand(s)
	}
}




func (p *AqlParser) LikeOperand() (localctx ILikeOperandContext) {
	localctx = NewLikeOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, AqlParserRULE_likeOperand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(311)
		_la = p.GetTokenStream().LA(1)

		if !(_la == AqlParserPARAMETER || _la == AqlParserSTRING) {
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


// IMatchesOperandContext is an interface to support dynamic dispatch.
type IMatchesOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SYM_LEFT_CURLY() antlr.TerminalNode
	AllValueListItem() []IValueListItemContext
	ValueListItem(i int) IValueListItemContext
	SYM_RIGHT_CURLY() antlr.TerminalNode
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode
	TerminologyFunction() ITerminologyFunctionContext
	URI() antlr.TerminalNode

	// IsMatchesOperandContext differentiates from other interfaces.
	IsMatchesOperandContext()
}

type MatchesOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMatchesOperandContext() *MatchesOperandContext {
	var p = new(MatchesOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_matchesOperand
	return p
}

func InitEmptyMatchesOperandContext(p *MatchesOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_matchesOperand
}

func (*MatchesOperandContext) IsMatchesOperandContext() {}

func NewMatchesOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MatchesOperandContext {
	var p = new(MatchesOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_matchesOperand

	return p
}

func (s *MatchesOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *MatchesOperandContext) SYM_LEFT_CURLY() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_CURLY, 0)
}

func (s *MatchesOperandContext) AllValueListItem() []IValueListItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueListItemContext); ok {
			len++
		}
	}

	tst := make([]IValueListItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueListItemContext); ok {
			tst[i] = t.(IValueListItemContext)
			i++
		}
	}

	return tst
}

func (s *MatchesOperandContext) ValueListItem(i int) IValueListItemContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueListItemContext); ok {
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

	return t.(IValueListItemContext)
}

func (s *MatchesOperandContext) SYM_RIGHT_CURLY() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_CURLY, 0)
}

func (s *MatchesOperandContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AqlParserSYM_COMMA)
}

func (s *MatchesOperandContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_COMMA, i)
}

func (s *MatchesOperandContext) TerminologyFunction() ITerminologyFunctionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITerminologyFunctionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITerminologyFunctionContext)
}

func (s *MatchesOperandContext) URI() antlr.TerminalNode {
	return s.GetToken(AqlParserURI, 0)
}

func (s *MatchesOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MatchesOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *MatchesOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterMatchesOperand(s)
	}
}

func (s *MatchesOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitMatchesOperand(s)
	}
}




func (p *AqlParser) MatchesOperand() (localctx IMatchesOperandContext) {
	localctx = NewMatchesOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, AqlParserRULE_matchesOperand)
	var _la int

	p.SetState(328)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(313)
			p.Match(AqlParserSYM_LEFT_CURLY)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(314)
			p.ValueListItem()
		}
		p.SetState(319)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		for _la == AqlParserSYM_COMMA {
			{
				p.SetState(315)
				p.Match(AqlParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(316)
				p.ValueListItem()
			}


			p.SetState(321)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
		    	goto errorExit
		    }
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(322)
			p.Match(AqlParserSYM_RIGHT_CURLY)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(324)
			p.TerminologyFunction()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(325)
			p.Match(AqlParserSYM_LEFT_CURLY)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(326)
			p.Match(AqlParserURI)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(327)
			p.Match(AqlParserSYM_RIGHT_CURLY)
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


// IValueListItemContext is an interface to support dynamic dispatch.
type IValueListItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primitive() IPrimitiveContext
	PARAMETER() antlr.TerminalNode
	TerminologyFunction() ITerminologyFunctionContext

	// IsValueListItemContext differentiates from other interfaces.
	IsValueListItemContext()
}

type ValueListItemContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueListItemContext() *ValueListItemContext {
	var p = new(ValueListItemContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_valueListItem
	return p
}

func InitEmptyValueListItemContext(p *ValueListItemContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_valueListItem
}

func (*ValueListItemContext) IsValueListItemContext() {}

func NewValueListItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueListItemContext {
	var p = new(ValueListItemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_valueListItem

	return p
}

func (s *ValueListItemContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueListItemContext) Primitive() IPrimitiveContext {
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

func (s *ValueListItemContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AqlParserPARAMETER, 0)
}

func (s *ValueListItemContext) TerminologyFunction() ITerminologyFunctionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITerminologyFunctionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITerminologyFunctionContext)
}

func (s *ValueListItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueListItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ValueListItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterValueListItem(s)
	}
}

func (s *ValueListItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitValueListItem(s)
	}
}




func (p *AqlParser) ValueListItem() (localctx IValueListItemContext) {
	localctx = NewValueListItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, AqlParserRULE_valueListItem)
	p.SetState(333)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserNULL, AqlParserBOOLEAN, AqlParserINTEGER, AqlParserREAL, AqlParserSCI_INTEGER, AqlParserSCI_REAL, AqlParserDATE, AqlParserTIME, AqlParserDATETIME, AqlParserSTRING, AqlParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(330)
			p.Primitive()
		}


	case AqlParserPARAMETER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(331)
			p.Match(AqlParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserTERMINOLOGY:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(332)
			p.TerminologyFunction()
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
	NumericPrimitive() INumericPrimitiveContext
	DATE() antlr.TerminalNode
	TIME() antlr.TerminalNode
	DATETIME() antlr.TerminalNode
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
	p.RuleIndex = AqlParserRULE_primitive
	return p
}

func InitEmptyPrimitiveContext(p *PrimitiveContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_primitive
}

func (*PrimitiveContext) IsPrimitiveContext() {}

func NewPrimitiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimitiveContext {
	var p = new(PrimitiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_primitive

	return p
}

func (s *PrimitiveContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimitiveContext) STRING() antlr.TerminalNode {
	return s.GetToken(AqlParserSTRING, 0)
}

func (s *PrimitiveContext) NumericPrimitive() INumericPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumericPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumericPrimitiveContext)
}

func (s *PrimitiveContext) DATE() antlr.TerminalNode {
	return s.GetToken(AqlParserDATE, 0)
}

func (s *PrimitiveContext) TIME() antlr.TerminalNode {
	return s.GetToken(AqlParserTIME, 0)
}

func (s *PrimitiveContext) DATETIME() antlr.TerminalNode {
	return s.GetToken(AqlParserDATETIME, 0)
}

func (s *PrimitiveContext) BOOLEAN() antlr.TerminalNode {
	return s.GetToken(AqlParserBOOLEAN, 0)
}

func (s *PrimitiveContext) NULL() antlr.TerminalNode {
	return s.GetToken(AqlParserNULL, 0)
}

func (s *PrimitiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimitiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PrimitiveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterPrimitive(s)
	}
}

func (s *PrimitiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitPrimitive(s)
	}
}




func (p *AqlParser) Primitive() (localctx IPrimitiveContext) {
	localctx = NewPrimitiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, AqlParserRULE_primitive)
	p.SetState(342)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserSTRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(335)
			p.Match(AqlParserSTRING)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserINTEGER, AqlParserREAL, AqlParserSCI_INTEGER, AqlParserSCI_REAL, AqlParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(336)
			p.NumericPrimitive()
		}


	case AqlParserDATE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(337)
			p.Match(AqlParserDATE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserTIME:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(338)
			p.Match(AqlParserTIME)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserDATETIME:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(339)
			p.Match(AqlParserDATETIME)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserBOOLEAN:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(340)
			p.Match(AqlParserBOOLEAN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserNULL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(341)
			p.Match(AqlParserNULL)
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


// INumericPrimitiveContext is an interface to support dynamic dispatch.
type INumericPrimitiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER() antlr.TerminalNode
	REAL() antlr.TerminalNode
	SCI_INTEGER() antlr.TerminalNode
	SCI_REAL() antlr.TerminalNode
	SYM_MINUS() antlr.TerminalNode
	NumericPrimitive() INumericPrimitiveContext

	// IsNumericPrimitiveContext differentiates from other interfaces.
	IsNumericPrimitiveContext()
}

type NumericPrimitiveContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumericPrimitiveContext() *NumericPrimitiveContext {
	var p = new(NumericPrimitiveContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_numericPrimitive
	return p
}

func InitEmptyNumericPrimitiveContext(p *NumericPrimitiveContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_numericPrimitive
}

func (*NumericPrimitiveContext) IsNumericPrimitiveContext() {}

func NewNumericPrimitiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumericPrimitiveContext {
	var p = new(NumericPrimitiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_numericPrimitive

	return p
}

func (s *NumericPrimitiveContext) GetParser() antlr.Parser { return s.parser }

func (s *NumericPrimitiveContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(AqlParserINTEGER, 0)
}

func (s *NumericPrimitiveContext) REAL() antlr.TerminalNode {
	return s.GetToken(AqlParserREAL, 0)
}

func (s *NumericPrimitiveContext) SCI_INTEGER() antlr.TerminalNode {
	return s.GetToken(AqlParserSCI_INTEGER, 0)
}

func (s *NumericPrimitiveContext) SCI_REAL() antlr.TerminalNode {
	return s.GetToken(AqlParserSCI_REAL, 0)
}

func (s *NumericPrimitiveContext) SYM_MINUS() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_MINUS, 0)
}

func (s *NumericPrimitiveContext) NumericPrimitive() INumericPrimitiveContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumericPrimitiveContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumericPrimitiveContext)
}

func (s *NumericPrimitiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumericPrimitiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NumericPrimitiveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterNumericPrimitive(s)
	}
}

func (s *NumericPrimitiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitNumericPrimitive(s)
	}
}




func (p *AqlParser) NumericPrimitive() (localctx INumericPrimitiveContext) {
	localctx = NewNumericPrimitiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, AqlParserRULE_numericPrimitive)
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserINTEGER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(344)
			p.Match(AqlParserINTEGER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserREAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(345)
			p.Match(AqlParserREAL)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserSCI_INTEGER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(346)
			p.Match(AqlParserSCI_INTEGER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserSCI_REAL:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(347)
			p.Match(AqlParserSCI_REAL)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(348)
			p.Match(AqlParserSYM_MINUS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(349)
			p.NumericPrimitive()
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


// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name token.
	GetName() antlr.Token 


	// SetName sets the name token.
	SetName(antlr.Token) 


	// Getter signatures
	TerminologyFunction() ITerminologyFunctionContext
	SYM_LEFT_PAREN() antlr.TerminalNode
	SYM_RIGHT_PAREN() antlr.TerminalNode
	STRING_FUNCTION_ID() antlr.TerminalNode
	NUMERIC_FUNCTION_ID() antlr.TerminalNode
	DATE_TIME_FUNCTION_ID() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	AllTerminal() []ITerminalContext
	Terminal(i int) ITerminalContext
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	name antlr.Token
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) GetName() antlr.Token { return s.name }


func (s *FunctionCallContext) SetName(v antlr.Token) { s.name = v }


func (s *FunctionCallContext) TerminologyFunction() ITerminologyFunctionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITerminologyFunctionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITerminologyFunctionContext)
}

func (s *FunctionCallContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_PAREN, 0)
}

func (s *FunctionCallContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_PAREN, 0)
}

func (s *FunctionCallContext) STRING_FUNCTION_ID() antlr.TerminalNode {
	return s.GetToken(AqlParserSTRING_FUNCTION_ID, 0)
}

func (s *FunctionCallContext) NUMERIC_FUNCTION_ID() antlr.TerminalNode {
	return s.GetToken(AqlParserNUMERIC_FUNCTION_ID, 0)
}

func (s *FunctionCallContext) DATE_TIME_FUNCTION_ID() antlr.TerminalNode {
	return s.GetToken(AqlParserDATE_TIME_FUNCTION_ID, 0)
}

func (s *FunctionCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AqlParserIDENTIFIER, 0)
}

func (s *FunctionCallContext) AllTerminal() []ITerminalContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITerminalContext); ok {
			len++
		}
	}

	tst := make([]ITerminalContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITerminalContext); ok {
			tst[i] = t.(ITerminalContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) Terminal(i int) ITerminalContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITerminalContext); ok {
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

	return t.(ITerminalContext)
}

func (s *FunctionCallContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AqlParserSYM_COMMA)
}

func (s *FunctionCallContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_COMMA, i)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}




func (p *AqlParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, AqlParserRULE_functionCall)
	var _la int

	p.SetState(366)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserTERMINOLOGY:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(352)
			p.TerminologyFunction()
		}


	case AqlParserSTRING_FUNCTION_ID, AqlParserNUMERIC_FUNCTION_ID, AqlParserDATE_TIME_FUNCTION_ID, AqlParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(353)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*FunctionCallContext).name = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 4611686078556930048) != 0)) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*FunctionCallContext).name = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(354)
			p.Match(AqlParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(363)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if ((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 4827858860673859584) != 0) || ((int64((_la - 65)) & ^0x3f) == 0 && ((int64(1) << (_la - 65)) & 2097407) != 0) {
			{
				p.SetState(355)
				p.Terminal()
			}
			p.SetState(360)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			for _la == AqlParserSYM_COMMA {
				{
					p.SetState(356)
					p.Match(AqlParserSYM_COMMA)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(357)
					p.Terminal()
				}


				p.SetState(362)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
			    	goto errorExit
			    }
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(365)
			p.Match(AqlParserSYM_RIGHT_PAREN)
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
	IdentifiedPath() IIdentifiedPathContext
	SYM_ASTERISK() antlr.TerminalNode
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
	p.RuleIndex = AqlParserRULE_aggregateFunctionCall
	return p
}

func InitEmptyAggregateFunctionCallContext(p *AggregateFunctionCallContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_aggregateFunctionCall
}

func (*AggregateFunctionCallContext) IsAggregateFunctionCallContext() {}

func NewAggregateFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregateFunctionCallContext {
	var p = new(AggregateFunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_aggregateFunctionCall

	return p
}

func (s *AggregateFunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregateFunctionCallContext) GetName() antlr.Token { return s.name }


func (s *AggregateFunctionCallContext) SetName(v antlr.Token) { s.name = v }


func (s *AggregateFunctionCallContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_PAREN, 0)
}

func (s *AggregateFunctionCallContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_PAREN, 0)
}

func (s *AggregateFunctionCallContext) COUNT() antlr.TerminalNode {
	return s.GetToken(AqlParserCOUNT, 0)
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

func (s *AggregateFunctionCallContext) SYM_ASTERISK() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_ASTERISK, 0)
}

func (s *AggregateFunctionCallContext) DISTINCT() antlr.TerminalNode {
	return s.GetToken(AqlParserDISTINCT, 0)
}

func (s *AggregateFunctionCallContext) MIN() antlr.TerminalNode {
	return s.GetToken(AqlParserMIN, 0)
}

func (s *AggregateFunctionCallContext) MAX() antlr.TerminalNode {
	return s.GetToken(AqlParserMAX, 0)
}

func (s *AggregateFunctionCallContext) SUM() antlr.TerminalNode {
	return s.GetToken(AqlParserSUM, 0)
}

func (s *AggregateFunctionCallContext) AVG() antlr.TerminalNode {
	return s.GetToken(AqlParserAVG, 0)
}

func (s *AggregateFunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregateFunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *AggregateFunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterAggregateFunctionCall(s)
	}
}

func (s *AggregateFunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitAggregateFunctionCall(s)
	}
}




func (p *AqlParser) AggregateFunctionCall() (localctx IAggregateFunctionCallContext) {
	localctx = NewAggregateFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, AqlParserRULE_aggregateFunctionCall)
	var _la int

	p.SetState(383)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AqlParserCOUNT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(368)

			var _m = p.Match(AqlParserCOUNT)

			localctx.(*AggregateFunctionCallContext).name = _m
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(369)
			p.Match(AqlParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(375)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case AqlParserDISTINCT, AqlParserIDENTIFIER:
			p.SetState(371)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			if _la == AqlParserDISTINCT {
				{
					p.SetState(370)
					p.Match(AqlParserDISTINCT)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}

			}
			{
				p.SetState(373)
				p.IdentifiedPath()
			}


		case AqlParserSYM_ASTERISK:
			{
				p.SetState(374)
				p.Match(AqlParserSYM_ASTERISK)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}



		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(377)
			p.Match(AqlParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AqlParserMIN, AqlParserMAX, AqlParserSUM, AqlParserAVG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(378)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*AggregateFunctionCallContext).name = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 67553994410557440) != 0)) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*AggregateFunctionCallContext).name = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(379)
			p.Match(AqlParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(380)
			p.IdentifiedPath()
		}
		{
			p.SetState(381)
			p.Match(AqlParserSYM_RIGHT_PAREN)
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


// ITerminologyFunctionContext is an interface to support dynamic dispatch.
type ITerminologyFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TERMINOLOGY() antlr.TerminalNode
	SYM_LEFT_PAREN() antlr.TerminalNode
	AllSTRING() []antlr.TerminalNode
	STRING(i int) antlr.TerminalNode
	AllSYM_COMMA() []antlr.TerminalNode
	SYM_COMMA(i int) antlr.TerminalNode
	SYM_RIGHT_PAREN() antlr.TerminalNode

	// IsTerminologyFunctionContext differentiates from other interfaces.
	IsTerminologyFunctionContext()
}

type TerminologyFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTerminologyFunctionContext() *TerminologyFunctionContext {
	var p = new(TerminologyFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_terminologyFunction
	return p
}

func InitEmptyTerminologyFunctionContext(p *TerminologyFunctionContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_terminologyFunction
}

func (*TerminologyFunctionContext) IsTerminologyFunctionContext() {}

func NewTerminologyFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TerminologyFunctionContext {
	var p = new(TerminologyFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_terminologyFunction

	return p
}

func (s *TerminologyFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *TerminologyFunctionContext) TERMINOLOGY() antlr.TerminalNode {
	return s.GetToken(AqlParserTERMINOLOGY, 0)
}

func (s *TerminologyFunctionContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_LEFT_PAREN, 0)
}

func (s *TerminologyFunctionContext) AllSTRING() []antlr.TerminalNode {
	return s.GetTokens(AqlParserSTRING)
}

func (s *TerminologyFunctionContext) STRING(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserSTRING, i)
}

func (s *TerminologyFunctionContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AqlParserSYM_COMMA)
}

func (s *TerminologyFunctionContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_COMMA, i)
}

func (s *TerminologyFunctionContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AqlParserSYM_RIGHT_PAREN, 0)
}

func (s *TerminologyFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TerminologyFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TerminologyFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterTerminologyFunction(s)
	}
}

func (s *TerminologyFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitTerminologyFunction(s)
	}
}




func (p *AqlParser) TerminologyFunction() (localctx ITerminologyFunctionContext) {
	localctx = NewTerminologyFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, AqlParserRULE_terminologyFunction)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(385)
		p.Match(AqlParserTERMINOLOGY)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(386)
		p.Match(AqlParserSYM_LEFT_PAREN)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(387)
		p.Match(AqlParserSTRING)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(388)
		p.Match(AqlParserSYM_COMMA)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(389)
		p.Match(AqlParserSTRING)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(390)
		p.Match(AqlParserSYM_COMMA)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(391)
		p.Match(AqlParserSTRING)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(392)
		p.Match(AqlParserSYM_RIGHT_PAREN)
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


// ITopContext is an interface to support dynamic dispatch.
type ITopContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDirection returns the direction token.
	GetDirection() antlr.Token 


	// SetDirection sets the direction token.
	SetDirection(antlr.Token) 


	// Getter signatures
	TOP() antlr.TerminalNode
	INTEGER() antlr.TerminalNode
	FORWARD() antlr.TerminalNode
	BACKWARD() antlr.TerminalNode

	// IsTopContext differentiates from other interfaces.
	IsTopContext()
}

type TopContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	direction antlr.Token
}

func NewEmptyTopContext() *TopContext {
	var p = new(TopContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_top
	return p
}

func InitEmptyTopContext(p *TopContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AqlParserRULE_top
}

func (*TopContext) IsTopContext() {}

func NewTopContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopContext {
	var p = new(TopContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AqlParserRULE_top

	return p
}

func (s *TopContext) GetParser() antlr.Parser { return s.parser }

func (s *TopContext) GetDirection() antlr.Token { return s.direction }


func (s *TopContext) SetDirection(v antlr.Token) { s.direction = v }


func (s *TopContext) TOP() antlr.TerminalNode {
	return s.GetToken(AqlParserTOP, 0)
}

func (s *TopContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(AqlParserINTEGER, 0)
}

func (s *TopContext) FORWARD() antlr.TerminalNode {
	return s.GetToken(AqlParserFORWARD, 0)
}

func (s *TopContext) BACKWARD() antlr.TerminalNode {
	return s.GetToken(AqlParserBACKWARD, 0)
}

func (s *TopContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TopContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TopContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.EnterTop(s)
	}
}

func (s *TopContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AqlParserListener); ok {
		listenerT.ExitTop(s)
	}
}




func (p *AqlParser) Top() (localctx ITopContext) {
	localctx = NewTopContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, AqlParserRULE_top)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(394)
		p.Match(AqlParserTOP)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(395)
		p.Match(AqlParserINTEGER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(397)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AqlParserFORWARD || _la == AqlParserBACKWARD {
		{
			p.SetState(396)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*TopContext).direction = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == AqlParserFORWARD || _la == AqlParserBACKWARD) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*TopContext).direction = _ri
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


func (p *AqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 8:
			var t *WhereExprContext = nil
			if localctx != nil { t = localctx.(*WhereExprContext) }
			return p.WhereExpr_Sempred(t, predIndex)

	case 11:
			var t *ContainsExprContext = nil
			if localctx != nil { t = localctx.(*ContainsExprContext) }
			return p.ContainsExpr_Sempred(t, predIndex)

	case 19:
			var t *NodePredicateContext = nil
			if localctx != nil { t = localctx.(*NodePredicateContext) }
			return p.NodePredicate_Sempred(t, predIndex)


	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *AqlParser) WhereExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
			return p.Precpred(p.GetParserRuleContext(), 3)

	case 1:
			return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *AqlParser) ContainsExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
			return p.Precpred(p.GetParserRuleContext(), 3)

	case 3:
			return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *AqlParser) NodePredicate_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 4:
			return p.Precpred(p.GetParserRuleContext(), 2)

	case 5:
			return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

