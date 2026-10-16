
### 4.0dv 停车点（**2026-10-16 13:0x +08**）｜编排者｜第十波收完两枚（`326-v1` 结案 ∥ `165-r3` 转裁）⇒ 台账 `A1007` ＋ ★放行 `270-r2` ＋ 第十一波五枚在飞（写 1＋只读 3＋裁 1）＋ ★推送前提已复尺（那枚红绿了）

**① 盘上状态（接手先复跑这三把尺，⛔ 信本节）**：`git log -1`＝`1bb2fcee`（A1007 收档笔，票 326 翻六枚＋改名 `-done`）；台账末节＝`## A1007`（尺＝`grep -c '^## A1007' docs/reports/pending-and-issues.md`＝1，文件 18856 行）；本地领先 `origin/dev`＝**127 枚**（尺＝`git rev-list --count origin/dev..HEAD`，⛔ 裸推会去 cnb、⛔ force，机主授权 `A813` 只覆盖 `origin/dev`）。

**② ★推送前提复尺（这条是两周来的那把锁，今天开了）**：按住推送的理由＝dev 上压着 `TestTicket255RosterStillMatchesTheActualReadSites` 一枚已知可修的红（归口票 324）。现跑＝**绿**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp -run 'TestTicket255RosterStillMatchesTheActualReadSites' -count=1 -timeout 30m -v` ⇒ `--- PASS (0.73s)`／`ok 0.820s`／`rc=0`，14/14 段 `match=true`（件 `logs/orch-255roster-check2-1bb2fcee.txt`）。⚠ **第一次跑我拿回 `exit status 0xc0000135`＋`rc=1`、`--- FAIL` 零枚＝那包根本没跑**（少 DLL 的台面形，同族坑见 `wisp-measurement-traps`），补 PATH 才成读数 ⇒ **引 `cmd/wisp` 任何颜色前先核 PATH**。⇒ **仍⛔ 现在推**：本波写腿在飞，推一次会自启 `slo-full` 抢本机 CPU（机主 70% 那把尺）、且我更想一次推完再取 CI 色。

**③ 在飞（第十一波，13:0x 放出，五枚）**：`270-r2`（写·`internal/tools/subagent_197.go` **只注释**，带"产码字节零改动"尺＋正控）∥ `165-v2`（裁决·票 165 `AC#6`，独占 `internal/agent` 编译面）∥ `327-a1`／`328-a1`／`216-a2`（三枚只读普查，⛔ 跑 Go，写面各只自家 probes＋本票末尾一条 bullet）。放腿判据＝CPU 连发三发 **44/42/43 ≤70%** ＋ 进程名册除机主 `wisp.exe` **PID 21892** 外 `go.exe`/`*.test.exe` 命中 0。派单存档＝`2e706255`。

**④ ★共享编译面本波空出**：`cmd/wisp/**` 与 `internal/panel/**` **零写腿**（`220-r2` 已死于 150 轮帽、遗产由我代提入 `A1005`；`216-r1`／`220-v1`／`326-r1` 全收完）。⇒ 下一波**必须**补写腿（机主两次点名"⛔ 一腿独跑"），候选按料齐程度排：`313` 的 `cmd/wisp` 落点包（⚠ 三格未勾是**界面侧**，那份归机主的外部前端 agent、本编队零写面，别派腿去动 `frontend/**`）、票 252 残余、`220` 工作区那半。

**⑤ 队列（每枚带释放条件，⛔ 无条件待派）**：`216-v2`（裁票 216 `AC#5`＋`AC#3` 未被打的那一半）＝**按住等 `216-a2` 交件**（没现量不许裁）；`328-a2`（票 328 `AC#1` 选形料）＝等 `328-a1`；`327-a2`＝等 `327-a1`；`329-c1`（`326-v1` 交回的缺母行／迟到戳两枚，先查重再决定立票）＝队列尾，⛔ 现在不写票面判据。

**⑥ 机械手册（本波新增／复用，接手直接用）**：ⓐ ★**台面配方第二改**＝八件套只用于 `go build`/`go test`；**`d22scan`／长度门／`gofumpt` 只在工作树或全树 `git archive`（不带 pathspec）上跑**——`tools/d22scan/main.go:581-582` 自己声明了 `design/`＋`frontend/` 两枚射程，声明了的射程扫到 0 枚是**硬退出 `rc=2`**，长得像红（本波 `165-r3` 缺 `design/`、`326-v1` 缺 `.git` 各撞一次，两枚腿都具名上报没算进判据）。ⓑ 翻框照旧走机械件（`txt/flip-326.sh`：尺 6/0→0/6，diff 只见 `[ ]`→`[x]`，判据文字一字未动）。ⓒ `git mv` 之后⛔ 再把旧路径放进 `git add`（`fatal: pathspec … did not match any files`＝整条 `&&` 链断掉、一枚都没暂存）；正确形状＝`git mv` 已把重命名放进索引，`git add` 只列**存在的路径**，`git commit -- <旧 新 全部>` 两处都带 pathspec。ⓓ 结案改名 `-done` 前后各跑一次长度门（本波改名后 rc=0，最后一名组件 95 字／相对 116 字 ≤ 帽 121 ⇒ 不需进名册）。

**⑦ 待人项（照旧只一枚，默认不催）**：**Q-85**＝需要机主给约 30 分钟、让他自己退出他 10-15 19:03 起的 `build\wisp.exe`（PID 21892，台账 `A988`），我才能出票 220 `AC#5` 那格要的**整包终态名册**；⛔ 我或任何腿代他关这个进程。另：工作树 `go.mod` 有一枚**不是任何腿的活**的脏改动（`github.com/jchv/go-webview2` 由 indirect 提为 direct，mtime `2026-10-10 22:28`），`go.mod` 在我这儿是禁区，⛔ 由腿顺手 `go mod tidy`，只具名留给机主。

**你需要做什么：没有。** 五枚腿在跑；Q-85 那 30 分钟你哪天方便再说，我不催。
