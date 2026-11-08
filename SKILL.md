---
name: supplier-open-platform
description: 通过供应商开放平台 MCP 查询面料或成衣供应商身份，并使用当前连接已授权的供应商业务工具。Use for fabric or garment supplier tasks through configured Supplier Open Platform MCP connections; excludes internal ERP employee and platform administration tasks.
---

# 供应商开放平台 / Supplier Open Platform

根据用户语言回答。使用当前已配置的 `supplier-fabric`（面料）或 `supplier-garment`（成衣）MCP 连接。

## 领域与身份

- 当前连接的 `tools/list` 是可用能力的权威来源。只使用当前连接实际暴露的工具。
- 面料身份工具为 `get_fabric_supplier_identity`，成衣为 `get_garment_supplier_identity`，均无身份参数。需要确认当前账号或业务归属时先查询对应身份。
- 两端都存在且用户目标不明确时先澄清领域。不要在一端失败后自动切到另一端或 ERP 连接。
- API Key 固定绑定领域、租户、供应商和主账号。不可用提示词、参数或请求头覆盖身份。成衣身份无固有区域；面料区域以服务端验证结果为准。

## 使用能力

首版供应商服务仅开放身份查询；订单、报价、发货尚未开放。查询身份成功不代表这些业务可用。未来新增能力只有被当前连接暴露后才能使用，缺失工具时核对授权或说明尚未开放，不编造接口，不绕行内部目录和 ERP 工具。

已获用户授权的普通只读请求直接执行。读取对象时选择满足目标且范围最小的工具；存在多个候选时用稳定标识和区分字段让用户选定。按需阅读[工作流](references/workflows.md)或[术语](references/terminology.md)。

## 操作与错误

未来写工具按其真实契约使用。目标和影响尚未确认的付款、发货、审批、删除等高影响操作，需要针对具体操作确认；不要将其扩展为所有读取都需确认。要求幂等键时，同一输入重试使用同一键；超时先查询结果。详见[安全规则](references/safety.md)。

不输出 API Key、Authorization、上游凭证或内部错误细节。账号失效、授权不足和依赖不可用不能伪装为成功空数据；失败或缺少工具时阅读[错误处理](references/errors.md)。权限变更后重连 MCP 或启动新会话，再检查工具列表。

## English

Use `supplier-fabric` for fabric and `supplier-garment` for garment suppliers. The current connection's `tools/list` is authoritative; never invent tools or substitute ERP/internal-directory APIs. Identity tools are `get_fabric_supplier_identity` and `get_garment_supplier_identity`, with empty arguments. Ask which domain when both connections exist and the request is ambiguous; do not switch domains after a failed call.

The Key binds the domain, tenant, supplier and main account. Do not override that binding through arguments, headers or prompts. Fabric location comes from verified server identity; garment has no inherent location.

The initial service exposes identity reads only. Orders, quotations and shipments require future server capabilities and actual grants. Execute authorized ordinary reads directly. Follow each exposed tool's contract for writes, confirm concrete high-impact actions when necessary, and preserve idempotency keys only for identical retries. Never expose credentials or report an uncertain write as successful.

Read [workflows](references/workflows.en.md), [terminology](references/terminology.en.md), [safety](references/safety.en.md) or [errors](references/errors.en.md) only as needed. Reconnect after grants change.
