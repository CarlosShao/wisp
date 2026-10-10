#!/usr/bin/env bash
# 303-v1  ⓐ M4: the shape a lazy leg COULD land instead of the JS wrapper -- unwrap the
# double-wrapped frame on the GO side (accept a frame whose method is the binding name and
# re-dispatch its single param).  Does any nail in the repo bite it?  Export tree only.
set -u
MOTHER="/d/work/workspace/projects plans/Wisp"
LOG="$MOTHER/.scratch/wisp/probes/303/v1/logs"
T="$HOME/wisp-303-v1-mut-M4"
export PATH="$MOTHER/third_party/sherpa-onnx:$MOTHER/build:$PATH"
rm -rf "$T"; mkdir -p "$T"; (cd "$MOTHER" && git archive HEAD | tar -x -C "$T")

python - "$T/cmd/wisp/panel_host_windows.go" <<'PY'
import sys, io
p = sys.argv[1]; s = io.open(p, encoding='utf-8').read()
old = '''func (m *PanelManager) dispatchRaw(ctx context.Context, raw string) (string, error) {
	if m.disp == nil {'''
new = '''func (m *PanelManager) dispatchRaw(ctx context.Context, raw string) (string, error) {
	// M4 MUTATION (303-v1, export tree only): a Go-side "helpful" unwrap of the
	// double-wrapped frame the fb2fb802 regression produced, instead of the landed
	// external.invoke wrapper. If no nail bites this, the net has a hole.
	var probe253 struct {
		ID     int               `json:"id"`
		Method string            `json:"method"`
		Params []json.RawMessage `json:"params"`
	}
	if json.Unmarshal([]byte(raw), &probe253) == nil && probe253.Method == panelDispatchBinding && len(probe253.Params) == 1 {
		var inner string
		if json.Unmarshal(probe253.Params[0], &inner) == nil && inner != "" {
			return m.dispatchRaw(ctx, inner)
		}
	}
	if m.disp == nil {'''
assert s.count(old) == 1, ("M4 anchor missing", s.count(old))
s = s.replace(old, new)
if '\t"encoding/json"' not in s:
    s = s.replace('import (\n\t"context"\n', 'import (\n\t"context"\n\t"encoding/json"\n', 1)
io.open(p, 'w', encoding='utf-8', newline='').write(s)
print("M4 applied")
PY

# Revert the landed JS wrapper too, so the ONLY thing carrying nail1 is the Go-side unwrap.
python - "$T/cmd/wisp/panel_host_windows.go" <<'PY'
import sys, io
p = sys.argv[1]; s = io.open(p, encoding='utf-8').read()
old = '''  var ex = window.external;
  if (ex && typeof ex.invoke === "function") {
    var nativeInvoke = ex.invoke;
    ex.invoke = function (s) {
      inside = true;
      try { return nativeInvoke(s); } finally { inside = false; }
    };
  }
'''
assert s.count(old) == 1, ("M4 wrapper-removal anchor missing", s.count(old))
io.open(p, 'w', encoding='utf-8', newline='').write(s.replace(old, ''))
print("M4: landed JS wrapper removed as well")
PY

{ echo "### M4 build (in the exported tree only)"; go build -C "$T" ./cmd/wisp/ 2>&1; echo "rc=$?"; } > "$LOG/mut-M4-build.txt" 2>&1
go vet -C "$T" ./cmd/wisp/ > "$LOG/mut-M4-vet.txt" 2>&1; echo "rc=$?" >> "$LOG/mut-M4-vet.txt"
go test -C "$T" ./cmd/wisp/ -count=1 -timeout 300s -v -run 'Transport|ShapeA3|LegacySubShape|Unforwarded|ForwardingHook' > "$LOG/mut-M4-yard.txt" 2>&1
echo "rc=$?" >> "$LOG/mut-M4-yard.txt"
go test -C "$T" ./cmd/wisp/ -count=1 -timeout 300s -v -run 'TestAC14AwaitedBindingReplyReachesThePage|TestPanelHostRealWindowHopAndLifecycle|TestAC14GoSideEvalPushReachesThePage' > "$LOG/mut-M4-real.txt" 2>&1
echo "rc=$?" >> "$LOG/mut-M4-real.txt"
echo "M4_DONE" >> "$LOG/mut.status"
