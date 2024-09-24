package aql

import (
	"errors"
	"fmt"
	"github.com/acontrolfreak/openehr/aql/gen"
	"regexp"
	"strings"
)

type QueryBuilder struct {
	listener   *ParserListener
	parameters map[string]interface{}
}

func NewQueryBuilder(listener *ParserListener, parameters map[string]interface{}) *QueryBuilder {
	return &QueryBuilder{
		listener:   listener,
		parameters: parameters,
	}
}

func (q *QueryBuilder) Build() (string, error) {
	dataSourcingQuery, err := q.buildDataSourcingQuery()
	if err != nil {
		return "", err
	}

	// SELECT
	dataExtractionQuery, err := q.buildDataExtractionQuery(dataSourcingQuery)
	if err != nil {
		return "", err
	}

	return dataExtractionQuery, nil
}

func (q *QueryBuilder) buildDataSourcingQuery() (string, error) {
	selectQuery := "SELECT "

	i := 0
	for _, source := range q.listener.sources {
		if i > 0 {
			selectQuery += ","
		}

		selectQuery += fmt.Sprintf("data_%d", source.Number)
		i++
	}

	joinsQuery, conditions, err := q.buildFromLayer(nil, q.listener.nodes[0])
	if err != nil {
		return "", err
	}

	if q.listener.where != nil {
		whereQuery, err := q.buildWhereExpr(q.listener.where)
		if err != nil {
			return "", err
		}

		conditions += " AND " + whereQuery
	}

	return fmt.Sprintf("%s FROM ehr %s WHERE %s", selectQuery, joinsQuery, conditions), nil
}

func (q *QueryBuilder) buildWhereExpr(ctx *gen.WhereExprContext) (string, error) {
	if ctx.IdentifiedExpr() != nil {
		return q.buildIdentifiedExpr(ctx.IdentifiedExpr().(*gen.IdentifiedExprContext))
	} else if ctx.AND() != nil {
		leftWhereExpr, err := q.buildWhereExpr(ctx.WhereExpr(0).(*gen.WhereExprContext))
		if err != nil {
			return "", err
		}

		rightWhereExpr, err := q.buildWhereExpr(ctx.WhereExpr(1).(*gen.WhereExprContext))
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("(%s) AND (%s)", leftWhereExpr, rightWhereExpr), nil
	} else if ctx.OR() != nil {
		leftWhereExpr, err := q.buildWhereExpr(ctx.WhereExpr(0).(*gen.WhereExprContext))
		if err != nil {
			return "", err
		}

		rightWhereExpr, err := q.buildWhereExpr(ctx.WhereExpr(1).(*gen.WhereExprContext))
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("(%s) OR (%s)", leftWhereExpr, rightWhereExpr), nil
	} else if ctx.NOT() != nil {
		whereExprQuery, err := q.buildWhereExpr(ctx.WhereExpr(0).(*gen.WhereExprContext))
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("NOT (%s)", whereExprQuery), nil
	} else if ctx.SYM_LEFT_PAREN() != nil {
		return q.buildWhereExpr(ctx.WhereExpr(0).(*gen.WhereExprContext))
	}

	return "", errors.New("unknown WHERE expr")
}

