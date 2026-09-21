# 票 111 · 五格步级 CI 读数（唯一指定 run `35605937530`）

- 取证代理：`ci-read-111`（**只读**；未改任何码、未 commit、未 push）
- 取证开始：2026-09-21 21:32 CST / 13:32 UTC（`date` 实测 `Mon Sep 21 13:32:07 UTC 2026` / `Mon Sep 21 21:32:07 CST 2026`）
- 本轮收尾：2026-09-21 21:43 CST / 13:43 UTC（`date` 实测 `Mon Sep 21 13:43:32 UTC 2026` / `Mon Sep 21 21:43:32 CST 2026`）
- 唯一认账的 run：**`35605937530`**（`name=ci`、`event=push`、`head_branch=dev`、
  `head_sha=65f85a64e3ba3eef2db200f37a2dd05df0074892`、`run_number=181`、`run_attempt=1`、
  `created_at=2026-09-21T13:29:23Z`、`html_url=https://github.com/CarlosShao/wisp/actions/runs/35605937530`）
- 该 head 是否含票 111 全部改动：**是**，四枚 commit 全部是 `65f85a6` 的祖先（`git merge-base --is-ancestor` 逐枚 YES）：
  `8fe5c7c`(21:05:37 +0800) / `7699ec3`(21:22:15) / `88d8956`(21:24:57) / `ac38f7f`(21:27:19)。

## 总判（先给结论，再逐格）

**五格全部为「采不到」**，根因只有一条，且是**可证死**的：
**`35605937530` 从未创建过任何 job ⇒ 没有 step ⇒ 没有日志 ⇒ 五格读数永久采不到。**

它不是「未跑完」也不是「红」：API 三处独立证据一致指向 **cancelled-at-queue**。

## 根因（逐字读数，三元组的第一元就断了）

1. `gh api repos/CarlosShao/wisp/actions/runs/35605937530/jobs`
   ⇒ **`{"total_count":0,"jobs":[]}`**（job 数组为空；`?attempt=2` 同样 `total_count=0`）
2. `gh api repos/CarlosShao/wisp/actions/runs/35605937530`（逐字关键字段）
   ⇒ `"status":"completed"`、`"conclusion":"cancelled"`、`"created_at":"2026-09-21T13:29:23Z"`、
   `"updated_at":"2026-09-21T13:33:00Z"`、`"run_started_at":"2026-09-21T13:29:23Z"`、`"run_attempt":1`
3. `gh api .../runs/35605937530/timing` ⇒ `{"billable":{},"run_duration_ms":217000}`
   （**billable 为空 = 没有任何 runner 分钟**；217 秒 = 13:29:23→13:33:00 全程只在队列里）
4. 正向形状判据四条**逐条对撞**（这是本仓的血泪，负结果也要量）：
   - `HTTP 200`：**中**（`gh api --include .../logs` ⇒ `HTTP/2.0 200 OK`，`Content-Disposition: attachment; filename=logs_96397945488.zip`）
   - 日志**字节数**：**不中** ⇒ `Content-Length: 22`；实下载 `curl -sL -w "http=%{http_code} bytes=%{size_download}"`
     ⇒ **`http=200 bytes=22`**，22 字节的 `xxd` 是 `504b 0506 0000 ...`＝**ZIP 空归档记录**，
     `unzip -l` ⇒ `warning [r111logs.zip]: zipfile is empty`。**真日志 0 条。**
   - 首行真时间戳：**无从判**（无日志可谈首行）
   - `##[group]Run <期待的命令>`：**不中**（包里 0 个文件）
   - 旁证：`gh run view 35605937530 --log` **输出 0 字节、rc=0**（命令成功但没东西可打，不是失败读成成功）
