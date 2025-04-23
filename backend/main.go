package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DebdipWritesCode/VisitorManagementSystem/api"
	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/util"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	// Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Connect to the database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	// Create the store and server
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization",
		},
	})

	// Use CORS handler with the server router
	handler := corsHandler.Handler(server.GetRouter())

	// Start the HTTP server
	err = http.ListenAndServe(config.ServerAddress, handler)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
