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
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "';'", "'<'", "'>'", "'<='", "'>='", "'!='", "'='", 
    "'('", "')'", "','", "'/'", "'*'", "'+'", "'-'", "'['", "']'", "'{'", 
    "'}'", "'--'",
  }
  staticData.SymbolicNames = []string{
    "", "WS", "UNICODE_BOM", "COMMENT", "SELECT", "AS", "FROM", "WHERE", 
    "ORDER", "BY", "DESC", "DESCENDING", "ASC", "ASCENDING", "LIMIT", "OFFSET", 
    "DISTINCT", "LATEST_VERSION", "ALL_VERSIONS", "NULL", "BOOLEAN", "TOP", 
    "FORWARD", "BACKWARD", "CONTAINS", "AND", "OR", "NOT", "EXISTS", "COMPARISON_OPERATOR", 
    "LIKE", "MATCHES", "STRING_FUNCTION_ID", "NUMERIC_FUNCTION_ID", "DATE_TIME_FUNCTION_ID", 
    "LENGTH", "POSITION", "SUBSTRING", "CONCAT", "CONCAT_WS", "ABS", "MOD", 
    "CEIL", "FLOOR", "ROUND", "CURRENT_DATE", "CURRENT_TIME", "CURRENT_DATE_TIME", 
    "NOW", "CURRENT_TIMEZONE", "COUNT", "MIN", "MAX", "SUM", "AVG", "TERMINOLOGY", 
    "PARAMETER", "ID_CODE", "AT_CODE", "CONTAINED_REGEX", "ARCHETYPE_HRID", 
    "IDENTIFIER", "TERM_CODE", "URI", "INTEGER", "REAL", "SCI_INTEGER", 
    "SCI_REAL", "DATE", "TIME", "DATETIME", "STRING", "SYM_SEMICOLON", "SYM_LT", 
    "SYM_GT", "SYM_LE", "SYM_GE", "SYM_NE", "SYM_EQ", "SYM_LEFT_PAREN", 
    "SYM_RIGHT_PAREN", "SYM_COMMA", "SYM_SLASH", "SYM_ASTERISK", "SYM_PLUS", 
    "SYM_MINUS", "SYM_LEFT_BRACKET", "SYM_RIGHT_BRACKET", "SYM_LEFT_CURLY", 
    "SYM_RIGHT_CURLY", "SYM_DOUBLE_DASH",
  }
  staticData.RuleNames = []string{
    "query", "selectQuery", "selectClause", "fromClause", "whereClause", 
    "orderByClause", "limitClause", "selectExpr", "fromExpr", "whereExpr", 
    "orderByExpr", "columnExpr", "containsExpr", "identifiedExpr", "classExprOperand", 
    "terminal", "identifiedPath", "pathPredicate", "nodePredicate", "pathPredicateOperand", 
    "objectPath", "pathPart", "likeOperand", "matchesOperand", "valueListItem", 
    "primitive", "numericPrimitive", "functionCall", "aggregateFunctionCall", 
    "terminologyFunction", "limitOperand",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 90, 398, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 
	4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 
	10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15, 
	2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2, 
	21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26, 
	7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 1, 0, 1, 
	0, 1, 1, 1, 1, 1, 1, 3, 1, 68, 8, 1, 1, 1, 3, 1, 71, 8, 1, 1, 1, 3, 1, 
	74, 8, 1, 1, 2, 1, 2, 3, 2, 78, 8, 2, 1, 2, 1, 2, 1, 2, 5, 2, 83, 8, 2, 
	10, 2, 12, 2, 86, 9, 2, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 
	1, 5, 1, 5, 1, 5, 5, 5, 99, 8, 5, 10, 5, 12, 5, 102, 9, 5, 1, 6, 1, 6, 
	1, 6, 1, 6, 3, 6, 108, 8, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 114, 8, 6, 3, 
	6, 116, 8, 6, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 122, 8, 7, 3, 7, 124, 8, 7, 
	1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 136, 
	8, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 5, 9, 144, 8, 9, 10, 9, 12, 9, 
	147, 9, 9, 1, 10, 1, 10, 3, 10, 151, 8, 10, 1, 11, 1, 11, 1, 11, 1, 11, 
	3, 11, 157, 8, 11, 1, 12, 1, 12, 1, 12, 3, 12, 162, 8, 12, 1, 12, 1, 12, 
	3, 12, 166, 8, 12, 1, 12, 1, 12, 1, 12, 1, 12, 3, 12, 172, 8, 12, 1, 12, 
	1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 5, 12, 180, 8, 12, 10, 12, 12, 12, 183, 
	9, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 
	13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 
	1, 13, 1, 13, 3, 13, 207, 8, 13, 1, 14, 1, 14, 3, 14, 211, 8, 14, 1, 14, 
	3, 14, 214, 8, 14, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 220, 8, 15, 1, 16, 
	1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 227, 8, 16, 1, 16, 1, 16, 3, 16, 231, 
	8, 16, 1, 17, 1, 17, 1, 17, 1, 17, 3, 17, 237, 8, 17, 1, 17, 1, 17, 1, 
	17, 1, 17, 1, 17, 3, 17, 244, 8, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 
	3, 17, 251, 8, 17, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 3, 18, 258, 8, 18, 
	1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 3, 18, 266, 8, 18, 3, 18, 268, 
	8, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 3, 18, 278, 
	8, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 5, 18, 286, 8, 18, 10, 
	18, 12, 18, 289, 9, 18, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 296, 
	8, 19, 1, 20, 1, 20, 1, 20, 5, 20, 301, 8, 20, 10, 20, 12, 20, 304, 9, 
	20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 3, 21, 311, 8, 21, 1, 22, 1, 22, 
	1, 23, 1, 23, 1, 23, 1, 23, 5, 23, 319, 8, 23, 10, 23, 12, 23, 322, 9, 
	23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 3, 23, 330, 8, 23, 1, 24, 
	1, 24, 1, 24, 3, 24, 335, 8, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 
	25, 1, 25, 3, 25, 344, 8, 25, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 
	3, 26, 352, 8, 26, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 5, 27, 360, 
	8, 27, 10, 27, 12, 27, 363, 9, 27, 3, 27, 365, 8, 27, 1, 27, 3, 27, 368, 
	8, 27, 1, 28, 1, 28, 1, 28, 3, 28, 373, 8, 28, 1, 28, 1, 28, 3, 28, 377, 
	8, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 3, 28, 385, 8, 28, 1, 
	29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 30, 1, 30, 
	1, 30, 0, 3, 18, 24, 36, 31, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 
	24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 
	60, 0, 5, 2, 0, 10, 10, 12, 12, 2, 0, 56, 56, 71, 71, 1, 0, 32, 34, 1, 
	0, 51, 54, 2, 0, 56, 56, 64, 64, 445, 0, 62, 1, 0, 0, 0, 2, 64, 1, 0, 0, 
	0, 4, 75, 1, 0, 0, 0, 6, 87, 1, 0, 0, 0, 8, 90, 1, 0, 0, 0, 10, 93, 1, 
	0, 0, 0, 12, 115, 1, 0, 0, 0, 14, 123, 1, 0, 0, 0, 16, 125, 1, 0, 0, 0, 
	18, 135, 1, 0, 0, 0, 20, 148, 1, 0, 0, 0, 22, 156, 1, 0, 0, 0, 24, 171, 
	1, 0, 0, 0, 26, 206, 1, 0, 0, 0, 28, 208, 1, 0, 0, 0, 30, 219, 1, 0, 0, 
	0, 32, 221, 1, 0, 0, 0, 34, 250, 1, 0, 0, 0, 36, 277, 1, 0, 0, 0, 38, 295, 
	1, 0, 0, 0, 40, 297, 1, 0, 0, 0, 42, 305, 1, 0, 0, 0, 44, 312, 1, 0, 0, 
	0, 46, 329, 1, 0, 0, 0, 48, 334, 1, 0, 0, 0, 50, 343, 1, 0, 0, 0, 52, 351, 
	1, 0, 0, 0, 54, 367, 1, 0, 0, 0, 56, 384, 1, 0, 0, 0, 58, 386, 1, 0, 0, 
	0, 60, 395, 1, 0, 0, 0, 62, 63, 3, 2, 1, 0, 63, 1, 1, 0, 0, 0, 64, 65, 
	3, 4, 2, 0, 65, 67, 3, 6, 3, 0, 66, 68, 3, 8, 4, 0, 67, 66, 1, 0, 0, 0, 
	67, 68, 1, 0, 0, 0, 68, 70, 1, 0, 0, 0, 69, 71, 3, 10, 5, 0, 70, 69, 1, 
	0, 0, 0, 70, 71, 1, 0, 0, 0, 71, 73, 1, 0, 0, 0, 72, 74, 3, 12, 6, 0, 73, 
	72, 1, 0, 0, 0, 73, 74, 1, 0, 0, 0, 74, 3, 1, 0, 0, 0, 75, 77, 5, 4, 0, 
	0, 76, 78, 5, 16, 0, 0, 77, 76, 1, 0, 0, 0, 77, 78, 1, 0, 0, 0, 78, 79, 
	1, 0, 0, 0, 79, 84, 3, 14, 7, 0, 80, 81, 5, 81, 0, 0, 81, 83, 3, 14, 7, 
	0, 82, 80, 1, 0, 0, 0, 83, 86, 1, 0, 0, 0, 84, 82, 1, 0, 0, 0, 84, 85, 
	1, 0, 0, 0, 85, 5, 1, 0, 0, 0, 86, 84, 1, 0, 0, 0, 87, 88, 5, 6, 0, 0, 
	88, 89, 3, 16, 8, 0, 89, 7, 1, 0, 0, 0, 90, 91, 5, 7, 0, 0, 91, 92, 3, 
	18, 9, 0, 92, 9, 1, 0, 0, 0, 93, 94, 5, 8, 0, 0, 94, 95, 5, 9, 0, 0, 95, 
	100, 3, 20, 10, 0, 96, 97, 5, 81, 0, 0, 97, 99, 3, 20, 10, 0, 98, 96, 1, 
	0, 0, 0, 99, 102, 1, 0, 0, 0, 100, 98, 1, 0, 0, 0, 100, 101, 1, 0, 0, 0, 
	101, 11, 1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 103, 104, 5, 14, 0, 0, 104, 
	107, 3, 60, 30, 0, 105, 106, 5, 15, 0, 0, 106, 108, 3, 60, 30, 0, 107, 
	105, 1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108, 116, 1, 0, 0, 0, 109, 110, 
	5, 15, 0, 0, 110, 113, 3, 60, 30, 0, 111, 112, 5, 14, 0, 0, 112, 114, 3, 
	60, 30, 0, 113, 111, 1, 0, 0, 0, 113, 114, 1, 0, 0, 0, 114, 116, 1, 0, 
	0, 0, 115, 103, 1, 0, 0, 0, 115, 109, 1, 0, 0, 0, 116, 13, 1, 0, 0, 0, 
	117, 124, 5, 83, 0, 0, 118, 121, 3, 22, 11, 0, 119, 120, 5, 5, 0, 0, 120, 
	122, 5, 61, 0, 0, 121, 119, 1, 0, 0, 0, 121, 122, 1, 0, 0, 0, 122, 124, 
	1, 0, 0, 0, 123, 117, 1, 0, 0, 0, 123, 118, 1, 0, 0, 0, 124, 15, 1, 0, 
	0, 0, 125, 126, 3, 24, 12, 0, 126, 17, 1, 0, 0, 0, 127, 128, 6, 9, -1, 
	0, 128, 136, 3, 26, 13, 0, 129, 130, 5, 27, 0, 0, 130, 136, 3, 18, 9, 4, 
	131, 132, 5, 79, 0, 0, 132, 133, 3, 18, 9, 0, 133, 134, 5, 80, 0, 0, 134, 
	136, 1, 0, 0, 0, 135, 127, 1, 0, 0, 0, 135, 129, 1, 0, 0, 0, 135, 131, 
	1, 0, 0, 0, 136, 145, 1, 0, 0, 0, 137, 138, 10, 3, 0, 0, 138, 139, 5, 25, 
	0, 0, 139, 144, 3, 18, 9, 4, 140, 141, 10, 2, 0, 0, 141, 142, 5, 26, 0, 
	0, 142, 144, 3, 18, 9, 3, 143, 137, 1, 0, 0, 0, 143, 140, 1, 0, 0, 0, 144, 
	147, 1, 0, 0, 0, 145, 143, 1, 0, 0, 0, 145, 146, 1, 0, 0, 0, 146, 19, 1, 
	0, 0, 0, 147, 145, 1, 0, 0, 0, 148, 150, 3, 32, 16, 0, 149, 151, 7, 0, 
	0, 0, 150, 149, 1, 0, 0, 0, 150, 151, 1, 0, 0, 0, 151, 21, 1, 0, 0, 0, 
	152, 157, 3, 32, 16, 0, 153, 157, 3, 50, 25, 0, 154, 157, 3, 56, 28, 0, 
	155, 157, 3, 54, 27, 0, 156, 152, 1, 0, 0, 0, 156, 153, 1, 0, 0, 0, 156, 
	154, 1, 0, 0, 0, 156, 155, 1, 0, 0, 0, 157, 23, 1, 0, 0, 0, 158, 159, 6, 
	12, -1, 0, 159, 165, 3, 28, 14, 0, 160, 162, 5, 27, 0, 0, 161, 160, 1, 
	0, 0, 0, 161, 162, 1, 0, 0, 0, 162, 163, 1, 0, 0, 0, 163, 164, 5, 24, 0, 
	0, 164, 166, 3, 24, 12, 0, 165, 161, 1, 0, 0, 0, 165, 166, 1, 0, 0, 0, 
	166, 172, 1, 0, 0, 0, 167, 168, 5, 79, 0, 0, 168, 169, 3, 24, 12, 0, 169, 
	170, 5, 80, 0, 0, 170, 172, 1, 0, 0, 0, 171, 158, 1, 0, 0, 0, 171, 167, 
	1, 0, 0, 0, 172, 181, 1, 0, 0, 0, 173, 174, 10, 3, 0, 0, 174, 175, 5, 25, 
	0, 0, 175, 180, 3, 24, 12, 4, 176, 177, 10, 2, 0, 0, 177, 178, 5, 26, 0, 
	0, 178, 180, 3, 24, 12, 3, 179, 173, 1, 0, 0, 0, 179, 176, 1, 0, 0, 0, 
	180, 183, 1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 181, 182, 1, 0, 0, 0, 182, 
	25, 1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 184, 185, 5, 28, 0, 0, 185, 207, 
	3, 32, 16, 0, 186, 187, 3, 32, 16, 0, 187, 188, 5, 29, 0, 0, 188, 189, 
	3, 30, 15, 0, 189, 207, 1, 0, 0, 0, 190, 191, 3, 54, 27, 0, 191, 192, 5, 
	29, 0, 0, 192, 193, 3, 30, 15, 0, 193, 207, 1, 0, 0, 0, 194, 195, 3, 32, 
	16, 0, 195, 196, 5, 30, 0, 0, 196, 197, 3, 44, 22, 0, 197, 207, 1, 0, 0, 
	0, 198, 199, 3, 32, 16, 0, 199, 200, 5, 31, 0, 0, 200, 201, 3, 46, 23, 
	0, 201, 207, 1, 0, 0, 0, 202, 203, 5, 79, 0, 0, 203, 204, 3, 26, 13, 0, 
	204, 205, 5, 80, 0, 0, 205, 207, 1, 0, 0, 0, 206, 184, 1, 0, 0, 0, 206, 
	186, 1, 0, 0, 0, 206, 190, 1, 0, 0, 0, 206, 194, 1, 0, 0, 0, 206, 198, 
	1, 0, 0, 0, 206, 202, 1, 0, 0, 0, 207, 27, 1, 0, 0, 0, 208, 210, 5, 61, 
	0, 0, 209, 211, 5, 61, 0, 0, 210, 209, 1, 0, 0, 0, 210, 211, 1, 0, 0, 0, 
	211, 213, 1, 0, 0, 0, 212, 214, 3, 34, 17, 0, 213, 212, 1, 0, 0, 0, 213, 
	214, 1, 0, 0, 0, 214, 29, 1, 0, 0, 0, 215, 220, 3, 50, 25, 0, 216, 220, 
	5, 56, 0, 0, 217, 220, 3, 32, 16, 0, 218, 220, 3, 54, 27, 0, 219, 215, 
	1, 0, 0, 0, 219, 216, 1, 0, 0, 0, 219, 217, 1, 0, 0, 0, 219, 218, 1, 0, 
	0, 0, 220, 31, 1, 0, 0, 0, 221, 226, 5, 61, 0, 0, 222, 223, 5, 86, 0, 0, 
	223, 224, 3, 36, 18, 0, 224, 225, 5, 87, 0, 0, 225, 227, 1, 0, 0, 0, 226, 
	222, 1, 0, 0, 0, 226, 227, 1, 0, 0, 0, 227, 230, 1, 0, 0, 0, 228, 229, 
	5, 82, 0, 0, 229, 231, 3, 40, 20, 0, 230, 228, 1, 0, 0, 0, 230, 231, 1, 
	0, 0, 0, 231, 33, 1, 0, 0, 0, 232, 233, 5, 86, 0, 0, 233, 236, 5, 18, 0, 
	0, 234, 235, 5, 81, 0, 0, 235, 237, 3, 36, 18, 0, 236, 234, 1, 0, 0, 0, 
	236, 237, 1, 0, 0, 0, 237, 238, 1, 0, 0, 0, 238, 251, 5, 87, 0, 0, 239, 
	240, 5, 86, 0, 0, 240, 243, 5, 17, 0, 0, 241, 242, 5, 81, 0, 0, 242, 244, 
	3, 36, 18, 0, 243, 241, 1, 0, 0, 0, 243, 244, 1, 0, 0, 0, 244, 245, 1, 
	0, 0, 0, 245, 251, 5, 87, 0, 0, 246, 247, 5, 86, 0, 0, 247, 248, 3, 36, 
	18, 0, 248, 249, 5, 87, 0, 0, 249, 251, 1, 0, 0, 0, 250, 232, 1, 0, 0, 
	0, 250, 239, 1, 0, 0, 0, 250, 246, 1, 0, 0, 0, 251, 35, 1, 0, 0, 0, 252, 
	257, 6, 18, -1, 0, 253, 258, 5, 57, 0, 0, 254, 258, 5, 58, 0, 0, 255, 258, 
	5, 56, 0, 0, 256, 258, 5, 60, 0, 0, 257, 253, 1, 0, 0, 0, 257, 254, 1, 
	0, 0, 0, 257, 255, 1, 0, 0, 0, 257, 256, 1, 0, 0, 0, 258, 267, 1, 0, 0, 
	0, 259, 265, 5, 81, 0, 0, 260, 266, 5, 58, 0, 0, 261, 266, 5, 57, 0, 0, 
	262, 266, 5, 56, 0, 0, 263, 266, 5, 71, 0, 0, 264, 266, 5, 62, 0, 0, 265, 
	260, 1, 0, 0, 0, 265, 261, 1, 0, 0, 0, 265, 262, 1, 0, 0, 0, 265, 263, 
	1, 0, 0, 0, 265, 264, 1, 0, 0, 0, 266, 268, 1, 0, 0, 0, 267, 259, 1, 0, 
	0, 0, 267, 268, 1, 0, 0, 0, 268, 278, 1, 0, 0, 0, 269, 270, 3, 40, 20, 
	0, 270, 271, 5, 29, 0, 0, 271, 272, 3, 38, 19, 0, 272, 278, 1, 0, 0, 0, 
	273, 274, 3, 40, 20, 0, 274, 275, 5, 31, 0, 0, 275, 276, 5, 59, 0, 0, 276, 
	278, 1, 0, 0, 0, 277, 252, 1, 0, 0, 0, 277, 269, 1, 0, 0, 0, 277, 273, 
	1, 0, 0, 0, 278, 287, 1, 0, 0, 0, 279, 280, 10, 2, 0, 0, 280, 281, 5, 25, 
	0, 0, 281, 286, 3, 36, 18, 3, 282, 283, 10, 1, 0, 0, 283, 284, 5, 26, 0, 
	0, 284, 286, 3, 36, 18, 2, 285, 279, 1, 0, 0, 0, 285, 282, 1, 0, 0, 0, 
	286, 289, 1, 0, 0, 0, 287, 285, 1, 0, 0, 0, 287, 288, 1, 0, 0, 0, 288, 
	37, 1, 0, 0, 0, 289, 287, 1, 0, 0, 0, 290, 296, 3, 50, 25, 0, 291, 296, 
	3, 40, 20, 0, 292, 296, 5, 56, 0, 0, 293, 296, 5, 57, 0, 0, 294, 296, 5, 
	58, 0, 0, 295, 290, 1, 0, 0, 0, 295, 291, 1, 0, 0, 0, 295, 292, 1, 0, 0, 
	0, 295, 293, 1, 0, 0, 0, 295, 294, 1, 0, 0, 0, 296, 39, 1, 0, 0, 0, 297, 
	302, 3, 42, 21, 0, 298, 299, 5, 82, 0, 0, 299, 301, 3, 42, 21, 0, 300, 
	298, 1, 0, 0, 0, 301, 304, 1, 0, 0, 0, 302, 300, 1, 0, 0, 0, 302, 303, 
	1, 0, 0, 0, 303, 41, 1, 0, 0, 0, 304, 302, 1, 0, 0, 0, 305, 310, 5, 61, 
	0, 0, 306, 307, 5, 86, 0, 0, 307, 308, 3, 36, 18, 0, 308, 309, 5, 87, 0, 
	0, 309, 311, 1, 0, 0, 0, 310, 306, 1, 0, 0, 0, 310, 311, 1, 0, 0, 0, 311, 
	43, 1, 0, 0, 0, 312, 313, 7, 1, 0, 0, 313, 45, 1, 0, 0, 0, 314, 315, 5, 
	88, 0, 0, 315, 320, 3, 48, 24, 0, 316, 317, 5, 81, 0, 0, 317, 319, 3, 48, 
	24, 0, 318, 316, 1, 0, 0, 0, 319, 322, 1, 0, 0, 0, 320, 318, 1, 0, 0, 0, 
	320, 321, 1, 0, 0, 0, 321, 323, 1, 0, 0, 0, 322, 320, 1, 0, 0, 0, 323, 
	324, 5, 89, 0, 0, 324, 330, 1, 0, 0, 0, 325, 330, 3, 58, 29, 0, 326, 327, 
	5, 88, 0, 0, 327, 328, 5, 63, 0, 0, 328, 330, 5, 89, 0, 0, 329, 314, 1, 
	0, 0, 0, 329, 325, 1, 0, 0, 0, 329, 326, 1, 0, 0, 0, 330, 47, 1, 0, 0, 
	0, 331, 335, 3, 50, 25, 0, 332, 335, 5, 56, 0, 0, 333, 335, 3, 58, 29, 
	0, 334, 331, 1, 0, 0, 0, 334, 332, 1, 0, 0, 0, 334, 333, 1, 0, 0, 0, 335, 
	49, 1, 0, 0, 0, 336, 344, 5, 71, 0, 0, 337, 344, 3, 52, 26, 0, 338, 344, 
	5, 68, 0, 0, 339, 344, 5, 69, 0, 0, 340, 344, 5, 70, 0, 0, 341, 344, 5, 
	20, 0, 0, 342, 344, 5, 19, 0, 0, 343, 336, 1, 0, 0, 0, 343, 337, 1, 0, 
	0, 0, 343, 338, 1, 0, 0, 0, 343, 339, 1, 0, 0, 0, 343, 340, 1, 0, 0, 0, 
	343, 341, 1, 0, 0, 0, 343, 342, 1, 0, 0, 0, 344, 51, 1, 0, 0, 0, 345, 352, 
	5, 64, 0, 0, 346, 352, 5, 65, 0, 0, 347, 352, 5, 66, 0, 0, 348, 352, 5, 
	67, 0, 0, 349, 350, 5, 85, 0, 0, 350, 352, 3, 52, 26, 0, 351, 345, 1, 0, 
	0, 0, 351, 346, 1, 0, 0, 0, 351, 347, 1, 0, 0, 0, 351, 348, 1, 0, 0, 0, 
	351, 349, 1, 0, 0, 0, 352, 53, 1, 0, 0, 0, 353, 368, 3, 58, 29, 0, 354, 
	355, 7, 2, 0, 0, 355, 364, 5, 79, 0, 0, 356, 361, 3, 30, 15, 0, 357, 358, 
	5, 81, 0, 0, 358, 360, 3, 30, 15, 0, 359, 357, 1, 0, 0, 0, 360, 363, 1, 
	0, 0, 0, 361, 359, 1, 0, 0, 0, 361, 362, 1, 0, 0, 0, 362, 365, 1, 0, 0, 
	0, 363, 361, 1, 0, 0, 0, 364, 356, 1, 0, 0, 0, 364, 365, 1, 0, 0, 0, 365, 
	366, 1, 0, 0, 0, 366, 368, 5, 80, 0, 0, 367, 353, 1, 0, 0, 0, 367, 354, 
	1, 0, 0, 0, 368, 55, 1, 0, 0, 0, 369, 370, 5, 50, 0, 0, 370, 376, 5, 79, 
	0, 0, 371, 373, 5, 16, 0, 0, 372, 371, 1, 0, 0, 0, 372, 373, 1, 0, 0, 0, 
	373, 374, 1, 0, 0, 0, 374, 377, 3, 32, 16, 0, 375, 377, 5, 83, 0, 0, 376, 
	372, 1, 0, 0, 0, 376, 375, 1, 0, 0, 0, 377, 378, 1, 0, 0, 0, 378, 385, 
	5, 80, 0, 0, 379, 380, 7, 3, 0, 0, 380, 381, 5, 79, 0, 0, 381, 382, 3, 
	32, 16, 0, 382, 383, 5, 80, 0, 0, 383, 385, 1, 0, 0, 0, 384, 369, 1, 0, 
	0, 0, 384, 379, 1, 0, 0, 0, 385, 57, 1, 0, 0, 0, 386, 387, 5, 55, 0, 0, 
	387, 388, 5, 79, 0, 0, 388, 389, 5, 71, 0, 0, 389, 390, 5, 81, 0, 0, 390, 
	391, 5, 71, 0, 0, 391, 392, 5, 81, 0, 0, 392, 393, 5, 71, 0, 0, 393, 394, 
	5, 80, 0, 0, 394, 59, 1, 0, 0, 0, 395, 396, 7, 4, 0, 0, 396, 61, 1, 0, 
	0, 0, 50, 67, 70, 73, 77, 84, 100, 107, 113, 115, 121, 123, 135, 143, 145, 
	150, 156, 161, 165, 171, 179, 181, 206, 210, 213, 219, 226, 230, 236, 243, 
	250, 257, 265, 267, 277, 285, 287, 295, 302, 310, 320, 329, 334, 343, 351, 
	361, 364, 367, 372, 376, 384,
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
	AQLParserCOMMENT = 3
	AQLParserSELECT = 4
	AQLParserAS = 5
	AQLParserFROM = 6
	AQLParserWHERE = 7
	AQLParserORDER = 8
	AQLParserBY = 9
	AQLParserDESC = 10
	AQLParserDESCENDING = 11
	AQLParserASC = 12
	AQLParserASCENDING = 13
	AQLParserLIMIT = 14
	AQLParserOFFSET = 15
	AQLParserDISTINCT = 16
	AQLParserLATEST_VERSION = 17
	AQLParserALL_VERSIONS = 18
	AQLParserNULL = 19
	AQLParserBOOLEAN = 20
	AQLParserTOP = 21
	AQLParserFORWARD = 22
	AQLParserBACKWARD = 23
	AQLParserCONTAINS = 24
	AQLParserAND = 25
	AQLParserOR = 26
	AQLParserNOT = 27
	AQLParserEXISTS = 28
	AQLParserCOMPARISON_OPERATOR = 29
	AQLParserLIKE = 30
	AQLParserMATCHES = 31
	AQLParserSTRING_FUNCTION_ID = 32
	AQLParserNUMERIC_FUNCTION_ID = 33
	AQLParserDATE_TIME_FUNCTION_ID = 34
	AQLParserLENGTH = 35
	AQLParserPOSITION = 36
	AQLParserSUBSTRING = 37
	AQLParserCONCAT = 38
	AQLParserCONCAT_WS = 39
	AQLParserABS = 40
	AQLParserMOD = 41
	AQLParserCEIL = 42
	AQLParserFLOOR = 43
	AQLParserROUND = 44
	AQLParserCURRENT_DATE = 45
	AQLParserCURRENT_TIME = 46
	AQLParserCURRENT_DATE_TIME = 47
	AQLParserNOW = 48
	AQLParserCURRENT_TIMEZONE = 49
	AQLParserCOUNT = 50
	AQLParserMIN = 51
	AQLParserMAX = 52
	AQLParserSUM = 53
	AQLParserAVG = 54
	AQLParserTERMINOLOGY = 55
	AQLParserPARAMETER = 56
	AQLParserID_CODE = 57
	AQLParserAT_CODE = 58
	AQLParserCONTAINED_REGEX = 59
	AQLParserARCHETYPE_HRID = 60
	AQLParserIDENTIFIER = 61
	AQLParserTERM_CODE = 62
	AQLParserURI = 63
	AQLParserINTEGER = 64
	AQLParserREAL = 65
	AQLParserSCI_INTEGER = 66
	AQLParserSCI_REAL = 67
	AQLParserDATE = 68
	AQLParserTIME = 69
	AQLParserDATETIME = 70
	AQLParserSTRING = 71
	AQLParserSYM_SEMICOLON = 72
	AQLParserSYM_LT = 73
	AQLParserSYM_GT = 74
	AQLParserSYM_LE = 75
	AQLParserSYM_GE = 76
	AQLParserSYM_NE = 77
	AQLParserSYM_EQ = 78
	AQLParserSYM_LEFT_PAREN = 79
	AQLParserSYM_RIGHT_PAREN = 80
	AQLParserSYM_COMMA = 81
	AQLParserSYM_SLASH = 82
	AQLParserSYM_ASTERISK = 83
	AQLParserSYM_PLUS = 84
	AQLParserSYM_MINUS = 85
	AQLParserSYM_LEFT_BRACKET = 86
	AQLParserSYM_RIGHT_BRACKET = 87
	AQLParserSYM_LEFT_CURLY = 88
	AQLParserSYM_RIGHT_CURLY = 89
	AQLParserSYM_DOUBLE_DASH = 90
)

