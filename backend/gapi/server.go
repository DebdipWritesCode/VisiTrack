package gapi

import (
	"fmt"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/DebdipWritesCode/VisitorManagementSystem/util"
)

type Server struct {
	pb.UnimplementedVisiTrackServer
	config util.Config
	store  db.Store
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	fmt.Println("Initializing Twilio client...")
	util.InitTwilio()

	server := &Server{
		config: config,
		store:  store,
	}

	return server, nil
}
