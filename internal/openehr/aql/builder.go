package aql

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/openehr/aql/gen"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

type Source struct {
	Model string
	Table string
	Alias string
}

func ToSQL(aqlQuery string, params map[string]any) (string, []string, error) {
	listener := NewTreeShapeListener()
	errorListener := NewErrorListener()

	input := antlr.NewInputStream(aqlQuery)
	lexer := gen.NewAQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	p := gen.NewAQLParser(stream)
	p.AddErrorListener(errorListener)

	antlr.ParseTreeWalkerDefault.Walk(listener, p.Query())

	if len(errorListener.Errors) > 0 {
		return "", nil, errors.Join(errorListener.Errors...)
	}

	query, columnNames, err := BuildSelectQuery(listener.Query.SelectQuery(), params)
	if err != nil {
		return "", nil, err
	}

	// Wrap the query to return a JSON array
	query = fmt.Sprintf("SELECT jsonb_build_array(%s) FROM (%s) AS result", strings.Join(columnNames, ", "), query)

	return query, columnNames, nil
}

func BuildSelectQuery(ctx gen.ISelectQueryContext, params map[string]any) (string, []string, error) {
	// FROM
	fromClause, additionalWhereExpressions, sources, err := BuildFromClause(ctx.FromClause(), params)
	if err != nil {
		return "", nil, err
	}

	// WHERE
	whereClause, err := BuildWhereClause(ctx.WhereClause(), params, sources, additionalWhereExpressions)
	if err != nil {
		return "", nil, err
	}

	// SELECT
	selectClause, columnNames, selectHelperTables, singleRow, err := BuildSelectClause(ctx.SelectClause(), params, sources)
	if err != nil {
		return "", nil, err
	}

	// todo
	orderByClause := ""

	// LIMIT / OFFSET
	limitOffsetClause, err := BuildLimitOffsetClause(ctx.LimitClause(), params, singleRow)
	if err != nil {
		return "", nil, err
	}

	// Add helper tables to FROM clause
	if len(selectHelperTables) > 0 {
		fromClause += " " + strings.Join(selectHelperTables, " ")
	}

	// Query
	query := fmt.Sprintf("SELECT * FROM (%s %s %s) dataset %s %s", selectClause, fromClause, whereClause, orderByClause, limitOffsetClause)

	return query, columnNames, nil
}

func BuildSelectClause(ctx gen.ISelectClauseContext, params map[string]any, sources []Source) (string, []string, []string, bool, error) {
	clause := "SELECT "
	if ctx.DISTINCT() != nil {
		clause += "DISTINCT "
	}

	columnNames := make([]string, 0)
	helperTables := make([]string, 0)
	singleRow := true

	for idx, expr := range ctx.AllSelectExpr() {
		if idx > 0 {
			clause += ", "
		}

		selectExpressions, selectColumnNames, helperTable, singleRowExpr, err := BuildSelectExpr(expr, params, sources, len(columnNames))
		if err != nil {
			return "", nil, nil, singleRow, err
		}

		clause += strings.Join(selectExpressions, ", ")
		if helperTable != "" {
			helperTables = append(helperTables, helperTable)
		}

		columnNames = append(columnNames, selectColumnNames...)
		if !singleRowExpr {
			singleRow = false
		}
	}

	return clause, columnNames, helperTables, singleRow, nil
}

func BuildSelectExpr(ctx gen.ISelectExprContext, params map[string]any, sources []Source, columnNumber int) ([]string, []string, string, bool, error) {
	switch true {
	case ctx.SYM_ASTERISK() != nil:
		expressions := make([]string, 0)
		columnNames := make([]string, 0)
		for i, source := range sources {
			name := fmt.Sprintf(`f%d`, columnNumber+i)
			if source.Alias != "" {
				name = source.Alias
			}
			expr := fmt.Sprintf("%s.data AS %s", source.Table, name)
			expressions = append(expressions, expr)
			columnNames = append(columnNames, name)
		}
		return expressions, columnNames, "", false, nil
	case ctx.ColumnExpr() != nil:
		column, helperTable, singleRow, err := BuildColumnExpr(ctx.ColumnExpr(), params, sources, columnNumber)
		if err != nil {
			return nil, nil, "", false, err
		}

		name := fmt.Sprintf(`f%d`, columnNumber)
		if ctx.GetAliasName() != nil {
			name = ctx.GetAliasName().GetText()
		}
		expr := fmt.Sprintf("%s AS %s", column, name)
		return []string{expr}, []string{name}, helperTable, singleRow, nil
	default:
		return nil, nil, "", false, fmt.Errorf("unsupported select expression")
	}
}

func BuildColumnExpr(ctx gen.IColumnExprContext, params map[string]any, sources []Source, columnNumber int) (string, string, bool, error) {
	switch true {
	case ctx.Primitive() != nil:
		value, err := BuildPrimitive(ctx.Primitive(), true)
		return value, "", true, err
	case ctx.IdentifiedPath() != nil:
		source, path, _, _, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", "", false, err
		}

		return fmt.Sprintf("jsonb_path_query(%s.data, '%s')", source.Table, path), "", false, nil
	case ctx.FunctionCall() != nil:
		value, err := BuildFunctionCall(ctx.FunctionCall(), params, sources)
		return value, "", false, err
	case ctx.AggregateFunctionCall() != nil:
		source, path, err := BuildAggregateFunctionCall(ctx.AggregateFunctionCall(), params, sources, columnNumber)
		if err != nil {
			return "", "", false, err
		}

		return source, path, true, nil
	default:
		return "", "", false, fmt.Errorf("unsupported column expression")
	}
}