// AQLParser rules.
const (
	AQLParserRULE_query = 0
	AQLParserRULE_selectQuery = 1
	AQLParserRULE_selectClause = 2
	AQLParserRULE_fromClause = 3
	AQLParserRULE_whereClause = 4
	AQLParserRULE_orderByClause = 5
	AQLParserRULE_limitClause = 6
	AQLParserRULE_selectExpr = 7
	AQLParserRULE_fromExpr = 8
	AQLParserRULE_whereExpr = 9
	AQLParserRULE_orderByExpr = 10
	AQLParserRULE_columnExpr = 11
	AQLParserRULE_containsExpr = 12
	AQLParserRULE_identifiedExpr = 13
	AQLParserRULE_classExprOperand = 14
	AQLParserRULE_terminal = 15
	AQLParserRULE_identifiedPath = 16
	AQLParserRULE_pathPredicate = 17
	AQLParserRULE_nodePredicate = 18
	AQLParserRULE_pathPredicateOperand = 19
	AQLParserRULE_objectPath = 20
	AQLParserRULE_pathPart = 21
	AQLParserRULE_likeOperand = 22
	AQLParserRULE_matchesOperand = 23
	AQLParserRULE_valueListItem = 24
	AQLParserRULE_primitive = 25
	AQLParserRULE_numericPrimitive = 26
	AQLParserRULE_functionCall = 27
	AQLParserRULE_aggregateFunctionCall = 28
	AQLParserRULE_terminologyFunction = 29
	AQLParserRULE_limitOperand = 30
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SelectQuery() ISelectQueryContext

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
		p.SetState(62)
		p.SelectQuery()
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
	WhereClause() IWhereClauseContext
	OrderByClause() IOrderByClauseContext
	LimitClause() ILimitClauseContext

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
		p.SetState(64)
		p.SelectClause()
	}
	{
		p.SetState(65)
		p.FromClause()
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserWHERE {
		{
			p.SetState(66)
			p.WhereClause()
		}

	}
	p.SetState(70)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserORDER {
		{
			p.SetState(69)
			p.OrderByClause()
		}

	}
	p.SetState(73)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserLIMIT || _la == AQLParserOFFSET {
		{
			p.SetState(72)
			p.LimitClause()
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
		p.SetState(75)
		p.Match(AQLParserSELECT)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(77)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserDISTINCT {
		{
			p.SetState(76)
			p.Match(AQLParserDISTINCT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	}
	{
		p.SetState(79)
		p.SelectExpr()
	}
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AQLParserSYM_COMMA {
		{
			p.SetState(80)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(81)
			p.SelectExpr()
		}


		p.SetState(86)
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
		p.SetState(87)
		p.Match(AQLParserFROM)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(88)
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
	p.EnterRule(localctx, 8, AQLParserRULE_whereClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(90)
		p.Match(AQLParserWHERE)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(91)
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

func (s *OrderByClauseContext) ORDER() antlr.TerminalNode {
	return s.GetToken(AQLParserORDER, 0)
}

func (s *OrderByClauseContext) BY() antlr.TerminalNode {
	return s.GetToken(AQLParserBY, 0)
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
	p.EnterRule(localctx, 10, AQLParserRULE_orderByClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(93)
		p.Match(AQLParserORDER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(94)
		p.Match(AQLParserBY)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(95)
		p.OrderByExpr()
	}
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == AQLParserSYM_COMMA {
		{
			p.SetState(96)
			p.Match(AQLParserSYM_COMMA)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(97)
			p.OrderByExpr()
		}


		p.SetState(102)
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

	// GetLimit returns the limit rule contexts.
	GetLimit() ILimitOperandContext

	// GetOffset returns the offset rule contexts.
	GetOffset() ILimitOperandContext


	// SetLimit sets the limit rule contexts.
	SetLimit(ILimitOperandContext)

	// SetOffset sets the offset rule contexts.
	SetOffset(ILimitOperandContext)


	// Getter signatures
	LIMIT() antlr.TerminalNode
	AllLimitOperand() []ILimitOperandContext
	LimitOperand(i int) ILimitOperandContext
	OFFSET() antlr.TerminalNode

	// IsLimitClauseContext differentiates from other interfaces.
	IsLimitClauseContext()
}

type LimitClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	limit ILimitOperandContext 
	offset ILimitOperandContext 
}

func NewEmptyLimitClauseContext() *LimitClauseContext {
	var p = new(LimitClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_limitClause
	return p
}

func InitEmptyLimitClauseContext(p *LimitClauseContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_limitClause
}

func (*LimitClauseContext) IsLimitClauseContext() {}

func NewLimitClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseContext {
	var p = new(LimitClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_limitClause

	return p
}

func (s *LimitClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseContext) GetLimit() ILimitOperandContext { return s.limit }

func (s *LimitClauseContext) GetOffset() ILimitOperandContext { return s.offset }


func (s *LimitClauseContext) SetLimit(v ILimitOperandContext) { s.limit = v }

func (s *LimitClauseContext) SetOffset(v ILimitOperandContext) { s.offset = v }


func (s *LimitClauseContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(AQLParserLIMIT, 0)
}

func (s *LimitClauseContext) AllLimitOperand() []ILimitOperandContext {
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

func (s *LimitClauseContext) LimitOperand(i int) ILimitOperandContext {
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

func (s *LimitClauseContext) OFFSET() antlr.TerminalNode {
	return s.GetToken(AQLParserOFFSET, 0)
}

func (s *LimitClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LimitClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterLimitClause(s)
	}
}

func (s *LimitClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitLimitClause(s)
	}
}




func (p *AQLParser) LimitClause() (localctx ILimitClauseContext) {
	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, AQLParserRULE_limitClause)
	var _la int

	p.SetState(115)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserLIMIT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(103)
			p.Match(AQLParserLIMIT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(104)

			var _x = p.LimitOperand()


			localctx.(*LimitClauseContext).limit = _x
		}
		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserOFFSET {
			{
				p.SetState(105)
				p.Match(AQLParserOFFSET)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(106)

				var _x = p.LimitOperand()


				localctx.(*LimitClauseContext).offset = _x
			}

		}


	case AQLParserOFFSET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(109)
			p.Match(AQLParserOFFSET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(110)

			var _x = p.LimitOperand()


			localctx.(*LimitClauseContext).offset = _x
		}
		p.SetState(113)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserLIMIT {
			{
				p.SetState(111)
				p.Match(AQLParserLIMIT)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(112)

				var _x = p.LimitOperand()


				localctx.(*LimitClauseContext).limit = _x
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

	// GetAliasName returns the aliasName token.
	GetAliasName() antlr.Token 


	// SetAliasName sets the aliasName token.
	SetAliasName(antlr.Token) 


	// Getter signatures
	SYM_ASTERISK() antlr.TerminalNode
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

func (s *SelectExprContext) GetAliasName() antlr.Token { return s.aliasName }


func (s *SelectExprContext) SetAliasName(v antlr.Token) { s.aliasName = v }


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

func (s *SelectExprContext) AS() antlr.TerminalNode {
	return s.GetToken(AQLParserAS, 0)
}

func (s *SelectExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, 0)
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
	p.EnterRule(localctx, 14, AQLParserRULE_selectExpr)
	var _la int

	p.SetState(123)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserSYM_ASTERISK:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(117)
			p.Match(AQLParserSYM_ASTERISK)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserNULL, AQLParserBOOLEAN, AQLParserSTRING_FUNCTION_ID, AQLParserNUMERIC_FUNCTION_ID, AQLParserDATE_TIME_FUNCTION_ID, AQLParserCOUNT, AQLParserMIN, AQLParserMAX, AQLParserSUM, AQLParserAVG, AQLParserTERMINOLOGY, AQLParserIDENTIFIER, AQLParserINTEGER, AQLParserREAL, AQLParserSCI_INTEGER, AQLParserSCI_REAL, AQLParserDATE, AQLParserTIME, AQLParserDATETIME, AQLParserSTRING, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(118)
			p.ColumnExpr()
		}
		p.SetState(121)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserAS {
			{
				p.SetState(119)
				p.Match(AQLParserAS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(120)

				var _m = p.Match(AQLParserIDENTIFIER)

				localctx.(*SelectExprContext).aliasName = _m
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
	p.EnterRule(localctx, 16, AQLParserRULE_fromExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(125)
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
	_startState := 18
	p.EnterRecursionRule(localctx, 18, AQLParserRULE_whereExpr, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(135)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(128)
			p.IdentifiedExpr()
		}


	case 2:
		{
			p.SetState(129)
			p.Match(AQLParserNOT)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(130)
			p.whereExpr(4)
		}


	case 3:
		{
			p.SetState(131)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(132)
			p.whereExpr(0)
		}
		{
			p.SetState(133)
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
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(143)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext()) {
			case 1:
				localctx = NewWhereExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_whereExpr)
				p.SetState(137)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(138)
					p.Match(AQLParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(139)
					p.whereExpr(4)
				}


			case 2:
				localctx = NewWhereExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_whereExpr)
				p.SetState(140)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(141)
					p.Match(AQLParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(142)
					p.whereExpr(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 20, AQLParserRULE_orderByExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(148)
		p.IdentifiedPath()
	}
	p.SetState(150)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	if _la == AQLParserDESC || _la == AQLParserASC {
		{
			p.SetState(149)

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
	p.EnterRule(localctx, 22, AQLParserRULE_columnExpr)
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(152)
			p.IdentifiedPath()
		}


	case AQLParserNULL, AQLParserBOOLEAN, AQLParserINTEGER, AQLParserREAL, AQLParserSCI_INTEGER, AQLParserSCI_REAL, AQLParserDATE, AQLParserTIME, AQLParserDATETIME, AQLParserSTRING, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(153)
			p.Primitive()
		}


	case AQLParserCOUNT, AQLParserMIN, AQLParserMAX, AQLParserSUM, AQLParserAVG:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(154)
			p.AggregateFunctionCall()
		}


	case AQLParserSTRING_FUNCTION_ID, AQLParserNUMERIC_FUNCTION_ID, AQLParserDATE_TIME_FUNCTION_ID, AQLParserTERMINOLOGY:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(155)
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
	p.RuleIndex = AQLParserRULE_containsExpr
	return p
}

func InitEmptyContainsExprContext(p *ContainsExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_containsExpr
}

func (*ContainsExprContext) IsContainsExprContext() {}

func NewContainsExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContainsExprContext {
	var p = new(ContainsExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_containsExpr

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
	return s.GetToken(AQLParserCONTAINS, 0)
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
	return s.GetToken(AQLParserNOT, 0)
}

func (s *ContainsExprContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *ContainsExprContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *ContainsExprContext) AND() antlr.TerminalNode {
	return s.GetToken(AQLParserAND, 0)
}

func (s *ContainsExprContext) OR() antlr.TerminalNode {
	return s.GetToken(AQLParserOR, 0)
}

func (s *ContainsExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContainsExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ContainsExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterContainsExpr(s)
	}
}

func (s *ContainsExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitContainsExpr(s)
	}
}





func (p *AQLParser) ContainsExpr() (localctx IContainsExprContext) {
	return p.containsExpr(0)
}

func (p *AQLParser) containsExpr(_p int) (localctx IContainsExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewContainsExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IContainsExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 24
	p.EnterRecursionRule(localctx, 24, AQLParserRULE_containsExpr, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(171)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserIDENTIFIER:
		{
			p.SetState(159)
			p.ClassExprOperand()
		}
		p.SetState(165)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 17, p.GetParserRuleContext()) == 1 {
			p.SetState(161)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			if _la == AQLParserNOT {
				{
					p.SetState(160)
					p.Match(AQLParserNOT)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}

			}
			{
				p.SetState(163)
				p.Match(AQLParserCONTAINS)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(164)
				p.containsExpr(0)
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}


	case AQLParserSYM_LEFT_PAREN:
		{
			p.SetState(167)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(168)
			p.containsExpr(0)
		}
		{
			p.SetState(169)
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
	p.SetState(181)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(179)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext()) {
			case 1:
				localctx = NewContainsExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_containsExpr)
				p.SetState(173)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(174)
					p.Match(AQLParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(175)
					p.containsExpr(4)
				}


			case 2:
				localctx = NewContainsExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_containsExpr)
				p.SetState(176)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(177)
					p.Match(AQLParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(178)
					p.containsExpr(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(183)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext())
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
	p.RuleIndex = AQLParserRULE_identifiedExpr
	return p
}

func InitEmptyIdentifiedExprContext(p *IdentifiedExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_identifiedExpr
}

func (*IdentifiedExprContext) IsIdentifiedExprContext() {}

func NewIdentifiedExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifiedExprContext {
	var p = new(IdentifiedExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_identifiedExpr

	return p
}

func (s *IdentifiedExprContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifiedExprContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(AQLParserEXISTS, 0)
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
	return s.GetToken(AQLParserCOMPARISON_OPERATOR, 0)
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
	return s.GetToken(AQLParserLIKE, 0)
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
	return s.GetToken(AQLParserMATCHES, 0)
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
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
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
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *IdentifiedExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifiedExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IdentifiedExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterIdentifiedExpr(s)
	}
}

func (s *IdentifiedExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitIdentifiedExpr(s)
	}
}




func (p *AQLParser) IdentifiedExpr() (localctx IIdentifiedExprContext) {
	localctx = NewIdentifiedExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, AQLParserRULE_identifiedExpr)
	p.SetState(206)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(184)
			p.Match(AQLParserEXISTS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(185)
			p.IdentifiedPath()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(186)
			p.IdentifiedPath()
		}
		{
			p.SetState(187)
			p.Match(AQLParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(188)
			p.Terminal()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(190)
			p.FunctionCall()
		}
		{
			p.SetState(191)
			p.Match(AQLParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(192)
			p.Terminal()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(194)
			p.IdentifiedPath()
		}
		{
			p.SetState(195)
			p.Match(AQLParserLIKE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(196)
			p.LikeOperand()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(198)
			p.IdentifiedPath()
		}
		{
			p.SetState(199)
			p.Match(AQLParserMATCHES)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(200)
			p.MatchesOperand()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(202)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(203)
			p.IdentifiedExpr()
		}
		{
			p.SetState(204)
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


// IClassExprOperandContext is an interface to support dynamic dispatch.
type IClassExprOperandContext interface {
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
	PathPredicate() IPathPredicateContext

	// IsClassExprOperandContext differentiates from other interfaces.
	IsClassExprOperandContext()
}

type ClassExprOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	alias antlr.Token
}

func NewEmptyClassExprOperandContext() *ClassExprOperandContext {
	var p = new(ClassExprOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_classExprOperand
	return p
}

func InitEmptyClassExprOperandContext(p *ClassExprOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_classExprOperand
}

func (*ClassExprOperandContext) IsClassExprOperandContext() {}

func NewClassExprOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassExprOperandContext {
	var p = new(ClassExprOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_classExprOperand

	return p
}

func (s *ClassExprOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *ClassExprOperandContext) GetAlias() antlr.Token { return s.alias }


func (s *ClassExprOperandContext) SetAlias(v antlr.Token) { s.alias = v }


func (s *ClassExprOperandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(AQLParserIDENTIFIER)
}

func (s *ClassExprOperandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserIDENTIFIER, i)
}

func (s *ClassExprOperandContext) PathPredicate() IPathPredicateContext {
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

func (s *ClassExprOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClassExprOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ClassExprOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterClassExprOperand(s)
	}
}

func (s *ClassExprOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitClassExprOperand(s)
	}
}




func (p *AQLParser) ClassExprOperand() (localctx IClassExprOperandContext) {
	localctx = NewClassExprOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, AQLParserRULE_classExprOperand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(208)
		p.Match(AQLParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(210)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 22, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(209)

			var _m = p.Match(AQLParserIDENTIFIER)

			localctx.(*ClassExprOperandContext).alias = _m
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

		} else if p.HasError() { // JIM
			goto errorExit
	}
	p.SetState(213)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(212)
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
	p.RuleIndex = AQLParserRULE_terminal
	return p
}

func InitEmptyTerminalContext(p *TerminalContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_terminal
}

func (*TerminalContext) IsTerminalContext() {}

func NewTerminalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TerminalContext {
	var p = new(TerminalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_terminal

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
	return s.GetToken(AQLParserPARAMETER, 0)
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
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterTerminal(s)
	}
}

func (s *TerminalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitTerminal(s)
	}
}




func (p *AQLParser) Terminal() (localctx ITerminalContext) {
	localctx = NewTerminalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, AQLParserRULE_terminal)
	p.SetState(219)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserINTEGER, AQLParserREAL, AQLParserSCI_INTEGER, AQLParserSCI_REAL, AQLParserDATE, AQLParserTIME, AQLParserDATETIME, AQLParserSTRING, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(215)
			p.Primitive()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(216)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(217)
			p.IdentifiedPath()
		}


	case AQLParserSTRING_FUNCTION_ID, AQLParserNUMERIC_FUNCTION_ID, AQLParserDATE_TIME_FUNCTION_ID, AQLParserTERMINOLOGY:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(218)
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
	NodePredicate() INodePredicateContext
	SYM_RIGHT_BRACKET() antlr.TerminalNode
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

func (s *IdentifiedPathContext) NodePredicate() INodePredicateContext {
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
	p.EnterRule(localctx, 32, AQLParserRULE_identifiedPath)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(221)
		p.Match(AQLParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(226)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(222)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(223)
			p.nodePredicate(0)
		}
		{
			p.SetState(224)
			p.Match(AQLParserSYM_RIGHT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

		} else if p.HasError() { // JIM
			goto errorExit
	}
	p.SetState(230)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(228)
			p.Match(AQLParserSYM_SLASH)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(229)
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
	ALL_VERSIONS() antlr.TerminalNode
	SYM_RIGHT_BRACKET() antlr.TerminalNode
	SYM_COMMA() antlr.TerminalNode
	NodePredicate() INodePredicateContext
	LATEST_VERSION() antlr.TerminalNode

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
	p.RuleIndex = AQLParserRULE_pathPredicate
	return p
}

func InitEmptyPathPredicateContext(p *PathPredicateContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathPredicate
}

func (*PathPredicateContext) IsPathPredicateContext() {}

func NewPathPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathPredicateContext {
	var p = new(PathPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_pathPredicate

	return p
}

func (s *PathPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PathPredicateContext) SYM_LEFT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_BRACKET, 0)
}

func (s *PathPredicateContext) ALL_VERSIONS() antlr.TerminalNode {
	return s.GetToken(AQLParserALL_VERSIONS, 0)
}

func (s *PathPredicateContext) SYM_RIGHT_BRACKET() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_BRACKET, 0)
}

func (s *PathPredicateContext) SYM_COMMA() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, 0)
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

func (s *PathPredicateContext) LATEST_VERSION() antlr.TerminalNode {
	return s.GetToken(AQLParserLATEST_VERSION, 0)
}

func (s *PathPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterPathPredicate(s)
	}
}

func (s *PathPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitPathPredicate(s)
	}
}




func (p *AQLParser) PathPredicate() (localctx IPathPredicateContext) {
	localctx = NewPathPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, AQLParserRULE_pathPredicate)
	var _la int

	p.SetState(250)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(232)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(233)
			p.Match(AQLParserALL_VERSIONS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(236)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserSYM_COMMA {
			{
				p.SetState(234)
				p.Match(AQLParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(235)
				p.nodePredicate(0)
			}

		}
		{
			p.SetState(238)
			p.Match(AQLParserSYM_RIGHT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(239)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(240)
			p.Match(AQLParserLATEST_VERSION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(243)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if _la == AQLParserSYM_COMMA {
			{
				p.SetState(241)
				p.Match(AQLParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(242)
				p.nodePredicate(0)
			}

		}
		{
			p.SetState(245)
			p.Match(AQLParserSYM_RIGHT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(246)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(247)
			p.nodePredicate(0)
		}
		{
			p.SetState(248)
			p.Match(AQLParserSYM_RIGHT_BRACKET)
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


// INodePredicateContext is an interface to support dynamic dispatch.
type INodePredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeftIdCode returns the leftIdCode token.
	GetLeftIdCode() antlr.Token 

	// GetLeftAtCode returns the leftAtCode token.
	GetLeftAtCode() antlr.Token 

	// GetLeftParamter returns the leftParamter token.
	GetLeftParamter() antlr.Token 

	// GetRightAtCode returns the rightAtCode token.
	GetRightAtCode() antlr.Token 

	// GetRightIdCode returns the rightIdCode token.
	GetRightIdCode() antlr.Token 

	// GetRightParameter returns the rightParameter token.
	GetRightParameter() antlr.Token 


	// SetLeftIdCode sets the leftIdCode token.
	SetLeftIdCode(antlr.Token) 

	// SetLeftAtCode sets the leftAtCode token.
	SetLeftAtCode(antlr.Token) 

	// SetLeftParamter sets the leftParamter token.
	SetLeftParamter(antlr.Token) 

	// SetRightAtCode sets the rightAtCode token.
	SetRightAtCode(antlr.Token) 

	// SetRightIdCode sets the rightIdCode token.
	SetRightIdCode(antlr.Token) 

	// SetRightParameter sets the rightParameter token.
	SetRightParameter(antlr.Token) 


	// Getter signatures
	ARCHETYPE_HRID() antlr.TerminalNode
	AllID_CODE() []antlr.TerminalNode
	ID_CODE(i int) antlr.TerminalNode
	AllAT_CODE() []antlr.TerminalNode
	AT_CODE(i int) antlr.TerminalNode
	AllPARAMETER() []antlr.TerminalNode
	PARAMETER(i int) antlr.TerminalNode
	SYM_COMMA() antlr.TerminalNode
	STRING() antlr.TerminalNode
	TERM_CODE() antlr.TerminalNode
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
	leftIdCode antlr.Token
	leftAtCode antlr.Token
	leftParamter antlr.Token
	rightAtCode antlr.Token
	rightIdCode antlr.Token
	rightParameter antlr.Token
}

func NewEmptyNodePredicateContext() *NodePredicateContext {
	var p = new(NodePredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_nodePredicate
	return p
}

func InitEmptyNodePredicateContext(p *NodePredicateContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_nodePredicate
}

func (*NodePredicateContext) IsNodePredicateContext() {}

func NewNodePredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NodePredicateContext {
	var p = new(NodePredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_nodePredicate

	return p
}

func (s *NodePredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *NodePredicateContext) GetLeftIdCode() antlr.Token { return s.leftIdCode }

func (s *NodePredicateContext) GetLeftAtCode() antlr.Token { return s.leftAtCode }

func (s *NodePredicateContext) GetLeftParamter() antlr.Token { return s.leftParamter }

func (s *NodePredicateContext) GetRightAtCode() antlr.Token { return s.rightAtCode }

func (s *NodePredicateContext) GetRightIdCode() antlr.Token { return s.rightIdCode }

func (s *NodePredicateContext) GetRightParameter() antlr.Token { return s.rightParameter }


func (s *NodePredicateContext) SetLeftIdCode(v antlr.Token) { s.leftIdCode = v }

func (s *NodePredicateContext) SetLeftAtCode(v antlr.Token) { s.leftAtCode = v }

func (s *NodePredicateContext) SetLeftParamter(v antlr.Token) { s.leftParamter = v }

func (s *NodePredicateContext) SetRightAtCode(v antlr.Token) { s.rightAtCode = v }

func (s *NodePredicateContext) SetRightIdCode(v antlr.Token) { s.rightIdCode = v }

func (s *NodePredicateContext) SetRightParameter(v antlr.Token) { s.rightParameter = v }


func (s *NodePredicateContext) ARCHETYPE_HRID() antlr.TerminalNode {
	return s.GetToken(AQLParserARCHETYPE_HRID, 0)
}

func (s *NodePredicateContext) AllID_CODE() []antlr.TerminalNode {
	return s.GetTokens(AQLParserID_CODE)
}

func (s *NodePredicateContext) ID_CODE(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserID_CODE, i)
}

func (s *NodePredicateContext) AllAT_CODE() []antlr.TerminalNode {
	return s.GetTokens(AQLParserAT_CODE)
}

func (s *NodePredicateContext) AT_CODE(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserAT_CODE, i)
}

func (s *NodePredicateContext) AllPARAMETER() []antlr.TerminalNode {
	return s.GetTokens(AQLParserPARAMETER)
}

func (s *NodePredicateContext) PARAMETER(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, i)
}

func (s *NodePredicateContext) SYM_COMMA() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, 0)
}

func (s *NodePredicateContext) STRING() antlr.TerminalNode {
	return s.GetToken(AQLParserSTRING, 0)
}

func (s *NodePredicateContext) TERM_CODE() antlr.TerminalNode {
	return s.GetToken(AQLParserTERM_CODE, 0)
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
	return s.GetToken(AQLParserCOMPARISON_OPERATOR, 0)
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
	return s.GetToken(AQLParserMATCHES, 0)
}

func (s *NodePredicateContext) CONTAINED_REGEX() antlr.TerminalNode {
	return s.GetToken(AQLParserCONTAINED_REGEX, 0)
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
	return s.GetToken(AQLParserAND, 0)
}

func (s *NodePredicateContext) OR() antlr.TerminalNode {
	return s.GetToken(AQLParserOR, 0)
}

func (s *NodePredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NodePredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NodePredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterNodePredicate(s)
	}
}

func (s *NodePredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitNodePredicate(s)
	}
}





func (p *AQLParser) NodePredicate() (localctx INodePredicateContext) {
	return p.nodePredicate(0)
}

func (p *AQLParser) nodePredicate(_p int) (localctx INodePredicateContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewNodePredicateContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx INodePredicateContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 36
	p.EnterRecursionRule(localctx, 36, AQLParserRULE_nodePredicate, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(277)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext()) {
	case 1:
		p.SetState(257)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case AQLParserID_CODE:
			{
				p.SetState(253)

				var _m = p.Match(AQLParserID_CODE)

				localctx.(*NodePredicateContext).leftIdCode = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}


		case AQLParserAT_CODE:
			{
				p.SetState(254)

				var _m = p.Match(AQLParserAT_CODE)

				localctx.(*NodePredicateContext).leftAtCode = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}


		case AQLParserPARAMETER:
			{
				p.SetState(255)

				var _m = p.Match(AQLParserPARAMETER)

				localctx.(*NodePredicateContext).leftParamter = _m
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}


		case AQLParserARCHETYPE_HRID:
			{
				p.SetState(256)
				p.Match(AQLParserARCHETYPE_HRID)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}



		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		p.SetState(267)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(259)
				p.Match(AQLParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			p.SetState(265)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case AQLParserAT_CODE:
				{
					p.SetState(260)

					var _m = p.Match(AQLParserAT_CODE)

					localctx.(*NodePredicateContext).rightAtCode = _m
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}


			case AQLParserID_CODE:
				{
					p.SetState(261)

					var _m = p.Match(AQLParserID_CODE)

					localctx.(*NodePredicateContext).rightIdCode = _m
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}


			case AQLParserPARAMETER:
				{
					p.SetState(262)

					var _m = p.Match(AQLParserPARAMETER)

					localctx.(*NodePredicateContext).rightParameter = _m
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}


			case AQLParserSTRING:
				{
					p.SetState(263)
					p.Match(AQLParserSTRING)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}


			case AQLParserTERM_CODE:
				{
					p.SetState(264)
					p.Match(AQLParserTERM_CODE)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}



			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			} else if p.HasError() { // JIM
				goto errorExit
		}


	case 2:
		{
			p.SetState(269)
			p.ObjectPath()
		}
		{
			p.SetState(270)
			p.Match(AQLParserCOMPARISON_OPERATOR)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(271)
			p.PathPredicateOperand()
		}


	case 3:
		{
			p.SetState(273)
			p.ObjectPath()
		}
		{
			p.SetState(274)
			p.Match(AQLParserMATCHES)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(275)
			p.Match(AQLParserCONTAINED_REGEX)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(285)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext()) {
			case 1:
				localctx = NewNodePredicateContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_nodePredicate)
				p.SetState(279)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(280)
					p.Match(AQLParserAND)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(281)
					p.nodePredicate(3)
				}


			case 2:
				localctx = NewNodePredicateContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, AQLParserRULE_nodePredicate)
				p.SetState(282)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(283)
					p.Match(AQLParserOR)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(284)
					p.nodePredicate(2)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(289)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext())
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
	p.RuleIndex = AQLParserRULE_pathPredicateOperand
	return p
}

func InitEmptyPathPredicateOperandContext(p *PathPredicateOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_pathPredicateOperand
}

func (*PathPredicateOperandContext) IsPathPredicateOperandContext() {}

func NewPathPredicateOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathPredicateOperandContext {
	var p = new(PathPredicateOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_pathPredicateOperand

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
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *PathPredicateOperandContext) ID_CODE() antlr.TerminalNode {
	return s.GetToken(AQLParserID_CODE, 0)
}

func (s *PathPredicateOperandContext) AT_CODE() antlr.TerminalNode {
	return s.GetToken(AQLParserAT_CODE, 0)
}

func (s *PathPredicateOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathPredicateOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PathPredicateOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterPathPredicateOperand(s)
	}
}

func (s *PathPredicateOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitPathPredicateOperand(s)
	}
}




func (p *AQLParser) PathPredicateOperand() (localctx IPathPredicateOperandContext) {
	localctx = NewPathPredicateOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, AQLParserRULE_pathPredicateOperand)
	p.SetState(295)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserINTEGER, AQLParserREAL, AQLParserSCI_INTEGER, AQLParserSCI_REAL, AQLParserDATE, AQLParserTIME, AQLParserDATETIME, AQLParserSTRING, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(290)
			p.Primitive()
		}


	case AQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(291)
			p.ObjectPath()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(292)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserID_CODE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(293)
			p.Match(AQLParserID_CODE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserAT_CODE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(294)
			p.Match(AQLParserAT_CODE)
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
	p.EnterRule(localctx, 40, AQLParserRULE_objectPath)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(297)
		p.PathPart()
	}
	p.SetState(302)
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
				p.SetState(298)
				p.Match(AQLParserSYM_SLASH)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(299)
				p.PathPart()
			}


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
	NodePredicate() INodePredicateContext
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

func (s *PathPartContext) NodePredicate() INodePredicateContext {
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
	p.EnterRule(localctx, 42, AQLParserRULE_pathPart)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(305)
		p.Match(AQLParserIDENTIFIER)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(310)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(306)
			p.Match(AQLParserSYM_LEFT_BRACKET)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(307)
			p.nodePredicate(0)
		}
		{
			p.SetState(308)
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
	p.RuleIndex = AQLParserRULE_likeOperand
	return p
}

func InitEmptyLikeOperandContext(p *LikeOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_likeOperand
}

func (*LikeOperandContext) IsLikeOperandContext() {}

func NewLikeOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LikeOperandContext {
	var p = new(LikeOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_likeOperand

	return p
}

func (s *LikeOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *LikeOperandContext) STRING() antlr.TerminalNode {
	return s.GetToken(AQLParserSTRING, 0)
}

func (s *LikeOperandContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(AQLParserPARAMETER, 0)
}

func (s *LikeOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LikeOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LikeOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterLikeOperand(s)
	}
}

func (s *LikeOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitLikeOperand(s)
	}
}




func (p *AQLParser) LikeOperand() (localctx ILikeOperandContext) {
	localctx = NewLikeOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, AQLParserRULE_likeOperand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(312)
		_la = p.GetTokenStream().LA(1)

		if !(_la == AQLParserPARAMETER || _la == AQLParserSTRING) {
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
	p.RuleIndex = AQLParserRULE_matchesOperand
	return p
}

func InitEmptyMatchesOperandContext(p *MatchesOperandContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_matchesOperand
}

func (*MatchesOperandContext) IsMatchesOperandContext() {}

func NewMatchesOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MatchesOperandContext {
	var p = new(MatchesOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_matchesOperand

	return p
}

func (s *MatchesOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *MatchesOperandContext) SYM_LEFT_CURLY() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_CURLY, 0)
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
	return s.GetToken(AQLParserSYM_RIGHT_CURLY, 0)
}

func (s *MatchesOperandContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *MatchesOperandContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
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
	return s.GetToken(AQLParserURI, 0)
}

func (s *MatchesOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MatchesOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *MatchesOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterMatchesOperand(s)
	}
}

func (s *MatchesOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitMatchesOperand(s)
	}
}




func (p *AQLParser) MatchesOperand() (localctx IMatchesOperandContext) {
	localctx = NewMatchesOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, AQLParserRULE_matchesOperand)
	var _la int

	p.SetState(329)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(314)
			p.Match(AQLParserSYM_LEFT_CURLY)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(315)
			p.ValueListItem()
		}
		p.SetState(320)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		for _la == AQLParserSYM_COMMA {
			{
				p.SetState(316)
				p.Match(AQLParserSYM_COMMA)
				if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
				}
			}
			{
				p.SetState(317)
				p.ValueListItem()
			}


			p.SetState(322)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
		    	goto errorExit
		    }
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(323)
			p.Match(AQLParserSYM_RIGHT_CURLY)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(325)
			p.TerminologyFunction()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(326)
			p.Match(AQLParserSYM_LEFT_CURLY)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(327)
			p.Match(AQLParserURI)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(328)
			p.Match(AQLParserSYM_RIGHT_CURLY)
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
	p.RuleIndex = AQLParserRULE_valueListItem
	return p
}

