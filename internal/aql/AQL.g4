grammar AQL;

query
    : selectQuery EOF
    ;

selectQuery
    : selectClause fromClause joinClause* whereClause? groupByClause? ((UNION ALL? selectQuery) | (orderByClause? limitOffsetClause?))
    ;

selectClause
    : SELECT DISTINCT? selectExpr (SYM_COMMA selectExpr)*
    ;

fromClause
    : FROM fromExpr
    ;

joinClause
    : LEFT? JOIN joinExpr
    ;

whereClause
    : WHERE whereExpr
    ;

groupByClause
    : GROUP_BY identifiedPath (SYM_COMMA identifiedPath)*
    ;

orderByClause
    : ORDER_BY orderByExpr (SYM_COMMA orderByExpr)*
    ;

limitOffsetClause
    : LIMIT leftLimit=limitOperand (OFFSET rightOffset=limitOperand)?
    | OFFSET leftOffset=limitOperand (LIMIT rightLimit=limitOperand)?
    ;

selectExpr
    : SYM_ASTERISK
    | columnExpr (AS? IDENTIFIER)?
    ;

fromExpr
    : IDENTIFIER (SYM_LEFT_BRACKET ALL_VERSIONS SYM_RIGHT_BRACKET)? (alias=IDENTIFIER)?
    ;

joinExpr
    : IDENTIFIER (SYM_LEFT_BRACKET ALL_VERSIONS SYM_RIGHT_BRACKET)? (alias=IDENTIFIER)? ON source=IDENTIFIER
    | IDENTIFIER (alias=IDENTIFIER)? IN source=IDENTIFIER
    | IDENTIFIER (alias=IDENTIFIER)? AT identifiedPath
    ;

whereExpr
    : booleanCondition
    | NOT whereExpr
    | whereExpr AND whereExpr
    | whereExpr OR whereExpr
    | SYM_LEFT_PAREN whereExpr SYM_RIGHT_PAREN
    ;

orderByExpr
    : IDENTIFIER order=(DESC|ASC)?
    ;

columnExpr
    : primitive
    | PARAMETER
    | identifiedPath
    | aggregateFunctionCall
    | functionCall
    ;

identifiedPath
    : IDENTIFIER (SYM_LEFT_BRACKET pathCondition SYM_RIGHT_BRACKET)? (SYM_SLASH objectPath)? CAST?
    ;

objectPath
    : pathPart (SYM_SLASH pathPart)*
    ;

pathPart
    : IDENTIFIER (SYM_LEFT_BRACKET pathCondition SYM_RIGHT_BRACKET)?
    ;

pathCondition
    : pathConditionOperand COMPARISON_OPERATOR pathConditionOperand
    | pathCondition AND pathCondition
    | pathCondition OR pathCondition
    | SYM_LEFT_PAREN pathCondition SYM_RIGHT_PAREN
    ;

booleanCondition
    : identifiedPath EXISTS
    | comparisonOperand COMPARISON_OPERATOR comparisonOperand
    | comparisonOperand LIKE stringOperand
    | comparisonOperand IN SYM_LEFT_PAREN inOperand SYM_RIGHT_PAREN
    | IDENTIFIER CONTAINS IDENTIFIER (AT objectPath)?
    | SYM_LEFT_PAREN booleanCondition SYM_RIGHT_PAREN
    ;

pathConditionOperand
    : primitive
    | objectPath
    | PARAMETER
    ;

comparisonOperand
    : primitive
    | identifiedPath
    | PARAMETER
    ;

inOperand
    : selectQuery
    | inOperandValue (SYM_COMMA inOperandValue)*
    ;

inOperandValue
    : primitive
    | PARAMETER
    ;

primitive
    : STRING
    | intPrimitive
    | floatPrimitive
    | BOOLEAN
    | NULL
    ;

intPrimitive
    : SYM_MINUS? INTEGER
    | SYM_MINUS? SCI_INTEGER
    ;

floatPrimitive
    : SYM_MINUS? FLOAT
    | SYM_MINUS? SCI_FLOAT
    ;

