# 用 Agent 安装供应商 Skill

[English](agent-install-prompts.en.md)

## 用户直接复制这段

发给 Codex、Claude Code 或 WorkBuddy；个人 API Key 不放进提示词：

```text
请按下面的安装指南，帮我把供应商开放平台 Skill 安装到当前客户端，并完成连接验证：
https://github.com/xsy-tech/supplier-agent-skill/blob/main/docs/agent-install-prompts.md
缺少供应商类型或 MCP 地址时，请一次向我确认。
API Key 请让我在本机安全输入，不要让我粘贴到聊天中。
```

已经知道信息时，可在末尾补充“客户端：WorkBuddy；供应商类型：成衣；MCP 地址：管理员提供的完整地址”。面料选 `fabric`，成衣选 `garment`；不用自己填写配置文件。

## 给执行安装的 Agent

目标是使用官方安装器持久化配置 `supplier-open-platform` Skill 与对应领域的 MCP 连接，并诚实报告验证进度。下面的内容就是执行指南；无需再让用户复制上面的提示词。

### 1. 确认安装目标

识别用户当前客户端：Codex、Claude Code 或 WorkBuddy，分别对应 `codex`、`claude-code`、`workbuddy`。无法识别时再询问，不默认向三个客户端都安装。只在缺失或不明确时，一次询问客户端、供应商类型及 MCP 完整地址；不要索要聊天中的 Key。

| 供应商类型 | 领域参数 | MCP 地址后缀 | 安装后的连接名 | 身份工具 |
| --- | --- | --- | --- | --- |
| 面料 | `fabric` | `/fabric/mcp` | `supplier-fabric` | `get_fabric_supplier_identity` |
| 成衣 | `garment` | `/garment/mcp` | `supplier-garment` | `get_garment_supplier_identity` |

地址使用管理员提供的实际值。非本机回环地址要求 HTTPS；不接受含用户名、密码、Query 或片段的 URL。两个领域各用自己的 Key。用户要装两个领域时分别安装，保留 ERP `erp-open` 与其他现有连接。

### 2. 获取并校验本机安装包

仓库：`https://github.com/xsy-tech/supplier-agent-skill`；当前已发布版本：`v0.1.1`。优先用预编译包，普通用户无需安装 Go 或从源码编译。

根据本机系统/架构选择 `supplier-agent-skill_0.1.1_<平台>.tar.gz` 或 `.zip`：macOS 对应 `darwin_arm64` / `darwin_amd64`，Linux 对应 `linux_arm64` / `linux_amd64`，Windows 为 `windows_amd64.zip`。不支持的本机平台应明确说明。

这是私有仓库，先检查现有 GitHub 登录是否具有访问权限，再用 `gh release download v0.1.1 --repo xsy-tech/supplier-agent-skill` 下载目标包与 `checksums.txt`。可以用 `--pattern` 仅下载这两个文件；只下载一个平台时，核对清单中该包的 SHA-256，不直接运行要求所有五个平台文件都存在的全清单校验。也可校验管理员提供的离线包与清单。无法访问时引导用户在本机完成授权或索取离线包，不索要 GitHub Token，不反复匿名下载私有资源。

校验通过后解压，将 `supplier-skill` / `supplier-skill.exe` 放入用户可写目录；可先用绝对路径运行，执行 `version` 确认版本。避免为安装 CLI 索取管理员权限。手动操作参考[安装指南](installation.md)。

### 3. 安全输入 Key 并安装

使用用户可操作的本机隐藏输入或受保护输入渠道，通过标准输入运行：

```text
supplier-skill install <客户端> --domain <fabric或garment> --url <管理员提供的完整MCP地址> --key-stdin
```

尖括号内容为待替换参数。禁止把 Key 写进聊天、命令行参数、脚本常量或工具输出。隐藏输入的 Bash 与 PowerShell 例子见[手动安装指南](installation.md)。如果 Agent 当前环境无法提供用户可操作的安全输入，给出本机执行命令，并明确标记“等待本机输入 Key”，不要声称已完成。

安装器会写入各客户端私有用户配置，Unix 权限为 `0600`；两个供应商连接共享 Skill。不要重写整个 MCP 配置，也不要在冲突时直接使用 `--force`；先核对已有文件。

### 4. 验证并告诉用户结果

运行 `supplier-skill doctor <客户端> --domain <fabric或garment>`。`VERIFIED` 表示安装器实际身份查询已通过；`INSTALLED_AUTHORIZATION_REQUIRED` 表示仍需管理员授予本域身份能力。其他认证、网络、错域或字段错误按实际结果报告。

然后重启或重连客户端，在原生会话中加载 `supplier-open-platform`，发现并调用所选领域的身份工具，输入 `{}`，核对返回的域、主账号、租户、公司和 Client；面料还核对区域，成衣不推断区域。WorkBuddy 首次连接时按其提示信任对应服务器。需要用户重启或点选信任而尚未完成时，明确报告“配置完成，等待重启/信任及原生验证”；不要把 doctor 通过当成原生会话已调用成功。

最终只向用户说明：安装到哪个客户端/领域、诊断与原生验证结果、尚待处理的一步。不要展示完整 Key 或配置文件。当前服务仅开放供应商身份查询；订单、报价、发货能力未开放。
