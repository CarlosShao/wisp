#!/bin/sh
# 【故障件·只建不删】本版把仓路径写成 MSYS 形 `/d/work/workspace/projects plans/Wisp`，
# 本机 sh 的 git 拿它做 `git -C` 当场 `fatal: cannot change to ... No such file or directory`
# （rc=128，零读数）＝**台面故障、⛔ 计入任何格的判据**。
# 有效版＝同目录 host-winsec.ps1（原生 PowerShell，⛔ MSYS 路径形）。本件留底供复现坑。
# ---------------------------------------------------------------------------
# 322-r2 本机 internal/winsec 取证：成对两枚临时根台面（D: 预建 GOTMPDIR ⇄ 默认 C 卷 TMP），
# 各自整包 -count=1 -v 一发；定向 -run 补 estate 色（已知本机格：整包有真窗区整进程摆动，
# estate 色取定向 -run，⛔ 拿整包当"全绿"）；Windows witnessTruncated 复认。
# 机主 2026-10-18 10:1x 新铁律：落点唯一允许 D:\tmp\wisp\20261018\322r2\*，盘符根零污染。
# ⛔ 在共享工作树跑 go test —— 全部在仓外导出树 $TREE 里。
set -e
REPO='/d/work/workspace/projects plans/Wisp'
ROOT=/d/tmp/wisp/20261018/322r2
TREE=$(cat "$ROOT/tree-head.txt")
GT="$ROOT/gtD"
LOGS="$REPO/.scratch/wisp/probes/322/r2/logs"
mkdir -p "$LOGS" "$GT"
HEAD=$(git -C "$REPO" rev-parse HEAD)
SUM="$LOGS/host-winsec-summary.txt"
: >"$SUM"
say() { echo "$@" | tee -a "$SUM"; }
say "### 322-r2 host winsec started $(date '+%F %T %z') HEAD=$HEAD"
say "TREE=$TREE  GT-D=$GT  (D: 卷·预建·8.3 别名)"
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
say "PATH-HEAD=$(echo "$PATH" | cut -c1-160)"

four() { # <bodylog>  —— 口径具名：RUN 含子测试，PASS/FAIL/SKIP 只数顶层
  b=$1
  printf 'RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s\n' \
    "$(grep -c '^=== RUN' "$b" || true)" "$(grep -c '^--- PASS' "$b" || true)" \
    "$(grep -c '^--- FAIL' "$b" || true)" "$(grep -c '^--- SKIP' "$b" || true)"
}
toolwords() { # <bodylog>  —— 任一非 0 且句子像台面 ⇒ 判台面故障，⛔ 计入判据
  b=$1
  printf 'walk=%s buildfailed=%s noteststorun=%s setupfailed=%s exit0xc0000135=%s\n' \
    "$(grep -ci walk "$b" || true)" "$(grep -ci 'build failed' "$b" || true)" \
    "$(grep -ci 'no tests to run' "$b" || true)" "$(grep -ci 'setup failed' "$b" || true)" \
    "$(grep -c '0xc0000135' "$b" || true)"
}

host_run() { # <label> <tmpdir-or-NONE> <filter-or-empty>
  label=$1; tmp=$2; filt=$3; rc=0
  body="$LOGS/$label.log"
  ( cd "$TREE"
    if [ "$tmp" != NONE ]; then export GOTMPDIR="$tmp"; mkdir -p "$tmp"; else unset GOTMPDIR; fi
    if [ -n "$filt" ]; then
      go test ./internal/winsec/ -count=1 -timeout 30m -v -run "$filt" >"$body" 2>&1
    else
      go test ./internal/winsec/ -count=1 -timeout 30m -v >"$body" 2>&1
    fi ) || rc=$?
  { echo "=== $label rc=$rc GOTMPDIR=$tmp filter='$filt'"
    four "$body"; toolwords "$body"
    echo "--- top FAIL names:"; grep -E '^--- FAIL' "$body" || echo "(none)"
    echo "--- top SKIP names:"; grep -E '^--- SKIP' "$body" || echo "(none)"
    echo "--- witnessTruncated rows:"; grep -h 'witnessTruncated' "$body" || echo "(no witnessTruncated row)"
  } | tee -a "$SUM"
}

PAT_ESTATE='TestGateJudgesThePrivateSetByResolvedSID|TestSealNarrowsAndNamesThePrincipalItRemovedBySID'
case "${STEP:-all}" in
  a)  host_run wh-D "$GT" '' ;;
  b)  host_run wh-C NONE '' ;;
  c)  host_run wd-D "$GT" "$PAT_ESTATE" ;;
  d)  host_run wd-C NONE "$PAT_ESTATE" ;;
  e)  host_run wt-D "$GT" 'TestAC3Windows' ;;
  all) host_run wh-D "$GT" ''; host_run wh-C NONE ''
       host_run wd-D "$GT" "$PAT_ESTATE"
       host_run wd-C NONE "$PAT_ESTATE"
       host_run wt-D "$GT" 'TestAC3Windows' ;;
esac
say "### done $(date '+%F %T %z')"
