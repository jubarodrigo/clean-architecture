package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/rodrigocavalhero/clean_arch_orders_list/configs"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/infra/graph"
	orderv1 "github.com/rodrigocavalhero/clean_arch_orders_list/internal/infra/grpc/pb/order/v1"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/infra/grpc/service"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/infra/web/webserver"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxDBAttempts = 60
	retryInterval = 1 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("ordersystem: %v", err)
	}
}

func run() error {
	cfg, err := configs.Load(".")
	if err != nil {
		return err
	}
	migPath := cfg.MigrationsDir
	absMig, err := filepath.Abs(migPath)
	if err != nil {
		return fmt.Errorf("migrações: %w", err)
	}
	migURL := "file://" + absMig
	if err := runMigrations(cfg.DatabaseURL, migURL, maxDBAttempts, retryInterval); err != nil {
		return err
	}

	db, err := openSQLWithRetry(cfg.DatabaseURL, maxDBAttempts, retryInterval)
	if err != nil {
		return err
	}
	defer db.Close()

	listUC := NewListOrdersUseCase(db)
	createUC := NewCreateOrderUseCase(db)
	webH := NewWebOrderHandler(db)
	grpcSvc := service.NewOrderService(listUC)

	ws := webserver.NewWebServer(cfg.WebServerPort)
	ws.Router.Get("/order", webH.List)
	ws.Router.Post("/order", webH.Create)

	grpcLis, err := net.Listen("tcp", cfg.GRPCAddress())
	if err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}
	grpcS := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(grpcS, grpcSvc)
	reflection.Register(grpcS)

	srv := graphql_handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			ListOrdersUseCase:  listUC,
			CreateOrderUseCase: createUC,
		}}),
	)
	gqlMux := http.NewServeMux()
	gqlMux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	gqlMux.Handle("/query", srv)
	graphqlServer := &http.Server{Addr: cfg.GraphQLAddress(), Handler: gqlMux}

	errCh := make(chan error, 3)
	go func() {
		log.Printf("REST (Chi) em %s — GET/POST /order\n", cfg.WebServerPort)
		if err := ws.Start(); err != nil {
			errCh <- err
		}
	}()
	go func() {
		log.Printf("gRPC em %s — order.v1.OrderService/ListOrders (reflection ativada)\n", cfg.GRPCAddress())
		errCh <- grpcS.Serve(grpcLis)
	}()
	go func() {
		log.Printf("GraphQL (gqlgen) em %s — playground em /, API em /query\n", cfg.GraphQLAddress())
		errCh <- graphqlServer.ListenAndServe()
	}()
	return <-errCh
}

func runMigrations(dsn, migFileURL string, maxAttempts int, every time.Duration) error {
	if err := waitPostgres(dsn, maxAttempts, every); err != nil {
		return err
	}
	m, err := migrate.New(migFileURL, dsn)
	if err != nil {
		return fmt.Errorf("migrações (new): %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrações (up): %w", err)
	}
	return nil
}

func waitPostgres(dsn string, maxAttempts int, every time.Duration) error {
	for i := 0; i < maxAttempts; i++ {
		pctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		db, nerr := sql.Open("pgx", dsn)
		if nerr == nil {
			nerr = db.PingContext(pctx)
		}
		cancel()
		if nerr == nil {
			_ = db.Close()
			if i > 0 {
				log.Println("banco de dados respondeu, seguindo com migrações")
			}
			return nil
		}
		if db != nil {
			_ = db.Close()
		}
		if i == 0 {
			log.Println("aguardando Postgres (sincronização com o container)...")
		} else if (i+1)%10 == 0 {
			log.Printf("ainda aguardando banco: tentativa %d/%d\n", i+1, maxAttempts)
		}
		time.Sleep(every)
	}
	return fmt.Errorf("timeout: banco inacessível após %d tentativas", maxAttempts)
}

func openSQLWithRetry(dsn string, maxAttempts int, every time.Duration) (*sql.DB, error) {
	for i := 0; i < maxAttempts; i++ {
		db, err := sql.Open("pgx", dsn)
		if err == nil {
			if err = db.PingContext(context.Background()); err == nil {
				return db, nil
			}
			_ = db.Close()
		}
		if (i+1)%10 == 0 {
			log.Printf("sql: retentando conexão %d/%d\n", i+1, maxAttempts)
		}
		time.Sleep(every)
	}
	return nil, fmt.Errorf("timeout: *sql.DB indisponível após %d tentativas", maxAttempts)
}
