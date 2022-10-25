package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/controllers"
	"github.com/meneses-pt/go_als.zone/database"
	"github.com/meneses-pt/go_als.zone/util"
	"log"
	"net/http"
)

func registerProductRoutes(router *mux.Router, logger *log.Logger, pool *pgxpool.Pool) {
	matchController := &controllers.MatchController{DBPool: pool, Logger: logger}
	router.HandleFunc("/api/matches", matchController.GetMatches).Methods("GET")
}

func main() {
	ctx := context.Background()
	logger := log.Default()

	appConfig, err := util.LoadAppConfig()
	if err != nil {
		log.Fatal("Exiting because of error reading configuration")
	}

	// Initialize Database
	pool, err := database.Connect(ctx, logger, appConfig)
	if err != nil {
		log.Fatal("Exiting because of error creating DB connection")
	}
	defer pool.Close()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)
	// Register Routes
	registerProductRoutes(router, logger, pool)
	// Start the server
	log.Printf("Starting Server on %s", appConfig.HTTPAddr)
	log.Fatal(http.ListenAndServe(appConfig.HTTPAddr, router))
}
