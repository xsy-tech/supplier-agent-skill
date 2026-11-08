# Build and release

`make check` runs Go tests and vet. `make build` builds the native CLI. `make dist` cross-compiles five platform archives; `bash scripts/verify-dist.sh` verifies checksums; automated tests check package contents. Archives carry a self-contained CLI with embedded Skill, references and metadata. Use `LC_ALL=C LANG=C` on hosts without `C.UTF-8`.

Source version, manifest and Connector metadata start at `0.1.0`. No remote is configured and no Release exists. Retain the copyright and follow [LICENSE](../LICENSE). Every commit must include GYXT-35342 and the full work-item URL in AGENTS.md; run `make install-hooks`.

After the approved GitHub target exists and this commit has been pushed there, `make release REPO=owner/repository VERSION=0.1.0` runs checks, creates a draft with all five packages and checksums, then publishes it. Existing tags are never overwritten. Do not execute until remote publication is authorized. For other hosting, distribute the same verified packages using that platform's release process. Only document real download locations.
