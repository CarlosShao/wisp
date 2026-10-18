#!/bin/sh
# 337-a1：交件自证三把尺（逐枚提交号＋name-only／件绝对路径＋wc -c＋行数＋占位符尺／决定性尺逐字）
cd "/d/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/337/a1/logs
D=.scratch/wisp/probes/337/a1
{
  echo "### 337-a1 交件自证 rev 起手锚＝c11163c896b09e7b6863d347dadbcb65702585cc；落笔侧 HEAD＝$(git log -1 --format='%H %ci')"
  echo "date=$(date '+%a %b %d %H:%M:%S %z')"
  echo ""
  echo "== 尺一（本腿全部提交＋逐枚 name-only＋.go 命中枚数）=="
  for c in $(git log --format='%H' --author="$(git config user.name)" --grep='^票 337 普查腿 337-a1' | head -20); do
    echo "--- $c  $(git log -1 --format='%s' $c)"
    git show --name-only --format='' $c
    echo "    go_hits=$(git show --name-only --format='' $c | grep -c '\.go$' || true)  files=$(git show --name-only --format='' $c | grep -c . )"
  done
  echo "★提交主题前缀尺＝git log --grep '^票 337 普查腿 337-a1'（采样真实形状后用的严格前缀）命中枚数＝$(git log --oneline --grep='^票 337 普查腿 337-a1' | grep -c . )"
  echo "★越面尺＝本腿全部 name-only 去重后，不在两处允许写面（probes/337/a1/** ＋ 票 337 那枚件）的枚数："
  git log --grep='^票 337 普查腿 337-a1' --name-only --format='' | grep -v '^$' | sort -u | grep -vcE '^\.scratch/wisp/(probes/337/a1/|issues/337-)' || true
  echo ""
  echo "== 尺二（证据件绝对路径＋wc -c＋行数＋占位符尺 0 命中）=="
  echo "根＝$(pwd -W 2>/dev/null || pwd)"
  for f in "$D/census.md" $(ls "$D"/logs/*.log "$D"/scripts/*.sh "$D"/msg/*.txt "$D"/msg/*.md 2>/dev/null); do
    printf '%-78s bytes=%-7s lines=%-5s placeholder=%s\n' "$f" "$(wc -c < "$f")" "$(wc -l < "$f")" "$(grep -cE 'TODO|FIXME|placeholder|占位' "$f" || true)"
  done
  echo "★零字节件枚数（应 0）＝$(find "$D" -type f -size 0 | wc -l)"
  echo "正控（同一把占位符尺在别处应≥1）＝docs/evidence/s1/334-v1-verdict.md 命中 $(grep -cE 'TODO|FIXME|placeholder|占位' docs/evidence/s1/334-v1-verdict.md || true) 枚"
  echo ""
  echo "== 尺三（决定性尺逐字：票面七格＋写面＋票尾一条）=="
  T=.scratch/wisp/issues/337-fail-closed-reason-has-no-reader-wiring-census.md
  echo "票面框数现跑＝未勾 $(grep -c -- '^- \[ \]' "$T")／已勾 $(grep -c -- '^- \[x\]' "$T")（尺＝grep -c -- '^- [ ]' 于票面）"
  echo "票尾本腿那条＝$(grep -c '337-a1' "$T") 行；票面相对 c11163c8 的差＝$(git diff --numstat c11163c8 HEAD -- "$T" | cat)"
  echo "★本腿交件时对 internal/**／cmd/**／frontend/** 的字节改动＝"
  git diff --name-only c11163c8 HEAD | grep -E '^(internal|cmd|frontend|design|tools|scripts)/' | wc -l
  echo ""
  echo "== 必答逐字读数汇总（尺见各 logs）=="
  echo "(i)  c 跳正控 Decision.Reason 非测试枚数＝$(git grep -nE '(dec|it\.Decision|decision)\.Reason[^A-Za-z]' c11163c8 -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:' | wc -l)"
  echo "(i)  c 跳被测尺（BlacklistNote.Reason 读形）命中枚数＝$(git grep -nE '(note|n|Blacklist)\.Reason[^A-Za-z]' c11163c8 -- '*.go' | grep -v '^[^:]*:\.scratch/' | wc -l)（唯一那枚＝mode.go:111 写侧）"
  echo "(iv) 最远字节所在行逐字＝$(git show c11163c8:internal/agent/approval/pending_read.go | sed -n '67p')"
  echo "(iv) 写侧行逐字＝$(git show c11163c8:internal/tools/mode.go | sed -n '111p')"
  echo "gate 调用点逐字＝$(git show c11163c8:internal/tools/mode.go | sed -n '98p')"
  echo "risk.Logf 生产注入枚数＝$(git grep -nE 'risk\.Logf[[:space:]]*=[^=]' c11163c8 -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:' | wc -l)"
  echo "ToolOutcomeText 全仓命中枚数（票面引用名）＝$(git grep -n 'ToolOutcomeText' c11163c8 | grep -v '^[^:]*:\.scratch/' | wc -l)"
} > "$L/95-final-selfcheck.log"
wc -c "$L/95-final-selfcheck.log"
echo "cr=$(grep -c "$(printf '\r')" "$L/95-final-selfcheck.log" || true)"
