# Supplier Agent Skill / 供应商 Agent Skill

中文 | [English](README.en.md)

供应商专用的独立 Agent Skill 与 Go 安装器，参考 ERP Agent Skill 的设计，支持 Codex、Claude Code 和 WorkBuddy。版本 `0.1.0`；当前交付源码和本地安装包，尚未发布远程仓库或 Release。

一份 `supplier-open-platform` Skill 支持两个独立 MCP 连接，可与 ERP 的 `erp-open` 同时安装：

| 领域参数 | MCP 地址后缀 | 配置名 | 身份工具 |
| --- | --- | --- | --- |
| `fabric` | `/fabric/mcp` | `supplier-fabric` | `get_fabric_supplier_identity` |
| `garment` | `/garment/mcp` | `supplier-garment` | `get_garment_supplier_identity` |

**首版服务只开放当前供应商身份查询。** 订单、报价、发货能力尚未开放；后续能力以当前 MCP `tools/list` 和服务端授权为准。

## 构建和本地安装

```bash
make install-hooks
make check
make build
./bin/supplier-skill version
```

Go 版本以 `go.mod` 为准。可直接运行 `./bin/supplier-skill`，也可将它复制到个人 PATH 目录。

`make dist` 生成 macOS arm64/amd64、Linux arm64/amd64、Windows amd64 五个平台包及 `dist/checksums.txt`；`bash scripts/verify-dist.sh` 校验 SHA-256。安装预编译包无需 Go，按[安装指南](docs/installation.md)校验和解压。

## 连接 Agent

先从管理员获取对应领域的 MCP 完整地址、供应商 Key 和身份查询授权。以下为 Bash 示例；地址是示例，需要替换，**不要把个人 Key 粘贴到聊天中**。

```bash
MCP_URL='https://supplier-open.example.com/garment/mcp'
IFS= read -r -s -p '供应商 API Key: ' SUPPLIER_INSTALL_KEY; printf '\n'
printf '%s\n' "$SUPPLIER_INSTALL_KEY" | ./bin/supplier-skill install codex --domain garment --url "$MCP_URL" --key-stdin
unset SUPPLIER_INSTALL_KEY
./bin/supplier-skill doctor codex --domain garment
```

Claude Code 将 `codex` 换成 `claude-code`；WorkBuddy 换成 `workbuddy`。面料使用 `--domain fabric` 和 `/fabric/mcp` 地址。两端分别安装、分别使用对应 Key；不得混用。

安装将凭证保存在用户级私有 MCP 配置，Unix 文件权限 `0600`，不经过客户端命令参数。重启客户端后正常打开 Codex、Claude Code 或 WorkBuddy 即可；不需要特殊启动器。安装不改变其他 MCP 配置。

`doctor` 和安装预检只调用选定领域的身份工具。`VERIFIED` 表示身份查询已实际验证；`INSTALLED_AUTHORIZATION_REQUIRED` 表示缺少身份能力，需要管理员授权。网络、认证、错域或畸形身份响应不会报告已验证。

## 使用与卸载

连接后可直接说：“查询当前成衣供应商身份，列出账号、租户和公司编码。”若两端都存在，应说明使用哪端。修改能力后重连 MCP 或新建会话。

```bash
./bin/supplier-skill uninstall codex --domain garment
```

只删除选定领域连接及安装记录；另一供应商领域和 ERP 保留。最后一个供应商连接卸载时才移除共享 Skill；手工修改过的 Skill 文件会保留并报告冲突。

## 文档

- [安装、升级、卸载](docs/installation.md)与[Agent 安装提示词](docs/agent-install-prompts.md)
- [Agent 使用](docs/agent-usage.md)、[构建和发布](docs/enterprise-release.md)
- [Skill 指令](SKILL.md)、[工作流](references/workflows.md)、[错误处理](references/errors.md)
- [本地验证记录](docs/verification.md)

维护者发布前需明确远程目标。安装脚本要求提供 `SUPPLIER_SKILL_RELEASE_BASE_URL`（PowerShell 为 `-BaseUrl`），发布命令要求 `REPO=owner/repository`；没有默认的假定公开下载地址。软件版权和分发条件见 [LICENSE](LICENSE)。
