package main

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/controllers"
	"github.com/meneses-pt/go_als.zone/util"
	"log"
)

func registerMatchesRoutes(router *mux.Router, logger *log.Logger, pool *pgxpool.Pool, appConfig *util.Config) {
	matchController := &controllers.Controller{DBPool: pool, Logger: logger, AppConfig: appConfig}
	router.HandleFunc("/api/matches/", matchController.GetMatches).Methods("GET")
	router.HandleFunc("/api/matches/{slug}", matchController.GetMatch).Methods("GET")
	router.HandleFunc("/api/matches-search-week/", matchController.GetMatchesSearchWeek).Methods("GET")
}

func registerTeamsRoutes(router *mux.Router, logger *log.Logger, pool *pgxpool.Pool, appConfig *util.Config) {
	teamController := &controllers.Controller{DBPool: pool, Logger: logger, AppConfig: appConfig}
	router.HandleFunc("/api/teams/", teamController.GetTeams).Methods("GET")
	router.HandleFunc("/api/teams/{slug}", teamController.GetTeam).Methods("GET")
}
