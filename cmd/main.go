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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Fatal(err)
	}

	dbUser := "postgres"
	dbPassword := "password"
	dbName := "shop"
	dbHost := "db"
	dbPort := "5432"
	sslMode := "disable"

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbPort, sslMode)

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

	err = avito.RegisterAvitoShopHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
	if err != nil {
		log.Fatal("Failed to register gRPC Gateway:", err)
	}

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting HTTP server on :8080")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
