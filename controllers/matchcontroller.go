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

	vmSqlQuery := fmt.Sprintf(`
		SELECT %s, %s
		FROM matches_videogoal v
				INNER JOIN matches_postmatch mp ON v.id = mp.videogoal_id                                 
         		LEFT JOIN matches_videogoalmirror vm on v.id = vm.videogoal_id
		WHERE v.match_id = $2
		ORDER BY v.id, vm.id;
	`, videoFields, mirrorFields)
	vmRows, err := c.DBPool.Query(r.Context(), vmSqlQuery, c.AppConfig.RedditRoot, m.ID)
	if err != nil {
		c.Logger.Println(err)
	}
	for vmRows.Next() {
		var v entities.Video
		var url string
		var mr entities.Mirror
		err := vmRows.Scan(
			&v.ID,
			&v.Title,
			&v.RedditLink,
			&v.SimplePermalink,
			&url,
			&mr.Title,
			&mr.Url,
		)
		if err != nil {
			c.Logger.Println(err)
		}
		if len(m.Videos) == 0 || m.Videos[len(m.Videos)-1].ID != v.ID {
			title := "Original Link"
			var firstMirror = entities.Mirror{
				Title: &title,
				Url:   &url,
			}
			v.Mirrors = append(v.Mirrors, firstMirror)
			if mr.Url != nil {
				v.Mirrors = append(v.Mirrors, mr)
			}
			m.Videos = append(m.Videos, v)
		} else {
			m.Videos[len(m.Videos)-1].Mirrors = append(m.Videos[len(m.Videos)-1].Mirrors, mr)
		}
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
