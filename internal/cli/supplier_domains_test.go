package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-saas/supplier-agent-skill/internal/mcpclient"
	"github.com/go-saas/supplier-agent-skill/internal/profile"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSupplierDomainsCoexistAndUninstallIndependently(t *testing.T) {
	for _, agent := range []string{"codex", "claude-code", "workbuddy"} {
		t.Run(agent, func(t *testing.T) {
			app, _, _, _, _, _ := newTestApp(t)
			app.Profiles = profile.NewFileStore(filepath.Join(app.HomeDir, "profiles.json"))
			configPath, err := app.clientConfig().Path(agent)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
				t.Fatal(err)
			}
			seed := `{"mcpServers":{"erp-open":{"url":"https://erp.example/mcp","headers":{"Authorization":"Bearer erp-preserve"}}},"custom":"preserve"}`
			if agent == "codex" {
				seed = "custom = \"preserve\"\n[mcp_servers.erp-open]\nurl = \"https://erp.example/mcp\"\n[mcp_servers.erp-open.http_headers]\nAuthorization = \"Bearer erp-preserve\"\n"
			}
			if err := os.WriteFile(configPath, []byte(seed), 0600); err != nil {
				t.Fatal(err)
			}

			for _, d := range []string{"fabric", "garment"} {
				app.In = strings.NewReader("private-" + d + "\n")
				if err := app.Run(t.Context(), []string{"install", agent, "--domain", d, "--url", "https://supplier.example.com/" + d + "/mcp", "--key-stdin"}); err != nil {
					t.Fatal(err)
				}
			}
			path, err := app.clientConfig().Path(agent)
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"supplier-fabric", "supplier-garment", "private-fabric", "private-garment"} {
				if !strings.Contains(string(data), name) {
					t.Fatalf("missing %s", name)
				}
			}
			if err := app.Run(t.Context(), []string{"uninstall", agent, "--domain", "fabric"}); err != nil {
				t.Fatal(err)
			}
			data, _ = os.ReadFile(path)
			if strings.Contains(string(data), "supplier-fabric") || !strings.Contains(string(data), "supplier-garment") {
				t.Fatal("uninstall crossed domain boundary")
			}
			skillRoot := filepath.Join(app.HomeDir, ".codex", "skills", "supplier-open-platform")
			if agent == "claude-code" {
				skillRoot = filepath.Join(app.HomeDir, ".claude", "skills", "supplier-open-platform")
			}
			if agent == "workbuddy" {
				skillRoot = filepath.Join(app.HomeDir, ".workbuddy", "skills", "supplier-open-platform")
			}
			if _, err := os.Stat(filepath.Join(skillRoot, "SKILL.md")); err != nil {
				t.Fatal("shared skill removed too early", err)
			}
			if err := app.Run(t.Context(), []string{"doctor", agent, "--domain", "garment"}); err != nil {
				t.Fatal(err)
			}
			if err := app.Run(t.Context(), []string{"uninstall", agent, "--domain", "garment"}); err != nil {
				t.Fatal(err)
			}
			data, err = os.ReadFile(configPath)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "erp-preserve") || !strings.Contains(string(data), "erp-open") || !strings.Contains(string(data), "custom") {
				t.Fatal("ERP or unrelated configuration removed")
			}
			if _, err := os.Stat(skillRoot); !os.IsNotExist(err) {
				t.Fatal("shared skill remains after last uninstall")
			}
		})
	}
}

func TestDomainIsExplicitAndMatchesEndpoint(t *testing.T) {
	for _, args := range [][]string{
		{"install", "codex", "--url", "https://supplier.example.com/garment/mcp", "--key-stdin"},
		{"install", "codex", "--domain", "fabric", "--url", "https://supplier.example.com/garment/mcp", "--key-stdin"},
		{"install", "codex", "--domain", "garment", "--url", "https://supplier.example.com/mcp", "--key-stdin"},
	} {
		app, _, _, _, _, _ := newTestApp(t)
		app.In = strings.NewReader("secret")
		if err := app.Run(t.Context(), args); err == nil {
			t.Fatal("accepted missing or conflicting domain")
		}
	}
}

func TestSupplierCLIUsesOwnCredentialsOverHTTP(t *testing.T) {
	for _, agent := range []string{"codex", "claude-code", "workbuddy"} {
		t.Run(agent, func(t *testing.T) {
			handlers := map[string]http.Handler{}
			for _, name := range []string{"fabric", "garment"} {
				d := name
				access := "GARMENT_SUPPLIER"
				if d == "fabric" {
					access = "FABRIC_SUPPLIER"
				}
				identity := map[string]any{"access_domain": access, "subject_type": "SUPPLIER_ACCOUNT", "subject_id": "91", "tenant_id": "T01", "supplier_code": "S01"}
				if d == "fabric" {
					identity["location"] = "CN"
				}
				server := mcp.NewServer(&mcp.Implementation{Name: "supplier", Version: "test"}, nil)
				server.AddTool(&mcp.Tool{Name: "get_" + d + "_supplier_identity", InputSchema: json.RawMessage(`{"type":"object"}`)}, func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					return &mcp.CallToolResult{StructuredContent: map[string]any{"client_code": "C01", "identity": identity}}, nil
				})
				handlers["/"+d+"/mcp"] = mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				d := strings.Split(r.URL.Path, "/")[1]
				if r.Header.Get("Authorization") != "Bearer private-"+d {
					http.Error(w, "wrong credential", 401)
					return
				}
				h, ok := handlers[r.URL.Path]
				if !ok {
					http.NotFound(w, r)
					return
				}
				h.ServeHTTP(w, r)
			}))
			defer ts.Close()
			app, out, _, _, _, _ := newTestApp(t)
			app.Doctor = mcpclient.New()
			for _, d := range []string{"fabric", "garment"} {
				app.In = strings.NewReader("private-" + d)
				if err := app.Run(t.Context(), []string{"install", agent, "--domain", d, "--url", ts.URL + "/" + d + "/mcp", "--key-stdin"}); err != nil {
					t.Fatal(err)
				}
			}
			for _, d := range []string{"fabric", "garment"} {
				out.Reset()
				if err := app.Run(t.Context(), []string{"doctor", agent, "--domain", d}); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(out.String(), "VERIFIED") || strings.Contains(out.String(), "private-") {
					t.Fatal("diagnostic failed or disclosed credential")
				}
			}
		})
	}
}

func TestInstallRejectsInvalidClientConfigBeforeInstallingSkill(t *testing.T) {
	app, _, _, _, _, _ := newTestApp(t)
	path, err := app.clientConfig().Path("codex")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("broken = ["), 0600); err != nil {
		t.Fatal(err)
	}
	app.In = strings.NewReader("private-key")
	if err := app.Run(t.Context(), []string{"install", "codex", "--domain", "garment", "--url", "https://supplier.example/garment/mcp", "--key-stdin"}); err == nil {
		t.Fatal("invalid existing configuration accepted")
	}
	if _, err := os.Stat(filepath.Join(app.CodexHome, "skills", "supplier-open-platform")); !os.IsNotExist(err) {
		t.Fatal("failed configuration preflight left a partial skill installation")
	}
}
