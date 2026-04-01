package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/yafeng-Soong/aokiji/registry"

	"google.golang.org/grpc"
)

var listenAddr string

func init() {
	localIP := getLocalAddr()
	if localIP == "" {
		log.Fatal("failed to get local IP address")
	}
	listenAddr = fmt.Sprintf("%s:0", localIP)
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
		registry: registry.NewEmptyRegistry(),
	}

	for _, apply := range opts {
		apply(options)
	}

	if options.serviceName == "" {
		log.Fatal("service name is required")
	}

	serviceName := path.Join(registry.Prefix, options.serviceName)
	return &Server{
		serviceName: serviceName,
		grpcServer:  options.grpcServer,
		registry:    options.registry,
	}
}

// Start starts the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if s.grpcServer == nil {
		log.Fatal("grpc server is not initialized")
	}

	lc := net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen, error: %v", err)
	}

	addr := lis.Addr().String()
	if err = s.registry.Register(ctx, s.serviceName, addr); err != nil {
		log.Fatalf("failed to register service: %v", err)
	}
	defer func() {
		if err = s.registry.Deregister(ctx, s.serviceName, addr); err != nil {
			log.Printf("failed to deregister service: %v", err)
		}
	}()

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
