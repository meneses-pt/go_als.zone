package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/meneses-pt/go_als.zone/controllers"
	"github.com/meneses-pt/go_als.zone/database"
	"log"
	"net/http"
)

func registerProductRoutes(router *mux.Router) {
	router.HandleFunc("/api/matches", controllers.GetMatches).Methods("GET")
}

func main() {
	// Initialize Database
	database.Connect()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)
	// Register Routes
	registerProductRoutes(router)
	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", "3000"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", 3000), router))
}
