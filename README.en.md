# Supplier Agent Skill

[中文](README.md) | English

Connect Codex, Claude Code or WorkBuddy to the Supplier Open Platform. Currently supports fabric and garment supplier identity reads.

## Copy this prompt to your Agent

Send the following to your current Agent. It handles downloading, extraction and configuration:

```text
Install the Supplier Open Platform Skill in my current client and verify the connection using this guide:
https://github.com/xsy-tech/supplier-agent-skill/blob/main/docs/agent-install-prompts.en.md
If my supplier domain or MCP URL is missing, ask me for both together.
Let me enter the API Key securely on my machine, never in chat.
```

Have your administrator's **supplier domain, MCP URL and API Key** ready. Share the domain and URL with the Agent; enter the Key locally. This repository is private. If the guide is inaccessible, sign in to an authorized GitHub account or request an offline package and guide from your administrator.

The Agent selects and verifies your platform package, installs the connection and checks its identity. It explains any required restart or first server trust; use the Skill after native connection verification succeeds.

## After installation

Ask: “Use the Supplier Open Platform Skill to show my current supplier account, tenant and company code.” Specify fabric or garment if both are connected. Orders, quotations and shipments are not yet available.

## Further reading

- [Agent installation guide](docs/agent-install-prompts.en.md)
- [Manual installation, upgrades and removal](docs/installation.en.md)
- [Usage](docs/agent-usage.en.md) · [Errors](references/errors.en.md)
- [Release v0.1.1](https://github.com/xsy-tech/supplier-agent-skill/releases/tag/v0.1.1): five platform packages, SHA-256 and offline guides; repository access required.
- [Native WorkBuddy acceptance](docs/verification-workbuddy-2026-11-09.en.md): v0.1.0 passed 5/5 scoped cases; v0.1.1 documentation and package checks are recorded separately.
- [Build and release](docs/enterprise-release.en.md) · [Skill](SKILL.md) · [Verification](docs/verification.md) · [LICENSE](LICENSE)
