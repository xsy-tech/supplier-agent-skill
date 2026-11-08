# Identity and credentials

A Key binds one domain, tenant, supplier and main account. Never substitute ERP or another supplier domain, override identity, forge headers or call internal directories to bypass that binding. Treat server rejection of inactive accounts, suppliers or changed authentication versions as authoritative.

Local installation may save a personal Key in private user-level MCP configuration. That permission does not extend to chat, source code, business reports, shared configuration or exported Connector packages. Do not expose SQL, upstream machine tokens or stack traces.

The initial service is read-only. For future writes, explain the concrete target and impact and obtain confirmation where the action requires it. Execute authorized ordinary reads directly. After a write times out, query its state; reuse an idempotency key only for identical input when the contract permits retry.
