package pgx_crud

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

const (
	PGDSNENV = "PG_DSN"
)

var (
	schema = `
CREATE TABLE IF NOT EXISTS "model" (
    "id" BIGSERIAL NOT NULL,
    "parent_id" BIGINT NULL DEFAULT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NULL DEFAULT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP NULL DEFAULT NULL,
    "from" TEXT                               
)
`
	errNotFoundExample = errors.New("not found error example")
)

func TestCrud(t *testing.T) {
	suite.Run(t, new(Suite))
}

type BaseModel struct {
	ID        int64 `dbSkip:"insert,update"`
	CreatedAt time.Time
	DeletedAt *time.Time
}

type Model struct {
	BaseModel
	ParentID    *int64
	Name        string
	Description *string
	From        string
}

type Suite struct {
	suite.Suite

	ctx  context.Context
	crud *Crud[Model]
}

func (suite *Suite) SetupSuite() {
	//pgDsn := os.Getenv(PGDSNENV)
	//if pgDsn == "" {
	//	suite.FailNow(PGDSNENV + " ENVIRONMENT VARIABLE IS REQUIRED")
	//}
	pgDsn := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	suite.ctx = context.Background()

	db, err := pgxpool.New(suite.ctx, pgDsn)
	if err != nil {
		suite.FailNow(err.Error())
	}

	suite.crud, err = New[Model](db, "model")
	if err != nil {
		suite.FailNow(err.Error())
	}

	_, err = db.Exec(suite.ctx, "DROP TABLE IF EXISTS model")
	if err != nil {
		suite.FailNow(err.Error())
	}
	_, err = db.Exec(suite.ctx, schema)
	if err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *Suite) TearDownSuite() {
	//_, err := suite.crud.db.GetRunner(suite.ctx).Exec(suite.ctx, `DROP TABLE "model"`)
	//if err != nil {
	//	suite.FailNow("cant teardown env: ", err.Error())
	//}
}
