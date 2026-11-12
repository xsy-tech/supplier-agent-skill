# Identity and credentials

A Key binds one domain, tenant, supplier and main account. Never substitute ERP or another supplier domain, override identity, forge headers or call internal directories to bypass that binding. Treat server rejection of inactive accounts, suppliers or changed authentication versions as authoritative.

Local installation may save a personal Key in private user-level MCP configuration. That permission does not extend to chat, source code, business reports, shared configuration or exported Connector packages. Do not expose SQL, upstream machine tokens or stack traces.

For fabric inventory writes, follow the [inventory workflow](fabric-inventory.md): preview, obtain explicit confirmation and use the separate import or adjustment commit tool. There is no preview ID. Generate a new idempotency key for each approved submission; reuse the original key and identical input only for a network retry of that submission. Never automatically retry an ambiguous timeout; query current stock and ask the user. For other writes, explain the concrete target and impact and obtain confirmation where required. Execute authorized ordinary reads directly.
