#!/usr/bin/env bash
# 183-v2 §2 恒真性矩阵：三发自取变异 × 逐名判据。
# 还原只用 git cat-file blob HEAD:<path> > <path>（派单 §4）。
set -u
cd "$(git rev-parse --show-toplevel)"
OUT=.scratch/wisp/probes/183/v2
mkdir -p "$OUT/logs"
PY=$(command -v python || command -v python3)
if [ -z "$PY" ]; then echo "NO PYTHON AVAILABLE"; exit 3; fi

apply() { # apply <pyexpr-name>
  "$PY" - "$1" <<'PYEOF'
import sys, io, re, pathlib
name = sys.argv[1]
def edit(path, old, new):
    with io.open(path, encoding='utf-8', newline='') as fh:
        t = fh.read()
    if t.count(old) != 1:
        print("PATTERN-NOT-UNIQUE-ABORT", path, t.count(old)); sys.exit(2)
    with io.open(path, 'w', encoding='utf-8', newline='') as fh:
        fh.write(t.replace(old, new))
if name == 'M1':
    # 摘掉 contains() 里 spellsDeclaredPath 那段跳过（值规则整支删除）
    edit('internal/risk/taintmatch.go',
"""		if len(f.declaredPath) >= n && spellsDeclaredPath(f.declaredPath, win) {
			// Ticket 183: this candidate window literally spells part of the path
			// the host itself wrote into THIS mark, so it is no evidence against
			// the model that was handed that path - at any position, which is
			// exactly what the positional skip span above could not promise.
			continue
		}
""", "")
elif name == 'M2':
    # 摘掉 provenance.go 里 attachDeclaredPath 那一支（保留位置性 skip）
    edit('internal/risk/provenance.go',
"""	if declaredNorm != "" {
		// Ticket 183 (root cause 183-a1 judged as (d)): the span above excludes
		// windows by POSITION while the index hashes and the verification corpus
		// compare by STRING SET, so one repetition of the path's own 8-rune
		// spelling anywhere else in this body re-indexed the exemption and the
		// model's reread of the pointer the HOST wrote hit R4 (20000 bytes of
		// artifact, 0 bytes read back). Attaching the declared string to this
		// index makes its own windows non-evidence at every position, while
		// every fragment that spells nothing of it stays indexed exactly as
		// before - and the mark is still local here, so nothing observes a
		// half-built index.
		m.idx.attachDeclaredPath(declaredNorm)
	}
""", """	_ = declaredNorm // MUTATION M2 (183-v2): positional span kept, value rule never attached
""")
elif name == 'M4':
    # 把"逐字拼出声明串"放宽成"差一枚 rune 也算"（前缀/近似放宽那一族）
    edit('internal/risk/taintmatch.go',
"""	for i := 0; i+len(window) <= len(declared); i++ {
		match := true
		for j := range window {
			if declared[i+j] != window[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false""",
"""	for i := 0; i+len(window) <= len(declared); i++ {
		diffs := 0
		for j := range window {
			if declared[i+j] != window[j] {
				diffs++
			}
		}
		if diffs <= 1 {
			return true
		}
	}
	return false""")
elif name == 'M5':
    # 洗戳族：声明过路径的整枚 mark 不再算证据（值规则退化成整枚豁免）
    edit('internal/risk/taintmatch.go',
"""	if len(rp) < n || len(f.hashes) == 0 {
		return "", false
	}
""", """	if len(rp) < n || len(f.hashes) == 0 {
		return "", false
	}
	if len(f.declaredPath) > 0 {
		return "", false // MUTATION M5 (183-v2): whole-mark laundering
	}
""")
elif name == 'M0':
    # 总开关：两半豁免（位置跨度 + 值规则）一起摘掉＝票 177 之前的状态
    edit('internal/risk/provenance.go',
"""			skip = [][2]int{{lo, lo + len(nhp)}}
			declaredNorm = hp
""", """			// MUTATION M0 (183-v2): shape-A exemption switched off entirely
""")
elif name == 'M3':
    # 不要求 runeIndexOf 命中就无条件挂 declaredNorm（fail-closed 那一支有没有牙）
    edit('internal/risk/provenance.go',
"""			// The declared path is not in this content: the exemption claims
			// nothing it can prove, so it excludes nothing (fail-closed).
			logf("risk/C25: declared host path absent from mark %s origin=%q content after normalization: nothing excluded", tool, origin)""",
"""			// MUTATION M3 (183-v2): declared path absent, attach anyway.
			logf("risk/C25: MUTATION M3 declared host path absent but attached", tool, origin)
			declaredNorm = hp""")
else:
    print("UNKNOWN-MUTATION", name); sys.exit(4)
print("APPLIED", name)
PYEOF
}

restore() { # restore <paths...>
  for f in "$@"; do git cat-file blob "HEAD:$f" > "$f"; done
}

run_case() { # run_case <tag> <mutspec>
  local tag="$1" spec="$2"
  if [ "$spec" = "HEAD" ]; then
    :
  else
    apply "$spec" || { echo "APPLY-FAILED $spec"; return 1; }
  fi
  go test ./internal/risk/ -run 'TestPointer183' -v > "$OUT/logs/mut-$tag.txt" 2>&1
  echo "== $tag exit=$? ==" >> "$OUT/logs/mut-$tag.txt"
  grep -E '^(=== RUN|--- (PASS|FAIL)|ok|FAIL|PASS|== )' "$OUT/logs/mut-$tag.txt" >> "$OUT/logs/summary-$tag.txt"
  restore internal/risk/provenance.go internal/risk/taintmatch.go
}

SPECS=${*:-HEAD M1 M2 M3 M4}
for spec in $SPECS; do
  run_case "$spec" "$spec"
done

echo "=== per-case red names (FAIL lines) ==="
for spec in $SPECS; do
  echo "--- $spec ---"
  grep -E '^--- FAIL|^    --- FAIL|^ok|^FAIL|^== ' "$OUT/logs/mut-$spec.txt" | sort -u
done
echo "=== final git status internal/ cmd/ ==="
git status --porcelain -- internal/ cmd/
echo "(empty above = restored)"
