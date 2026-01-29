package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terminal-echoware/pkg/types"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) CallAPI(req types.APIRequest) (*types.APIResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if apiResp.Error != "" {
		return nil, fmt.Errorf("api error: %s", apiResp.Error)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error: %d - %s", resp.StatusCode, apiResp.Error)
	}

	return &apiResp, nil
}

func (c *Client) ListProducts(params types.ProductListParams) ([]types.Product, int, error) {
	req := types.APIRequest{
		Type:      types.OperationTypeQuery,
		Operation: "product.list",
		Params:    params,
	}

	resp, err := c.CallAPI(req)
	if err != nil {
		return nil, 0, err
	}

	productsData, ok := resp.Data.([]interface{})
	if !ok {
		productsJSON, _ := json.Marshal(resp.Data)
		var products []types.Product
		if err := json.Unmarshal(productsJSON, &products); err != nil {
			return nil, 0, fmt.Errorf("unmarshal products: %w", err)
		}
		return products, resp.Count, nil
	}

	productsJSON, _ := json.Marshal(productsData)
	var products []types.Product
	if err := json.Unmarshal(productsJSON, &products); err != nil {
		return nil, 0, fmt.Errorf("unmarshal products: %w", err)
	}

	return products, resp.Count, nil
}

func (c *Client) GetProduct(id string) (*types.Product, error) {
	req := types.APIRequest{
		Type:      types.OperationTypeQuery,
		Operation: "product.get",
		Params:    types.ProductGetParams{ID: id},
	}

	resp, err := c.CallAPI(req)
	if err != nil {
		return nil, err
	}

	productJSON, _ := json.Marshal(resp.Data)
	var product types.Product
	if err := json.Unmarshal(productJSON, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product: %w", err)
	}

	return &product, nil
}

func (c *Client) SearchProducts(params types.ProductSearchParams) ([]types.Product, int, error) {
	req := types.APIRequest{
		Type:      types.OperationTypeQuery,
		Operation: "product.search",
		Params:    params,
	}

	resp, err := c.CallAPI(req)
	if err != nil {
		return nil, 0, err
	}

	productsData, ok := resp.Data.([]interface{})
	if !ok {
		productsJSON, _ := json.Marshal(resp.Data)
		var products []types.Product
		if err := json.Unmarshal(productsJSON, &products); err != nil {
			return nil, 0, fmt.Errorf("unmarshal products: %w", err)
		}
		return products, resp.Count, nil
	}

	productsJSON, _ := json.Marshal(productsData)
	var products []types.Product
	if err := json.Unmarshal(productsJSON, &products); err != nil {
		return nil, 0, fmt.Errorf("unmarshal products: %w", err)
	}

	return products, resp.Count, nil
}

func (c *Client) ListCategories(params types.CategoryListParams) ([]types.Category, int, error) {
	req := types.APIRequest{
		Type:      types.OperationTypeQuery,
		Operation: "category.list",
		Params:    params,
	}

	resp, err := c.CallAPI(req)
	if err != nil {
		return nil, 0, err
	}

	categoriesData, ok := resp.Data.([]interface{})
	if !ok {
		categoriesJSON, _ := json.Marshal(resp.Data)
		var categories []types.Category
		if err := json.Unmarshal(categoriesJSON, &categories); err != nil {
			return nil, 0, fmt.Errorf("unmarshal categories: %w", err)
		}
		return categories, resp.Count, nil
	}

	categoriesJSON, _ := json.Marshal(categoriesData)
	var categories []types.Category
	if err := json.Unmarshal(categoriesJSON, &categories); err != nil {
		return nil, 0, fmt.Errorf("unmarshal categories: %w", err)
	}

	return categories, resp.Count, nil
}

func (c *Client) CreateOrder(params types.OrderCreateParams) (*types.Order, error) {
	req := types.APIRequest{
		Type:      types.OperationTypeMutation,
		Operation: "order.create",
		Params:    params,
	}

	resp, err := c.CallAPI(req)
	if err != nil {
		return nil, err
	}

	orderJSON, _ := json.Marshal(resp.Data)
	var order types.Order
	if err := json.Unmarshal(orderJSON, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order: %w", err)
	}

	return &order, nil
}
