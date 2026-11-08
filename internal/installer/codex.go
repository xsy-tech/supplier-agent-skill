package installer

import "context"

type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) error
}

// Deprecated: new installations use ClientConfig.Register for desktop-ready authentication.
func RegisterCodex(ctx context.Context, runner CommandRunner, endpoint string) error {
	return runner.Run(ctx, "codex", "mcp", "add", "supplier-garment", "--url", endpoint, "--bearer-token-env-var", "SUPPLIER_OPEN_CODEX_KEY")
}

func UnregisterCodex(ctx context.Context, runner CommandRunner) error {
	return runner.Run(ctx, "codex", "mcp", "remove", "supplier-garment")
}
