package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NavaneethWKT/CapStone_GO_Lang/client/handler"
	"github.com/NavaneethWKT/CapStone_GO_Lang/client/service"
	"github.com/gin-gonic/gin"
)

const (
	httpPort     = ":8080"
	grpcServerAddress = "localhost:50051"
)

func main() {
	log.Println("Starting REST API Gateway...")

	// Step 1: Connect to gRPC server
	grpcClient, err := service.NewGRPCClient(grpcServerAddress)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer grpcClient.Close()
	log.Println("Connected to gRPC server")

	// Step 2: Initialize handlers
	voucherHandler := handler.NewVoucherHandler(grpcClient)
	paymentHandler := handler.NewPaymentHandler(grpcClient)
	walletHandler := handler.NewWalletHandler(grpcClient)
	transactionHandler := handler.NewTransactionHandler(grpcClient)
	log.Println("Handlers initialized")

	// Step 3: Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "voucher-payment-api",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Voucher routes
		vouchers := api.Group("/vouchers")
		{
			vouchers.GET("/search", voucherHandler.Search)
			vouchers.POST("/buy", paymentHandler.BuyVoucher)
		}

		// Wallet routes
		wallet := api.Group("/wallet")
		{
			wallet.GET("/balance/:user_id", walletHandler.GetBalance)
		}

		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.GET("/:user_id", transactionHandler.ListTransactions)
		}
	}

	log.Printf("REST API Gateway listening on %s", httpPort)

	// Step 4: Start server in goroutine
	go func() {
		if err := router.Run(httpPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Step 5: Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down REST API Gateway...")
	log.Println("Server stopped")
}

