# Installation, upgrade and removal

[中文](installation.md)

For Agent-assisted installation, copy the [installation prompt](agent-install-prompts.en.md). This page is the detailed manual reference.

Build with `make build`, then run `./bin/supplier-skill`. Alternatively verify the supplied `dist/checksums.txt` and extract the appropriate precompiled archive. Use `sha256sum -c checksums.txt`, macOS `shasum -a 256 -c checksums.txt`, or Windows `Get-FileHash -Algorithm SHA256`. Place the binary on your personal PATH.

[Release v0.1.1](https://github.com/xsy-tech/supplier-agent-skill/releases/tag/v0.1.1) provides precompiled packages in the private repository. Use authenticated `gh release download` as shown below and verify downloaded packages. For a directly accessible administrator mirror, set `SUPPLIER_SKILL_RELEASE_BASE_URL` before `bash scripts/install.sh`, or pass `-BaseUrl` to `scripts/install.ps1`. Version defaults to `0.1.1`; `SUPPLIER_SKILL_INSTALL_DIR` overrides the default personal `.local/bin`. Mirror scripts do not carry GitHub account tokens; use authenticated `gh` for private GitHub assets.

Use the local hidden-input example below. The connection commands require `--domain fabric|garment`, a matching `/fabric/mcp` or `/garment/mcp` URL and `--key-stdin`. Select `codex`, `claude-code` or `workbuddy`. URLs must use HTTPS except loopback development, and cannot embed credentials, query strings or fragments. Literal Key command arguments are rejected.

Download with an authorized GitHub login:

```bash
gh release download v0.1.1 --repo xsy-tech/supplier-agent-skill --dir ./supplier-release
cd supplier-release
shasum -a 256 -c checksums.txt
```

This downloads all assets. If downloading only your platform archive and the manifest, verify that archive's checksum individually.

With the CLI on your PATH, replace the example URL and enter the Key locally:

```bash
MCP_URL='https://supplier-open.example.com/garment/mcp'
IFS= read -r -s -p 'Supplier API Key: ' SUPPLIER_INSTALL_KEY; printf '\n'
printf '%s\n' "$SUPPLIER_INSTALL_KEY" | supplier-skill install codex --domain garment --url "$MCP_URL" --key-stdin
unset SUPPLIER_INSTALL_KEY
supplier-skill doctor codex --domain garment
```

For Windows PowerShell, the [Chinese installation page](installation.md) includes a language-independent SecureString/BSTR example, with explicit buffer cleanup. Never put the Key in chat or command history.

User configurations are `$CODEX_HOME/config.toml` (default `~/.codex/config.toml`), `~/.claude.json` (under `CLAUDE_CONFIG_DIR` when set), and `~/.workbuddy/mcp.json`. The Skill is installed under each client's user `skills/supplier-open-platform` directory. Unix credential-bearing files are `0600`; do not share or commit them. This targets Claude Code, not the separate Claude desktop chat configuration.

Restart the client after installation. `supplier-fabric` and `supplier-garment` remain independent and preserve ERP's `erp-open`.

```bash
supplier-skill doctor codex --domain garment
supplier-skill uninstall codex --domain garment
```

These read the saved domain-specific endpoint; `--url` can be explicit. WorkBuddy also accepts `doctor ... --key-stdin`. `launch codex|claude-code --domain garment -- <client arguments>` is a compatibility command; persistent installations normally start clients directly.

Replace the CLI to upgrade, then rerun install. Inspect manually edited Skill files before using `--force`. Removal preserves other domains and other MCP servers, and only removes the shared Skill after the final supplier connection. Edited or unmanaged Skill files are protected.

Optional WorkBuddy export: `supplier-skill install workbuddy --domain fabric --url https://supplier-open.example.com/fabric/mcp --output ./supplier-fabric-connector`. This contains no real Key; import it and enter the token locally in WorkBuddy. Do not combine export with `--key-stdin`.

### First trust in WorkBuddy

Restart after installation and inspect custom connectors under connector management. Approve each supplier server when prompted. Each currently exposes only its domain identity tool. In a new session, load `supplier-open-platform`, call that native tool and verify domain, account, tenant, company and Client; fabric also requires location. Installer `VERIFIED` alone does not prove native desktop execution. Reconnect after grants or credential changes. See [native acceptance](verification-workbuddy-2026-11-09.en.md). Never expose Keys in the UI, logs or screenshots.