func (q *QueryBuilder) buildIdentifiedExpr(ctx *gen.IdentifiedExprContext) (string, error) {
	if ctx.EXISTS() != nil {
		return q.buildIdentifiedPath(ctx.IdentifiedPath().(*gen.IdentifiedPathContext), true)
	} else if ctx.IdentifiedPath() != nil && ctx.COMPARISON_OPERATOR() != nil {
		identifiedPathQuery, err := q.buildIdentifiedPath(ctx.IdentifiedPath().(*gen.IdentifiedPathContext), false)
		if err != nil {
			return "", err
		}

		terminalQuery, err := q.buildTerminal(ctx.Terminal().(*gen.TerminalContext), true)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("(SELECT %s) %s '%s'::jsonb", identifiedPathQuery, ctx.COMPARISON_OPERATOR().GetText(), terminalQuery), nil
	} else if ctx.FunctionCall() != nil && ctx.COMPARISON_OPERATOR() != nil {
		functionCallQuery, err := q.buildFunctionCall(ctx.FunctionCall().(*gen.FunctionCallContext))
		if err != nil {
			return "", err
		}

		terminalQuery, err := q.buildTerminal(ctx.Terminal().(*gen.TerminalContext), true)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("(SELECT %s) %s %s::jsonb", functionCallQuery, ctx.COMPARISON_OPERATOR().GetText(), terminalQuery), nil
	} else if ctx.LIKE() != nil {
		identifiedPathQuery, err := q.buildIdentifiedPath(ctx.IdentifiedPath().(*gen.IdentifiedPathContext), false)
		if err != nil {
			return "", err
		}

		var val string
		if ctx.LikeOperand().PARAMETER() != nil {
			parameterQuery, err := q.buildParameter(ctx.LikeOperand().PARAMETER().GetText(), false)
			if err != nil {
				return "", err
			}

			val = parameterQuery
		} else {
			val = ctx.LikeOperand().STRING().GetText()
		}

		singleCharacterRegexp := regexp.MustCompile(`(?<!\\)(\?)`)
		val = singleCharacterRegexp.ReplaceAllString(val, "_")

		anyNumberOfCharactersRegexp := regexp.MustCompile(`(?<!\\)(\*)`)
		val = anyNumberOfCharactersRegexp.ReplaceAllString(val, "%")

		return fmt.Sprintf("%s LIKE '%s'", identifiedPathQuery, val), nil
	} else if ctx.MATCHES() != nil {
		identifiedPathQuery, err := q.buildIdentifiedPath(ctx.IdentifiedPath().(*gen.IdentifiedPathContext), false)
		if err != nil {
			return "", err
		}

		if ctx.MatchesOperand().TerminologyFunction() != nil {
			return "", errors.New("terminology function, not implemented")
		}

		var items []string
		for _, valueListItem := range ctx.MatchesOperand().AllValueListItem() {
			if valueListItem.Primitive() != nil {
				primitiveQuery, err := q.buildPrimitive(valueListItem.Primitive().(*gen.PrimitiveContext), false)
				if err != nil {
					return "", err
				}

				items = append(items, primitiveQuery)
			} else if valueListItem.PARAMETER() != nil {
				parameterQuery, err := q.buildParameter(valueListItem.PARAMETER().GetText(), false)
				if err != nil {
					return "", err
				}

				items = append(items, parameterQuery)
			} else if valueListItem.TerminologyFunction() != nil {
				return "", errors.New("terminology function, not implemented")
			}
		}

		return fmt.Sprintf("%s IN (%s)", identifiedPathQuery, strings.Join(items, ",")), nil
	} else if ctx.IdentifiedExpr() != nil {
		return q.buildIdentifiedExpr(ctx.IdentifiedExpr().(*gen.IdentifiedExprContext))
	}

	return "", errors.New("unknown identified expr")
}

func (q *QueryBuilder) buildDataExtractionQuery(dataSourceQuery string) (string, error) {
	var selectQuery string
	var joinsQuery string

	for i, s := range q.listener.selects {
		if i > 0 {
			selectQuery += ", "
			joinsQuery += " "
		}

		if s.ColumnExpr().IdentifiedPath() != nil {
			selectQuery += fmt.Sprintf("\"#%d\"", i)
			exprQuery, err := q.buildIdentifiedPath(s.ColumnExpr().IdentifiedPath().(*gen.IdentifiedPathContext), false)
			if err != nil {
				return "", err
			}

			joinsQuery += fmt.Sprintf("LEFT JOIN %s \"#%d\" ON TRUE", exprQuery, i)
			continue
		} else if s.ColumnExpr().Primitive() != nil {
			exprQuery, err := q.buildPrimitive(s.ColumnExpr().Primitive().(*gen.PrimitiveContext), false)
			if err != nil {
				return "", err
			}

			selectQuery += fmt.Sprintf("%s \"#%d\"", exprQuery, i)
			continue
		} else if s.ColumnExpr().AggregateFunctionCall() != nil {
			identifiedPathQuery, err := q.buildIdentifiedPath(s.ColumnExpr().AggregateFunctionCall().IdentifiedPath().(*gen.IdentifiedPathContext), false)
			if err != nil {
				return "", err
			}

			exprQuery, err := q.buildAggregateFunctionCall(s.ColumnExpr().AggregateFunctionCall().(*gen.AggregateFunctionCallContext), i)
			if err != nil {
				return "", err
			}

			selectQuery += fmt.Sprintf("%s \"#%d\"", exprQuery, i)
			joinsQuery += fmt.Sprintf("LEFT JOIN %s \"#%d\" ON TRUE", identifiedPathQuery, i)
			continue
		} else if s.ColumnExpr().FunctionCall() != nil {
			exprQuery, err := q.buildFunctionCall(s.ColumnExpr().FunctionCall().(*gen.FunctionCallContext))
			if err != nil {
				return "", err
			}

			selectQuery += fmt.Sprintf("%s \"#%d\"", exprQuery, i)
			continue
		}

		return "", errors.New("undefined column expression")
	}

	return fmt.Sprintf("SELECT %s FROM (%s) data_sources %s", selectQuery, dataSourceQuery, joinsQuery), nil
}

