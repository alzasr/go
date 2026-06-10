package pgx_crud

import (
	"github.com/Masterminds/squirrel"
	"time"
)

func (suite *Suite) TestDeleteByFilter() {
	suite.Run("smoke", func() {
		createdBatch, err := suite.crud.CreateBatch(suite.ctx, []*Model{
			{
				BaseModel{
					0,
					time.Date(2024, 6, 16, 1, 0, 0, 0, time.UTC),
					Ptr(time.Date(2024, 6, 17, 1, 0, 0, 0, time.UTC)),
				},
				Ptr(int64(1)),
				"smoke-1",
				Ptr("smoke desc 1"),
				"test-from-1",
			},
			{
				BaseModel{
					0,
					time.Date(2024, 6, 16, 2, 0, 0, 0, time.UTC),
					Ptr(time.Date(2024, 6, 17, 2, 0, 0, 0, time.UTC)),
				},
				Ptr(int64(2)),
				"smoke-2",
				Ptr("smoke desc 2"),
				"test-from-2",
			},
			{
				BaseModel{
					0,
					time.Date(2024, 6, 16, 3, 0, 0, 0, time.UTC),
					Ptr(time.Date(2024, 6, 17, 3, 0, 0, 0, time.UTC)),
				},
				Ptr(int64(3)),
				"smoke-3",
				Ptr("smoke desc 3"),
				"test-from-3",
			},
		})
		suite.Require().NoError(err)

		filter := func(builder squirrel.DeleteBuilder) squirrel.DeleteBuilder {
			return builder.Where(squirrel.NotEq{idFieldName: createdBatch[0].ID})
		}

		res, err := suite.crud.DeleteByFilter(suite.ctx, DeleteFilterFn(filter))

		suite.Require().NoError(err)
		suite.GreaterOrEqual(res, int64(2)) // могут быть строки от других тестов, поэтому не менее 2
		_, err = suite.crud.Get(suite.ctx, createdBatch[0].ID)
		suite.NoError(err)
		_, err = suite.crud.Get(suite.ctx, createdBatch[1].ID)
		suite.Error(err)
		_, err = suite.crud.Get(suite.ctx, createdBatch[2].ID)
		suite.Error(err)
	})
}

type DeleteFilterFn func(builder squirrel.DeleteBuilder) squirrel.DeleteBuilder

func (f DeleteFilterFn) ModifyQuery(builder squirrel.DeleteBuilder) squirrel.DeleteBuilder {
	return f(builder)
}
