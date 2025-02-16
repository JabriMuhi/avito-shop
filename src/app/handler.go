package app

import (
	"avito-shop/avito"
	"avito-shop/src/models"
	"avito-shop/src/service"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	avito.UnimplementedAvitoShopServer
	userService        service.UserService
	purchaseService    service.PurchaseService
	transactionService service.TransactionService
	authService        service.AuthServiceImpl
}

func NewService(userService service.UserService, purchaseService service.PurchaseService, transactionService service.TransactionService, authService service.AuthServiceImpl) *Handler {
	return &Handler{userService: userService, purchaseService: purchaseService, transactionService: transactionService, authService: authService}
}

func (h *Handler) GetInfo(ctx context.Context, req *avito.InfoRequest) (*avito.InfoResponse, error) {
	userID, ok := ctx.Value("user").(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	userInfo, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	inventory, err := h.purchaseService.GetPurchasedMerchByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	senderHistory, receivedHistory, err := h.transactionService.GetCoinTransactionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	receivedTransactions := []*avito.Transaction{}

	for name, amount := range receivedHistory {
		receivedTransaction := avito.Transaction{
			User:   name,
			Amount: amount,
		}

		receivedTransactions = append(receivedTransactions, &receivedTransaction)
	}

	senderTransactions := []*avito.Transaction{}

	for name, amount := range senderHistory {
		senderTransaction := avito.Transaction{
			User:   name,
			Amount: amount,
		}

		senderTransactions = append(senderTransactions, &senderTransaction)
	}

	coinHistory := avito.CoinHistory{
		Received: receivedTransactions,
		Sent:     senderTransactions,
	}

	return &avito.InfoResponse{
		Coins:       userInfo.Balance,
		Inventory:   inventory,
		CoinHistory: &coinHistory,
	}, nil
}

func (h *Handler) SendCoin(ctx context.Context, req *avito.SendCoinRequest) (*avito.SendCoinResponse, error) {
	userSenderID, ok := ctx.Value("user").(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	userReceiver, err := h.userService.GetUserByName(ctx, req.ToUser)
	if err != nil {
		if errors.Is(err, models.ErrNoSuchUser) || errors.Is(err, models.ErrEmptyName) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	err = h.transactionService.TransferCoins(ctx, userSenderID, userReceiver.ID, int64(req.Amount))
	if err != nil {
		if errors.Is(err, models.ErrInvalidAmount) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (h *Handler) BuyItem(ctx context.Context, req *avito.BuyItemRequest) (*avito.BuyItemResponse, error) {
	userID, ok := ctx.Value("user").(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	err := h.purchaseService.BuyMerch(ctx, userID, req.Item)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return nil, nil
}

func (h *Handler) Authenticate(ctx context.Context, req *avito.AuthRequest) (*avito.AuthResponse, error) {
	return h.authService.Authenticate(ctx, req)
}