func (q *QueryBuilder) buildFromLayer(prev, node *Node) (string, string, error) {
	switch node.Type {
	case "CLASS":
		{
			query, condition, err := q.buildFromLayerClass(prev, node)
			if err != nil {
				return "", "", err
			}

			return query, condition, nil
		}
	case "CONTAINS":
		{
			// TODO exists

			left := node.Children[0]
			right := node.Children[1]

			leftQuery, leftCondition, leftErr := q.buildFromLayer(prev, left)
			if leftErr != nil {
				return "", "", leftErr
			}

			rightQuery, rightCondition, rightErr := q.buildFromLayer(left, right)
			if rightErr != nil {
				return "", "", rightErr
			}

			query := leftQuery + " " + rightQuery
			condition := leftCondition + " AND " + rightCondition
			return query, condition, nil
		}
	}

	return "", "", fmt.Errorf("no case found for node type %s", node.Type)
}

func (q *QueryBuilder) buildFromLayerClass(prev, curr *Node) (string, string, error) {
	switch curr.ClassName {
	case "EHR":
		{
			query := fmt.Sprintf("LEFT JOIN (SELECT id, data data_%[1]d FROM ehr) tbl_%[1]d ON ehr.id = tbl_%[1]d.id", curr.Number)
			condition := fmt.Sprintf("data_%[1]d IS NOT NULL", curr.Number)

			if curr.Predicate != nil {
				predicate, err := q.buildFromLayerClassPredicate(curr)

				if err != nil {
					return "", "", err
				}

				condition = fmt.Sprintf("%s AND %s", condition, predicate)
			}

			return query, condition, nil
		}
	case "EHR_STATUS":
		{
			// TODO
		}
	case "EHR_ACCESS":
		{
			// TODO
		}
	case "COMPOSITION":
		{
			query := fmt.Sprintf("LEFT JOIN (SELECT ehr_id, data data_%[1]d FROM composition) tbl_%[1]d ON ehr.id = tbl_%[1]d.ehr_id", curr.Number)
			condition := fmt.Sprintf("data_%[1]d IS NOT NULL", curr.Number)

			if curr.Predicate != nil {
				predicate, err := q.buildFromLayerClassPredicate(curr)

				if err != nil {
					return "", "", err
				}

				condition = fmt.Sprintf("%s AND %s", condition, predicate)
			}

			return query, condition, nil
		}
	default:
		{
			if prev == nil {
				return "", "", fmt.Errorf("cannot start FROM clause with %s", curr.ClassName)
			}

			query := fmt.Sprintf("LEFT JOIN jsonb_path_query(data_%d, '$.** ? (@._type == \"%s\")') data_%d ON TRUE", prev.Number, curr.ClassName, curr.Number)
			condition := fmt.Sprintf("data_%[1]d IS NOT NULL", curr.Number)

			if curr.Predicate != nil {
				predicate, err := q.buildFromLayerClassPredicate(curr)

				if err != nil {
					return "", "", err
				}

				condition = fmt.Sprintf("%s AND %s", condition, predicate)
			}

			return query, condition, nil
		}
	}

	return "", "", errors.New("class name does not match any case")
}

func (q *QueryBuilder) buildFromLayerClassPredicate(node *Node) (string, error) {
	pathPredicateQuery, err := q.buildPathPredicate(node.Predicate)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("jsonb_path_exists(data_%d, '$ ? %s')", node.Number, pathPredicateQuery), nil
}

