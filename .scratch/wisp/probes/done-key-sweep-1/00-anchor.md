# done-key-sweep-1 · 00 起手锚

腿名：`done-key-sweep-1`（只读账目腿）
仓库：`D:\work\workspace\projects plans\Wisp`，分支 `dev`
射程：`.scratch/wisp/issues/*.md` 全池（每一枚 `NN-*.md`，含 `-done` 的）

## 起手三读数（现取，写前一刻）

- `date`：`Thu Oct  8 19:10:08 CST 2026`
- `git log -1 --format='%H %ad %s'`：
  `b8b3937a2e601abd4602e4221e5cae5f06b17937 Thu Oct 8 19:09:54 2026 +0800 A738 落账＋票 186 面钉入（原句未改、不翻框）：收 next-instruments-brief-1（c219693c/a6ec48d3）两案简案，落点我逐处核对（approval.go:564-567 成员扫描 / queue.go:390 state 卫 / 预判红处 ticket259_denial_rulers_test.go:139:208:259 / workspace.go:120 Actable 闸+:144-146 死枝 全对上）。案一（181 :141 盲区）判语＝成功路径可达但值域单点{false}；A 现在装=给死枝打钉只锁字面、B 等 186 接线后装⇒★我裁走 B 并把三读数（调用者 0→?/Actable 闸/可达性）钉进票 186 新节点（不改任何既有框）。案二（259 AC#2 两发残余突变）＝发1 approval.go:564-567 会红 :139/:259、发2 queue.go:390 会红 :208；种前 hash 67fb1468/66fec7ae＋四件套尺在件里⇒留给下一枚 approval 面裁决腿（时机＝comment-fix-land-1 退出后同图互斥）；它没种＝照收未种即未证。它具名报回我两处转述差异我采（state 检查不在 spend 内／原话在 A734 非 v1 evidence）。⛔ 零翻框零 push`
- `git status --porcelain | head -20`（只登记，不动任何一件）：

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
```

（以上为 head -20 截断；全量未取。本腿不改其中任何一件。）

## 硬性禁区（本腿自缚，逐条抄自派单）

- 零 Go 命令（有腿 `comment-fix-land-1` 在 `cmd/wisp`／`internal/panel`／`internal/risk` 作业）。
- 只改只新建 `.scratch/wisp/probes/done-key-sweep-1/*.md`；不新建票、不改台账、不帮改名、零翻框、不 push。
- commit 必带显式 pathspec；禁 `add -A`/`.`/`-a`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`。
- 不读不引 `frontend/**`／`design/**`；凭据不进对话。
- 范围尺不用 `2>/dev/null`；跑完检 rc；空输出仅在 rc=0 时读成"没有"。
