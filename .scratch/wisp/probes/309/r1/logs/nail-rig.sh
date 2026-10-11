set -u
LOGS=.scratch/wisp/probes/309/r1/logs
SCRATCH=internal/tools/zz309rig_scratch.go
cleanup() { rm -f "$SCRATCH"; }
trap cleanup EXIT

echo "########## BEFORE (no extra construction point) ##########"
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 \
  -run 'Test309ToolRequestProductionConstructionNail' -v 2>&1 | tail -4
echo "BEFORE_rc=${PIPESTATUS[0]}"

# add ONE non-test production construction point under the scanned root, in my own write face
cat > "$SCRATCH" <<'EOF'
package tools

import "github.com/CarlosShao/wisp/internal/agent"

func rig309EmptyTaskIDShape() agent.ToolRequest {
	return agent.ToolRequest{TaskID: "", Name: "fs.read"}
}
EOF

echo "########## RIG: added agent.ToolRequest{TaskID:\"\"} in a non-test prod file ##########"
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 \
  -run 'Test309ToolRequestProductionConstructionNail' -v 2>&1 | tail -6
echo "RIG_rc=${PIPESTATUS[0]}"

cleanup
echo "########## AFTER restore ##########"
test -f "$SCRATCH" && echo "SCRATCH_STILL_PRESENT" || echo "SCRATCH_REMOVED=yes"
git status --porcelain -- internal/tools
