package firefly

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ryanlemes/firefly-iii-mcp-server/foundation/firefly"
)

type Store struct {
	log    *slog.Logger
	client *firefly.Client
}

func NewStore(log *slog.Logger, client *firefly.Client) *Store {
	return &Store{
		log:    log,
		client: client,
	}
}

func (s *Store) List(ctx context.Context) ([]Account, error) {
	// TODO implement pagination and filtering options
	const page = 1
	const limit = 50
	paginated, err := s.client.ListAccounts(ctx, limit, page)
	if err != nil {
		return nil, fmt.Errorf("listing accounts from api: %w", err)
	}

	accounts := make([]Account, len(paginated.Data))
	for i, a := range paginated.Data {
		accounts[i] = toBusinessAccount(a)
	}

	return accounts, nil
}

// TODO Switch the create to return a single transaction instead of a slice.
func (s *Store) Create(ctx context.Context, newTx NewTransaction) (Transaction, error) {
	apiTx := toAPITransactionStore(newTx)

	result, err := s.client.CreateTransaction(ctx, apiTx)
	if err != nil {
		return Transaction{}, fmt.Errorf("creating transaction in api: %w", err)
	}
	return toBusinessTransaction(result.Data), nil
}

// =======================================================================
// Model Conversion Helpers

func toBusinessAccount(acc firefly.AccountRead) Account {
	return Account{
		ID:             acc.ID,
		Name:           acc.Attributes.Name,
		Type:           acc.Attributes.Type,
		CurrentBalance: acc.Attributes.CurrentBalance,
		CurrencyCode:   acc.Attributes.CurrencyCode,
		Active:         acc.Attributes.Active,
	}
}

func toAPITransactionStore(newTx NewTransaction) firefly.TransactionStore {
	return firefly.TransactionStore{
		ApplyRules:   true,
		FireWebhooks: true,
		Transactions: firefly.TransactionSplit{
			Type:          newTx.Type,
			Date:          newTx.Date.Format("2006-01-02"),
			Amount:        newTx.Amount,
			Description:   newTx.Description,
			SourceID:      &newTx.SourceAccountID,
			DestinationID: &newTx.DestAccountID,
			CategoryName:  &newTx.Category,
		},
	}
}

func toBusinessTransaction(tx firefly.TransactionRead) Transaction {
	split := tx.Attributes.Transactions

	if len(split) == 0 {
		return Transaction{}
	}

	// for _, entry := range split {
	// 	if entry.TransactionJournalID != tx.ID {
	// 		continue
	// 	}
	// 	return []Transaction{
	// 		{
	// 			ID:          entry.TransactionJournalID,
	// 			Description: entry.Description,
	// 			Amount:      entry.Amount,
	// 			Date:        entry.Date,
	// 		},
	// 	}
	// }

	return Transaction{
		ID:          tx.ID,
		Description: split[0].Description,
		Amount:      split[0].Amount,
		Date:        split[0].Date,
	}
}
