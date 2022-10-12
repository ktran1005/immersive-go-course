package auth

import (
	"fmt"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is meant to be used by other services to talk with the Auth service.
type Client struct {
	conn   *grpc.ClientConn
	cancel context.CancelFunc
	aC     pb.AuthClient
}

// Create a new Client for the auth service.
// Call Close() to release resources associated with this Client.
func NewClient(ctx context.Context, target string) (*Client, error) {
	return newClientWithOpts(ctx, target, defaultOpts()...)
}

// Call Close() to release resources associated with this Client.
func (c *Client) Close() error {
	// We cancel the context in case the connection is still being formed...
	c.cancel()
	// ...but according to grpc.DialContext docs, we still need to call conn.Close()
	return c.conn.Close()
}

func defaultOpts() []grpc.DialOption {
	return []grpc.DialOption{
		// TODO: insecure connection should move to TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

// Use this function in tests to configure the underlying client with options
func newClientWithOpts(ctx context.Context, target string, opts ...grpc.DialOption) (*Client, error) {
	// Wrapping the context WithCancel allows us to cancel the connection if the caller chooses to
	// immediately Close() the Client.
	ctx, cancel := context.WithCancel(ctx)
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &Client{
		conn:   conn,
		cancel: cancel,
		aC:     pb.NewAuthClient(conn),
	}, nil
}