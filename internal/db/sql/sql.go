package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/db/dbtypes"
	"reflect"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	mapColumns sync.Map = sync.Map{}
)

// SQLConnection implements DatabaseConnection for MySQL and PostgreSQL
type SQLConnection struct {
	conn   *sql.DB
	config *config.DatabaseConfig
}

func NewDBConnection(config *config.DatabaseConfig) connection.DatabaseConnection {
	return &SQLConnection{
		config: config,
	}
}

func (p *SQLConnection) GetConnection() any {
	return p.conn
}

func (p *SQLConnection) GetDriver() string {
	return p.config.Driver
}

func (p *SQLConnection) GetName() string {
	switch p.config.Driver {
	case "mysql":
		return "MySQL"
	default:
		return "PostgreSQL"
	}
}

func (p *SQLConnection) Connect(ctx context.Context) error {
	if p.conn != nil {
		return nil
	}

	// Build connection string
	var connectionString string
	if p.config.Username != "" && p.config.Password != "" {
		if p.config.Driver == "mysql" {
			connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
				p.config.Username,
				p.config.Password,
				p.config.Host,
				p.config.Port,
				p.config.Database)
		} else {
			connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				p.config.Host,
				p.config.Port,
				p.config.Username,
				p.config.Password,
				p.config.Database,
			)
		}
	} else {
		if p.config.Driver == "mysql" {
			connectionString = fmt.Sprintf("tcp(%s:%d)/%s",
				p.config.Host,
				p.config.Port,
				p.config.Database)
		} else {
			connectionString = fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable",
				p.config.Host,
				p.config.Port,
				p.config.Database,
			)
		}
	}

	if p.config.Attributes != "" {
		connectionString += " " + p.config.Attributes
	}

	slog.Debug("Attempting to connect to "+p.GetName()+" with", "connectionString", connectionString)

	// Connect to Database using the driver from config
	db, err := sql.Open(p.config.Driver, connectionString)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping"+p.GetName(), "error", err)
		return err
	}

	p.conn = db

	slog.Info("Successfully connected to " + p.GetName())
	return nil
}

func (p *SQLConnection) Close(ctx context.Context) error {
	if p.conn != nil {
		err := p.conn.Close()
		if err == nil {
			slog.Info("Successfully disconnected from " + p.GetName())
		}
		return err
	}
	return nil
}

func (p *SQLConnection) Ping(ctx context.Context) error {
	if p.conn != nil {
		return p.conn.Ping()
	}
	return nil
}

func (m *SQLConnection) Find(ctx context.Context, table string, column []string, filter dbtypes.M, sort map[string]int, limit int64, skip int64) (any, error) {
	colSelect := m.buildSelectClause(column)
	whereClause, args := m.buildWhereClause(filter)
	orderClause := m.buildOrderClause(sort, limit, skip)

	query := fmt.Sprintf("SELECT %s FROM %s%s%s", colSelect, table, whereClause, orderClause)
	slog.Debug("Query: " + query)

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %v", table, err)
	}
	defer rows.Close()

	var results []dbtypes.M
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		entry := make(dbtypes.M)
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	str, _ := json.Marshal(results)
	slog.Debug("Hasil: " + string(str))

	return results, nil
}

