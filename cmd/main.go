package main

import (
	"avito-shop/avito"
	"avito-shop/src/app"
	"avito-shop/src/interceptor"
	"avito-shop/src/repository"
	db2 "avito-shop/src/repository/pg/db"
	"avito-shop/src/service"
	"avito-shop/tools"
	"context"
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	grpcPort := os.Getenv("GRPC_SERVER_PORT")
	httpPort := os.Getenv("HTTP_SERVER_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))

	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка при открытии соединения с базой данных:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	fmt.Println("Успешное подключение к базе данных")

	txCommitter := db2.NewTxProvider(db)

	merchRepo := repository.NewMerchRepository(txCommitter)
	merchSlice, err := merchRepo.GetAllMerch()

	if err != nil {
		log.Fatal(err)
	}

	merchMap := tools.SliceToMap(merchSlice)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.CreateAuthInterceptor()))

	userRepo := repository.NewUserRepo(txCommitter)
	userService := service.NewUserService(userRepo)

	purchaseRepo := repository.NewPurchaseRepo(txCommitter, merchMap)
	purchaseService := service.NewPurchaseService(purchaseRepo)

	transactionRepo := repository.NewTransactionRepo(txCommitter)
	transactionService := service.NewTransactionService(transactionRepo)

	authService := service.NewAuthService(userRepo)

	avito.RegisterAvitoShopServer(grpcServer, app.NewService(*userService, *purchaseService, *transactionService, *authService))
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	mux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = avito.RegisterAvitoShopHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%v", grpcPort), opts)
	if err != nil {
		log.Fatal("Failed to register gRPC Gateway:", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", httpPort),
		Handler: mux,
	}

	fmt.Println("Starting HTTP server on :" + httpPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
