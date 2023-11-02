package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	cfg      Config
	services []interface{}
	log      *zap.Logger
	exitFunc func(code int)
}

type GRPCRegisterer interface {
	GRPCRegisterer() func(s *grpc.Server)
}

func New(cfg Config) *Server {
	s := &Server{
		cfg: cfg,
	}

	s.log = stdoutLogger()

	s.AddExitFunc(func(_ int) {
		if err := s.log.Sync(); err != nil {
			println(fmt.Sprintf("failed to sync logger: %v", err))
		}
	})

	return s
}

func (s *Server) Serve(ctx context.Context, authInterceptor grpc.UnaryServerInterceptor) {
	var cancelFunc func()
	ctx, cancelFunc = context.WithCancel(ctx)
	go s.watchShutdown(cancelFunc)

	opts := []grpc.ServerOption{
		s.unaryMiddlewares(authInterceptor),
	}

	srv := grpc.NewServer(opts...)

	for _, se := range s.services {
		if r, ok := se.(GRPCRegisterer); ok {
			r.GRPCRegisterer()(srv)
		}
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(s.cfg.Port))
	if err != nil {
		panic(fmt.Sprintf("binding to %d: %v", s.cfg.Port, err))
	}

	s.log.Info("serving started")

	go func() {
		if err := srv.Serve(l); err != nil {
			s.log.Sugar().Infof("server exited: %v", err)
		}
	}()

	<-ctx.Done()

	srv.GracefulStop()
	s.log.Sugar().Info("service gracefully stopped")
	s.exit(0)
}

func (s *Server) AddExitFunc(fn func(code int)) {
	exitFunc := s.exitFunc
	if exitFunc == nil {
		exitFunc = os.Exit
	}

	s.exitFunc = func(code int) {
		fn(code)
		exitFunc(code)
	}
}

func (s *Server) exit(code int) {
	if s.exitFunc == nil {
		os.Exit(code)
	} else {
		s.exitFunc(code)
	}
}

func (s *Server) watchShutdown(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	s.log.Sugar().Infof("received %s signal from OS\n", sig.String())
	cancelFunc()
}

func (s *Server) Register(dependencies ...interface{}) {
	for _, dep := range dependencies {
		switch d := dep.(type) {
		case GRPCRegisterer:
			s.services = append(s.services, d)
		}
	}
}

func (s *Server) unaryMiddlewares(authInterceptor grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		grpc_validator.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		authInterceptor,
	)
}

func (s *Server) Log() *zap.Logger {
	return s.log
}

func stdoutLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	log, _ := config.Build()
	return log
}