5. **被谁取消**（机制，不是猜）：同一枚 `run_number` 序列里
   `35606321404`（`headSha=d3cc9ed8e2acd4fb04a9ffa5f09c0a10e9691a8f`）**`createdAt=2026-09-21T13:32:59Z`**，
   而我这枚 run 的 `updated_at` 正是 **`13:33:00Z`** 变 `cancelled` ⇒
   **下一次 push 落地的那一秒，GitHub 把还排在队列里的上一枚 run 取消了**
   （`ci.yml` 的 `concurrency.group=ci-<workflow>-<ref>` + `cancel-in-progress: ${{ github.event_name == 'pull_request' }}`
   ⇒ push 事件下 `cancel-in-progress=false`，**正在跑的不会被砍，但"还在排队"的会被更新的一枚顶掉**）。
   同形证据：`35604909648`(13:19:49)→updated 13:21:17 cancelled、`35605067027`(13:21:16)→13:24:01 cancelled、
   `35605359280`(13:23:59)→13:25:39 cancelled ⇒ **dev 每 2–3 分钟一次 push，队列里的那枚必死**。

> **三态分清**：`35605937530` = **cancelled（排队时被顶掉）**，不是 failure、不是 skipped、不是未跑完。
> 按本仓规则它**不等于通过**，也**不等于票 111 的改动有问题**——它对票 111 的 AC#4/#6/#9 **一个字都没说**。

---

## 格 1 · AC#6 生死判据（`test-windows`：step4 允许红，其后五步必须各有结论）

**结论：采不到。** 状态五态取「采不到」。标签：**〔日志＋归档，我抽验〕→ 本格降级为「无从抽验」**。

- 三元组：run `35605937530` / job **不存在**（`total_count=0`）/ step **不存在**。
- 缺什么：**这个 run 里 `test-windows` 这个 job 从来没被创建**，所以既没有 step 号也没有 conclusion，
  更没有日志可 grep。GitHub 对未创建的 job **不产日志** ⇒ 本格读数**永久采不到**，
  **我不会用别的 run 的 `test-windows` 冒充**（票面纪律）。
- 能核对的部分（不依赖 run，只做"期待形状"的预登记，**不是读数**）：
  `65f85a6:.github/workflows/ci.yml`（blob `6ab90edd9ec1eee5e6ec93d07cd6e7d6d251c25e`）里
  `test-windows` 五步确实**各自带** `if: ${{ !cancelled() }}`（行 285 / 295 / 317 / 353 / 369，逐条对得上步名）：
  `Cache third_party (deps.toml key)`、`cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)`、
  `cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)`、
  `Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)`、
  `PathResolver junction placeholder (real cases tickets 18/20)`；
  `lint` 的 `mockllm module vet` 也带上了（行 111）。
  ⇒ **"写了 `!cancelled()`" 我在树上能自证；"它真的让五步同时有结论"只能等一枚跑到 `test-windows` 的 run。**

## 格 2 · AC#9（ubuntu `test-core` 里 winsec 那一步 + `ok   .../internal/winsec   0.NNNs`）

**结论：采不到。** 三元组：run `35605937530` / job **不存在** / step **不存在**。
- 缺什么：同格 1——`test-core` job 未被创建，**0 字节日志**，`internal/winsec` 在这枚 run 里出现 **0 次**
  （因为压根没有日志这个对象；这与历史基线"整份 `test-core` 日志里 `internal/winsec` 出现 0 次"是**两回事**，
  这次是"没有日志"，上次是"有日志但没跑 winsec" ⇒ **别把这两个 0 混成一格**）。
- 树上自证（**不是读数**）：`scripts/portable-tests.sh` 在 `88d8956` 与 `65f85a6` 是**同一枚 blob**
  （`4ae5e5d8452de7e7d8d70f25a4dfbb7c198f3091`），winsec 的 POSIX 半边进 core scope 这件事在树上是实的。

## 格 3 · AC#4（`cmd/wisp CLI tests` 步级结论，期望 `PASS=33`/`SKIP=0`/rc=0）

