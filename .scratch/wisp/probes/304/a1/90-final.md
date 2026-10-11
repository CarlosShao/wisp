# 304-a1 — 终件（一页结论＋三张表指回件名＋欠读数名册）

腿＝`304-a1`（票 304 `AC#0` 只读普查的一部分）；锚＝**`d75db9ee`**（现跑相符）；`git status --porcelain -- internal cmd tools`＝**0 行**（现跑）；branch `dev`；时刻 2026-10-11 15:1x–15:3x。派单＝`.scratch/wisp/probes/orch/20261011/txt/304a1-dispatch.txt`；票面整件读过（48 行）。

## 1. 交了什么、没交什么

- **交了**：(A) 当前残留名册重量 → `10-residual.md`；表① 全名册 → `20-roster.md`；表③ 三形代价表＋必答 (I)(II) → `30-costs.md`。
- **⛔ 交（本派单整格⛔ 做）**：表② 三形**分开实测**——它要真起子进程＋`go test -timeout` 那一形＝Go 运行面，本腿全程⛔ 碰编译面/运行面（派单 §0.3）。⇒ **⛔ 一件"读文档／读码推出来的结论"被冒充成②的读数**（票面 `:24` 自己点名的那枚雷）；②的欠账具名＝`Q-304-a1-A`。

## 2. 三张表的要点（细节逐枚在件里）

**残留重量（现量 ⇄ 快照）**：`mockllm.exe` **0 枚** ⇄ 派单快照 2 枚（pid 20904/25660 现量**都不在**）⇄ 票面 `:17` 立案后自量 0 枚；`wisp.exe`/`balldebug.exe` **0/0**；`msedgewebview2.exe` **18 枚＝3 棵×6**（`SearchHost.exe`／`wetype_update.exe`／`clipsync-desktop.exe` 三枚父**全在**，逐枚命令行 `--webview-exe-name=` 背书，wisp 系 0）⇄ 票面 `:17` 的 24（四棵，含 magpie）＝**magpie 本人还在（pid 24208）但名下 0 枚**；`%TEMP%` 的 `Test*` 目录 **21 枚** ⇄ 票面 `:16` 的 21（同值）。⚠ 本腿**一枚进程⛔ 停**（处置三形的裁还没裁）。⚠ 新学到一条尺的可靠性：`11724` 现为 `RuntimeBroker.exe`（10-11 07:22 起）＝**PID 会复用** ⇒ 归属只认命令行具名路径。

**表①**：R2 全 scope `.go`＝**72 枚** spawn 语句＝R1 的 `_test.go` **66 枚**＋**非测试文件 6 枚**；66 枚用集合差切成 **18 枚"父不等它"＋48 枚同步 inline**（`logs/partition-66-18-48.txt`，两向作差皆空）。★**两枚 `mockllm` 源**＝`internal/llm/adaptertest/mockllm.go:72`（**⛔ 是 `_test.go` ⇒ R1 那把尺整档看不见它**）＋`internal/llm/openaichat/mockllm_integ_test.go:73`（词面尺 `print-addr` 复量＝恰两枚，与票面 `:15` 同值）。回收形名册：`t.Cleanup`（含 Kill）＝#1/#4/#9/#10/#11/#12/#13、`t.Cleanup` **只 Wait、⛔ Kill**＝#17、**`defer`**＝#5/#6/#8（★外力 `TerminateProcess` 下**一行都⛔ 跑**＝本票缺陷面，单独标出）、**体内内联**＝#2/#14/#15/#16、**Job Object**＝#3（唯一一枚 OS 兜底）、OS 关句柄即端＝见 `30-costs.md`。4 枚 spawn 点带 `CREATE_NEW_PROCESS_GROUP`（R6 现量 3 枚 `SysProcAttr`：`cmd/balldebug/diff_windows.go:412`、`resident_sink_nail_127:222`、`resident_task_source_246:460`）＝控制台信号端⛔ 着它们。

