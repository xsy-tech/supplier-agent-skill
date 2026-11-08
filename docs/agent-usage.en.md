# Supplier Agent usage

After installation, restart your client. Ask: “Use the garment connection to show my current supplier account, tenant and company code.” Fabric uses its own identity tool and verified location. If both are configured, specify the domain.

The current service exposes identity reads only. An order, quotation or shipment request requires a real exposed business tool and server grant. The Skill cannot create those capabilities locally. Reconnect after grants change. Diagnostics do not invoke writes; errors must not be presented as successful empty results. See [Skill](../SKILL.md) and [errors](../references/errors.en.md).
