# 182 r4 leg — AC#4 rulers (script 01): ticket 181 vs content-review division
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
T181=$(ls .scratch/wisp/issues/181-*.md)
T189=$(ls .scratch/wisp/issues/189-*.md)
echo "== T181 file: $T181 =="
echo "-- T181 box state (checked/unchecked counts) --"
grep -c "^\s*- \[x\]" "$T181" || true
grep -c "^\s*- \[ \]" "$T181" || true
echo "-- T181 AC titles --"
grep -nE "^\s*-\s*\[[x ]\]\s*\*\*AC" "$T181"
echo "== T189 file: $T189 =="
echo "-- T189 box state --"
grep -c "^\s*- \[x\]" "$T189" || true
grep -c "^\s*- \[ \]" "$T189" || true
grep -nE "^\s*-\s*\[[x ]\]\s*\*\*AC" "$T189"
echo "== git.go GitView json fields (display vs content axes) =="
grep -n "json:" internal/panel/git.go
echo "== git.go Uncommitted/Diff anchors =="
grep -n "Uncommitted\|Diff " internal/panel/git.go | head -12
echo "== git.go runGitReadOnly + write-verb literals (0 expected; positive control: query literals) =="
grep -cn "runGitReadOnly(" internal/panel/git.go || true
grep -nE "exec\.Command|commit|push|add |checkout|rebase|reset" internal/panel/git.go | head -8
echo "-- positive control: read-only query literals in git.go --"
grep -cE "status|diff|rev-parse|ls-files" internal/panel/git.go || true
echo "== '归口不是闭格手段' rule source =="
grep -rn "归口" .scratch/wisp/issues/README.md AGENTS.md 2>/dev/null | head -6
grep -c "归口" .scratch/wisp/issues/181-*.md || true
exit 0
