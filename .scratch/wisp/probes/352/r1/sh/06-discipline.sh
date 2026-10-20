#!/bin/sh
# 352-r1 AC#3 纪律自证尺（名册＋写面，两把独立锚对撞；★⛔ 拿 A..B 区间数"谁动了某文件"＝区间里有别腿的笔）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
OUT=/d/tmp/wisp/20261020/352r1/logs/discipline.txt
: > "$OUT"
{
  echo "== R1 名册尺＝锚定行首（⛔ 假设其后任何字符、⛔ 非锚定子串）=="
  git log --grep='^票 352 落地腿 352-r1' --format='%h %ad %s' | tee -a "$OUT.r1"
  echo "ANCHORED_COUNT=$(git log --grep='^票 352 落地腿 352-r1' --format='%h' | grep -c .)"
  echo "NON_ANCHORED_COUNT=$(git log --grep='票 352 落地腿 352-r1' --format='%h' | grep -c .)  （非锚定＝会吃进正文引用了本前缀的别腿笔＝A1045 那味）"
  echo
  echo "== R2 第二把独立尺＝单次 %H%x09%s 导出＋内存内 index(\$2,P)==1 =="
  git log --format='%H%x09%s' > /d/tmp/wisp/20261020/352r1/logs/subjects.tsv
  perl -e '
    my $P = "票 352 落地腿 352-r1";
    my @m;
    open my $fh, "<", "/d/tmp/wisp/20261020/352r1/logs/subjects.tsv" or die;
    while (<$fh>) { chomp; my ($h,$s) = split /\t/, $_, 2; next unless defined $s; push @m, "$h\t$s" if index($s, $P) == 0; }
    print "INMEMORY_ANCHOR0_COUNT=", scalar(@m), "\n";
    print "  MY: $_\n" for @m;
  ' >> "$OUT" 2>&1
  echo
  echo "== R3 写面自证＝逐枚 name-only（只数本腿那几笔，⛔ 区间）=="
  for c in $(git log --grep='^票 352 落地腿 352-r1' --format='%H'); do
    echo "--- commit $c ---"
    git show --name-only --format='%s' "$c"
  done | grep -v '^$' > "$OUT.r3"
  { for c in $(git log --grep='^票 352 落地腿 352-r1' --format='%H'); do git show --name-only --format='' "$c"; done; } | sort -u | grep -v '^$' > /d/tmp/wisp/20261020/352r1/logs/mypaths.txt
  echo "PATHS_TOTAL=$(grep -c . /d/tmp/wisp/20261020/352r1/logs/mypaths.txt)"
  echo "OUTSIDE_ALLOWED_WRITEFACE=$(grep -v -E '^(internal/winsec/winsec\.go|\.scratch/wisp/probes/352/r1/|\.scratch/wisp/issues/352-)' /d/tmp/wisp/20261020/352r1/logs/mypaths.txt | grep -c . || true)"
  echo "GO_FILES_TOUCHED=$(grep -c '\.go$' /d/tmp/wisp/20261020/352r1/logs/mypaths.txt || true)  （★产码只 1 枚＝internal/winsec/winsec.go；⛔ 一枚 .go 是别腿的）"
  echo "WINSEC_GO_LINES=$(grep -c 'internal/winsec/winsec.go' /d/tmp/wisp/20261020/352r1/logs/mypaths.txt || true)"
  echo
  echo "== R4 那一行那笔的 -U0 全貌＝只动那一枚格式串（⛔ 任何别枚字节）=="
  git show -U0 --format='COMMIT=%s' ecffb93be03bd86f647c1de468783bcef35a9fc9
  echo "NET_DIFFSTAT=$(git show --numstat --format='' ecffb93be03bd86f647c1de468783bcef35a9fc9 | awk '{print $1" insertions "$2" deletions"}')"
  echo
  echo "== R5 禁区逐枚现验＝⛔ 动过=="
  for P in internal/winsec/pinnedleaf_342_windows.go internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go internal/winsec/pinnedleaf_342_unix.go internal/winsec/pinnedleaf_342_fallback.go internal/risk internal/config frontend design cmd/wisp internal/agent internal/panel; do
    n=$(for c in $(git log --grep='^票 352 落地腿 352-r1' --format='%H'); do git show --name-only --format='' "$c"; done | grep -c -F "$P" || true)
    printf '  %-45s 命中本腿笔数=%s\n' "$P" "$n"
  done
  echo
  echo "== R6 票 352 框尺（现量，⛔ 采信任何旧叙述）=="
  # ★本腿自抓第一版这把尺＝`^  - \[x\]` 带两枚前导空格⇒恒 0 命中＝钝尺（票面那四枚框是行首 `- [ ]`，
  #   零命中的"0 勾"与真值 0 恰好撞脸＝最危险那味）；现改成行首锚＋另打一把"总框数"独立尺对撞。
  T=.scratch/wisp/issues/352-wrappath-errno-half-stringified-so-fs-ernotexist-cannot-punch-through-private-route.md
  echo "BOX_CHECKED=$(grep -cE '^- \[x\] ' "$T" || true)"
  echo "BOX_UNCHECKED=$(grep -cE '^- \[ \] ' "$T" || true)"
  echo "BOX_ANY=$(grep -cE '^- \[[ x]\] ' "$T" || true)   （★CHECKED+UNCHECKED 必＝ANY，否则那把尺钝）"
  echo "AC_NAMES=$(grep -oE 'AC#[0-9]' "$T" | sort -u | tr '\n' ' ')  （★凡报几格先跑框尺＝本票分母＝4 格 AC#0..AC#3）"
  echo "  ★正控（报 0 命中前先跑一枚已知存在的正控）＝同两把尺数票 348（应有勾）："
  T2=.scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md
  echo "    CONTROL_348_CHECKED=$(grep -cE '^- \[x\] ' "$T2" || true) CONTROL_348_UNCHECKED=$(grep -cE '^- \[ \] ' "$T2" || true) CONTROL_348_ANY=$(grep -cE '^- \[[ x]\] ' "$T2" || true)"
  echo
  echo "== R7 票尾 numstat＝须纯追加（0 deletions）=="
  git log --format='%h %s' --grep='^票 352 落地腿 352-r1' | while IFS= read -r l; do
    h=${l%% *}; echo "  $l"
    git show --numstat --format='' "$h" | grep -F 'issues/352-' || echo "     （本笔⛔ 碰票面件）"
  done
  echo
  echo "== R8 被禁字面字符（码位 U+00A7）尺＝只数我写面的字节；正控＝同尺数 HANDOVER.md =="
  perl -e '
    my $band = "\xc2\xa7";
    my $tot = 0;
    open(my $lh, "<", "/d/tmp/wisp/20261020/352r1/logs/mypaths.txt") or die;
    my @paths; while (<$lh>) { chomp; push @paths, $_ if length; } close $lh;
    for my $p (@paths) { next unless -f $p; open(my $f, "<", $p) or next; local $/; my $s = <$f>; close $f;
      my $n = () = $s =~ /$band/g; $tot += $n; print "  HIT $p = $n\n" if $n; }
    print "MY_WRITEFACE_U00A7_TOTAL=$tot\n";
    open(my $h, "<", "docs/reports/HANDOVER.md") or die; local $/; my $hs = <$h>; close $h;
    my $c = () = $hs =~ /$band/g; print "POSITIVE_CONTROL_HANDOVER_U00A7=$c  （正控在场＝尺灵）\n";
  ' >> "$OUT" 2>&1
  echo
  echo "== R9 零字节尺＋.out 尺（我写面逐枚）=="
  Z=0; O=0
  while IFS= read -r p; do
    [ -f "$p" ] || continue
    s=$(wc -c < "$p"); [ "$s" = 0 ] && { echo "  ZERO_BYTE $p"; Z=$((Z+1)); }
    case "$p" in *.out) echo "  OUT_EXT $p"; O=$((O+1));; esac
  done < /d/tmp/wisp/20261020/352r1/logs/mypaths.txt
  echo "ZERO_BYTE_FILES=$Z  OUT_EXT_FILES=$O"
  echo
  echo "== R10 起手锚⇄现 HEAD＋winsec.go blob 各笔串（供编排者判 351-r2 读数可比性）=="
  git log --grep='^票 352 落地腿 352-r1' --format='%h' | tac | while IFS= read -r c; do
    printf '  %s winsec_blob=%s\n' "$c" "$(git rev-parse "$c:internal/winsec/winsec.go")"
  done
  printf '  START_ANCHOR=bf509bb837aa1328617320608be38e4bca5a5d35 winsec_blob=%s\n' "$(git rev-parse bf509bb8:internal/winsec/winsec.go)"
  printf '  NOW_HEAD=%s winsec_blob=%s worktree_hash=%s\n' "$(git rev-parse HEAD)" "$(git rev-parse HEAD:internal/winsec/winsec.go)" "$(git hash-object internal/winsec/winsec.go)"
  echo
  echo "== R11 ⛔ push 自证＝本地 dev ⇄ origin/dev =="
  git rev-list --count origin/dev..HEAD | sed 's/^/UNPUSHED_COMMITS_TOTAL=/'
  git rev-list --count origin/dev..HEAD --grep='^票 352 落地腿 352-r1' | sed 's/^/UNPUSHED_COMMITS_MINE=/'
  echo
  echo "== R12 ⛔ 越界碰别腿写面（本腿笔里应⛔ 含 cmd/wisp／agent/approval 任何字节）=="
  for c in $(git log --grep='^票 352 落地腿 352-r1' --format='%H'); do
    git show --name-only --format='' "$c" | grep -E '^(cmd/wisp|internal/agent/approval)/' && echo "  VIOLATION in $c"
  done; echo "  CROSSLEG_HITS=0（上面无输出＝0）"
} >> "$OUT" 2>&1
echo "WROTE $OUT"
wc -l "$OUT"
