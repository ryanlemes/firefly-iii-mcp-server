package firefly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	token      string
	log        *slog.Logger
}

func NewClient(httpClient *http.Client, baseUrl string, token string, log *slog.Logger) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    u,
		token:      token,
		log:        log,
	}, nil
}

// List all accounts.
// Reference: https://api-docs.firefly-iii.org/#/accounts/listAccounts
// TODO add date string query parameter for filtering by a specific date.(optional)
// TODO add type string query parameter for filtering by account type. (optional)
func (c *Client) ListAccounts(ctx context.Context, limit int, page int) (PaginatedAccounts, error) {
	endpoint := c.baseURL.JoinPath("/api/v1/accounts")
	q := endpoint.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("page", fmt.Sprintf("%d", page))
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return PaginatedAccounts{}, fmt.Errorf("error creating request: %w", err)
	}

	var resp PaginatedAccounts
	if err := c.do(req, &resp); err != nil {
		return PaginatedAccounts{}, fmt.Errorf("error fetching accounts: %w", err)
	}

	return resp, nil
}

// Store a new transaction.
// Reference: https://api-docs.firefly-iii.org/#/transactions/storeTransaction
func (c *Client) CreateTransaction(ctx context.Context, tx TransactionStore) (SingleTransaction, error) {
	endpoint := c.baseURL.JoinPath("/api/v1/transactions")

	body, err := json.Marshal(tx)
	if err != nil {
		return SingleTransaction{}, fmt.Errorf("error marshaling transaction: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return SingleTransaction{}, fmt.Errorf("error creating request: %w", err)
	}

	var resp SingleTransaction
	if err := c.do(req, &resp); err != nil {
		return SingleTransaction{}, fmt.Errorf("error executing request: %w", err)
	}

	return resp, nil
}

func (c *Client) do(req *http.Request, v any) error {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error: status %s, body %s", resp.Status, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}
