package dbrepo

import (
	"database/sql"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/repository"
)

type dbPostgresRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates new postgres db repository
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DBRepo {
	return &dbPostgresRepo{
		App: a,
		DB:  conn,
	}
}

type testdbPostgresRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewTestingRepo
func NewTestingRepo(a *config.AppConfig) repository.DBRepo {
	return &testdbPostgresRepo{
		App: a,
	}
}