func (q *QueryBuilder) buildPrimitive(ctx *gen.PrimitiveContext, inlineJson bool) (string, error) {
	if ctx.STRING() != nil {
		val := ctx.STRING().GetText()
		val = val[1 : len(val)-1] // Remove quotes
		if inlineJson {
			return fmt.Sprintf("\"%s\"", val), nil
		} else {
			return fmt.Sprintf("'%s'", val), nil
		}
	} else if ctx.NumericPrimitive() != nil {
		return fmt.Sprintf("%s", ctx.NumericPrimitive().GetText()), nil
	} else if ctx.DATE() != nil {
		val := ctx.DATE().GetText()

		if inlineJson {
			var format string
			if strings.Contains(val, "-") {
				format = "YY-MM-DD"
			} else {
				format = "YYMMDD"
			}

			return fmt.Sprintf("\"%s\".datetime(\"%s\")", val, format), nil
		} else {
			return fmt.Sprintf("'%s'::date", val), nil
		}
	} else if ctx.TIME() != nil {
		val := ctx.TIME().GetText()
		if inlineJson {
			if strings.Contains(val, "+") || strings.Contains(val, "-") || strings.Contains(val, "Z") {
				return "", errors.New("no timezone support for predicates")
			}

			var format string
			if strings.Contains(val, ":") {
				format = "HH24:MI:SS"
			} else {
				format = "HH24MISS"
			}

			if strings.Contains(val, ".") {
				format += ".US"
			}

			return fmt.Sprintf("\"%s\".datetime(\"%s\")", val, format), nil
		} else {
			return fmt.Sprintf("TIME WITH TIME ZONE '%s'", val), nil
		}
	} else if ctx.DATETIME() != nil {
		val := ctx.DATETIME().GetText()

		if inlineJson {
			if strings.Contains(val, "+") || strings.Contains(val, "-") || strings.Contains(val, "Z") {
				return "", errors.New("no timezone support for predicates")
			}

			var format string
			if strings.Contains(val, "-") {
				format = "YY-MM-DD"
			} else {
				format = "YYMMDD"
			}

			if strings.Contains(val, "T") {
				format += "\\\"T\\\"HH24:MI:SS"
			} else {
				format += "\\\"T\\\"HH24MISS"
			}

			if strings.Contains(val, ".") {
				format += ".US"
			}

			return fmt.Sprintf("\"%s\".datetime(\"%s\")", val, format), nil
		} else {
			return fmt.Sprintf("'%s'::time", val), nil
		}
	} else if ctx.BOOLEAN() != nil {
		return fmt.Sprintf("%s", ctx.NumericPrimitive().GetText()), nil
	} else if ctx.NULL() != nil {
		return fmt.Sprintf("%s", ctx.NumericPrimitive().GetText()), nil
	}

	return "", errors.New("not implemented")
}

func (q *QueryBuilder) buildIdentifiedPath(ctx *gen.IdentifiedPathContext, exists bool) (string, error) {
	var query string

	source, found := q.listener.sources[ctx.IDENTIFIER().GetText()]

	if !found {
		return "", errors.New("no sources found")
	}

	if ctx.ObjectPath() != nil {
		subQuery, err := q.buildObjectPath(ctx.ObjectPath().(*gen.ObjectPathContext))
		if err != nil {
			return "", err
		}

		query += subQuery
	}

	if ctx.PathPredicate() != nil {
		pathPredicateQuery, err := q.buildPathPredicate(ctx.PathPredicate().(*gen.PathPredicateContext))
		if err != nil {
			return "", err
		}

		query = fmt.Sprintf("%s %s", pathPredicateQuery, query)
	}

	if exists {
		query = fmt.Sprintf("jsonb_path_exists(data_%d, '$ %s')", source.Number, query)
	} else {
		query = fmt.Sprintf("jsonb_path_query(data_%d, '$ %s')", source.Number, query)
	}
	return query, nil
}

