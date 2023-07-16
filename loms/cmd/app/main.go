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
	"route256/loms/internal/kafka"
	"route256/loms/internal/pkg/logger"
	"route256/loms/internal/pkg/metrics"
	"route256/loms/internal/pkg/tracer"
	"route256/loms/internal/repository/postgres"
	"route256/loms/internal/sender"
	"route256/loms/pkg/loms_v1"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort    = 50052
	httpPort    = 8081
	kafkaTopic  = "orders"
	serviceName = "loms"
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
		return errors.Wrap(err, "read config")
	}

	// Init tracer
	if err := tracer.InitGlobal(serviceName, cfg.Jaeger.Host, cfg.Jaeger.Port); err != nil {
		return err
	}

	// Connect to database
	pool, err := pgxpool.Connect(context.Background(), cfg.Postgres.ConnectionString)
	if err != nil {
		return errors.Wrap(err, "connect to postgres")
	}

	// Create TCP connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	// Create new kafka producer
	producer, err := kafka.NewProducer(cfg.Brokers)
	if err != nil {
		return errors.Wrap(err, "fail connect to kafka broker")
	}

	service := domain.New(
		postgres.NewOrderRepository(pool),
		postgres.NewStockRepository(pool),
		sender.NewKafkaSender(producer, kafkaTopic),
	)

	// Create new gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.MiddlewareGRPC,
			tracer.MiddlewareGRPC,
			metrics.MiddlewareGRPC,
		),
	)
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
		return errors.Wrap(err, "Failed to dial server")
	}

	mux := runtime.NewServeMux()
	err = loms_v1.RegisterLomsHandler(context.Background(), mux, conn)
	if err != nil {
		return errors.Wrap(err, "Failed to register gateway")
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
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}
