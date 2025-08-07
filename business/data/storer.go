package data

import (
	"context"

	"github.com/ryanlemes/firefly-iii-mcp-server/business/data/firefly"
)

type AccountStorer interface {
	List(ctx context.Context) ([]firefly.Account, error)
}

type TransactionStorer interface {
	Create(ctx context.Context, newTx firefly.NewTransaction) (firefly.Transaction, error)
}
