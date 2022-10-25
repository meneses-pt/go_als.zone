package controllers

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/entities"
	"log"
	"net/http"
)

// MatchController is a controller to get matches from the database.
type MatchController struct {
	DBPool *pgxpool.Pool
	Logger *log.Logger
}

// GetMatches handles retrieving matches from the database and renders them to JSON.
func (c *MatchController) GetMatches(w http.ResponseWriter, r *http.Request) {
	var matches []entities.Match
	rows, err := c.DBPool.Query(r.Context(), "SELECT id, slug, datetime, score from matches_match order by datetime desc LIMIT 50")
	for rows.Next() {
		var m entities.Match
		err := rows.Scan(&m.ID, &m.Slug, &m.Datetime, &m.Score)
		if err != nil {
			c.Logger.Println(err)
		}
		matches = append(matches, m)
	}
	if err != nil {
		c.Logger.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(matches)
	if err != nil {
		c.Logger.Println(err)
	}
}
