package pg_crud

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const idFieldName = "id"

// Filter - контракт для фильтра
type Filter interface {
	ModifyQuery(squirrel.SelectBuilder) squirrel.SelectBuilder
}

// New конструктор
func New[T any](
	table string,
	db *pgx.Conn,
	options ...Options,
) (*Crud[T], error) {
	crud := &Crud[T]{
		table,
		db,
		squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		nil,
		nil,
		nil,
		&crudOptions{},
	}
	err := crud.init()
	if err != nil {
		return nil, fmt.Errorf("init: %w", err)
	}
	for _, opt := range options {
		opt(crud.options)
	}
	return crud, nil
}

// Crud реализация базовых операций
type Crud[T any] struct {
	table        string
	db           *pgx.Conn
	queriBuilder squirrel.StatementBuilderType
	selectFields []string
	insertFields []string
	updateFields []string
	options      *crudOptions
}

type crudOptions struct {
	NotFoundError error
	TracingPrefix string
}

func (c *Crud[T]) init() error {
	err := c.fillQueryFields()
	if err != nil {
		return fmt.Errorf("fillQueryFields: %w", err)
	}
	return nil
}

func (c *Crud[T]) fillQueryFields() error {
	var model T
	reflectModel := reflect.ValueOf(model)
	if reflectModel.Kind() != reflect.Struct {
		return fmt.Errorf("model has kind %v, needed Struct", reflectModel.Kind())
	}

	fields := c.getFields(&model)

	idField, ok := fields[idFieldName]
	if !ok {
		return fmt.Errorf("model must contains id field with type int64")
	}
	if idField.Value.Kind() != reflect.Int64 {
		return fmt.Errorf("ID field of model must be type int64 have %v", idField.Value.Kind())
	}

	c.selectFields = make([]string, 0, len(fields))
	c.insertFields = make([]string, 0, len(fields))
	c.updateFields = make([]string, 0, len(fields))

	for _, field := range fields {
		fieldName := `"` + field.Name + `"`
		c.selectFields = append(c.selectFields, fieldName)
		if !field.Skips["insert"] && field.Name != idFieldName {
			c.insertFields = append(c.insertFields, fieldName)
		}
		if !field.Skips["update"] && field.Name != idFieldName {
			c.updateFields = append(c.updateFields, fieldName)
		}
	}

	return nil
}

func (c *Crud[T]) getFields(model *T) map[string]DBField {
	reflectModel := reflect.ValueOf(model)
	reflectModel = reflectModel.Elem()
	if reflectModel.Kind() != reflect.Struct {
		panic(fmt.Errorf("model has kind %v, needed Struct", reflectModel.Kind()))
	}

	result := map[string]DBField{}

	structQueue := []reflect.Value{reflectModel}
	for len(structQueue) > 0 {
		structItem := structQueue[0]
		structQueue = structQueue[1:]

		for i := 0; i < structItem.NumField(); i++ {
			field := structItem.Field(i)
			fieldType := structItem.Type().Field(i)
			if field.Kind() == reflect.Struct && fieldType.Anonymous {
				structQueue = append(structQueue, field)
				continue
			}

			dbField := DBField{
				Value: &field,
			}

			tag, ok := fieldType.Tag.Lookup("db")
			if ok {
				if tag == "-" {
					continue
				}
				dbField.Name = strings.TrimSpace(strings.Split(tag, ",")[0])
			} else {
				dbField.Name = PascalToSnake(fieldType.Name)
			}

			skipTag, _ := fieldType.Tag.Lookup("dbSkip")
			dbField.Skips = SliceToMap(strings.Split(skipTag, ","))
			result[dbField.Name] = dbField
		}
	}
	return result
}

// DBField информация о поле в базе
type DBField struct {
	Name  string
	Skips map[string]bool
	Value *reflect.Value
}

func (c *Crud[T]) returning() string {
	return "returning " + strings.Join(c.selectFields, ", ")
}

func (c *Crud[T]) valuesMap(model *T, fields []string) (map[string]any, int64) {
	values := c.getFields(model)
	res := make(map[string]any, len(fields))
	for _, field := range fields {
		clearField := strings.Trim(field, `"`)
		value := values[clearField]
		if value.Value.Kind() == reflect.Ptr && !value.Value.IsZero() {
			res[field] = value.Value.Elem().Interface()
		} else {
			res[field] = value.Value.Interface()
		}
	}
	idField := values[idFieldName]
	return res, idField.Value.Int()
}

func (c *Crud[T]) valuesSlice(model *T, fields []string) []any {
	values, _ := c.valuesMap(model, fields)
	res := make([]any, 0, len(fields))
	for _, field := range fields {
		res = append(res, values[field])
	}
	return res
}

// WithNorFoundError используется как опция для установки кастомной ошибки
func WithNorFoundError(err error) Options {
	return func(c *crudOptions) {
		c.NotFoundError = err
	}
}

func WithTracingPrefix(prefix string) Options {
	return func(c *crudOptions) {
		c.TracingPrefix = prefix
	}
}

// Options опциональные параметры конструктора
type Options func(c *crudOptions)

// возращает кастомную ошибку если она имплементирована, либо возращает базовую ошибку
func (c *Crud[T]) getNotFoundErrorOrErr(err error) error {
	if c.options.NotFoundError != nil {
		return c.options.NotFoundError
	}
	return err
}
