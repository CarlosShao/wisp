#!/usr/bin/env bash
# 303-v1  ⓐ attack rig: does today's instrument net catch the edit shapes that would turn
# ticket 303's red into a FAKE GREEN?  Runs ONLY in exported trees under $HOME (never the
# repo, never the shared worktree, never the clone the package pair is running in).
# Source tree = `git archive HEAD | tar -x` of the mother repo at 34e1962b.
set -u
MOTHER="/d/work/workspace/projects plans/Wisp"
OUT="$MOTHER/.scratch/wisp/probes/303/v1"
LOG="$OUT/logs"
BASE="$HOME/wisp-303-v1-mut-base"
export PATH="$MOTHER/third_party/sherpa-onnx:$MOTHER/build:$PATH"
YARD='Transport|ShapeA3|LegacySubShape|Unforwarded|ForwardingHook'
NAIL='TestAC14AwaitedBindingReplyReachesThePage'
LIFE='TestPanelHostRealWindowHopAndLifecycle'

fresh() { # $1 = tree dir
  rm -rf "$1"; mkdir -p "$1"
  git -C "$MOTHER" archive HEAD | tar -x -C "$1"
}

yard() { # $1=tree $2=label
  go test -C "$1" ./cmd/wisp/ -count=1 -timeout 300s -v -run "$YARD" > "$LOG/mut-$2-yard.txt" 2>&1
  echo "rc=$?" >> "$LOG/mut-$2-yard.txt"
}
realwin() { # $1=tree $2=label $3=-run regex
  go test -C "$1" ./cmd/wisp/ -count=1 -timeout 300s -v -run "$3" > "$LOG/mut-$2-real.txt" 2>&1
  echo "rc=$?" >> "$LOG/mut-$2-real.txt"
}

# ---- pristine control (HEAD, no mutation) ------------------------------------
fresh "$BASE"
yard "$BASE" "pristine"
realwin "$BASE" "pristine" "$NAIL|$LIFE"

# ---- M1: neuter the fold branch (postMessage always goes to the native exit) --
T="$HOME/wisp-303-v1-mut-M1"; fresh "$T"
F="$T/cmd/wisp/panel_host_windows.go"
python - "$F" <<'PY'
import sys,io
p=sys.argv[1]; s=io.open(p,encoding='utf-8').read()
old='''  cw.postMessage = function (message) {
    if (inside || typeof window.%[1]s !== "function") { return native.call(cw, message); }'''
new='''  cw.postMessage = function (message) {
    return native.call(cw, message);
    if (inside || typeof window.%[1]s !== "function") { return native.call(cw, message); }'''
assert s.count(old)==1, ("M1 anchor not found", s.count(old))
io.open(p,'w',encoding='utf-8',newline='').write(s.replace(old,new))
print("M1 applied")
PY
yard "$T" "M1"; realwin "$T" "M1" "$NAIL|$LIFE"

# ---- M2: move the inside-raise point (flag stuck true from install time) ------
T="$HOME/wisp-303-v1-mut-M2"; fresh "$T"
F="$T/cmd/wisp/panel_host_windows.go"
python - "$F" <<'PY'
import sys,io
p=sys.argv[1]; s=io.open(p,encoding='utf-8').read()
old='  var inside = false;\n  var ex = window.external;'
new='  var inside = true;\n  var ex = window.external;'
assert s.count(old)==1, ("M2 anchor not found", s.count(old))
io.open(p,'w',encoding='utf-8',newline='').write(s.replace(old,new))
print("M2 applied")
PY
yard "$T" "M2"

# ---- M3: fixture stops simulating the direct-binding entry (posts via
#      chrome.webview.postMessage instead) -- test-file edit, still a red->green shape
T="$HOME/wisp-303-v1-mut-M3"; fresh "$T"
F="$T/cmd/wisp/panel_resident_windows_test.go"
python - "$F" <<'PY'
import sys,io
p=sys.argv[1]; s=io.open(p,encoding='utf-8').read()
a='return "(async function(){" +'
assert a in s, "M3 anchor a missing"
old='"window.wispDispatch(" + env("\'ac14r-beacon\'", "\'SCRIPT-RAN\'") + ");" +'
new='"window.chrome.webview.postMessage(" + env("\'ac14r-beacon\'", "\'SCRIPT-RAN\'") + ");" +'
assert s.count(old)==1, ("M3 anchor beacon missing", s.count(old))
s=s.replace(old,new)
old2='"Promise.resolve(window.wispDispatch(ENVS[i])).then(function(v){ return \'REPLIED\'; })," +'
new2='"Promise.resolve(window.chrome.webview.postMessage(ENVS[i])).then(function(v){ return \'REPLIED\'; })," +'
assert s.count(old2)==1, ("M3 anchor awaited missing", s.count(old2))
s=s.replace(old2,new2)
old3='"window.wispDispatch(" + env("REP[i]", "got.replace(/[^A-Za-z0-9-]/g, \'\\')") + ");" +'
new3='"window.chrome.webview.postMessage(" + env("REP[i]", "got.replace(/[^A-Za-z0-9-]/g, \'\')") + ");" +'
assert s.count(old3)==1, ("M3 anchor report missing", s.count(old3))
s=s.replace(old3,new3)
io.open(p,'w',encoding='utf-8',newline='').write(s)
print("M3 applied")
PY
# M3 keeps the shipping code but only reverts the overlay to the pre-fix shape so the
# claim under test is exactly "fixture swap alone turns nail1 green with the fix ABSENT".
python - "$T/cmd/wisp/panel_host_windows.go" <<'PY'
import sys,io
p=sys.argv[1]; s=io.open(p,encoding='utf-8').read()
old='''  var ex = window.external;
  if (ex && typeof ex.invoke === "function") {
    var nativeInvoke = ex.invoke;
    ex.invoke = function (s) {
      inside = true;
      try { return nativeInvoke(s); } finally { inside = false; }
    };
  }
'''
assert s.count(old)==1, ("M3 fix-removal anchor missing", s.count(old))
io.open(p,'w',encoding='utf-8',newline='').write(s.replace(old,''))
print("M3: fix removed as well")
PY
yard "$T" "M3"; realwin "$T" "M3" "$NAIL|$LIFE"

echo "ALL_DONE" >> "$LOG/mut.status"