func (q *QueryBuilder) buildPathPredicate(ctx *gen.PathPredicateContext) (string, error) {
	if ctx.StandardPredicate() != nil {
		objectPathQuery, err := q.buildObjectPath(ctx.StandardPredicate().ObjectPath().(*gen.ObjectPathContext))
		if err != nil {
			return "", err
		}

		comparisonOperator := ctx.StandardPredicate().COMPARISON_OPERATOR().GetText()

		pathPredicateOperandQuery, err := q.buildPathPredicateOperand(ctx.StandardPredicate().PathPredicateOperand().(*gen.PathPredicateOperandContext))
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("@%s %s %s", objectPathQuery, comparisonOperator, pathPredicateOperandQuery), nil
	}

	if ctx.ArchetypePredicate() != nil {
		return fmt.Sprintf("(@.archetype_node_id == \"%s\")", ctx.ArchetypePredicate().GetText()), nil
	}

	if ctx.NodePredicate() != nil {
		if ctx.NodePredicate().ARCHETYPE_HRID() != nil || ctx.NodePredicate().ID_CODE(0) != nil || ctx.NodePredicate().AT_CODE(0) != nil {
			var query string

			if ctx.NodePredicate().ARCHETYPE_HRID() != nil {
				query = fmt.Sprintf("@.archetype_details.archetype_id.value == \"%s\"", ctx.NodePredicate().ARCHETYPE_HRID().GetText())
			} else if ctx.NodePredicate().GetChild(0) == ctx.NodePredicate().ID_CODE(0) {
				query = fmt.Sprintf("@.archetype_node_id == \"%s\"", ctx.NodePredicate().ID_CODE(0).GetText())
			} else if ctx.NodePredicate().GetChild(0) == ctx.NodePredicate().AT_CODE(0) {
				query = fmt.Sprintf("@.archetype_node_id == \"%s\"", ctx.NodePredicate().AT_CODE(0).GetText())
			} else {
				return "", errors.New("path predicate, invalid node predicate")
			}

			if ctx.NodePredicate().SYM_COMMA() != nil {
				if ctx.NodePredicate().STRING() != nil {
					query = fmt.Sprintf("%s && @.name.value == \"%s\"", query, ctx.NodePredicate().STRING().GetText())
				} else if ctx.NodePredicate().PARAMETER() != nil {
					parameterQuery, err := q.buildParameter(ctx.NodePredicate().PARAMETER().GetText(), true)
					if err != nil {
						return "", err
					}

					query = fmt.Sprintf("%s && @.name.value == %s", query, parameterQuery)
				} else if ctx.NodePredicate().TERM_CODE() != nil {
					values := strings.Split(ctx.NodePredicate().TERM_CODE().GetText(), "::")
					query = fmt.Sprintf("%s && @.name.defining_code.code_string == \"%s\" && @.name.defining_code.terminology_id.value == \"%s\"", query, values[1], values[0])
				} else if ctx.NodePredicate().GetStop() == ctx.NodePredicate().AllAT_CODE()[len(ctx.NodePredicate().AllAT_CODE())].GetSymbol() {
					query = fmt.Sprintf("%s && @.archetype_node_id == \"%s\"", query, ctx.NodePredicate().AllAT_CODE()[len(ctx.NodePredicate().AllAT_CODE())].GetText())
				} else if ctx.NodePredicate().GetStop() == ctx.NodePredicate().AllID_CODE()[len(ctx.NodePredicate().AllID_CODE())].GetSymbol() {
					query = fmt.Sprintf("%s && @.archetype_node_id == \"%s\"", query, ctx.NodePredicate().AllID_CODE()[len(ctx.NodePredicate().AllID_CODE())].GetText())
				} else {
					return "", errors.New("path predicate, invalid node predicate after comma")
				}
			}

			return fmt.Sprintf("(%s)", query), nil

		} else if ctx.NodePredicate().ID_CODE(0) != nil || ctx.NodePredicate().AT_CODE(0) != nil {

		} else if ctx.NodePredicate().ID_CODE(0) != nil && ctx.NodePredicate().GetChild(0) == ctx.NodePredicate().ID_CODE(0) {

		} else if ctx.NodePredicate().PARAMETER() != nil {

		} else if ctx.NodePredicate().ObjectPath() != nil {

		} else if ctx.NodePredicate().AND != nil {

		} else if ctx.NodePredicate().OR != nil {
		}
	}

	return "", errors.New("path predicate, no predicate type matching expected cases")
}

func (q *QueryBuilder) buildObjectPath(ctx *gen.ObjectPathContext) (string, error) {
	var query string

	for i, part := range ctx.AllPathPart() {
		if i > 0 {
			query += " "
		}

		query += fmt.Sprintf(".%s", part.IDENTIFIER().GetText())

		if part.PathPredicate() != nil {
			predicateQuery, err := q.buildPathPredicate(part.PathPredicate().(*gen.PathPredicateContext))

			if err != nil {
				return "", err
			}

			query += fmt.Sprintf(" ? %s", predicateQuery)
		}
	}

	return query, nil
}

