# 24 — doc.read + local search tools: PDF/docx extraction, search.files/content/apps

**Status:** **已结案 done**（2026-11-03 13:0x 编排者 24-v3：PDF/docx/扫描件窄过滤 rc=0 全绿〔24r5 0950ddd7c + 24r6 52ed59cff〕；search.content/files/apps 前验 PASS；DEFERRED xlsx/OCR 与 SPEC-12 §5 对上）
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 20-host-bridge-fs-tools
**Parallel slots:** ≤2 sub-agents (A: doc.read PDF/docx; B: search family)
**Spec refs:** SPEC-07 §3, D34, D30② (search taint), D35 spill interplay

## What to build
Document ingestion and local search tools: `doc.read` (PDF + docx → text), `search.files`,
`search.content`, `search.apps` — the knowledge/comprehension base for scenarios ①②, all
taint-aware and spill-integrated.

## Key constraints
- `doc.read` L0, capability fs.read; PDF text extraction + docx XML text extraction with pure-Go
  or whitelisted-dep parsers (propose + justify any new dependency — whitelist change needs
  approval); output capped by spill rules (>4000 tok → artifacts file + head/tail in context);
  content taint-marked; **xlsx + OCR NOT implemented — DEFERRED(D-34) markers + registry entry**.
- Scanned-PDF (no text layer) → clear `doc` error "无可提取文本（扫描件暂不支持 OCR）" —
  no silent empty result.
- `search.files` L0 (name glob within allowlist), `search.content` L0 (full-text within
  allowlist; results taint-marked; huge result sets → spill), `search.apps` L0 (installed-app
  enumeration feeding app.launch; enumeration only, never launches).
- All searches path-resolved via C26; blacklisted targets excluded even from search RESULTS
  listing (A-list filenames must not leak metadata).
- Timeout per tool from config; cancelled mid-scan → partial results marked truncated.

## Out of scope
- Indexing daemon/persistent index (not in plan — on-demand search only); xlsx/OCR (deferred);
  panel result rendering (36).

## Acceptance criteria
- [x] doc.read: text PDF, multi-page PDF >4k tokens (spill path), docx, scanned PDF error case.
- [x] search.content: fixture corpus query → results with taint; blacklist file never appears
      even in listings; timeout → truncated marker.
- [x] search.files glob + allowlist scoping (outside → rejected via R2).
- [x] search.apps enumerates fixture-installed app; feeding result to app.launch works (22).
- [x] DEFERRED(D-34) markers present for xlsx/OCR + registry cross-check passes.

## Progress log (append-only, newest last)

- [2026-11-03 13:0x +08] agent=编排者（结案）did=24-v3 五格全PASS(PDF/docx/scan窄过滤rc=0全绿;search前验;DEFERRED xlsx/OCR对上);24r5 0950ddd7c+24r6 52ed59cff;PDF从推迟→实现.