func InitEmptyValueListItemContext(p *ValueListItemContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_valueListItem
}

func (*ValueListItemContext) IsValueListItemContext() {}

func NewValueListItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueListItemContext {
	var p = new(ValueListItemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_valueListItem

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
	return s.GetToken(AQLParserPARAMETER, 0)
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
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterValueListItem(s)
	}
}

func (s *ValueListItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitValueListItem(s)
	}
}




func (p *AQLParser) ValueListItem() (localctx IValueListItemContext) {
	localctx = NewValueListItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, AQLParserRULE_valueListItem)
	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserNULL, AQLParserBOOLEAN, AQLParserINTEGER, AQLParserREAL, AQLParserSCI_INTEGER, AQLParserSCI_REAL, AQLParserDATE, AQLParserTIME, AQLParserDATETIME, AQLParserSTRING, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(331)
			p.Primitive()
		}


	case AQLParserPARAMETER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(332)
			p.Match(AQLParserPARAMETER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserTERMINOLOGY:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(333)
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
	return s.GetToken(AQLParserDATE, 0)
}

func (s *PrimitiveContext) TIME() antlr.TerminalNode {
	return s.GetToken(AQLParserTIME, 0)
}

func (s *PrimitiveContext) DATETIME() antlr.TerminalNode {
	return s.GetToken(AQLParserDATETIME, 0)
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
	p.EnterRule(localctx, 50, AQLParserRULE_primitive)
	p.SetState(343)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserSTRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(336)
			p.Match(AQLParserSTRING)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserINTEGER, AQLParserREAL, AQLParserSCI_INTEGER, AQLParserSCI_REAL, AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(337)
			p.NumericPrimitive()
		}


	case AQLParserDATE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(338)
			p.Match(AQLParserDATE)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserTIME:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(339)
			p.Match(AQLParserTIME)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserDATETIME:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(340)
			p.Match(AQLParserDATETIME)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserBOOLEAN:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(341)
			p.Match(AQLParserBOOLEAN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserNULL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(342)
			p.Match(AQLParserNULL)
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
	p.RuleIndex = AQLParserRULE_numericPrimitive
	return p
}