func (q *QueryBuilder) buildPathPredicateOperand(ctx *gen.PathPredicateOperandContext) (string, error) {
	if ctx.Primitive() != nil {
		return q.buildPrimitive(ctx.Primitive().(*gen.PrimitiveContext), true)
	} else if ctx.ObjectPath() != nil {
		return q.buildObjectPath(ctx.ObjectPath().(*gen.ObjectPathContext))
	} else if ctx.PARAMETER() != nil {
		return q.buildParameter(ctx.PARAMETER().GetText(), false)
	} else if ctx.ID_CODE() != nil {
		return ctx.ID_CODE().GetText(), nil
	} else if ctx.AT_CODE() != nil {
		return ctx.AT_CODE().GetText(), nil
	}

	return "", errors.New("path predicate operand, no predicate type matching expected cases")
}

func (q *QueryBuilder) buildParameter(name string, inlineJson bool) (string, error) {
	val, found := q.parameters[name]

	if !found {
		return "", fmt.Errorf("no parameter found for %s", name)
	}

	if _, ok := val.(string); ok {
		if inlineJson {
			return fmt.Sprintf("\"%s\"", val), nil
		}

		return fmt.Sprintf("'%s'", val), nil
	}

	return fmt.Sprintf("%v", val), nil
}

func (q *QueryBuilder) buildAggregateFunctionCall(ctx *gen.AggregateFunctionCallContext, number int) (string, error) {
	if ctx.COUNT() != nil {
		if ctx.DISTINCT() != nil {
			return fmt.Sprintf("COUNT(DISTINCT \"#%d\") OVER ()", number), nil
		}
		return fmt.Sprintf("COUNT(\"#%d\") OVER ()", number), nil
	} else if ctx.MIN() != nil {
		return fmt.Sprintf("MIN(\"#%d\") OVER ()", number), nil
	} else if ctx.MAX() != nil {
		return fmt.Sprintf("MAX(\"#%d\") OVER ()", number), nil
	} else if ctx.SUM() != nil {
		return fmt.Sprintf("SUM(\"#%d\") OVER ()", number), nil
	} else if ctx.AVG() != nil {
		return fmt.Sprintf("AVG(\"#%d\") OVER ()", number), nil
	}

	return "", errors.New("aggregate function call called without any arguments")
}

