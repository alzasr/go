package pgx_crud

import (
	"time"
)

func (suite *Suite) TestCreate() {
	suite.Run("smoke", func() {
		model := &Model{
			BaseModel{
				0,
				time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC),
				Ptr(time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)),
			},
			Ptr(int64(12)),
			"smoke",
			Ptr("smoke desc"),
			"test-from",
		}

		res, err := suite.crud.Create(suite.ctx, model)

		suite.Require().NoError(err)
		suite.NotEmpty(res.ID)
		model.ID = res.ID
		suite.Equal(model, res)
	})

	suite.Run("empty_string_nullable_value", func() {
		model := &Model{
			Description: Ptr(""),
		}

		res, err := suite.crud.Create(suite.ctx, model)

		suite.Require().NoError(err)
		suite.NotEmpty(res.ID)
		model.ID = res.ID
		suite.Equal(model, res)
	})
}
