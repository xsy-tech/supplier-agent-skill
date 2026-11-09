# Supplier Agent Skill

[中文](README.md) | English

Version `0.1.1` provides an independent supplier Skill and Go installer for Codex, Claude Code and WorkBuddy, based on the ERP Agent Skill design. [Source](https://github.com/xsy-tech/supplier-agent-skill) and [v0.1.1 release packages](https://github.com/xsy-tech/supplier-agent-skill/releases/tag/v0.1.1) are hosted in a private repository and require access.

Use `--domain fabric` with `/fabric/mcp` and `supplier-fabric`, or `--domain garment` with `/garment/mcp` and `supplier-garment`. Both connections share `supplier-open-platform` and coexist with ERP's `erp-open`.

Initially the service exposes only `get_fabric_supplier_identity` or `get_garment_supplier_identity`. Order, quotation and shipment tools are not available yet. Runtime discovery and server grants determine future availability.

WorkBuddy 5.3.11 passed **5/5 scoped native cases** using v0.1.0: Skill loading, tool discovery and both supplier identity calls. See [acceptance evidence](docs/verification-workbuddy-2026-11-09.en.md). Version v0.1.1 packages this evidence and updated guides; native and package verification versions are distinguished.

Download all platform packages and checksums with an authorized GitHub account:

```bash
gh auth login
gh release download v0.1.1 --repo xsy-tech/supplier-agent-skill --dir ./supplier-release
cd supplier-release
shasum -a 256 -c checksums.txt  # or sha256sum -c checksums.txt
```

Extract the archive for your platform and put `supplier-skill` on your personal PATH. Windows uses ZIP and `Get-FileHash -Algorithm SHA256`. Private Release downloads require authentication; ask the administrator for access or an approved mirror.

For source builds:

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

If installed from Release onto PATH, replace `./bin/supplier-skill` with `supplier-skill`. Substitute `claude-code` or `workbuddy` for the client. Fabric uses its own URL, Key and `--domain fabric`. Restart the client after installation; authentication is saved in private user-level configuration (Unix `0600`). Other connections are preserved.

`VERIFIED` means the identity call and its output were validated. `INSTALLED_AUTHORIZATION_REQUIRED` means the identity grant is missing. Diagnostics never call writes or depend on ERP's platform-info tool.

```bash
./bin/supplier-skill uninstall codex --domain garment
```

Uninstall removes the selected connection only. The shared Skill is removed after the final supplier connection, unless its files have been manually changed.

See [installation](docs/installation.en.md), [installation prompts](docs/agent-install-prompts.en.md), [usage](docs/agent-usage.en.md), [release](docs/enterprise-release.en.md), [Skill](SKILL.md) and [verification](docs/verification.md). The release target defaults to `xsy-tech/supplier-agent-skill`. Prefer authenticated `gh` downloads for this private repository. Download scripts accept an administrator-provided directly accessible mirror directory. Distribution terms are in [LICENSE](LICENSE).
