#!/usr/bin/env bash
# 341-v1 final state + self-discipline record.
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
W=/d/tmp/wisp/20261018/341v1

echo "########## FINAL STATE (341-v1) ##########"
echo "-- HEAD now (drift is expected: other legs commit on dev) --"; git log -1 --format='%H %s' | cut -c1-72
echo "-- my own roster, strict per-hash prefix 「票 341 终裁腿 341-v1·」 --"
git log --all --grep='票 341 终裁腿 341-v1·' --format='%H %s' | cut -c1-90
echo "-- my commit 95b1015d is an ancestor of HEAD? --"
git merge-base --is-ancestor 95b1015dcccf04ac0b8c1b1ac161ae09ae24d9e6 HEAD && echo "   YES" || echo "   NO"
echo "-- leg's test-file blob: HEAD vs the leg's self-reported c577f1cc --"
echo "   HEAD blob   = $(git rev-parse HEAD:internal/config/writeguard_unnailed_341_test.go)"
echo "   leg reported= c577f1ccd4587cd48c2168f7eb1eb767b720450f"
echo "-- the two production files I mutated: tree bytes vs pristine blob (want match) --"
for x in writeguard.go loader.go migrate.go; do
  printf "   %-14s treeD_sha256=%s blob_sha256=%s\n" "$x" \
    "$(sha256sum $W/treeD/internal/config/$x | cut -d' ' -f1 | cut -c1-12)" \
    "$(git show HEAD:internal/config/$x | sha256sum | cut -d' ' -f1 | cut -c1-12)"
done
echo "-- the leg's start anchor .. HEAD: only non-test .go in internal/config (want 1 = new test file) --"
git log --name-only --format='' 98dbf20a..HEAD -- 'internal/config/*.go' | grep -v '^$' | sort -u | sed 's/^/   /'
echo "-- worktree clean in my lanes? (want empty) --"; git status --porcelain -- internal/config cmd/wisp internal/risk
echo "-- shared index now (want 0) --"; git diff --cached --name-only | wc -l

echo
echo "########## 票框终量尺 (want 7 unchecked / 0 checked, title+Status untouched) ##########"
T=".scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md"
echo "unchecked=$(grep -c '^- \[ \] \*\*AC#' "$T" || true) checked=$(grep -c '^- \[x\]' "$T" || true) lines=$(wc -l < "$T")"
echo "title    = $(head -1 "$T" | cut -c1-40)..."
echo "Status   = $(sed -n '3p' "$T" | cut -c1-40)..."

echo
echo "########## 341-r1 两处纪律事故：我这一趟的独立复跑 ##########"
echo "-- (1) repo-root junk files: tracked? size? in any gate's walk scope? --"
for n in '⛔' '主料＝（26' '前一腿' '本件＝票'; do
  if [ -e "$n" ]; then
    printf "   name=%s bytes=%s tracked=%s\n" "$n" "$(wc -c < "$n")" "$(git ls-files -- "$n" | wc -l)"
  fi
done
echo "   d22scan walk scopes (from tools/d22scan/main.go) = internal/ cmd/ frontend/ internal/tools/ design/ ; 仓库根散件⛔ 在射程内"
grep -n 'filepath.Join(root,' tools/d22scan/main.go | sed 's/^/   /' | head -8
echo "-- (2) the hang episode: those 7 files are all committed now (index empty above); --"
echo "   决定性读数全部在事故之后由我现跑（blob 同一性 ⇄ 腿自报锚逐枚对撞，见上）。"

echo
echo "########## 341-v1 self-discipline: my own U+00A7 hit and fix (named, ⛔ 掩盖) ##########"
echo "logs/06 抓 226 判语件原句时把两处 U+00A7 抄进了我的证据面＝派单明令⛔ 出现在我写面。"
echo "处置＝perl -CSD 逐枚替换为字面 [U+00A7]（原文⛔ 改、只在我这枚镜像件里遮字符），改后尺＝"
printf "   U+00A7 hits in .scratch/wisp/probes/341/v1 = %s\n" "$(grep -rlP '\x{00A7}' .scratch/wisp/probes/341/v1 | wc -l)"
echo "   positive control, same ruler on docs/reports/HANDOVER.md = $(grep -rlP '\x{00A7}' docs/reports/HANDOVER.md | wc -l) file(s)"
echo "   mask markers now in logs/06 = $(grep -c '\[U+00A7\]' .scratch/wisp/probes/341/v1/logs/06-ac2-quote-shape-and-ac3.txt || true)"
echo "   被遮原句仍在盘上原件＝docs/evidence/s1/226-config-write-no-clobber-r1.md:111 与 v2.md:162（我一字未动）"

echo
echo "########## 台面与突变卫生 ##########"
echo "-- export trees (仓外, D: 卷, ⛔ 盘符根直下, 只建⛔ 删) --"
ls -d "$W/treeD" "$W/treeB" "$W/pristine" | sed 's/^/   /'
echo "-- GOTMPDIR / GOCACHE used by runs (all under D:\tmp\wisp\20261018\341v1) --"
ls -d "$W"/run-*  "$W"/log-p* 2>/dev/null | sed 's/^/   /'
echo "-- PATH passed to go test in bash form was cygpath -w'ed (公共块 A1042 定式)：见 scripts/run2-color.sh / run3-mu.sh --"
