package generator

import (
	"fmt"
	"go/format"
	"reflect"
	"strings"
	"text/template"
)

// CodeGenerator generates repository implementation code
type CodeGenerator struct {
	analyzer *Analyzer
	entityType reflect.Type
	tableName  string
	fieldToColumn map[string]string
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator(entityType reflect.Type) (*CodeGenerator, error) {
	analyzer, err := NewAnalyzer(entityType)
	if err != nil {
		return nil, err
	}

	// Get table name from entity
	tableName := toSnakeCase(entityType.Name())
	if entityType.Kind() == reflect.Ptr {
		tableName = toSnakeCase(entityType.Elem().Name())
	}

	// Build field to column mapping
	fieldToColumn := make(map[string]string)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			dbTag = toSnakeCase(field.Name)
		}
		fieldToColumn[field.Name] = dbTag
	}

	return &CodeGenerator{
		analyzer:      analyzer,
		entityType:    entityType,
		tableName:     tableName,
		fieldToColumn: fieldToColumn,
	}, nil
}

// GenerateMethod generates code for a single query method
func (g *CodeGenerator) GenerateMethod(method *QueryMethod, entityName string, idType string) (string, error) {
	tmpl := `func (r *{{.RepositoryName}}) {{.MethodName}}(ctx context.Context{{.Params}}) {{.Returns}} {
	{{.Body}}
}
`

	// Build parameters string
	var params []string
	for _, param := range method.Parameters {
		params = append(params, fmt.Sprintf("%s %s", param.Name, param.Type))
	}
	paramsStr := ""
	if len(params) > 0 {
		paramsStr = ", " + strings.Join(params, ", ")
	}

	// Build returns string
	var returns []string
	switch method.ReturnType {
	case ReturnSingle:
		returns = []string{fmt.Sprintf("*%s", entityName), "error"}
	case ReturnSlice:
		returns = []string{fmt.Sprintf("[]*%s", entityName), "error"}
	case ReturnInt64:
		returns = []string{"int64", "error"}
	case ReturnBool:
		returns = []string{"bool", "error"}
	}
	returnsStr := "(" + strings.Join(returns, ", ") + ")"

	// Generate method body
	body := g.generateMethodBody(method, entityName)

	data := map[string]interface{}{
		"RepositoryName": fmt.Sprintf("%sRepository", entityName),
		"MethodName":     method.Name,
		"Params":         paramsStr,
		"Returns":        returnsStr,
		"Body":           body,
	}

	t, err := template.New("method").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	// Format the generated code
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return buf.String(), nil // Return unformatted if formatting fails
	}

	return string(formatted), nil
}

