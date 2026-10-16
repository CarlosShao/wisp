#!/bin/sh
# 325-v1 AC#4 门禁：四枚逐枚自落 rc ＋ 两枚台面雷的处理（具名，⛔ 当红、⛔ 顶包）
# 用法：sh v1-ac4-gates.sh <repo> <tree> <out>
set -u
repo=${1:?repo}; tree=${2:?tree}; out=${3:?out}
GOFUMPT=/d/work/base/gopath/bin/gofumpt.exe

one() { # label, then command in the exported tree
    echo "@@@@ $1"
    ( cd "$tree" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT PATH="$tree/third_party/sherpa-onnx:$PATH" "${@:2}" ) > "$out.$1.txt" 2>&1
    rc=$?
    echo "  RC=$rc   bytes=$(wc -c < "$out.$1.txt")"
    tail -6 "$out.$1.txt" | sed 's/^/  | /'
    echo
}

{
echo "V1 AC#4 gates"
echo "start=$(date '+%F %T %z')   tree=$tree   repo=$repo"
echo
one bash-n-runtests bash -n tools/d22scan/runtests.sh
one d22scan-sh      sh scripts/d22scan.sh
one go-build-all    env GOFLAGS= go build ./...
echo '@@@@ length-budget gate --with-self-test  (NEEDS the git index, so this one runs read-only IN THE REPO)'
( cd "$repo" && sh scripts/check-path-length-budget.sh --with-self-test ) > "$out.length-budget.txt" 2>&1
echo "  RC=$?   bytes=$(wc -c < "$out.length-budget.txt")"
grep -E 'VERDICT|over-budget|self-test' "$out.length-budget.txt" | tail -8 | sed 's/^/  | /'
echo
echo '@@@@ gofumpt: the strict ruler - find its real body first'
echo "  which gofumpt rc=$(which gofumpt >/dev/null 2>&1; echo $?)   (1 = not on PATH, which is why 325-r1 wrote 找不到真身)"
echo "  absolute body present? $([ -x "$GOFUMPT" ] && echo YES || echo NO)  -> $GOFUMPT"
"$GOFUMPT" --version 2>&1 | sed 's/^/  | /'
echo "  .go files touched by the ticket commit 428ef21e = $(git -C "$repo" show --stat --format='' 428ef21e | grep -c '\.go |' || true)   (the diff is one .sh)"
echo "  gofumpt -l over tools/ scripts/ (pre-existing repo volume, NOT this ticket's debt):"
( cd "$tree" && "$GOFUMPT" -l tools scripts ) > "$out.gofumpt-l.txt" 2>&1
echo "    RC=$?  offenders=$(grep -c . "$out.gofumpt-l.txt" || true)"
sed 's/^/    | /' "$out.gofumpt-l.txt"
echo "  contrast ruler (the one that must NOT be passed off): gofmt -l over the same dirs gives $( cd "$tree" && gofmt -l tools scripts | wc -l ) files"
echo
echo '@@@@ ticket 311 volume guard: did any run above get the rc=2 refusal?'
grep -rn 'ticket 311 AC#1 shape B: GUARD' "$out".*.txt 2>/dev/null | sed 's/^/  /' || echo "  no GUARD line in any gate output (root stayed on the %LOCALAPPDATA% volume C:) - control: the string above is searched with grep -rn, and the files listed here exist: $(ls "$out".*.txt 2>/dev/null | wc -l)"
echo
echo '@@@@ table-fault holes scan over every gate capture'
for pat in 'walk' 'build failed' 'no tests to run' '0xc0000135'; do
    printf "  pat=%-20s hits=%s\n" "$pat" "$(grep -rlF "$pat" "$out".*.txt 2>/dev/null | wc -l)"
done
echo "  (positive control for the line above: captures that DO contain text: $(grep -rlc 'runtests.sh' "$out".*.txt 2>/dev/null | wc -l))"
echo "end=$(date '+%F %T %z')"
} >> "$out" 2>&1
echo "rc=$? out=$out"
