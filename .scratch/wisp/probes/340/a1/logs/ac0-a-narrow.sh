#!/bin/sh
# 票 340 普查腿 340-a1 · AC#0 a) 真产叠改形用例枚数（票面窄尺，具名提交树）
# 可逐字重跑：sh .scratch/wisp/probes/340/a1/logs/ac0-a-narrow.sh
ANCHOR=feac8f77c2fbb15addbc36fb1bbaf13ed56588be
TICKET=f705b082f452754af477dff096e7000e590c5d82

echo "== ruler A1 (票面窄尺) git grep -nE '\\\"[^\\\"]*#[0-9]+\\\"' -- internal/agent/approval"
echo "-- at 票面锚 f705b082:"
git grep -nE '"[^"]*#[0-9]+"' "$TICKET" -- internal/agent/approval || echo "rc=$? (grep -c 0 命中会回 rc=1)"
echo "-- at 本腿起手锚 feac8f77:"
git grep -nE '"[^"]*#[0-9]+"' "$ANCHOR" -- internal/agent/approval || echo "rc=$? (同上)"

echo
echo "== 正控 1：同一把尺，射程收窄到 queue.go（那枚改写串自己就该命中）"
git grep -nE '"[^"]*#[0-9]+"' "$ANCHOR" -- internal/agent/approval/queue.go || echo "rc=$? (正控若也 0＝这把尺失明，⛔ 用它的 0)"

echo
echo "== 正控 2：同一把尺，全仓（不限包）的命中枚数（证明尺在别的包看得见东西）"
git grep -nE '"[^"]*#[0-9]+"' "$ANCHOR" | wc -l

echo
echo "== 正控 3：尺的形状对照（把 #[0-9] 放宽成任意 # ，approval 包内命中枚数）"
git grep -nE '"[^"]*#' "$ANCHOR" -- internal/agent/approval | wc -l
