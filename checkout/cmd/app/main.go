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
	"route256/checkout/internal/pkg/logger"
	"route256/checkout/internal/pkg/metrics"
	"route256/checkout/internal/pkg/ratelimit"
	"route256/checkout/internal/pkg/tracer"
	"route256/checkout/internal/repository/postgres"
	"route256/checkout/pkg/cart_v1"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort        = 50051
	httpPort        = 8090
	productRPSLimit = 10
	maxConcurrency  = 5
	serviceName     = "checkout"
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

	// Init tracer
	if err := tracer.InitGlobal(serviceName, cfg.Jaeger.Host, cfg.Jaeger.Port); err != nil {
		return err
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
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.MiddlewareGRPC,
			tracer.MiddlewareGRPC,
			metrics.MiddlewareGRPC,
		),
	)
	reflection.Register(s)

	d := domain.New(
		loms.New(cfg.Services.Loms),
		products.New(cfg.Services.Products, cfg.Token, *ratelimit.New(productRPSLimit, maxConcurrency)),
		postgres.New(pool),
	)
	cart_v1.RegisterCartServer(s, api.New(d))

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

	if err := mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		return fmt.Errorf("something wrong with metrics handler: %w", err)
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
