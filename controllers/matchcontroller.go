package controllers

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/entities"
	"github.com/meneses-pt/go_als.zone/util"
	"log"
	"net/http"
)

// MatchController is a controller to get matches from the database.
type MatchController struct {
	DBPool    *pgxpool.Pool
	Logger    *log.Logger
	AppConfig *util.Config
}

// GetMatches handles retrieving matches from the database and renders them to JSON.
func (c *MatchController) GetMatches(w http.ResponseWriter, r *http.Request) {
	limit, offset := HandleLimitOffsetParams(r)
	var matches []entities.Match
	sqlQuery := `SELECT
    			DISTINCT
    			mm.id,
    			mm.score,
    			split_part(mm.score, ':', 1) "home_team_score",
    			split_part(mm.score, ':', 2) "home_team_score",
    			mm.datetime,
    			mm.slug,
    			mht.name,
    			CONCAT($1::text, mht.logo_file),
    			mht.slug,
    			mat.name,
    			CONCAT($1::text, mat.logo_file),
    			mat.slug
				FROM matches_match mm
				INNER JOIN matches_videogoal mv ON mm.id = mv.match_id
				INNER JOIN matches_team mht ON mht.id = mm.home_team_id
				INNER JOIN matches_team mat ON mat.id = mm.away_team_id
				ORDER BY mm.datetime DESC
				LIMIT $2
				OFFSET $3;`
	rows, err := c.DBPool.Query(r.Context(), sqlQuery, c.AppConfig.MediaRoot, limit, offset)
	for rows.Next() {
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
