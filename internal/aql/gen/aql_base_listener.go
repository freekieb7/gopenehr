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

// EnterJoinClause is called when production joinClause is entered.
func (s *BaseAQLListener) EnterJoinClause(ctx *JoinClauseContext) {}

// ExitJoinClause is called when production joinClause is exited.
func (s *BaseAQLListener) ExitJoinClause(ctx *JoinClauseContext) {}

// EnterWhereClause is called when production whereClause is entered.
func (s *BaseAQLListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production whereClause is exited.
func (s *BaseAQLListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterGroupByClause is called when production groupByClause is entered.
func (s *BaseAQLListener) EnterGroupByClause(ctx *GroupByClauseContext) {}

// ExitGroupByClause is called when production groupByClause is exited.
func (s *BaseAQLListener) ExitGroupByClause(ctx *GroupByClauseContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (s *BaseAQLListener) EnterOrderByClause(ctx *OrderByClauseContext) {}

// ExitOrderByClause is called when production orderByClause is exited.
func (s *BaseAQLListener) ExitOrderByClause(ctx *OrderByClauseContext) {}

// EnterLimitOffsetClause is called when production limitOffsetClause is entered.
func (s *BaseAQLListener) EnterLimitOffsetClause(ctx *LimitOffsetClauseContext) {}

// ExitLimitOffsetClause is called when production limitOffsetClause is exited.
func (s *BaseAQLListener) ExitLimitOffsetClause(ctx *LimitOffsetClauseContext) {}

// EnterSelectExpr is called when production selectExpr is entered.
func (s *BaseAQLListener) EnterSelectExpr(ctx *SelectExprContext) {}

// ExitSelectExpr is called when production selectExpr is exited.
func (s *BaseAQLListener) ExitSelectExpr(ctx *SelectExprContext) {}

// EnterFromExpr is called when production fromExpr is entered.
func (s *BaseAQLListener) EnterFromExpr(ctx *FromExprContext) {}

// ExitFromExpr is called when production fromExpr is exited.
func (s *BaseAQLListener) ExitFromExpr(ctx *FromExprContext) {}

// EnterJoinExpr is called when production joinExpr is entered.
func (s *BaseAQLListener) EnterJoinExpr(ctx *JoinExprContext) {}

// ExitJoinExpr is called when production joinExpr is exited.
func (s *BaseAQLListener) ExitJoinExpr(ctx *JoinExprContext) {}

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

// EnterIdentifiedPath is called when production identifiedPath is entered.
func (s *BaseAQLListener) EnterIdentifiedPath(ctx *IdentifiedPathContext) {}

// ExitIdentifiedPath is called when production identifiedPath is exited.
func (s *BaseAQLListener) ExitIdentifiedPath(ctx *IdentifiedPathContext) {}

// EnterObjectPath is called when production objectPath is entered.
func (s *BaseAQLListener) EnterObjectPath(ctx *ObjectPathContext) {}

// ExitObjectPath is called when production objectPath is exited.
func (s *BaseAQLListener) ExitObjectPath(ctx *ObjectPathContext) {}

// EnterPathPart is called when production pathPart is entered.
func (s *BaseAQLListener) EnterPathPart(ctx *PathPartContext) {}

// ExitPathPart is called when production pathPart is exited.
func (s *BaseAQLListener) ExitPathPart(ctx *PathPartContext) {}

// EnterPathCondition is called when production pathCondition is entered.
func (s *BaseAQLListener) EnterPathCondition(ctx *PathConditionContext) {}

// ExitPathCondition is called when production pathCondition is exited.
func (s *BaseAQLListener) ExitPathCondition(ctx *PathConditionContext) {}

// EnterBooleanCondition is called when production booleanCondition is entered.
func (s *BaseAQLListener) EnterBooleanCondition(ctx *BooleanConditionContext) {}

// ExitBooleanCondition is called when production booleanCondition is exited.
func (s *BaseAQLListener) ExitBooleanCondition(ctx *BooleanConditionContext) {}

// EnterPathConditionOperand is called when production pathConditionOperand is entered.
func (s *BaseAQLListener) EnterPathConditionOperand(ctx *PathConditionOperandContext) {}

// ExitPathConditionOperand is called when production pathConditionOperand is exited.
func (s *BaseAQLListener) ExitPathConditionOperand(ctx *PathConditionOperandContext) {}

// EnterComparisonOperand is called when production comparisonOperand is entered.
func (s *BaseAQLListener) EnterComparisonOperand(ctx *ComparisonOperandContext) {}

// ExitComparisonOperand is called when production comparisonOperand is exited.
func (s *BaseAQLListener) ExitComparisonOperand(ctx *ComparisonOperandContext) {}

// EnterInOperand is called when production inOperand is entered.
func (s *BaseAQLListener) EnterInOperand(ctx *InOperandContext) {}

// ExitInOperand is called when production inOperand is exited.
func (s *BaseAQLListener) ExitInOperand(ctx *InOperandContext) {}

// EnterInOperandValue is called when production inOperandValue is entered.
func (s *BaseAQLListener) EnterInOperandValue(ctx *InOperandValueContext) {}

// ExitInOperandValue is called when production inOperandValue is exited.
func (s *BaseAQLListener) ExitInOperandValue(ctx *InOperandValueContext) {}

// EnterPrimitive is called when production primitive is entered.
func (s *BaseAQLListener) EnterPrimitive(ctx *PrimitiveContext) {}

// ExitPrimitive is called when production primitive is exited.
func (s *BaseAQLListener) ExitPrimitive(ctx *PrimitiveContext) {}

// EnterIntPrimitive is called when production intPrimitive is entered.
func (s *BaseAQLListener) EnterIntPrimitive(ctx *IntPrimitiveContext) {}

// ExitIntPrimitive is called when production intPrimitive is exited.
func (s *BaseAQLListener) ExitIntPrimitive(ctx *IntPrimitiveContext) {}

// EnterFloatPrimitive is called when production floatPrimitive is entered.
func (s *BaseAQLListener) EnterFloatPrimitive(ctx *FloatPrimitiveContext) {}

// ExitFloatPrimitive is called when production floatPrimitive is exited.
func (s *BaseAQLListener) ExitFloatPrimitive(ctx *FloatPrimitiveContext) {}

// EnterStringOperand is called when production stringOperand is entered.
func (s *BaseAQLListener) EnterStringOperand(ctx *StringOperandContext) {}

// ExitStringOperand is called when production stringOperand is exited.
func (s *BaseAQLListener) ExitStringOperand(ctx *StringOperandContext) {}

// EnterIntOperand is called when production intOperand is entered.
func (s *BaseAQLListener) EnterIntOperand(ctx *IntOperandContext) {}

// ExitIntOperand is called when production intOperand is exited.
func (s *BaseAQLListener) ExitIntOperand(ctx *IntOperandContext) {}

// EnterNumbericOperand is called when production numbericOperand is entered.
func (s *BaseAQLListener) EnterNumbericOperand(ctx *NumbericOperandContext) {}

// ExitNumbericOperand is called when production numbericOperand is exited.
func (s *BaseAQLListener) ExitNumbericOperand(ctx *NumbericOperandContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseAQLListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseAQLListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterAggregateFunctionCall is called when production aggregateFunctionCall is entered.
func (s *BaseAQLListener) EnterAggregateFunctionCall(ctx *AggregateFunctionCallContext) {}

// ExitAggregateFunctionCall is called when production aggregateFunctionCall is exited.
func (s *BaseAQLListener) ExitAggregateFunctionCall(ctx *AggregateFunctionCallContext) {}

// EnterLimitOperand is called when production limitOperand is entered.
func (s *BaseAQLListener) EnterLimitOperand(ctx *LimitOperandContext) {}

// ExitLimitOperand is called when production limitOperand is exited.
func (s *BaseAQLListener) ExitLimitOperand(ctx *LimitOperandContext) {}
