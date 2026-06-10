package pgx_crud

import "time"

func (suite *Suite) TestGet() {
	suite.Run("smoke", func() {
		createdValue, err := suite.crud.Create(suite.ctx, &Model{
			BaseModel{
				0,
				time.Date(2024, 6, 16, 1, 0, 0, 0, time.UTC),
				Ptr(time.Date(2024, 6, 17, 1, 0, 0, 0, time.UTC)),
			},
			Ptr(int64(1)),
			"smoke-get",
			Ptr("smoke desc get"),
			"test-from",
		})
		suite.Require().NoError(err)

		dbValue, err := suite.crud.Get(suite.ctx, createdValue.ID)

		suite.Require().NoError(err)
		suite.Equal(createdValue, dbValue)

	})
}
