package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/util"
	"os"
)

var DBPool *pgxpool.Pool
var _ error

func Connect() {
	var c = util.AppConfig
	var databaseUrl = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
	_, err := fmt.Fprintf(os.Stderr, "%s", databaseUrl)
	if err != nil {
		return
	}
	DBPool, err = pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
