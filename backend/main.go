package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/DebdipWritesCode/VisitorManagementSystem/api"
	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/gapi"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/DebdipWritesCode/VisitorManagementSystem/util"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGrpcServer(store, config)
}

func runGrpcServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create gRPC server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterVisiTrackServer(grpcServer, server)
	reflection.Register(grpcServer) // This command registers the server for reflection; Reflection allows clients to discover the services and methods available on the server.

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("Starting GRPC server on %s...\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGinServer(store db.Store, config util.Config) {
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
	err = http.ListenAndServe(config.HTTPServerAddress, handler)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
