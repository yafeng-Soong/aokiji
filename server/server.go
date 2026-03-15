package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/yafeng-Soong/aokiji/registry"

	"google.golang.org/grpc"
)

var listenAddr string

func init() {
	if localIP := getLocalAddr(); localIP == "" {
		log.Fatal("failed to get local IP address")
	} else {
		listenAddr = fmt.Sprintf("%s:0", localIP)
	}
}

// Server represents a gRPC server instance.
type Server struct {
	serviceName string
	grpcServer  *grpc.Server
	registry    registry.Registry
}

// NewServer creates a new Server instance.
func NewServer(opts ...Option) *Server {
	options := &serverOption{
		registryBuilder: emptyRegistryBuilder,
	}

	for _, apply := range opts {
		apply(options)
	}

	serverRegistry, err := options.registryBuilder.build(
		registry.RegistryConfig{
			ServiceName: options.serviceName,
		},
	)
	if err != nil {
		log.Fatalf("failed to create registry: %v", err)
	}

	return &Server{
		serviceName: options.serviceName,
		grpcServer:  options.grpcServer,
		registry:    serverRegistry,
	}
}

// Start starts the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if s.grpcServer == nil {
		log.Fatal("grpc server is not initialized")
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen, error: %v", err)
	}

	addr := lis.Addr().String()
	s.registry.Register(ctx, addr)
	defer s.registry.Deregister(ctx)

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-stop:
			log.Println("shutting down gRPC server...")
		case <-ctx.Done():
			log.Println("server cancelled by context, shutting down...")
		}

		s.grpcServer.GracefulStop()
		log.Println("server stopped")
	}()

	log.Printf("gRPC echo server listening on %s", addr)
	return s.grpcServer.Serve(lis)
}

func getLocalAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				return ip.String()
			}
		}
	}

	return ""
}
