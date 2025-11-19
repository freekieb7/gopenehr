// Code generated from AQL.g4 by ANTLR 4.13.2. DO NOT EDIT.

package gen // AQL
import "github.com/antlr4-go/antlr/v4"

// BaseAQLListener is a complete listener for a parse tree produced by AQLParser.
type BaseAQLListener struct{}

var _ AQLListener = &BaseAQLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseAQLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseAQLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseAQLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseAQLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseAQLListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseAQLListener) ExitQuery(ctx *QueryContext) {}

// EnterSelectQuery is called when production selectQuery is entered.
func (s *BaseAQLListener) EnterSelectQuery(ctx *SelectQueryContext) {}

// ExitSelectQuery is called when production selectQuery is exited.
func (s *BaseAQLListener) ExitSelectQuery(ctx *SelectQueryContext) {}

// EnterSelectClause is called when production selectClause is entered.
func (s *BaseAQLListener) EnterSelectClause(ctx *SelectClauseContext) {}

// ExitSelectClause is called when production selectClause is exited.
func (s *BaseAQLListener) ExitSelectClause(ctx *SelectClauseContext) {}

// EnterFromClause is called when production fromClause is entered.
func (s *BaseAQLListener) EnterFromClause(ctx *FromClauseContext) {}

// ExitFromClause is called when production fromClause is exited.
func (s *BaseAQLListener) ExitFromClause(ctx *FromClauseContext) {}

