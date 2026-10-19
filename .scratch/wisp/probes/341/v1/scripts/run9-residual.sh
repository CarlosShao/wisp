#!/usr/bin/env bash
# 341-v1 AC#2/AC#5/AC#6 residual gauges.
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9

echo "########## positive control for the single-quote canonical shape (another resident test) ##########"
grep -n "model = 'deepseek-reasoner'" internal/config/settings_313_test.go
echo "hits of the double-quote shape in resident assertions about re-laid text:"
grep -rn 'Contains(text, "\(base_url\|model\|protocol\) = "' internal/config/*_test.go | wc -l
grep -rn "Contains(text, \"model = \|base_url = \\\\\"" internal/config/*_test.go | sed 's/^/   /' | head -5

echo
echo "########## assertion-count arithmetic: pre-fix red messages vs today's file ##########"
echo "-- pre-fix (r1 logs/06 + logs/07): distinct messages, with their file:line --"
grep -oE 'writeguard_unnailed_341_test\.go:[0-9]+: [a-z].*' .scratch/wisp/probes/341/r1/logs/07-p6-raw.txt | sort -u | sed 's/^/   /'
echo "-- today's file: same three messages with line numbers --"
grep -n 'came back as something other than the built-in preset\|did not come back as an explicit empty string\|did not survive the re-lay' internal/config/writeguard_unnailed_341_test.go
echo "-- delta in lines (want 6 = the inserted comment block) --"
echo "314 -> 320 = +6 ; 368 -> 374 = +6 ; 374 -> 380 = +6"
echo "-- the inserted block, verbatim (lines 313-318) --"
sed -n '313,318p' internal/config/writeguard_unnailed_341_test.go
echo "-- t.Errorf/t.Fatal counts per case (today) --"
awk '/^func Test/{name=$2; n=0} /t\.(Errorf|Error|Fatal|Fatalf)\(/{n++} /^}$/{if(name!=""){printf "   %-70s assertion-sites=%s\n", name, n; name=""}}' internal/config/writeguard_unnailed_341_test.go

echo
echo "########## AC#2(iii): P2 case body -- does it pin P2's OWN judgement? ##########"
echo "-- P2's unique assertion (:100-:112) verbatim --"
sed -n '96,117p' internal/config/writeguard_unnailed_341_test.go
echo "-- P1's counterpart (227 file :55-:65) verbatim --"
sed -n '53,65p' internal/config/writeguard_handedit_227_test.go

echo
echo "########## AC#5: P5 case = 观测今天这形 or 夹带行为改动（乙形）? ##########"
echo "-- P5 body lines 209-270 --"
sed -n '209,232p' internal/config/writeguard_unnailed_341_test.go
echo "-- any write-behaviour call introduced? (SaveFile/applyMigrations callers in the test file) --"
grep -n "SaveFile\|applyMigrations\|readConfigFile\|Migrate" internal/config/writeguard_unnailed_341_test.go | sed 's/^/   /'
echo "-- production diff for migrate.go / loader.go since the leg's anchor: (want zero) --"
git diff --name-only 98dbf20a..HEAD -- internal/config | sed 's/^/   /'

echo
echo "########## AC#6: discipline gauges ##########"
echo "-- index clean now? (git diff --cached --name-only) --"
git diff --cached --name-only | sed 's/^/   /'; echo "   count=$(git diff --cached --name-only | wc -l)"
echo "-- the 4 junk files at repo root (the leg's own incident) --"
for n in '⛔' '主料＝（26' '前一腿' '本件＝票'; do
  [ -e "$n" ] && printf "   EXISTS %s bytes=%s mtime=%s\n" "$n" "$(wc -c < "$n")" "$(stat -c %y "$n" | cut -d. -f1)" || printf "   MISSING %s\n" "$n"
done
echo "-- are they tracked? (want no) --"
git ls-files -- '⛔' '主料＝（26' '前一腿' '本件＝票' | wc -l
echo "-- do they fall in any gate's scope? d22scan scope dirs --"
grep -nE 'addPath|dirs|"internal"|"cmd"|"design"|"frontend"' tools/d22scan/main.go | head -12
echo "-- length gate denominator = tracked paths (junk is untracked) --"
echo "   tracked=$(git ls-files | wc -l)"
echo "-- leg's 6 commits: any file outside its allowed面? --"
git show --name-only --format='' 6948eb08 fd82f786 f710d104 920537ee 46c920c8 cef2a6c9 | sort -u | sed 's/^/   /' | grep -vE '^ *$' | grep -vE '\.scratch/wisp/probes/341/r1/|^   internal/config/writeguard_unnailed_341_test\.go$|^   \.scratch/wisp/issues/341-' | sed 's/^/   OUTSIDE: /'
echo "-- count of outside files --"
git show --name-only --format='' 6948eb08 fd82f786 f710d104 920537ee 46c920c8 cef2a6c9 | sort -u | grep -vE '^$' | grep -vE '\.scratch/wisp/probes/341/r1/|internal/config/writeguard_unnailed_341_test\.go$|^\.scratch/wisp/issues/341-' | wc -l
echo "-- ticket tail bullet: numstat and append-only check --"
git show --numstat --format='' 920537ee -- .scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
git show 920537ee -- .scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md | grep -c '^-[^-]' || true
