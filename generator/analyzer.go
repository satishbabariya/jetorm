package generator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// QueryMethod represents a parsed query method
type QueryMethod struct {
	Name           string
	Operation      Operation
	Fields         []FieldCondition
	SortFields     []SortField
	Limit          int
	ReturnType     ReturnType
	Parameters     []Parameter
	GeneratedSQL   string
}

// Operation represents the type of query operation
type Operation int

const (
	OpFind Operation = iota
	OpCount
	OpExists
	OpDelete
)

// ReturnType represents the return type of a method
type ReturnType int

const (
	ReturnSingle ReturnType = iota
	ReturnSlice
	ReturnInt64
	ReturnBool
	ReturnError
)

// FieldCondition represents a condition on a field
type FieldCondition struct {
	FieldName string
	Operator  Operator
	AndOr     string // "AND" or "OR"
}

// Operator represents a comparison operator
type Operator int

const (
	OpEqual Operator = iota
	OpNotEqual
	OpGreaterThan
	OpGreaterThanEqual
	OpLessThan
	OpLessThanEqual
	OpLike
	OpNotLike
	OpIn
	OpNotIn
	OpIsNull
	OpIsNotNull
	OpBetween
	OpContaining
	OpStartingWith
	OpEndingWith
	OpIgnoreCase
	OpTrue
	OpFalse
)

// SortField represents a sort field
type SortField struct {
	FieldName string
	Direction string // "ASC" or "DESC"
}

// Parameter represents a method parameter
type Parameter struct {
	Name string
	Type string
}

// Analyzer analyzes method names and generates query methods
type Analyzer struct {
	entityType reflect.Type
	fields     map[string]reflect.StructField
}

// NewAnalyzer creates a new analyzer for an entity type
func NewAnalyzer(entityType reflect.Type) (*Analyzer, error) {
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	if entityType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity type must be a struct")
	}

	fields := make(map[string]reflect.StructField)
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		fields[field.Name] = field
	}

	return &Analyzer{
		entityType: entityType,
		fields:     fields,
	}, nil
}

// AnalyzeMethod analyzes a method name and returns a QueryMethod
func (a *Analyzer) AnalyzeMethod(methodName string) (*QueryMethod, error) {
	method := &QueryMethod{
		Name: methodName,
	}

	// Determine operation type
	if strings.HasPrefix(methodName, "Find") {
		method.Operation = OpFind
	} else if strings.HasPrefix(methodName, "Count") {
		method.Operation = OpCount
	} else if strings.HasPrefix(methodName, "Exists") {
		method.Operation = OpExists
	} else if strings.HasPrefix(methodName, "Delete") {
		method.Operation = OpDelete
	} else {
		return nil, fmt.Errorf("unsupported method prefix: %s", methodName)
	}

	// Parse method name
	remaining := methodName
	var err error

	// Check for First/TopN
	if strings.HasPrefix(remaining, "FindFirst") {
		method.Limit = 1
		remaining = strings.TrimPrefix(remaining, "FindFirst")
	} else if strings.HasPrefix(remaining, "FindTop") {
		// Extract number from FindTop{N}
		re := regexp.MustCompile(`^FindTop(\d+)`)
		matches := re.FindStringSubmatch(remaining)
		if len(matches) > 1 {
			fmt.Sscanf(matches[1], "%d", &method.Limit)
			remaining = strings.TrimPrefix(remaining, matches[0])
		}
	} else if strings.HasPrefix(remaining, "Find") {
		remaining = strings.TrimPrefix(remaining, "Find")
	} else if strings.HasPrefix(remaining, "Count") {
		remaining = strings.TrimPrefix(remaining, "Count")
	} else if strings.HasPrefix(remaining, "Exists") {
		remaining = strings.TrimPrefix(remaining, "Exists")
	} else if strings.HasPrefix(remaining, "Delete") {
		remaining = strings.TrimPrefix(remaining, "Delete")
	}

	// Parse "By" conditions
	if strings.HasPrefix(remaining, "By") {
		remaining = strings.TrimPrefix(remaining, "By")
		// Check if there's an OrderBy clause before parsing conditions
		orderByPos := strings.Index(remaining, "OrderBy")
		if orderByPos > 0 {
			// Parse conditions up to OrderBy
			conditionsPart := remaining[:orderByPos]
			_, err = a.parseConditions(conditionsPart, method)
			if err != nil {
				return nil, err
			}
			remaining = remaining[orderByPos:]
		} else {
			remaining, err = a.parseConditions(remaining, method)
			if err != nil {
				return nil, err
			}
		}
	}

	// Parse OrderBy clause
	if strings.HasPrefix(remaining, "OrderBy") {
		remaining = strings.TrimPrefix(remaining, "OrderBy")
		remaining, err = a.parseOrderBy(remaining, method)
		if err != nil {
			return nil, err
		}
	}

	// Determine return type based on operation
	switch method.Operation {
	case OpFind:
		if method.Limit == 1 {
			method.ReturnType = ReturnSingle
		} else {
			method.ReturnType = ReturnSlice
		}
	case OpCount:
		method.ReturnType = ReturnInt64
	case OpExists:
		method.ReturnType = ReturnBool
	case OpDelete:
		method.ReturnType = ReturnInt64
	}

	// Generate parameters based on conditions
	method.Parameters = a.generateParameters(method)

	return method, nil
}

