package service

import (
	"log"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	voucherClient protoc.VoucherServiceClient
	conn          *grpc.ClientConn
}

// NewGRPCClient creates a new gRPC client connection
func NewGRPCClient(serverAddress string) (*GRPCClient, error) {
	// Connect to gRPC server
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Create client
	voucherClient := protoc.NewVoucherServiceClient(conn)

	log.Printf("Connected to gRPC server at %s", serverAddress)

	return &GRPCClient{
		voucherClient: voucherClient,
		conn:          conn,
	}, nil
}

// GetVoucherClient returns the voucher service client
func (c *GRPCClient) GetVoucherClient() protoc.VoucherServiceClient {
	return c.voucherClient
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

