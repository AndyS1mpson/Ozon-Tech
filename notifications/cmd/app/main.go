package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	api "route256/notifications/internal/api/notifications"
	"route256/notifications/internal/clients/telegram"
	"route256/notifications/internal/config"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/kafka"
	"route256/notifications/internal/pkg/cache/lru"
	"route256/notifications/internal/pkg/logger"
	"route256/notifications/internal/pkg/metrics"
	"route256/notifications/internal/pkg/tracer"
	"route256/notifications/internal/repository/postgres"
	"route256/notifications/pkg/notifications_v1"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort      = 50053
	httpPort      = 8082
	topic         = "orders"
	groupID       = "notifications"
	cacheCapacity = 100
)

func main() {
	keepRunning := true
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	tgClient, err := telegram.New(cfg.Telegram.APIKey, cfg.Telegram.ChatID)
	if err != nil {
		log.Fatalf("Cannot connect to telegram: %v", err)
	}

	// Connect to database
	pool, err := pgxpool.Connect(context.Background(), cfg.Postgres.ConnectionString)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	service := domain.NewService(
		postgres.NewMessageRepository(pool),
		tgClient,
		lru.NewLRUCache[domain.CacheKey, domain.CacheVal](cacheCapacity),
	)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.MiddlewareGRPC,
			tracer.MiddlewareGRPC,
			metrics.MiddlewareGRPC,
		),
	)

	reflection.Register(s)
	notifications_v1.RegisterNotificationsServer(s, api.New(service))

	// Start and listen server
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// run kafka consumer
	go func() {
		consumer := kafka.NewConsumerGroupHandler(service)
		client, err := initConsumerGroup(cfg.Brokers, groupID)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		consumptionIsPaused := false
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if err := client.Consume(ctx, []string{topic}, &consumer); err != nil {
					log.Fatalf("Error from consumer: %v", err)
				}

				if ctx.Err() != nil {
					return
				}
			}
		}()

		<-consumer.Ready()
		log.Println("Consumer up and running")

		sigusr1 := make(chan os.Signal, 1)
		signal.Notify(sigusr1, syscall.SIGUSR1)

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

		for keepRunning {
			select {
			case <-ctx.Done():
				log.Println("terminating: context cancelled")
				keepRunning = false
			case <-sigterm:
				log.Println("terminating: via signal")
				keepRunning = false
			case <-sigusr1:
				toggleConsumptionFlow(client, &consumptionIsPaused)
			}
		}

		cancel()
		wg.Wait()

		if err = client.Close(); err != nil {
			log.Fatalf("error closing client: %v", err)
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
		log.Fatalf("Failed to dial server: %v", err)
	}

	mux := runtime.NewServeMux()
	err = notifications_v1.RegisterNotificationsHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	if err := mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		log.Fatalf("something wrong with metrics handler: %v", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: mux,
	}

	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	err = gwServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initConsumerGroup(brokers []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