func BuildAggregateFunctionCall(ctx gen.IAggregateFunctionCallContext, params map[string]any, sources []Source, columnNumber int) (string, string, error) {
	switch true {
	case ctx.COUNT() != nil:
		switch true {
		case ctx.SYM_ASTERISK() != nil:
			if ctx.DISTINCT() != nil {
				return "COUNT(DISTINCT *)", "", nil
			}

			return "COUNT(*)", "", nil
		case ctx.IdentifiedPath() != nil:
			source, path, _, endsWith, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
			if err != nil {
				return "", "", err
			}

			expression := fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')", source.Table, path)
			if endsWith == nil {
				// Doing a count distinct when there are possibly multiple values in the path is not supported
				expression = fmt.Sprintf("%s.data", source.Table)
			}

			if ctx.DISTINCT() != nil {
				return fmt.Sprintf("COUNT(DISTINCT %s)", expression), "", nil
			}

			return fmt.Sprintf("SUM(ag_source_%d.data)", columnNumber),
				fmt.Sprintf("LEFT JOIN LATERAL (SELECT COALESCE((jsonb_path_query_first(%s.data, '%s.size()') #>> '{}')::int, 0) data) agg_source_%d ON TRUE", source.Table, path, columnNumber),
				nil
		default:
			return "", "", fmt.Errorf("unsupported COUNT argument")
		}
	case ctx.MIN() != nil:
		source, _, pathWithoutEnd, endsWith, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", "", err
		}

		// Worst case scenario, we need to do a switch case for every possible comparable object model.
		switchExpression := BuildSlowValueExtractionExpr(endsWith)
		if endsWith == nil {
			return fmt.Sprintf("FIRST_VALUE(agg_source_%d.data) OVER (ORDER BY agg_source_%d.data ASC)", columnNumber, columnNumber),
				fmt.Sprintf("LEFT JOIN LATERAL (SELECT (%s) data FROM %s target) agg_source_%d ON TRUE", switchExpression, source.Table, columnNumber),
				nil

		}

		pathEnding, err := BuildPathEnding(endsWith, params)
		if err != nil {
			return "", "", err
		}

		return fmt.Sprintf("FIRST_VALUE(agg_source_%d.data) OVER (ORDER BY agg_source_%d.sortable_data ASC)", columnNumber, columnNumber),
			fmt.Sprintf("LEFT JOIN LATERAL (SELECT target.data data, (%s) sortable_data FROM JSON_TABLE(%s.data, '%s' COLUMNS(data JSONB PATH '$')) parent, JSON_TABLE(parent.data, '%s' COLUMNS(data JSONB PATH '$')) target) agg_source_%d ON TRUE", switchExpression, source.Table, pathWithoutEnd, pathEnding, columnNumber),
			nil
	case ctx.MAX() != nil:
		source, _, pathWithoutEnd, endsWith, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", "", err
		}

		// Worst case scenario, we need to do a switch case for every possible comparable object model.
		switchExpression := BuildSlowValueExtractionExpr(endsWith)
		if endsWith == nil {
			return fmt.Sprintf("FIRST_VALUE(agg_source_%d.data) OVER (ORDER BY agg_source_%d.data DESC)", columnNumber, columnNumber),
				fmt.Sprintf("LEFT JOIN LATERAL (SELECT (%s) data FROM %s target) agg_source_%d ON TRUE", switchExpression, source.Table, columnNumber),
				nil
		}

		pathEnding, err := BuildPathEnding(endsWith, params)
		if err != nil {
			return "", "", err
		}

		return fmt.Sprintf("FIRST_VALUE(agg_source_%d.data) OVER (ORDER BY agg_source_%d.sortable_data DESC)", columnNumber, columnNumber),
			fmt.Sprintf("LEFT JOIN LATERAL (SELECT target.data data, (%s) sortable_data FROM JSON_TABLE(%s.data, '%s' COLUMNS(data JSONB PATH '$')) parent, JSON_TABLE(parent.data, '%s' COLUMNS(data JSONB PATH '$')) target) agg_source_%d ON TRUE", switchExpression, source.Table, pathWithoutEnd, pathEnding, columnNumber),
			nil
	case ctx.SUM() != nil:
		source, _, pathWithoutEnd, endsWith, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", "", err
		}

		switchExpression := BuildSlowValueExtractionExpr(endsWith)

		if endsWith == nil {
			return fmt.Sprintf("SUM(agg_source_%d.data)", columnNumber),
				fmt.Sprintf("LEFT JOIN LATERAL (SELECT ((%s) #>> '{}')::decimal data FROM %s target) agg_source_%d ON TRUE", switchExpression, source.Table, columnNumber),
				nil
		}

		pathEnding, err := BuildPathEnding(endsWith, params)
		if err != nil {
			return "", "", err
		}

		return fmt.Sprintf("SUM(agg_source_%d.data)", columnNumber),
			fmt.Sprintf("LEFT JOIN LATERAL (SELECT ((%s) #>> '{}')::decimal data FROM JSON_TABLE(%s.data, '%s' COLUMNS(data JSONB PATH '$')) parent, JSON_TABLE(parent.data, '%s' COLUMNS(data JSONB PATH '$')) target) agg_source_%d ON TRUE", switchExpression, source.Table, pathWithoutEnd, pathEnding, columnNumber),
			nil
	case ctx.AVG() != nil:
		source, _, pathWithoutEnd, endsWith, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", "", err
		}

		switchExpression := BuildSlowValueExtractionExpr(endsWith)
		if endsWith == nil {
			return fmt.Sprintf("AVG(agg_source_%d.data)", columnNumber),
				fmt.Sprintf("LEFT JOIN LATERAL (SELECT ((%s) #>> '{}')::decimal data FROM %s target) agg_source_%d ON TRUE", switchExpression, source.Table, columnNumber),
				nil
		}

		pathEnding, err := BuildPathEnding(endsWith, params)
		if err != nil {
			return "", "", err
		}

		return fmt.Sprintf("AVG(agg_source_%d.data)", columnNumber),
			fmt.Sprintf("LEFT JOIN LATERAL (SELECT ((%s) #>> '{}')::decimal data FROM JSON_TABLE(%s.data, '%s' COLUMNS(data JSONB PATH '$')) parent, JSON_TABLE(parent.data, '%s' COLUMNS(data JSONB PATH '$')) target) agg_source_%d ON TRUE", switchExpression, source.Table, pathWithoutEnd, pathEnding, columnNumber),
			nil
	default:
		return "", "", fmt.Errorf("unsupported aggregate function call")
	}
}

func BuildFromClause(ctx gen.IFromClauseContext, params map[string]any) (string, string, []Source, error) {
	sourceNumber := 0
	fromExpr, whereExpr, sources, err := BuildContainsExpr(ctx.FromExpr().ContainsExpr(), params, util.None[Source](), false, &sourceNumber)
	return fmt.Sprintf("FROM %s", fromExpr), whereExpr, sources, err
}

func BuildContainsExpr(ctx gen.IContainsExprContext, params map[string]any, prevSource util.Optional[Source], searchInModel bool, sourceNumber *int) (string, string, []Source, error) {
	switch true {
	case ctx.ClassExprOperand() != nil:
		source := Source{
			Model: ctx.ClassExprOperand().IDENTIFIER(0).GetText(),
			Table: fmt.Sprintf("source_%d", *sourceNumber),
		}
		if ctx.ClassExprOperand().GetAlias() != nil {
			source.Alias = ctx.ClassExprOperand().GetAlias().GetText()
		}

		*sourceNumber++ // Increment for next source

		// If null is returned, it means didn't find a special model relation.
		// So, we should change the search strategy to search in the model itself
		fromExpression, err := BuildClassExprOperand(ctx.ClassExprOperand(), params, source, prevSource, searchInModel)
		if err != nil {
			return "", "", nil, err
		}
		if fromExpression == "" {
			searchInModel = true
			fromExpression, err = BuildClassExprOperand(ctx.ClassExprOperand(), params, source, prevSource, searchInModel)
			if err != nil {
				return "", "", nil, err
			}
		}

		whereExpression := fmt.Sprintf("%s.data IS NOT NULL", source.Table)

		// Nothing to join anymore
		if ctx.CONTAINS() == nil {
			return fromExpression, whereExpression, []Source{source}, nil
		}

		// Continue with next contains expression
		nextFrom, nextWhereExpression, nextSources, err := BuildContainsExpr(ctx.ContainsExpr(0), params, util.Some(source), searchInModel, sourceNumber)
		if err != nil {
			return "", "", nil, err
		}

		// Merge sources
		sources := []Source{source}
		for _, nextSource := range nextSources {
			if nextSource.Alias != "" {
				// Check for duplicate use of alias
				for _, source := range sources {
					if source.Alias == nextSource.Alias {
						return "", "", nil, fmt.Errorf("duplicate use of source alias: %s", nextSource.Alias)
					}
				}
			}

			sources = append(sources, nextSource)
		}

		return fmt.Sprintf("%s %s", fromExpression, nextFrom), fmt.Sprintf("(%s) AND (%s)", whereExpression, nextWhereExpression), sources, nil
	case ctx.AND() != nil:
		leftFrom, leftWhere, leftSources, err := BuildContainsExpr(ctx.ContainsExpr(0), params, prevSource, searchInModel, sourceNumber)
		if err != nil {
			return "", "", nil, err
		}
		rightFrom, rightWhere, rightSources, err := BuildContainsExpr(ctx.ContainsExpr(1), params, prevSource, searchInModel, sourceNumber)
		if err != nil {
			return "", "", nil, err
		}

		// Merge sources
		sources := leftSources
		for _, rightSource := range rightSources {
			if rightSource.Alias != "" {
				// Check for duplicate use of alias
				for _, source := range sources {
					if source.Alias == rightSource.Alias {
						return "", "", nil, fmt.Errorf("duplicate use of source alias: %s", rightSource.Alias)
					}
				}
			}

			sources = append(sources, rightSource)
		}

		return fmt.Sprintf("%s %s", leftFrom, rightFrom), fmt.Sprintf("(%s) AND (%s)", leftWhere, rightWhere), sources, nil
	case ctx.OR() != nil:
		leftFrom, leftWhere, leftSources, err := BuildContainsExpr(ctx.ContainsExpr(0), params, prevSource, searchInModel, sourceNumber)
		if err != nil {
			return "", "", nil, err
		}
		rightFrom, rightWhere, rightSources, err := BuildContainsExpr(ctx.ContainsExpr(1), params, prevSource, searchInModel, sourceNumber)
		if err != nil {
			return "", "", nil, err
		}

		// Merge sources
		sources := leftSources
		for _, rightSource := range rightSources {
			if rightSource.Alias != "" {
				// Check for duplicate use of alias
				for _, source := range sources {
					if source.Alias == rightSource.Alias {
						return "", "", nil, fmt.Errorf("duplicate use of source alias: %s", rightSource.Alias)
					}
				}
			}

			sources = append(sources, rightSource)
		}

		return fmt.Sprintf("%s %s", leftFrom, rightFrom), fmt.Sprintf("(%s) OR (%s)", leftWhere, rightWhere), sources, nil
	case ctx.SYM_LEFT_PAREN() != nil:
		return BuildContainsExpr(ctx.ContainsExpr(0), params, prevSource, searchInModel, sourceNumber)
	default:
		return "", "", nil, fmt.Errorf("unsupported contains expression")
	}
}

