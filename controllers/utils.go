package controllers

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/meneses-pt/go_als.zone/entities"
	"net/http"
)

func handleLimitOffsetParams(r *http.Request) (string, string) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "50"
	}
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	return limit, offset
}

func encodeResultIntoJson(w http.ResponseWriter, err error, o any, c *Controller) {
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(o)
	if err != nil {
		c.Logger.Println(err)
	}
}

func scanMatch(rows pgx.Rows, c *Controller) entities.Match {
	var m entities.Match
	err := rows.Scan(
		&m.ID,
		&m.Score,
		&m.HomeTeamScore,
		&m.AwayTeamScore,
		&m.Datetime,
		&m.Slug,
		&m.HomeTeam.Name,
		&m.HomeTeam.LogoFile,
		&m.HomeTeam.Slug,
		&m.AwayTeam.Name,
		&m.AwayTeam.LogoFile,
		&m.AwayTeam.Slug,
	)
	if err != nil {
		c.Logger.Println(err)
	}
	return m
}

func scanVideo(vRows pgx.Rows, c *Controller) (entities.Video, string) {
	var v entities.Video
	var url string
	err := vRows.Scan(
		&v.ID,
		&v.Title,
		&v.RedditLink,
		&url,
	)
	if err != nil {
		c.Logger.Println(err)
	}
	return v, url
}

func scanMirror(mRows pgx.Rows, c *Controller) entities.Mirror {
	var m entities.Mirror
	err := mRows.Scan(
		&m.Title,
		&m.Url,
	)
	if err != nil {
		c.Logger.Println(err)
	}
	return m
}

func scanTeam(rows pgx.Rows, c *Controller) entities.Team {
	var t entities.Team
	err := rows.Scan(
		&t.ID,
		&t.Name,
		&t.LogoFile,
		&t.Slug,
		nil,
	)
	if err != nil {
		c.Logger.Println(err)
	}
	return t
}

var matchFields = `
m.id,
m.score,
split_part(m.score, ':', 1) "home_team_score",
split_part(m.score, ':', 2) "away_team_score",
m.datetime,
m.slug,
mht.name,
CONCAT($1::text, mht.logo_file),
mht.slug,
mat.name,
CONCAT($1::text, mat.logo_file),
mat.slug
`

var videoFields = `v.id, v.title, CONCAT($1::text, mp.permalink), RIGHT(LEFT(mp.permalink, 25), 6), v.url`

var teamFields = `
t.id, 
t.name,
CONCAT($1::text, t.logo_file),
t.slug
`

var mirrorFields = `vm.title, vm.url`
