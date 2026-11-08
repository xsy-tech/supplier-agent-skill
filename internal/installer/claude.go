package installer

import "context"

// Deprecated: new installations use ClientConfig.Register for desktop-ready authentication.
func RegisterClaudeCode(ctx context.Context, runner CommandRunner, endpoint string) error {
	return runner.Run(ctx, "claude", "mcp", "add", "--scope", "user", "--transport", "http", "supplier-garment", endpoint, "--header", "Authorization: Bearer ${SUPPLIER_OPEN_CLAUDE_KEY}")
}

func UnregisterClaudeCode(ctx context.Context, runner CommandRunner) error {
	return runner.Run(ctx, "claude", "mcp", "remove", "--scope", "user", "supplier-garment")
}
