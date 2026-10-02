# 167-r1b — 面板三枚 Go 侧接线：占用条／排队序号／停止（写腿交件）

> 本文件由写腿 `167-r1b` 产出。票面：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md`（AC 框一枚未碰、票面原句一字未改）。
> 前程 `167-r1` 预检中途断流，零 commit 零证据件，无可继承——本件从零起。
> 只读取证腿 `167-a1` 的普查件 `.scratch/wisp/probes/167/a1/census.md` 已完整读过，本件的落地面以它为底（引用处具名）。
> 写面＝`cmd/wisp/**`＋`internal/panel/**`（测试为主）；`internal/panel` 写面只到票 145 已批解冻范围的"新增字段"为止。

## 0. 起手锚（逐字读数）

`date -Iseconds`：

```
2026-10-02T20:29:11+08:00
```

分支（`git branch --show-current`）：

```
dev
```

`git log -1 --format='%H %ad' --date=iso`：

```
11a9c510e365a3e571d9a65c10c9600d32ae4ccf 2026-10-02 20:07:56 +0800
```

`git status --porcelain -- cmd/wisp internal/panel`（起手读数，逐字）：

```
（空）
```

同机在飞：`242-r1`（写 `internal/agent/approval` 测试）。凡跑 `cmd/wisp` 测试前均先复尺 `git status --porcelain -- cmd/wisp`，测试严格串行。

## 1. 三枚修法与判据（带 file:line）

（骨架占位——取数后逐枚填：修法、判据、每枚"有／无／有但没接"的现量结论。）

## 2. 门禁四数

（骨架占位——build／gofumpt -l／`go test ./cmd/wisp ./internal/panel -count=1`（带 sherpa PATH）／d22scan exe。）

## 3. 变异自证表

（骨架占位——每发：原档 backup 台件、突变内容、预期红、实测红句逐字、还原后 porcelain=0。）

## 4. 我可能写错的条目

（骨架占位。）

## 5. 判不动的地方

（骨架占位。）

## 6. 交件判语

（骨架占位。）