func BuildClassExprOperand(ctx gen.IClassExprOperandContext, params map[string]any, source Source, prevSource util.Optional[Source], searchInModel bool) (string, error) {
	model := strings.ToUpper(ctx.IDENTIFIER(0).GetText())

	// Take care of predicates
	whereExpression := ""
	allVersions := false
	if ctx.PathPredicate() != nil {
		predicate := ctx.PathPredicate()
		switch true {
		case predicate.ALL_VERSIONS() != nil:
			allVersions = true
		case predicate.LATEST_VERSION() != nil:
			allVersions = false
		case predicate.NodePredicate() != nil:
			condition, err := BuildNodePredicate(predicate.NodePredicate(), params)
			if err != nil {
				return "", err
			}
			whereExpression = fmt.Sprintf("data @?? '$ ? (%s)'", condition)
		}
	}

	if searchInModel {
		// [Freek] Allow for generic searches where you want inheriting models to be included
		relatedModels := ModelInheritanceTable(model)
		typeFilter := ""
		for i, relatedModel := range relatedModels {
			if i > 0 {
				typeFilter += " || "
			}
			typeFilter += fmt.Sprintf(`@._type == "%s"`, relatedModel)
		}

		// Search in the model itself
		query := fmt.Sprintf("SELECT * FROM JSON_TABLE(%s.data, 'strict $.*.** ? (%s)' COLUMNS(data JSONB PATH '$')) data", prevSource.V.Table, typeFilter)
		if whereExpression != "" {
			query += " WHERE " + whereExpression
		}
		return fmt.Sprintf("LEFT JOIN LATERAL (%s) %s ON TRUE", query, source.Table), nil
	}

	switch model {
	case openehr.EHR_MODEL_NAME:
		expression := "SELECT id, data FROM openehr.tbl_ehr_data"
		if whereExpression != "" {
			expression += " WHERE " + whereExpression
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		return "", nil
	case openehr.CONTRIBUTION_MODEL_NAME:
		expression := "SELECT c.id, c.ehr_id, cd.data FROM openehr.tbl_contribution c JOIN openehr.tbl_contribution_data cd ON c.id = cd.id"
		if whereExpression != "" {
			expression += " WHERE " + whereExpression
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.ehr_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.VERSIONED_EHR_STATUS_MODEL_NAME:
		expression := fmt.Sprintf("SELECT vo.id, vo.ehr_id, vo.contribution_id, vod.data FROM openehr.tbl_versioned_object vo JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id WHERE vo.type = '%s'", openehr.VERSIONED_EHR_STATUS_MODEL_NAME)
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.ehr_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.VERSIONED_EHR_ACCESS_MODEL_NAME:
		expression := fmt.Sprintf("SELECT vo.id, vo.ehr_id, vo.contribution_id, vod.data FROM openehr.tbl_versioned_object vo JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id WHERE vo.type = '%s'", openehr.VERSIONED_EHR_ACCESS_MODEL_NAME)
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.ehr_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.VERSIONED_COMPOSITION_MODEL_NAME:
		expression := fmt.Sprintf("SELECT vo.id, vo.ehr_id, vo.contribution_id, vod.data FROM openehr.tbl_versioned_object vo JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id WHERE vo.type = '%s'", openehr.VERSIONED_COMPOSITION_MODEL_NAME)
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.ehr_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.VERSIONED_FOLDER_MODEL_NAME:
		expression := fmt.Sprintf("SELECT vo.id, vo.ehr_id, vo.contribution_id, vod.data FROM openehr.tbl_versioned_object vo JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id WHERE vo.type = '%s'", openehr.VERSIONED_FOLDER_MODEL_NAME)
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.ehr_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.VERSIONED_PARTY_MODEL_NAME:
		expression := fmt.Sprintf("SELECT vo.id, vo.ehr_id, vo.contribution_id, vod.data FROM openehr.tbl_versioned_object vo JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id WHERE vo.type = '%s'", openehr.VERSIONED_PARTY_MODEL_NAME)
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf(`
				LEFT JOIN openehr.tbl_object_version tmp_es_%[2]s ON tmp_es_%[2]s.ehr_id = %[3]s.id AND tmp_es_%[2]s.type = '%[4]s'
				LEFT JOIN openehr.tbl_object_version_data tmp_esd_%[2]s ON tmp_es_%[2]s.id = tmp_esd_%[2]s.id
				LEFT JOIN %[1]s
					ON %[2]s.id = tmp_esd_%[2]s.object_data->'subject'->'external_ref'->'id'->>'value'
					AND tmp_esd_%[2]s.object_data->'subject'->'external_ref'->>'namespace' = 'local'
					AND tmp_esd_%[2]s.object_data->'subject'->'external_ref'->>'type' = '%[5]s'
			`, expression, source.Table, prevSource.V.Table, openehr.EHR_STATUS_MODEL_NAME, openehr.VERSIONED_PARTY_MODEL_NAME), nil
		case openehr.EHR_STATUS_MODEL_NAME:
			return fmt.Sprintf(`
				LEFT JOIN %[1]s
					ON %[2]s.id = %[3]s.data->'subject'->'external_ref'->'id'->>'value'
					AND %[3]s.data->'subject'->'external_ref'->>'namespace' = 'local'
					AND %[3]s.data->'subject'->'external_ref'->>'type' = '%[4]s'
			`, expression, source.Table, prevSource.V.Table, openehr.VERSIONED_PARTY_MODEL_NAME), nil
		default:
			return "", nil
		}
	case openehr.COMPOSITION_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (vo.versioned_object_id) "
		}
		expression += "vo.id, vo.ehr_id, vo.versioned_object_id, vo.contribution_id, ovd.object_data data FROM openehr.tbl_object_version vo JOIN openehr.tbl_object_version_data ovd ON vo.id = ovd.id WHERE vo.type = 'COMPOSITION'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY vo.versioned_object_id, vo.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.ehr_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.VERSIONED_COMPOSITION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.FOLDER_MODEL_NAME:
			return fmt.Sprintf(`LEFT JOIN LATERAL (SELECT composition_id FROM JSON_TABLE(%s.data, '$.**.items ? (@.type == "COMPOSITION")' COLUMNS (composition_id text PATH '$.id.value'))) AS tmp_%s ON TRUE LEFT JOIN %s ON %s.id = tmp_%s.composition_id`, prevSource.V.Table, source.Table, expression, source.Table, source.Table), nil
		default:
			return "", nil
		}
	case openehr.EHR_STATUS_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'EHR_STATUS'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.ehr_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.VERSIONED_EHR_STATUS_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.EHR_ACCESS_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'EHR_ACCESS'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.ehr_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.VERSIONED_EHR_ACCESS_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.FOLDER_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'FOLDER'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		if !prevSource.E {
			return expression, nil
		}
		switch prevSource.V.Model {
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.ehr_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.VERSIONED_FOLDER_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.ROLE_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'ROLE'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.VERSIONED_PARTY_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.PERSON_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'PERSON'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.VERSIONED_PARTY_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.EHR_MODEL_NAME:
			return fmt.Sprintf(`
				LEFT JOIN openehr.tbl_object_version tmp_es_%[2]s ON tmp_es_%[2]s.ehr_id = %[3]s.id AND tmp_es_%[2]s.type = '%[4]s'
				LEFT JOIN openehr.tbl_object_version_data tmp_esd_%[2]s ON tmp_es_%[2]s.id = tmp_esd_%[2]s.id
				LEFT JOIN %[1]s
					ON %[2]s.id = tmp_esd_%[2]s.object_data->'subject'->'external_ref'->'id'->>'value'
					AND tmp_esd_%[2]s.object_data->'subject'->'external_ref'->>'namespace' = 'local'
					AND tmp_esd_%[2]s.object_data->'subject'->'external_ref'->>'type' = '%[5]s'
			`, expression, source.Table, prevSource.V.Table, openehr.EHR_STATUS_MODEL_NAME, openehr.PERSON_MODEL_NAME), nil
		case openehr.EHR_STATUS_MODEL_NAME:
			return fmt.Sprintf(`
				LEFT JOIN %[1]s
					ON %[2]s.id = %[3]s.data->'subject'->'external_ref'->'id'->>'value'
					AND %[3]s.data->'subject'->'external_ref'->>'namespace' = 'local'
					AND %[3]s.data->'subject'->'external_ref'->>'type' = '%[4]s'
			`, expression, source.Table, prevSource.V.Table, openehr.PERSON_MODEL_NAME), nil
		default:
			return "", nil
		}
	case openehr.AGENT_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'AGENT'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.VERSIONED_PARTY_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.ORGANISATION_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'ORGANISATION'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.VERSIONED_PARTY_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	case openehr.GROUP_MODEL_NAME:
		expression := "SELECT "
		if !allVersions {
			expression += "DISTINCT ON (ov.versioned_object_id) "
		}
		expression += "ov.id, ov.versioned_object_id, ov.ehr_id, ov.contribution_id, ovd.object_data data FROM openehr.tbl_object_version ov JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id WHERE ov.type = 'GROUP'"
		if whereExpression != "" {
			expression += " AND " + whereExpression
		}
		if !allVersions {
			expression += " ORDER BY ov.versioned_object_id, ov.id DESC"
		}
		expression = "(" + expression + ") " + source.Table

		if !prevSource.E {
			return expression, nil
		}

		switch prevSource.V.Model {
		case openehr.VERSIONED_PARTY_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %s ON %s.versioned_object_id = %s.id", expression, source.Table, prevSource.V.Table), nil
		case openehr.CONTRIBUTION_MODEL_NAME:
			return fmt.Sprintf("LEFT JOIN %[1]s ON %[2]s.contribution_id = %[3]s.id", expression, source.Table, prevSource.V.Table), nil
		default:
			return "", nil
		}
	default:
		if !prevSource.E {
			return "", fmt.Errorf("unsupported operantion: FROM %s", model)
		}

		return "", nil
	}
}

