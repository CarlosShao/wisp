# 111-ci1 — 推送后两发 CI run 的步级读数（只取数，不裁框）

> 本件是**只读取证腿** `111-ci1` 的交件。范围＝把票 111 里逐字要求"CI 侧真读数"的那几格
> （AC#1 第④列 / AC#3 前提今判 / AC#6 skipped 全名册 / 常红归因）喂到**步级**。
> **本腿不翻任何 `- [ ]` 框、不改票面、不裁 AC 的生死**；每张表末尾的〔建议位〕只是给编排者的理由，
> 判语一律留给编排者。

---

## §0 起手锚（三把尺原文）

| 尺 | 命令原文 | 读数 |
|---|---|---|
| HEAD | `git rev-parse --short HEAD` | `c98fc5db` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev` |
| 时刻 | `date '+%m-%d %H:%M'` | `10-06 11:46` |
| 自证零 go 命令 | 见下 | 本腿全程未调用 `go test` / `go build` / `go vet` / `go run` 任何一枚；枚数与行数一律走 `git ls-files`／`wc`／`sed`，CI 侧读数一律走 `gh`。go 使用权整轮留给写腿 `232-r2`。 |

**本腿的尺原文清单（可复量）**
- `gh run list --branch dev -L 8`
- `gh run view 37405698188 --json jobs -q '...'`（步级 conclusion）
- `gh run view 37406757402 --json jobs -q '...'`
- 计数／名单：`git ls-files` 系列、`wc -l`、`sed -n`、`grep -n`
- ⛔ 无 `go`／无 `go test`／无 `go vet`

**run 号存活确认**（`gh run list --branch dev -L 8` 原文摘录，10-06 11:4x）
```
completed  failure  ## A634                                              ci  dev  push  37406757402  9m36s  2026-10-06T02:58:39Z
completed  failure  docs(ticket-269): 把已撤销就地标到 Status 行         ci  dev  push  37406422380  9m29s  2026-10-06T02:54:30Z
completed  failure  ## 5. 编排者翻勾记录（2026-10-06 10:3x +08…          ci  dev  push  37405698188  8m34s  2026-10-06T02:45:36Z
completed  failure  ci                                                 ci  dev  sched   37396530365  9m6s   2026-10-06T00:55:36Z
completed  success  slo-fresh                                 slo-fresh dev  sched   37392522160  15s    2026-10-06T00:09:51Z
```
- 主口径两发（任务给的）：`37405698188`（推送后新配置第一发）、`37406757402`（headSha `cc315261`，推送后第二发）。
- **推送后其实有第三发** `37406422380`（10-06 02:54:30Z，夹在两发之间）＝任务前文未列；本件主表仍按任务口径只用那两发，
  第三发只在它能让某一问的读数更硬时才补，补了会在对应小节具名。
- 推前旧配置对照发：`37396530365`（schedule，10-06 00:55:36Z）＝§5 作差的左端。

---

## §1 票面逐字（问 1）

<!-- 待填：AC#1/2/3/4/6/7/8/9/10 十枚框原文，sed 行号已核对与否要具名 -->

---

## §2 AC#1 第④列名册（问 2）

<!-- 待填：包 × (有无测试文件 / 在不在 CI 某一步 / 那一步真给过结论的 run+job+step)，含 census 步打印的四数逐作业 -->

---

## §3 GUARD D 与空分母（问 3）

<!-- 待填：session/watchdog 今日跟踪的 _test.go 枚数、在册与否、census success 的成因（步级证据 + portable-tests.sh 退出码路径） -->

---

## §4 skipped 全名册与两分（问 4，自我对抗节 A）

<!-- 先写满：两发 run 全部 skipped 步逐枚；①被前一步失败吃掉的实质测试步 ②setup post 清理步；谁承载读数谁不承载 -->

---

## §5 红名作差（问 5）

<!-- 待填：test-core step7 / lint step8+step12 staticcheck / test-windows step7+step9 红句逐名，与 37396530365 同口径作差 -->

---

## §6 判不动的地方（问 6，自我对抗节 B）

<!-- 先写满：真写内容；拿不到的读数要具名到 run+job+step，不许用本地或旧归档顶替 -->
