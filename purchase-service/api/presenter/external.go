package presenter

// External service response DTOs
type ProductResponse struct {
	ID       string  `json:"productId"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Qty      int     `json:"qty"`
	Price    float64 `json:"price"`
	SKU      string  `json:"sku"`
	FileID   string  `json:"fileId"`
	FileURI  string  `json:"fileUri"`
	FileThumbnailURI string `json:"fileThumbnailUri"`
	SellerID string  `json:"sellerId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ExternalUserResponse struct {
	ID                string `json:"id"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

type SellerResponse struct {
	ID                string `json:"id"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
}
