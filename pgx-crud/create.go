package pgx_crud

import (
	"context"
	"fmt"
)

// Create сохранение элемента в базе
func (c *Crud[T]) Create(ctx context.Context, model *T) (*T, error) {

	query := c.builder.
		Insert(c.table).
		SetMap(c.createValues(model)).
		Suffix(c.returning())

	res, err := c.queryOne(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("queryOne: %w", err)
	}
	return res, nil
}

func (c *Crud[T]) createValues(model *T) map[string]any {
	createValues, _ := c.valuesMap(model, c.insertFields)
	return createValues
}
