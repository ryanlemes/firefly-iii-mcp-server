package firefly_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ryanlemes/firefly-iii-mcp-server/foundation/firefly"
)

func TestNewClient(t *testing.T) {
	t.Run("should return a new client", func(t *testing.T) {
		client, err := firefly.NewClient(http.DefaultClient, "http://localhost", "token", slog.New(slog.NewJSONHandler(io.Discard, nil)))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("client is nil")
		}
	})

	t.Run("should return an error if the base url is invalid", func(t *testing.T) {
		_, err := firefly.NewClient(http.DefaultClient, "://invalid url", "token", slog.New(slog.NewJSONHandler(io.Discard, nil)))
		if err == nil {
			t.Fatal("expected an error, but got nil")
		}
	})
}

func TestClient_ListAccounts(t *testing.T) {
	t.Run("should return a list of accounts", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/accounts" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data": [{"id": "1", "type": "accounts", "attributes": {"name": "test account"}}]}`)
		}))
		defer server.Close()

		client, err := firefly.NewClient(server.Client(), server.URL, "token", slog.New(slog.NewJSONHandler(io.Discard, nil)))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		accounts, err := client.ListAccounts(context.Background(), 10, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(accounts.Data) != 1 {
			t.Fatalf("expected 1 account, but got %d", len(accounts.Data))
		}

		if accounts.Data[0].Attributes.Name != "test account" {
			t.Fatalf("unexpected account name: %s", accounts.Data[0].Attributes.Name)
		}
	})
}

func TestClient_CreateTransaction(t *testing.T) {
	t.Run("should create a new transaction", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/transactions" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data": {"id": "1", "type": "transactions", "attributes": {"transactions": {"description": "test transaction"}}}}`)
		}))
		defer server.Close()

		client, err := firefly.NewClient(server.Client(), server.URL, "token", slog.New(slog.NewJSONHandler(io.Discard, nil)))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tx := firefly.TransactionStore{
			Transactions: firefly.TransactionSplit{
				Type:        "withdrawal",
				Date:        "2025-01-01T12:00:00Z",
				Amount:      "100.00",
				Description: "test transaction",
			},
		}

		transaction, err := client.CreateTransaction(context.Background(), tx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if transaction.Data.Attributes.Transactions.Description != "test transaction" {
			t.Fatalf("unexpected transaction description: %s", transaction.Data.Attributes.Transactions.Description)
		}
	})

	t.Run("should return an error if the request fails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client, err := firefly.NewClient(server.Client(), server.URL, "token", slog.New(slog.NewJSONHandler(io.Discard, nil)))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = client.CreateTransaction(context.Background(), firefly.TransactionStore{})
		if err == nil {
			t.Fatal("expected an error, but got nil")
		}
	})
}
