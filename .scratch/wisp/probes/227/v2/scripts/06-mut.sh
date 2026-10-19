#!/bin/bash
# AC#6 去牙台（★只打在仓外导出树的物理字节上，⛔ 改跟踪文件、⛔ -overlay）。
# 牙的定义＝cmd/wisp/config_reload_223_test.go:333 TestTicket223RunArmsTheReloadTick
#   「真轮询下手改仍看得见」那枚常驻用例。去牙＝摘掉生产侧的轮询接缝，两发落点不同：
#   MU-1＝摘 Spawn 那一行（tick 根本不起）；MU-2＝摘 reloadOnce 里的 CheckAndReload 调用（tick 醒着但不重读）。
#   ⇒ 每一发该用例都必须红；若仍绿＝那枚用例⛔ 有牙。
# 尺＝内容锚（函数体内那枚 token），落刀前验命中数＝期望值、落刀后再验（新锚命中＋旧锚归零）。
set -u
SRC=/d/tmp/wisp/20261019/227v2/deliv-tree-ac5
WORK=/d/tmp/wisp/20261019/227v2/mut
mkdir -p "$WORK"
case "$1" in
  pristine)
    sha256sum "$SRC/cmd/wisp/config_reload.go" | sed 's# .*# #'
    echo "PRISTINE_COPY=$WORK/config_reload.go.pristine"
    cp -p "$SRC/cmd/wisp/config_reload.go" "$WORK/config_reload.go.pristine"
    ;;
  mu1)   # 摘掉 Spawn 那一行：tick 不起
    F="$SRC/cmd/wisp/config_reload.go"
    echo "before_sha=$(sha256sum "$F" | cut -d' ' -f1)"
    echo "pre_anchor_hits=$(grep -c 'rt.reloadHandle = observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)' "$F")"
    cp -p "$WORK/config_reload.go.pristine" "$F"
    perl -0pi -e 's/\trt\.reloadHandle = observe\.Default\.Spawn\("watchdog", "config", rt\.reloadRoot, rt\.configReloadTick\)\n//s' "$F"
    echo "post_old_anchor_hits=$(grep -c 'rt.reloadHandle = observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)' "$F")"
    echo "post_spawn_token_lines=$(grep -c 'configReloadTick' "$F")"
    echo "after_sha=$(sha256sum "$F" | cut -d' ' -f1)"
    ;;
  mu2)   # 摘掉 CheckAndReload 那一调用：tick 醒着、不重读
    F="$SRC/cmd/wisp/config_reload.go"
    echo "before_sha=$(sha256sum "$F" | cut -d' ' -f1)"
    echo "pre_anchor_hits=$(grep -c 'rep, err := rt.mgr.CheckAndReload()' "$F")"
    cp -p "$WORK/config_reload.go.pristine" "$F"
    perl -0pi -e 's/\trep, err := rt\.mgr\.CheckAndReload\(\)/\tvar rep *config.Report\n\tvar err error\n\t_ = rep; _ = err \/\/ MUB: seam removed by 227-v2/s' "$F"
    echo "post_old_anchor_hits=$(grep -c 'rep, err := rt.mgr.CheckAndReload()' "$F")"
    echo "post_new_anchor_hits=$(grep -c 'MUB: seam removed by 227-v2' "$F")"
    echo "after_sha=$(sha256sum "$F" | cut -d' ' -f1)"
    ;;
  restore)
    F="$SRC/cmd/wisp/config_reload.go"
    cp -p "$WORK/config_reload.go.pristine" "$F"
    echo "restored_sha=$(sha256sum "$F" | cut -d' ' -f1)"
    echo "pristine_sha=$(sha256sum "$WORK/config_reload.go.pristine" | cut -d' ' -f1)"
    echo "head_blob=$(cd '/d/work/workspace/projects plans/Wisp' && git rev-parse HEAD:cmd/wisp/config_reload.go)"
    echo "restore_equals_pristine=$([ "$(sha256sum "$F" | cut -d' ' -f1)" = "$(sha256sum "$WORK/config_reload.go.pristine" | cut -d' ' -f1)" ] && echo yes || echo NO)"
    ;;
  *) echo "usage: $0 pristine|mu1|mu2|restore"; exit 9 ;;
esac
