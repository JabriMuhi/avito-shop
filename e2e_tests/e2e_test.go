package e2e_test

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
	"log"
	"net"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis       *bufconn.Listener
	client    avito.AvitoShopClient
	db        *sql.DB
	authToken string
)

func init() {
	var err error

	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.CreateAuthInterceptor()))

	db, err = sql.Open("postgres", "user=postgres password=password dbname=shop host=localhost port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}

	txCommitter := db2.NewTxProvider(db)

	merchRepo := repository.NewMerchRepository(txCommitter)
	merchSlice, err := merchRepo.GetAllMerch()

	if err != nil {
		log.Fatal(err)
	}

	merchMap := tools.SliceToMap(merchSlice)

	userRepo := repository.NewUserRepo(txCommitter)
	userService := service.NewUserService(userRepo)
	purchaseRepo := repository.NewPurchaseRepo(txCommitter, merchMap)
	purchaseService := service.NewPurchaseService(purchaseRepo)
	transactionRepo := repository.NewTransactionRepo(txCommitter)
	transactionService := service.NewTransactionService(transactionRepo)
	authService := service.NewAuthService(userRepo)

	handler := app.NewService(*userService, *purchaseService, *transactionService, *authService)
	avito.RegisterAvitoShopServer(server, handler)

	go server.Serve(lis)

	conn, _ := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	client = avito.NewAvitoShopClient(conn)
}

func getAuthContext() context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", authToken))
}

func TestFullFlow(t *testing.T) {
	testUserFullFlowName1 := "test_full_flow_1"
	testUserFullFlowName2 := "test_full_flow_2"

	defer cleanUsersData(testUserFullFlowName1, testUserFullFlowName2)

	resp2, err := client.Authenticate(context.Background(), &avito.AuthRequest{
		Username: testUserFullFlowName1,
		Password: "password",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp2.Token)

	resp, err := client.Authenticate(context.Background(), &avito.AuthRequest{
		Username: testUserFullFlowName2,
		Password: "password",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)

	authToken = resp.Token

	ctx := getAuthContext()

	respInfo, err := client.GetInfo(ctx, &avito.InfoRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, respInfo)

	respSendCoint, err := client.SendCoin(ctx, &avito.SendCoinRequest{
		ToUser: testUserFullFlowName1,
		Amount: 10,
	})
	assert.NoError(t, err)
	assert.NotNil(t, respSendCoint)

	respBuyItem, err := client.BuyItem(ctx, &avito.BuyItemRequest{Item: "cup"})
	assert.NoError(t, err)
	assert.NotNil(t, respBuyItem)
}

func cleanUsersData(user1, user2 string) error {
	var id int64

	err := db.QueryRow(`SELECT id FROM users WHERE name=$1`, user1).Scan(&id)

	var id2 int64

	err = db.QueryRow(`SELECT id FROM users WHERE name=$1`, user2).Scan(&id2)

	_, err = db.Exec(`
		DELETE FROM users WHERE id=$1 OR id=$2;
	`, id, id2)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DELETE FROM transactions WHERE sender_id=$1 OR receiver_id=$1 OR receiver_id=$2 OR sender_id=$2;
	`, id, id2)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DELETE FROM purchases WHERE user_id=$1 OR user_id=$2;
	`, id, id2)

	return err
}
