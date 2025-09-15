package presenter

// Purchase Response DTOs
type PurchaseResponse struct {
	PurchaseID     string                `json:"purchaseId"`
	PurchasedItems []PurchaseItemResponse `json:"purchasedItems"`
	TotalPrice     float64               `json:"totalPrice"`
	PaymentDetails []PaymentDetail       `json:"paymentDetails"`
}

type GetPurchaseResponse struct {
	PurchaseID       string                `json:"purchaseId"`
	UserID           string                `json:"userId"`
	PaymentProofIds  []string              `json:"paymentProofIds"`
	PurchasedItems   []PurchaseItemResponse `json:"purchasedItems"`
	TotalPrice       float64               `json:"totalPrice"`
	PaymentDetails   []PaymentDetail       `json:"paymentDetails"`
	SenderInfo       SenderInfo            `json:"senderInfo"`
	CreatedAt        string                `json:"createdAt"`
	UpdatedAt        string                `json:"updatedAt"`
}

type ListPurchasesResponse struct {
	Purchases []GetPurchaseResponse `json:"purchases"`
	Total     int                   `json:"total"`
	Page      int                   `json:"page"`
	Limit     int                   `json:"limit"`
}

type PurchaseItemResponse struct {
	ProductID         string  `json:"productId"`
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	Qty               int     `json:"qty"`
	Price             float64 `json:"price"`
	SKU               string  `json:"sku"`
	FileID            string  `json:"fileId"`
	FileURI           string  `json:"fileUri"`
	FileThumbnailURI  string  `json:"fileThumbnailUri"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

type PaymentDetail struct {
	BankAccountName   string  `json:"bankAccountName"`
	BankAccountHolder string  `json:"bankAccountHolder"`
	BankAccountNumber string  `json:"bankAccountNumber"`
	TotalPrice        float64 `json:"totalPrice"`
}

type SenderInfo struct {
	SenderName          string `json:"senderName"`
	SenderContactType   string `json:"senderContactType"`
	SenderContactDetail string `json:"senderContactDetail"`
}
