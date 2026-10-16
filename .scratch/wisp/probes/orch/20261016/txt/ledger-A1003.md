## A1003（2026-10-16 10:5x +08，编排者）收 322-v1／216-a1／166-v1 三枚交件（166-v1 死于 150 轮帽但裁决表已在盘上）⇒ 翻框十一格＋新立票 326＋我自抓一枚台面错

**一、三枚腿的收件＝盘上三把尺＋决定性读数我自己复跑（⛔ 转述腿的数）**
- `322-v1`（非实现者终裁，票 322 六格）：五笔提交 `755fb974`／`8b8fa0be`／`f918d81e`／`6eeb71ea`／`da6fda0d`；判语件 `docs/evidence/s1/322-v1-acceptance.md` 现量 **163 行／16776 字节**；读数件根 `.scratch/wisp/probes/322/v1/`（含 `logs/` 38 枚＋台件 3 枚）。越界尺**逐枚**跑（⛔ 用区间——`755fb974~1..da6fda0d` 那个区间里有别腿提交，我第一次就是这么读出一枚假"139 枚越界"，改成 `git show --name-only` 逐枚作差后＝**55 枚路径、越界 0**）。
- `216-a1`（票 216 AC#1 普查，只读）：三笔 `63fff90b`／`a9886142`／`35cfbad9`；名册件 `docs/evidence/s1/216-census.md` **134 行／14673 字节**；读数根 `.scratch/wisp/probes/216/a1/`（`logs/` 45 枚）。★它的来源口径值得抄进后续派单：**工作树进场即脏＋写腿在动 `cmd/wisp` ⇒ 名册一律读提交树快照 `71002d0`**（`git grep ... <rev>`），⛔ 混在飞态。
- `166-v1`（票 166 六格终裁）：**死于 150 轮帽**，但裁决表与它自己那笔票尾 bullet 已入库＝`2c679634`；`docs/evidence/s1/166-session-verdict.md` **98 行／13583 字节**；没交完的只有第五批读数的 commit 正文一枚未跟踪件 `probes/166/v1/msg-b5.txt` ⇒ 我按原样代提入库、⛔ 替它补内容。

**二、翻框（裁决者≠实现者；每格都附我自己的复跑）**
- 票 322 ⇒ 翻 **AC#0／AC#1／AC#2／AC#3／AC#5**，**AC#4 ⛔ 翻**（腿判"部分成立"，缺的正是 `cmd/wisp` 那一发整包读数＝票 314 死亡区、本波归唯一写腿 `220-r2`）；现量尺＝5 未勾变 1 未勾／5 已勾。我的复跑：`winsec.PrivateFile(` 生产 **5 枚**调用者＋`PrivateFileExclusive(` **2 枚**（另有 1 枚在 `.scratch/wisp/probes/174/r4/pristine-spill.go`＝归档拷贝、⛔ 算生产）；`booked separately` 在 HEAD **0 命中** ⇄ 正控 `8baec646^` 在 `winsec_other.go:82` **1 命中**；AC#3 本机两枚用例 `-count=1` 全绿 rc=0 且 `witnessTruncated=true witnessSize=35` 同色复现。★票面第二形〔仅本机可量〕我照腿的判语认，⛔ 任何"Windows 窗口已关"的读法。
- 票 166 ⇒ 翻 **AC#1／AC#3／AC#4／AC#5／AC#6**，**AC#2 ⛔ 翻＝接线缺口**（我的复跑：`git grep -c 'session\.NewManager' -- '*.go'` ⇒ **命中 0 枚、rc=1、stderr 0 字节**，正控 `func NewManager` 在 `internal/session/manage.go` **1 枚**；`grep -cE '"sessions"|"history"' cmd/wisp/main.go` ＝ **0 枚**；`go test ./internal/session/ -run TestSessionListByWorkspace -count=20 -v` ⇒ **`RUN=20 TOPPASS=20 TOPFAIL=0` rc=0**）。AC#6 按腿的**范围修正**记账：真正关掉硬币的只有 `:69` 那一枚，`:84` 实测⛔ 是硬币（改动属防未来）。⇒ 票 166 ⛔ `-done`：产品面上"用户回不到昨天那次对话"照旧为真，那一跳撞共享编译面，释放条件＝`220-r2` 交件。
- 票 216 ⇒ 翻 **AC#1**（这一格是"先出名册"的普查格，凭据＝名册件＋我逐条复跑：`unicode.IsControl` 整树 1 枚、`findMetachars` 第二份、`sanitizeDisplayName` 只 trim＋取末段、`ContainsAny(name,"\x00\r\n\t")` 只挡 4 枚字符 ⇒ ESC／C1 今天能穿过进菜单、`composer.go` 里 `plugin|skill|command` **0 命中**＝三枚字段今天零源零出向）。

