package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/lemissel/GRPCService/internal/svc"
	greeting_v1 "github.com/lemissel/GRPCService/pkg/pb/proto/greeting/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	listener net.Listener
	server   *grpc.Server
	logger   *zap.Logger
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	initListener()

	server = grpc.NewServer()

	greeting_v1.RegisterGreeterServiceServer(server, &svc.GreeterService{})
	logger.Info("Handlers registred")

	go signalsListener(server)

	logger.Info("Starting GRPC server")
	if err := server.Serve(listener); err != nil {
		logger.Panic("Failed to start GRPC server", zap.Error(err))
	}
}

func initListener() {

	var err error

	addr := "localhost:50051"

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		logger.Panic("Failed to listen",
			zap.String("address", addr),
			zap.Error(err),
		)
	}

	logger.Info("Started listening... ", zap.String("address", addr))

	return

}

func signalsListener(server *grpc.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	_ = <-sigs

	logger.Info("Gracefull stopping server")
	server.GracefulStop()
}