func (m *SQLConnection) FindOne(ctx context.Context, result any, table string, column []string, filter dbtypes.M, sort map[string]int) error {
	colSelect := m.buildSelectClause(column)
	whereClause, args := m.buildWhereClause(filter)
	orderClause := m.buildOrderClause(sort, 1, 0)

	query := fmt.Sprintf("SELECT %s FROM %s%s%s", colSelect, table, whereClause, orderClause)

	row := m.conn.QueryRowContext(ctx, query, args...)

	var columns []string
	keyMap := table + "Columns"
	cols, _ := mapColumns.Load(keyMap)
	if cols != nil {
		// Load column names from cache
		columns = cols.([]string)
	} else {
		// Get the column names by querying the table schema
		if m.config.Driver == "mysql" {
			schemaQuery := fmt.Sprintf("SHOW COLUMNS FROM %s", table)
			rows, err := m.conn.QueryContext(ctx, schemaQuery)
			if err != nil {
				return fmt.Errorf("failed to get columns: %v", err)
			}
			defer rows.Close()

			for rows.Next() {
				var column string
				if err := rows.Scan(&column); err != nil {
					return fmt.Errorf("failed to scan column: %v", err)
				}
				columns = append(columns, column)
			}
		} else {
			// For PostgreSQL
			schemaQuery := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", table)
			rows, err := m.conn.QueryContext(ctx, schemaQuery)
			if err != nil {
				return fmt.Errorf("failed to get columns: %v", err)
			}
			defer rows.Close()

			for rows.Next() {
				var column string
				if err := rows.Scan(&column); err != nil {
					return fmt.Errorf("failed to scan column: %v", err)
				}
				columns = append(columns, column)
			}
		}
		mapColumns.Store(keyMap, columns)
	}

	if len(columns) == 0 {
		return fmt.Errorf("no columns found for table %s", table)
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	if err := row.Scan(valuePtrs...); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no rows found")
		}
		return fmt.Errorf("failed to scan row: %v", err)
	}

	// Convert to map first
	resultMap := make(dbtypes.M)
	for i, col := range columns {
		val := values[i]
		b, ok := val.([]byte)
		if ok {
			resultMap[col] = string(b)
		} else {
			resultMap[col] = val
		}
	}

	// Convert to the desired result type
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr || resultValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("result must be a pointer to a struct")
	}

	resultElem := resultValue.Elem()
	resultType := resultElem.Type()

	for i := 0; i < resultElem.NumField(); i++ {
		field := resultElem.Field(i)
		fieldType := resultType.Field(i)

		// Get the JSON tag as the column name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = fieldType.Name
		}

		// Remove comma if it's part of the JSON tag
		if len(jsonTag) > 0 && jsonTag[len(jsonTag)-1] == ',' {
			jsonTag = jsonTag[:len(jsonTag)-1]
		}

		if val, exists := resultMap[jsonTag]; exists && field.CanSet() {
			valValue := reflect.ValueOf(val)
			if valValue.Type().ConvertibleTo(field.Type()) {
				field.Set(valValue.Convert(field.Type()))
			}
		}
	}

	return nil
}

func (m *SQLConnection) InsertOne(ctx context.Context, table string, data any) (any, error) {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	if dataValue.Kind() != reflect.Map && dataValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data must be a map or struct")
	}

	var columns []string
	var placeholders []string
	var values []any

	if dataValue.Kind() == reflect.Map {
		for _, key := range dataValue.MapKeys() {
			columns = append(columns, key.String())
			val := dataValue.MapIndex(key).Interface()
			values = append(values, val)
		}
	} else {
		// Handle struct
		dataType := dataValue.Type()
		for i := 0; i < dataValue.NumField(); i++ {
			field := dataValue.Field(i)
			if field.CanInterface() {
				fieldType := dataType.Field(i)
				// Get the JSON tag as the column name
				jsonTag := fieldType.Tag.Get("json")
				bsonTag := fieldType.Tag.Get("bson")
				if bsonTag != "" {
					jsonTag = bsonTag
				} else if jsonTag == "" {
					jsonTag = fieldType.Name
				}

				// Remove comma if it's part of the JSON tag
				if len(jsonTag) > 0 && jsonTag[len(jsonTag)-1] == ',' {
					jsonTag = jsonTag[:len(jsonTag)-1]
				}

				columns = append(columns, jsonTag)
				values = append(values, field.Interface())
			}
		}
	}

	for range columns {
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	result, err := m.conn.ExecContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into table %s: %v", table, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return dbtypes.M{"id": id}, nil
}

