# Supplier workflows

Choose fabric or garment before checking the account. Inspect tools on that connection and call its identity tool with empty arguments. Report the returned domain, main-account ID, tenant, supplier company code and Client code. Use returned location for fabric; do not invent one for garment or infer company names from codes.

Initially only identity reads are available. If the user requests an order, quotation or shipment without a corresponding exposed tool, explain the unavailable capability and ask the administrator to check availability and grants. An identity result does not authorize business actions. Use future business tools only when exposed; resolve ambiguous objects before writing and follow the actual tool contract.

After the administrator changes grants, reconnect or start a new session. This Skill neither issues Keys nor modifies portal accounts or permissions.
