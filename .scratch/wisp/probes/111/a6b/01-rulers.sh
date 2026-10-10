#!/usr/bin/env bash
# 111-a6b rulers R1-R6. Read-only: everything measured against HEAD via
# `git archive HEAD` into a dir OUTSIDE the repo (object layer; no working tree used
# as HEAD; ZERO Go compile surface — no go build/test/vet/run/list).
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 2
REL=".scratch/wisp/probes/111/a6b"
MOD="github.com/CarlosShao/wisp"
mkdir -p "$REL"
OUT="$(pwd)/$REL"
tmpd="$(mktemp -d)"
awkf="$(mktemp)"
git archive HEAD -- cmd internal | tar -x -C "$tmpd" || { echo "R0 rc=2 archive failed"; exit 2; }
wtest=$(git ls-files -- 'cmd/wisp/*_test.go' | head -1)
{
  echo "R0 archive-vs-HEAD self-check (cksum must match per row)"
  echo "internal/tools/fs_write_test.go  git=$(git show HEAD:internal/tools/fs_write_test.go | cksum)  disk=$(cksum "$tmpd/internal/tools/fs_write_test.go")"
  echo "$wtest  git=$(git show "HEAD:$wtest" | cksum)  disk=$(cksum "$tmpd/$wtest")"
} > "$OUT/00b-archive-check.txt"

# ---------------- R1 + R2 + R1f : one awk pass over every module _test.go -------
cat > "$awkf" <<'AWK'
function classify(t) {
  if (t == "") return "NONE"
  if (index(t, "!windows")) return "POSIX"
  if (index(t, "windows") && index(t, "winlive")) return "WIN_LIVE"
  if (index(t, "windows")) return "WIN"
  if (index(t, "linux") || index(t, "darwin")) return "POSIX_OTHER"
  return "TAG_OTHER"
}
function fnameclass(f) {
  if (index(f, "_windows_test.go")) return "FN_WIN"
  if (index(f, "_posix_test.go") || index(f, "_other_test.go")) return "FN_POSIX"
  return "FN_PLAIN"
}
function emit() {
  printf "%s\t%s\t%s\t%d\t%s\n", file, classify(tag), (tag == "" ? "-" : tag), nf, fnameclass(file)
  tag = ""; nf = 0; file = ""
}
{
  if (FNR == 1) { if (NR > 1) emit(); file = FILENAME }
  if (FNR <= 3 && /^\/\/go:build/) tag = $0
  if ($0 ~ /^func (Test|Benchmark|Fuzz|Example)/) nf++
}
END { if (file != "") emit() }
AWK
find "$tmpd" -name '*_test.go' | sed "s#^$tmpd/##" | sort > /tmp/files.txt
xargs -a /tmp/files.txt awk -f "$awkf" > "$OUT/01-file-class.tsv"
r1=$?
echo "rc=$r1  # R1/R2/R1f. cols: path | content-tag-class | raw //go:build line (searched in FIRST 3 LINES ONLY) | top-level Test/Benchmark/Fuzz/Example count | filename-class. Universe=$(grep -c . /tmp/files.txt) module _test.go under cmd internal from git archive HEAD (no .scratch, no sub-module dirs)" >> "$OUT/01-file-class.tsv"

# ---------------- R5 : tier pins (HEAD blob) + per-package rollup --------------
git show HEAD:scripts/portable-tests.sh > /tmp/pt.sh || exit 2
declare -A PINVAR=([core]=core_pin [windows]=win_pin [cli]=cli_pin [winsec]=winsec_pin)
for t in core windows cli winsec; do
  v=${PINVAR[$t]}
  sed -n "/^${v}='/,/^'$/p" /tmp/pt.sh | grep -v -e "${v}='\$" -e "^'\$" | grep . > "/tmp/pin.$t"
  echo "tier=$t var=${v} rows=$(grep -c . /tmp/pin.$t || true)"
done > /tmp/pincounts.txt

awk -F'\t' '
$1 ~ /^rc=/ { next }
{
  p = $1; sub(/\/[^\/]*$/, "", p); pkg = p
  if (!(pkg in seen)) { seen[pkg] = 1; order[++n] = pkg }
  cls = $2; nf = $4 + 0
  files[pkg"_"cls]++; cases[pkg"_"cls] += nf; tot[pkg] += nf
  pkgs[pkg] = 1
}
END {
  for (i = 1; i <= n; i++) {
    p = order[i]
    printf "%s\t%d\t%d/%d\t%d/%d\t%d/%d\t%d/%d\t%d/%d\t%d/%d\n", p, tot[p],
      files[p"_WIN"], cases[p"_WIN"],
      files[p"_WIN_LIVE"], cases[p"_WIN_LIVE"],
      files[p"_POSIX"], cases[p"_POSIX"],
      files[p"_POSIX_OTHER"], cases[p"_POSIX_OTHER"],
      files[p"_TAG_OTHER"], cases[p"_TAG_OTHER"],
      files[p"_NONE"], cases[p"_NONE"]
  }
}' "$OUT/01-file-class.tsv" | sort > "$OUT/02-pkg-summary.tsv"
r2=$?
cat /tmp/pincounts.txt
echo "rc=$r2  # R5 rollup of R1. cols: pkg | total-top-level-cases | WIN(files/cases) | WIN_LIVE(f/c) | POSIX(f/c) | POSIX_OTHER(f/c) | TAG_OTHER(f/c) | NONE-untagged(f/c)" >> "$OUT/02-pkg-summary.tsv"

