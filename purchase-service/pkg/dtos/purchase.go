package dtos

// API Request DTOs
type CreatePurchaseRequest struct {
	PurchasedItems []PurchaseItemRequest `json:"purchasedItems" validate:"required,min=1,dive"`
	SenderName     string                `json:"senderName" validate:"required,min=4,max=55"`
	SenderContactType    string          `json:"senderContactType" validate:"required,oneof=email phone"`
	SenderContactDetail  string          `json:"senderContactDetail" validate:"required"`
}

type PurchaseItemRequest struct {
	ProductID string `json:"productId" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=1"`
}

type PaymentProofRequest struct {
	FileIds []string `json:"fileIds" validate:"required,min=1,dive,required"`
}
