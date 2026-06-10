package pgx_crud

import (
	"time"

	"github.com/Masterminds/squirrel"
)

var _ Filter = IDFilter{}

type IDFilter struct {
	ids []int64
}

func (I IDFilter) ModifyQuery(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return builder.Where(squirrel.Eq{"id": I.ids})
}

func (suite *Suite) TestList() {
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
				"test-from-1",
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
				"test-from-1",
			},
		})
		suite.Require().NoError(err)

		list, err := suite.crud.List(suite.ctx, IDFilter{ids: []int64{createdBatch[0].ID, createdBatch[1].ID, createdBatch[2].ID}})

		suite.Require().NoError(err)
		suite.ElementsMatch(createdBatch, list)
	})
}