stringOperand
    : STRING
    | PARAMETER
    | IDENTIFIER
    ;

intOperand
    : intPrimitive
    | PARAMETER
    | IDENTIFIER
    ;

numbericOperand
    : intPrimitive
    | floatPrimitive
    | PARAMETER
    | IDENTIFIER
    ;

functionCall
    : LENGTH SYM_LEFT_PAREN stringOperand SYM_RIGHT_PAREN
    | POSITION SYM_LEFT_PAREN stringOperand SYM_COMMA intOperand SYM_RIGHT_PAREN
    | SUBSTRING SYM_LEFT_PAREN stringOperand SYM_COMMA intOperand SYM_COMMA intOperand SYM_RIGHT_PAREN
    | CONCAT SYM_LEFT_PAREN stringOperand (SYM_COMMA stringOperand)* SYM_RIGHT_PAREN
    | CONCAT_WS SYM_LEFT_PAREN stringOperand SYM_COMMA stringOperand (SYM_COMMA stringOperand)* SYM_RIGHT_PAREN
    | ABS SYM_LEFT_PAREN numbericOperand SYM_RIGHT_PAREN
    | MOD SYM_LEFT_PAREN numbericOperand SYM_COMMA numbericOperand SYM_RIGHT_PAREN
    | CEIL SYM_LEFT_PAREN numbericOperand SYM_RIGHT_PAREN
    | FLOOR SYM_LEFT_PAREN numbericOperand SYM_RIGHT_PAREN
    | ROUND SYM_LEFT_PAREN numbericOperand SYM_COMMA intOperand SYM_RIGHT_PAREN
    | CURRENT_DATE
    | CURRENT_TIME
    | CURRENT_DATE_TIME
    ;

aggregateFunctionCall
    : name=COUNT SYM_LEFT_PAREN (DISTINCT? SYM_ASTERISK | identifiedPath) SYM_RIGHT_PAREN
    | name=(MIN | MAX | SUM | AVG) SYM_LEFT_PAREN identifiedPath SYM_RIGHT_PAREN
    ;

limitOperand
    : INTEGER
    | PARAMETER
    ;

// SKIP
WS: [ \t\r\n]+ -> skip;
UNICODE_BOM: (
    '\uEFBBBF' // UTF-8 BOM
    | '\uFEFF' // UTF16_BOM
    | '\u0000FEFF' // UTF32_BOM
    ) -> skip;

// Keywords
// Common Keywords
SELECT: S E L E C T ;
AS: A S ;
FROM: F R O M ;
WHERE: W H E R E ;
fragment ORDER: O R D E R ;
fragment GROUP: G R O U P ;
fragment BY: B Y ;
ORDER_BY: ORDER BY ;
GROUP_BY: GROUP BY ;

DESC: D E S C ;
ASC: A S C ;
LIMIT: L I M I T ;
OFFSET: O F F S E T ;
// other keywords
DISTINCT: D I S T I N C T ;
NULL: N U L L ;
TRUE: T R U E ;
FALSE: F A L S E ;
BOOLEAN: TRUE | FALSE ;
// unique keywords
ALL_VERSIONS : A L L '_' V E R S I O N S ;

// Operators
// Containment operator
CONTAINS : C O N T A I N S ;
// Logical operators
AND : A N D ;
OR : O R ;
NOT : N O T ;
EXISTS: E X I S T S ;
// Comparison operators
COMPARISON_OPERATOR: SYM_EQ | SYM_NE | SYM_GT | SYM_GE | SYM_LT | SYM_LE ;
LIKE: L I K E ;
IN: I N ;
ON: O N ;
AT: A T ;
// Join operators
JOIN: J O I N ;
LEFT: L E F T ;
UNION: U N I O N ;
ALL: A L L ;

