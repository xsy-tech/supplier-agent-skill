package domain

import "testing"

func TestSupplierEndpointBoundaries(t *testing.T) {
	for _, raw := range []string{"https://open.example/fabric/mcp", "https://open.example/garment/mcp/", "http://127.0.0.1:18085/garment/mcp", "http://[::1]:8080/fabric/mcp", "http://localhost/fabric/mcp"} {
		if _, err := FromEndpoint(raw); err != nil {
			t.Errorf("valid endpoint %q: %v", raw, err)
		}
	}
	for _, raw := range []string{"https://open.example/mcp", "http://open.example/garment/mcp", "https://user:secret@open.example/garment/mcp", "https://open.example/garment/mcp?tenant=x", "https://open.example/garment/mcp?", "https://open.example/garment/mcp#fragment", "https://open.example/%67arment/mcp", "https://open.example/fabric/mcp/../garment/mcp", "file:///garment/mcp", "https:///garment/mcp"} {
		if _, err := FromEndpoint(raw); err == nil {
			t.Errorf("accepted invalid endpoint %q", raw)
		}
	}
	for _, name := range []string{"", "erp", "GARMENT_SUPPLIER"} {
		if _, err := Parse(name); err == nil {
			t.Errorf("accepted domain %q", name)
		}
	}
}
