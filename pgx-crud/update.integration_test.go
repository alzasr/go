package pgx_crud

import "time"

func (suite *Suite) TestUpdate() {
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

		createdValue.ParentID = Ptr(int64(12))
		createdValue.Name += "updated"
		createdValue.Description = Ptr(*createdValue.Description + "updated")
		createdValue.CreatedAt = createdValue.CreatedAt.Add(time.Hour)
		createdValue.DeletedAt = Ptr(createdValue.DeletedAt.Add(time.Hour))

		res, err := suite.crud.Update(suite.ctx, createdValue)

		suite.Require().NoError(err)
		suite.Equal(createdValue, res)
	})

	suite.Run("empty_string", func() {
		createdValue, err := suite.crud.Create(suite.ctx, &Model{
			Description: Ptr("smoke desc get"),
		})
		suite.Require().NoError(err)

		createdValue.Description = Ptr("")

		res, err := suite.crud.Update(suite.ctx, createdValue)

		suite.Require().NoError(err)
		suite.Equal(createdValue, res)
	})
}
