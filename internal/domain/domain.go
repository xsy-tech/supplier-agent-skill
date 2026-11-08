// Package domain defines the two independent supplier access domains.
package domain

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

type Config struct{ Name, AccessDomain, ServerName, Tool, Path string }

func Parse(name string) (Config, error) {
	switch name {
	case "fabric":
		return Config{name, "FABRIC_SUPPLIER", "supplier-fabric", "get_fabric_supplier_identity", "/fabric/mcp"}, nil
	case "garment":
		return Config{name, "GARMENT_SUPPLIER", "supplier-garment", "get_garment_supplier_identity", "/garment/mcp"}, nil
	default:
		return Config{}, errors.New("--domain must be fabric or garment")
	}
}

func FromEndpoint(raw string) (Config, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Hostname() == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || parsed.RawPath != "" {
		return Config{}, errors.New("invalid supplier MCP URL")
	}
	if parsed.Scheme != "https" {
		ip := net.ParseIP(parsed.Hostname())
		if parsed.Scheme != "http" || (parsed.Hostname() != "localhost" && (ip == nil || !ip.IsLoopback())) {
			return Config{}, errors.New("supplier MCP URL requires HTTPS outside loopback")
		}
	}
	for _, name := range []string{"fabric", "garment"} {
		d, _ := Parse(name)
		if strings.TrimSuffix(parsed.Path, "/") == d.Path {
			return d, nil
		}
	}
	return Config{}, errors.New("supplier MCP URL must end in /fabric/mcp or /garment/mcp")
}

func (d Config) ValidateEndpoint(raw string) error {
	actual, err := FromEndpoint(raw)
	if err != nil {
		return err
	}
	if actual.Name != d.Name {
		return errors.New("MCP URL does not match --domain")
	}
	return nil
}
