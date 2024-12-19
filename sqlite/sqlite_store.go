package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/Silencevoice/go-store"
)

type SQLiteStore[T store.SQLModel] struct {
	db *sql.DB
}

func NewSQLiteStore[T store.SQLModel](db *sql.DB) *SQLiteStore[T] {
	return &SQLiteStore[T]{
		db: db,
	}
}

func (s *SQLiteStore[T]) GetByID(ctx context.Context, id string) (*T, error) {
	tableName := store.TableName[T]()
	query := fmt.Sprintf(store.SelectAllFieldsWhereIdEquals, tableName)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if row == nil {
		return nil, errors.New("no row found")
	}

	var entity T
	err = row.Scan(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (s *SQLiteStore[T]) GetMultipleByID(ctx context.Context, ids []string) ([]*T, error) {
	tableName := store.TableName[T]()

	query := fmt.Sprintf(store.SelectAllFieldsWhereIdIn, tableName)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return buildEntities[T](rows)
}

func (s *SQLiteStore[T]) GetAll(ctx context.Context) ([]*T, error) {
	tableName := store.TableName[T]()
	query := fmt.Sprintf(store.SelectAllFieldsFromTable, tableName)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return buildEntities[T](rows)
}

func buildEntity[T store.SQLModel](rows *sql.Rows) (*T, error) {
	var r T
	value := reflect.ValueOf(&r).Elem()
	numCols := value.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := value.Field(i)
		columns[i] = field.Addr().Interface()
	}

	err := rows.Scan(columns...)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func buildEntities[T store.SQLModel](rows *sql.Rows) ([]*T, error) {
	ret := []*T{}
	for rows.Next() {
		r, err := buildEntity[T](rows)
		if err != nil {
			return nil, err
		}
		ret = append(ret, r)
	}

	return ret, nil
}
