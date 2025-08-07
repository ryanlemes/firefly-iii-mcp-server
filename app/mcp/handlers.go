package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ryanlemes/firefly-iii-mcp-server/business/core/firefly"
	data "github.com/ryanlemes/firefly-iii-mcp-server/business/data/firefly"
)

type Handlers struct {
	log  *slog.Logger
	core *firefly.Core
}

func NewHandlers(log *slog.Logger, core *firefly.Core) *Handlers {
	return &Handlers{
		log:  log,
		core: core,
	}
}

func (h *Handlers) ListAccounts(ctx context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[ListAccountsParams]) (*mcp.CallToolResultFor[any], error) {
	accounts, err := h.core.ListAccounts(ctx)
	if err != nil {
		h.log.Error("failed to list accounts", "error", err)
	}

	jsonData, err := json.Marshal(accounts)
	if err != nil {
		return nil, fmt.Errorf("marshaling accounts to json: %w", err)
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil
}

func (h *Handlers) CreateTransaction(ctx context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[data.NewTransaction]) (*mcp.CallToolResultFor[any], error) {
	args := params.Arguments

	date, err := time.Parse("2006-01-02", args.Date.String())
	if err != nil {
		// If parsing fails, default to current date.
		date = time.Now()
	}

	newTx := data.NewTransaction{
		Type:            args.Type,
		Date:            date,
		Amount:          args.Amount,
		Description:     args.Description,
		SourceAccountID: args.SourceAccountID,
		DestAccountID:   args.DestAccountID,
		Category:        args.Category,
	}

	tx, err := h.core.CreateTransaction(ctx, newTx)
	if err != nil {
		h.log.Error("failed to create transaction", "error", err)
		return nil, fmt.Errorf("creating transaction: %w", err)
	}

	resultText := fmt.Sprintf("transaction '%s' created with ID: %s", tx.Description, tx.ID)

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultText},
		},
	}, nil
}
