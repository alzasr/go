package pgx_crud

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillSelectFields(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		type testStruct struct {
			ID   int64
			Name string
		}

		crud, err := New[testStruct](nil, "")
		if err != nil {
			t.Error(err)
		}
		assert.ElementsMatch(t, []string{`"id"`, `"name"`}, crud.selectFields)
	})

	t.Run("db_tag", func(t *testing.T) {
		type testStruct struct {
			ID     int64
			MyName string `db:"other_name"`
		}

		crud, err := New[testStruct](nil, "")
		if err != nil {
			t.Error(err)
		}
		assert.ElementsMatch(t, []string{`"id"`, `"other_name"`}, crud.selectFields)
	})

	t.Run("case_converter", func(t *testing.T) {
		type testStruct struct {
			ID     int64
			MyName string
		}

		crud, err := New[testStruct](nil, "")
		if err != nil {
			t.Error(err)
		}
		assert.ElementsMatch(t, []string{`"id"`, `"my_name"`}, crud.selectFields)
	})

	t.Run("parent_fields", func(t *testing.T) {
		type parent struct {
			ID int64
		}

		type child struct {
			parent
			Name string
		}

		crud, err := New[child](nil, "")
		if err != nil {
			t.Error(err)
		}
		assert.ElementsMatch(t, []string{`"id"`, `"name"`}, crud.selectFields)
	})
}

func TestCustomError(t *testing.T) {
	t.Run("set_custom_not_found_error", func(t *testing.T) {
		type TestStruct struct {
			ID int64
		}
		crud, err := New[TestStruct](nil, "", WithNorFoundError(errNotFoundExample))
		if err != nil {
			t.Error(err)
		}
		err = crud.getNotFoundErrorOrErr(nil)
		if !errors.Is(err, errNotFoundExample) {
			t.Error("cant set not found error, not valid error in struct")
		}
	})
	t.Run("default_not_found_error", func(t *testing.T) {
		type TestStruct struct {
			ID int64
		}
		crud, err := New[TestStruct](nil, "")
		if err != nil {
			t.Error(err)
		}
		err = crud.getNotFoundErrorOrErr(errNotFoundExample)
		errors.Is(err, errNotFoundExample)
		if !errors.Is(err, errNotFoundExample) {
			t.Error("cant set not found error, not valid error in struct")
		}
	})
}
