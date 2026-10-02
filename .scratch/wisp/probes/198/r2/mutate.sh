#!/usr/bin/env bash
# 198-r2 变异台件：每枚先落文件、finally 还原（题面「变异台件纪律」）。
# 每枚的判语口径：rc 取自 go test 本体（不是管道），红名从落下来的文件里读。
set -u
cd "$(git rev-parse --show-toplevel)" || exit 1
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"

PROBE=.scratch/wisp/probes/198/r2
BK=$PROBE/backup
mkdir -p "$BK"
FILES=(cmd/wisp/firstrun.go cmd/wisp/run.go internal/config/defaults.go)
for f in "${FILES[@]}"; do cp "$f" "$BK/$(basename "$f")"; done

restore_all() {
  for f in "${FILES[@]}"; do
    [ -f "$BK/$(basename "$f")" ] && cp "$BK/$(basename "$f")" "$f"
  done
}
trap restore_all EXIT

apply_and_run() { # $1=tag
  echo "==================== MUTATION $1 ===================="
  python - "$1" <<'PY'
import pathlib, sys
tag = sys.argv[1]
edits = {
 "M1": [("cmd/wisp/firstrun.go",
   "\tif err := config.SaveFile(cfgPath, config.NewDefaults()); err != nil {",
   "\tc := config.NewDefaults() // 198r2-M1 planted legal non-default\n\tc.App.Theme = \"light\"\n\tif err := config.SaveFile(cfgPath, c); err != nil {")],
 "M2": [("internal/config/defaults.go",
   "\tc := &Config{}\n\tapplyDefaults(reflect.ValueOf(c).Elem())\n\treturn c",
   "\tc := &Config{}\n\tapplyDefaults(reflect.ValueOf(c).Elem())\n\tc.Session.WarmTimeoutSec = 60 // 198r2-M2 edited default table\n\treturn c")],
 "M3": [("cmd/wisp/run.go",
   "\tif _, frErr := ensureFirstRunConfig(s.dataDir, s.stderr); frErr != nil {\n\t\tfmt.Fprintf(s.stderr, \"wisp run: 首次配置建不出来（%v）：本轮仍按缺配置响亮失败，不会拿半份配置顶上\\n\", frErr)\n\t}",
   "\t_, _ = ensureFirstRunConfig(s.dataDir, s.stderr) // 198r2-M3 swallowed (198-v1 V4)")],
 "M4": [("cmd/wisp/firstrun.go",
   "\tif _, err := secret.NewStore(dataDir); err != nil {\n\t\treturn false, fmt.Errorf(\"数据根不可建（%s）：%w\", dataDir, err)\n\t}",
   "\tif err := secret.NewStore(dataDir); err != nil {\n\t\t_ = err\n\t\treturn true, nil // 198r2-M4 claimed success while creating nothing\n\t}")],
 "M6": [("cmd/wisp/firstrun.go",
   "wisp providers discover 与 probe 读这份文件去问真实端点",
   "wisp catalog sync 读这份文件去问真实端点 // 198r2-M6 invented entry point")],
 "M7": [("cmd/wisp/run.go",
   "\tif _, frErr := ensureFirstRunConfig(s.dataDir, s.stderr); frErr != nil {\n\t\tfmt.Fprintf(s.stderr, \"wisp run: 首次配置建不出来（%v）：本轮仍按缺配置响亮失败，不会拿半份配置顶上\\n\", frErr)\n\t}",
   "\t// 198r2-M7 the J1 call moved out of the run entry"),
  ("cmd/wisp/run.go",
   "\trt := &agentRuntime{spec: s, stdout: s.stdout, stderr: s.stderr, notify: s.notify}",
   "\tif _, frErr := ensureFirstRunConfig(s.dataDir, s.stderr); frErr != nil { // 198r2-M7 inside the shared root\n\t\tfmt.Fprintf(s.stderr, \"wisp run: 首次配置建不出来（%v）：本轮仍按缺配置响亮失败，不会拿半份配置顶上\\n\", frErr)\n\t}\n\trt := &agentRuntime{spec: s, stdout: s.stdout, stderr: s.stderr, notify: s.notify}")],
 "M8": [("cmd/wisp/firstrun.go",
   "\tif _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist) {",
   "\tif _, err := os.Stat(cfgPath); err == nil { // 198r2-M8 the V2b shape, any stat failure creates")],
}
if tag == "M5":
    p = pathlib.Path("cmd/wisp/firstrun.go")
    t = p.read_text(encoding="utf-8")
    start = t.index('\tfmt.Fprint(stderr,\n\t\t"wisp run: 缺的两样各有各的入口')
    end = t.index("\treturn true, nil", start)
    p.write_text(t[:start] + "\t// 198r2-M5 the where-to-fix-it guidance was deleted\n" + t[end:], encoding="utf-8")
    print("applied M5")
    sys.exit(0)
for path, old, new in edits[tag]:
    p = pathlib.Path(path)
    text = p.read_text(encoding="utf-8")
    n = text.count(old)
    if n != 1:
        sys.exit("MUTATION %s NOT APPLIED on %s: anchor found %d times (want 1)" % (tag, path, n))
    p.write_text(text.replace(old, new), encoding="utf-8")
    print("applied %s -> %s" % (tag, path))
PY
  if [ $? -ne 0 ]; then
    echo "MUTATION $1 NOT APPLIED - no reading taken"
    return
  fi
  go test ./cmd/wisp -run '^TestTicket198' -count=1 -v > "$PROBE/mut-$1.txt" 2>&1
  echo "GOTEST_RC=$? (full output: $PROBE/mut-$1.txt)"
  echo "top-level PASS=$(grep -cE '^--- PASS' "$PROBE/mut-$1.txt") FAIL=$(grep -cE '^--- FAIL' "$PROBE/mut-$1.txt") SKIP=$(grep -cE '^--- SKIP' "$PROBE/mut-$1.txt")"
  grep -E '^--- FAIL' "$PROBE/mut-$1.txt"
  grep -E '^(ok|FAIL)\s+github' "$PROBE/mut-$1.txt" | head -2
  restore_all
}

for tag in M1 M2 M3 M4 M5 M6 M7 M8; do apply_and_run "$tag"; done

echo "==================== FINAL STATE ===================="
git status --porcelain -- cmd internal
echo "porcelain-lines=$(git status --porcelain -- cmd internal | wc -l)"
