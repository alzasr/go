package pg_crud

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
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
)

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
	pgDsn := "postgresql://postgres:postgres@localhost:5432/pgxcrud?sslmode=disable"
	suite.ctx = context.Background()

	dbConn, err := pgx.Connect(suite.ctx, pgDsn)
	if err != nil {
		suite.FailNow(err.Error())
	}

	suite.crud, err = New[Model]("model", dbConn)
	if err != nil {
		suite.FailNow(err.Error())
	}
	_, err = dbConn.Exec(suite.ctx, "DROP TABLE model")
	if err != nil {
		suite.FailNow(err.Error())
	}
	_, err = dbConn.Exec(suite.ctx, schema)
	if err != nil {
		suite.FailNow(err.Error())
	}
}

func TestCrud(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) TearDownSuite() {
}
