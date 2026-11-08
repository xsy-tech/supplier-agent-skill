# 供应商术语

| 术语 | 含义 |
| --- | --- |
| 面料 / fabric | `FABRIC_SUPPLIER`，MCP `/fabric/mcp`，别名 `supplier-fabric` |
| 成衣 / garment | `GARMENT_SUPPLIER`，MCP `/garment/mcp`，别名 `supplier-garment` |
| Client | 调用应用，固定绑定领域、租户、具体供应商及正式主账号 |
| API Key | Client 下签发的个人调用凭证，继承绑定及明确授予的能力 |
| 主账号 / subject_id | 服务端确认的供应商正式主账号，不是 ERP 员工 ID |
| supplier_code | 该领域、租户内的门户公司归属编码，不等于可互换的 SCM 供应商 ID |
| Capability | 服务端注册且授予 Key 的能力，业务操作还校验当前权限与对象范围 |
| location | 面料当前有效账号的区域；成衣身份无固有区域 |

公司编码和账号 ID 在不同领域或租户相同，不表示是同一身份。
