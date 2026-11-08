# Supplier Agent Skill

[中文](README.md) | English

Version `0.1.0` provides an independent supplier Skill and Go installer for Codex, Claude Code and WorkBuddy, based on the ERP Agent Skill design. Source and local packages are delivered; no remote repository or Release has been published.

Use `--domain fabric` with `/fabric/mcp` and `supplier-fabric`, or `--domain garment` with `/garment/mcp` and `supplier-garment`. Both connections share `supplier-open-platform` and coexist with ERP's `erp-open`.

Initially the service exposes only `get_fabric_supplier_identity` or `get_garment_supplier_identity`. Order, quotation and shipment tools are not available yet. Runtime discovery and server grants determine future availability.

```bash
make install-hooks
make check
make build
./bin/supplier-skill version
```

Use the Go version specified in `go.mod`. Run the built binary directly or place it on your PATH. `make dist` creates five platform packages (macOS/Linux amd64/arm64, Windows amd64) and SHA-256 checksums. Precompiled installation does not need Go.

Obtain a domain-specific MCP URL, Key and identity grant from your administrator. In Bash, replace the example URL and enter the Key locally, never in chat:

```bash
MCP_URL='https://supplier-open.example.com/garment/mcp'
IFS= read -r -s -p 'Supplier API Key: ' SUPPLIER_INSTALL_KEY; printf '\n'
printf '%s\n' "$SUPPLIER_INSTALL_KEY" | ./bin/supplier-skill install codex --domain garment --url "$MCP_URL" --key-stdin
unset SUPPLIER_INSTALL_KEY
./bin/supplier-skill doctor codex --domain garment
```

Substitute `claude-code` or `workbuddy` for the client. Fabric uses its own URL, Key and `--domain fabric`. Restart the client after installation; authentication is saved in private user-level configuration (Unix `0600`). Other connections are preserved.

`VERIFIED` means the identity call and its output were validated. `INSTALLED_AUTHORIZATION_REQUIRED` means the identity grant is missing. Diagnostics never call writes or depend on ERP's platform-info tool.

```bash
./bin/supplier-skill uninstall codex --domain garment
```

Uninstall removes the selected connection only. The shared Skill is removed after the final supplier connection, unless its files have been manually changed.

See [installation](docs/installation.en.md), [installation prompts](docs/agent-install-prompts.en.md), [usage](docs/agent-usage.en.md), [release](docs/enterprise-release.en.md), [Skill](SKILL.md) and [verification](docs/verification.md). Download scripts require an administrator-provided release directory; no public download location is assumed. Distribution terms are in [LICENSE](LICENSE).
