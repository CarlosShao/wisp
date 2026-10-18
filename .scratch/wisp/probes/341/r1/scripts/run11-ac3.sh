#!/bin/bash
# 341-r1 AC#3: name the over-claiming sentences in the ticket-226 verdict piece. ⛔ 落笔 (read-only 指认).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
H=$(git rev-parse HEAD)
out=.scratch/wisp/probes/341/r1/logs/10-ac3-overclaim.txt
{
  echo "HEAD at measurement = $H   (行号会漂 ⇒ 每次引用都带取数时刻的短号; 本件用 $(git rev-parse --short HEAD))"
  echo
  echo "[which verdict files exist]"
  git ls-tree --name-only HEAD docs/evidence/s1/ | grep '^docs/evidence/s1/226' || echo "(none tracked)"
  ls -1 docs/evidence/s1/ | grep '^226' || echo "(none on disk)"
  echo
  echo "[the sentences: anything claiming comments / key order / bytes are KEPT, in the 226 verdict piece]"
  for f in docs/evidence/s1/226*.md; do
    echo "---- $f ----"
    grep -nE '注释|键序|不白重排|不白丢|重排|保留' "$f"
  done
  echo
  echo "[section headers of that piece, to name section 1 and section 5 precisely]"
  for f in docs/evidence/s1/226*.md; do echo "---- $f ----"; grep -nE '^#{1,3} ' "$f"; done
  echo
  echo "[what the production code says TODAY, verbatim, with this anchor's line numbers]"
  git grep -n 're-laid in canonical form' HEAD -- internal/config
  echo
  echo "[and what ticket 227 AC#1③ added: the delete-shape clause, verbatim]"
  grep -n 'a key line the' internal/config/writeguard.go
  echo
  echo "[measured in THIS leg's runs: do comments survive a guarded write?]"
  echo "  P1 case asserts the comment is GONE and is green at baseline: writeguard_handedit_227_test.go:60-62"
  echo "  P2 case (this leg) asserts the same for the INFO branch:      writeguard_unnailed_341_test.go (default-to-me check)"
  echo "  P5 case (this leg) asserts the comment dies on the write branch too (hand note check)"
  echo "  hand key order: b2 line 8/16/25/34/42/50 all say \"[ball] before [app] in the file: false\""
} > "$out" 2>&1
# Ruler fix, named: the first run of this gauge wrote an INVALID-UTF-8 log because `cut -c1-220` slices a
# multi-byte CJK character in half - the same defect as run0. The cuts are removed rather than patched over:
# a file quoting a verdict piece must never be byte-truncated at all. The first corrupt file is kept on disk.
echo "written $(wc -c < "$out") bytes lines $(wc -l < "$out")"
printf 'valid UTF-8? %s\n' "$(iconv -f utf-8 -t utf-8 "$out" >/dev/null 2>&1 && echo YES || echo NO)"
