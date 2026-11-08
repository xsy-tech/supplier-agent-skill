# Supplier terminology

| Term | Meaning |
| --- | --- |
| fabric | `FABRIC_SUPPLIER`, `/fabric/mcp`, connection `supplier-fabric` |
| garment | `GARMENT_SUPPLIER`, `/garment/mcp`, connection `supplier-garment` |
| Client | Application bound to a domain, tenant, supplier and main account |
| API Key | Credential inheriting the Client binding and explicit grants |
| subject_id | Verified supplier main-account ID; never an ERP employee ID |
| supplier_code | Portal company code within the domain and tenant; not interchangeable with an SCM supplier ID |
| Capability | Registered and granted operation, further bounded by current business permissions |
| location | Current fabric account's verified region; absent for garment identity |

Identical account IDs or company codes across domains or tenants do not merge identities.
