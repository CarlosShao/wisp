#!/usr/bin/env bash
# second pass: fix my own path guess + pin the four claims my first ruler left open
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/orch/20261011"
L="$D/logs/orch-A850-304a1-verify2.txt"
: > "$L"
cd "$R" || exit 1

{ echo "=== A) the real mockllm spawner (my first ruler guessed the wrong dir)"
  f=internal/llm/adaptertest/mockllm.go
  echo "exists=$([ -f "$f" ] && echo yes || echo no) bytes=$(wc -c <"$f" 2>/dev/null)"
  echo "--- :55-:95"; sed -n '55,95p' "$f"
  echo "--- does the parent wait/kill? (grep in this file)"
  grep -nE 'Wait\(|Kill\(|t\.Cleanup|defer|ProcessState|Exit\(' "$f" | cut -c1-120
} >> "$L"

{ echo "=== B) jobscope: real Test funcs + who calls it in production"
  echo "--- func Test* in jobscope_windows_test.go:"
  grep -nE '^func Test' internal/proc/jobscope_windows_test.go | cut -c1-110
  echo "--- the 3 non-test files that mention JobScope-ish names:"
  git -c core.quotePath=false grep -lE 'proc\.(NewJob|AssignProcessToJob|CreateJobObject|JobScope)|JobScope' -- internal cmd tools > "$D/logs/orch-A850-jobscope-files.txt"
  cat "$D/logs/orch-A850-jobscope-files.txt"
  echo "--- production (non _test.go) call sites, with line:"
  git -c core.quotePath=false grep -nE 'JobScope' -- internal cmd tools | grep -vE '_test[.]go:' | cut -c1-150
  echo "--- does anything already use job scope to reap children in production?"
  grep -rnE 'OpenJobScope|NewJobScope' --include='*.go' internal cmd tools | grep -vE '_test[.]go:' | cut -c1-150
} >> "$L"

{ echo "=== C) ci.yml: the 11 mockllm references (what do they actually do)"
  grep -niE 'mockllm' .github/workflows/ci.yml | cut -c1-170
  echo "--- any cleanup/kill step in ci.yml at all?"
  grep -niE 'if:.*always\(\)|\bkill\b|kill-process|cleanup' .github/workflows/ci.yml | cut -c1-140
} >> "$L"

{ echo "=== D) does HANDOVER already count mockllm.exe (ticket :5 one staleness claim)"
  for n in 2431 2455 2466 2474 2490; do
    printf "H:%s has_mockllm=%s :: %s\n" "$n" "$(sed -n "${n}p" docs/reports/HANDOVER.md | grep -ciE 'mockllm' || true)" "$(sed -n "${n}p" docs/reports/HANDOVER.md | cut -c1-150)"
  done
  echo "handover_lines_with_mockllm=$(grep -ciE 'mockllm' docs/reports/HANDOVER.md || true)"
  echo "--- every HANDOVER line naming mockllm.exe AND tasklist in the same line:"
  grep -niE 'mockllm' docs/reports/HANDOVER.md | grep -iE 'tasklist|IMAGENAME' | grep -c ''
} >> "$L"

{ echo "=== E) leg roster denominators (row counts of its own tables)"
  p=.scratch/wisp/probes/304/a1/20-roster.md
  echo "--- heading lines:"; grep -nE '^#{2,4} ' "$p" | cut -c1-120
  echo "--- markdown table rows per section:"; grep -cE '^\|' "$p"
  echo "--- separator rows:"; grep -cE '^\|[ :-]+\|' "$p"
  echo "claims_18_48=$(grep -cE '\b18\b.*\b48\b' "$p" || true)"
  for pat in '18' '48' '66' '72'; do echo "occurrences_of_$pat=$(grep -ow "$pat" "$p" | grep -c '' || true)"; done
  echo "--- where 18/48 appear (context):"
  grep -nE '\b(18|48)\b' "$p" | cut -c1-160 | head -20
} >> "$L"

{ echo "=== F) what the leg literally wrote about ci.yml / grep 0"
  for p in .scratch/wisp/probes/304/a1/30-costs.md .scratch/wisp/probes/304/a1/90-final.md .scratch/wisp/probes/304/a1/20-roster.md; do
    echo "--- $p"; grep -niE 'ci\.yml|tasklist' "$p" | cut -c1-170 | head -12
  done
} >> "$L"

echo "=== hygiene"
echo "head=$(git rev-parse --short HEAD) staged=$(git diff --cached --name-only | grep -c '' || true) dirty_internal_cmd_tools=$(git status --porcelain -- internal cmd tools | wc -l | tr -d ' ')"
cat "$L"
