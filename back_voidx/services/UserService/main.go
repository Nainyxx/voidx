package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"example.com/m/config"
	"example.com/m/handler"
	"example.com/m/repository"
	userservice "example.com/m/userService"
	userv1 "example.com/voidx/proto/UserService"
)

func main() {
	cfg := config.Load()

	db, err := repository.ConnectDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepo(db)
	svc := userservice.NewUserService(repo)
	grpcHandler := handler.NewGrpcHandler(svc)

	server := grpc.NewServer()
	userv1.RegisterUserServiceServer(server, grpcHandler)
	if cfg.Environment == "dev" {
		reflection.Register(server)
	}

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("UserService listening on :%s", cfg.GRPCPort)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("grpc server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down UserService")
	server.GracefulStop()
}
