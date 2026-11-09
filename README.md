# 供应商 Agent Skill

中文 | [English](README.en.md)

让 Codex、Claude Code 或 WorkBuddy 连接供应商开放平台。当前支持面料、成衣供应商身份查询。

## 复制给 Agent，自动安装

将下面整段发给你正在使用的 Agent，无需自己下载、解压或编辑配置：

```text
请按下面的安装指南，帮我把供应商开放平台 Skill 安装到当前客户端，并完成连接验证：
https://github.com/xsy-tech/supplier-agent-skill/blob/main/docs/agent-install-prompts.md
缺少供应商类型或 MCP 地址时，请一次向我确认。
API Key 请让我在本机安全输入，不要让我粘贴到聊天中。
```

准备好管理员提供的 **供应商类型、MCP 地址和 API Key** 即可。类型和地址可以告诉 Agent，Key 在本机输入。安装指南所在仓库为私有；无法访问时，先登录有权限的 GitHub 账号，或向管理员索取安装包与指南。

Agent 会选择本机安装包、校验、安装并检查身份连接。若客户端需要重启或首次信任，会告诉你下一步；完成原生验证后即可使用。

## 安装后怎么用

在客户端中说：

> 使用供应商开放平台 Skill，查询我的当前供应商身份，列出账号、租户和公司编码。

同时安装了两类供应商连接时，请说明面料或成衣。当前只开放身份查询，订单、报价、发货尚未开放。

## 需要时再看

- [Agent 安装执行指南](docs/agent-install-prompts.md)：供 Agent 读取的安装流程。
- [手动安装、升级、卸载](docs/installation.md)：终端操作和各客户端配置。
- [使用说明](docs/agent-usage.md) · [故障处理](references/errors.md)
- [Release v0.1.1](https://github.com/xsy-tech/supplier-agent-skill/releases/tag/v0.1.1)：五个平台包、SHA-256 和离线文档，下载需仓库权限。
- [WorkBuddy 原生实测](docs/verification-workbuddy-2026-11-09.md)：v0.1.0 限定场景 5/5 通过；v0.1.1 文档与包验证单独记录。
- [构建与发布](docs/enterprise-release.md) · [Skill 指令](SKILL.md) · [验证记录](docs/verification.md) · [LICENSE](LICENSE)