func BuildWhereClause(ctx gen.IWhereClauseContext, params map[string]any, sources []Source, additionalExpression string) (string, error) {
	if ctx == nil {
		if additionalExpression == "" {
			return "", nil
		}

		return "WHERE " + additionalExpression, nil
	}

	expression, err := BuildWhereExpr(ctx.WhereExpr(), params, sources)
	if err != nil {
		return "", err
	}

	query := "WHERE " + expression
	if additionalExpression != "" {
		query += " AND " + additionalExpression
	}

	return query, nil
}

func BuildWhereExpr(ctx gen.IWhereExprContext, params map[string]any, sources []Source) (string, error) {
	switch true {
	case ctx.IdentifiedExpr() != nil:
		return BuildIdentifiedExpr(ctx.IdentifiedExpr(), params, sources)
	case ctx.AND() != nil:
		leftExpr, err := BuildWhereExpr(ctx.WhereExpr(0), params, sources)
		if err != nil {
			return "", err
		}
		rightExpr, err := BuildWhereExpr(ctx.WhereExpr(1), params, sources)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s) AND (%s)", leftExpr, rightExpr), nil
	case ctx.OR() != nil:
		leftExpr, err := BuildWhereExpr(ctx.WhereExpr(0), params, sources)
		if err != nil {
			return "", err
		}
		rightExpr, err := BuildWhereExpr(ctx.WhereExpr(1), params, sources)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s) OR (%s)", leftExpr, rightExpr), nil
	case ctx.SYM_LEFT_PAREN() != nil:
		return BuildWhereExpr(ctx.WhereExpr(0), params, sources)
	default:
		return "", fmt.Errorf("unsupported where expression")
	}
}

func BuildLimitOffsetClause(ctx gen.ILimitClauseContext, params map[string]any, singleRow bool) (string, error) {
	if ctx == nil {
		if singleRow {
			return "LIMIT 1", nil
		}
		return "", nil
	}

	query := " "
	if ctx.LIMIT() != nil {
		limitOperand, err := BuildLimitOperand(ctx.GetLimit(), params)
		if err != nil {
			return "", err
		}

		if singleRow && limitOperand > "1" {
			limitOperand = "1"
		}

		query += fmt.Sprintf("LIMIT %s ", limitOperand)
	}

	if ctx.OFFSET() != nil {
		offsetOperand, err := BuildLimitOperand(ctx.GetOffset(), params)
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf("OFFSET %s ", offsetOperand)
	}

	return query, nil
}

func BuildLimitOperand(ctx gen.ILimitOperandContext, params map[string]any) (string, error) {
	switch true {
	case ctx.INTEGER() != nil:
		return ctx.INTEGER().GetText(), nil
	case ctx.PARAMETER() != nil:
		return BuildParameter(ctx.PARAMETER(), params, false, "int")
	default:
		return "", fmt.Errorf("unsupported limit operand")
	}
}

func BuildIdentifiedExpr(ctx gen.IIdentifiedExprContext, params map[string]any, sources []Source) (string, error) {
	switch true {
	case ctx.EXISTS() != nil:
		source, path, _, _, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("jsonb_path_exists(%s.data, '%s')", source.Table, path), nil
	case ctx.IdentifiedPath() != nil && ctx.COMPARISON_OPERATOR() != nil:
		source, _, pathWithoutEnd, endsWith, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", err
		}

		comparison := ctx.COMPARISON_OPERATOR().GetText()

		value, err := BuildTerminal(ctx.Terminal(), params, sources)
		if err != nil {
			return "", err
		}

		fastPath, err := BuildFastValueExtractionExpr(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", err
		}
		if fastPath != "" {
			return fmt.Sprintf("(%s %s %s)", fastPath, comparison, value), nil
		}

		// Worst case scenario, try every possible combination
		switchExpression := BuildSlowValueExtractionExpr(endsWith)
		if endsWith == nil {
			return fmt.Sprintf("EXISTS(SELECT 1 FROM %s target WHERE (%s) %s %s)", source.Table, switchExpression, comparison, value), nil
		}

		pathEnding, err := BuildPathEnding(endsWith, params)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("EXISTS(SELECT 1 FROM JSON_TABLE(%s.data, '%s' COLUMNS(data JSONB PATH '$')) parent, JSON_TABLE(parent.data, '%s' COLUMNS(data JSONB PATH '$')) target WHERE (%s) %s %s)", source.Table, pathWithoutEnd, pathEnding, switchExpression, comparison, value), nil
	case ctx.FunctionCall() != nil && ctx.COMPARISON_OPERATOR() != nil:
		return BuildFunctionCall(ctx.FunctionCall(), params, sources)
	case ctx.LIKE() != nil:
		source, path, _, _, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", err
		}

		operand, err := BuildLikeOperand(ctx.LikeOperand(), params)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("EXISTS (SELECT 1 FROM JSON_TABLE(%s.data, '%s' COLUMNS(data JSONB PATH '$')) WHERE data #>> '{}' LIKE %s)", source.Table, path, operand), nil

	case ctx.MATCHES() != nil:
		source, path, _, _, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", err
		}

		values, err := BuildMatchedOperand(ctx.MatchesOperand(), params)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("jsonb_path_query_array(%s.data, '%s') <@ jsonb_build_array(%s)", source.Table, path, values), nil
	case ctx.SYM_LEFT_PAREN() != nil:
		return BuildIdentifiedExpr(ctx.IdentifiedExpr(), params, sources)
	default:
		return "", fmt.Errorf("unsupported identified expression")
	}
}

// convertAQLWildcardToPostgres converts AQL wildcard patterns to PostgreSQL LIKE patterns
// AQL wildcards:
//   - ? matches any single character
//   - * matches any sequence of zero or more characters
//
// PostgreSQL LIKE wildcards:
//   - _ matches any single character
//   - % matches any sequence of zero or more characters
func convertAQLWildcardToPostgres(aqlPattern string) string {
	var result strings.Builder
	escaped := false

	for i, char := range aqlPattern {
		// Handle escape sequences
		if char == '\\' && i+1 < len(aqlPattern) {
			nextChar := rune(aqlPattern[i+1])
			// If backslash is escaping a wildcard or another backslash
			if nextChar == '?' || nextChar == '*' || nextChar == '\\' {
				escaped = true
				continue
			}
		}

		if escaped {
			// Escaped characters are treated literally
			// We need to escape them for PostgreSQL LIKE
			switch char {
			case '?', '*':
				// Escaped AQL wildcards become literal characters
				result.WriteRune('\\')
				result.WriteRune(char)
			case '\\':
				// Escaped backslash
				result.WriteString("\\\\")
			case '%', '_':
				// These need escaping in PostgreSQL LIKE
				result.WriteRune('\\')
				result.WriteRune(char)
			default:
				result.WriteRune(char)
			}
			escaped = false
			continue
		}

		// Convert unescaped wildcards
		switch char {
		case '?':
			// AQL single-character wildcard → PostgreSQL single-character wildcard
			result.WriteRune('_')
		case '*':
			// AQL multi-character wildcard → PostgreSQL multi-character wildcard
			result.WriteRune('%')
		case '%', '_':
			// Escape PostgreSQL special characters that appear literally in the pattern
			result.WriteRune('\\')
			result.WriteRune(char)
		case '\\':
			// Literal backslash in PostgreSQL LIKE
			result.WriteString("\\\\")
		default:
			result.WriteRune(char)
		}
	}

	return result.String()
}

