#!/bin/sh
# 票 338 落地腿 338-r1 · AC#0 追加尺：panel 枚的 gofumpt 非幂等＋非空白变化到底是什么（⛔ 不碰工作树，只在 /d/tmp 副本上量）
set -u
GOFUMPT="/d/work/base/gopath/bin/gofumpt.exe"
W=/d/tmp/wisp/20261018/338r1/wtest2
B=/d/tmp/wisp/20261018/338r1/blobcheck
mkdir -p "$W"
cd "/d/work/workspace/projects plans/Wisp" || exit 9
P=cmd/wisp/panel_transport_35r2_test.go
T=internal/tools/task_output_pointer_notice_errtext_174r7_test.go
flatp=cmd_wisp_panel_transport_35r2_test.go
flatt=internal_tools_task_output_pointer_notice_errtext_174r7_test.go

echo '===== [A] HEAD blob（仓外 LF 形）gofmt -l ⇄ gofumpt -l 对照 ====='
for f in "$T" "$P"; do
  flat=$(printf '%s' "$f" | tr '/' '_')
  git show "HEAD:$f" > "$B/$flat"
  printf '%s | gofmt_l_rows=%s | gofumpt_l_rows=%s\n' "$f" \
    "$(gofmt -l "$B/$flat" | grep -c . || true)" \
    "$("$GOFUMPT" -l "$B/$flat" | grep -c . || true)"
done

echo '===== [B] panel 枚 gofumpt -w 幂等性：连跑三发逐发哈希（起点＝工作树字节，工作树是 LF） ====='
cp "$B/$flatp" "$W/p"
for i in 1 2 3; do
  "$GOFUMPT" -w "$W/p"
  echo "w$i sha=$(sha256sum < "$W/p" | cut -c1-16) gofumpt_l_rows=$("$GOFUMPT" -l "$W/p" | grep -c . || true) gofmt_l_rows=$(gofmt -l "$W/p" | grep -c . || true) bytes=$(wc -c < "$W/p")"
done
echo '--- 第三发之后 gofumpt -d 全文（＝仍不收敛的那部分） ---'
"$GOFUMPT" -d "$W/p"
echo "gofumpt_d_rc=$?"

echo '===== [C] panel 枚：非空白尺逐字输出（blob 起点 ⇄ 一发 -w） ====='
cp "$B/$flatp" "$W/b2"; cp "$W/b2" "$W/a2"; "$GOFUMPT" -w "$W/a2"
echo '--- diff <(tr -d " \t" < before) <(tr -d " \t" < after) ---'
diff <(tr -d ' \t' < "$W/b2") <(tr -d ' \t' < "$W/a2")
echo "RULER_RC=$?"
echo "nonblank_bytes before=$(tr -d ' \t' < "$W/b2" | wc -c) after=$(tr -d ' \t' < "$W/a2" | wc -c)"

echo '===== [D] panel 枚：手工版最小改动（只把注释列对齐）是否足够 ====='
cp "$B/$flatp" "$W/m"
gofmt -l "$W/m" | grep -c . || true
echo "gofmt_l_on_blob=$? rows=$(gofmt -l "$W/m" | grep -c . || true)"

echo '===== [E] tools 枚：同一把尺（对照，必须空输出＋收敛） ====='
cp "$B/$flatt" "$W/b3"; cp "$W/b3" "$W/a3"; "$GOFUMPT" -w "$W/a3"
diff <(tr -d ' \t' < "$W/b3") <(tr -d ' \t' < "$W/a3")
echo "RULER_RC=$? gofumpt_l_after=$("$GOFUMPT" -l "$W/a3" | grep -c . || true)"
cp "$W/a3" "$W/a3b"; "$GOFUMPT" -w "$W/a3b"
echo "second_w_identical=$([ "$(sha256sum < "$W/a3"|cut -d' ' -f1)" = "$(sha256sum < "$W/a3b"|cut -d' ' -f1)" ] && echo YES || echo NO)"