**表③**：**甲形在本仓⛔ 是空想**——仪器与验收都在（`internal/proc/jobscope_windows.go:71-181`＋`internal/proc/jobscope_windows_test.go:41`/`:92-104` 那句判据逐字 "helper survived for %v after Job close"），代价＝两枚 spawn 同批＋一枚新 import 边（现量 `internal/llm` 下只有 `tempdir_resolved_124_test.go:15` 进 `internal/proc`）＋ubuntu 那一支⛔ 受它保护。**乙形**＝`mockllm` 今天**零自退路**（`main.go:49-66` 只 listen→Serve；`control.go` 8 枚端点⛔ shutdown）、flag 现量 4 枚 ⇒ 新增第 5 枚、仓内**零**父 PID helper（尺 0 命中）、且 `tools/mockllm` 是**⛔ require/⛔ replace/⛔ go.work 的独立 stdlib-only 模块** ⇒ 用⛔ 那枚单调钟 helper `observe.NewTimeout`。**丙形**＝票面 `:45` 明写合法，代价具名＝台面继续累积假大脑＋**CI 那侧零进程计数**（`grep -cE "tasklist|IMAGENAME|Get-CimInstance" .github/workflows/ci.yml` 现量 **0**）。⚠ **一处票面已过期**：起手闸**实践**已在数 `mockllm.exe`（`HANDOVER.md:2431/:2455/:2466/:2474/:2490` 五行逐字），仍**只列两名的模板件**＝`probes/orch/2026-10-09-live-window-runsheet.md:45`、`296-ac3-runsheet.md:7`、`a793-append.md:7`、`handover-40as.md:4`、`handover-40at.md:4` ⇒ 那是 `AC#4` 的落点（本腿⛔ 改，只报名册）。

**必答 (I)**：甲形"⛔ 会端别人"停在**读码级**（job 只有 `Assign`/`StartInJob` 两枚入口、只端显式成员；`CreateJobObject(nil,nil)` ⇒ 句柄不可继承；嵌套 job 与 Win8 前失败形都收敛到"要么⛔ 动别人、要么那一枚子代自己没跑成"）。**最后一层需真机/需 CI 色**＝`Q-304-a1-B`；本腿给出可复量的方向：`portable-tests.sh:249-251` 把 `./internal/proc/` 放进 **hosted** `windows-latest`（`ci.yml:534`/`:797`）的射程 ⇒ 若甲形会端 runner，job 用例早该现形。
**必答 (II)**：仪器字面只认 **`\.Sub\(time\.Now\(\)\)`**（`tools/d22scan/main.go:166`，应用点 `:781-792`：空行与 `//` 开头豁免 ⇒ 注释豁免、代码⛔ 豁免）＋第二支 `time.Now().Unix*(` 与 `timeout|deadline|expire|ttl|budget|until` **同行**。**存在性轮询／`time.After`／`time.Now().Add`＋`After`／`observe.NewTimeout`＋`Expired()` 全⛔ 在射程内**（仓内四件现量先例逐枚具名，见 `30-costs.md` §3 表），只有写成"墙钟差比大小"才撞。⇒ 规格侧"存在性轮询算⛔ 算超时实现"**没有判例**＝`Q-304-a1-D`，⛔ 归本腿裁；本腿也⛔ 跑过门禁＝`Q-304-a1-C`。

## 3. 本腿欠的读数（具名逐枚，交编排者）