func (q *QueryBuilder) buildFunctionCall(ctx *gen.FunctionCallContext) (string, error) {
	if ctx.TerminologyFunction() != nil {
		return "", errors.New("function call with terminology function, is not implemented")
	}

	if ctx.STRING_FUNCTION_ID() != nil {
		switch strings.ToUpper(ctx.STRING_FUNCTION_ID().GetText()) {
		case "LENGTH":
			{
				if len(ctx.AllTerminal()) != 1 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				terminalQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("LENGTH(%s)", terminalQuery), nil
			}
		case "CONTAINS":
			{
				if len(ctx.AllTerminal()) != 2 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				haystackQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				needleQuery, err := q.buildTerminal(ctx.Terminal(1).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("%s LIKE concat('%%', %s , '%%'", haystackQuery, needleQuery), nil
			}
		case "POSITION":
			{
				if len(ctx.AllTerminal()) != 2 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				haystackQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				needleQuery, err := q.buildTerminal(ctx.Terminal(1).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("POSTITION(%s IN %s)", haystackQuery, needleQuery), nil
			}
		case "SUBSTRING":
			{
				l := len(ctx.AllTerminal())
				if l < 2 || l > 3 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				expressionQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				positionQuery, err := q.buildTerminal(ctx.Terminal(1).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				if l == 2 {
					return fmt.Sprintf("SUBSTRING(%s, %s)", expressionQuery, positionQuery), nil
				}

				lengthQuery, err := q.buildTerminal(ctx.Terminal(2).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("SUBSTRING(%s, %s, %s)", expressionQuery, positionQuery, lengthQuery), nil
			}
		case "CONCAT_WS":
			{
				if len(ctx.AllTerminal()) < 2 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				var expressions []string

				for _, terminal := range ctx.AllTerminal() {
					terminalQuery, err := q.buildTerminal(terminal.(*gen.TerminalContext), false)
					if err != nil {
						return "", err
					}

					expressions = append(expressions, terminalQuery)
				}

				return fmt.Sprintf("CONCAT_WS(%s)", strings.Join(expressions, ",")), nil
			}
		case "CONCAT":
			{
				if len(ctx.AllTerminal()) < 1 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				var expressions []string

				for _, terminal := range ctx.AllTerminal() {
					terminalQuery, err := q.buildTerminal(terminal.(*gen.TerminalContext), false)
					if err != nil {
						return "", err
					}

					expressions = append(expressions, terminalQuery)
				}

				return fmt.Sprintf("CONCAT(%s)", strings.Join(expressions, ",")), nil
			}
		default:
			{
				return "", errors.New("string function call but no matching case found")
			}
		}
	} else if ctx.NUMERIC_FUNCTION_ID() != nil {
		switch strings.ToUpper(ctx.NUMERIC_FUNCTION_ID().GetText()) {
		case "ABS":
			{
				if len(ctx.AllTerminal()) != 1 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				terminalQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("ABS(%s)", terminalQuery), nil
			}
		case "MOD":
			{
				if len(ctx.AllTerminal()) != 2 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				xQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				yQuery, err := q.buildTerminal(ctx.Terminal(1).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("MOD(%s, %s)", xQuery, yQuery), nil
			}
		case "CEIL":
			{
				if len(ctx.AllTerminal()) != 1 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				terminalQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("CEIL(%s)", terminalQuery), nil
			}
		case "FLOOR":
			{
				if len(ctx.AllTerminal()) != 1 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				terminalQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("FLOOR(%s)", terminalQuery), nil
			}
		case "ROUND":
			{
				if len(ctx.AllTerminal()) != 2 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				expressionQuery, err := q.buildTerminal(ctx.Terminal(0).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				decimalQuery, err := q.buildTerminal(ctx.Terminal(1).(*gen.TerminalContext), false)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("ROUND(%s, %s)", expressionQuery, decimalQuery), nil
			}
		default:
			{
				return "", errors.New("numeric function call but no matching case found")
			}
		}
	} else if ctx.DATE_TIME_FUNCTION_ID() != nil {
		switch strings.ToUpper(ctx.DATE_TIME_FUNCTION_ID().GetText()) {
		case "CURRENT_DATE":
			{
				if len(ctx.AllTerminal()) != 0 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				return "CURRENT_DATE()", nil
			}
		case "CURRENT_TIME":
			{
				if len(ctx.AllTerminal()) != 0 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				return "CURRENT_TIME()", nil
			}
		case "NOW":
			fallthrough
		case "CURRENT_DATE_TIME":
			{
				if len(ctx.AllTerminal()) != 0 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				return "NOW()", nil
			}
		case "CURRENT_TIMEZONE":
			{
				if len(ctx.AllTerminal()) != 0 {
					return "", errors.New("terminology function call called with invalid number of arguments")
				}

				return "(SELECT TO_CHAR(utc_offset, 'HH24:MM') from pg_timezone_names where name = current_setting('TIMEZONE'))", nil
			}
		default:
			{
				return "", errors.New("numeric function call but no matching case found")
			}
		}
	} else if ctx.IDENTIFIER() != nil {
		return "", fmt.Errorf("function call with identifier '%s', not implemented", ctx.IDENTIFIER().GetText())
	}

	return "", errors.New("function call called but no case matches")
}

func (q *QueryBuilder) buildTerminal(ctx *gen.TerminalContext, inlineJson bool) (string, error) {
	if ctx.Primitive() != nil {
		return q.buildPrimitive(ctx.Primitive().(*gen.PrimitiveContext), inlineJson)
	} else if ctx.PARAMETER() != nil {
		return q.buildParameter(ctx.PARAMETER().GetText(), inlineJson)
	} else if ctx.IdentifiedPath() != nil {
		return q.buildIdentifiedPath(ctx.IdentifiedPath().(*gen.IdentifiedPathContext), false)
	} else if ctx.FunctionCall() != nil {
		return q.buildFunctionCall(ctx.FunctionCall().(*gen.FunctionCallContext))
	}

	return "", errors.New("terminology function call called but no case matches")
}
