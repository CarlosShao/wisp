# 326 — 会话里**继续聊不会把它顶到前面**：`AppendSessionTurn` 只往 `session_turn` 插一行，⛔ 碰 `session.updated_at`，而列表排的就是那一列

- Status: **已立，待派**（2026-10-16 10:5x 编排者立；来源＝非实现者验收腿 `166-v1` 交回来的第三形＋票 166 的 `AC#2` 接线缺口；台账 `A1003`）
- 类型：**行为缺陷**（不是判据仪器缺陷、不是契约变更）。射程＝`internal/memory/**` ＋ `internal/session/**`，⛔ 动 `cmd/wisp`／`internal/panel`（共享编译面本波归 `220-r2`）。
- ⚠ 行号只当**读数**、⛔ 当判据（这仓的行号锚会随别人加行腐烂）；下面每条都附可复跑的尺。

## 这是什么（人话）

你在跟 Wisp 说话，每说一句，它都会把这句话存进那次对话的记录里。可是"这次对话排第几"看的是**它上一次被动过的时间**，而"存一句话"这一步**没去更新那个时间**。
结果就是：你在旧对话里聊了半天，列表上它⛔ 不会跳到最前面；你要找它，还得往下翻。
今天唯一会把它顶上去的动作是**改名、归档、分叉**——都是"你手动折腾它"，⛔ "你在用它"。

## 现量（起手逐条复算，⛔ 信这里的号）

| 事实 | 读数（编排者 2026-10-16 10:5x 现跑，尺附在旁边） |
|---|---|
| 追加那一轮只写子表 | `internal/memory/dao_session.go:253` 是 `func (s *Store) AppendSessionTurn(...)`；该函数体内**只有** `INSERT INTO session_turn(session_id, seq, role, content, created_at)` ⇒ 同函数内⛔ 一条 `UPDATE session`（尺＝`grep -n 'func (s \*Store) AppendSessionTurn' internal/memory/dao_session.go` ＋ 把那一段整块读出来数 `UPDATE`） |
| 真正写 `session.updated_at` 的只有四处 | 同一文件：建会话时（`:84` 的 `INSERT`）／改名（`:177`）／归档（`:198`）／分叉另建（`:231`）⇒ **说话不在里面**（尺＝`grep -n 'updated_at' internal/memory/dao_session.go`） |
| 列表确实按这一列排 | 同文件两处 `ORDER BY updated_at DESC, id DESC`（`:125` 会话列、`:144` 搜索列） |
| 时间源 | 那几处用的都是 `observe.NowWallUTC().Unix()`＝**秒**精度；`internal/session/manage.go:77` 默认 `time.Now` ⇒ "同一秒动过两条"是日常形状（同秒先后＝票 166 `AC#6` 那枚硬币，本票⛔ 治它，只⛔ 让每一轮都成为一次"动过"） |

## 判据（每一格都要答"这一发在未修码上响不响"）

- [ ] **AC#0 起手复算**：把上面四行逐条现跑（⛔ 用我的号），并把尺与射程落 `rc=`。
- [ ] **AC#1 每一轮都算"动过"**：`AppendSessionTurn` 在同一枚事务里刷新该会话的 `session.updated_at`（时间源沿既有那一族，⛔ 新造第三种时钟、⛔ 用墙钟时间差实现任何超时）。完成判据＝成对读数：种"会话 A 先建、B 后建，然后在 A 里追加一轮" ⇒ 改前列表首个仍是 B、改后是 A；**未修码上这一发必红**，红句连文案逐字落盘。
- [ ] **AC#2 ⛔ 顺手治 tie-break**：同秒两条的先后今天由 `id DESC` 决定，而 `session.id` 是 `rand.Read` 铸的随机 hex（⛔ 单调）⇒ 那是票 166 `AC#6` 那一族的另一半，**换 tie-break 词＝行为变更**，本票射程外。要动就具名上报、由我另裁。判据＝本腿交付里 `ORDER BY` 那两行**逐字节不动**（尺＝`git diff` 对 `:125`／`:144` 两行命中 0）。
- [ ] **AC#3 原子性与迁移不破**：多这一条 `UPDATE` ⛔ 改变 v2→v3 迁移那族的崩溃形状（中途被杀之后要么全有要么全无，沿票 166 `AC#3` 已立的真杀进程台件形）；要一发"摘掉这条 `UPDATE` 必须回到那对颜色"的去牙反控，⛔ `-overlay`。
- [ ] **AC#4 门禁＋整包颜色**：`go build ./...`／`go vet ./internal/memory/ ./internal/session/`／`sh scripts/d22scan.sh`／长度门 `--with-self-test`／格式化尺（本机 `gofumpt` ⛔ 在 PATH＝rc=127，但盘上有 `/d/work/base/gopath/bin/gofumpt.exe` ⇒ 按绝对路径跑，⛔ 拿 `gofmt -l` 的 0 冒充）＋`go test ./internal/memory/ ./internal/session/ -count=1` 逐包四数（`RUN=`／`PASS_TOP=`／`FAIL_TOP=`／`SKIP=`，⛔ 混用两把尺）。
- [ ] **AC#5 生产可见性说实话**：会话列表今天**压根没接进产品**（票 166 `AC#2` 判"接线缺口"：`git grep -c 'session\.NewManager' -- '*.go'` 命中 0 枚、`cmd/wisp/main.go` 里⛔ `sessions`/`history` 子命令）⇒ 本票修完在产品面上**仍然看不见**。这一格要具名写"修完之后界面上还差哪一跳"，⛔ 顺手接线（那一跳撞共享编译面，归我另排）。

## 禁区

⛔ `docs/PLAN.md`／`docs/specs/**`／SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／三枚冻结件／`frontend/**`／`design/**`／`.github/**`／`ci.yml`／`go.mod`／`go.sum`／`.gitignore`；D1–D47／C1–C32／R1–R9／D43 转移表 ⛔ 动（本票⛔ 新造状态名、⛔ 改状态拼写）。
⛔ 为变绿放宽任何断言、⛔ 删用例、⛔ 加 `t.Skip`。⛔ `cmd/wisp/**` 与 `internal/panel/**` 一个字节（本波写腿＝`220-r2`）。
⛔ `go test ./cmd/wisp` 整包（票 314 死亡区）；⛔ `go mod tidy`／`go get`。

## 落点候选（⛔ 预设，落地腿自己现读后选并具名理由）

ⓐ 在 `AppendSessionTurn` 的同一枚事务里补一条 `UPDATE session SET updated_at=? WHERE id=?`；ⓑ 让它复用既有那枚"动过"的写侧助手（如果有）。
判据：① ⛔ 动既有签名（要动就停手上报）；② 时间源沿既有那一族；③ 写失败要 propagate，⛔ 静默降级成"只插子表"。
