package mcpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-saas/supplier-agent-skill/internal/domain"

	bundle "github.com/go-saas/supplier-agent-skill"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	ErrUnauthorized = errors.New("MCP endpoint rejected the API key")
	ErrNoTools      = errors.New("MCP endpoint returned no visible tools")
	ErrUnavailable  = errors.New("MCP endpoint is unavailable")
)

type DoctorStatus string

const (
	DoctorStatusVerified              DoctorStatus = "VERIFIED"
	DoctorStatusAuthorizationRequired DoctorStatus = "INSTALLED_AUTHORIZATION_REQUIRED"
)

type VerificationResult struct {
	Tool   string `json:"tool"`
	Output any    `json:"output"`
}

type DoctorResult struct {
	AccessDomain       string              `json:"access_domain"`
	RequiredCapability string              `json:"required_capability,omitempty"`
	Status             DoctorStatus        `json:"status"`
	ServerName         string              `json:"server_name"`
	ServerVersion      string              `json:"server_version"`
	Instructions       string              `json:"instructions,omitempty"`
	Tools              []string            `json:"tools"`
	Verification       *VerificationResult `json:"verification,omitempty"`
}

func (c *Client) Doctor(ctx context.Context, endpoint, key string) (DoctorResult, error) {
	d, err := domain.FromEndpoint(endpoint)
	if err != nil {
		return DoctorResult{}, err
	}
	parsed, err := url.Parse(strings.TrimSpace(endpoint))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return DoctorResult{}, errors.New("MCP endpoint must be an absolute URL")
	}
	if strings.TrimSpace(key) == "" {
		return DoctorResult{}, errors.New("MCP API key is required")
	}
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	doctorContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client := mcp.NewClient(&mcp.Implementation{Name: "supplier-agent-skill", Version: bundle.Version}, nil)
	session, err := client.Connect(doctorContext, &mcp.StreamableClientTransport{
		Endpoint: endpoint, HTTPClient: c.authenticatedHTTPClient(key), DisableStandaloneSSE: true, MaxRetries: -1,
	}, nil)
	if err != nil {
		return DoctorResult{}, doctorError(doctorContext, parsed.Host, err)
	}
	defer session.Close()

	initialize := session.InitializeResult()
	if initialize == nil {
		return DoctorResult{}, fmt.Errorf("MCP endpoint %s returned no initialize result", parsed.Host)
	}
	result := DoctorResult{
		AccessDomain:  d.AccessDomain,
		ServerName:    initialize.ServerInfo.Name,
		ServerVersion: initialize.ServerInfo.Version,
		Instructions:  initialize.Instructions,
		Tools:         []string{},
	}
	hasIdentity := false
	seenCursors := make(map[string]struct{})
	seenTools := make(map[string]struct{})
	params := &mcp.ListToolsParams{}
	for {
		page, listErr := session.ListTools(doctorContext, params)
		if listErr != nil {
			return DoctorResult{}, doctorError(doctorContext, parsed.Host, listErr)
		}
		for _, tool := range page.Tools {
			if _, exists := seenTools[tool.Name]; exists {
				continue
			}
			seenTools[tool.Name] = struct{}{}
			result.Tools = append(result.Tools, tool.Name)
			if tool.Name == d.Tool {
				hasIdentity = true
			}
		}
		if page.NextCursor == "" {
			break
		}
		if _, exists := seenCursors[page.NextCursor]; exists {
			return DoctorResult{}, doctorError(doctorContext, parsed.Host, errors.New("tools list returned a repeated cursor"))
		}
		seenCursors[page.NextCursor] = struct{}{}
		params.Cursor = page.NextCursor
	}
	sort.Strings(result.Tools)
	if !hasIdentity {
		result.Status = DoctorStatusAuthorizationRequired
		result.RequiredCapability = "supplier." + d.Name + ".identity.read"
		return result, nil
	}
	verification, callErr := session.CallTool(doctorContext, &mcp.CallToolParams{
		Name:      d.Tool,
		Arguments: map[string]any{},
	})
	if callErr != nil {
		return DoctorResult{}, doctorError(doctorContext, parsed.Host, callErr)
	}
	if verification == nil || verification.IsError {
		return DoctorResult{}, doctorError(doctorContext, parsed.Host, errors.New("supplier identity verification failed"))
	}
	output := verification.StructuredContent
	if output == nil {
		for _, content := range verification.Content {
			if text, ok := content.(*mcp.TextContent); ok {
				if output != nil {
					return DoctorResult{}, errors.New("ambiguous supplier identity output")
				}
				if err := json.Unmarshal([]byte(text.Text), &output); err != nil {
					return DoctorResult{}, errors.New("invalid supplier identity JSON")
				}
			}
		}
	}
	if err := validateIdentity(output, d); err != nil {
		return DoctorResult{}, err
	}
	result.Status = DoctorStatusVerified
	result.Verification = &VerificationResult{Tool: d.Tool, Output: output}
	return result, nil
}

func doctorError(ctx context.Context, host string, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "401") || strings.Contains(message, "unauthorized") {
		return fmt.Errorf("MCP endpoint %s: %w", host, ErrUnauthorized)
	}
	return fmt.Errorf("MCP endpoint %s: %w", host, ErrUnavailable)
}

func validateIdentity(output any, d domain.Config) error {
	data, err := json.Marshal(output)
	if err != nil {
		return errors.New("invalid supplier identity output")
	}
	var value struct {
		ClientCode string `json:"client_code"`
		Identity   struct {
			AccessDomain string `json:"access_domain"`
			SubjectType  string `json:"subject_type"`
			SubjectID    string `json:"subject_id"`
			TenantID     string `json:"tenant_id"`
			SupplierCode string `json:"supplier_code"`
			Location     string `json:"location"`
		} `json:"identity"`
	}
	if json.Unmarshal(data, &value) != nil {
		return errors.New("invalid supplier identity output")
	}
	id := value.Identity
	if id.AccessDomain != d.AccessDomain || id.SubjectType != "SUPPLIER_ACCOUNT" || strings.TrimSpace(id.SubjectID) == "" || strings.TrimSpace(id.TenantID) == "" || strings.TrimSpace(id.SupplierCode) == "" || strings.TrimSpace(value.ClientCode) == "" {
		return errors.New("supplier identity does not match the selected domain or is incomplete")
	}
	if d.Name == "fabric" && strings.TrimSpace(id.Location) == "" {
		return errors.New("fabric supplier identity is missing its location")
	}
	if d.Name == "garment" && id.Location != "" {
		return errors.New("garment supplier identity unexpectedly contains a location")
	}
	return nil
}
