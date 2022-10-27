package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/meneses-pt/go_als.zone/entities"
	"net/http"
	"time"
)

// GetMatches handles retrieving matches from the database and renders them to JSON.
func (c *Controller) GetMatches(w http.ResponseWriter, r *http.Request) {
	limit, offset := handleLimitOffsetParams(r)
	var matches []entities.Match
	sqlQuery := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM matches_match m
				 INNER JOIN matches_videogoal mv ON m.id = mv.match_id
				 INNER JOIN matches_team mht ON mht.id = m.home_team_id
				 INNER JOIN matches_team mat ON mat.id = m.away_team_id
		ORDER BY m.datetime DESC
		LIMIT $2 OFFSET $3;
	`, matchFields)
	rows, err := c.DBPool.Query(r.Context(), sqlQuery, c.AppConfig.MediaRoot, limit, offset)
	if err != nil {
		c.Logger.Println(err)
		return
	}
	for rows.Next() {
		m := scanMatch(rows, c)
		matches = append(matches, m)
	}
	encodeResultIntoJson(w, err, matches, c)
}

// GetMatch handles retrieving a single match from the database and renders them to JSON.
func (c *Controller) GetMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	var m entities.Match
	matchSqlQuery := fmt.Sprintf(`
		SELECT %s
		FROM matches_match m
				 INNER JOIN matches_team mht ON mht.id = m.home_team_id
				 INNER JOIN matches_team mat ON mat.id = m.away_team_id
		WHERE m.slug = $2
		LIMIT 1;
	`, matchFields)
	rows, err := c.DBPool.Query(r.Context(), matchSqlQuery, c.AppConfig.MediaRoot, slug)
	if err != nil {
		c.Logger.Println(err)
	}
	for rows.Next() {
		m = scanMatch(rows, c)
	}

	videosSqlQuery := fmt.Sprintf(`
		SELECT %s
		FROM matches_videogoal v
				 INNER JOIN matches_postmatch mp ON v.id = mp.videogoal_id
		WHERE v.match_id = $2
		ORDER BY v.minute;
	`, videoFields)
	mirrorsSqlQuery := fmt.Sprintf(`
		SELECT %s
		FROM matches_videogoalmirror vm
		WHERE vm.videogoal_id = $1;
	`, mirrorFields)
	vRows, err := c.DBPool.Query(r.Context(), videosSqlQuery, c.AppConfig.RedditRoot, m.ID)
	if err != nil {
		c.Logger.Println(err)
	}
	for vRows.Next() {
		v, url := scanVideo(vRows, c)
		var firstMirror = entities.Mirror{
			Title: "Original Link",
			Url:   url,
		}
		v.Mirrors = append(v.Mirrors, firstMirror)
		mRows, err := c.DBPool.Query(r.Context(), mirrorsSqlQuery, v.ID)
		if err != nil {
			c.Logger.Println(err)
		}
		for mRows.Next() {
			m := scanMirror(mRows, c)
			v.Mirrors = append(v.Mirrors, m)
		}
		m.Videos = append(m.Videos, v)
	}

	if err != nil {
		c.Logger.Println(err)
		return
	}
	encodeResultIntoJson(w, err, m, c)
}

// GetMatchesSearchWeek handles retrieving last week matches from the database and renders them to JSON.
func (c *Controller) GetMatchesSearchWeek(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	startTime := time.Now().AddDate(0, 0, -7)
	var matches []entities.Match
	sqlQuery := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM matches_match m
				 INNER JOIN matches_videogoal mv ON m.id = mv.match_id
				 INNER JOIN matches_team mht ON m.home_team_id = mht.id
				 INNER JOIN matches_team mat ON m.away_team_id = mat.id
		WHERE m.datetime >= $2
		  AND (
		      	UPPER(UNACCENT(mht.name)::text) LIKE
				'%%' || 
				UPPER(REPLACE(REPLACE(REPLACE((UNACCENT($3)), E'\\', E'\\\\'), E'%%', E'\\%%'), E'_', E'\\_')) ||
				'%%'
				OR UPPER(UNACCENT(mat.name)::text) LIKE
				   '%%' || 
				   UPPER(REPLACE(REPLACE(REPLACE((UNACCENT($3)), E'\\', E'\\\\'), E'%%', E'\\%%'), E'_', E'\\_')) || 
				   '%%'
			)
		ORDER BY m.datetime DESC
		LIMIT 50;
	`, matchFields)
	rows, err := c.DBPool.Query(r.Context(), sqlQuery, c.AppConfig.MediaRoot, startTime, filter)
	if err != nil {
		c.Logger.Println(err)
		return
	}
	for rows.Next() {
		m := scanMatch(rows, c)
		matches = append(matches, m)
	}
	encodeResultIntoJson(w, err, matches, c)
}
