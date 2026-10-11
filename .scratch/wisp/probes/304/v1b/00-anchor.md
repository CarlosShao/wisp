# 304-v1b 00 anchor (own readings, leg 304-v1b)

- Own HEAD reading (this leg, 2026-10-11): `git rev-parse --short HEAD` = `3b8066a2`, full `3b8066a27032881ccbed08ecf3b4a613b863a8cb`, branch `dev`.
- staged = `git diff --cached --name-only` piped `wc -l` = **0**.
- Production scope clean: `git status --porcelain -- internal cmd tools scripts` = **0** lines.
- Shared worktree is dirty elsewhere (other legs / design / .gitignore); this leg touches none of it.
- Predecessor leg 304-v1: committed only its anchor `c27928d6`; its uncommitted gate logs were swept by orchestrator commit `3b8066a2`. On-disk v1 tree = `00-anchor.md`, `15-gates-ac13.md`, `logs/a00..a03,g1-*,g2-*,g3-*`.
- Zero-byte census (own reading, `find .scratch/wisp/probes/304 -type f -size 0`): **6** files ->
  v1: `logs/g2-fmt-worktree.txt`, `logs/g3-build-body.txt`, `logs/g3-vet-body.txt`;
  r1: `logs/build-1.txt`, `logs/vet-1.txt`, `logs/vet-reverse.txt`.
- Discrepancy vs continuation brief: brief names 4 zero-byte files incl. `labels.stderr`; on disk no `labels.stderr` exists anywhere under probes/304 (find, whole tree = 100 files). The orchestrator note names 3 zero-byte in 304-r1 (matches on-disk r1 three). Carried forward.
- Addenda (2026-10-11 19:0x, committed as 1st addendum commit): re-measured HEAD for pandoc-free exactness = `3b8066a2 编排者代提：两枚验收腿（额度掐断）＋304-r1／309-r3 的未提交小件入库（判语⛔ 代填）`; v1 tree file count = 12 in `probes/304/v1/**` (2 .md + 10 logs .txt).
- Write surface declared: `.scratch/wisp/probes/304/v1b/**`, `docs/evidence/s1/**`, one append-only bullet at ticket end. Nothing else.
- A world_facts checksum of this reading lives in `logs/c00-anchor.txt` (self-rc landing).
