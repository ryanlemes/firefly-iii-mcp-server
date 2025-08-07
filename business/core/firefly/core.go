package firefly

import (
	"context"
	"log/slog"

	data "github.com/ryanlemes/firefly-iii-mcp-server/business/data"
	"github.com/ryanlemes/firefly-iii-mcp-server/business/data/firefly"
)

type Core struct {
	log         *slog.Logger
	account     data.AccountStorer
	transaction data.TransactionStorer
}

func NewCore(log *slog.Logger, accountStore data.AccountStorer, transaction data.TransactionStorer) *Core {
	return &Core{
		log:         log,
		account:     accountStore,
		transaction: transaction,
	}
}

func (c *Core) ListAccounts(ctx context.Context) ([]firefly.Account, error) {
	return c.account.List(ctx)
}

func (c *Core) CreateTransaction(ctx context.Context, newTx firefly.NewTransaction) (firefly.Transaction, error) {
	return c.transaction.Create(ctx, newTx)
}
