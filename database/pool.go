package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/util"
	"log"
)

// Connect establishes a connection to the database from the given configuration file.
func Connect(ctx context.Context, logger *log.Logger, c *util.Config) (*pgxpool.Pool, error) {
	databaseUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
	logger.Printf("Going to connect to database: %s", c.DBName)
	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		logger.Printf("Unable to create connection pool: %v", err)
		return nil, err
	}
	return pool, nil
}
