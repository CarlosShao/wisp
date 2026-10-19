#!/bin/sh
# 342-v2 必答(ii) 决定性对照＝【拔掉产码那一只手】的仪器空转发。
# 问题：落地侧那 1~2/400 的 truncated 漏出，是【产码真的跟随链接截断了外来对象】，
#       还是【仪器自己的采样/复位串扰】（racer 的 rename 序列与 restock/readWitness 在同一枚名字上抢）？
# 做法：把 `w.run(..., func(){ return PrivateFile(...) })` 换成 `func(){ return nil }`
#       ＝racer／readWitness／restockWitness 全照跑、产码一次都⛔ 被调用。
#       残余若照旧 >0 ⇒ 台损（仪器形状）；残余归 0 ⇒ 那几枚是产码造成的＝真漏、⛔ 归因到台。
# ★只打仓外拷贝树 mut-X-* 的物理字节；落刀前验旧锚命中＝期望、落刀后旧锚 0＋新锚≥1；跑完当场还原、报还原哈希。
set -u
ROOTW='D:\tmp\wisp\20261019\342v2'
ROOT=/d/tmp/wisp/20261019/342v2
LOGS="$ROOT/logs"
T=internal/winsec/ancestoropen_342_windows_test.go
ANCHOR='return PrivateFile(w.name, []byte(winAncPayload342r1), 0o600)'
for SHOT in 1 2; do
  MUT="mut-X-open-treeL-$SHOT"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/treeL" "$ROOT/$MUT"
  PRISTINE=$(sha256sum "$ROOT/treeL/$T" | cut -c1-16)
  BEFORE=$(sha256sum "$ROOT/$MUT/$T" | cut -c1-16)
  echo "== $MUT PRISTINE_SHA16=$PRISTINE BEFORE=$BEFORE MATCH=$([ "$PRISTINE" = "$BEFORE" ] && echo YES || echo NO)"
  echo "OLD_ANCHOR_HITS_BEFORE=$(grep -c "$ANCHOR" "$ROOT/$MUT/$T")"
  perl -0pi -e "s/\t\t\treturn PrivateFile\(w\.name, \[\]byte\(winAncPayload342r1\), 0o600\)/\t\t\treturn nil \/\/ MUTATION-X: the production hand is unplugged/" "$ROOT/$MUT/$T"
  echo "OLD_ANCHOR_HITS_AFTER=$(grep -c "$ANCHOR" "$ROOT/$MUT/$T" || true)"
  echo "NEW_ANCHOR_HITS_AFTER=$(grep -c 'MUTATION-X: the production hand is unplugged' "$ROOT/$MUT/$T")"
  LABEL="ii03-${MUT}"
  DLLP=$(cygpath -w "$ROOTW\\$MUT\\build"); SHP=$(cygpath -w "$ROOTW\\$MUT\\third_party\\sherpa-onnx")
  ( cd "$ROOTW\\$MUT" && PATH="$SHP;$DLLP;$PATH" GOCACHE='D:\tmp\wisp\20261019\342v2\gocache-win' \
      go test ./internal/winsec/ -v -count=1 -timeout 20m -run 'TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1' ) > "$LOGS/$LABEL.raw.txt" 2>&1
  echo "RC=$?"
  grep -E 'form=[a-z-]+ attempts=400' "$LOGS/$LABEL.raw.txt" | sed -E 's/^ +[a-z0-9_]+\.go:[0-9]+: //'
  grep -E 'RED:' "$LOGS/$LABEL.raw.txt" | sed -E 's/^ +[a-z0-9_]+\.go:[0-9]+: //' | head -4
  grep -E '^--- (FAIL|PASS)' "$LOGS/$LABEL.raw.txt" || true
  cp -f "$ROOT/treeL/$T" "$ROOT/$MUT/$T"
  echo "RESTORE_SHA16=$(sha256sum "$ROOT/$MUT/$T" | cut -c1-16) PRISTINE_MATCH_AFTER=$([ "$PRISTINE" = "$(sha256sum "$ROOT/$MUT/$T" | cut -c1-16)" ] && echo YES || echo NO)"
  echo "-----"
done
