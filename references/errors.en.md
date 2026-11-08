# Supplier errors

`VERIFIED` means the selected supplier identity tool was called and its output validated. It does not imply order or shipment capability.

`INSTALLED_AUTHORIZATION_REQUIRED` means initialization and tool discovery succeeded but the expected identity tool is missing. Ask the administrator to grant `supplier.fabric.identity.read` or `supplier.garment.identity.read` for the selected domain, then reconnect.

For 401 or invalid bindings, check the Key, Client, main-account and supplier status with the administrator. For 403, check grants and current account permissions. Do not switch to ERP or another domain. Respect rate-limit guidance for 429. Report timeouts, directory failures and 5xx as failures; never invent empty data or cached verification. Reject wrong-domain or incomplete identity output.

Share only redacted status, time and request identifiers when available. Never share the personal Key or Authorization header.