// string functions
LENGTH: L E N G T H ;
POSITION: P O S I T I O N ;
SUBSTRING: S U B S T R I N G ;
CONCAT: C O N C A T ;
CONCAT_WS: C O N C A T '_' W S ;
// numeric functions
ABS: A B S ;
MOD: M O D ;
CEIL: C E I L ;
FLOOR: F L O O R ;
ROUND: R O U N D ;
// date and time functions
CURRENT_DATE: C U R R E N T '_' D A T E ;
CURRENT_TIME: C U R R E N T '_' T I M E ;
CURRENT_DATE_TIME: C U R R E N T '_' D A T E '_' T I M E ;
CURRENT_TIMEZONE: C U R R E N T '_' T I M E Z O N E ;
// aggregate function
COUNT: C O U N T ;
MIN: M I N ;
MAX: M A X ;
SUM: S U M ;
AVG: A V G ;

// other, identifiers
PARAMETER: '$' IDENTIFIER_CHAR ;
CAST: SYM_COLON SYM_COLON IDENTIFIER_CHAR ;

// ---------------------- Identifiers ---------------------

IDENTIFIER: IDENTIFIER_CHAR;
fragment IDENTIFIER_CHAR : ALPHA_CHAR WORD_CHAR* ;

// --------------------- atomic primitive types -------------------

STRING: SYM_SINGLE_QUOTE ( ESCAPE_SEQ | UTF8CHAR | OCTAL_ESC | ~('\\'|'\'') )* SYM_SINGLE_QUOTE ;
INTEGER: DIGIT+ ;
FLOAT: DIGIT* '.' DIGIT+ ;
SCI_INTEGER: INTEGER E_SUFFIX ;
SCI_FLOAT: FLOAT E_SUFFIX ;
fragment E_SUFFIX: E [-+]? DIGIT+ ;
fragment ESCAPE_SEQ: '\\' ['"*?abfnrtv\\] ;

// ------------------- character fragments ------------------

fragment NAME_CHAR: WORD_CHAR | '-' ;
fragment WORD_CHAR: ALPHANUM_CHAR | '_' ;
fragment ALPHANUM_CHAR: ALPHA_CHAR | DIGIT ;

fragment ALPHA_CHAR: [a-zA-Z];
fragment UTF8CHAR: '\\u' HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT ;

fragment DIGIT: [0-9];
fragment HEX_DIGIT: [0-9a-fA-F];

fragment OCTAL_ESC: '\\' [0-3] OCTAL_DIGIT OCTAL_DIGIT | '\\' OCTAL_DIGIT OCTAL_DIGIT | '\\' OCTAL_DIGIT;
fragment OCTAL_DIGIT: [0-7];

// ---------- symbols ----------

SYM_COLON: ':' ;
SYM_SEMICOLON: ';' ;
SYM_LT: '<' ;
SYM_GT: '>' ;
SYM_LE: '<=' ;
SYM_GE: '>=' ;
SYM_NE: '!=' ;
SYM_EQ: '=' ;
SYM_LEFT_PAREN: '(' ;
SYM_RIGHT_PAREN: ')' ;
SYM_COMMA: ',';

SYM_SLASH: '/';
SYM_ASTERISK: '*';
SYM_PLUS: '+';
SYM_MINUS: '-';

SYM_LEFT_BRACKET: '[';
SYM_RIGHT_BRACKET: ']';
SYM_LEFT_CURLY: '{';
SYM_RIGHT_CURLY: '}';
SYM_DOUBLE_DASH: '--';

fragment SYM_SINGLE_QUOTE: '\'';
fragment SYM_DOUBLE_QUOTE: '"';

// ------------------- Fragment letters ---------------------
fragment A: [aA];
fragment B: [bB];
fragment C: [cC];
fragment D: [dD];
fragment E: [eE];
fragment F: [fF];
fragment G: [gG];
fragment H: [hH];
fragment I: [iI];
fragment J: [jJ];
fragment K: [kK];
fragment L: [lL];
fragment M: [mM];
fragment N: [nN];
fragment O: [oO];
fragment P: [pP];
fragment Q: [qQ];
fragment R: [rR];
fragment S: [sS];
fragment T: [tT];
fragment U: [uU];
fragment V: [vV];
fragment W: [wW];
fragment X: [xX];
fragment Y: [yY];
fragment Z: [zZ];
