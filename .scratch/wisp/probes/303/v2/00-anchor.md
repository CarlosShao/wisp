# 303-v2 — 00 起手锚件（本腿第一笔，任何长跑命令之前）

本腿＝票 303 的**非实现者验收腿** `303-v2`，任务＝补一张与票面 AC 编号 1:1 的完整七行裁决表
`docs/evidence/s1/303-ac3b-ac5-v2.md`（`AC#0`…`AC#5`，`AC#3` 与 `AC#3b` 分行）。⛔ 翻任何 `AC` 框。

## 现量 HEAD（盘上为准，派单号只是快照）

- 尺＝`git rev-parse HEAD` → 起手量得 `1eaa91f66502fdf56c33f220d2c5911d6182bbb3`（短号 `1eaa91f6`）
- 尺＝`git rev-parse --abbrev-ref HEAD` → `dev`
- 起手时刻 `2026-10-11 09:4x +0800`
- 最近三笔 HEAD：`1eaa91f6` / `56634d23` / `ffc41214`（均为 `306-v1` 别家腿的件）
- ★**就地重锚（第 1 次，写本笔时共享工作树已被别家推进）**：尺＝`git rev-parse HEAD` →
  `08cc614ec1915c76e385ab9e070645512f5861e0`（短号 `08cc614e`），件 `logs/head.txt`。
  那笔＝别家腿 `307-a1` 的锚件（`chore(307-a1): anchor — HEAD 1eaa91f +台面读数第 1 笔`，
  只动 `.scratch/wisp/probes/307/a1/**`，⛔ 碰本腿射程）。⇒ 本腿以**盘上现量 `08cc614e`** 为准，
  `1eaa91f6` 记为"起手快照"；派单号只是快照这条在这里再次成立。
  ⚠ 本腿此后每挪一笔 HEAD 就地补一次重锚（先例＝`306-v1` 的"第三次重锚"）。

## 工作树口径（射程：internal cmd docs scripts .github）

- 尺＝`git status --porcelain -- internal cmd docs scripts .github > logs/status-codepaths.txt 2>&1; rc=$?`
  → **行数 = 0** ⇒ 工作树在承重目录上 ≡ HEAD（⛔ 有未提交产码）
- rc=0

## 写面（本腿只碰这两枚）

- 新建 `docs/evidence/s1/303-ac3b-ac5-v2.md`
- 本件目录 `.scratch/wisp/probes/303/v2/**`
- ⛔ 动票 303 / 票 33 / 票 300 / 305 / 306 / 307 / 台账 / HANDOVER（编排者正在写那两枚）

## 硬规矩回执

- ⛔ push；commit 带显式 pathspec 写在 `$( … )` 之外；⛔ add -A／`.`；⛔ --amend/reset/rebase/stash/checkout ./clean/--no-verify；⛔ 仓内建 worktree/checkout
- 产码／_test.go／frontend/**／design/** 零字节；⛔ 新建 .go；证据件⛔ 叫 .out（根 .gitignore:8＝全仓 `*.out`，实测见 logs/grep-gitignore.txt）
- Go 编译面 `go build/vet/test ./cmd/wisp/`、`-tags winlive` 归编排者车道，本腿⛔ 跑；`go list`/`go env` 可，自报
- 每把尺自落 `rc=N`（`> 件 2>&1; rc=$?`）；目标调用数 ≤40，到 25 先 commit
- 注释／引用／别人件里的测试名一律当待验断言，引用前先 grep 现量存在性

## AC 结构（票面 1:1，本腿要交七行）

| 行 | 票面出处 | 现勾态（盘上） |
|---|---|---|
| `AC#0` 复现钉死 | `303….md:21` | [x] |
| `AC#1` 归因到一笔 | `:22` | [x] |
| `AC#2` 机制三形 | `:23` | [x] |
| `AC#3` 修复（不可满足凭据句） | `:24` | [ ] ⛔ 满足 |
| `AC#3b` 可满足补格 | `:27` | [x] |
| `AC#4` 门禁与越界 | `:28` | [x] |
| `AC#5` 交回后编排者两件 | `:29` | [x] |