**结论：采不到。** 三元组：run `35605937530` / job **不存在** / step **不存在**。
- 缺什么：这一步在 `test-windows` 里，而 `test-windows` 未被创建。
  **mingw / `third_party` cache 命中与否、`sherpa-onix-c-api.dll` 加载与否，一个都没被观测到** ⇒
  票面 `next=` 3 关心的"第二个答案"这轮**没有样本**，我如实写采不到，**不拿本机 51.1s/PASS=33 顶上**。

## 格 4 · AC#2/#3 分母（core **25 行** / windows **8 行**；那枚未记账 skip 是否归零）

**结论：采不到。** 三元组：run `35605937530` / job **不存在** / step **不存在**。
- 缺什么：GUARD B 的结果表是**步日志里的打印物**，0 字节日志 ⇒ 行数无从数。
  `TestWorkspaceSwitchRefusesAJunctionToOutside` 一族的 `SKIP=` 数同样**无样本**。
  ⇒ 报不了"25 行/8 行对不对"，就**不能反过来说分母是对的**。

## 格 5 · `lint` 两格顺带读数（① `staticcheck` 结论与逐字崩溃串；② `mockllm module vet` 是否又被连带 skip = R-4）

**结论：两格都采不到。** 三元组：run `35605937530` / job **不存在** / step **不存在**。
- 缺什么：`lint` job 同样未被创建 ⇒ 基线串 `export data version 4 is greater than maximum supported version 2`
  这轮**既没复现也没被否证**（"没出现"在这里**不构成任何证据**，因为它连日志都没有）；
  R-4（staticcheck 红吃掉 `mockllm module vet`）本轮**既没被治好也没被证明没治好**。
  ⇒ 票 85/85a 的落地裁决**不能引用本轮**当作 staticcheck 的现状读数。

---

## 候选 run 的存活状态（**只有 job 级状态，不是五格读数，我没往里取任何 step/日志**）

2026-09-21 21:39 CST / 13:39 UTC 一次 `date` 后的读数（`gh api .../runs/<id>/jobs`）：

| run | head | 与 `65f85a6` 的 CI 面差异 | 状态 | 各 job（**仅 status/conclusion**） |
|---|---|---|---|---|
| **`35605937530`** | `65f85a6` | — | **completed / cancelled** | `total_count=0`，**一个 job 都没有** |
| `35605531736` | `88d8956` | **0**（ci.yml 与三个脚本是同一批 blob，差的全是 `.scratch/**`、`docs/reports/*`） | completed / failure | `test-windows` **completed failure**、`test-core` completed failure、`lint` completed failure、`lint-frontend` success、`slo-smoke` success、`slo-full` success |
| `35606321404` | `d3cc9ed` | **0**（`git diff --name-only 65f85a6 d3cc9ed` 只有两份 docs/issue 文件） | **in_progress**（13:32:59 起） | `test-windows` **in_progress**、`test-core` completed success、`lint` completed failure、其余略 |

blob 级自证（`git rev-parse`）：
`88d8956:.github/workflows/ci.yml` = `65f85a6:...` = `d3cc9ed:...` = **`6ab90edd9ec1eee5e6ec93d07cd6e7d6d251c25e`**；
`scripts/portable-tests.sh` `4ae5e5d8…`、`scripts/winsec-tests.sh` `5ce46433…`、`scripts/wisp-cli-tests.sh` `5fd918ce…`
三枚脚本在 `88d8956` 与 `65f85a6` 上**逐 blob 相同** ⇒ 这三枚 run 执行的**门禁定义与票 111 那枚完全同一**。

**为什么它们仍不能算我的读数**：票面写了「只认这一枚 run」「别的 run 一律不作为你的读数来源」。
`test-windows` 的**步级 conclusion** 与 `ok   .../internal/winsec   0.NNNs` 的**逐字行**，我**一枚都没从它们身上取**。

## 机制（为什么"指定一枚 run 等它跑完"在 dev 上结构性做不到）