func BuildLikeOperand(ctx gen.ILikeOperandContext, params map[string]any) (string, error) {
	var value string
	switch true {
	case ctx.STRING() != nil:
		value = ctx.STRING().GetText()
		value = value[1 : len(value)-1] // Remove quotes
	case ctx.PARAMETER() != nil:
		paramName := ctx.PARAMETER().GetText()
		paramName = paramName[1:] // Remove leading '$'
		paramValue, ok := params[paramName]
		if !ok {
			return "", fmt.Errorf("missing parameter: %s", paramName)
		}

		strValue, ok := paramValue.(string)
		if !ok {
			return "", fmt.Errorf("parameter %s must be a string for LIKE operation", paramName)
		}
		value = strValue
	default:
		return "", fmt.Errorf("unsupported like operand")
	}

	// Convert AQL wildcards to PostgreSQL LIKE patterns
	// AQL: ? = single character, * = zero or more characters
	// PostgreSQL: _ = single character, % = zero or more characters
	// We also need to escape existing % and _ in the value

	pgPattern := convertAQLWildcardToPostgres(value)

	return fmt.Sprintf("'%s'", pgPattern), nil
}

func BuildMatchedOperand(ctx gen.IMatchesOperandContext, params map[string]any) (string, error) {
	switch true {
	case ctx.AllValueListItem() != nil:
		cases := ""
		for idx, item := range ctx.AllValueListItem() {
			if idx > 0 {
				cases += ", "
			}

			value, err := BuildValueListItem(item, params)
			if err != nil {
				return "", err
			}
			cases += value
		}

		return cases, nil
	case ctx.TerminologyFunction() != nil:
		return "", fmt.Errorf("terminology functions not yet supported in MATCHES")
	case ctx.URI() != nil:
		return "", fmt.Errorf("URI not yet supported in MATCHES")
	default:
		return "", fmt.Errorf("unsupported matches operand")
	}
}

func BuildValueListItem(ctx gen.IValueListItemContext, params map[string]any) (string, error) {
	switch true {
	case ctx.Primitive() != nil:
		return BuildPrimitive(ctx.Primitive(), false)
	case ctx.PARAMETER() != nil:
		return BuildParameter(ctx.PARAMETER(), params, false, "any")
	case ctx.TerminologyFunction() != nil:
		return "", fmt.Errorf("terminology functions not yet supported in value list item")
	default:
		return "", fmt.Errorf("unsupported value list item")
	}
}

func BuildIdentifiedPath(ctx gen.IIdentifiedPathContext, params map[string]any, sources []Source) (Source, string, string, gen.IPathPartContext, error) {
	root := ctx.IDENTIFIER().GetText()

	source := Source{}
	for _, s := range sources {
		if s.Alias == root {
			source = s
			break
		}
	}
	if source.Model == "" {
		return Source{}, "", "", nil, fmt.Errorf("unknown source alias: %s", root)
	}

	path := "$"
	pathWithoutEnd := ""
	var endsWith gen.IPathPartContext
	if ctx.NodePredicate() != nil {
		condition, err := BuildNodePredicate(ctx.NodePredicate(), params)
		if err != nil {
			return Source{}, "", "", nil, err
		}
		path += fmt.Sprintf(" ? (%s)", condition)
	}

	if ctx.ObjectPath() != nil {
		objectPath, objectPathWithoutEnd, objectPathEndsWith, err := BuildObjectPath(ctx.ObjectPath(), params)
		if err != nil {
			return Source{}, "", "", nil, err
		}

		pathWithoutEnd += path + objectPathWithoutEnd
		path += objectPath
		endsWith = objectPathEndsWith
	}

	return source, path, pathWithoutEnd, endsWith, nil
}

func BuildObjectPath(ctx gen.IObjectPathContext, params map[string]any) (string, string, gen.IPathPartContext, error) {
	pathParts := ctx.AllPathPart()

	path := ""
	partialPath := ""
	var endsWith gen.IPathPartContext
	endsAt := len(pathParts) - 1
	for idx, part := range pathParts {
		identifier := part.IDENTIFIER().GetText()
		path += fmt.Sprintf(".%s", identifier)

		if part.NodePredicate() != nil {
			condition, err := BuildNodePredicate(part.NodePredicate(), params)
			if err != nil {
				return "", "", nil, err
			}
			path += fmt.Sprintf(" ? (%s)", condition)
		}

		endsWith = part
		if idx < endsAt {
			partialPath = path
		}
	}

	return path, partialPath, endsWith, nil
}

func BuildNodePredicate(ctx gen.INodePredicateContext, params map[string]any) (string, error) {
	buildSymComma := func() (string, error) {
		switch true {
		case ctx.GetRightAtCode() != nil:
			atCode := ctx.AT_CODE(0).GetText()
			if ctx.AT_CODE(1) != nil {
				atCode = ctx.AT_CODE(1).GetText()
			}

			return fmt.Sprintf(`@.name.defining_code.code_string == "%s" && @.name.defining_code.terminology_id.value == "local"`, atCode), nil
		case ctx.GetRightIdCode() != nil:
			idCode := ctx.ID_CODE(0).GetText()
			if ctx.ID_CODE(1) != nil {
				idCode = ctx.ID_CODE(1).GetText()
			}

			return fmt.Sprintf(`@.name.defining_code.code_string == "%s"`, idCode), nil
		case ctx.GetRightParameter() != nil:
			parameter := ctx.PARAMETER(0)
			if ctx.PARAMETER(1) != nil {
				parameter = ctx.PARAMETER(1)
			}

			value, err := BuildParameter(parameter, params, true, "string")
			if err != nil {
				return "", err
			}

			return fmt.Sprintf(`@.name.defining_code.code_string == %s && @.name.defining_code.terminology_id.value == "local"`, value), nil
		case ctx.STRING() != nil:
			value := ctx.STRING().GetText()
			value = value[1 : len(value)-1] // Remove quotes
			return fmt.Sprintf(`@.name.value == "%s"`, value), nil
		case ctx.TERM_CODE() != nil:
			termCode := ctx.TERM_CODE().GetText()
			termCodeSplit := strings.SplitN(termCode, "::", 2)
			if len(termCodeSplit) != 2 {
				return "", fmt.Errorf("invalid TERM_CODE format: %s", termCode)
			}

			return fmt.Sprintf(`@.name.defining_code.code_string == "%s" && @.name.defining_code.terminology_id.value == "%s"`, termCodeSplit[0], termCodeSplit[1]), nil
		default:
			return "", fmt.Errorf("unsupported SYM_COMMA right operand")
		}
	}

	switch true {
	case ctx.ARCHETYPE_HRID() != nil:
		value := ctx.ARCHETYPE_HRID().GetText()
		query := fmt.Sprintf(`@.archetype_node_id == "%s"`, value)
		if ctx.SYM_COMMA() != nil {
			symCommaExpr, err := buildSymComma()
			if err != nil {
				return "", err
			}
			query += " && " + symCommaExpr
		}
		return query, nil
	case ctx.GetLeftIdCode() != nil:
		idCode := ctx.ID_CODE(0).GetText()
		query := fmt.Sprintf(`@.archetype_node_id == "%s"`, idCode)
		if ctx.SYM_COMMA() != nil {
			symCommaExpr, err := buildSymComma()
			if err != nil {
				return "", err
			}
			query += " && " + symCommaExpr
		}

		return query, nil
	case ctx.GetLeftAtCode() != nil:
		atCode := ctx.AT_CODE(0).GetText()
		query := fmt.Sprintf(`@.archetype_node_id == "%s"`, atCode)
		if ctx.SYM_COMMA() != nil {
			symCommaExpr, err := buildSymComma()
			if err != nil {
				return "", err
			}
			query += " && " + symCommaExpr
		}

		return query, nil
	case ctx.GetLeftParamter() != nil:
		parameter := ctx.PARAMETER(0)
		value, err := BuildParameter(parameter, params, true, "string")
		if err != nil {
			return "", err
		}

		query := fmt.Sprintf("@.archetype_node_id == %s", value)
		if ctx.SYM_COMMA() != nil {
			symCommaExpr, err := buildSymComma()
			if err != nil {
				return "", err
			}
			query += " && " + symCommaExpr
		}

		return query, nil
	case ctx.COMPARISON_OPERATOR() != nil:
		path, _, _, err := BuildObjectPath(ctx.ObjectPath(), params)
		if err != nil {
			return "", err
		}

		comparison := ctx.COMPARISON_OPERATOR().GetText()
		if comparison == "=" {
			comparison = "=="
		}

		condition, err := BuildPathPredicateOperand(ctx.PathPredicateOperand(), params)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("@%s %s %s", path, comparison, condition), nil
	case ctx.MATCHES() != nil:
		return "", fmt.Errorf("MATCHES not yet supported in node predicate")
	case ctx.AND() != nil:
		leftExpr, err := BuildNodePredicate(ctx.NodePredicate(0), params)
		if err != nil {
			return "", err
		}
		rightExpr, err := BuildNodePredicate(ctx.NodePredicate(1), params)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s && %s", leftExpr, rightExpr), nil
	case ctx.OR() != nil:
		leftExpr, err := BuildNodePredicate(ctx.NodePredicate(0), params)
		if err != nil {
			return "", err
		}
		rightExpr, err := BuildNodePredicate(ctx.NodePredicate(1), params)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s || %s", leftExpr, rightExpr), nil
	default:
		return "", fmt.Errorf("unsupported node predicate")
	}
}

