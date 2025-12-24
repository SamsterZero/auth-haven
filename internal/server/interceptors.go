package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	// Tenant ID
	if tenantIDs := md.Get("tenant-id"); len(tenantIDs) > 0 {
		ctx = context.WithValue(ctx, "tenant-id", tenantIDs[0])
	}

	// JWT Authentication
	if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
		token := authHeaders[0]
		userID, err := validateJWT(token)
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
		}
		ctx = context.WithValue(ctx, "user-id", userID)
	}

	log.Printf("handling %s", info.FullMethod)
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	log.Printf("method=%s duration=%s error=%v", info.FullMethod, duration, err)

	// TODO: Push audit log asynchronously
	return resp, err
}
