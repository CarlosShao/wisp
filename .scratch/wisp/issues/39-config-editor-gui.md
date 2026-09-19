# 39 — Config editor GUI: sections ↔ structs round-trip, three-tier badges, security locks

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 05-config-model, 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: schema-driven form rendering; B: security-section flow +
restart/reload UX)
**Spec refs:** SPEC-08 §5.2, D6, D36, §14.5, design/screens/config.html

## What to build
The config editor as a pure VIEW over config.toml: schema-driven forms for every section,
three-tier effect badges, security-section lock affordances with L2 re-confirm flow, and
strict single-source-of-truth discipline (GUI holds nothing TOML doesn't).

## Key constraints
- Forms generated from the SAME Go schema structs (bridge `config.get/set`); GUI never owns
  defaults or extra state (D36 rule 4; Zod for input feedback only — Go stays truth).
- Three-tier badges: hot (applies now) / reload (subsystem reload note) / restart (banner
  "需重启生效"); after save, applied-tier feedback surfaced from the reload events.
- 🔒 sections (risk/fs/net/plugins): lock icon; loosening edits trigger the L2 re-confirm
  (native card), reject → keep old value; tightening applies hot; all directions logged.
  Danger-styled affordance, never a silent toggle.
- api_key fields: write-only inputs (never render stored values); set → creates/updates
  `api_key_ref` via SecretStore (06); placeholder shows last-4 only.
- Unknown-key handling: file edited externally with unknown keys → validation errors listed
  with line numbers (05 errors surfaced verbatim).
- Migration/mismatch: `schema_version` mismatch → guided flow, no silent overwrite.
- Round-trip guarantee: open → save unchanged → file diff empty (formatting excepted).

## Out of scope
- First-run wizard (firstrun screen lands with 41/14 polish; basic flow here); models page
  heavy UX (40 covers listing).

## Acceptance criteria
- [ ] Every D36 section renders from schema (golden snapshot per section); edit → save →
      hot/reload/restart behavior matches tier (events observed).
- [ ] Security loosening → native L2 confirm → accept applies + logs; reject keeps old
      (e2e with 21's card).
- [ ] api_key write-only + last-4 display + ref creation verified; plaintext never rendered
      or stored by GUI.
- [ ] Round-trip no-op save produces no diff; external-edit + unknown key shows line-numbered
      errors.
- [ ] Visual sign-off vs design/screens/config.html (lock icons on security sections).

## Progress log (append-only, newest last)
