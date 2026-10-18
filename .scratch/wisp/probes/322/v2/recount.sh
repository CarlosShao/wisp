#!/usr/bin/env bash
# 322-v2 复数尺：逐枚【指名那一个原始 log】，⛔ 通配（A1029 那枚 `cmdwisp-summary*` 假矛盾的定式）。
# 尺＝^=== RUN（含子测试）⇄ ^--- PASS（只顶格）⇄ ^--- FAIL ⇄ ^--- SKIP，⛔ 混用⛔ 相减。
# 台面故障尺＝0xc0000135／build failed／no tests to run／setup failed 逐发验零。
# 来源两目录：R＝322-r2 入库 log（blob 已逐枚对过 HEAD，MATCH）；V＝本腿 13:0x 自己现跑的新发。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 99
L=".scratch/wisp/probes/322/r2/logs"
V=".scratch/wisp/probes/322/v2/logs"
{
echo "### 322-v2 复数尺（本腿自己现跑，逐枚指名单文件，⛔ 通配） $(date '+%F %T %z')"
echo "### 起手 git log -1 = $(git log -1 --format='%H %ad' --date=iso-local)"
count_one() {
  local tag="$1" p="$2"
  [ -f "$p" ] || { printf "%-22s MISSING %s\n" "$tag" "$p"; return; }
  printf "%-22s RUN=%-5s PASS_TOP=%-5s FAIL_TOP=%-4s SKIP=%-4s | c0000135=%s buildfailed=%s noteststorun=%s setupfailed=%s | sha=%s\n" \
    "$tag" "$(grep -c '^=== RUN' "$p" || true)" "$(grep -c '^--- PASS' "$p" || true)" \
    "$(grep -c '^--- FAIL' "$p" || true)" "$(grep -c '^--- SKIP' "$p" || true)" \
    "$(grep -c '0xc0000135' "$p" || true)" "$(grep -c 'build failed' "$p" || true)" \
    "$(grep -c 'no tests to run' "$p" || true)" "$(grep -c 'setup failed' "$p" || true)" \
    "$(git hash-object "$p" | cut -c1-12)"
}
count_one full-D              "$L/cwisp-full-D.log"
count_one full-C              "$L/cwisp-full-C.log"
count_one estate-resident     "$L/cwisp-estate-resident.log"
count_one estate-deathpoints  "$L/cwisp-estate-deathpoints.log"
count_one estate-live-winelive "$L/cwisp-estate-live-winelive.log"
count_one estate-live-FAULT   "$L/cwisp-estate-live.log"
count_one cc2-head-r2        "$L/cc2-head.body.log"
count_one cc2-mine-r2        "$L/cc2-mine.body.log"
count_one cw2-head-r2        "$L/cw2-head.body.log"
count_one cw2-mine-r2        "$L/cw2-mine.body.log"
count_one orch-10-16-full     ".scratch/wisp/probes/orch/20261016/logs/orch-322v1-ac4-cmdwisp-full.txt"
count_one v2head-minefresh "$V/v2head.body.log"
count_one v2mine-minefresh "$V/v2mine.body.log"
count_one v2winsec-minefresh "$V/v2winsec.body.log"
count_one v2cons-minefresh "$V/v2cons.body.log"
echo "### done $(date '+%F %T %z')"
} | tee .scratch/wisp/probes/322/v2/logs/recount.txt
