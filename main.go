package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meneses-pt/go_als.zone/controllers"
	"github.com/meneses-pt/go_als.zone/database"
	"github.com/meneses-pt/go_als.zone/util"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

func registerProductRoutes(router *mux.Router, logger *log.Logger, pool *pgxpool.Pool, appConfig *util.Config) {
	matchController := &controllers.MatchController{DBPool: pool, Logger: logger, AppConfig: appConfig}
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

	//Memory Cache
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(30*time.Second),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router.Use(cacheClient.Middleware)

	// Register Routes
	registerProductRoutes(router, logger, pool, appConfig)
	// Start the server
	log.Printf("Starting Server on %s", appConfig.HTTPAddr)
	log.Fatal(http.ListenAndServe(appConfig.HTTPAddr, router))
}
