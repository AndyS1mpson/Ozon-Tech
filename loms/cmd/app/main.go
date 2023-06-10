// Project launch
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	api "route256/loms/internal/api/loms"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/postgres"
	"route256/loms/pkg/loms_v1"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50052
	httpPort = 8081
)

// Service start point
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

// Start server and listen gRPC and HTTP requests
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	// Connect to database
	pool, err := pgxpool.Connect(context.Background(), cfg.Postgres.ConnectionString)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}

	// Create TCP connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	service := domain.New(
		postgres.NewOrderRepository(pool),
		postgres.NewStockRepository(pool),
	)

	// Create new gRPC server
	s := grpc.NewServer()
	reflection.Register(s)
	loms_v1.RegisterLomsServer(s, api.New(service))

	// Start and listen server
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Create connection to gRPC-gateway
	conn, err := grpc.DialContext(
		context.Background(),
		lis.Addr().String(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("Failed to dial server: %w", err)
	}

	mux := runtime.NewServeMux()
	err = loms_v1.RegisterLomsHandler(context.Background(), mux, conn)
	if err != nil {
		return fmt.Errorf("Failed to register gateway: %w", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: mux,
	}

	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	err = gwServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
