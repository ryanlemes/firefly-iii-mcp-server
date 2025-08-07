package main

import (
	"log/slog"
	"os"

	"github.com/ryanlemes/firefly-iii-mcp-server/foundation/logger"
)

func main() {
	// Construct the application logger.
	log := logger.New("firefly-iii-mcp-server")

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Error("startup", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	return nil
}
