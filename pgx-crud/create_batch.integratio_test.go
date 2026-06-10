package pgx_crud

import (
	"time"
)

func (suite *Suite) TestCreateBatch() {
	suite.Run("smoke", func() {
		batch := []*Model{
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
		}

		res, err := suite.crud.CreateBatch(suite.ctx, batch)

		suite.Require().NoError(err)
		suite.Equal(len(batch), len(res))
		for i := range res {
			suite.NotEmpty(res[i].ID)
			batch[i].ID = res[i].ID
		}
		suite.Equal(batch, res)
	})

}
