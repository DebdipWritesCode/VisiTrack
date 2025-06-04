package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"

	"github.com/DebdipWritesCode/VisitorManagementSystem/api"
	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	_ "github.com/DebdipWritesCode/VisitorManagementSystem/doc/statik"
	"github.com/DebdipWritesCode/VisitorManagementSystem/gapi"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/DebdipWritesCode/VisitorManagementSystem/util"
	"github.com/DebdipWritesCode/VisitorManagementSystem/worker"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config:")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Connect to the database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database:")
	}

	// Create the store and server
	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt) // Initialize the task distributor with Redis options.

	go runTaskProcessor(redisOpt, store)                // Run the task processor in a separate goroutine.
	go runGatewayServer(store, config, taskDistributor) // Run the HTTP Gateway server in a separate goroutine.
	runGrpcServer(store, config, taskDistributor)       // Run the gRPC server.
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store) // Create a new task processor with Redis options and the store.
	log.Info().Msg("Starting task processor...")
	err := taskProcessor.Start() // Start the task processor to listen for tasks.
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start task processor:")
	}
}

func runGrpcServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger) // This creates a new gRPC server with the logger interceptor. The logger interceptor is used to log incoming requests and outgoing responses.
	pb.RegisterVisiTrackServer(grpcServer, server)
	reflection.Register(grpcServer) // This command registers the server for reflection; Reflection allows clients to discover the services and methods available on the server.

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener:")
	}

	log.Info().Msgf("Starting GRPC server on %s...\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server:")
	}
}

func runGatewayServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server:")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{ // This option controls how the JSON is marshaled. Marshaled means converting the protobuf message to JSON format.
			UseProtoNames: true,
			// EmitUnpopulated: true, // This option ensures that fields with zero values are included in the JSON output.
		},
		UnmarshalOptions: protojson.UnmarshalOptions{ // This option controls how the JSON is unmarshaled. Unmarshaled means converting the JSON format back to a protobuf message.
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterVisiTrackHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register gRPC handler:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statisFs, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik filesystem:")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statisFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener:")
	}

	log.Info().Msgf("Starting HTTP Gateway server on %s...\n", listener.Addr().String())
	handler := gapi.HttpLogger(mux) // This wraps the HTTP handler with the logger middleware to log incoming requests and outgoing responses.

	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP Gateway server:")
	}
}

func runGinServer(store db.Store, config util.Config) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
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
		log.Fatal().Err(err).Msg("cannot start server:")
	}
}
