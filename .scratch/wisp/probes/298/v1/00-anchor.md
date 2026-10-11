# 298-v1 — 开工锚件（非实现者验收腿；⛔ 本腿零翻框、零实现修复、零 push）

钟点尺＝`date '+%Y-%m-%d %H:%M:%S %z'` 的 stdout：`2026-10-11 12:00:15 +0800`
工作目录＝`/d/work/workspace/projects plans/Wisp`（win32 / git bash）

## G1 双锚

尺＝`git log -1 --format='%h %ad %s' --date=iso-strict`

    27059960 2026-10-11T11:57:39+08:00 probe(309-a1): AC#0 census …
    rc=0

⇒ 起手 HEAD＝`27059960`（⚠ 共享工作树，HEAD 会随别人的腿漂；本腿一切读数按"对象锚"钉：
落地腿四笔＝`24ae75b7`（锚）／`18d43f8e`（唯一碰 `.go`）／`c14fd204`（凭据）／第 4 笔（30-gates＋票面追加，
号码由 `git log --format=%h -1 -- .scratch/wisp/probes/298/r1/30-gates.md` 现取，逐枚具名见 `10-ac-readings.md`）。
`18d43f8e` 的父＝`48e7b0a4`＝"改前那一棵"，AC#0/AC#1 的 blob 全部从 `48e7b0a4:<path>` 取。）

## G2 停手旁证（票面 AC#1 起手那条，验收腿同尺复跑）

尺＝`git status --porcelain -- cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go`

    （无输出）
    rc=0 lines=0

尺＝`git status --porcelain -- cmd/wisp | wc -l` ⇒ `0`（整个 cmd/wisp 工作树干净＝落地腿改动已全部入库）。

## G3 进程尺（起手一发；此后每一发 go test 前再逐发跑）

- 尺＝`tasklist //FI "IMAGENAME eq wisp.exe" | grep -ci '^wisp\.exe'` ⇒ **0**（原文＝`INFO: No tasks are running…`）
- 尺＝`tasklist //FI "IMAGENAME eq balldebug.exe" | grep -ci '^balldebug\.exe'` ⇒ **0**

## G4 共享索引现量（⛔ 本腿不动别人的暂存件）

尺＝`git diff --cached --name-only | wc -l` ⇒ **20**，名册全部在 `.scratch/wisp/probes/orch/20261011/**`
（编排者自己的暂存件）。⇒ 本腿每一笔 commit 前重数这一把，commit 一律**显式 pathspec**，⛔ 带出别人的 20 枚。

## G5 起手基线（与落地腿"改后"声称对拉的起点）

- 尺＝`gofmt -l cmd/wisp`（工作树）⇒ 逐字 `cmd\wisp\models.go` 一枚，`rc=0`
  ＝落地腿 `30-gates` 表声称的"改后 1 枚（＝假枚）"同形；**不是 3 枚**＝票面现量节那把工作树尺已被 `18d43f8e` 改掉两枚。
- 尺＝`test -x "$(go env GOPATH)/bin/gofumpt.exe"` ⇒ `yes rc=0`（⛔ 裸 `gofumpt` rc=127 那一形不存在，门禁尺可用全路径）。
- 尺＝`find .scratch/wisp/probes/298/r1 -type f -size 0 | wc -l` ⇒ **0**（落地腿"无 0 字节件"的声称初核；逐件 rc 行复验见 `10-ac-readings.md`）。

## 写面声明（本腿自缚，＝AGENTS §1.5 三处落点）

`.scratch/wisp/probes/298/v1/**`（本件＋logs/，临时件只建不删、⛔ 不落仓根、⛔ 不叫 `.out`）＋
`docs/evidence/s1/**`（裁决表一枚，命名照齐 `ls docs/evidence/s1` 现风格）＋
票 298 末尾一条 append-only bullet（1 加／0 删）。
⛔ 任何源码／测试／`.gitattributes`／git 配置／`frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt`。
⛔ 翻任何一格框。
