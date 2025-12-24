package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"auth-haven/internal/config"

	"google.golang.org/grpc"
)

func StartGRPC(cfg *config.Config, db *sql.DB) error {
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(Unary),
	}

	s := grpc.NewServer(opts...)

	// TODO: Register services
	// proto.RegisterAuthServiceServer(s, auth.NewService(db, cfg.JWTSecret))
	// proto.RegisterTenantServiceServer(s, tenant.NewService(db))

	log.Printf("gRPC server running on %s", cfg.GRPCPort)
	return s.Serve(lis)
}
