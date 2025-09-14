package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"purchase-service/api/presenter"
	"time"
)

type Client struct {
	baseURL        string
	httpClient     *http.Client
	internalSecret string
}

func NewClient(baseURL, internalSecret string) *Client {
	return &Client{
		baseURL:        baseURL,
		internalSecret: internalSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// setInternalHeaders sets the required headers for internal service communication
func (c *Client) setInternalHeaders(req *http.Request, userID string) {
	req.Header.Set("X-Secret", c.internalSecret)
	req.Header.Set("X-Auth-Gateway", "backend-infra")
	req.Header.Set("X-User-ID", userID)
}

// GetUserDetail fetches user details from user service
func (c *Client) GetUserDetail(ctx context.Context, userID, authenticatedUserID string) (*presenter.ExternalUserResponse, error) {
	url := fmt.Sprintf("%s/user/%s", c.baseURL, userID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers for internal service communication
	c.setInternalHeaders(req, authenticatedUserID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned status %d: %s", resp.StatusCode, string(body))
	}

	var user presenter.ExternalUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// GetProductDetail fetches product details from product service
func (c *Client) GetProductDetail(ctx context.Context, productID, authenticatedUserID string) (*presenter.ProductResponse, error) {
	url := fmt.Sprintf("%s/product/%s", c.baseURL, productID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers for internal service communication
	c.setInternalHeaders(req, authenticatedUserID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product service returned status %d: %s", resp.StatusCode, string(body))
	}

	var product presenter.ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &product, nil
}

// GetProductDetails fetches multiple product details in parallel
func (c *Client) GetProductDetails(ctx context.Context, productIDs []string, authenticatedUserID string) (map[string]*presenter.ProductResponse, error) {
	type result struct {
		productID string
		product   *presenter.ProductResponse
		err       error
	}

	results := make(chan result, len(productIDs))
	
	// Fetch all products in parallel
	for _, productID := range productIDs {
		go func(id string) {
			product, err := c.GetProductDetail(ctx, id, authenticatedUserID)
			results <- result{productID: id, product: product, err: err}
		}(productID)
	}

	// Collect results
	productMap := make(map[string]*presenter.ProductResponse)
	for i := 0; i < len(productIDs); i++ {
		res := <-results
		if res.err != nil {
			return nil, fmt.Errorf("failed to fetch product %s: %w", res.productID, res.err)
		}
		productMap[res.productID] = res.product
	}

	return productMap, nil
}

// DecreaseProductQuantity decreases the quantity of a product
func (c *Client) DecreaseProductQuantity(ctx context.Context, productID string, quantity int, authenticatedUserID string) error {
	url := fmt.Sprintf("%s/product/%s/decrease-quantity", c.baseURL, productID)
	
	requestBody := map[string]interface{}{
		"quantity": quantity,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	c.setInternalHeaders(req, authenticatedUserID)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("product service returned status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}
