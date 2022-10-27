package controllers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/util"
	"log"
)

// Controller is a controller to get matches from the database.
type Controller struct {
	DBPool    *pgxpool.Pool
	Logger    *log.Logger
	AppConfig *util.Config
}