func BuildPathPredicateOperand(ctx gen.IPathPredicateOperandContext, params map[string]any) (string, error) {
	switch true {
	case ctx.Primitive() != nil:
		return BuildPrimitive(ctx.Primitive(), true)
	case ctx.ObjectPath() != nil:
		path, _, _, err := BuildObjectPath(ctx.ObjectPath(), params)
		if err != nil {
			return "", err
		}
		return path, nil
	case ctx.PARAMETER() != nil:
		return BuildParameter(ctx.PARAMETER(), params, false, "any")
	case ctx.ID_CODE() != nil:
		value := ctx.ID_CODE().GetText()
		return fmt.Sprintf(`"%s"`, value), nil
	case ctx.AT_CODE() != nil:
		termCode := ctx.AT_CODE().GetText()
		return fmt.Sprintf(`"%s"`, termCode), nil
	default:
		return "", fmt.Errorf("unsupported path predicate operand")
	}
}

func BuildTerminal(ctx gen.ITerminalContext, params map[string]any, sources []Source) (string, error) {
	switch true {
	case ctx.Primitive() != nil:
		return BuildPrimitive(ctx.Primitive(), false)
	case ctx.PARAMETER() != nil:
		return BuildParameter(ctx.PARAMETER(), params, false, "any")
	case ctx.IdentifiedPath() != nil:
		source, path, _, _, err := BuildIdentifiedPath(ctx.IdentifiedPath(), params, sources)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')", source.Table, path), nil
	case ctx.FunctionCall() != nil:
		return BuildFunctionCall(ctx.FunctionCall(), params, sources)
	default:
		return "", fmt.Errorf("unsupported terminal")
	}
}

func BuildFunctionCall(ctx gen.IFunctionCallContext, params map[string]any, sources []Source) (string, error) {
	switch true {
	case ctx.TerminologyFunction() != nil:
		return "", fmt.Errorf("terminology functions not yet supported")
	case ctx.STRING_FUNCTION_ID() != nil:
		switch strings.ToUpper(ctx.STRING_FUNCTION_ID().GetText()) {
		case "LENGTH":
			if len(ctx.AllTerminal()) != 1 {
				return "", fmt.Errorf("LENGTH function requires exactly one argument")
			}

			value, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("jsonb_array_length(%s)", value), nil
		case "POSITION":
			if len(ctx.AllTerminal()) != 2 {
				return "", fmt.Errorf("POSITION function requires exactly two arguments")
			}

			substr, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}
			str, err := BuildTerminal(ctx.Terminal(1), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("POSITION((%s #>> '{}') IN (%s #>> '{}'))", substr, str), nil
		case "SUBSTRING":
			if len(ctx.AllTerminal()) != 3 {
				return "", fmt.Errorf("SUBSTRING function requires exactly three arguments")
			}

			str, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}
			start, err := BuildTerminal(ctx.Terminal(1), params, sources)
			if err != nil {
				return "", err
			}
			length, err := BuildTerminal(ctx.Terminal(2), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("SUBSTRING((%s #>> '{}') FROM (%s)::int FOR (%s)::int)", str, start, length), nil
		case "CONCAT":
			if len(ctx.AllTerminal()) < 2 {
				return "", fmt.Errorf("CONCAT function requires at least two arguments")
			}

			parts := []string{}
			for _, termCtx := range ctx.AllTerminal() {
				part, err := BuildTerminal(termCtx, params, sources)
				if err != nil {
					return "", err
				}
				parts = append(parts, fmt.Sprintf("(%s #>> '{}')", part))
			}

			return fmt.Sprintf("CONCAT(%s)", strings.Join(parts, ", ")), nil
		case "CONCAT_WS":
			if len(ctx.AllTerminal()) < 3 {
				return "", fmt.Errorf("CONCAT_WS function requires at least three arguments")
			}

			separator, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}

			parts := []string{}
			for _, termCtx := range ctx.AllTerminal()[1:] {
				part, err := BuildTerminal(termCtx, params, sources)
				if err != nil {
					return "", err
				}
				parts = append(parts, fmt.Sprintf("(%s #>> '{}')", part))
			}

			return fmt.Sprintf("CONCAT_WS((%s #>> '{}'), %s)", separator, strings.Join(parts, ", ")), nil
		default:
			return "", fmt.Errorf("unsupported function call")
		}
	case ctx.NUMERIC_FUNCTION_ID() != nil:
		switch strings.ToUpper(ctx.NUMERIC_FUNCTION_ID().GetText()) {
		case "ABS":
			if len(ctx.AllTerminal()) != 1 {
				return "", fmt.Errorf("ABS function requires exactly one argument")
			}

			value, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("ABS((%s #>> '{}')::numeric)", value), nil
		case "MOD":
			if len(ctx.AllTerminal()) != 2 {
				return "", fmt.Errorf("MOD function requires exactly two arguments")
			}

			numerator, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}
			denominator, err := BuildTerminal(ctx.Terminal(1), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("MOD((%s #>> '{}')::numeric, (%s)::numeric)", numerator, denominator), nil
		case "CEIL":
			if len(ctx.AllTerminal()) != 1 {
				return "", fmt.Errorf("CEIL function requires exactly one argument")
			}

			value, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("CEIL((%s #>> '{}')::numeric)", value), nil
		case "FLOOR":
			if len(ctx.AllTerminal()) != 1 {
				return "", fmt.Errorf("FLOOR function requires exactly one argument")
			}

			value, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("FLOOR((%s #>> '{}')::numeric)", value), nil
		case "ROUND":
			if len(ctx.AllTerminal()) != 2 {
				return "", fmt.Errorf("ROUND function requires exactly two arguments")
			}

			value, err := BuildTerminal(ctx.Terminal(0), params, sources)
			if err != nil {
				return "", err
			}
			precision, err := BuildTerminal(ctx.Terminal(1), params, sources)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("ROUND((%s #>> '{}')::numeric, (%s)::int)", value, precision), nil
		default:
			return "", fmt.Errorf("unsupported function call")
		}
	case ctx.DATE_TIME_FUNCTION_ID() != nil:
		switch strings.ToUpper(ctx.DATE_TIME_FUNCTION_ID().GetText()) {
		case "CURRENT_DATE":
			if len(ctx.AllTerminal()) != 0 {
				return "", fmt.Errorf("CURRENT_DATE function requires no arguments")
			}

			return "CURRENT_DATE", nil
		case "CURRENT_TIME":
			if len(ctx.AllTerminal()) != 0 {
				return "", fmt.Errorf("CURRENT_TIME function requires no arguments")
			}

			return "CURRENT_TIME", nil
		case "CURRENT_DATETIME":
			if len(ctx.AllTerminal()) != 0 {
				return "", fmt.Errorf("CURRENT_DATETIME function requires no arguments")
			}

			return "NOW()", nil
		case "NOW":
			if len(ctx.AllTerminal()) != 0 {
				return "", fmt.Errorf("NOW function requires no arguments")
			}

			return "to_jsonb(now())", nil
		case "CURRENT_TIMEZONE":
			if len(ctx.AllTerminal()) != 0 {
				return "", fmt.Errorf("CURRENT_TIMEZONE function requires no arguments")
			}

			return "SELECT to_char (NOW()::timestamp, 'OF')", nil // (-)hh:mm
		default:
			return "", fmt.Errorf("unsupported function call")
		}
	default:
		return "", fmt.Errorf("unsupported function call")
	}
}

