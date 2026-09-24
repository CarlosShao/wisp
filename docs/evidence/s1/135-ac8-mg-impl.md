# 票 135 AC#8（家族第 7 发 `M-G`）实现方证据 —— 变异自证三态 ＋ 门禁账

**执行方**：`worker-ticket135-ac8-r2`（**接续程**，非 `worker-ticket135-ac8-mg`）
**被验版本**：commit **`8369b24`**（`8369b24e20eb85b5c9afb45198499dc4b9d9d804`，含上一程的 `fa35557`）
**本格判据来源**：票 135 面 AC#8（09-24 00:0x 编排者追加）＋ 票 133 第二任验收方的
`docs/evidence/s1/133-ac2-r2-acceptance.md` §1.11／§1.12／§3.3／§3.4（`R-133-9`＝`M-G` 的出处）
**本程地界**：一枚新建证据文件（本文件）。**生产码零接触、上一程那把尺零改动**（§1.4 字节为证）。
**本程不翻任何勾**（AC#8 由编排者按本表翻）。

## 0. 争用闸门（每一批读数之前跑一次，逐轮原样登记）

### 0.1 为什么本程把这道闸门当第一条规矩

这台机器上 **self-hosted CI runner 与编队同机**：`.github/workflows/ci.yml:538` 的 `slo-full`
那一枚 job 的 `runs-on` 是 `[self-hosted, wisp-slo]`。编排者在 **10:50:36 +08** 推了一次 dev
（推的正是被验那枚 `8369b24`，`gh run list` 的 `headSha=8369b24e20eb85b5c9afb45198499dc4b9d9d804`），
当场起了一枚 `ci` run ⇒ **"编队里没有别的代理"不等于"机器是空的"**。

### 0.2 固定动作（三行）与本程用的台件

派单写死的三行：

```
tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'
tasklist | grep -iE '^go\.exe|compile\.exe|cgo\.exe'
gh run list --limit 3 --json status,conclusion,name,createdAt
```

本程把它们包成一枚**有上界的轮询**：`/d/tmp/wisp135mg-r2-gate/poll.sh`
（`INTERVAL=45`s × `MAXROUNDS=26` ≈ 19 分钟一枚上界；命中即 `BUSY` 继续等，
全空即 `CLEAR` 退出 0；到点上界仍忙则退出 3 并照实报。
**没有用 `time.Sleep` 糊窗口**——等的是"条件成立"，且带次数与时长上界。
**没有杀过任何 runner 进程、没有改过 `.github/workflows/**`、没有取消过任何 run**：只观察。）

### 0.3 逐轮原始读数

**手工头两轮（本程自己跑的三行原样）**

```
$ tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'   2026-09-24 10:51:50 +0800
Runner.Listener.exe          50052 Console                    1     86,552 K
Runner.Worker.exe            36068 Console                    1     98,024 K      ← Worker 活着 = 远程 run 在本机落地
$ tasklist | grep -iE '^go\.exe|compile\.exe|cgo\.exe'          （无输出，rc=1）
$ gh run list --limit 3 --json status,conclusion,name,createdAt
[{"conclusion":"","createdAt":"2026-09-24T02:50:36Z","name":"ci","status":"in_progress"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:24:40Z","name":"ci","status":"completed"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:14:55Z","name":"ci","status":"completed"}]

$ tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'   2026-09-24 10:55:01 +0800
Runner.Listener.exe          50052 Console                    1     86,936 K      ← Worker 已退，Listener 常驻（本机 runner 空转）
$ gh run list --limit 3 --json status,conclusion,name,createdAt
（同上：仍有一枚 in_progress 的 ci）
```

**轮询逐轮**（全文在 `/d/tmp/wisp135mg-r2-gate/gate-rounds.log`，下面按轮摘 VERDICT 行＋该轮时刻；
每轮的三行完整输出都留在那枚日志里）

```
---- round=1 at=2026-09-24 10:56:55 +0800 ----  worker-hits=0 listener-hits=1 toolchain-hits=0 IN_PROGRESS_COUNT=1  VERDICT=BUSY
---- round=1 at=2026-09-24 10:57:10 +0800 ----  worker-hits=0 listener-hits=1 toolchain-hits=0 IN_PROGRESS_COUNT=1  VERDICT=BUSY
```

⚠ 本节此刻**只有这两轮**：本程后面每一批发作之前都按 0.2 那三行复跑一次，
VERDICT 与该轮时刻记在**对应读数那一节**里（§2 起），不在这里预填。
**读到本文件时若 §2…§7 任一节缺失或未结，那一节的读数就是没采**，不是"采了没写"。

### 0.4 判读规则与本程的处置

- `Runner.Listener.exe` **常驻不算命中**：它是等活的进程，只有 `Runner.Worker.exe` 代表"真有 job 在这台机器上跑"。
  这一条本程按派单原文的形状写（派单第 2 节：「有 Worker 就说明远程 run 活着」）。
- **`gh run list` 里有一枚 `in_progress` 的 ci＝命中**，即使该 run 此刻落在 GitHub 托管的 runner 上
  （本程 10:53 现量到那枚 in_progress run 的 job 分布：`lint`/`slo-full`/`lint-frontend`/`test-core`/`slo-smoke`
  已 completed、`test-windows` in_progress，`runs-on: windows-latest` 即 GitHub 托管）——
  **本程不赌"它不在这台机器上"，只要闸门命中就重采**。
- ⚠ 一条本程自己造出来的噪声要交代：**一旦本程开始 `go test`，第 2 行的 `^go\.exe|compile\.exe` 必然命中**
  （那是本程自己的编译）。⇒ 闸门只在**每批读数开始之前**跑一次，批次内的命中不属争用。
- **本程没有把任何一发读数发落在已知 BUSY 的轮次里**（§1.2 起每一批前都复跑一次并登记 VERDICT）。

**结论（§0）**：闸门**确实命中了**——头两轮一枚 `Runner.Worker.exe`、加上一直挂着的 `in_progress` ci run。
本程因此**先等再采**，等待用的是有上界的轮询，逐轮留档。
〔独立复现〕本程自己跑的三行与 `poll.sh`；日志 `/d/tmp/wisp135mg-r2-gate/gate-rounds.log` 可逐轮重看。
