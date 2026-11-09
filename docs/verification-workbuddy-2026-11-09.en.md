# Native WorkBuddy supplier MCP acceptance

2026-11-09 (UTC+08), GYXT-35342. **5/5 scoped cases passed.** See the [full record and screenshots](verification-workbuddy-2026-11-09.md).

The verified GitHub v0.1.0 macOS arm64 package installed two temporary, identity-only supplier connections in real WorkBuddy 5.3.11 against the project1 local platform. After restart and individual server trust, the native session loaded `supplier-open-platform`, discovered both domain-specific tools and called each once with `{}`. Installer preflight calls were counted separately; no shell or substitute HTTP client supplied the native acceptance result.

| Case | Result | Evidence |
| --- | --- | --- |
| Skill loading | PASS | Native Skill loading trace |
| Tool discovery and ERP coexistence | PASS | ERP 11/11; fabric 1/1; garment 1/1 |
| Garment identity | PASS | Account 32, company 18, no inherent location; MCP SUCCESS audit |
| Fabric identity | PASS | Account 2035, company SYM19438, QUANZHOU; MCP SUCCESS audit |
| Domain identity separation | PASS | Each server exposed its own tool and returned its expected tenant, account and Client |

Both temporary Keys were revoked and archived; requests with them then returned HTTP 401. The original ERP configuration was restored, temporary supplier connections/Skill/trust entries removed, and WorkBuddy reopened without debugging. Audit and archived Key history remain by design.

Native acceptance belongs to **v0.1.0**. Version **v0.1.1** adds documentation, evidence and version metadata without changing supplier call logic; automated regression and archive verification are performed separately. These five cases do not cover production, other accounts, negative cross-tenant attacks, native account-state transitions or future order/quotation/shipment tools.
