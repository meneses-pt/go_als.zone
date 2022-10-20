package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/meneses-pt/go_als.zone/controllers"
	"github.com/meneses-pt/go_als.zone/database"
	"github.com/meneses-pt/go_als.zone/util"
	"log"
	"net/http"
)

func RegisterProductRoutes(router *mux.Router) {
	router.HandleFunc("/api/matches", controllers.GetMatches).Methods("GET")
}

func main() {
	// Load Configurations from config.json using Viper
	util.LoadAppConfig()
	// Initialize Database
	database.Connect()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)
	// Register Routes
	RegisterProductRoutes(router)
	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", "3000"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", 3000), router))
}