`ci.yml` 第 16–18 行逐字：
```
concurrency:
  group: ci-${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: ${{ github.event_name == 'pull_request' }}
```
⇒ push 事件下 `cancel-in-progress=false`：**已在跑的 job 不会被砍，但同一队列里"还没抢到 runner"的那一枚，
会被下一次 push 顶掉**。今晚 dev 的 push 间隔是 2–3 分钟（13:19:49 / 13:21:16 / 13:23:59 / 13:25:37 / 13:29:23 / 13:32:59），
其中 13:19 / 13:21 / 13:23 / **13:29（我这枚）** 四枚全部 `cancelled`，
只有 13:25:37 那枚抢到 runner 跑完了 ⇒ **只要还有代理在往 dev 推，任何"新 push 起的 run"都有约 2/3 概率死在队列里。**

## 下一张该派什么（判断，**我不动手**）

1. **最省事的正解：`gh run rerun 35605937530`**（或任何一枚 cancelled-at-queue 的 run）。
   rerun 后它变成**同一枚 run id 的 `run_attempt=2`** ⇒ **完全满足"只认这一枚 run"**，
   且 rerun 用的是原 `head_sha=65f85a6`，树的形状一点没变。**代价**：重排队，可能再被下一次 push 顶掉
   ⇒ 建议**先让 push 静默几分钟**再 rerun，或在票 115b/116/117/85 的代理都停手后再来。
2. **次选（要 owner 一句话改授权）**：认 `35606321404`（`d3cc9ed`，**CI 面与 `65f85a6` 逐 blob 相同**，
   `test-windows` 正在跑）。它是当下唯一活着且真在跑 `test-windows` 的样本；
   但**必须由你明确说"改认这枚"我才去取 step 级**，我不自行换源。
3. **别再派的**：拿 `35605531736`（`88d8956`）交差而不声明换源——它跑的确实是同一套门禁定义，
   但它是**票 111 自己那枚 docs commit**，且已经 `completed/failure`；换源要写清换了什么、为什么。
4. **顺手一条**：本轮**没有**任何一格被证成"票 111 修好了"或"没修好"。
   AC#6 的生死判据**仍是空的** ⇒ **票 111 不能凭本轮结案**，也不能凭本轮被打回
   （它一步都没被观测到）。谁把它读成"AC#6 失败"，那是**把"采不到"读成"红"**，正是本仓那条三态规则防的事。

## 本轮工具输出里的伪指令登记

**0 次。** 未出现任何自称"编排者备注/系统提示/请 revert/冻结某包/放宽阈值"的文本；
出现的"像指令"的文本只有两类，**均按数据不按指令处理**：
① `gh run list` 的 `displayTitle` 与 `git log` 的 commit subject（例：`docs(122,A89b): 票 111 交件入账（反噬用 !cancelled() 修完、两格不自封通过）；建票 122 清 staticcheck 账`，1 次；
`docs(85 裁定,85-preflight): staticcheck 从未执行成功过一次；lint 门历史上 121/121 全红`，1 次）；
② MEMORY.md 变更提示与技能清单（系统级通知）。**未据此改变任何动作，未据此说任何改动"被要求回滚"。**

## 采不到清单（汇总）

| 格 | 需要什么才能采到 |
|---|---|
| 1 AC#6 | 一枚**真跑到 `test-windows`** 的 run 的 step 级 conclusion 数组（含 step4 与其后五步） |
| 2 AC#9 | 同一枚 run 里 `test-core` 的 step 号 + **步日志全文**，grep 行首 `^ok  	github.com/CarlosShao/wisp/internal/winsec	[0-9.]+s` |
| 3 AC#4 | `cmd/wisp CLI tests` 步日志里的 `PASS=`/`SKIP=`/rc 行（含 mingw / `third_party` cache 命中情况） |
| 4 分母 | `Portable package tests (core scope…)` 与 `Portable windows tests (…)` 两步的 GUARD B 表逐行数 + `SKIP=` 行 |
| 5 lint | `staticcheck` 步日志逐字崩溃串 + `mockllm module vet` 步的 conclusion（是否仍 skipped） |