| 号 | 欠什么 | 为什么本腿交⛔ 出 |
|---|---|---|
| **`Q-304-a1-A`** | **表②整格**：(a) 用例正常结束／(b) `go test -timeout` panic／(c) 外力 `TerminateProcess` 三形分开实测，每一形量"那一发的子进程还在不在"（尺＝具名 CIM 名册）；顺带表①里"外力杀父时子代命运"那一列**全部**〔读码级〕判断、与同步档 48 枚"跟着一起没"那一半 | 派单 §0.3 硬边界（编译/运行面被在飞腿占着）；票面 `:24` 明写前两形即便绿也要交读数 |
| **`Q-304-a1-B`** | (I) 的最后一层：hosted `windows-latest` 上测试二进制的**真实 job 成员关系**＋job 用例在那台机器上的 CI 色 | ⛔ 真机、⛔ CI 现量（⛔ push、⛔ 读 run 日志不在本腿射程） |
| **`Q-304-a1-C`** | 乙形**候选写法**在 `sh scripts/d22scan.sh` 下的 rc＝0 证明（含 `--with-self-test` 那把自测尺） | 派单 §0.3 明禁本腿跑仓里门禁脚本；本腿结论止于"正则字面 ⇒ 词面判断" |
| **`Q-304-a1-D`** | 规格侧未定案：**存在性轮询**算⛔ 算 D22 那句"用墙钟时间差实现超时"；权威文本（`PLAN.md` D22 正文／SPEC-06）里本腿找⛔ 到判例 | 属规格解释＝人工批准面（`AGENTS.md §0` 第 2 句）；⛔ 由普查腿填 |
| **`Q-304-a1-E`** | `cmd/wisp/resident_audio_247_live_windows_test.go:239`（tag＝`windows`，**⛔ winlive**）起的 `powershell` 声源在托管色上到底跑⛔ 跑 | 需要 skip 逻辑的真机/CI 色；本腿只读了 tag 尺与 tier 射程尺 |
| **`Q-304-a1-F`** | 票面 `:10` 那两枚孤儿的占用/CPU/两枚 Listen 套接字读数 | 进程已经没了 ⇒ 历史值⛔ 再复量（本腿只量能证的"PID 已不存在＋PID 已复用"） |
| **`Q-304-a1-G`** | 24→18 那 6 枚（magpie 那一棵）**为什么**没 | 本腿只量到"magpie 还在（pid 24208）＋名下 0 枚"；成因需另尺 |
| **`Q-304-a1-H`** | `%TEMP%` 21 枚 `Test*` 目录的**逐枚归属**（哪一次杀留下哪一枚）＋目录 mtime 名册 | 起它们的进程⛔ 在名册里；本腿只量枚数（票面 `:16` 已警告那半把尺落空） |
| **`Q-304-a1-I`** | ⛔ 走 `os/exec` 的起法（`syscall` 直调 `CreateProcess` 那一族）的名册 | 本腿只打了 `os.StartProcess`（0 命中），⛔ 打直调尺 ⇒ 表①在这一面上⛔ 闭 |
| **`Q-304-a1-J`** | `ci.yml:430` 那枚 compose 起的 mock-llm 算⛔ 算本票"测试里起子进程"的射程 | 本腿按票面字面（"测试里起子进程"）⛔ 入册，需编排者一句裁定（同一族残留的**另一种来源**） |

## 4. 越界自报（⛔ 自我盖章）

- 写面＝**只新建** `.scratch/wisp/probes/304/a1/**`（`.md` **5** 枚＋`logs/*.txt` **16** 枚＋交件消息件 1 枚）＋票尾**纯 ＋1 行 −0 行**一枚 Progress log bullet；⛔ 动任何原句／⛔ 翻任何框／⛔ 碰 `Status:`／⛔ 改名 `-done`。
- ⛔ 跑过任何 Go 编译面/运行面命令（`go build`／`go vet`／`go test`／任何 `.exe`／`scripts/d22scan.sh`／`scripts/portable-tests.sh`）；**`go env`／`go list` 也一枚没用**（＝上面 `Q-304-a1-C/E` 那两枚欠账的理由）。
- ⛔ 杀任何进程；⛔ 删任何件；⛔ 在仓内建 worktree/checkout；⛔ 在 `probes/**` 下建 `.go`；⛔ 碰 `tools/d22scan/allowlist.txt`／`docs/reports/**`／`docs/evidence/s1/**`／`internal/**`／`cmd/**`／`tools/**`／`scripts/**`／`.github/**`／`AGENTS.md`／`docs/PLAN.md`／`docs/specs/**`（票面行只**追加**一行 Progress log）。
- ⛔ push（A813 那枚授权与"推送归编排者"两条同向）。
- PowerShell 全程⛔ `$_`（内联只用 `Where-Object <名> -match/-like` 简化形）；⛔ 把带反引号的正文塞进双引号 shell 串（写件一律走 Write 工具）。
- ⚠ 换行形状自报：`git add` 本腿那 17 枚 `logs/*.txt`＋1 枚消息件时报了 14 行 `LF will be replaced by CRLF the next time Git touches it`＝仓内既有 `core.autocrlf` 那一形（别腿具名过＝"格式名册里的**换行假枚**"，例 `308-r2` 那笔的 `text: set / eol: lf`＋`core.autocrlf=true`）。本腿⛔ 动 `.gitattributes`、⛔ 改任何件的字节、⛔ 把这一条算成自己的越界；只把它登记成"下一位量 `gofmt`/长度名册时可能看见的别家形状"。
