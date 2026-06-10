package pgx_crud

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// Update изменить элемент
func (c *Crud[T]) Update(ctx context.Context, model *T) (*T, error) {

	values, id := c.updateValues(model)
	query := c.builder.
		Update(c.table).
		SetMap(values).
		Where(squirrel.Eq{idFieldName: id}).
		Suffix(c.returning())
	res, err := c.queryOne(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("queryOne: %w", err)
	}
	return res, nil
}

func (c *Crud[T]) updateValues(model *T) (map[string]any, int64) {
	return c.valuesMap(model, c.updateFields)
}