func InitEmptyNumericPrimitiveContext(p *NumericPrimitiveContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_numericPrimitive
}

func (*NumericPrimitiveContext) IsNumericPrimitiveContext() {}

func NewNumericPrimitiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumericPrimitiveContext {
	var p = new(NumericPrimitiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_numericPrimitive

	return p
}

func (s *NumericPrimitiveContext) GetParser() antlr.Parser { return s.parser }

func (s *NumericPrimitiveContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(AQLParserINTEGER, 0)
}

func (s *NumericPrimitiveContext) REAL() antlr.TerminalNode {
	return s.GetToken(AQLParserREAL, 0)
}

func (s *NumericPrimitiveContext) SCI_INTEGER() antlr.TerminalNode {
	return s.GetToken(AQLParserSCI_INTEGER, 0)
}

func (s *NumericPrimitiveContext) SCI_REAL() antlr.TerminalNode {
	return s.GetToken(AQLParserSCI_REAL, 0)
}

func (s *NumericPrimitiveContext) SYM_MINUS() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_MINUS, 0)
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
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterNumericPrimitive(s)
	}
}

func (s *NumericPrimitiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitNumericPrimitive(s)
	}
}




func (p *AQLParser) NumericPrimitive() (localctx INumericPrimitiveContext) {
	localctx = NewNumericPrimitiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, AQLParserRULE_numericPrimitive)
	p.SetState(351)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserINTEGER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(345)
			p.Match(AQLParserINTEGER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserREAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(346)
			p.Match(AQLParserREAL)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserSCI_INTEGER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(347)
			p.Match(AQLParserSCI_INTEGER)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserSCI_REAL:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(348)
			p.Match(AQLParserSCI_REAL)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserSYM_MINUS:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(349)
			p.Match(AQLParserSYM_MINUS)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(350)
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
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *FunctionCallContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *FunctionCallContext) STRING_FUNCTION_ID() antlr.TerminalNode {
	return s.GetToken(AQLParserSTRING_FUNCTION_ID, 0)
}

