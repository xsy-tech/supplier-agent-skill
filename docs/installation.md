# 安装、升级与卸载

[English](installation.en.md)

## 安装 CLI

当前提供本地源码和 `dist/` 安装包，尚未发布远程 Release。从源码：`make build`，然后使用 `./bin/supplier-skill`。从预编译包：先用 `sha256sum -c checksums.txt` 或 macOS `shasum -a 256 -c checksums.txt` 校验管理员提供的安装包与清单，再解压适合本机的包，把 `supplier-skill` 放到个人 PATH。Windows 解压 ZIP，将 `supplier-skill.exe` 放到个人 PATH；可用 `Get-FileHash -Algorithm SHA256` 比较清单。

管理员部署下载目录后，Bash 可设置 `SUPPLIER_SKILL_RELEASE_BASE_URL` 并运行 `bash scripts/install.sh`；Windows 运行 `./scripts/install.ps1 -BaseUrl <管理员下载目录>`。脚本选择平台、下载并校验 SHA-256，再安装至用户 `.local/bin`。可配置 `SUPPLIER_SKILL_INSTALL_DIR`，版本默认 `0.1.0`。无下载地址时直接报错，不能使用推测的远程仓库。

## 连接三个客户端

[README](../README.md) 提供 Bash 隐藏输入示例。参数 `--domain fabric|garment` 必填，地址必须匹配 `/fabric/mcp` 或 `/garment/mcp`；仅本机回环允许 HTTP，其他地址要求 HTTPS。不允许 URL 中携带用户名、密码、查询串和片段，不接受命令行明文 Key。

```text
supplier-skill install codex --domain garment --url https://supplier-open.example.com/garment/mcp --key-stdin
supplier-skill install claude-code --domain garment --url https://supplier-open.example.com/garment/mcp --key-stdin
supplier-skill install workbuddy --domain garment --url https://supplier-open.example.com/garment/mcp --key-stdin
```

标准输入必须来自本地隐藏输入。Windows PowerShell 示例：

```powershell
$secret = Read-Host 'Supplier API Key' -AsSecureString
$ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secret)
try {
  $plain = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr)
  $plain | .\bin\supplier-skill.exe install codex --domain garment --url 'https://supplier-open.example.com/garment/mcp' --key-stdin
} finally {
  [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr)
  Remove-Variable plain, secret -ErrorAction SilentlyContinue
}
```

将客户端替换为 `claude-code` 或 `workbuddy`，面料替换领域与地址。安装完成后重启客户端。普通 Claude 桌面聊天与 Claude Code 不同，本安装器针对 Claude Code 本地配置。

| 客户端 | MCP 配置 | Skill 目录 |
| --- | --- | --- |
| Codex | `$CODEX_HOME/config.toml`，默认 `~/.codex/config.toml` | 对应目录下 `skills/supplier-open-platform` |
| Claude Code | `~/.claude.json`，设置 `CLAUDE_CONFIG_DIR` 时使用其下 `.claude.json` | 对应 Claude 配置目录下 `skills/supplier-open-platform` |
| WorkBuddy | `~/.workbuddy/mcp.json` | `~/.workbuddy/skills/supplier-open-platform` |

配置包含个人 Key，Unix 权限为 `0600`。不要提交或分享这些文件。供应商连接名称分别为 `supplier-fabric`、`supplier-garment`；不会覆盖 ERP `erp-open`。每个客户端和领域有独立安装记录。缺少身份工具时安装会明确返回 `INSTALLED_AUTHORIZATION_REQUIRED`，不宣称身份验证成功。

## 诊断、升级、卸载

```bash
supplier-skill doctor codex --domain garment
supplier-skill uninstall codex --domain garment
```

上述命令默认读取安装记录；也可显式 `--url`。WorkBuddy 支持 `doctor ... --key-stdin` 做独立验证。`launch codex|claude-code --domain garment -- <客户端参数>` 为兼容入口；持久化安装后通常直接打开客户端。

升级替换 CLI 后重跑 install；安装文件若经过手工修改，会要求 `--force`，应先检查差异。卸载只移除选定领域及其私有配置，最后一个供应商连接卸载时才清理共享 Skill。无法识别或修改过的 Skill 文件不会被静默删除。

## 可选 WorkBuddy Connector

```bash
supplier-skill install workbuddy --domain fabric --url https://supplier-open.example.com/fabric/mcp --output ./supplier-fabric-connector
```

此方式导出无真实 Key 的包；在 WorkBuddy 导入并通过 Token 界面填写。不可同时使用 `--key-stdin`。两领域使用不同 Connector 标识。默认仍推荐直接安装。
