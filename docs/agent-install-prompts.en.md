# Install the supplier Skill with an Agent

[中文](agent-install-prompts.md)

## Prompt for users

Send this to Codex, Claude Code or WorkBuddy without including an API Key:

```text
Install the Supplier Open Platform Skill in my current client and verify the connection using this guide:
https://github.com/xsy-tech/supplier-agent-skill/blob/main/docs/agent-install-prompts.en.md
If my supplier domain or MCP URL is missing, ask me for both together.
Let me enter the API Key securely on my machine, never in chat.
```

Optionally append your client, supplier domain and administrator-provided complete MCP URL. No manual configuration editing is needed.

## Instructions for the installing Agent

Follow this guide directly; do not ask the user to copy the prompt again.

1. Identify the current client (`codex`, `claude-code` or `workbuddy`). Ask once for missing or ambiguous client, supplier type and complete MCP URL. Fabric uses `--domain fabric`, `/fabric/mcp`, `supplier-fabric`, and `get_fabric_supplier_identity`; garment uses `--domain garment`, `/garment/mcp`, `supplier-garment`, and `get_garment_supplier_identity`. Each domain needs its own Key. Preserve existing connections, including ERP's `erp-open`. Require HTTPS except loopback development; reject URL credentials, query strings and fragments.
2. Prefer the published **v0.1.1** binary from the private `xsy-tech/supplier-agent-skill` GitHub repository. Detect the local platform: macOS/Linux arm64/amd64 or Windows amd64. Check existing authenticated access and download the matching `supplier-agent-skill_0.1.1_<platform>.tar.gz` or Windows `.zip`, plus `checksums.txt`, using `gh release download`. With `--pattern`, download just those two files and verify the selected archive's checksum; do not run a full-manifest check requiring all five archives. An administrator-provided offline package and checksum manifest also work. If access is unavailable, guide local sign-in or request the offline package, never a Token in chat. Ordinary installation does not require Go or a source build. Extract only after verification, use a user-writable directory or absolute binary path, and confirm `version`.
3. Obtain the Key through a user-operable local hidden-input terminal or protected input channel, then run `supplier-skill install <client> --domain <fabric|garment> --url <complete-MCP-URL> --key-stdin` with placeholders replaced. Never place credentials in chat, arguments, script constants or tool output. See [manual installation](installation.en.md) for local Bash/PowerShell input. If secure user input is unavailable, provide a local command and report “waiting for local Key entry”; do not claim completion. Let the installer manage private configuration (Unix `0600`); inspect conflicts before considering `--force`.
4. Run `supplier-skill doctor <client> --domain <fabric|garment>`. `VERIFIED` confirms its actual identity call; `INSTALLED_AUTHORIZATION_REQUIRED` means the administrator must grant that domain's identity capability. Report all other errors accurately. Restart/reconnect the client and load `supplier-open-platform` in a native session. Discover and call the matching identity tool with `{}`; verify domain, main account, tenant, company and Client, plus location for fabric. Do not infer a garment location. WorkBuddy may require first trust for each server. If user restart/trust is pending, report “configured; awaiting restart/trust and native verification.” Doctor success alone does not prove native execution.

Finish with the selected client/domain, diagnostic and native verification results, and the next required step, if any. Do not display full credentials or configuration files. The current service supports identity reads only; order, quotation and shipment tools are not yet available.