func (s *FunctionCallContext) NUMERIC_FUNCTION_ID() antlr.TerminalNode {
	return s.GetToken(AQLParserNUMERIC_FUNCTION_ID, 0)
}

func (s *FunctionCallContext) DATE_TIME_FUNCTION_ID() antlr.TerminalNode {
	return s.GetToken(AQLParserDATE_TIME_FUNCTION_ID, 0)
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
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *FunctionCallContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
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
	p.EnterRule(localctx, 54, AQLParserRULE_functionCall)
	var _la int

	p.SetState(367)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserTERMINOLOGY:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(353)
			p.TerminologyFunction()
		}


	case AQLParserSTRING_FUNCTION_ID, AQLParserNUMERIC_FUNCTION_ID, AQLParserDATE_TIME_FUNCTION_ID:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(354)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*FunctionCallContext).name = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 30064771072) != 0)) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*FunctionCallContext).name = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(355)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(364)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)


		if ((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 2413929430336929792) != 0) || ((int64((_la - 64)) & ^0x3f) == 0 && ((int64(1) << (_la - 64)) & 2097407) != 0) {
			{
				p.SetState(356)
				p.Terminal()
			}
			p.SetState(361)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			for _la == AQLParserSYM_COMMA {
				{
					p.SetState(357)
					p.Match(AQLParserSYM_COMMA)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}
				{
					p.SetState(358)
					p.Terminal()
				}


				p.SetState(363)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
			    	goto errorExit
			    }
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(366)
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
	return s.GetToken(AQLParserSYM_ASTERISK, 0)
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
	p.EnterRule(localctx, 56, AQLParserRULE_aggregateFunctionCall)
	var _la int

	p.SetState(384)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case AQLParserCOUNT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(369)

			var _m = p.Match(AQLParserCOUNT)

			localctx.(*AggregateFunctionCallContext).name = _m
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(370)
			p.Match(AQLParserSYM_LEFT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		p.SetState(376)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case AQLParserDISTINCT, AQLParserIDENTIFIER:
			p.SetState(372)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)


			if _la == AQLParserDISTINCT {
				{
					p.SetState(371)
					p.Match(AQLParserDISTINCT)
					if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
					}
				}

			}
			{
				p.SetState(374)
				p.IdentifiedPath()
			}


		case AQLParserSYM_ASTERISK:
			{
				p.SetState(375)
				p.Match(AQLParserSYM_ASTERISK)
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
			p.SetState(378)
			p.Match(AQLParserSYM_RIGHT_PAREN)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case AQLParserMIN, AQLParserMAX, AQLParserSUM, AQLParserAVG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(379)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*AggregateFunctionCallContext).name = _lt

			_la = p.GetTokenStream().LA(1)

			if !(((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 33776997205278720) != 0)) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*AggregateFunctionCallContext).name = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
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
			p.IdentifiedPath()
		}
		{
			p.SetState(382)
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
	p.RuleIndex = AQLParserRULE_terminologyFunction
	return p
}

