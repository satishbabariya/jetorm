package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseRepository provides the base implementation for Repository interface
type BaseRepository[T any, ID comparable] struct {
	db       *Database
	tx       *Tx
	entity   *Entity
	tableName string
	pkField  string
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T any, ID comparable](db *Database) (*BaseRepository[T, ID], error) {
	var zero T
	entity, err := EntityMetadata(zero)
	if err != nil {
		return nil, err
	}

	if entity.PrimaryKey == nil {
		return nil, ErrNoPrimaryKey
	}

	return &BaseRepository[T, ID]{
		db:        db,
		entity:    entity,
		tableName: entity.TableName,
		pkField:   entity.PrimaryKey.DBName,
	}, nil
}

// Save inserts or updates an entity
func (r *BaseRepository[T, ID]) Save(ctx context.Context, entity *T) (*T, error) {
	if r.tx != nil {
		return r.saveWithTx(ctx, entity)
	}
	return r.saveWithPool(ctx, entity)
}

func (r *BaseRepository[T, ID]) saveWithPool(ctx context.Context, entity *T) (*T, error) {
	// Get primary key value
	pkValue := r.getPKValue(entity)
	
	// Check if entity exists (has non-zero primary key)
	if r.isZeroValue(pkValue) {
		// Insert
		return r.insert(ctx, entity, r.db.pool)
	}
	
	// Update
	return r.update(ctx, entity, r.db.pool)
}

func (r *BaseRepository[T, ID]) saveWithTx(ctx context.Context, entity *T) (*T, error) {
	tx := r.tx.tx
	
	// Get primary key value
	pkValue := r.getPKValue(entity)
	
	// Check if entity exists (has non-zero primary key)
	if r.isZeroValue(pkValue) {
		// Insert
		return r.insertTx(ctx, entity, tx)
	}
	
	// Update
	return r.updateTx(ctx, entity, tx)
}

