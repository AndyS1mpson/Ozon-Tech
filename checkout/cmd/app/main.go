// Project launch
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	api "route256/checkout/internal/api/cart"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/products"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/postgres"
	"route256/checkout/pkg/cart_v1"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
	httpPort = 8080
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
	pool, err := pgxpool.Connect(context.Background(), cfg.Postgres.ConnectionString)
	if err != nil {
		log.Fatalln("connect to postgres: ", err)
	}

	// Create TCP connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create new gRPC server
	s := grpc.NewServer()
	reflection.Register(s)
	cart_v1.RegisterCartServer(s, api.New(domain.New(
		loms.New(cfg.Services.Loms),
		products.New(cfg.Services.Products, cfg.Token),
		postgres.New(pool),
	)))

	// Start and listen gRPC server
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
	err = cart_v1.RegisterCartHandler(context.Background(), mux, conn)
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
		return fmt.Errorf("Failed to serve: %w", err)
	}

	return nil
}
