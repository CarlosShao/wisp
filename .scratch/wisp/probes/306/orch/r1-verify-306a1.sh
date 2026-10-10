#!/usr/bin/env bash
# 306 编排者复跑：收 306-a1 之前的独立取证（只读产码面；突变那一发在仓外导出树）
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT_DIR=".scratch/wisp/probes/306/orch/logs"
mkdir -p "$OUT_DIR"
STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="$OUT_DIR/r1-orch-verify-306a1-$STAMP.txt"

exec > "$OUT" 2>&1

echo "# 306 编排者独立复跑（收 306-a1 之前）  时刻=$(date '+%Y-%m-%d %H:%M:%S %z')"
echo
echo "## 0. 台面"
echo "HEAD_sha=$(git rev-parse --short HEAD)"
echo "HEAD_full=$(git rev-parse HEAD)"
git status --porcelain -- internal cmd docs scripts .github > /tmp/306o-porcelain.txt 2>&1
echo "rc_porcelain_scoped=$?"
echo "porcelain_internal_cmd_docs_scripts_github_lines=$(wc -l < /tmp/306o-porcelain.txt)"
git show HEAD:internal/audio/wavinjector.go > /tmp/306o-blob.go 2>&1
cmp /tmp/306o-blob.go internal/audio/wavinjector.go
echo "rc_blob_vs_worktree_cmp=$?"
echo
echo "## 1. 表②那两支的调用链（我现跑）"
grep -rn --include=*.go "NewWavInjector(" internal cmd tools > /tmp/306o-newinj.txt 2>&1
echo "rc_grep_newinj=$?"
echo "newinj_all_lines=$(wc -l < /tmp/306o-newinj.txt)"
echo "newinj_nondeff_lines=$(grep -vc 'func NewWavInjector' /tmp/306o-newinj.txt)"
echo "newinj_nontest_nonself_lines=$(grep -v '_test.go' /tmp/306o-newinj.txt | grep -vc '^internal/audio/wavinjector.go')"
echo "--- 逐字（非 _test.go 命中）---"
grep -rn --include=*.go "NewWavInjector\|WavInjector" internal cmd tools | grep -v _test.go | grep -v '^internal/audio/wavinjector.go'
echo "rc_grep_comments=$?"
echo
echo "## 1b. writeWav 枚数两把尺并排（腿报 11＝含定义行）"
grep -rn --include=*.go "writeWav(" internal cmd tools > /tmp/306o-wav.txt 2>&1
echo "writewiv_incl_def_lines=$(wc -l < /tmp/306o-wav.txt)"
echo "writewiv_callsites_excl_def=$(grep -vc 'func writeWav' /tmp/306o-wav.txt)"
echo "--- 逐字调用点（剥定义行）---"
grep -v "func writeWav" /tmp/306o-wav.txt
echo "--- 按文件计数 ---"
grep -v "func writeWav" /tmp/306o-wav.txt | cut -d: -f1 | sort | uniq -c
echo
echo "## 1c. 旁路入口与夹具硬编码（内容锚逐字）"
sed -n '581p' internal/audio/hotplug_test.go
echo "rc_sed_hotplug581=$?"
sed -n '26p;27p;32p' internal/audio/wavinjector_test.go
echo "rc_sed_writewav_lines=$?"
echo "--- parseWav doc comment 实际跨 162-164（腿引 163-164）---"
sed -n '162,164p' internal/audio/wavinjector.go
echo
echo "## 2. 覆盖尺我自己复跑（go test ./internal/audio/ -coverprofile）"
COV="/tmp/306o-cover-$STAMP.txt"
GOFLAGS=-mod=mod go test ./internal/audio/ -count=1 -coverprofile="$COV" > /tmp/306o-test.txt 2>&1
echo "rc_test=$?"
tail -2 /tmp/306o-test.txt
go tool cover -func="$COV" > /tmp/306o-func.txt 2>&1
echo "rc_coverfunc=$?"
grep -E "wavinjector.go:165: parseWav|wavinjector.go:56: NewWavInjector|wasapi_windows.go:376: Drain|wasapi_windows.go:409: convertPacket|wasapi_windows.go:200: parseWaveFormat|total:" /tmp/306o-func.txt
echo "rc_grep_func=$?"
echo "--- wavinjector.go 行段块计数（188-226）---"
grep -E "wavinjector\.go:(18[89]|19[0-9]|2[0-2][0-9])\." "$COV"
echo "rc_grep_blocks=$?"
cp "$COV" "$OUT_DIR/r1-orch-cover-$STAMP.txt"
echo "saved_cover_copy_rc=$?"
echo
echo "## 3. 戊类同形异族 6 枚具名排除（我逐枚 sed）"
echo "[1] internal/tools/recycle_windows.go:37"; sed -n '37p' internal/tools/recycle_windows.go
echo "[2] cmd/wisp/panel_host_windows_test.go:181"; sed -n '181p' cmd/wisp/panel_host_windows_test.go
echo "[3] cmd/wisp/notify_windows.go:27"; sed -n '27p' cmd/wisp/notify_windows.go
echo "[4] cmd/wisp/notify_windows.go:30"; sed -n '30p' cmd/wisp/notify_windows.go
echo "[5] cmd/wisp/notify_windows.go:35"; sed -n '35p' cmd/wisp/notify_windows.go
echo "rc_sed_exclusions=$?"
echo "--- 仓里 0x0003/0x00000001 全部命中枚数（不加族尺=虚报的源头）---"
grep -rn --include=*.go -E "0x0003|0x00000001" internal cmd tools > /tmp/306o-0003.txt 2>&1
echo "rc_grep_0003=$?"
echo "lines=$(wc -l < /tmp/306o-0003.txt)"
cat /tmp/306o-0003.txt
echo
echo "## 4. 甲/乙两形的成本尺（我现跑）"
git ls-files > /tmp/306o-lsfiles.txt 2>&1
echo "rc_lsfiles=$?"
echo "tracked_wav_count=$(grep -icE '\.wav$' /tmp/306o-lsfiles.txt)"
find internal cmd -type d -name testdata > /tmp/306o-testdata.txt 2>&1
echo "rc_find_testdata=$?"
cat /tmp/306o-testdata.txt
grep -nE "wav" .gitignore > /tmp/306o-ignore.txt 2>&1
echo "rc_grep_gitignore_wav=$? (1=零命中=不会被静默跳过)"
wc -l < /tmp/306o-ignore.txt
echo
echo "## 5. CI 可见性：core 与 windows 两档都含 audio（我现跑）"
sed -n '242p;253p' scripts/portable-tests.sh
echo "rc_sed_scope=$?"
sed -n '164p;168p;193p;195p' scripts/portable-tests.sh
echo "rc_sed_pins=$?"
sed -n '419,420p;487p;533,534p;797p' .github/workflows/ci.yml
echo "rc_sed_ci=$?"
echo
echo "## 6. 硬前提决定性两发（平台中立件不能引那三枚常量）"
EXPORT="/c/Users/swq/tmp/306o-refconst-$STAMP"
mkdir -p "$EXPORT"
git archive HEAD | tar -x -C "$EXPORT"
echo "rc_archive_export=$?"
cd "$EXPORT" || exit 8
sed -i 's/if fmtTag == 0xFFFE {/if fmtTag == waveFormatExt {/; s/case fmtTag == 3 \&\& bits == 32:/case fmtTag == waveFormatFloat \&\& bits == 32:/' internal/audio/wavinjector.go
echo "rc_sed_mutate=$?"
echo "--- 改后那两行逐字 ---"
sed -n '192p;218p' internal/audio/wavinjector.go
GOOS=linux GOFLAGS= go build ./internal/audio/ > /tmp/306o-mut-linux.txt 2>&1
echo "rc_linux_mutated=$?"
cat /tmp/306o-mut-linux.txt
GOOS=windows GOFLAGS= go build ./internal/audio/ > /tmp/306o-mut-win.txt 2>&1
echo "rc_windows_mutated=$?"
cat /tmp/306o-mut-win.txt
cd "D:/work/workspace/projects plans/Wisp" || exit 9
GOOS=linux GOFLAGS= go build ./internal/audio/ > /tmp/306o-pristine-linux.txt 2>&1
echo "rc_linux_build_pristine_worktree=$?"
cat /tmp/306o-pristine-linux.txt
git show HEAD:internal/audio/wavinjector.go > /tmp/306o-blob2.go 2>&1
cmp /tmp/306o-blob2.go internal/audio/wavinjector.go
echo "rc_final_blob_cmp=$? (0=产码面一字未动)"
echo
echo "## 7. 越界自查"
git status --porcelain -- internal cmd tools docs scripts .github > /tmp/306o-final-porcelain.txt 2>&1
echo "rc_final_porcelain=$?"
echo "final_porcelain_lines=$(wc -l < /tmp/306o-final-porcelain.txt)"
echo "new_go_files_in_probes=$(find .scratch/wisp/probes/306 -name '*.go' | wc -l)"
echo "out_files_in_probes=$(find .scratch/wisp/probes/306 -name '*.out' | wc -l)"
echo "ticket_box_lines_changed=$(git status --porcelain -- .scratch/wisp/issues | wc -l)"