# ---------------- blind-spot join: which tiers claim each pkg, on which OS -----
printf 'core\tubuntu\nwindows\twindows\ncli\twindows\nwinsec\twindows\n' > /tmp/tier_os.txt
: > "$OUT/03-blind-roster.tsv"
while IFS=$'\t' read -r pkg tot wpair lpair ppair popair topair npair; do
  case "$pkg" in rc=*|"") continue ;; esac
  winf=${wpair%%/*}; wc_=${wpair##*/}
  lu=${lpair%%/*}; lc=${lpair##*/}
  pu=${ppair%%/*}; pc=${ppair##*/}
  po_=${popair%%/*}; pofc=${popair##*/}
  tu=${topair%%/*}; uc=${topair##*/}
  nu=${npair%%/*}; nc=${npair##*/}
  ip="$MOD/$pkg"
  where=""
  for t in core windows cli winsec; do
    if grep -qxF "$ip" "/tmp/pin.$t" 2>/dev/null; then where="$where,$t"; fi
  done
  where="${where#,}"
  osset=""
  for t in ${where//,/ }; do
    os="$(awk -F'\t' -v T="$t" '$1==T{print $2}' /tmp/tier_os.txt)"
    osset="$osset,$os"
  done
  osset="${osset#,}"; [ -z "$osset" ] && osset="NO-SCOPE"
  win_cases=$((wc_ + lc)); posix_cases=$((pc + pofc))
  bwin=0; if [ "$win_cases" -gt 0 ] && ! printf '%s' "$osset" | grep -q -e windows; then bwin=1; fi
  bpos=0; if [ "$posix_cases" -gt 0 ] && ! printf '%s' "$osset" | grep -q -e ubuntu; then bpos=1; fi
  printf '%s\ttot=%s\tWIN=%s/%s\tWIN_LIVE=%s/%s\tPOSIX=%s/%s\tPOSIX_OTHER=%s/%s\tTAG_OTHER=%s/%s\tUNTAGGED=%s/%s\ttiers=%s\trunnerOS=%s\tBLIND_WIN=%d\tBLIND_POSIX=%d\n' \
    "$pkg" "$tot" "$winf" "$wc_" "$lu" "$lc" "$pu" "$pc" "$po_" "$pofc" "$tu" "$uc" "$nu" "$nc" "$where" "$osset" "$bwin" "$bpos"
done < "$OUT/02-pkg-summary.tsv" >> "$OUT/03-blind-roster.tsv"
echo "rc=0  # R3+R5 join. tier->runner-OS read from ci.yml runs-on (see 05-ruler-meta.txt). BLIND_WIN=pkg carries windows-tagged cases but NO claiming tier runs on windows. BLIND_POSIX=mirror." >> "$OUT/03-blind-roster.tsv"

# ---------------- R4 : per-case precondition scan (content, not name) ----------
BLIND_PKGS=$(awk -F'\t' '$0 ~ /BLIND_WIN=1/ {print $1}' "$OUT/03-blind-roster.tsv")
echo "BLIND_PKGS: $BLIND_PKGS"
: > "$OUT/04-case-preconditions.tsv"
for pkg in $BLIND_PKGS; do
  while IFS=$'\t' read -r f cls tag nf fnc; do
    [ -f "$tmpd/$f" ] || continue
    awk -v FILE="$f" -v CLS="$cls" '
      /^func (Test|Benchmark|Fuzz|Example)/ {
        if (infn) flushline()
        nm = $2; sub(/\(.*/, "", nm); infn = 1; body = ""
      }
      infn { body = body " " $0 }
      /^}/ { if (infn) flushline() }
      END { if (infn) flushline() }
      function flushline() {
        P[1] = "WISP_LIVE_MIC"; P[2] = "t.Skip("; P[3] = "t.Setenv"; P[4] = "net.Listen"
        P[5] = "net.Dial";      P[6] = "http.Get";      P[7] = "exec.Command"; P[8] = "os/exec"
        P[9] = "registry.";     P[10] = "syscall.";     P[11] = "HWND";        P[12] = "webview"
        P[13] = "frontend/dist"; P[14] = "testing.Short"; P[15] = "time.Sleep"; P[16] = "18080"
        P[17] = "CreateWindow"; P[18] = "VirtualAlloc";  P[19] = "mic";         P[20] = "watcher"
        hit = ""
        for (i = 1; i <= 20; i++) if (index(body, P[i])) hit = hit ";" P[i]
        printf "%s\t%s\t%s\t%s\n", FILE, CLS, nm, (hit == "" ? "-" : substr(hit, 2))
        infn = 0; body = ""; nm = ""
      }
    ' "$tmpd/$f" >> "$OUT/04-case-preconditions.tsv"
  done < <(awk -F'\t' -v p="$pkg" 'index($1, p "/")==1' "$OUT/01-file-class.tsv")
done
echo "rc=0  # R4. cols: file | file-tag-class | top-level case name | precondition tokens found INSIDE THE CASE BODY (content anchors: mic/env/skip/net/exec/registry/window/frontend). '-' = none of the 20 tokens in the body." >> "$OUT/04-case-preconditions.tsv"

rm -rf "$tmpd" "$awkf"
