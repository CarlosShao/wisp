# ledger-hash-sweep-1 · 起手锚

- 腿名：`ledger-hash-sweep-1`（只读对账腿：A720–A736 区段的 commit 号 / 文件路径机械核对）
- `date`：2026-10-08 19:01:45 +0800
- HEAD：`5f52f310f1940a9cdb4b043f2024d8d71eafee7f`（`5f52f310`）2026-10-08 19:01 `A736 落账：收 gate-snapshot-1…`
- 分支：`dev`
- `git status --porcelain | head -20`（只登记）：
  - `M .gitignore`
  - `M .scratch/wisp/probes/152/my152.py`
  - `M .scratch/wisp/probes/161/r6/logs/flip-1.txt` …（flip-2..6、flip-baseline、flip-restored 共 8 枚）
  - `M .scratch/wisp/probes/242/r3/logs/probe-routed.txt`
  - `M .scratch/wisp/probes/268/v1/evidence.md`
  - `M cmd/wisp/subagent_selfapproval_197_test.go`
  - `D design/assets/base.css` / `design/assets/icons.js` / `design/assets/theme.js` / `design/assets/tokens.css`
  - `M design/doubao/README.md` / `M design/doubao/demo/app.js` / `M design/doubao/demo/index.html`
  - （以上只登记，不处置；`259-r4` 腿在 `cmd/wisp` 写测试中）

## 目标区段定位

- `grep -n '^## A7[2-3][0-9]' docs/reports/pending-and-issues.md`：
  - A720 节头 = `14093`；A736 节头 = `14389`；文件共 `14393` 行（A736 为末节，14389–14393）。
- 核对区段 = `14093,14393`。

## 纪律

- 零 Go 命令；只用 git log/cat-file/ls-tree/show + wc/grep/sed/date。
- 写面＝只新建本目录 `*.md`；不改任何文件；零翻框；不 push。
- 临时件只建不删、建仓库外。
- 范围尺不接 `2>/dev/null`；跑完必检 rc；空输出只有 rc=0 才读成"没有"。
