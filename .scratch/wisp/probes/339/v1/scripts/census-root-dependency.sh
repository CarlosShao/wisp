#!/bin/bash
# 339-v1 台件⑤c：★同族普查的"对形的尺"（更正 ⑤b 那把尺的判据形状，⛔ 删旧件）。
# 教训＝票 339 的缺陷从来不是环境变量问题（e6fd0f50 上那件 :86 早就在设 HOME，已现量），
#       所以"按 env 拼写筛无档件"这把尺结构上抓不到它＝0 命中是必然、不是好消息。
# 真同族形＝**仪器断言依赖 C26 handle-resolved root（formsOf/anchorForms 的 .root、
#       uncertainAnchorMiss 的可燃性）而该件没有任何平台档** ⇒ POSIX 上 root 恒空那一维会咬人。
# 逐枚给 file:line；正控＝同一把尺在 known 在场件上必须开火。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
LOG="$REPO/.scratch/wisp/probes/339/v1/logs/v1-census-root-dependency.txt"
{
  echo "== 尺甲＝internal/risk/** 里引用 formsOf|anchorForms|uncertainAnchorMiss 的 _test 件，逐枚验档 =="
  for f in $(git grep -lE 'formsOf|anchorForms|uncertainAnchorMiss' HEAD -- 'internal/risk/*_test.go' | sed 's/^HEAD://'); do
    T=$(grep -c '^//go:build' "$f" || true)
    TAGV=$(grep -m1 '^//go:build' "$f" || true)
    HITS=$(grep -nE 'formsOf|anchorForms|uncertainAnchorMiss' "$f" | cut -d: -f1 | paste -sd, -)
    case "$f" in
      *_windows_test.go|*_other_test.go) SFX=suffixed ;;
      *) SFX=no-suffix ;;
    esac
    echo "FILE=$f tag_lines=$T tag=\"$TAGV\" name=$SFX ref_lines=$HITS"
  done
  echo
  echo "== 尺乙＝上述件里既无 //go:build 行、文件名也无后缀的（＝真同族未档件）=="
  for f in $(git grep -lE 'formsOf|anchorForms|uncertainAnchorMiss' HEAD -- 'internal/risk/*_test.go' | sed 's/^HEAD://'); do
    T=$(grep -c '^//go:build' "$f" || true)
    case "$f" in *_windows_test.go|*_other_test.go) T=9;; esac
    [ "$T" -eq 0 ] && echo "UNTAGGED-DEPENDENT file=$f lines=$(grep -nE 'formsOf|anchorForms|uncertainAnchorMiss' "$f" | cut -d: -f1 | paste -sd, -)"
  done
  echo
  echo "== 尺丙＝正控：同一把尺打在票 339 落码前的旧锚 e6fd0f50（那时 334 件无档＝必须开火）=="
  for f in $(git grep -lE 'formsOf|anchorForms|uncertainAnchorMiss' e6fd0f50 -- 'internal/risk/*_test.go' | sed 's/^e6fd0f50://'); do
    T=$(git show "e6fd0f50:$f" | grep -c '^//go:build' || true)
    [ "$T" -eq 0 ] && echo "  OLD-UNTAGGED file=$f (tag_lines=$T)"
  done
  echo
  echo "== 尺丁＝internal/risk 之外有没有同族依赖（必答 v）=="
  for f in $(git grep -lE 'formsOf|anchorForms|uncertainAnchorMiss' HEAD -- '*.go' ':!internal/risk' | sed 's/^HEAD://'); do
    echo "  OUTSIDE file=$f"
  done
  echo "outside_count=$(git grep -lE 'formsOf|anchorForms|uncertainAnchorMiss' HEAD -- '*.go' ':!internal/risk' | wc -l)"
  echo
  echo "== 尺戊＝跨包消费者尺：哪些包直接读 pathForms/.root 形状（只读数）=="
  echo "  (pathForms 是包内类型，导出面只有 Classify/Gate/Resolve ⇒ 外包零依赖是结构必然)"
  git grep -lE 'risk\.(Classify|Gate)\(' HEAD -- 'internal/*/*.go' ':!internal/risk' | wc -l
} > "$LOG" 2>&1
cat "$LOG"
