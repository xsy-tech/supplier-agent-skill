package mcpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-saas/supplier-agent-skill/internal/domain"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSupplierDoctorValidatesIdentity(t *testing.T) {
	for _, d := range []string{"fabric", "garment"} {
		t.Run(d, func(t *testing.T) {
			server := mcp.NewServer(&mcp.Implementation{Name: "supplier", Version: "1"}, nil)
			domain := "GARMENT_SUPPLIER"
			if d == "fabric" {
				domain = "FABRIC_SUPPLIER"
			}
			identity := map[string]any{"access_domain": domain, "subject_type": "SUPPLIER_ACCOUNT", "subject_id": "91", "supplier_code": "S01", "tenant_id": "T01"}
			if d == "fabric" {
				identity["location"] = "CN"
			}
			mcp.AddTool(server, &mcp.Tool{Name: "get_" + d + "_supplier_identity"}, func(context.Context, *mcp.CallToolRequest, noInput) (*mcp.CallToolResult, map[string]any, error) {
				return nil, map[string]any{"identity": identity, "client_code": "C01"}, nil
			})
			handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})
			ts := httptest.NewServer(handler)
			defer ts.Close()
			result, err := New().Doctor(t.Context(), ts.URL+"/"+d+"/mcp", "secret")
			if err != nil || result.Status != DoctorStatusVerified {
				t.Fatalf("valid supplier identity rejected: %v %s", err, result.Status)
			}
			identity["access_domain"] = "ERP"
			result, err = New().Doctor(t.Context(), ts.URL+"/"+d+"/mcp", "secret")
			if err == nil || result.Status == DoctorStatusVerified {
				t.Fatal("wrong domain accepted")
			}
		})
	}
}

func TestRejectIncompleteSupplierIdentity(t *testing.T) {
	for _, name := range []string{"fabric", "garment"} {
		t.Run(name, func(t *testing.T) {
			d, _ := domain.Parse(name)
			valid := map[string]any{"client_code": "C01", "identity": map[string]any{"access_domain": d.AccessDomain, "subject_type": "SUPPLIER_ACCOUNT", "subject_id": "91", "tenant_id": "T01", "supplier_code": "S01"}}
			if name == "fabric" {
				valid["identity"].(map[string]any)["location"] = "CN"
			}
			for _, field := range []string{"client_code", "access_domain", "subject_type", "subject_id", "tenant_id", "supplier_code", "location"} {
				data, _ := json.Marshal(valid)
				var output map[string]any
				json.Unmarshal(data, &output)
				if field == "client_code" {
					output[field] = ""
				} else if field == "location" && name == "garment" {
					output["identity"].(map[string]any)[field] = "unexpected"
				} else {
					output["identity"].(map[string]any)[field] = ""
				}
				if err := validateIdentity(output, d); err == nil {
					t.Fatalf("accepted invalid %s", field)
				}
			}
			for _, output := range []any{nil, "secret", []any{}, map[string]any{"identity": "secret"}} {
				if err := validateIdentity(output, d); err == nil || strings.Contains(err.Error(), "secret") {
					t.Fatal("accepted malformed identity or disclosed its body")
				}
			}
		})
	}
}

func TestDoctorDoesNotFollowRedirectsWithCredential(t *testing.T) {
	var received bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { received = true; http.Error(w, "unexpected", 500) }))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/garment/mcp", http.StatusTemporaryRedirect)
	}))
	defer source.Close()
	result, err := New().Doctor(t.Context(), source.URL+"/garment/mcp", "private-key")
	if err == nil || result.Status == DoctorStatusVerified || received {
		t.Fatal("redirect accepted or credential forwarded")
	}
}
