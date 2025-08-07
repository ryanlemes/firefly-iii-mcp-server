package firefly

import "time"

// Account represents a financial account in the business domain.
type Account struct {
	ID             string
	Name           string
	Type           string
	CurrentBalance string
	CurrencyCode   string
	Active         bool
}

// NewTransaction represents the information needed to create a transaction.
type NewTransaction struct {
	Type            string
	Date            time.Time
	Amount          string
	Description     string
	SourceAccountID string
	DestAccountID   string
	Category        string
}

// Transaction represents a created transaction in the business domain.
type Transaction struct {
	ID          string
	Description string
	Amount      string
	Date        time.Time
}
