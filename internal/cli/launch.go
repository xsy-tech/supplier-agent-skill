package cli

import (
	"context"
	"errors"
	"strings"
)

func (a *App) runLaunch(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("launch requires codex or claude-code")
	}
	agent := args[0]
	command := ""
	environmentKey := ""
	switch agent {
	case "codex":
		command = "codex"
		environmentKey = "SUPPLIER_OPEN_CODEX_KEY"
	case "claude-code":
		command = "claude"
		environmentKey = "SUPPLIER_OPEN_CLAUDE_KEY"
	default:
		return errors.New("launch supports codex or claude-code")
	}
	endpoint, err := a.Profiles.Get(a.profileKey(agent))
	if err != nil {
		return err
	}
	secret, err := a.storedCredential(agent, endpoint)
	if err != nil {
		return err
	}
	agentArgs := args[1:]
	if len(agentArgs) > 0 && agentArgs[0] == "--" {
		agentArgs = agentArgs[1:]
	}
	if err := a.Launch(ctx, command, agentArgs, environmentKey+"_"+strings.ToUpper(a.Domain.Name), secret); err != nil {
		return redactError(err, secret)
	}
	return nil
}
