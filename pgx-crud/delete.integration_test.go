package pgx_crud

import "time"

func (suite *Suite) TestDelete() {
	suite.Run("smoke", func() {
		created, err := suite.crud.Create(suite.ctx, &Model{
			BaseModel{
				0,
				time.Date(2024, 6, 16, 1, 0, 0, 0, time.UTC),
				Ptr(time.Date(2024, 6, 17, 1, 0, 0, 0, time.UTC)),
			},
			Ptr(int64(1)),
			"smoke-1",
			Ptr("smoke desc 1"),
			"test-from-1",
		})
		suite.Require().NoError(err)

		affectedRows, err := suite.crud.Delete(suite.ctx, created.ID)

		suite.EqualValues(1, affectedRows)

		suite.Require().NoError(err)
		_, err = suite.crud.Get(suite.ctx, created.ID)
		suite.Error(err)
	})
}