func InitEmptyTerminologyFunctionContext(p *TerminologyFunctionContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AQLParserRULE_terminologyFunction
}

func (*TerminologyFunctionContext) IsTerminologyFunctionContext() {}

func NewTerminologyFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TerminologyFunctionContext {
	var p = new(TerminologyFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AQLParserRULE_terminologyFunction

	return p
}

func (s *TerminologyFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *TerminologyFunctionContext) TERMINOLOGY() antlr.TerminalNode {
	return s.GetToken(AQLParserTERMINOLOGY, 0)
}

func (s *TerminologyFunctionContext) SYM_LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_LEFT_PAREN, 0)
}

func (s *TerminologyFunctionContext) AllSTRING() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSTRING)
}

func (s *TerminologyFunctionContext) STRING(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSTRING, i)
}

func (s *TerminologyFunctionContext) AllSYM_COMMA() []antlr.TerminalNode {
	return s.GetTokens(AQLParserSYM_COMMA)
}

func (s *TerminologyFunctionContext) SYM_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_COMMA, i)
}

func (s *TerminologyFunctionContext) SYM_RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(AQLParserSYM_RIGHT_PAREN, 0)
}

func (s *TerminologyFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TerminologyFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TerminologyFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.EnterTerminologyFunction(s)
	}
}

