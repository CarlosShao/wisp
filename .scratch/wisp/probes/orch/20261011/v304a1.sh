#!/usr/bin/env bash
# orchestrator re-run of 304-a1's load-bearing numbers (independent rulers)
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/orch/20261011"
L="$D/logs/orch-A850-304a1-verify.txt"
: > "$L"
cd "$R" || exit 1

{ echo "=== HEAD / branch"
  echo "head=$(git rev-parse --short HEAD) branch=$(git rev-parse --abbrev-ref HEAD)"
  echo "anchor_leg_reported=d75db9ee exists=$(git cat-file -t d75db9ee 2>&1)"
} >> "$L"

{ echo "=== ruler 1: spawn statements in internal cmd tools (.go)"
  git -c core.quotePath=false ls-files -- internal cmd tools | grep '[.]go$' > "$D/logs/orch-A850-gofiles.txt"
  echo "go_files_total=$(grep -c '' "$D/logs/orch-A850-gofiles.txt")"
  xargs -a "$D/logs/orch-A850-gofiles.txt" -d '\n' grep -nHE 'exec\.Command(Context)?\(' 2>/dev/null > "$D/logs/orch-A850-spawn-all.txt"
  echo "spawn_all=$(grep -c '' "$D/logs/orch-A850-spawn-all.txt")"
  grep -E '_test[.]go:' "$D/logs/orch-A850-spawn-all.txt" > "$D/logs/orch-A850-spawn-test.txt"
  grep -vE '_test[.]go:' "$D/logs/orch-A850-spawn-all.txt" > "$D/logs/orch-A850-spawn-nontest.txt"
  echo "spawn_test=$(grep -c '' "$D/logs/orch-A850-spawn-test.txt") spawn_nontest=$(grep -c '' "$D/logs/orch-A850-spawn-nontest.txt")"
  echo "--- non-test sites (name:line):"
  cat "$D/logs/orch-A850-spawn-nontest.txt"
  echo "--- distinct files with spawns:"
  cut -d: -f1 "$D/logs/orch-A850-spawn-all.txt" | sort -u | grep -c ''
} >> "$L"

{ echo "=== ruler 2: the two mockllm spawn sites the leg named"
  sed -n '70,74p' internal/adaptertest/mockllm.go 2>/dev/null || echo "MISSING internal/adaptertest/mockllm.go"
  find . -name 'mockllm_integ_test.go' -not -path './.git/*' | while read -r f; do echo "-- $f"; sed -n '71,75p' "$f"; done
} >> "$L"

{ echo "=== ruler 3: jobscope instrument in-repo (甲)"
  ls -l internal/proc/jobscope_windows.go internal/proc/jobscope_windows_test.go 2>&1 | sed 's/  */ /g'
  echo "--- jobscope_windows.go:71 (claim: 71-181 block)"
  sed -n '71p' internal/proc/jobscope_windows.go
  echo "--- jobscope_windows_test.go:92 (claim: 92-104 case)"
  sed -n '92p' internal/proc/jobscope_windows_test.go
  echo "jobscope_exported_fns=$(grep -cE '^func ' internal/proc/jobscope_windows.go)"
  echo "callers_outside_test=$(git -c core.quotePath=false grep -lE 'proc\.(NewJob|AssignProcessToJob|CreateJobObject|JobScope)' -- internal cmd tools 2>/dev/null | grep -vc '_test.go$' || true)"
} >> "$L"

{ echo "=== ruler 4: is tools/mockllm a separate module (乙 premise)"
  echo "mockllm_gomod=$(ls -1 tools/mockllm/go.mod 2>&1)"
  echo "root_gomod_requires_mockllm=$(grep -c 'mockllm' go.mod || true)"
  echo "replace_or_require=$(grep -nE 'mockllm' go.mod 2>/dev/null | tr '\n' ' ')"
  echo "go.work=$(ls -1 go.work 2>&1)"
  echo "module_line=$(head -1 tools/mockllm/go.mod 2>/dev/null)"
  echo "mockllm_imports_observe=$(grep -cE 'wisp/internal/observe|NewTimeout' tools/mockllm/*.go 2>/dev/null | tr '\n' ' ')"
} >> "$L"

