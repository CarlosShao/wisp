# 309-a1 · 00 起手锚（只读普查腿，格＝`AC#0`）

- 时刻：`Sun Oct 11 11:46:45 CST 2026`（尺＝`date`，本机时区 +08）
- HEAD 全号：`7dc753e37e894d5c7fcc2e560b68e0731a0eeb13`（尺＝`git rev-parse HEAD`）
- 分支：`dev`（尺＝`git rev-parse --abbrev-ref HEAD`）

## git log --oneline -6（尺＝`git log --oneline -6`）

```
7dc753e3 台账 A842＝收 `taskid-a1` 并结掉 `A838` §3 我自己登记那枚欠读数（`bridge.go:925` 空任务号那一肢）＋本轮编排者件与门禁读数
c9e20a1e 票 309 新立＝`assessorFor` 空任务号那一肢静默摘掉 R4 探测器（六格全未勾，AC#0 未裁形之前⛔ 派落地腿）＋本轮编排者件与门禁读数
79e4974d 72-r1 class B: add negative nails to the two 307-debt sync fallback cases
f278c57f 72-r1 class A: move 7 expectation sides off hardcoded temp spellings onto verdicts
a9ddc020 台账 A841＝车队卫生第二笔：台账引的 275 枚证据件名里 12 枚只在盘上从未提交＋2 枚"查无此件"逐枚具名（⛔ 删 ⛔ 搬 ⛔ 补提交别人的件）
93643a29 probe(taskid-a1): read-only census of the empty-taskID limb at bridge.go:925
```

## 工作树干净度（三态判据）

- 尺＝`git status --porcelain -- internal cmd | wc -l` → **0 行**。
  ⇒ 起手时刻 `internal/**` 与 `cmd/**` 的工作树 **与 HEAD 逐字节一致**（`72-r1` 的两笔 `f278c57f`／`79e4974d` 已进 HEAD，不在未提交态）。
  ⇒ 本腿读 `internal/**`／`cmd/**` 一律走 HEAD blob（`git show HEAD:<path>`／`git grep <pat> HEAD`），三态一律记作 **HEAD blob＝工作树**（0 行差异是实测而非假设）。
- 尺＝`git status --porcelain | wc -l` → **847 行**（全仓未跟踪/未提交件总数，属别家的活与临时件，本腿 ⛔ 动）。

## 我跑过的 Go 编译面命令（硬约束 1 具名自报）

- `go env`：**0 条**（未跑）
- `go list`：**0 条**（未跑）
- `go build`／`go vet`／`go test`／`gofmt`／`gofumpt`：**0 条**（⛔ 禁，未碰）
- ⇒ 本腿没有任何实测色读数；"哪枚既有钉会红"全部由静态尺（读断言正文＋它断言的对象）给出，拿不到的写进 `90-final.md` 欠读数清单。

## 写面声明

- 只写 `.scratch/wisp/probes/309/a1/**`（本目录）；证据件一律 `.md` 或 `logs/` 下非 `.out` 后缀（根 `.gitignore` 第 8 行是全仓 `*.out`）。
- ⛔ 改源码／测试／票面／台账／`docs/**`；⛔ 新建 `.go`；⛔ push。