func BuildPrimitive(ctx gen.IPrimitiveContext, jsonPathCompatible bool) (string, error) {
	switch true {
	case ctx.STRING() != nil:
		value := ctx.STRING().GetText()
		value = value[1 : len(value)-1] // Remove quotes
		if jsonPathCompatible {
			return fmt.Sprintf(`"%s"`, value), nil
		}
		return fmt.Sprintf(`'"%s"'::jsonb`, value), nil
	case ctx.NumericPrimitive() != nil:
		value := ctx.NumericPrimitive().GetText()
		if jsonPathCompatible {
			return value, nil
		}
		return fmt.Sprintf("'%s'::jsonb", value), nil
	case ctx.DATE() != nil:
		value := ctx.DATE().GetText()
		value = value[1 : len(value)-1] // Remove quotes
		if jsonPathCompatible {
			return fmt.Sprintf(`"%s"`, value), nil
		}
		return fmt.Sprintf("to_jsonb('%s'::date)", value), nil
	case ctx.TIME() != nil:
		value := ctx.TIME().GetText()
		value = value[1 : len(value)-1] // Remove quotes
		if jsonPathCompatible {
			return fmt.Sprintf(`"%s"`, value), nil
		}
		return fmt.Sprintf("to_jsonb('%s'::timetz)", value), nil
	case ctx.DATETIME() != nil:
		value := ctx.DATETIME().GetText()
		value = value[1 : len(value)-1] // Remove quotes
		if jsonPathCompatible {
			return fmt.Sprintf(`"%s"`, value), nil
		}
		return fmt.Sprintf("to_jsonb('%s'::timestamptz)", value), nil
	case ctx.BOOLEAN() != nil:
		value := ctx.BOOLEAN().GetText()
		if jsonPathCompatible {
			return value, nil
		}
		return fmt.Sprintf("to_jsonb(%s)", value), nil
	case ctx.NULL() != nil:
		if jsonPathCompatible {
			return "null", nil
		}
		return "to_jsonb(NULL)", nil
	default:
		return "", fmt.Errorf("unsupported primitive")
	}
}

func BuildParameter(ctx antlr.TerminalNode, params map[string]any, jsonPathCompatible bool, expectedType string) (string, error) {
	paramName := ctx.GetText()
	paramName = paramName[1:] // Remove leading '$'
	paramValue, ok := params[paramName]
	if !ok {
		return "", fmt.Errorf("missing parameter: %s", paramName)
	}

	switch paramValue.(type) {
	case string:
		if expectedType != "string" && expectedType != "any" {
			return "", fmt.Errorf("parameter %s expected to be of type %s, got string", paramName, expectedType)
		}

		if jsonPathCompatible {
			return fmt.Sprintf(`"%v"`, paramValue), nil
		}
		return fmt.Sprintf(`'"%v"'::jsonb`, paramValue), nil
	case int, int32, int64, float32, float64:
		if expectedType != "number" && expectedType != "any" {
			return "", fmt.Errorf("parameter %s expected to be of type %s, got number", paramName, expectedType)
		}

		if jsonPathCompatible {
			return fmt.Sprintf("%v", paramValue), nil
		}
		return fmt.Sprintf("'%v'::jsonb", paramValue), nil
	case bool:
		if expectedType != "boolean" && expectedType != "any" {
			return "", fmt.Errorf("parameter %s expected to be of type %s, got boolean", paramName, expectedType)
		}

		if jsonPathCompatible {
			return fmt.Sprintf("%v", paramValue), nil
		}
		return fmt.Sprintf("to_jsonb(%v)", paramValue), nil
	default:
		return "", fmt.Errorf("unsupported parameter type for %s", paramName)
	}
}

func BuildSlowValueExtractionExpr(ctx gen.IPathPartContext) string {
	expression := "CASE "

	expression += `
		WHEN target.data @> \'{"_type": "DV_DURATION"}\' THEN to_jsonb(now() + (target.data ->> \'value\')::interval)
		WHEN target.data @> \'{"_type": "DV_ELEMENT"}\' THEN target.data -> \'magnitude\' 
		WHEN target.data @> \'{"_type": "DV_QUANTITY"}\' THEN target.data -> \'magnitude\' 
		WHEN target.data @> \'{"_type": "DV_COUNT"}\' THEN target.data -> \'magnitude\' 
		WHEN target.data @> \'{"_type": "DV_DATE_TIME"}\' THEN to_jsonb((target.data ->> \'value\')::timestamptz)
		WHEN target.data @> \'{"_type": "DV_TIME"}\' THEN to_jsonb((target.data ->> \'value\')::timetz AT TIME ZONE \'UTC\')
		WHEN target.data @> \'{"_type": "DV_DATE"}\' THEN to_jsonb((target.data ->> \'value\')::date)
		WHEN target.data @> \'{"_type": "DV_ORDINAL"}\' THEN target.data -> \'value\' 
		WHEN target.data @> \'{"_type": "DV_PROPORTION"}\' THEN to_jsonb((target.data ->> \'numerator\')::numeric / (target.data ->> \'denominator\')::numeric) 
	`

	identifier := ctx.IDENTIFIER().GetText()
	if identifier == "value" {
		expression += `
			WHEN parent.data @> \'{"_type": "DV_DURATION"}\' THEN to_jsonb(now() + (target.data #>> \'{}\')::interval)
			WHEN parent.data @> \'{"_type": "DV_DATE_TIME"}\' THEN to_jsonb((target.data #>> \'{}\')::timestamptz)
			WHEN parent.data @> \'{"_type": "DV_TIME"}\' THEN to_jsonb((target.data #>> \'{}\')::timetz AT TIME ZONE \'UTC\')
			WHEN parent.data @> \'{"_type": "DV_DATE"}\' THEN to_jsonb((target.data #>> \'{}\')::date)
		`
	}

	expression += "ELSE target.data END"

	return expression
}

func BuildPathEnding(ctx gen.IPathPartContext, params map[string]any) (string, error) {
	path := fmt.Sprintf("$.%s", ctx.IDENTIFIER().GetText())
	if ctx.NodePredicate() != nil {
		condition, err := BuildNodePredicate(ctx.NodePredicate(), params)
		if err != nil {
			return "", err
		}
		path += fmt.Sprintf(" ? (%s)", condition)
	}

	return path, nil
}

func BuildFastValueExtractionExpr(ctx gen.IIdentifiedPathContext, params map[string]any, sources []Source) (string, error) {
	source, path, _, _, err := BuildIdentifiedPath(ctx, params, sources)
	if err != nil {
		return "", err
	}

	relatedModels := ModelInheritanceTable(source.Model)

	// Construct path
	var plainPath string
	for idx, part := range ctx.ObjectPath().AllPathPart() {
		if idx > 0 {
			plainPath += "/"
		}

		plainPath += part.IDENTIFIER().GetText()
	}

	switch true {
	case slices.Contains(relatedModels, openehr.EHR_MODEL_NAME):
		switch plainPath {
		case "system_id/value", "ehr_id/value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.COMPOSITION_MODEL_NAME):
		switch plainPath {
		case "start_time", "end_time":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.value')", source.Table, path), nil
		case "start_time/value", "end_time/value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_TIME_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.value')::text::timetz", source.Table, path), nil
		case "value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::timetz", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_DATE_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.value')::text::date", source.Table, path), nil
		case "value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::date", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_DATE_TIME_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.value')::text::timestamptz", source.Table, path), nil
		case "value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::timestamptz", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_DURATION_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.value')::text::interval", source.Table, path), nil
		case "value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::interval", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_PROPORTION_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.numerator')::text::float / jsonb_path_query_first(%s.data, '%s.denominator')::text::float", source.Table, path, source.Table, path), nil
		case "numerator", "denominator":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::float", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_ORDINAL_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.value')::text::int", source.Table, path), nil
		case "value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::int", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_COUNT_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.magnitude')::text::int", source.Table, path), nil
		case "magnitude":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::int", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_QUANTITY_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.magnitude')::text::float", source.Table, path), nil
		case "magnitude":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::float", source.Table, path), nil
		}
	case slices.Contains(relatedModels, openehr.DV_SCALE_MODEL_NAME):
		switch plainPath {
		case "":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s.value')::text::float", source.Table, path), nil
		case "value":
			return fmt.Sprintf("jsonb_path_query_first(%s.data, '%s')::text::float", source.Table, path), nil
		}
	}

	return "", nil
}

