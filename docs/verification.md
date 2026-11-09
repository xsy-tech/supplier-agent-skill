# 本地与发布验证 / Local and release verification

日期：2026-11-09；当前版本：0.1.1；工单：GYXT-35342。

发布目标为私有仓库 [xsy-tech/supplier-agent-skill](https://github.com/xsy-tech/supplier-agent-skill)，版本 `v0.1.1` 包含源码、嵌入式 Skill、五个平台安装包和 SHA-256 清单。发布前重新执行下列验证；发布后还会下载 Release 资产，逐一比对校验和并运行本机包。v0.1.0 发布后已使用真实供应商测试 Key 完成 WorkBuddy 原生 MCP 验收，见[实测记录](verification-workbuddy-2026-11-09.md)。v0.1.1 更新版本元数据、接入说明与验收证据，供应商调用逻辑未变；该版本另行执行自动化回归和 Release 包校验。未对生产环境执行联调。

## 验证范围

- `LC_ALL=C LANG=C make check`：Go 全量测试和 vet；覆盖配置保护、凭证脱敏、MCP 分页、超时、认证错误和可重复打包。
- `LC_ALL=C LANG=C go test -race ./internal/cli ./internal/mcpclient ./internal/installer ./internal/profile`：关键配置、诊断和安装模块竞态检测。
- `SUPPLIER_TEST_NATIVE_CLIENTS=1 go test ./internal/installer -run TestNativeClients -v -count=1`：本机真实 Codex、Claude Code CLI 读取隔离目录中的持久化 MCP URL 和鉴权配置，不调用推理，也不使用真实用户凭证。
- 三客户端经过本地真实 MCP/HTTP 测试服务完成两个领域安装、读取已保存凭证和实际身份查询。另有生命周期回归检查 ERP、面料、成衣三条配置并存，单域卸载保留另一域及共享 Skill，最后一个供应商域卸载再清理 Skill。
- 身份结果按领域和字段完整性验证；错域、错误主体类型、缺少主体/租户/公司/Client、错误区域和畸形结果均拒绝。诊断不依赖 ERP 平台信息工具，不调用业务写工具；HTTP 跳转不会转发凭证。
- 损坏的客户端配置在安装文件前预检并拒绝，已有配置和安装记录受保护。
- Skill frontmatter 由 skill-creator `quick_validate.py` 校验；安装资源包含 `agents/openai.yaml` 与中英文参考。
- 本机 CLI、五个平台包和 SHA-256 校验；安装包包含离线安装文档、Skill 和参考资料。全量打包回归检查不同目录、时区和 Git 上下文下的产物一致性。

## 边界

真实 Codex/Claude Code 仅验证本地配置加载。WorkBuddy 5.3.11 已在真实桌面客户端完成 Skill 加载、原生工具发现与两个供应商身份工具调用；使用已发布 v0.1.0 包，5/5 个限定场景通过。首次信任、重启、服务端 MCP SUCCESS 审计和撤销清理均有记录。此原生结果不等同于 v0.1.1 的桌面重测。Windows 安装器未在 Windows/PowerShell 执行；Windows amd64 包已交叉构建并参与归档校验。真实供应商账号失效和目录不可用场景需要在有测试凭证的实际服务环境继续验收；本地认证/错误响应测试不替代该验收。

系统不提供 `C.UTF-8`，测试命令显式使用 `LC_ALL=C LANG=C`。自动化测试使用临时目录。WorkBuddy 原生验收曾临时安装两个供应商连接并签发短期只读 Key；验收后已吊销归档 Key、验证 401，并还原原 ERP 配置、移除临时 Skill 和信任条目。
