package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-saas/supplier-agent-skill/internal/mcpclient"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestWorkBuddyInstallAndSavedCredentialDoctorOverHTTP(t *testing.T) {
	const key = "workbuddy-http-integration-test-key"
	var readCalls, writeCalls atomic.Int32
	server := mcp.NewServer(&mcp.Implementation{Name: "supplier-open-platform", Version: "test"}, nil)
	for _, name := range []string{"get_garment_supplier_identity", "create_order"} {
		server.AddTool(&mcp.Tool{Name: name, Description: "integration test", InputSchema: json.RawMessage(`{"type":"object"}`)},
			func(_ context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				if request.Params.Name == "get_garment_supplier_identity" {
					readCalls.Add(1)
				} else {
					writeCalls.Add(1)
				}
				return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: `{"identity":{"access_domain":"GARMENT_SUPPLIER","subject_type":"SUPPLIER_ACCOUNT","subject_id":"91","tenant_id":"T01","supplier_code":"S01"},"client_code":"C01"}`}}}, nil
			})
	}
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})
	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+key {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer endpoint.Close()
	app, stdout, _, _, _, _ := newTestApp(t)
	app.Doctor = mcpclient.New()
	app.In = strings.NewReader(key + "\n")
	if err := app.runGarmentTest(t.Context(), []string{"install", "workbuddy", "--url", endpoint.URL + "/garment/mcp", "--key-stdin"}); err != nil {
		t.Fatal(err)
	}
	app.In = strings.NewReader("")
	stdout.Reset()
	if err := app.runGarmentTest(t.Context(), []string{"doctor", "workbuddy", "--url", endpoint.URL + "/garment/mcp"}); err != nil {
		t.Fatal(err)
	}
	if readCalls.Load() != 2 || writeCalls.Load() != 0 || !strings.Contains(stdout.String(), "VERIFIED") || strings.Contains(stdout.String(), key) {
		t.Fatal("saved credential verification failed, invoked a write, or exposed a credential")
	}
}