func (r *BaseRepository[T, ID]) insert(ctx context.Context, entity *T, pool *pgxpool.Pool) (*T, error) {
	fields, values, placeholders := r.buildInsertQuery(entity)
	
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		r.tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)
	
	r.logQuery(query, values)
	
	row := pool.QueryRow(ctx, query, values...)
	result := new(T)
	if err := r.scanRow(row, result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (r *BaseRepository[T, ID]) insertTx(ctx context.Context, entity *T, tx pgx.Tx) (*T, error) {
	fields, values, placeholders := r.buildInsertQuery(entity)
	
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		r.tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)
	
	r.logQuery(query, values)
	
	row := tx.QueryRow(ctx, query, values...)
	result := new(T)
	if err := r.scanRow(row, result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (r *BaseRepository[T, ID]) update(ctx context.Context, entity *T, pool *pgxpool.Pool) (*T, error) {
	fields, values := r.buildUpdateQuery(entity)
	pkValue := r.getPKValue(entity)
	values = append(values, pkValue)
	
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = $%d RETURNING *",
		r.tableName,
		strings.Join(fields, ", "),
		r.pkField,
		len(values),
	)
	
	r.logQuery(query, values)
	
	row := pool.QueryRow(ctx, query, values...)
	result := new(T)
	if err := r.scanRow(row, result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (r *BaseRepository[T, ID]) updateTx(ctx context.Context, entity *T, tx pgx.Tx) (*T, error) {
	fields, values := r.buildUpdateQuery(entity)
	pkValue := r.getPKValue(entity)
	values = append(values, pkValue)
	
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = $%d RETURNING *",
		r.tableName,
		strings.Join(fields, ", "),
		r.pkField,
		len(values),
	)
	
	r.logQuery(query, values)
	
	row := tx.QueryRow(ctx, query, values...)
	result := new(T)
	if err := r.scanRow(row, result); err != nil {
		return nil, err
	}
	
	return result, nil
}

// SaveAll saves multiple entities
func (r *BaseRepository[T, ID]) SaveAll(ctx context.Context, entities []*T) ([]*T, error) {
	results := make([]*T, 0, len(entities))
	for _, entity := range entities {
		saved, err := r.Save(ctx, entity)
		if err != nil {
			return nil, err
		}
		results = append(results, saved)
	}
	return results, nil
}

// Update updates an existing entity (must have non-zero primary key)
func (r *BaseRepository[T, ID]) Update(ctx context.Context, entity *T) (*T, error) {
	pkValue := r.getPKValue(entity)
	if r.isZeroValue(pkValue) {
		return nil, ErrInvalidID
	}

	if r.tx != nil {
		tx := r.tx.tx
		return r.updateTx(ctx, entity, tx)
	}
	return r.update(ctx, entity, r.db.pool)
}

// UpdateAll updates multiple entities
func (r *BaseRepository[T, ID]) UpdateAll(ctx context.Context, entities []*T) ([]*T, error) {
	results := make([]*T, 0, len(entities))
	for _, entity := range entities {
		updated, err := r.Update(ctx, entity)
		if err != nil {
			return nil, err
		}
		results = append(results, updated)
	}
	return results, nil
}

// FindByID finds an entity by ID
func (r *BaseRepository[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", r.tableName, r.pkField)
	r.logQuery(query, []interface{}{id})
	
	var row pgx.Row
	if r.tx != nil {
		tx := r.tx.tx
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = r.db.pool.QueryRow(ctx, query, id)
	}
	
	result := new(T)
	if err := r.scanRow(row, result); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return result, nil
}

// FindAll finds all entities
func (r *BaseRepository[T, ID]) FindAll(ctx context.Context) ([]*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	r.logQuery(query, nil)
	
	var rows pgx.Rows
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		rows, err = tx.Query(ctx, query)
	} else {
		rows, err = r.db.pool.Query(ctx, query)
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanRows(rows)
}

// FindAllByIDs finds entities by IDs
func (r *BaseRepository[T, ID]) FindAllByIDs(ctx context.Context, ids []ID) ([]*T, error) {
	if len(ids) == 0 {
		return []*T{}, nil
	}
	
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s IN (%s)",
		r.tableName,
		r.pkField,
		strings.Join(placeholders, ", "),
	)
	r.logQuery(query, args)
	
	var rows pgx.Rows
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		rows, err = tx.Query(ctx, query, args...)
	} else {
		rows, err = r.db.pool.Query(ctx, query, args...)
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanRows(rows)
}

// Delete deletes an entity
func (r *BaseRepository[T, ID]) Delete(ctx context.Context, entity *T) error {
	pkValue := r.getPKValue(entity)
	return r.DeleteByID(ctx, pkValue.(ID))
}

// DeleteByID deletes an entity by ID
func (r *BaseRepository[T, ID]) DeleteByID(ctx context.Context, id ID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", r.tableName, r.pkField)
	r.logQuery(query, []interface{}{id})
	
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		_, err = tx.Exec(ctx, query, id)
	} else {
		_, err = r.db.pool.Exec(ctx, query, id)
	}
	
	return err
}

// DeleteAll deletes multiple entities
func (r *BaseRepository[T, ID]) DeleteAll(ctx context.Context, entities []*T) error {
	for _, entity := range entities {
		if err := r.Delete(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// DeleteAllByIDs deletes multiple entities by their IDs
func (r *BaseRepository[T, ID]) DeleteAllByIDs(ctx context.Context, ids []ID) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s IN (%s)",
		r.tableName,
		r.pkField,
		strings.Join(placeholders, ", "),
	)
	r.logQuery(query, args)

	var err error
	if r.tx != nil {
		tx := r.tx.tx
		_, err = tx.Exec(ctx, query, args...)
	} else {
		_, err = r.db.pool.Exec(ctx, query, args...)
	}

	return err
}

// Count counts all entities
func (r *BaseRepository[T, ID]) Count(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName)
	r.logQuery(query, nil)
	
	var count int64
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		err = tx.QueryRow(ctx, query).Scan(&count)
	} else {
		err = r.db.pool.QueryRow(ctx, query).Scan(&count)
	}
	
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// ExistsById checks if an entity exists by ID
func (r *BaseRepository[T, ID]) ExistsById(ctx context.Context, id ID) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", r.tableName, r.pkField)
	r.logQuery(query, []interface{}{id})
	
	var exists bool
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		err = tx.QueryRow(ctx, query, id).Scan(&exists)
	} else {
		err = r.db.pool.QueryRow(ctx, query, id).Scan(&exists)
	}
	
	if err != nil {
		return false, err
	}
	
	return exists, nil
}

