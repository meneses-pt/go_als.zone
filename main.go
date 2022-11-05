package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/meneses-pt/go_als.zone/database"
	"github.com/meneses-pt/go_als.zone/util"
	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	logger := log.Default()

	appConfig, err := util.LoadAppConfig()
	if err != nil {
		log.Fatalf("Exiting because of error reading configuration: %s", err)
	}

	// Initialize Database
	pool, err := database.Connect(ctx, logger, appConfig)
	if err != nil {
		log.Fatalf("Exiting because of error creating DB connection: %s", err)
	}
	defer pool.Close()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	//Memory Cache code
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		log.Fatalf("Exiting because of Memory Adapter configuration: %s", err)
	}
	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(30*time.Second),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		log.Fatalf("Exiting because of cache Adapter configuration: %s", err)
	}
	router.Use(cacheClient.Middleware)

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	corsObj := handlers.AllowedOrigins([]string{
		"https://goals.zone",
		"https://goals.africa",
		"https://gzreact.meneses.pt",
		"https://gz.meneses.pt",
		"https://videogoals.meneses.pt",
		"http://localhost:3000",
	})

	// Register Routes
	registerMatchesRoutes(router, logger, pool, appConfig)
	registerTeamsRoutes(router, logger, pool, appConfig)

	// Start the server
	log.Printf("Starting Server on %s", appConfig.HTTPAddr)
	log.Fatal(http.ListenAndServe(appConfig.HTTPAddr, handlers.CORS(corsObj)(router)))
}
