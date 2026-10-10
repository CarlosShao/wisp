#!/usr/bin/env bash
# 111-a6b rulers R1-R6. Read-only: everything is measured against the HEAD commit
# via `git archive HEAD` into a dir OUTSIDE the repo (object layer, no working tree,
# no `go list`/`go build`/`go test` => zero Go compile surface touched).
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 2
OUT=".scratch/wisp/probes/111/a6b"
MOD="github.com/CarlosShao/wisp"
mkdir -p "$OUT"
OUT="$(pwd)/$OUT"           # absolute: the awk stage runs with cwd = the archive dir
tmpd="$(mktemp -d)"
git archive HEAD -- cmd internal | tar -x -C "$tmpd" || { echo "R0 rc=2 archive failed"; exit 2; }
# verify the archive really is HEAD (two blobs, byte-for-byte via cksum)
{ echo "R0 archive-vs-HEAD self-check"
  echo "internal/tools/fs_write_test.go  git=$(git show HEAD:internal/tools/fs_write_test.go | cksum)  disk=$(cksum "$tmpd/internal/tools/fs_write_test.go")"
  echo "cmd/wisp/wiring_test.go          git=$(git show HEAD:cmd/wisp/wiring_test.go | cksum)  disk=$(cksum "$tmpd/cmd/wisp/wiring_test.go")"
} > "$OUT/00b-archive-check.txt"

# ---------------- R1 + R2 + R1f : one awk pass over every module _test.go -------
cd "$tmpd" || exit 2
find cmd internal -name '*_test.go' | sort > /tmp/files.txt
xargs -a /tmp/files.txt awk '
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
  printf "%s\t%s\t%s\t%d\t%s\n", file, classify(tag), tag, nf, fnameclass(file)
}
{
  if (FNR == 1 && NR > 1) emit()
  if (FNR <= 3 && /^\/\/go:build/) tag = substr($0, 4, 200)
  if ($0 ~ /^func (Test|Benchmark|Fuzz|Example)/) nf++
}
END { emit() }
' > "$OUT/01-file-class.tsv"
r1=$?
cd "$REPO" || exit 2
echo "rc=$r1  # R1/R2/R1f: 4 cols = path | content-tag-class | raw tagline (first 3 lines only) | top-level func Test/Benchmark/Fuzz/Example count | filename-class" >> "$OUT/01-file-class.tsv"

# ---------------- R5 : tier pins (from the HEAD blob) + rollup ------------------
git show HEAD:scripts/portable-tests.sh > /tmp/pt.sh || exit 2
for t in core windows cli winsec; do
  sed -n "/^${t}_pin='/,/^'$/p" /tmp/pt.sh | grep -v -e "_pin='\$" -e "^'\$" | grep . > "/tmp/pin.$t"
done

awk -F'\t' -v MOD="$MOD" '
$1 ~ /^rc=/ { next }
{
  p = $1; sub(/\/[^\/]*$/, "", p); pkg = p                 # internal/tools, internal/llm/anthropic, cmd/wisp
  if (!(pkg in seen)) { seen[pkg] = 1; order[++n] = pkg }
  cls = $2; nf = $4
  files[pkg"_"cls]++
  cases[pkg"_"cls] += nf
  tot[pkg] += nf
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
echo "rc=$r2  # R5 rollup. cols: pkg | total-top-level-cases | WIN(f/c) | WIN_LIVE(f/c) | POSIX(f/c) | POSIX_OTHER(f/c) | TAG_OTHER(f/c) | NONE-untagged(f/c). Universe = git archive HEAD -- cmd internal (module pkgs only, .scratch excluded by construction)" >> "$OUT/02-pkg-summary.tsv"

# ---------------- blind-spot join: which tiers claim each pkg, which OS --------
printf 'core\tubuntu\nwindows\twindows\ncli\twindows\nwinsec\twindows\n' > /tmp/tier_os.txt
: > "$OUT/03-blind-roster.tsv"
while IFS=$'\t' read -r pkg tot winf winc livef livec posixf posixc potherf potherc totherf totherc nonef nonec; do
  case "$pkg" in rc=*|"") continue ;; esac
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
  osset="${osset#,}"
  [ -z "$osset" ] && osset="NO-SCOPE"
  # direction 1: windows-tagged cases in a package no windows-running tier claims
  bwin=0; if [ $((winf + livef)) -gt 0 ] && ! printf '%s' "$osset" | grep -q windows; then bwin=1; fi
  # direction 2: posix-tagged cases in a package no ubuntu-running tier claims
  bpos=0; if [ $((posixf + potherf)) -gt 0 ] && ! printf '%s' "$osset" | grep -q ubuntu; then bpos=1; fi
  printf '%s\ttot=%s\twin=%s/%s\tlive=%s/%s\tposix=%s/%s\tuntagged=%s/%s\ttiers=%s\trunnerOS=%s\tBLIND_WIN=%d\tBLIND_POSIX=%d\n' \
    "$pkg" "$tot" "$winf" "$winc" "$livef" "$livec" "$posixf" "$posixc" "$nonef" "$nonec" "$where" "$osset" "$bwin" "$bpos"