// FindAllPaged finds entities with pagination
func (r *BaseRepository[T, ID]) FindAllPaged(ctx context.Context, pageable Pageable) (*Page[T], error) {
	// Build query with pagination
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	
	// Add sorting
	if len(pageable.Sort.Orders) > 0 {
		orderClauses := make([]string, len(pageable.Sort.Orders))
		for i, order := range pageable.Sort.Orders {
			direction := "ASC"
			if order.Direction == Desc {
				direction = "DESC"
			}
			orderClauses[i] = fmt.Sprintf("%s %s", order.Field, direction)
		}
		query += " ORDER BY " + strings.Join(orderClauses, ", ")
	}
	
	// Add pagination
	if pageable.Size > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageable.Size, pageable.Page*pageable.Size)
	}
	
	r.logQuery(query, nil)
	
	// Execute query
	var rows pgx.Rows
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		rows, err = tx.Query(ctx, query)
	} else {
		rows, err = r.db.pool.Query(ctx, query)
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	content, err := r.scanRows(rows)
	if err != nil {
		return nil, err
	}
	
	// Get total count
	totalElements, err := r.Count(ctx)
	if err != nil {
		return nil, err
	}
	
	// Calculate page info
	totalPages := 0
	if pageable.Size > 0 {
		totalPages = int((totalElements + int64(pageable.Size) - 1) / int64(pageable.Size))
	}
	
	numberOfElements := len(content)
	
	return &Page[T]{
		Content:          content,
		Pageable:         pageable,
		TotalElements:    totalElements,
		TotalPages:       totalPages,
		Size:             pageable.Size,
		Number:           pageable.Page,
		NumberOfElements: numberOfElements,
		First:            pageable.Page == 0,
		Last:             pageable.Page >= totalPages-1 || totalPages == 0,
		Empty:            numberOfElements == 0,
		Sort:             pageable.Sort,
	}, nil
}

// SaveBatch saves entities in batches
func (r *BaseRepository[T, ID]) SaveBatch(ctx context.Context, entities []*T, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100 // Default batch size
	}

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}

		batch := entities[i:end]
		_, err := r.SaveAll(ctx, batch)
		if err != nil {
			return fmt.Errorf("batch save failed at offset %d: %w", i, err)
		}
	}

	return nil
}

