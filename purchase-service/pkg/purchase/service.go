package purchase

import (
	"context"
	"encoding/json"
	"fmt"
	"purchase-service/api/presenter"
	"purchase-service/pkg/dtos"
	"purchase-service/pkg/entities"
	"purchase-service/pkg/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreatePurchase(ctx context.Context, req dtos.CreatePurchaseRequest) (*presenter.PurchaseResponse, error)
	UploadPaymentProof(ctx context.Context, purchaseID string, req dtos.PaymentProofRequest) error
	GetPurchaseByID(ctx context.Context, purchaseID string) (*presenter.GetPurchaseResponse, error)
	ListPurchases(ctx context.Context, page, limit int) (*presenter.ListPurchasesResponse, error)
}

type service struct {
	repo         Repository
	userClient   *http.Client
	productClient *http.Client
}

func NewService(repo Repository, userServiceURL, productServiceURL, internalSecret string) Service {
	return &service{
		repo:          repo,
		userClient:    http.NewClient(userServiceURL, internalSecret),
		productClient: http.NewClient(productServiceURL, internalSecret),
	}
}

func (s *service) CreatePurchase(ctx context.Context, req dtos.CreatePurchaseRequest) (*presenter.PurchaseResponse, error) {
	// Get authenticated user ID from context (set by gateway trust middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user context not found")
	}

	// Extract product IDs for batch fetching
	productIDs := make([]string, len(req.PurchasedItems))
	for i, item := range req.PurchasedItems {
		productIDs[i] = item.ProductID
	}

	// Fetch all products in parallel
	products, err := s.productClient.GetProductDetails(ctx, productIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	// Validate all products exist and collect seller IDs
	sellerIDs := make(map[string]bool)
	for _, item := range req.PurchasedItems {
		product, exists := products[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("product %s not found", item.ProductID)
		}
		sellerIDs[product.SellerID] = true
	}

	// Fetch seller details in parallel
	sellerDetails := make(map[string]*presenter.SellerResponse)
	for sellerID := range sellerIDs {
		user, err := s.userClient.GetUserDetail(ctx, sellerID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch seller %s: %w", sellerID, err)
		}
		// Convert ExternalUserResponse to SellerResponse
		sellerDetails[sellerID] = &presenter.SellerResponse{
			ID:                user.ID,
			BankAccountName:   user.BankAccountName,
			BankAccountHolder: user.BankAccountHolder,
			BankAccountNumber: user.BankAccountNumber,
		}
	}

	// Create purchase entity
	purchase := &entities.Purchase{
		UserID: uuid.MustParse(userID),
	}
	if err := s.repo.CreatePurchase(ctx, purchase); err != nil {
		return nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	// Create purchase items with copied product information
	var purchaseItems []*entities.PurchaseItem
	var totalPrice float64
	sellerTotals := make(map[string]float64)

	for _, item := range req.PurchasedItems {
		product := products[item.ProductID]
		
		// Copy product information (source of truth)
		purchaseItem := &entities.PurchaseItem{
			PurchaseID:        purchase.ID,
			ProductID:         product.ID,
			Name:              product.Name,
			Category:          product.Category,
			Qty:               product.Qty, // Original quantity before purchase
			Price:             product.Price,
			SKU:               product.SKU,
			FileID:            product.FileID,
			FileURI:           product.FileURI,
			FileThumbnailURI:  product.FileThumbnailURI,
		}

		purchaseItems = append(purchaseItems, purchaseItem)
		
		// Calculate totals
		itemTotal := product.Price * float64(item.Qty)
		totalPrice += itemTotal
		sellerTotals[product.SellerID] += itemTotal
	}

	// Save purchase items
	if err := s.repo.CreatePurchaseItems(ctx, purchaseItems); err != nil {
		return nil, fmt.Errorf("failed to create purchase items: %w", err)
	}

	// Create purchase sender
	sender := &entities.PurchaseSender{
		PurchaseID:          purchase.ID,
		SenderName:          req.SenderName,
		SenderContactType:   req.SenderContactType,
		SenderContactDetail: req.SenderContactDetail,
	}

	if err := s.repo.CreatePurchaseSender(ctx, sender); err != nil {
		return nil, fmt.Errorf("failed to create purchase sender: %w", err)
	}

	// Build payment details
	var paymentDetails []presenter.PaymentDetail
	for sellerID, total := range sellerTotals {
		seller := sellerDetails[sellerID]
		paymentDetails = append(paymentDetails, presenter.PaymentDetail{
			BankAccountName:   seller.BankAccountName,
			BankAccountHolder: seller.BankAccountHolder,
			BankAccountNumber: seller.BankAccountNumber,
			TotalPrice:        total,
		})
	}

	// Sort payment details by bank account name for consistency
	sort.Slice(paymentDetails, func(i, j int) bool {
		return paymentDetails[i].BankAccountName < paymentDetails[j].BankAccountName
	})

	// Build response
	var purchasedItems []presenter.PurchaseItemResponse
	for _, item := range purchaseItems {
		purchasedItems = append(purchasedItems, presenter.PurchaseItemResponse{
			ProductID:        item.ProductID,
			Name:             item.Name,
			Category:         item.Category,
			Qty:              item.Qty,
			Price:            item.Price,
			SKU:              item.SKU,
			FileID:           item.FileID,
			FileURI:          item.FileURI,
			FileThumbnailURI: item.FileThumbnailURI,
			CreatedAt:        item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &presenter.PurchaseResponse{
		PurchaseID:     purchase.ID.String(),
		PurchasedItems: purchasedItems,
		TotalPrice:     totalPrice,
		PaymentDetails: paymentDetails,
	}, nil
}

func (s *service) UploadPaymentProof(ctx context.Context, purchaseID string, req dtos.PaymentProofRequest) error {
	// Get authenticated user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return fmt.Errorf("user context not found")
	}

	// Validate purchase ID format
	if _, err := uuid.Parse(purchaseID); err != nil {
		return fmt.Errorf("invalid purchase ID format")
	}

	// Get purchase to verify it exists and belongs to the user
	purchase, err := s.repo.GetPurchaseByID(ctx, purchaseID)
	if err != nil {
		return fmt.Errorf("purchase not found: %w", err)
	}

	// Verify the purchase belongs to the authenticated user
	if purchase.UserID.String() != userID {
		return fmt.Errorf("unauthorized: purchase does not belong to user")
	}

	// Get purchase items to know which products to decrease
	_, err = s.repo.GetPurchaseItemsByPurchaseID(ctx, purchaseID)
	if err != nil {
		return fmt.Errorf("failed to get purchase items: %w", err)
	}

	// Decrease product quantities for each item
	// for _, item := range purchaseItems {
	// 	err := s.productClient.DecreaseProductQuantity(ctx, item.ProductID, item.Qty, userID)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to decrease quantity for product %s: %w", item.ProductID, err)
	// 	}
	// }

	// Convert file IDs to JSON string
	fileIdsJSON, err := json.Marshal(req.FileIds)
	if err != nil {
		return fmt.Errorf("failed to marshal file IDs: %w", err)
	}

	// Update purchase with payment proof
	if err := s.repo.UpdatePurchasePaymentProof(ctx, purchaseID, string(fileIdsJSON)); err != nil {
		return fmt.Errorf("failed to update purchase with payment proof: %w", err)
	}

	return nil
}

func (s *service) GetPurchaseByID(ctx context.Context, purchaseID string) (*presenter.GetPurchaseResponse, error) {
	// Get authenticated user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user context not found")
	}

	// Validate purchase ID format
	if _, err := uuid.Parse(purchaseID); err != nil {
		return nil, fmt.Errorf("invalid purchase ID format")
	}

	// Get purchase
	purchase, err := s.repo.GetPurchaseByID(ctx, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("purchase not found: %w", err)
	}

	// Verify the purchase belongs to the authenticated user
	if purchase.UserID.String() != userID {
		return nil, fmt.Errorf("unauthorized: purchase does not belong to user")
	}

	// Get purchase items
	purchaseItems, err := s.repo.GetPurchaseItemsByPurchaseID(ctx, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase items: %w", err)
	}

	// Get purchase sender
	sender, err := s.repo.GetPurchaseSenderByPurchaseID(ctx, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase sender: %w", err)
	}

	// Parse payment proof IDs
	var paymentProofIds []string
	if purchase.PaymentProofIds != "" {
		if err := json.Unmarshal([]byte(purchase.PaymentProofIds), &paymentProofIds); err != nil {
			// If parsing fails, just use empty array
			paymentProofIds = []string{}
		}
	}

	// Build purchase items response
	var purchasedItems []presenter.PurchaseItemResponse
	var totalPrice float64
	for _, item := range purchaseItems {
		purchasedItems = append(purchasedItems, presenter.PurchaseItemResponse{
			ProductID:         item.ProductID,
			Name:              item.Name,
			Category:          item.Category,
			Qty:               item.Qty,
			Price:             item.Price,
			SKU:               item.SKU,
			FileID:            item.FileID,
			FileURI:           item.FileURI,
			FileThumbnailURI:  item.FileThumbnailURI,
			CreatedAt:         item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         item.UpdatedAt.Format(time.RFC3339),
		})
		totalPrice += item.Price * float64(item.Qty)
	}

	// Build sender info
	senderInfo := presenter.SenderInfo{
		SenderName:          sender.SenderName,
		SenderContactType:   sender.SenderContactType,
		SenderContactDetail: sender.SenderContactDetail,
	}

	return &presenter.GetPurchaseResponse{
		PurchaseID:      purchase.ID.String(),
		UserID:          purchase.UserID.String(),
		PaymentProofIds: paymentProofIds,
		PurchasedItems:  purchasedItems,
		TotalPrice:      totalPrice,
		PaymentDetails:  []presenter.PaymentDetail{}, // This would need to be fetched from external services
		SenderInfo:      senderInfo,
		CreatedAt:       purchase.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       purchase.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) ListPurchases(ctx context.Context, page, limit int) (*presenter.ListPurchasesResponse, error) {
	// Get authenticated user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user context not found")
	}

	// Set default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Get purchases from repository
	purchases, total, err := s.repo.GetPurchasesByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchases: %w", err)
	}

	// Convert to response DTOs
	var purchaseResponses []presenter.GetPurchaseResponse
	for _, purchase := range purchases {
		// Get purchase items
		purchaseItems, err := s.repo.GetPurchaseItemsByPurchaseID(ctx, purchase.ID.String())
		if err != nil {
			continue // Skip this purchase if we can't get items
		}

		// Get purchase sender
		sender, err := s.repo.GetPurchaseSenderByPurchaseID(ctx, purchase.ID.String())
		if err != nil {
			continue // Skip this purchase if we can't get sender
		}

		// Parse payment proof IDs
		var paymentProofIds []string
		if purchase.PaymentProofIds != "" {
			if err := json.Unmarshal([]byte(purchase.PaymentProofIds), &paymentProofIds); err != nil {
				paymentProofIds = []string{}
			}
		}

		// Build purchase items response
		var purchasedItems []presenter.PurchaseItemResponse
		var totalPrice float64
		for _, item := range purchaseItems {
			purchasedItems = append(purchasedItems, presenter.PurchaseItemResponse{
				ProductID:         item.ProductID,
				Name:              item.Name,
				Category:          item.Category,
				Qty:               item.Qty,
				Price:             item.Price,
				SKU:               item.SKU,
				FileID:            item.FileID,
				FileURI:           item.FileURI,
				FileThumbnailURI:  item.FileThumbnailURI,
				CreatedAt:         item.CreatedAt.Format(time.RFC3339),
				UpdatedAt:         item.UpdatedAt.Format(time.RFC3339),
			})
			totalPrice += item.Price * float64(item.Qty)
		}

		// Build sender info
		senderInfo := presenter.SenderInfo{
			SenderName:          sender.SenderName,
			SenderContactType:   sender.SenderContactType,
			SenderContactDetail: sender.SenderContactDetail,
		}

		purchaseResponses = append(purchaseResponses, presenter.GetPurchaseResponse{
			PurchaseID:      purchase.ID.String(),
			PaymentProofIds: paymentProofIds,
			PurchasedItems:  purchasedItems,
			TotalPrice:      totalPrice,
			PaymentDetails:  []presenter.PaymentDetail{}, // This would need to be fetched from external services
			SenderInfo:      senderInfo,
			CreatedAt:       purchase.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       purchase.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &presenter.ListPurchasesResponse{
		Purchases: purchaseResponses,
		Total:     int(total),
		Page:      page,
		Limit:     limit,
	}, nil
}