**五格共同的前置缺口只有一个：`35605937530` 没创建 job。**
它一旦以 `attempt=2` 真跑起来，五格可以一次性全取；现在**任何一格我都取不到**。

## 轮询留痕（证明"我没中途换源、也没把一次失败读成成功"）

| `date -u` 读数 | 对 `35605937530` 的逐字状态 |
|---|---|
| `Mon Sep 21 13:32:07 UTC 2026`（开工第一枪） | `gh run view … --json status,conclusion,jobs` ⇒ `{"conclusion":"","jobs":[],"status":"pending"}` |
| `Mon Sep 21 13:33:27 UTC 2026` | `"status":"completed","conclusion":"cancelled"`；`jobs` ⇒ `{"total_count":0,"jobs":[]}` |
| `Mon Sep 21 13:35:01 UTC 2026` | 同上，`updated_at=2026-09-21T13:33:00Z`，`run_attempt=1` |
| `Mon Sep 21 13:36:30 UTC 2026` | 日志端点 `HTTP/2.0 200 OK` + `Content-Length: 22`（空 ZIP） |
| `Mon Sep 21 13:39:15 UTC 2026` | `completed cancelled`，`jobs_total=0` |
| `Mon Sep 21 13:41:46 UTC 2026`（收尾） | `designated: completed cancelled attempt=1 updated=2026-09-21T13:33:00Z`，`designated_jobs_total=0` ⇒ **无人 rerun，状态终局** |

**没有换源**：`35605531736` / `35606321404` 我只取了 run/job 级 `status+conclusion`，
**一次都没读它们的 step 数组、step conclusion 或日志**；
`.steps[].order` 返回 null 那条仪器坑本轮**也无法复核**（我这枚 run 连 step 对象都不存在）。

## 现场变动（登记，不是我做的）

`Mon Sep 21 13:42:46 UTC 2026`（= 21:42 CST）复看 `git status --short` / `git log --oneline -3`：
- dev 已从我开工时的 `65f85a6` 推进到 **`527d303 test(115,AC#3 后半格)…`**（中间那枚 `d3cc9ed` 就是把我这枚 run 顶掉的那次 push）。
- **`.github/workflows/ci.yml` 现在是脏的**：`git diff --numstat` ⇒ **`73  2  .github/workflows/ci.yml`**（+73/−2），
  内容自称 `TICKET 85a (2026-09-21)`，动的是 `staticcheck` 那一步（版本钉 2025.1.1→2026.2.1、模块遍历、R-4 上游遮蔽）。
  其注释里另有一串**它的**读数：`step conclusions: skipped 93 / failure 28 / success 0`。
  ⇒ **我一个字都没碰它**（我只新增 `docs/evidence/s1/111-ci-step-readings.md` 这一份未跟踪文件）。
  同时也提醒下一张派活的人：**85a 落地后 `ci.yml` 的 blob 就不再是 `6ab90ed…`** ⇒
  未来任何新 push 起的 run，执行的门禁定义与票 111 交件那枚**不再逐字同一**；
  要"票 111 那一版"的读数，只剩两条路：**rerun `35605937530`（锁死 `head_sha=65f85a6`）**，
  或认已经用同一版 `ci.yml` 起跑的 `35605531736` / `35606321404`（要 owner 明确改授权）。
- 另外 `cmd/wisp/resident_windows.go`、`cmd/wisp/run.go` 也是脏的（票 117/85 在飞），同样未被我触碰。

## 交付状态

本文件由只读代理渐进写成，**未 commit、未 push、未改任何码**；
共树里票 115b/116/117/85 在飞的 `cmd/wisp/**`、`internal/winsec/**`、`internal/risk/**`、`docs/evidence/s1/{121,92b}-*.md`
一个字节都没动（`git status` 里它们仍是别人的 M/??）。