// FindOne finds a single entity matching the specification
func (r *BaseRepository[T, ID]) FindOne(ctx context.Context, spec Specification[T]) (*T, error) {
	if spec == nil {
		return nil, ErrNotFound
	}

	whereClause, args := spec.ToSQL()
	if whereClause == "" {
		return nil, ErrNotFound
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", r.tableName, whereClause)
	r.logQuery(query, args)

	var row pgx.Row
	if r.tx != nil {
		row = r.tx.tx.QueryRow(ctx, query, args...)
	} else {
		row = r.db.pool.QueryRow(ctx, query, args...)
	}

	result := new(T)
	if err := r.scanRow(row, result); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return result, nil
}

// FindAllWithSpec finds all entities matching the specification
func (r *BaseRepository[T, ID]) FindAllWithSpec(ctx context.Context, spec Specification[T]) ([]*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	var args []interface{}

	if spec != nil {
		whereClause, specArgs := spec.ToSQL()
		if whereClause != "" {
			query += " WHERE " + whereClause
			args = specArgs
		}
	}

	r.logQuery(query, args)

	var rows pgx.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.tx.Query(ctx, query, args...)
	} else {
		rows, err = r.db.pool.Query(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// FindAllPagedWithSpec finds entities with pagination matching the specification
func (r *BaseRepository[T, ID]) FindAllPagedWithSpec(ctx context.Context, spec Specification[T], pageable Pageable) (*Page[T], error) {
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	var args []interface{}

	// Add WHERE clause if specification provided
	if spec != nil {
		whereClause, specArgs := spec.ToSQL()
		if whereClause != "" {
			query += " WHERE " + whereClause
			args = specArgs
		}
	}

	// Add sorting
	if len(pageable.Sort.Orders) > 0 {
		orderClauses := make([]string, len(pageable.Sort.Orders))
		for i, order := range pageable.Sort.Orders {
			direction := "ASC"
			if order.Direction == Desc {
				direction = "DESC"
			}
			orderClauses[i] = fmt.Sprintf("%s %s", order.Field, direction)
		}
		query += " ORDER BY " + strings.Join(orderClauses, ", ")
	}

	// Add pagination
	if pageable.Size > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageable.Size, pageable.Page*pageable.Size)
	}

	r.logQuery(query, args)

	// Execute query
	var rows pgx.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.tx.Query(ctx, query, args...)
	} else {
		rows, err = r.db.pool.Query(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	content, err := r.scanRows(rows)
	if err != nil {
		return nil, err
	}

	// Get total count with specification
	totalElements, err := r.CountWithSpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	// Calculate page info
	totalPages := 0
	if pageable.Size > 0 {
		totalPages = int((totalElements + int64(pageable.Size) - 1) / int64(pageable.Size))
	}

	numberOfElements := len(content)

	return &Page[T]{
		Content:          content,
		Pageable:         pageable,
		TotalElements:    totalElements,
		TotalPages:       totalPages,
		Size:             pageable.Size,
		Number:           pageable.Page,
		NumberOfElements: numberOfElements,
		First:            pageable.Page == 0,
		Last:             pageable.Page >= totalPages-1 || totalPages == 0,
		Empty:            numberOfElements == 0,
		Sort:             pageable.Sort,
	}, nil
}

// CountWithSpec counts entities matching the specification
func (r *BaseRepository[T, ID]) CountWithSpec(ctx context.Context, spec Specification[T]) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName)
	var args []interface{}

	if spec != nil {
		whereClause, specArgs := spec.ToSQL()
		if whereClause != "" {
			query += " WHERE " + whereClause
			args = specArgs
		}
	}

	r.logQuery(query, args)

	var count int64
	var err error
	if r.tx != nil {
		err = r.tx.tx.QueryRow(ctx, query, args...).Scan(&count)
	} else {
		err = r.db.pool.QueryRow(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

// ExistsWithSpec checks if any entity exists matching the specification
func (r *BaseRepository[T, ID]) ExistsWithSpec(ctx context.Context, spec Specification[T]) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s", r.tableName)
	var args []interface{}

	if spec != nil {
		whereClause, specArgs := spec.ToSQL()
		if whereClause != "" {
			query += " WHERE " + whereClause
			args = specArgs
		}
	}
	query += ")"

	r.logQuery(query, args)

	var exists bool
	var err error
	if r.tx != nil {
		err = r.tx.tx.QueryRow(ctx, query, args...).Scan(&exists)
	} else {
		err = r.db.pool.QueryRow(ctx, query, args...).Scan(&exists)
	}

	if err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteWithSpec deletes entities matching the specification and returns rows affected
func (r *BaseRepository[T, ID]) DeleteWithSpec(ctx context.Context, spec Specification[T]) (int64, error) {
	if spec == nil {
		return 0, fmt.Errorf("specification cannot be nil for delete")
	}

	whereClause, args := spec.ToSQL()
	if whereClause == "" {
		return 0, fmt.Errorf("specification must have a WHERE clause for delete")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", r.tableName, whereClause)
	r.logQuery(query, args)

	var result pgconn.CommandTag
	var err error
	if r.tx != nil {
		result, err = r.tx.tx.Exec(ctx, query, args...)
	} else {
		result, err = r.db.pool.Exec(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

// WithTx returns a repository bound to a transaction
func (r *BaseRepository[T, ID]) WithTx(tx *Tx) Repository[T, ID] {
	return &BaseRepository[T, ID]{
		db:        r.db,
		tx:        tx,
		entity:    r.entity,
		tableName: r.tableName,
		pkField:   r.pkField,
	}
}

// Query executes a raw SQL query and returns results
func (r *BaseRepository[T, ID]) Query(ctx context.Context, query string, args ...interface{}) ([]*T, error) {
	r.logQuery(query, args)

	var rows pgx.Rows
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		rows, err = tx.Query(ctx, query, args...)
	} else {
		rows, err = r.db.pool.Query(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// QueryOne executes a raw SQL query and returns a single result
func (r *BaseRepository[T, ID]) QueryOne(ctx context.Context, query string, args ...interface{}) (*T, error) {
	r.logQuery(query, args)

	var row pgx.Row
	if r.tx != nil {
		tx := r.tx.tx
		row = tx.QueryRow(ctx, query, args...)
	} else {
		row = r.db.pool.QueryRow(ctx, query, args...)
	}

	result := new(T)
	if err := r.scanRow(row, result); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return result, nil
}

// Exec executes a raw SQL statement and returns the number of rows affected
func (r *BaseRepository[T, ID]) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	r.logQuery(query, args)

	var result pgconn.CommandTag
	var err error
	if r.tx != nil {
		tx := r.tx.tx
		result, err = tx.Exec(ctx, query, args...)
	} else {
		result, err = r.db.pool.Exec(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

// Helper methods

func (r *BaseRepository[T, ID]) getPKValue(entity *T) interface{} {
	v := reflect.ValueOf(entity).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := r.entity.Type.Field(i)
		if field.Name == r.entity.PrimaryKey.Name {
			return v.Field(i).Interface()
		}
	}
	return nil
}

func (r *BaseRepository[T, ID]) isZeroValue(v interface{}) bool {
	return reflect.ValueOf(v).IsZero()
}

func (r *BaseRepository[T, ID]) buildInsertQuery(entity *T) ([]string, []interface{}, []string) {
	v := reflect.ValueOf(entity).Elem()
	
	fields := make([]string, 0)
	values := make([]interface{}, 0)
	placeholders := make([]string, 0)
	
	idx := 1
	for i := 0; i < v.NumField(); i++ {
		fieldMeta := r.entity.Fields[i]
		
		// Skip auto-increment primary keys
		if fieldMeta.AutoIncrement && fieldMeta.PrimaryKey {
			continue
		}
		
		// Skip auto-now fields (they should be handled by database)
		if fieldMeta.AutoNowAdd || fieldMeta.AutoNow {
			continue
		}
		
		fields = append(fields, fieldMeta.DBName)
		values = append(values, v.Field(i).Interface())
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
		idx++
	}
	
	return fields, values, placeholders
}

func (r *BaseRepository[T, ID]) buildUpdateQuery(entity *T) ([]string, []interface{}) {
	v := reflect.ValueOf(entity).Elem()
	
	fields := make([]string, 0)
	values := make([]interface{}, 0)
	
	idx := 1
	for i := 0; i < v.NumField(); i++ {
		fieldMeta := r.entity.Fields[i]
		
		// Skip primary key
		if fieldMeta.PrimaryKey {
			continue
		}
		
		// Skip auto-now-add fields
		if fieldMeta.AutoNowAdd {
			continue
		}
		
		fields = append(fields, fmt.Sprintf("%s = $%d", fieldMeta.DBName, idx))
		values = append(values, v.Field(i).Interface())
		idx++
	}
	
	return fields, values
}

func (r *BaseRepository[T, ID]) scanRow(row pgx.Row, dest *T) error {
	v := reflect.ValueOf(dest).Elem()
	
	// Create slice of pointers to struct fields
	fields := make([]interface{}, len(r.entity.Fields))
	for i := range r.entity.Fields {
		fields[i] = v.Field(i).Addr().Interface()
	}
	
	return row.Scan(fields...)
}

func (r *BaseRepository[T, ID]) scanRows(rows pgx.Rows) ([]*T, error) {
	results := make([]*T, 0)
	
	for rows.Next() {
		entity := new(T)
		if err := r.scanRow(rows, entity); err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return results, nil
}

func (r *BaseRepository[T, ID]) logQuery(query string, args []interface{}) {
	if r.db.config.LogSQL {
		r.db.logger.Debug("executing query", "query", query, "args", args)
	}
}

