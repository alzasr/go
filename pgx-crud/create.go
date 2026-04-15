package pg_crud

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Create сохранение элемента в базе
func (c *Crud[T]) Create(ctx context.Context, model *T) (*T, error) {
	query, params, err := c.queriBuilder.
		Insert(c.table).
		SetMap(c.createValues(model)).
		Suffix(c.returning()).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	rows, err := c.db.Query(ctx, query, params...)
	if err != nil {

	}
	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, fmt.Errorf("getx: %w", err)
	}
	return &res, nil
}

func (c *Crud[T]) createValues(model *T) map[string]any {
	createValues, _ := c.valuesMap(model, c.insertFields)
	return createValues
}
