package dbrepo

import (
	"database/sql"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/repository"
)

type dbPostgresRepository struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates new postgres db repository
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DBRepo {
	return &dbPostgresRepository{
		App: a,
		DB:  conn,
	}
}

type testdbPostgresRepository struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewTestingRepo
func NewTestingRepo(a *config.AppConfig) repository.DBRepo {
	return &testdbPostgresRepository{
		App: a,
	}
}