**三、★我自抓（三笔，原句⛔ 抹平）**
① **台面错（我自己派单里明令过的那条，我转头踩）**：第一次跑 322 那两枚用例前⛔ 先 `mkdir -p` 那枚 `GOTMPDIR` ⇒ 盘上只留两行 `go: creating work dir ... cannot find the file specified`；故障件 `.scratch/wisp/probes/orch/20261016/logs/orch-322v1-ac3-repro.txt` 只建不删、当场作废、⛔ 计入判据；第二次建好目录后 rc=0 出真读数（`orch-322v1-ac3-repro2.txt`）。同族教训＝记忆第 160 条。
② **一把尺的区间读法自抓**：`git log --name-only <锚>~1..<HEAD>` 当"这枚腿的路径清单"用＝错（区间里混着别腿的提交，我当场读出"139 枚越界"这种荒话）⇒ 定式：**逐枚 `git show --name-only` 再作差**，且作差前先证区间枚数与提交号一一对上。
③ **又一枚派单里的过期数被腿顶回**：我给 `216-a1` 的派单把票 200 那五状态的第五枚写成 `na`，盘上真值＝**`not_run`**（`internal/panel/instructions_200.go:40`，我复跑五枚常量逐枚对上；全仓 `\bna\b` 只命中 `ticket84_no_owner_test.go` 的测试局部变量 3 行）⇒ 票 216 尾部只加更正、⛔ 抹我原句。另记口径差：腿的正控写"`.Error()` 114"，我那两把尺读出 **374 行／144 枚文件**（`-n` 数行 vs `-c` 数文件）⇒ 只证尺⛔ 恒 0，⛔ 与它的数相减。

**四、新立票 326＋放一枚写腿（第二枚写腿？⛔——本波仍只一枚碰共享面）**
`326-append-session-turn-does-not-refresh-session-updated-so-a-live-chat-never-moves-up.md`＝会话里继续聊⛔ 把它顶到前面（`AppendSessionTurn` 只 `INSERT INTO session_turn`、⛔ 碰 `session.updated_at`，而列表排的就是那一列；写 `updated_at` 的既有四处全是"手动折腾"）。写面＝`internal/memory/**`＋`internal/session/**`，**与共享编译面无交** ⇒ 可与 `220-r2` 同波。落点约束我已裁：只做"每一轮都算动过"、⛔ 顺手治同秒 tie-break（AC#2 明令射程外）、⛔ 碰 `cmd/wisp`／`internal/panel`、⛔ 新造时钟。
队列变化：`216-r1` 待派，释放条件＝`220-r2` 交件（共享编译面每波只许一枚写腿，这条规矩我自己⛔ 破）；`220-v1` 按 `220-r2` 交件；`120-v2` 按 `322-v1` 交件（现已交，但两枚同碰 `internal/winsec` ⇒ 等本波只读腿散场再放，取干净读数）；票 166 的接线那一跳排 `220-r2` 之后。
未推＝本地领先 `origin/dev`（尺现量见 `git rev-list --count`），释放条件不变＝本波写腿收档之后一次推送（⛔ force、只 `origin/dev`、授权 `A813`）。
