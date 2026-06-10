package pgx_crud

import "time"

func (suite *Suite) TestDeleteBatch() {
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

		err = suite.crud.DeleteBatch(suite.ctx, []int64{createdBatch[0].ID, createdBatch[1].ID})

		suite.Require().NoError(err)
		_, err = suite.crud.Get(suite.ctx, createdBatch[0].ID)
		suite.Error(err)
		_, err = suite.crud.Get(suite.ctx, createdBatch[1].ID)
		suite.Error(err)
		_, err = suite.crud.Get(suite.ctx, createdBatch[2].ID)
		suite.NoError(err)
	})
}