func ModelInheritanceTable(model string) []string {
	models := []string{}
	queue := []string{model}
	for len(queue) > 0 {
		currModel := queue[0]
		queue = queue[1:]

		models = append(models, currModel)

		children, ok := MODEL_INHERITANCE_TABLE[currModel]
		if ok {
			queue = append(queue, children...)
		}
	}

	return models
}

func InheritanceTableTableReverse(model string) []string {
	models := []string{}
	queue := []string{model}
	for len(queue) > 0 {
		currModel := queue[0]
		queue = queue[1:]

		models = append(models, currModel)

		for parent, children := range MODEL_INHERITANCE_TABLE {
			for _, child := range children {
				if child == currModel {
					queue = append(queue, parent)
				}
			}
		}
	}

	return models
}

var MODEL_INHERITANCE_TABLE = map[string][]string{
	"AGENT":                   {},
	"ORGANISATION":            {},
	"GROUP":                   {},
	"PERSON":                  {},
	"ACTOR":                   {"AGENT", "ORGANISATION", "GROUP", "PERSON"},
	"PARTY_IDENTITY":          {},
	"PARTY_RELATIONSHIP":      {},
	"ROLE":                    {},
	"PARTY":                   {"ROLE", "ACTOR"},
	"CONTACT":                 {},
	"ADDRESS":                 {},
	"CAPABILITY":              {},
	"VERSIONED_PARTY":         {},
	"EHR":                     {},
	"VERSIONED_EHR_ACCESS":    {},
	"EHR_ACCESS":              {},
	"ACCESS_CONTROL_SETTINGS": {},
	"VERSIONED_EHR_STATUS":    {},
	"EHR_STATUS":              {},
	"VERSIONED_COMPOSITION":   {},
	"COMPOSITION":             {},
	"EVENT_CONTEXT":           {},
	"CONTENT_ITEM":            {"GENERIC_ENTRY", "SECTION", "ENTRY"},
	"SECTION":                 {},
	"ENTRY":                   {"ADMIN_ENTRY", "CARE_ENTRY"},
	"ADMIN_ENTRY":             {},
	"CARE_ENTRY":              {"OBSERVATION", "EVALUATION", "INSTRUCTION", "ACTION"},
	"OBSERVATION":             {},
	"EVALUATION":              {},
	"INSTRUCTION":             {},
	"ACTIVITY":                {},
	"ACTION":                  {},
	"INSTRUCTION_DETAILS":     {},
	"ISM_TRANSITION":          {},
	"PATHABLE":                {"LOCATABLE", "EVENT_CONTEXT", "INSTRUCTION_DETAILS", "ISM_TRANSITION"},
	"LOCATABLE":               {"EHR_ACCESS", "EHR_STATUS", "COMPOSITION", "CONTENT_ITEM", "ACTIVITY", "FOLDER", "DATA_STRUCTURE", "ITEM", "EVENT", "PARTY", "PARTY_RELATIONSHIP", "CONTACT", "ADDRESS", "CAPABILITY"},
	"ARCHETYPED":              {},
	"LINK":                    {},
	"FEEDER_AUDIT":            {},
	"FEEDER_AUDIT_DETAILS":    {},
	"PARTY_PROXY":             {"PARTY_SELF", "PARTY_IDENTIFIED"},
	"PARTY_SELF":              {},
	"PARTY_IDENTIFIED":        {"PARTY_RELATED", "PARTICIPATION"},
	"PARTY_RELATED":           {},
	"PARTICIPATION":           {},
	"AUDIT_DETAILS":           {"ATTESTATION"},
	"ATTESTATION":             {},
	"REVISION_HISTORY":        {},
	"REVISION_HISTORY_ITEM":   {},
	"VERSIONED_FOLDER":        {},
	"FOLDER":                  {},
	"VERSIONED_OBJECT":        {"VERSIONED_EHR_ACCESS", "VERSIONED_EHR_STATUS", "VERSIONED_COMPOSITION", "VERSIONED_FOLDER", "VERSIONED_PARTY"},
	"VERSION":                 {"ORIGINAL_VERSION", "IMPORTED_VERSION"},
	"ORIGINAL_VERSION":        {},
	"IMPORTED_VERSION":        {},
	"CONTRIBUTION":            {},
	//        "ITEM_TAG": {"ITEM_TAG"},
	//        "AUTHORED_RESOURCE": {},
	"TRANSLATION_DETAILS":       {},
	"RESOURCE_DESCRIPTION":      {},
	"RESOURCE_DESCRIPTION_ITEM": {},
	"DATA_STRUCTURE":            {"ITEM_SINGLE", "ITEM_LIST", "ITEM_TABLE", "ITEM_TREE", "HISTORY"},
	"ITEM_SINGLE":               {},
	"ITEM_LIST":                 {},
	"ITEM_TABLE":                {},
	"ITEM_TREE":                 {},
	"ITEM":                      {"CLUSTER", "ELEMENT"},
	"CLUSTER":                   {},
	"ELEMENT":                   {},
	"HISTORY":                   {},
	"EVENT":                     {"POINT_EVENT", "INTERVAL_EVENT"},
	"POINT_EVENT":               {},
	"INTERVAL_EVENT":            {},
	"DATA_VALUE":                {"DV_BOOLEAN", "DV_STATE", "DV_IDENTIFIER", "DV_TEXT", "DV_PARAGRAPH", "DV_ORDERED", "DV_INTERVAL"},
	"DV_BOOLEAN":                {},
	"DV_STATE":                  {},
	"DV_IDENTIFIER":             {},
	"DV_TEXT":                   {"DV_CODED_TEXT"},
	"TERM_MAPPING":              {},
	"CODE_PHRASE":               {},
	"DV_PARAGRAPH":              {},
	"DV_ORDERED":                {"DV_ORDINAL", "DV_SCALE", "DV_QUANTIFIED"},
	"DV_INTERVAL":               {},
	"REFERENCE_RANGE":           {},
	"DV_ORDINAL":                {},
	"DV_SCALE":                  {},
	"DV_QUANTIFIED":             {"DV_AMOUNT", "DV_ABSOLUTE_QUANTITY"},
	"DV_AMOUNT":                 {"DV_QUANTITY", "DV_COUNT", "DV_PROPORTION", "DV_DURATION"},
	"DV_QUANTITY":               {},
	"DV_COUNT":                  {},
	"DV_PROPORTION":             {},
	//        "PROPORTION_KIND": {"PROPORTION_KIND"},
	"DV_ABSOLUTE_QUANTITY":           {"DV_TEMPORAL"},
	"DV_TEMPORAL":                    {"DV_DATE", "DV_TIME", "DV_DATE_TIME"},
	"DV_DATE":                        {},
	"DV_TIME":                        {},
	"DV_DATE_TIME":                   {},
	"DV_DURATION":                    {},
	"DV_TIME_SPECIFICATION":          {"DV_PERIODIC_TIME_SPECIFICATION", "DV_GENERAL_TIME_SPECIFICATION"},
	"DV_PERIODIC_TIME_SPECIFICATION": {},
	"DV_GENERAL_TIME_SPECIFICATION":  {},
	"GENERIC_ENTRY":                  {},
	"UID":                            {"ISO_OID", "UUID", "INTERNET_ID"},
	"ISO_OID":                        {},
	"UUID":                           {},
	"INTERNET_ID":                    {},
	"OBJECT_ID":                      {"UID_BASED_ID", "ARCHETYPE_ID", "TEMPLATE_ID", "TERMINOLOGY_ID", "GENERIC_ID"},
	"UID_BASED_ID":                   {"HIER_OBJECT_ID", "OBJECT_VERSION_ID"},
	"HIER_OBJECT_ID":                 {},
	"OBJECT_VERSION_ID":              {},
	"VERSION_TREE_ID":                {},
	"ARCHETYPE_ID":                   {},
	"TEMPLATE_ID":                    {},
	"TERMINOLOGY_ID":                 {},
	"GENERIC_ID":                     {},
	"OBJECT_REF":                     {"PARTY_REF", "LOCATABLE_REF"},
	"PARTY_REF":                      {},
	"LOCATABLE_REF":                  {},
}
