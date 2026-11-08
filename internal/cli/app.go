package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-saas/supplier-agent-skill/internal/credential"
	"github.com/go-saas/supplier-agent-skill/internal/domain"
	"github.com/go-saas/supplier-agent-skill/internal/installer"
	"github.com/go-saas/supplier-agent-skill/internal/mcpclient"
	"github.com/go-saas/supplier-agent-skill/internal/profile"
	secureRunner "github.com/go-saas/supplier-agent-skill/internal/runner"
)

type Doctor interface {
	Doctor(context.Context, string, string) (mcpclient.DoctorResult, error)
}

type Launcher func(context.Context, string, []string, string, string) error

type App struct {
	Domain          domain.Config
	Version         string
	In              io.Reader
	Out             io.Writer
	Err             io.Writer
	HomeDir         string
	CodexHome       string
	ClaudeConfigDir string
	Credentials     credential.Store
	Profiles        profile.Store
	Doctor          Doctor
	Commands        installer.CommandRunner
	Launch          Launcher
}

func NewDefault(version string) (*App, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve user home: %w", err)
	}
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve user config directory: %w", err)
	}
	return &App{
		Version:         version,
		In:              os.Stdin,
		Out:             os.Stdout,
		Err:             os.Stderr,
		HomeDir:         home,
		CodexHome:       os.Getenv("CODEX_HOME"),
		ClaudeConfigDir: os.Getenv("CLAUDE_CONFIG_DIR"),
		Credentials:     credential.NewOSStore(),
		Profiles:        profile.NewFileStore(filepath.Join(configDirectory, "supplier-agent-skill", "profiles.json")),
		Doctor:          mcpclient.New(),
		Commands:        execCommandRunner{},
		Launch:          secureRunner.Run,
	}, nil
}

func (a *App) Run(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		a.writeHelp()
		return nil
	}
	switch args[0] {
	case "install", "doctor", "uninstall", "launch":
		var err error
		args, err = a.connectionArgs(args)
		if err != nil {
			return err
		}
	}
	switch args[0] {
	case "version":
		fmt.Fprintf(a.Out, "supplier-skill %s\n", a.Version)
		return nil
	case "install":
		return a.runInstall(ctx, args[1:])
	case "doctor":
		return a.runDoctor(ctx, args[1:])
	case "launch":
		return a.runLaunch(ctx, args[1:])
	case "uninstall":
		return a.runUninstall(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a *App) writeHelp() {
	fmt.Fprintln(a.Out, `supplier-skill 用于安装、验证和启动 供应商开放平台的 Agent 集成。
Install, verify, and launch Supplier Open Platform agent integrations (fabric / garment).
安装会保存用户级 MCP 鉴权配置；重启后照常使用客户端。
Installation saves user-level MCP authentication; restart and use your client normally.

命令 / Commands:
  install codex|claude-code|workbuddy --domain fabric|garment --url URL --key-stdin [--force]
  install workbuddy --domain fabric|garment --url URL --output DIRECTORY [--force]
  doctor codex|claude-code|workbuddy --domain fabric|garment [--url URL]
  doctor workbuddy --domain fabric|garment --url URL --key-stdin (可选独立检查 / optional independent check)
  launch codex|claude-code --domain fabric|garment [-- agent arguments]
  uninstall codex|claude-code|workbuddy --domain fabric|garment [--url URL]
  version`)
}

func readSecret(reader io.Reader) (string, error) {
	content, err := io.ReadAll(io.LimitReader(reader, 16*1024))
	if err != nil {
		return "", errors.New("read API key from stdin")
	}
	secret := strings.TrimSpace(string(content))
	if secret == "" {
		return "", errors.New("API key from stdin is empty")
	}
	if strings.ContainsAny(secret, "\x00\r\n") {
		return "", errors.New("API key from stdin is invalid")
	}
	return secret, nil
}

func redactError(err error, secret string) error {
	if err == nil {
		return nil
	}
	message := err.Error()
	if secret != "" {
		message = strings.ReplaceAll(message, secret, "[REDACTED]")
	}
	return errors.New(message)
}

type execCommandRunner struct{}

func (execCommandRunner) Run(ctx context.Context, name string, args ...string) error {
	command := exec.CommandContext(ctx, name, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func (a *App) clientConfig() installer.ClientConfig {
	return installer.ClientConfig{HomeDir: a.HomeDir, CodexHome: a.CodexHome, ClaudeConfigDir: a.ClaudeConfigDir, Domain: a.Domain.Name}
}

func (a *App) storedCredential(agent, endpoint string) (string, error) {
	secret, err := a.clientConfig().Credential(agent, endpoint)
	if errors.Is(err, installer.ErrConfigCredentialNotFound) && agent != "workbuddy" {
		return a.Credentials.Get(a.profileKey(agent), endpoint)
	}
	return secret, err
}

func restoreProfile(store profile.Store, agent, previous string, existed bool) {
	if existed {
		_ = store.Set(agent, previous)
		return
	}
	_ = store.Delete(agent)
}

func (a *App) profileKey(agent string) string { return agent + "/" + a.Domain.Name }

// Strip our domain flag before handing the remaining arguments to the command.
// Client arguments after -- belong to the launched client and are preserved.
func (a *App) connectionArgs(args []string) ([]string, error) {
	name := ""
	cleaned := []string{args[0]}
	for i := 1; i < len(args); i++ {
		if args[i] == "--" {
			cleaned = append(cleaned, args[i:]...)
			break
		}
		if args[i] == "--domain" {
			if name != "" || i+1 >= len(args) {
				return nil, errors.New("exactly one --domain is required")
			}
			i++
			name = args[i]
			continue
		}
		if strings.HasPrefix(args[i], "--domain=") {
			if name != "" {
				return nil, errors.New("exactly one --domain is required")
			}
			name = strings.TrimPrefix(args[i], "--domain=")
			continue
		}
		cleaned = append(cleaned, args[i])
	}
	d, err := domain.Parse(name)
	if err != nil {
		return nil, err
	}
	a.Domain = d
	if len(cleaned) < 2 {
		return nil, errors.New("a client is required")
	}
	endpoint := ""
	hasURL := false
	for i := 2; i < len(cleaned); i++ {
		if cleaned[i] == "--" {
			break
		}
		if cleaned[i] == "--url" {
			if hasURL || i+1 >= len(cleaned) {
				return nil, errors.New("exactly one --url is allowed")
			}
			i++
			endpoint = cleaned[i]
			hasURL = true
		} else if strings.HasPrefix(cleaned[i], "--url=") {
			if hasURL {
				return nil, errors.New("exactly one --url is allowed")
			}
			endpoint = strings.TrimPrefix(cleaned[i], "--url=")
			hasURL = true
		}
	}
	if !hasURL && (args[0] == "doctor" || args[0] == "uninstall") {
		endpoint, err = a.Profiles.Get(a.profileKey(cleaned[1]))
		if err != nil {
			return nil, err
		}
		cleaned = append(cleaned, "--url", endpoint)
		hasURL = true
	}
	if hasURL {
		if err := d.ValidateEndpoint(endpoint); err != nil {
			return nil, err
		}
	}
	return cleaned, nil
}