func (s *TerminologyFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AQLListener); ok {
		listenerT.ExitTerminologyFunction(s)
	}
}




func (p *AQLParser) TerminologyFunction() (localctx ITerminologyFunctionContext) {
	localctx = NewTerminologyFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, AQLParserRULE_terminologyFunction)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(386)
		p.Match(AQLParserTERMINOLOGY)
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
		p.Match(AQLParserSTRING)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
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
		p.Match(AQLParserSTRING)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
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
		p.Match(AQLParserSTRING)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(393)
		p.Match(AQLParserSYM_RIGHT_PAREN)
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
	p.EnterRule(localctx, 60, AQLParserRULE_limitOperand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(395)
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
	case 9:
			var t *WhereExprContext = nil
			if localctx != nil { t = localctx.(*WhereExprContext) }
			return p.WhereExpr_Sempred(t, predIndex)

	case 12:
			var t *ContainsExprContext = nil
			if localctx != nil { t = localctx.(*ContainsExprContext) }
			return p.ContainsExpr_Sempred(t, predIndex)

	case 18:
			var t *NodePredicateContext = nil
			if localctx != nil { t = localctx.(*NodePredicateContext) }
			return p.NodePredicate_Sempred(t, predIndex)


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

func (p *AQLParser) ContainsExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
			return p.Precpred(p.GetParserRuleContext(), 3)

	case 3:
			return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *AQLParser) NodePredicate_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 4:
			return p.Precpred(p.GetParserRuleContext(), 2)

	case 5:
			return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

