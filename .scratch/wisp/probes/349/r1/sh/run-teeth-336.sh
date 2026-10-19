#!/bin/sh
# 349-r1 复跑票 336 那 8 枚牙：一字⛔ 动台件（.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh，
# 我已 git cat-file -e HEAD:<path> 现验存在＝rc 0），env 指**我自己 git archive 重导的树**
# （treeDel＝交件锚 db2f9eab 的黑名单 blob 964ee75a），⛔ 复用别人的树、⛔ -overlay。
# 每发落刀前后各验锚命中（台件自带 落点命中枚数 before/after）与 pristine 哈希（台件自带 check_hash），
# 还原哈希逐发摊出来。
# ★ANCHOR 用我交件锚（⛔ 用 HEAD：HEAD 之后还会随我落笔漂，而树是 db2f9eab 导的）。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
S=.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh
echo "== step 0: 台件存在性＋一字未动的凭据 =="
git cat-file -e "HEAD:$S" && echo "cat-file -e rc=0 EXISTS"
printf '台件 blob at HEAD = %s\n' "$(git rev-parse "HEAD:$S")"
printf '台件 blob at 336-v1 交付锚 = %s\n' "$(git rev-parse 'c315fa4c:.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh' 2>/dev/null || echo NA)"
echo "== step 1: 树与锚的 blob 对撞（treeDel 里的 blacklist.go 必须 == ANCHOR 的 blob） =="
A=db2f9eab7bcffd827240763699aa0f4ad7ff3266
printf 'ANCHOR blob=%s\n' "$(git rev-parse "$A:internal/risk/blacklist.go")"
printf 'treeDel blob=%s\n' "$(cd /d/tmp/wisp/20261019/349r1/treeDel && git hash-object internal/risk/blacklist.go)"
echo "== step 2: 跑台件（串行、一发到底） =="
ANCHOR="$A" \
TREE='D:/tmp/wisp/20261019/349r1/treeDel' \
LOGDIR="$REPO/.scratch/wisp/probes/349/r1/logs/teeth" \
GOTMPROOT='/d/tmp/wisp/20261019/349r1/gt' \
  sh "$S"
echo "SWEEP_RC=$?"
