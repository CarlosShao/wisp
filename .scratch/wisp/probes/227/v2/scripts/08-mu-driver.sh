#!/bin/bash
# AC#6 去牙串行台：三发落点各不同，每发＝落刀验锚→跑定向 -run→当场还原→报还原哈希。
# 只打仓外导出树 D:\tmp\wisp\20261019\227v2\deliv-tree-ac5（⛔ 跟踪文件、⛔ -overlay）。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
SRC='/d/tmp/wisp/20261019/227v2/deliv-tree-ac5'
MUT='/d/tmp/wisp/20261019/227v2/mut'
CR="$SRC/cmd/wisp/config_reload.go"
mkdir -p "$MUT"
cp -p "$CR" "$MUT/config_reload.go.pristine"
PRISTINE_SHA=$(sha256sum "$CR" | cut -d' ' -f1)
echo "PRISTINE_SHA=$PRISTINE_SHA"

run_case() {
  local pattern="$1" tag="$2"
  pwsh.exe -NoProfile -File 'D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\227\v2\scripts\02-run.ps1' \
    -Gotmp "C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp\\227v2-mu" \
    -Log "D:\\tmp\\wisp\\20261019\\227v2\\mut-$tag.log" \
    -Tree 'D:\tmp\wisp\20261019\227v2\deliv-tree-ac5' \
    -Sel "./cmd/wisp -run=$pattern"
  echo "== $tag four numbers =="
  bash "$REPO/.scratch/wisp/probes/227/v2/scripts/03-four-numbers.sh" "/d/tmp/wisp/20261019/227v2/mut-$tag.log"
}

restore_now() {
  cp -p "$MUT/config_reload.go.pristine" "$CR"
  R=$(sha256sum "$CR" | cut -d' ' -f1)
  echo "RESTORED_SHA=$R"; echo "RESTORE_EQUALS_PRISTINE=$([ "$R" = "$PRISTINE_SHA" ] && echo yes || echo NO)"
}

case "$1" in
  mu1) # 摘掉 Spawn 那一行＝tick 根本不起
    echo "pre_old_anchor_hits=$(grep -c 'rt.reloadHandle = observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)' "$CR")"
    perl -0pi -e 's/\trt\.reloadHandle = observe\.Default\.Spawn\("watchdog", "config", rt\.reloadRoot, rt\.configReloadTick\)\n//' "$CR"
    echo "post_old_anchor_hits=$(grep -c 'rt.reloadHandle = observe.Default.Spawn' "$CR")"
    echo "mutated_sha=$(sha256sum "$CR" | cut -d' ' -f1)"
    run_case 'TestTicket223RunArmsTheReloadTick|TestAC1AlwaysBranchDoesNotRevertAHandEditedKey' 'mu1'
    restore_now
    ;;
  mu2) # 摘掉 reloadOnce 的 CheckAndReload 调用＝tick 醒着、从不重读
    echo "pre_old_anchor_hits=$(grep -c 'rep, err := rt.mgr.CheckAndReload()' "$CR")"
    perl -0pi -e 's/\trep, err := rt\.mgr\.CheckAndReload\(\)/\tvar rep *config.Report\n\tvar err error\n\t_, _ = rep, err \/\/ MUB 227-v2 seam removed/' "$CR"
    echo "post_old_anchor_hits=$(grep -c 'rt.mgr.CheckAndReload()' "$CR")"
    echo "post_new_anchor_hits=$(grep -c 'MUB 227-v2 seam removed' "$CR")"
    echo "mutated_sha=$(sha256sum "$CR" | cut -d' ' -f1)"
    run_case 'TestTicket223RunArmsTheReloadTick|TestAC1AlwaysBranchDoesNotRevertAHandEditedKey' 'mu2'
    restore_now
    ;;
  mu3) # 把 1s tick 改成 10ms＝让轮询真赶上手改与收口之间那一窗
    echo "pre_old_anchor_hits=$(grep -c 'configReloadPollInterval = time.Second' "$CR")"
    perl -0pi -e 's/configReloadPollInterval = time\.Second/configReloadPollInterval = 10 * time.Millisecond \/\/ MUC 227-v2/' "$CR"
    echo "post_old_anchor_hits=$(grep -c 'configReloadPollInterval = time.Second' "$CR")"
    echo "post_new_anchor_hits=$(grep -c 'MUC 227-v2' "$CR")"
    echo "mutated_sha=$(sha256sum "$CR" | cut -d' ' -f1)"
    run_case 'TestAC1AlwaysBranchDoesNotRevertAHandEditedKey' 'mu3'
    restore_now
    ;;
  control) # 未突变对照：同一 -run 形必须先绿
    run_case 'TestTicket223RunArmsTheReloadTick|TestAC1AlwaysBranchDoesNotRevertAHandEditedKey' 'control'
    ;;
esac
echo "final_restored_sha=$(sha256sum "$CR" | cut -d' ' -f1)"
