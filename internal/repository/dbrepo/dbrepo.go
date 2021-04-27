package dbrepo

import (
	"database/sql"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/config"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(app *config.AppConfig, conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: app,
		DB:  conn,
	}
}

func NewTestingRepo(app *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: app,
	}
}