// parseConditions parses field conditions from method name
func (a *Analyzer) parseConditions(remaining string, method *QueryMethod) (string, error) {
	firstField := true
	for remaining != "" {
		// Check for And/Or (but not on first field)
		andOr := ""
		if !firstField {
			if strings.HasPrefix(remaining, "And") {
				andOr = "AND"
				remaining = strings.TrimPrefix(remaining, "And")
			} else if strings.HasPrefix(remaining, "Or") {
				andOr = "OR"
				remaining = strings.TrimPrefix(remaining, "Or")
			} else if strings.HasPrefix(remaining, "OrderBy") {
				// Stop parsing conditions, OrderBy follows
				break
			}
		}

		// Parse field name and operator
		fieldName, operator, consumed, err := a.parseFieldCondition(remaining)
		if err != nil {
			return remaining, err
		}

		// Validate field exists
		if _, exists := a.fields[fieldName]; !exists {
			return remaining, fmt.Errorf("field %s not found in entity", fieldName)
		}

		method.Fields = append(method.Fields, FieldCondition{
			FieldName: fieldName,
			Operator:  operator,
			AndOr:     andOr,
		})

		remaining = remaining[consumed:]
		firstField = false
	}

	return remaining, nil
}

// parseFieldCondition parses a single field condition
// It stops at "And", "Or", or "OrderBy" to allow proper parsing of multiple conditions
func (a *Analyzer) parseFieldCondition(remaining string) (fieldName string, operator Operator, consumed int, err error) {
	// Find where the field condition ends (at And, Or, or OrderBy)
	endPos := len(remaining)
	if andPos := strings.Index(remaining, "And"); andPos > 0 && andPos < endPos {
		endPos = andPos
	}
	if orPos := strings.Index(remaining, "Or"); orPos > 0 && orPos < endPos {
		endPos = orPos
	}
	if orderByPos := strings.Index(remaining, "OrderBy"); orderByPos > 0 && orderByPos < endPos {
		endPos = orderByPos
	}

	// Extract the field condition part
	fieldPart := remaining[:endPos]

	// Try to match field name with various operators (in order of specificity)
	patterns := []struct {
		pattern  *regexp.Regexp
		operator Operator
	}{
		{regexp.MustCompile(`^(\w+)GreaterThanEqual$`), OpGreaterThanEqual},
		{regexp.MustCompile(`^(\w+)LessThanEqual$`), OpLessThanEqual},
		{regexp.MustCompile(`^(\w+)GreaterThan$`), OpGreaterThan},
		{regexp.MustCompile(`^(\w+)LessThan$`), OpLessThan},
		{regexp.MustCompile(`^(\w+)Containing$`), OpContaining},
		{regexp.MustCompile(`^(\w+)StartingWith$`), OpStartingWith},
		{regexp.MustCompile(`^(\w+)EndingWith$`), OpEndingWith},
		{regexp.MustCompile(`^(\w+)NotLike$`), OpNotLike},
		{regexp.MustCompile(`^(\w+)Like$`), OpLike},
		{regexp.MustCompile(`^(\w+)NotIn$`), OpNotIn},
		{regexp.MustCompile(`^(\w+)In$`), OpIn},
		{regexp.MustCompile(`^(\w+)IsNotNull$`), OpIsNotNull},
		{regexp.MustCompile(`^(\w+)IsNull$`), OpIsNull},
		{regexp.MustCompile(`^(\w+)Between$`), OpBetween},
		{regexp.MustCompile(`^(\w+)IgnoreCase$`), OpIgnoreCase},
		{regexp.MustCompile(`^(\w+)True$`), OpTrue},
		{regexp.MustCompile(`^(\w+)False$`), OpFalse},
		{regexp.MustCompile(`^(\w+)$`), OpEqual}, // Default to equal
	}

	for _, p := range patterns {
		matches := p.pattern.FindStringSubmatch(fieldPart)
		if len(matches) > 1 {
			fieldName = matches[1]
			operator = p.operator
			consumed = len(fieldPart)
			return
		}
	}

	return "", OpEqual, 0, fmt.Errorf("could not parse field condition from: %s", fieldPart)
}

// parseOrderBy parses OrderBy clause
func (a *Analyzer) parseOrderBy(remaining string, method *QueryMethod) (string, error) {
	// Parse field name
	re := regexp.MustCompile(`^(\w+)(Asc|Desc)`)
	matches := re.FindStringSubmatch(remaining)
	if len(matches) < 3 {
		return remaining, fmt.Errorf("invalid OrderBy format: %s", remaining)
	}

	fieldName := matches[1]
	direction := strings.ToUpper(matches[2])

	// Validate field exists
	if _, exists := a.fields[fieldName]; !exists {
		return remaining, fmt.Errorf("field %s not found in entity", fieldName)
	}

	method.SortFields = append(method.SortFields, SortField{
		FieldName: fieldName,
		Direction: direction,
	})

	remaining = remaining[len(matches[0]):]

	// Check for additional sort fields
	if remaining != "" && !strings.HasPrefix(remaining, "And") && !strings.HasPrefix(remaining, "Or") {
		// Try to parse another OrderBy field
		return a.parseOrderBy(remaining, method)
	}

	return remaining, nil
}

