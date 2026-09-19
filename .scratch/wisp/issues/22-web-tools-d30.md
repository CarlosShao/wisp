# 22 — Web & open tools: web.search/fetch/open, file.open, app.launch, D30 five layers

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 19-provenance-c25, 20-host-bridge-fs-tools
**Parallel slots:** ≤2 sub-agents (A: web.fetch/search + SSRF/redirect/URL rules; B: web.open/
file.open/app.launch + protocol whitelist + first-visit prompt)
**Spec refs:** SPEC-07 §3, SPEC-06 §6, D23, D30, 16.5.3, D33/F2 (no-HTML rule), §15#3

## What to build
The information-access tool family: web.search, web.fetch (→Markdown), web.open, file.open,
app.launch — with the D30 outbound constraints, the split protocol whitelists (16.5.3), and the
first-visit-new-domain prompt.

## Key constraints
- `web.search` (L0; taint in query → L2 via R4): implementation path is an OPEN decision
  (scrape results page vs search API) — **status: blocked-decision**; default implement behind
  a `SearchProvider` interface with the scrape provider + compose mock-search; record §15#3
  outcome when user decides; switching providers must not touch call-sites.
- `web.fetch` (L0): private-network ranges HARD-DENY (SSRF list: 10/8, 172.16/12, 192.168/16,
  127/8, 169.254/16 incl. 169.254.169.254, ::1, fc00::/7, fe80::/10); no redirect following to
  non-http(s); URL length cap 2048 (SPEC); response size/time caps (reuse spill caps); result →
  Markdown, **never rendered as HTML anywhere** (F2); response taint-marked; body never logged.
- `web.open` (L1): ONLY http/https — everything else rejected (D30③ strict side).
- `file.open` (L1, separate whitelist, default-deny): `file:` only inside allowlist via C26;
  `mailto:`; `ms-settings:`/known `ms-*`; explicit-deny: `search-ms:`, `mshta:`, `javascript:`,
  `vbscript:`, `shell:`, `tel:`, any unknown scheme.
- `app.launch` (L1): name/path/UWP AUMID resolution via `search.apps`; launched processes enter
  Job Object (C30).
- D30 layer 5: first visit to a new domain via builtin web tools → one-time prompt (rememberable,
  per-domain, stored in config/db).
- compose `mock-web` fixtures: clean page, injection-payload page (for 19/25 tests), redirect
  chain, private-IP host.

## Out of scope
- Panel rendering of results (36); robots/etiquette beyond caps; cloud search API integration
  pending user decision.

## Acceptance criteria
- [ ] SSRF suite: every private range form (incl. decimal/octal/hex IP encodings, DNS-rebinding
      fixture via mock) → denied BEFORE connection (receiver asserts zero hits).
- [ ] Redirect chain test: https→http and https→file/private → blocked; response→Markdown
      conversion verified; huge response → spill path.
- [ ] Protocol matrix: web.open/file.open allow/deny table per 16.5.3 — exhaustive test.
- [ ] app.launch launches (notepad fixture) under Job; AUMID + name + path paths work.
- [ ] Taint-in-query → L2 (R4 integration); first-visit prompt appears once, remembered.
- [ ] `web.search` OPEN decision: ticket updated with user's chosen path OR scrape-default
      noted with TODO block for 25 review.

## Progress log (append-only, newest last)