func (m *SQLConnection) UpdateOne(ctx context.Context, table string, filter dbtypes.M, data any) (int64, error) {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	if dataValue.Kind() != reflect.Map && dataValue.Kind() != reflect.Struct {
		return 0, fmt.Errorf("data must be a map or struct")
	}

	var setClauses []string
	var values []any

	if dataValue.Kind() == reflect.Map {
		for _, key := range dataValue.MapKeys() {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", key.String()))
			val := dataValue.MapIndex(key).Interface()
			values = append(values, val)
		}
	} else {
		// Handle struct
		dataType := dataValue.Type()
		for i := 0; i < dataValue.NumField(); i++ {
			field := dataValue.Field(i)
			if field.CanInterface() {
				fieldType := dataType.Field(i)
				// Get the JSON tag as the column name
				jsonTag := fieldType.Tag.Get("json")
				bsonTag := fieldType.Tag.Get("bson")
				if bsonTag != "" {
					jsonTag = bsonTag
				} else if jsonTag == "" {
					jsonTag = fieldType.Name
				}

				// Remove comma if it's part of the JSON tag
				if len(jsonTag) > 0 && jsonTag[len(jsonTag)-1] == ',' {
					jsonTag = jsonTag[:len(jsonTag)-1]
				}

				setClauses = append(setClauses, fmt.Sprintf("%s = ?", jsonTag))
				values = append(values, field.Interface())
			}
		}
	}

	whereClause, whereArgs := m.buildWhereClause(filter)
	values = append(values, whereArgs...)

	query := fmt.Sprintf("UPDATE %s SET %s%s",
		table,
		strings.Join(setClauses, ", "),
		whereClause)

	result, err := m.conn.ExecContext(ctx, query, values...)
	if err != nil {
		return 0, fmt.Errorf("failed to update table %s: %v", table, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %v", err)
	}

	return rowsAffected, nil
}

func (m *SQLConnection) DeleteOne(ctx context.Context, table string, filter dbtypes.M) (any, error) {
	whereClause, args := m.buildWhereClause(filter)

	query := fmt.Sprintf("DELETE FROM %s%s", table, whereClause)

	result, err := m.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to delete from table %s: %v", table, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	return dbtypes.M{"deleted_count": rowsAffected}, nil
}

// buildWhereClause converts MongoDB-style filter to SQL WHERE clause
func (m *SQLConnection) buildWhereClause(filter dbtypes.M) (string, []any) {
	if filter == nil {
		return "", nil
	}

	var whereClauses []string
	var args []any

	for key, value := range filter {
		switch v := value.(type) {
		case dbtypes.M:
			// Handle operators like $gt, $lt, $in, etc.
			for op, val := range v {
				switch op {
				case "$gt":
					whereClauses = append(whereClauses, fmt.Sprintf("%s > ?", key))
					args = append(args, val)
				case "$gte":
					whereClauses = append(whereClauses, fmt.Sprintf("%s >= ?", key))
					args = append(args, val)
				case "$lt":
					whereClauses = append(whereClauses, fmt.Sprintf("%s < ?", key))
					args = append(args, val)
				case "$lte":
					whereClauses = append(whereClauses, fmt.Sprintf("%s <= ?", key))
					args = append(args, val)
				case "$ne":
					whereClauses = append(whereClauses, fmt.Sprintf("%s != ?", key))
					args = append(args, val)
				case "$in":
					if vals, ok := val.([]any); ok {
						placeholders := make([]string, len(vals))
						for i := range vals {
							placeholders[i] = "?"
						}
						whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))
						args = append(args, vals...)
					}
				case "$nin":
					if vals, ok := val.([]any); ok {
						placeholders := make([]string, len(vals))
						for i := range vals {
							placeholders[i] = "?"
						}
						whereClauses = append(whereClauses, fmt.Sprintf("%s NOT IN (%s)", key, strings.Join(placeholders, ", ")))
						args = append(args, vals...)
					}
				case "$or":
					if orConditions, ok := val.([]any); ok {
						var orClauses []string
						for _, condition := range orConditions {
							if condMap, ok := condition.(dbtypes.M); ok {
								orWhere, orArgs := m.buildWhereClause(condMap)
								if orWhere != "" {
									orClauses = append(orClauses, fmt.Sprintf("(%s)", orWhere))
									args = append(args, orArgs...)
								}
							}
						}
						if len(orClauses) > 0 {
							whereClauses = append(whereClauses, fmt.Sprintf("(%s)", strings.Join(orClauses, " OR ")))
						}
					}
				case "$and":
					if andConditions, ok := val.([]any); ok {
						var andClauses []string
						for _, condition := range andConditions {
							if condMap, ok := condition.(dbtypes.M); ok {
								andWhere, andArgs := m.buildWhereClause(condMap)
								if andWhere != "" {
									andClauses = append(andClauses, fmt.Sprintf("(%s)", andWhere))
									args = append(args, andArgs...)
								}
							}
						}
						if len(andClauses) > 0 {
							whereClauses = append(whereClauses, fmt.Sprintf("(%s)", strings.Join(andClauses, " AND ")))
						}
					}
				}
			}
		default:
			// Simple equality
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", key))
			args = append(args, value)
		}
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return whereClause, args
}

func (m *SQLConnection) buildSelectClause(column []string) string {
	if column == nil {
		return "*"
	}

	colSelect := strings.Join(column, ", ")
	if colSelect == "" {
		colSelect = "*"
	}
	return colSelect
}

func (m *SQLConnection) buildOrderClause(sort map[string]int, limit int64, skip int64) string {
	if sort == nil {
		return ""
	}

	var orderClauses []string
	for key, value := range sort {
		if value == 1 {
			orderClauses = append(orderClauses, fmt.Sprintf("%s ASC", key))
		} else {
			orderClauses = append(orderClauses, fmt.Sprintf("%s DESC", key))
		}
	}
	orderClause := ""
	if len(orderClauses) > 0 {
		orderClause = " ORDER BY " + strings.Join(orderClauses, ", ")
	}
	if limit > 0 {
		orderClause += fmt.Sprintf(" LIMIT %d", limit)
	}
	if skip > 0 {
		orderClause += fmt.Sprintf(" OFFSET %d", skip)
	}
	return orderClause
}