// EnterWhereClause is called when production whereClause is entered.
func (s *BaseAQLListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production whereClause is exited.
func (s *BaseAQLListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (s *BaseAQLListener) EnterOrderByClause(ctx *OrderByClauseContext) {}

// ExitOrderByClause is called when production orderByClause is exited.
func (s *BaseAQLListener) ExitOrderByClause(ctx *OrderByClauseContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (s *BaseAQLListener) EnterLimitClause(ctx *LimitClauseContext) {}

// ExitLimitClause is called when production limitClause is exited.
func (s *BaseAQLListener) ExitLimitClause(ctx *LimitClauseContext) {}

// EnterSelectExpr is called when production selectExpr is entered.
func (s *BaseAQLListener) EnterSelectExpr(ctx *SelectExprContext) {}

// ExitSelectExpr is called when production selectExpr is exited.
func (s *BaseAQLListener) ExitSelectExpr(ctx *SelectExprContext) {}

// EnterFromExpr is called when production fromExpr is entered.
func (s *BaseAQLListener) EnterFromExpr(ctx *FromExprContext) {}

// ExitFromExpr is called when production fromExpr is exited.
func (s *BaseAQLListener) ExitFromExpr(ctx *FromExprContext) {}

// EnterWhereExpr is called when production whereExpr is entered.
func (s *BaseAQLListener) EnterWhereExpr(ctx *WhereExprContext) {}

// ExitWhereExpr is called when production whereExpr is exited.
func (s *BaseAQLListener) ExitWhereExpr(ctx *WhereExprContext) {}

// EnterOrderByExpr is called when production orderByExpr is entered.
func (s *BaseAQLListener) EnterOrderByExpr(ctx *OrderByExprContext) {}

// ExitOrderByExpr is called when production orderByExpr is exited.
func (s *BaseAQLListener) ExitOrderByExpr(ctx *OrderByExprContext) {}

// EnterColumnExpr is called when production columnExpr is entered.
func (s *BaseAQLListener) EnterColumnExpr(ctx *ColumnExprContext) {}

// ExitColumnExpr is called when production columnExpr is exited.
func (s *BaseAQLListener) ExitColumnExpr(ctx *ColumnExprContext) {}

// EnterContainsExpr is called when production containsExpr is entered.
func (s *BaseAQLListener) EnterContainsExpr(ctx *ContainsExprContext) {}

// ExitContainsExpr is called when production containsExpr is exited.
func (s *BaseAQLListener) ExitContainsExpr(ctx *ContainsExprContext) {}

// EnterIdentifiedExpr is called when production identifiedExpr is entered.
func (s *BaseAQLListener) EnterIdentifiedExpr(ctx *IdentifiedExprContext) {}

// ExitIdentifiedExpr is called when production identifiedExpr is exited.
func (s *BaseAQLListener) ExitIdentifiedExpr(ctx *IdentifiedExprContext) {}

// EnterClassExprOperand is called when production classExprOperand is entered.
func (s *BaseAQLListener) EnterClassExprOperand(ctx *ClassExprOperandContext) {}

// ExitClassExprOperand is called when production classExprOperand is exited.
func (s *BaseAQLListener) ExitClassExprOperand(ctx *ClassExprOperandContext) {}

// EnterTerminal is called when production terminal is entered.
func (s *BaseAQLListener) EnterTerminal(ctx *TerminalContext) {}

// ExitTerminal is called when production terminal is exited.
func (s *BaseAQLListener) ExitTerminal(ctx *TerminalContext) {}

// EnterIdentifiedPath is called when production identifiedPath is entered.
func (s *BaseAQLListener) EnterIdentifiedPath(ctx *IdentifiedPathContext) {}

// ExitIdentifiedPath is called when production identifiedPath is exited.
func (s *BaseAQLListener) ExitIdentifiedPath(ctx *IdentifiedPathContext) {}

// EnterPathPredicate is called when production pathPredicate is entered.
func (s *BaseAQLListener) EnterPathPredicate(ctx *PathPredicateContext) {}

// ExitPathPredicate is called when production pathPredicate is exited.
func (s *BaseAQLListener) ExitPathPredicate(ctx *PathPredicateContext) {}

// EnterNodePredicate is called when production nodePredicate is entered.
func (s *BaseAQLListener) EnterNodePredicate(ctx *NodePredicateContext) {}

// ExitNodePredicate is called when production nodePredicate is exited.
func (s *BaseAQLListener) ExitNodePredicate(ctx *NodePredicateContext) {}

// EnterPathPredicateOperand is called when production pathPredicateOperand is entered.
func (s *BaseAQLListener) EnterPathPredicateOperand(ctx *PathPredicateOperandContext) {}

// ExitPathPredicateOperand is called when production pathPredicateOperand is exited.
func (s *BaseAQLListener) ExitPathPredicateOperand(ctx *PathPredicateOperandContext) {}

// EnterObjectPath is called when production objectPath is entered.
func (s *BaseAQLListener) EnterObjectPath(ctx *ObjectPathContext) {}

// ExitObjectPath is called when production objectPath is exited.
func (s *BaseAQLListener) ExitObjectPath(ctx *ObjectPathContext) {}

// EnterPathPart is called when production pathPart is entered.
func (s *BaseAQLListener) EnterPathPart(ctx *PathPartContext) {}

// ExitPathPart is called when production pathPart is exited.
func (s *BaseAQLListener) ExitPathPart(ctx *PathPartContext) {}

// EnterLikeOperand is called when production likeOperand is entered.
func (s *BaseAQLListener) EnterLikeOperand(ctx *LikeOperandContext) {}

// ExitLikeOperand is called when production likeOperand is exited.
func (s *BaseAQLListener) ExitLikeOperand(ctx *LikeOperandContext) {}

// EnterMatchesOperand is called when production matchesOperand is entered.
func (s *BaseAQLListener) EnterMatchesOperand(ctx *MatchesOperandContext) {}

// ExitMatchesOperand is called when production matchesOperand is exited.
func (s *BaseAQLListener) ExitMatchesOperand(ctx *MatchesOperandContext) {}

// EnterValueListItem is called when production valueListItem is entered.
func (s *BaseAQLListener) EnterValueListItem(ctx *ValueListItemContext) {}

// ExitValueListItem is called when production valueListItem is exited.
func (s *BaseAQLListener) ExitValueListItem(ctx *ValueListItemContext) {}

// EnterPrimitive is called when production primitive is entered.
func (s *BaseAQLListener) EnterPrimitive(ctx *PrimitiveContext) {}

// ExitPrimitive is called when production primitive is exited.
func (s *BaseAQLListener) ExitPrimitive(ctx *PrimitiveContext) {}

// EnterNumericPrimitive is called when production numericPrimitive is entered.
func (s *BaseAQLListener) EnterNumericPrimitive(ctx *NumericPrimitiveContext) {}

// ExitNumericPrimitive is called when production numericPrimitive is exited.
func (s *BaseAQLListener) ExitNumericPrimitive(ctx *NumericPrimitiveContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseAQLListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseAQLListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterAggregateFunctionCall is called when production aggregateFunctionCall is entered.
func (s *BaseAQLListener) EnterAggregateFunctionCall(ctx *AggregateFunctionCallContext) {}

// ExitAggregateFunctionCall is called when production aggregateFunctionCall is exited.
func (s *BaseAQLListener) ExitAggregateFunctionCall(ctx *AggregateFunctionCallContext) {}

// EnterTerminologyFunction is called when production terminologyFunction is entered.
func (s *BaseAQLListener) EnterTerminologyFunction(ctx *TerminologyFunctionContext) {}

// ExitTerminologyFunction is called when production terminologyFunction is exited.
func (s *BaseAQLListener) ExitTerminologyFunction(ctx *TerminologyFunctionContext) {}

// EnterLimitOperand is called when production limitOperand is entered.
func (s *BaseAQLListener) EnterLimitOperand(ctx *LimitOperandContext) {}

// ExitLimitOperand is called when production limitOperand is exited.
func (s *BaseAQLListener) ExitLimitOperand(ctx *LimitOperandContext) {}
