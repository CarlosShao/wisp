#!/usr/bin/env bash
# 301-r1 roster extractor. ⛔ part of the product; it only reads log files.
#
# Two rulers are emitted side by side and are NEVER collapsed into one number
# (ticket 301 AC#2: "尺要写清…两把都交、⛔ 混成一枚数").
#
#   ruler A = top-level ONLY. Mode string anchored at column 0:
#             ^--- (PASS|FAIL|SKIP): <name>
#             Go prints subtest results indented by four spaces, so column 0 is
#             exactly the top level. ⛔ 含子测试.
#   ruler B = includes subtests. Mode string:
#             ^[[:space:]]*--- (PASS|FAIL|SKIP): <name>
#             <name> is [^ ]+, so a subtest keeps its `Parent/sub` spelling.
#
# Both rulers take the union of PASS|FAIL|SKIP = "was EVALUATED" (a name that
# produced a result line). A PASS-only roster is emitted too, because AC#4's
# "具名新增红" assertion is about FAILs and "新增求值" is about the union.
#
# Comment lines can never enter: every mode string requires `--- ` at the start
# of the line (column 0, or whitespace then `--- `), and Go test output has no
# comment syntax; the log's own `portable-tests.sh: ` banner lines do not match.
set -u

log="$1"
out="$2" # prefix, e.g. .../rosters/windows-before

sed -n 's/^--- \(PASS\|FAIL\|SKIP\): \([^ ]*\).*/\2/p' "$log" | sort -u >"$out.A-evaluated.txt"
sed -n 's/^[[:space:]]*--- \(PASS\|FAIL\|SKIP\): \([^ ]*\).*/\2/p' "$log" | sort -u >"$out.B-evaluated.txt"
sed -n 's/^--- PASS: \([^ ]*\).*/\1/p'                                 "$log" | sort -u >"$out.A-pass.txt"
sed -n 's/^--- FAIL: \([^ ]*\).*/\1/p'                                 "$log" | sort -u >"$out.A-fail.txt"
sed -n 's/^--- SKIP: \([^ ]*\).*/\1/p'                                 "$log" | sort -u >"$out.A-skip.txt"

printf '%s\tA-evaluated=%s\tB-evaluated=%s\tA-pass=%s\tA-fail=%s\tA-skip=%s\n' \
    "$(basename "$out")" \
    "$(wc -l <"$out.A-evaluated.txt" | tr -d '[:space:]')" \
    "$(wc -l <"$out.B-evaluated.txt" | tr -d '[:space:]')" \
    "$(wc -l <"$out.A-pass.txt" | tr -d '[:space:]')" \
    "$(wc -l <"$out.A-fail.txt" | tr -d '[:space:]')" \
    "$(wc -l <"$out.A-skip.txt" | tr -d '[:space:]')"