{ echo "=== ruler 5: does CI count/kill processes (丙)"
  for pat in mockllm tasklist 'Get-Process' 'Stop-Process' 'Get-CimInstance'; do
    echo "ci_yml_hits[$pat]=$(grep -ciE "$pat" .github/workflows/ci.yml || true)"
  done
  echo "scripts_dir_hits[mockllm]=$(grep -rilE 'mockllm' scripts/ 2>/dev/null | grep -c '' || true)"
} >> "$L"

{ echo "=== ruler 6: live process / temp residue (my own reading now)"
  echo "date=$(date '+%F %T %z')"
  for img in mockllm.exe wisp.exe balldebug.exe msedgewebview2.exe; do
    n=$(tasklist //FI "IMAGENAME eq $img" //FO CSV 2>/dev/null | grep -ciE "^\"$img" || true)
    echo "proc[$img]=$n"
  done
  echo "temp_test_dirs=$(ls -d "$LOCALAPPDATA"/Temp/Test* 2>/dev/null | grep -c '' || true)"
  echo "temp_test_dirs_ALT=$(ls -d /c/Users/swq/AppData/Local/Temp/Test* 2>/dev/null | grep -c '' || true)"
  echo "wisp_side_webview_parents=$(tasklist //FI "IMAGENAME eq wisp.exe" //FO CSV 2>/dev/null | grep -ciE '^"wisp.exe' || true)"
} >> "$L"

{ echo "=== ruler 7: d22scan wall-clock instrument literal"
  grep -n 'Sub(time.Now())' tools/d22scan/main.go | cut -c1-140
  echo "--- main.go:166 (leg claim)"; sed -n '166p' tools/d22scan/main.go
  echo "--- main.go:781,792 (leg claim)"; sed -n '781,792p' tools/d22scan/main.go
  echo "other_timeout_spellings_in_scope=$(git -c core.quotePath=false grep -cE 'time\.After\(|\.Wait\(' -- internal cmd tools 2>/dev/null | grep -c '' )"
} >> "$L"

{ echo "=== ruler 8: ticket 304 text + HANDOVER staleness claim"
  T=$(git -c core.quotePath=false ls-files -- .scratch/wisp/issues | grep -E '/304-' | head -1)
  echo "ticket_path=$T"
  echo "--- ticket :5 and :17"; sed -n '5p;17p' "$T"
  echo "--- HANDOVER lines that allegedly already count mockllm.exe:"
  for n in 2431 2455 2466 2474 2490; do printf "H:%s: %s\n" "$n" "$(sed -n "${n}p" docs/reports/HANDOVER.md | cut -c1-90)"; done
  echo "handover_mockllm_total=$(grep -ciE 'mockllm' docs/reports/HANDOVER.md || true)"
} >> "$L"

{ echo "=== ruler 9: leg artefact denominators (row counts, not prose)"
  for f in 00-anchor 10-residual 20-roster 30-costs 90-final; do
    p=".scratch/wisp/probes/304/a1/$f.md"
    [ -f "$p" ] && printf "%-12s bytes=%s lines=%s\n" "$f" "$(wc -c <"$p")" "$(grep -c '' "$p")" || printf "%-12s MISSING\n" "$f"
  done
  echo "logs_files=$(ls -1 .scratch/wisp/probes/304/a1/logs | grep -c '')"
  echo "zero_byte_logs=$(find .scratch/wisp/probes/304/a1/logs -type f -size 0 | grep -c '' || true)"
} >> "$L"

echo "=== hygiene"
echo "staged=$(git diff --cached --name-only | grep -c '' || true)"
echo "dirty_internal_cmd_tools=$(git status --porcelain -- internal cmd tools | wc -l | tr -d ' ')"
cat "$L"
