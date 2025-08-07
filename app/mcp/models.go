package mcp

type ListAccountsParams struct{}

type CreateTransactionParams struct {
	Type            string `json:"type" jsonschema:"The type of transaction. Must be 'withdrawal', 'deposit', or 'transfer'."`
	Description     string `json:"description" jsonschema:"A clear and concise description of the transaction."`
	Amount          string `json:"amount" jsonschema:"The monetary value of the transaction as a string (e.g., '123.45')."`
	Date            string `json:"date" jsonschema:"The date of the transaction in YYYY-MM-DD format. Defaults to today if not provided."`
	SourceAccountID string `json:"source_account_id" jsonschema:"The ID of the source account (for withdrawals and transfers)."`
	DestAccountID   string `json:"dest_account_id" jsonschema:"The ID of the destination account (for deposits and transfers)."`
	Category        string `json:"category,omitempty" jsonschema:"The expense or income category for this transaction."`
}