// generateParameters generates method parameters based on conditions
func (a *Analyzer) generateParameters(method *QueryMethod) []Parameter {
	var params []Parameter
	paramIndex := 1

	for _, field := range method.Fields {
		fieldType := a.fields[field.FieldName].Type
		typeStr := fieldType.String()

		switch field.Operator {
		case OpBetween:
			params = append(params, Parameter{
				Name: fmt.Sprintf("min%s", field.FieldName),
				Type: typeStr,
			})
			params = append(params, Parameter{
				Name: fmt.Sprintf("max%s", field.FieldName),
				Type: typeStr,
			})
		case OpIn, OpNotIn:
			// For In/NotIn, parameter is a slice
			params = append(params, Parameter{
				Name: fmt.Sprintf("%ss", strings.ToLower(field.FieldName)),
				Type: "[]" + typeStr,
			})
		case OpIsNull, OpIsNotNull, OpTrue, OpFalse:
			// No parameters for these operators
		default:
			params = append(params, Parameter{
				Name: strings.ToLower(field.FieldName),
				Type: typeStr,
			})
		}
		paramIndex++
	}

	return params
}

// ToSQL generates SQL WHERE clause from the method
func (m *QueryMethod) ToSQL(tableName string, fieldToColumn func(string) string) string {
	var conditions []string
	paramIndex := 1

	for i, field := range m.Fields {
		columnName := fieldToColumn(field.FieldName)
		var condition string

		switch field.Operator {
		case OpEqual:
			condition = fmt.Sprintf("%s = $%d", columnName, paramIndex)
			paramIndex++
		case OpNotEqual:
			condition = fmt.Sprintf("%s != $%d", columnName, paramIndex)
			paramIndex++
		case OpGreaterThan:
			condition = fmt.Sprintf("%s > $%d", columnName, paramIndex)
			paramIndex++
		case OpGreaterThanEqual:
			condition = fmt.Sprintf("%s >= $%d", columnName, paramIndex)
			paramIndex++
		case OpLessThan:
			condition = fmt.Sprintf("%s < $%d", columnName, paramIndex)
			paramIndex++
		case OpLessThanEqual:
			condition = fmt.Sprintf("%s <= $%d", columnName, paramIndex)
			paramIndex++
		case OpLike:
			condition = fmt.Sprintf("%s LIKE $%d", columnName, paramIndex)
			paramIndex++
		case OpNotLike:
			condition = fmt.Sprintf("%s NOT LIKE $%d", columnName, paramIndex)
			paramIndex++
		case OpIn:
			// For IN, we need to handle slice parameter
			condition = fmt.Sprintf("%s = ANY($%d)", columnName, paramIndex)
			paramIndex++
		case OpNotIn:
			condition = fmt.Sprintf("%s != ALL($%d)", columnName, paramIndex)
			paramIndex++
		case OpIsNull:
			condition = fmt.Sprintf("%s IS NULL", columnName)
		case OpIsNotNull:
			condition = fmt.Sprintf("%s IS NOT NULL", columnName)
		case OpBetween:
			condition = fmt.Sprintf("%s BETWEEN $%d AND $%d", columnName, paramIndex, paramIndex+1)
			paramIndex += 2
		case OpContaining:
			condition = fmt.Sprintf("%s LIKE $%d", columnName, paramIndex)
			paramIndex++
		case OpStartingWith:
			condition = fmt.Sprintf("%s LIKE $%d", columnName, paramIndex)
			paramIndex++
		case OpEndingWith:
			condition = fmt.Sprintf("%s LIKE $%d", columnName, paramIndex)
			paramIndex++
		case OpIgnoreCase:
			condition = fmt.Sprintf("LOWER(%s) = LOWER($%d)", columnName, paramIndex)
			paramIndex++
		case OpTrue:
			condition = fmt.Sprintf("%s = true", columnName)
		case OpFalse:
			condition = fmt.Sprintf("%s = false", columnName)
		}

		if i > 0 && field.AndOr != "" {
			conditions = append(conditions, field.AndOr)
		}
		conditions = append(conditions, condition)
	}

	whereClause := strings.Join(conditions, " ")

	// Build full query
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	// Add ORDER BY
	if len(m.SortFields) > 0 {
		orderClauses := make([]string, len(m.SortFields))
		for i, sf := range m.SortFields {
			orderClauses[i] = fmt.Sprintf("%s %s", fieldToColumn(sf.FieldName), sf.Direction)
		}
		query += " ORDER BY " + strings.Join(orderClauses, ", ")
	}

	// Add LIMIT
	if m.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", m.Limit)
	}

	return query
}

