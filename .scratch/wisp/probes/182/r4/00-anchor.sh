# 182 r4 leg — anchor + tree-movement verification (script 00)
# Run via Git Bash: "<gitbash>" .scratch/wisp/probes/182/r4/00-anchor.sh
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
unset GOTMPDIR TMPDIR TEMP TMP
echo "== bash =="; bash --version | head -1
echo "== HEAD =="; git rev-parse --short=9 HEAD; git log --oneline -5
echo "== r3 anchor ac2c1d903 exists? =="; git cat-file -t ac2c1d903 || echo MISSING
echo "== dispatch anchor 446831ea9 exists? =="; git cat-file -t 446831ea9 || echo MISSING
echo "== commits ac2c1d903..HEAD =="; git log --oneline ac2c1d903..HEAD
echo "== scope movement ac2c1d903..HEAD (internal/ cmd/ design/ PLAN specs go.mod go.sum) =="
git diff --name-only ac2c1d903..HEAD -- internal/ cmd/ design/ docs/PLAN.md docs/specs/ go.mod go.sum
echo "-- end movement (empty above = zero) --"
echo "== dffa9fd64 file list =="
git show --name-only --oneline dffa9fd64 | tail -n +2
echo "== untracked in-flight files touching compile scope =="
git status --porcelain -- internal/ cmd/ | head -20
echo "== ticket files (181/189/179/176/177/175/167/163/191/92/188/197/192/168/190/407) =="
ls .scratch/wisp/issues/ | grep -E "^181-|^189-|^179-|^176-|^177-|^175-|^167-|^163-|^191-|^92-|^188-|^197-|^192-|^168-|^190-|^407-"
exit 0
