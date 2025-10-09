// Code generated from AQL.g4 by ANTLR 4.13.2. DO NOT EDIT.

package gen // AQL
import "github.com/antlr4-go/antlr/v4"


// AQLListener is a complete listener for a parse tree produced by AQLParser.
type AQLListener interface {
	antlr.ParseTreeListener

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterSelectQuery is called when entering the selectQuery production.
	EnterSelectQuery(c *SelectQueryContext)

	// EnterSelectClause is called when entering the selectClause production.
	EnterSelectClause(c *SelectClauseContext)

	// EnterFromClause is called when entering the fromClause production.
	EnterFromClause(c *FromClauseContext)

	// EnterJoinClause is called when entering the joinClause production.
	EnterJoinClause(c *JoinClauseContext)

	// EnterWhereClause is called when entering the whereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterGroupByClause is called when entering the groupByClause production.
	EnterGroupByClause(c *GroupByClauseContext)

	// EnterOrderByClause is called when entering the orderByClause production.
	EnterOrderByClause(c *OrderByClauseContext)

	// EnterLimitOffsetClause is called when entering the limitOffsetClause production.
	EnterLimitOffsetClause(c *LimitOffsetClauseContext)

	// EnterSelectExpr is called when entering the selectExpr production.
	EnterSelectExpr(c *SelectExprContext)

	// EnterFromExpr is called when entering the fromExpr production.
	EnterFromExpr(c *FromExprContext)

	// EnterJoinExpr is called when entering the joinExpr production.
	EnterJoinExpr(c *JoinExprContext)

	// EnterWhereExpr is called when entering the whereExpr production.
	EnterWhereExpr(c *WhereExprContext)

	// EnterOrderByExpr is called when entering the orderByExpr production.
	EnterOrderByExpr(c *OrderByExprContext)

	// EnterColumnExpr is called when entering the columnExpr production.
	EnterColumnExpr(c *ColumnExprContext)

	// EnterIdentifiedPath is called when entering the identifiedPath production.
	EnterIdentifiedPath(c *IdentifiedPathContext)

	// EnterObjectPath is called when entering the objectPath production.
	EnterObjectPath(c *ObjectPathContext)

	// EnterPathPart is called when entering the pathPart production.
	EnterPathPart(c *PathPartContext)

	// EnterPathCondition is called when entering the pathCondition production.
	EnterPathCondition(c *PathConditionContext)

	// EnterBooleanCondition is called when entering the booleanCondition production.
	EnterBooleanCondition(c *BooleanConditionContext)

	// EnterPathConditionOperand is called when entering the pathConditionOperand production.
	EnterPathConditionOperand(c *PathConditionOperandContext)

	// EnterComparisonOperand is called when entering the comparisonOperand production.
	EnterComparisonOperand(c *ComparisonOperandContext)

	// EnterInOperand is called when entering the inOperand production.
	EnterInOperand(c *InOperandContext)

	// EnterInOperandValue is called when entering the inOperandValue production.
	EnterInOperandValue(c *InOperandValueContext)

	// EnterPrimitive is called when entering the primitive production.
	EnterPrimitive(c *PrimitiveContext)

	// EnterIntPrimitive is called when entering the intPrimitive production.
	EnterIntPrimitive(c *IntPrimitiveContext)

	// EnterFloatPrimitive is called when entering the floatPrimitive production.
	EnterFloatPrimitive(c *FloatPrimitiveContext)

	// EnterStringOperand is called when entering the stringOperand production.
	EnterStringOperand(c *StringOperandContext)

	// EnterIntOperand is called when entering the intOperand production.
	EnterIntOperand(c *IntOperandContext)

	// EnterNumbericOperand is called when entering the numbericOperand production.
	EnterNumbericOperand(c *NumbericOperandContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterAggregateFunctionCall is called when entering the aggregateFunctionCall production.
	EnterAggregateFunctionCall(c *AggregateFunctionCallContext)

	// EnterLimitOperand is called when entering the limitOperand production.
	EnterLimitOperand(c *LimitOperandContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitSelectQuery is called when exiting the selectQuery production.
	ExitSelectQuery(c *SelectQueryContext)

	// ExitSelectClause is called when exiting the selectClause production.
	ExitSelectClause(c *SelectClauseContext)

	// ExitFromClause is called when exiting the fromClause production.
	ExitFromClause(c *FromClauseContext)

	// ExitJoinClause is called when exiting the joinClause production.
	ExitJoinClause(c *JoinClauseContext)

	// ExitWhereClause is called when exiting the whereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitGroupByClause is called when exiting the groupByClause production.
	ExitGroupByClause(c *GroupByClauseContext)

	// ExitOrderByClause is called when exiting the orderByClause production.
	ExitOrderByClause(c *OrderByClauseContext)

	// ExitLimitOffsetClause is called when exiting the limitOffsetClause production.
	ExitLimitOffsetClause(c *LimitOffsetClauseContext)

	// ExitSelectExpr is called when exiting the selectExpr production.
	ExitSelectExpr(c *SelectExprContext)

	// ExitFromExpr is called when exiting the fromExpr production.
	ExitFromExpr(c *FromExprContext)

	// ExitJoinExpr is called when exiting the joinExpr production.
	ExitJoinExpr(c *JoinExprContext)

	// ExitWhereExpr is called when exiting the whereExpr production.
	ExitWhereExpr(c *WhereExprContext)

	// ExitOrderByExpr is called when exiting the orderByExpr production.
	ExitOrderByExpr(c *OrderByExprContext)

	// ExitColumnExpr is called when exiting the columnExpr production.
	ExitColumnExpr(c *ColumnExprContext)

	// ExitIdentifiedPath is called when exiting the identifiedPath production.
	ExitIdentifiedPath(c *IdentifiedPathContext)

	// ExitObjectPath is called when exiting the objectPath production.
	ExitObjectPath(c *ObjectPathContext)

	// ExitPathPart is called when exiting the pathPart production.
	ExitPathPart(c *PathPartContext)

	// ExitPathCondition is called when exiting the pathCondition production.
	ExitPathCondition(c *PathConditionContext)

	// ExitBooleanCondition is called when exiting the booleanCondition production.
	ExitBooleanCondition(c *BooleanConditionContext)

	// ExitPathConditionOperand is called when exiting the pathConditionOperand production.
	ExitPathConditionOperand(c *PathConditionOperandContext)

	// ExitComparisonOperand is called when exiting the comparisonOperand production.
	ExitComparisonOperand(c *ComparisonOperandContext)

	// ExitInOperand is called when exiting the inOperand production.
	ExitInOperand(c *InOperandContext)

	// ExitInOperandValue is called when exiting the inOperandValue production.
	ExitInOperandValue(c *InOperandValueContext)

	// ExitPrimitive is called when exiting the primitive production.
	ExitPrimitive(c *PrimitiveContext)

	// ExitIntPrimitive is called when exiting the intPrimitive production.
	ExitIntPrimitive(c *IntPrimitiveContext)

	// ExitFloatPrimitive is called when exiting the floatPrimitive production.
	ExitFloatPrimitive(c *FloatPrimitiveContext)

	// ExitStringOperand is called when exiting the stringOperand production.
	ExitStringOperand(c *StringOperandContext)

	// ExitIntOperand is called when exiting the intOperand production.
	ExitIntOperand(c *IntOperandContext)

	// ExitNumbericOperand is called when exiting the numbericOperand production.
	ExitNumbericOperand(c *NumbericOperandContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitAggregateFunctionCall is called when exiting the aggregateFunctionCall production.
	ExitAggregateFunctionCall(c *AggregateFunctionCallContext)

	// ExitLimitOperand is called when exiting the limitOperand production.
	ExitLimitOperand(c *LimitOperandContext)
}
