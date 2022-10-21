package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/meneses-pt/go_als.zone/database"
	"github.com/meneses-pt/go_als.zone/entities"
	"log"
	"net/http"
)

func GetMatches(w http.ResponseWriter, r *http.Request) {
	var matches []entities.Match
	rows, err := database.DBPool.Query(context.Background(),
		"SELECT id, slug, datetime, score from matches_match order by datetime desc LIMIT 50")
	for rows.Next() {
		var m entities.Match
		err := rows.Scan(&m.ID, &m.Slug, &m.Datetime, &m.Score)
		if err != nil {
			log.Fatal(err)
		}
		matches = append(matches, m)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(matches)
}
