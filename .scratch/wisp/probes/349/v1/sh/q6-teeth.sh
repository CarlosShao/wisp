#!/bin/sh
# 349-v1 复打票 336 那 8 枚牙：台件一字⛔ 动（.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh），
# env 指我自己 git archive 重导的 treeDel（交件锚 ed8db1ea，blacklist.go blob 964ee75a）。
# 起手先现验：台件存在性＋台件 blob == 336-v1 交付锚的 blob（一字未动的凭据）。
# GOTMPROOT 用 POSIX 形（台件内部 mkdir 需要），但这八枚牙件不建 junction／mklink＝
# 台损那形（POSIX GOTMPDIR 触发 mklink Invalid switch）打不到它们；本腿仍按"两形各验"
# 摊出：本轮＝POSIX 形 GOTMP（与 349-r1 同发对照），另跑一发干净形（Windows 形）做对撞。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
export PATH="/d/work/base/go/bin:$PATH"
export GOPATH='D:\work\base\gopath'
S=.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh
echo "== step 0: 台件存在性＋未动凭据 =="
git cat-file -e "HEAD:$S" && echo "cat-file -e rc=0 EXISTS"
printf '台件 blob at HEAD      = %s\n' "$(git rev-parse "HEAD:$S")"
printf '台件 blob at c315fa4c  = %s\n' "$(git rev-parse 'c315fa4c:.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh' 2>/dev/null || echo NA)"
echo "== step 1: 锚⇄树 blob 对撞 =="
A=ed8db1ea27eb50d3057dd0cd7d5c955a13e978fc
printf 'ANCHOR(ed8db1ea) blacklist.go blob=%s\n' "$(git rev-parse "$A:internal/risk/blacklist.go")"
printf 'treeDel blacklist.go hash-object=%s\n' "$(cd /d/tmp/wisp/20261019/349v1/treeDel && git hash-object internal/risk/blacklist.go)"
echo "== step 2: 跑台件（POSIX 形 GOTMP） =="
ANCHOR="$A" \
TREE='D:/tmp/wisp/20261019/349v1/treeDel' \
LOGDIR='/d/tmp/wisp/20261019/349v1/logsTeeth/posix' \
GOTMPROOT='/d/tmp/wisp/20261019/349v1/gtTeethP' \
  sh "$S"
echo "SWEEP_RC=$?"