done < "$OUT/02-pkg-summary.tsv" >> "$OUT/03-blind-roster.tsv"
echo "rc=0  # R3+R5 join. tier->runner-OS table above is hand-read from ci.yml runs-on (see 05-ruler-meta.txt for the call-site grep); BLIND_WIN = package carries windows-tagged cases yet NO claiming tier runs on windows; BLIND_POSIX mirror" >> "$OUT/03-blind-roster.tsv"

# ---------------- R4 : per-case precondition scan (content, not name) ----------
BLIND_PKGS=$(awk -F'\t' '$0 ~ /BLIND_WIN=1/ {print $1}' "$OUT/03-blind-roster.tsv")
: > "$OUT/04-case-preconditions.tsv"
for pkg in $BLIND_PKGS; do
  while IFS=$'\t' read -r f cls tag nf fnc; do
    case "$f" in "$pkg"/*) ;; *) continue ;; esac
    [ -f "$tmpd/$f" ] || continue
    awk -v FILE="$f" -v CLS="$cls" '
      /^\/\/go:build/ && FNR <= 3 { tagline = substr($0, 4, 200) }
      /^func (Test|Benchmark|Fuzz|Example)/ {
        if (infn) flushline()
        nm = $2; sub(/\(.*/, "", nm); infn = 1; body = ""
      }
      infn { body = body " " $0 }
      /^}/ { if (infn) flushline() }
      END { if (infn) flushline() }
      function flushline() {
        t = ""
        P[1] = "WISP_LIVE_MIC"; P[2] = "t.Skip("; P[3] = "t.Setenv"; P[4] = "net.Listen"
        P[5] = "net.Dial";      P[6] = "http.Get";      P[7] = "exec.Command"; P[8] = "os/exec"
        P[9] = "registry.";     P[10] = "syscall.";     P[11] = "HWND";        P[12] = "webview"
        P[13] = "frontend/dist"; P[14] = "testing.Short"; P[15] = "time.Sleep"; P[16] = "18080"
        P[17] = "CreateWindow"; P[18] = "VirtualAlloc";  P[19] = "audio";       P[20] = "watcher"
        np = 20
        hit = ""
        for (i = 1; i <= np; i++) if (index(body, P[i])) hit = hit ";" P[i]
        printf "%s\t%s\t%s\t%s\t%s\n", FILE, CLS, nm, (hit == "" ? "-" : substr(hit, 2)), tagline
        infn = 0; body = ""; nm = ""
      }
    ' "$tmpd/$f" >> "$OUT/04-case-preconditions.tsv"
  done < <(grep -F "$(printf '%s' "$pkg")/" "$OUT/01-file-class.tsv" | awk -F'\t' -v p="$pkg" 'index($1,p"/")==1')
done
echo "rc=0  # R4 case roster for every BLIND_WIN package. cols: file | file-tag-class | top-level case name | precondition tokens found IN THE BODY (content anchor, not the name) | raw build tagline" >> "$OUT/04-case-preconditions.tsv"

rm -rf "$tmpd"
