package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/meneses-pt/go_als.zone/entities"
	"net/http"
)

// GetTeams handles retrieving teams from the database and renders them to JSON.
func (c *Controller) GetTeams(w http.ResponseWriter, r *http.Request) {
	limit, offset := handleLimitOffsetParams(r)
	filter := r.URL.Query().Get("filter")
	var teams []entities.Team
	sqlQuery := fmt.Sprintf(`
		SELECT %s, count(DISTINCT m.id) as matches_count
		FROM matches_team t
				 INNER JOIN matches_match m ON t.id = m.home_team_id OR t.id = m.away_team_id
				 INNER JOIN matches_videogoal vg ON m.id = vg.match_id
		WHERE UPPER(UNACCENT(t.name)::text) LIKE
				'%%' || 
				UPPER(REPLACE(REPLACE(REPLACE((UNACCENT($2)), E'\\', E'\\\\'), E'%%', E'\\%%'), E'_', E'\\_')) || 
				'%%'
		GROUP BY t.id, t.logo_file, t.slug
		ORDER BY matches_count DESC
		LIMIT $3 OFFSET $4;
	`, teamFields)
	rows, err := c.DBPool.Query(r.Context(), sqlQuery, c.AppConfig.MediaRoot, filter, limit, offset)
	if err != nil {
		c.Logger.Println(err)
		return
	}
	for rows.Next() {
		t := scanTeam(rows, c)
		teams = append(teams, t)
	}
	encodeResultIntoJson(w, err, teams, c)
}

// GetTeam handles retrieving a single team from the database and renders them to JSON.
func (c *Controller) GetTeam(w http.ResponseWriter, r *http.Request) {
	limit, offset := handleLimitOffsetParams(r)
	params := mux.Vars(r)
	slug := params["slug"]
	var t entities.Team
	sqlQuery := fmt.Sprintf(`
		SELECT %s, 'nil'
		FROM matches_team t
		WHERE t.slug = $2
		LIMIT 1;
	`, teamFields)
	rows, err := c.DBPool.Query(r.Context(), sqlQuery, c.AppConfig.MediaRoot, slug)
	if err != nil {
		c.Logger.Println(err)
	}
	for rows.Next() {
		t = scanTeam(rows, c)
	}

	matchesSqlQuery := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM matches_match m
				 INNER JOIN matches_videogoal mv ON m.id = mv.match_id
				 INNER JOIN matches_team mht ON mht.id = m.home_team_id
				 INNER JOIN matches_team mat ON mat.id = m.away_team_id
		WHERE m.home_team_id = $2
		   OR m.away_team_id = $2
		ORDER BY m.datetime DESC
		LIMIT $3 OFFSET $4;
	`, matchFields)
	mRows, err := c.DBPool.Query(r.Context(), matchesSqlQuery, c.AppConfig.MediaRoot, t.ID, limit, offset)
	if err != nil {
		c.Logger.Println(err)
	}
	for mRows.Next() {
		m := scanMatch(mRows, c)
		t.Matches = append(t.Matches, m)
	}
	encodeResultIntoJson(w, err, t, c)
}
