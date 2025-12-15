package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/config"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/handler"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/repository"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = ":50051"
)

func main() {
	log.Println("Starting Voucher Payment Service...")

	// Step 1: Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Step 2: Connect to database
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// Step 3: Initialize repositories
	userRepo := repository.NewUserRepository(db)
	voucherRepo := repository.NewVoucherRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	log.Println("Repositories initialized")

	// Step 4: Initialize services
	mockUPI := service.NewMockUPI(0.95) // 95% success rate
	userService := service.NewUserService(userRepo)
	voucherService := service.NewVoucherService(voucherRepo)
	walletService := service.NewWalletService(walletRepo)
	transactionService := service.NewTransactionService(transactionRepo, userService)
	paymentService := service.NewPaymentService(
		db,
		userService,
		voucherService,
		walletService,
		voucherRepo,
		transactionRepo,
		mockUPI,
	)
	log.Println("Services initialized")

	// Step 5: Initialize handlers
	voucherServiceHandler := handler.NewVoucherServiceHandler(
		voucherService,
		paymentService,
		walletService,
		transactionService,
	)
	log.Println("Handlers initialized")

	// Step 6: Setup gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor), // Logging interceptor
	)

	// Step 7: Register gRPC service
	protoc.RegisterVoucherServiceServer(grpcServer, voucherServiceHandler)

	// Enable gRPC reflection for testing
	reflection.Register(grpcServer)
	log.Println("gRPC server configured")

	// Step 8: Start gRPC server
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcPort, err)
	}

	log.Printf("gRPC server listening on %s", grpcPort)

	// Step 9: Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("Server stopped gracefully")
	case <-ctx.Done():
		log.Println("Shutdown timeout, forcing stop")
		grpcServer.Stop()
	}
}

// unaryInterceptor is a logging interceptor for gRPC
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Log incoming request
	log.Printf("gRPC method: %s, request: %+v", info.FullMethod, req)

	// Call the handler
	resp, err := handler(ctx, req)

	// Log response and duration
	duration := time.Since(start)
	if err != nil {
		log.Printf("gRPC method: %s, error: %v, duration: %v", info.FullMethod, err, duration)
	} else {
		log.Printf("gRPC method: %s, success, duration: %v", info.FullMethod, duration)
	}

	return resp, err
}

