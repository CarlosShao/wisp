#!/bin/sh
# 编排者逐字重放 120-v2 的 tooth.sh —— ★第二发（第一发是我自己的台面错，见下）
# ★第一发为什么⛔ 算读数：我把树建在 mktemp -d＝/tmp/tmp.4xaEwDqEu5/head，
#   Docker Desktop 把容器外那个 `/tmp/…' 读成它自己 VM 里的路径、⛔ Windows 那枚 ⇒
#   挂载指到空目录、`go test` 连 go.mod 都看不见 ⇒ rc=1 且⛔ 落 `$tag.log`。
#   ⇒ 三发去牙的**突变与还原**部分真跑成了（LANDS 行＋RESTORE MATCHES×3 是我的读数），
#     四数那几行⛔ 是读数。定式＝**凡 rc 非 0 又拿不到日志，先证日志该出现的那枚目录真被挂上**。
# ★本发改回 `/d/tmp/…'（腿原树同一条卷＝Docker 认得 `/d/'）。突变逻辑仍⛔ 动一行。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
D=/d/tmp/orch-120v2-replay2
V=$D/head
LOG=$D/driver.log
mkdir -p "$V" "$D"
( cd "$REPO" && git archive HEAD go.mod go.sum internal cmd tools scripts frontend ) | tar -x -C "$V"
NGO=$(find "$V" -name '*.go' | wc -l)
{
  echo "tree=$V go_files=$NGO"
  echo "winsec_other.go sha256=$(sha256sum "$V/internal/winsec/winsec_other.go" | cut -d' ' -f1)"
} > "$LOG"

sed -e "s#^V=/d/tmp/probe120v2/head#V=$V#" \
    -e "s#^LOG=/d/tmp/probe120v2/v2-tooth-driver.log#LOG=$LOG#" \
    -e "s#BAK=\"\\\$F.bak-120v2-teeth-on\"#BAK=\"$D/winsec_other.go.bak-orch\"#" \
    -e "s#probe120v2build#orch120v2b#" \
    "$REPO/.scratch/wisp/probes/120/v2/tooth.sh" > "$D/tooth-orch.sh"
echo "diff-lines-vs-leg-script=$(diff "$REPO/.scratch/wisp/probes/120/v2/tooth.sh" "$D/tooth-orch.sh" | grep -c '^[<>]')" >> "$LOG"
sh "$D/tooth-orch.sh" >> "$LOG" 2>&1
echo "TOOTH_EXIT=$?" >> "$LOG"

# ★两发改的只有 -run／-count 那一段（同卷、同镜像、同 mod 缓存、同 GOPROXY=off）
docker_extra() {
  tag="$1"; shift
  rc=$(MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off \
      -v /d/work/base/gopath/pkg/mod:/gomod -e GOMODCACHE=/gomod \
      -v "$V":/work -v orch120v2b:/root/.cache/go-build -w /work golang:1.27 \
      sh -c "go test $* -count=1 -v ./internal/winsec/ > /work/$tag.log 2>&1; echo \$?" 2>>"$LOG")
  echo "$tag rc=$rc loglines=$(wc -l < "$V/$tag.log" 2>/dev/null || echo NOLOG)" >> "$LOG"
  echo "$tag RUN=$(grep -c '^=== RUN' "$V/$tag.log") TOPPASS=$(grep '^--- PASS' "$V/$tag.log" | awk '{print $3}' | sort -u | wc -l) TOPFAIL=$(grep -c '^--- FAIL' "$V/$tag.log") TOPSKIP=$(grep -c '^--- SKIP' "$V/$tag.log")" >> "$LOG"
  grep '^--- FAIL' "$V/$tag.log" | sed "s/^/$tag FAILED /" >> "$LOG"
}
# AC#3 那把尺：容器 -run 'TestAC[1-4]POSIX' -count=1 -v ⇒ 顶层枚数
docker_extra orch-ac3 "-run" "'TestAC[1-4]POSIX'"
# AC#5 那把尺：整包 -count=2 -v ⇒ 四数
rc=$(MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off \
    -v /d/work/base/gopath/pkg/mod:/gomod -e GOMODCACHE=/gomod \
    -v "$V":/work -v orch120v2b:/root/.cache/go-build -w /work golang:1.27 \
    sh -c "go test -count=2 -v ./internal/winsec/ > /work/orch-ac5.log 2>&1; echo \$?" 2>>"$LOG")
echo "orch-ac5 rc=$rc loglines=$(wc -l < "$V/orch-ac5.log" 2>/dev/null || echo NOLOG)" >> "$LOG"
echo "orch-ac5 RUN=$(grep -c '^=== RUN' "$V/orch-ac5.log") TOPPASS=$(grep -c '^--- PASS' "$V/orch-ac5.log") TOPFAIL=$(grep -c '^--- FAIL' "$V/orch-ac5.log") TOPSKIP=$(grep -c '^--- SKIP' "$V/orch-ac5.log")" >> "$LOG"
grep '^--- FAIL' "$V/orch-ac5.log" | sed 's/^/orch-ac5 FAILED /' >> "$LOG"
echo '===== driver tail ====='
tail -30 "$LOG"
echo "LOGPATH=$LOG"