// generateMethodBody generates the body of a query method
func (g *CodeGenerator) generateMethodBody(method *QueryMethod, entityName string) string {
	var body strings.Builder

	// Generate SQL query to extract WHERE clause
	fullQuery := method.ToSQL(g.tableName, func(fieldName string) string {
		return g.fieldToColumn[fieldName]
	})

	// Extract WHERE clause from full query
	wherePart := ""
	if idx := strings.Index(fullQuery, "WHERE"); idx > 0 {
		wherePart = fullQuery[idx+6:] // Skip "WHERE "
		// Remove ORDER BY and LIMIT if present
		if orderIdx := strings.Index(wherePart, " ORDER BY"); orderIdx > 0 {
			wherePart = wherePart[:orderIdx]
		}
		if limitIdx := strings.Index(wherePart, " LIMIT"); limitIdx > 0 {
			wherePart = wherePart[:limitIdx]
		}
	}

	// Build query based on operation
	var query string
	switch method.Operation {
	case OpFind:
		query = fmt.Sprintf("SELECT * FROM %s", g.tableName)
		if wherePart != "" {
			query += " WHERE " + wherePart
		}
		if len(method.SortFields) > 0 {
			orderClauses := make([]string, len(method.SortFields))
			for i, sf := range method.SortFields {
				orderClauses[i] = fmt.Sprintf("%s %s", g.fieldToColumn[sf.FieldName], sf.Direction)
			}
			query += " ORDER BY " + strings.Join(orderClauses, ", ")
		}
		if method.Limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", method.Limit)
		}
	case OpCount:
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s", g.tableName)
		if wherePart != "" {
			query += " WHERE " + wherePart
		}
	case OpExists:
		query = fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s", g.tableName)
		if wherePart != "" {
			query += " WHERE " + wherePart
		}
		query += ")"
	case OpDelete:
		query = fmt.Sprintf("DELETE FROM %s", g.tableName)
		if wherePart != "" {
			query += " WHERE " + wherePart
		}
	}

	// Build args list for logging and query execution
	argsList := make([]string, 0)
	for _, field := range method.Fields {
		switch field.Operator {
		case OpBetween:
			argsList = append(argsList, fmt.Sprintf("min%s", field.FieldName))
			argsList = append(argsList, fmt.Sprintf("max%s", field.FieldName))
		case OpIn, OpNotIn:
			argsList = append(argsList, fmt.Sprintf("%ss", strings.ToLower(field.FieldName)))
		case OpIsNull, OpIsNotNull, OpTrue, OpFalse:
			// No arguments
		case OpContaining:
			paramName := strings.ToLower(field.FieldName)
			argsList = append(argsList, fmt.Sprintf(`fmt.Sprintf("%%s%%", %s)`, paramName))
		case OpStartingWith:
			paramName := strings.ToLower(field.FieldName)
			argsList = append(argsList, fmt.Sprintf(`fmt.Sprintf("%%s%%", %s)`, paramName))
		case OpEndingWith:
			paramName := strings.ToLower(field.FieldName)
			argsList = append(argsList, fmt.Sprintf(`fmt.Sprintf("%%s%%", %s)`, paramName))
		default:
			argsList = append(argsList, strings.ToLower(field.FieldName))
		}
	}

	argsStr := ""
	if len(argsList) > 0 {
		argsStr = ", " + strings.Join(argsList, ", ")
	}

	// Generate body based on operation and return type
	switch method.Operation {
	case OpFind:
		if method.ReturnType == ReturnSingle {
			body.WriteString(fmt.Sprintf(`query := %q
	r.logQuery(query, []interface{}{%s})

	var row pgx.Row
	if r.tx != nil {
		row = r.tx.tx.QueryRow(ctx, query%s)
	} else {
		row = r.db.pool.QueryRow(ctx, query%s)
	}

	result := new(%s)
	if err := r.scanRow(row, result); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return result, nil`, query, strings.Join(argsList, ", "), argsStr, argsStr, entityName))
		} else {
			body.WriteString(fmt.Sprintf(`query := %q
	r.logQuery(query, []interface{}{%s})

	var rows pgx.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.tx.Query(ctx, query%s)
	} else {
		rows, err = r.db.pool.Query(ctx, query%s)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)`, query, strings.Join(argsList, ", "), argsStr, argsStr))
		}
	case OpCount:
		body.WriteString(fmt.Sprintf(`query := %q
	r.logQuery(query, []interface{}{%s})

	var count int64
	var err error
	if r.tx != nil {
		err = r.tx.tx.QueryRow(ctx, query%s).Scan(&count)
	} else {
		err = r.db.pool.QueryRow(ctx, query%s).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil`, query, strings.Join(argsList, ", "), argsStr, argsStr))
	case OpExists:
		body.WriteString(fmt.Sprintf(`query := %q
	r.logQuery(query, []interface{}{%s})

	var exists bool
	var err error
	if r.tx != nil {
		err = r.tx.tx.QueryRow(ctx, query%s).Scan(&exists)
	} else {
		err = r.db.pool.QueryRow(ctx, query%s).Scan(&exists)
	}

	if err != nil {
		return false, err
	}

	return exists, nil`, query, strings.Join(argsList, ", "), argsStr, argsStr))
	case OpDelete:
		body.WriteString(fmt.Sprintf(`query := %q
	r.logQuery(query, []interface{}{%s})

	var result pgconn.CommandTag
	var err error
	if r.tx != nil {
		result, err = r.tx.tx.Exec(ctx, query%s)
	} else {
		result, err = r.db.pool.Exec(ctx, query%s)
	}

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil`, query, strings.Join(argsList, ", "), argsStr, argsStr))
	}

	return body.String()
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		if r >= 'A' && r <= 'Z' {
			result.WriteRune(r + 32) // Convert to lowercase
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

