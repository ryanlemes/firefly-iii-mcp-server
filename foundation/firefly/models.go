package firefly

import "time"

// AccountRead represents a single account returned by the API.
type AccountRead struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
		Name            string    `json:"name"`
		Type            string    `json:"type"`
		IBAN            *string   `json:"iban"`
		BIC             *string   `json:"bic"`
		AccountNumber   *string   `json:"account_number"`
		CurrentBalance  string    `json:"current_balance"`
		CurrencyCode    string    `json:"currency_code"`
		CurrencySymbol  string    `json:"currency_symbol"`
		Active          bool      `json:"active"`
		Order           int32     `json:"order"`
		IncludeNetWorth bool      `json:"include_net_worth"`
	} `json:"attributes"`
}

// TransactionStore defines the structure for creating a new transaction.
type TransactionStore struct {
	ErrorIfDuplicateHash bool             `json:"error_if_duplicate_hash,omitempty"`
	ApplyRules           bool             `json:"apply_rules,omitempty"`
	FireWebhooks         bool             `json:"fire_webhooks,omitempty"`
	Transactions         TransactionSplit `json:"transactions"`
}

// TransactionSplit represents a single entry in a transaction journal.
type TransactionSplit struct {
	Type              string  `json:"type"`
	Date              string  `json:"date"` // Format: YYYY-MM-DD or RFC3339
	Amount            string  `json:"amount"`
	Description       string  `json:"description"`
	SourceID          *string `json:"source_id,omitempty"`
	SourceName        *string `json:"source_name,omitempty"`
	DestinationID     *string `json:"destination_id,omitempty"`
	DestinationName   *string `json:"destination_name,omitempty"`
	CategoryName      *string `json:"category_name,omitempty"`
	BudgetName        *string `json:"budget_name,omitempty"`
	Tags              string  `json:"tags,omitempty"`
	Reconciled        *bool   `json:"reconciled,omitempty"`
	ExternalURL       *string `json:"external_url,omitempty"`
	InternalReference *string `json:"internal_reference,omitempty"`
}

// TransactionRead represents a full transaction journal returned by the API.
type TransactionRead struct {
	ID         string `json:"id"`
	Attributes struct {
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		User         string    `json:"user"`
		GroupTitle   *string   `json:"group_title"`
		Transactions []struct {
			TransactionJournalID string    `json:"transaction_journal_id"`
			Type                 string    `json:"type"`
			Date                 time.Time `json:"date"`
			Amount               string    `json:"amount"`
			Description          string    `json:"description"`
			SourceID             string    `json:"source_id"`
			SourceName           string    `json:"source_name"`
			DestinationID        string    `json:"destination_id"`
			DestinationName      string    `json:"destination_name"`
		} `json:"transactions"`
	} `json:"attributes"`
}

// PaginatedAccounts represents the structure for a paginated list of accounts.
type PaginatedAccounts struct {
	Data []AccountRead `json:"data"`
	Meta struct {
		Pagination struct {
			Total       int `json:"total"`
			Count       int `json:"count"`
			PerPage     int `json:"per_page"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
		} `json:"pagination"`
	} `json:"meta"`
}

// SingleTransaction represents the structure for a single transaction response.
type SingleTransaction struct {
	Data TransactionRead `json:"data"`
}
