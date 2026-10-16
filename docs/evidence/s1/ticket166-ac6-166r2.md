# ticket 166 AC#6 · 166-r2 读数件（实现者自陈，⛔ 裁决／⛔ 勾框）

- 腿：166-r2（落地腿·写）｜票：`.scratch/wisp/issues/166-no-session-layer-list-search-rename-archive-fork-on-the-go-side.md` AC#6 ⓐ＋ⓑ
- 锚：dev `a9c7bea0`（派单写的 `17ca21ee` 是它前一笔；起树用当时盘上 HEAD）｜时刻 2026-10-15 23:0x–23:3x +08
- ★这份件的身份：**实现者本人的读数与自陈**。按 D22 双角色与 `SPEC-12 §4.3` #1，
  AC#6 那两半格成不成立⛔ 由我判——裁决与翻框归非实现者；本件只保证"每一格都有落盘的数、
  每一发都写得出台面名"。读数目录＝`.scratch/wisp/probes/166/r2/`（01 台面／02 改前／03 改后／04 牙／
  05 生产可见性／06 门禁／07 日志索引与坏尺自抓＋`logs/` 14 枚抽取件＋`00-ruler.sh`/`06-mutate-clock.sh` 两把尺）。

## 1. AC#6 ⓐ 仪器改成确定性（判据四发，全在⛔ 同一棵树上）

台面（01 §1.1）＝`git archive --format=tar HEAD | tar -x -C /d/tmp/wisp-166r2/tree`（仓外导出树，1087 枚 `.go` 先证非空）、
`GOTMPDIR=D:/tmp/wisp-166r2/gotmp`（预先 `mkdir -p`）、尺＝`go test ./internal/session/ -count=20 -run <regex> -v`。
**改前／改后／正控／还原四发⛔ 跨台面**（同一棵树、同卷、同 GOTMPDIR 写法）。

| 发 | 形态（sha256 前缀） | rc | RUN | 顶层 PASS | 顶层 FAIL | `TestSessionListByWorkspace` | 判定 |
|---|---|---|---|---|---|---|---|
| before-count20 | 57fcf6efc4624486（＝a9c7bea0 的字节） | 1 | 40 | 33 | 7 | PASS 13／**FAIL 7** | 改前复现 ✓ |
| before-count20-round2 | 57fcf6efc4624486 | 1 | 40 | 28 | 12 | PASS 8／**FAIL 12** | 硬币非恒定 ✓ |
| after-list-count20 | 2bf320df04e2819a（＝交付字节） | 0 | 40 | 40 | 0 | **PASS 20／FAIL 0** | 20/20 ✓ |
| after-list-count20-round2 | 2bf320df04e2819a | 0 | 40 | 40 | 0 | PASS 20／FAIL 0 | 再一发 ✓ |
| after-search-count20 | 2bf320df04e2819a | 0 | 20 | 20 | 0 | －（`TestSessionSearchTitleAndBody` PASS 20／FAIL 0） | 同形复跑 ✓ |
| after-all-count20（`-run TestSession`） | 2bf320df04e2819a | 0 | 100 | 100 | 0 | PASS 20／FAIL 0 | 本票五枚出口用例全绿 |
| after-fullpkg-count20（整包 `-v`） | 2bf320df04e2819a | 0 | 720 | 380（＋子测试 340） | 0 | PASS 20／FAIL 0 | 19 枚用例名册见 03 §3.4 |
| teeth-flatclock-list-count20 | 28f683c77893474d（物理摘掉步进） | 1 | 40 | 18 | 22 | PASS 12／**FAIL 8** | 牙成立 ✓ |
| teeth-flatclock-list-count20-round2 | 28f683c77893474d | 1 | 40 | 13 | 27 | PASS 11／**FAIL 9** | 再一发 ✓ |
| restore-list-count20 | 2bf320df04e2819a（`cmp` 与交付字节 0 差异） | 0 | 40 | 40 | 0 | PASS 20／FAIL 0 | 还原绿 ✓ |
| final-green-list-count20 | 2bf320df04e2819a | 0 | 40 | 40 | 0 | PASS 20／FAIL 0 | 盘上最后一发＝绿 |

单位名（⛔ 一个数三种读法）：RUN＝`^=== RUN` 行数（本用例无子测试，`^    --- ` 命中 0 枚 ⇒ 两把尺在此重合）；
PASS/FAIL＝**顶层用例**行数（`^--- PASS`／`^--- FAIL`）；按枚用例的那两列只数它自己的那枚顶层行。

改动本身（唯一产码面＝`internal/session/manage_test.go`，`git diff --numstat`＝14 加／2 删）：
`:69` 与 `:84` 两处 `newTestManager(..., nil)` 改成该件里已有的步进时钟形状 `ptrTime(time.Date(2026,5,1,0,0,0,0,time.UTC))`
（沿 `:48` 那一发，⛔ 新造第二枚时钟助手 ⇒ 签名⛔ 要扩，"没做完"那节因此是空的），另 12 行注释。
⛔ 放宽：`want` 一行未动、`if len(got) != 2 || got[0].ID != c.ID || got[1].ID != a.ID`（现 :86）逐字仍在盘上、
`t.Skip` 命中 0 枚、用例名册 5 枚未减（03 §3.6 那把尺即证据）。

红句逐字（连文案，⛔ 改写；全文更多见 02/04）：
```
    manage_test.go:78: workspace filter = [sess_5efe0ab090be1ce123fb8f59c4aa8ec4 sess_2d5404295ef429ff8313a516ea339a2e], want [sess_2d5404295ef429ff8313a516ea339a2e sess_5efe0ab090be1ce123fb8f59c4aa8ec4]
    manage_test.go:87: workspace filter = [sess_b8d775c6561e8bdf60f0d2f24fff92f3 sess_afa3e90f0f538446da3409cc42fb143b], want [sess_afa3e90f0f538446da3409cc42fb143b sess_b8d775c6561e8bdf60f0d2f24fff92f3]
```
第一行＝改前形态（断言在 :78）；第二行＝突变态（交付字节的同一条断言在 :87，那 9 行注释把它挪了位）。

## 2. AC#6 ⓑ 只答不改（完整链条在 05，这里给结论）

**答：看不见。** 链条不是"排错了但没人看"，是**压根没有一条调用链走到那枚清单**：
`git grep -n "session\.NewManager" -- '*.go'` 退码 **1＝0 枚命中** ⇒ 产码里⛔ 构造过 `Manager`；
`Manager.List`／`Search` 的调用点 9 枚**全在** `internal/session/manage_test.go`（:56 :82 :108 :115 :122 :125 :157 :160 :169）；
生产碰 `internal/session` 的只有 D45 授权台账那一支（`cmd/wisp/run.go:282/476/481`、
`cmd/wisp/resident_approval_windows.go:243/285` 全是 `session.Ledger`/`Mint`，⛔ `Manager`）；
`memory.Session` 类型在 `internal/panel`／`internal/tools`／`cmd/wisp` 的产码里 **0 行**；
CLI 子命令名册（`cmd/wisp/main.go:89-124` 十枚）与面板快照（`internal/panel/composer.go:57-61`＋`pump.go:138-155` 五枚源）
里都⛔ "把历史会话列出来"那一格。**⇒ 这一形今天是"看不见的债"，界面腿/票 145 一接线就立刻可见。**

数据侧它⛔ 是假想形（05 §5.2）：写 `session.updated_at` 的四枚点全是秒（`manage.go:93` 默认 `time.Now`、
`dao_session.go:175/:192/:218` 的 `observe.NowWallUTC().Unix()`）⇒ "同一秒动过两条"是日常形状。

三形料上交编排者裁 tie-break（⛔ 我动 `dao_session.go:125/:144`，⛔ 我替他裁）：
1. `SearchSessions` 排在 `dao_session.go:144` 同一枚 `ORDER BY updated_at DESC, id DESC` ⇒ 两枚命中会话同秒＝同一枚硬币。
2. ★`AppendSessionTurn`（`dao_session.go:253-278`）**不**刷新 `session.updated_at` ⇒ "按时间列"排的是
   最后一次 创建/改名/归档/分叉，⛔ 最后一次说话；会话里继续聊⛔ 把它顶到前面。这一形比 tie-break 更像产品缺陷，
   但它是行为变更，同样⛔ 本腿动。
3. 对照组（同一把尺量出来的"没病的那一支"）：`profile` 那两支 `ORDER BY updated_at …, id …`
   （`dao_profile.go:86/:111`）⛔ 是硬币——它的 `id` 是 `INTEGER PRIMARY KEY AUTOINCREMENT`（`schema.go:31`）＝单调；
   `session.id` 是 `TEXT PRIMARY KEY` 由 `session.go:96 rand.Read` 铸（`schema.go:221`）＝随机 hex。
   ⇒ 同形 SQL 在两枚表上一好一坏，坏的那一半⛔ 换个 tie-break 词就能治，⛔ 得靠仪器（本程已治）。

## 3. 门禁（六道，逐枚 rc 在 06）

rc=0 `go build ./...`（工作树⇄导出树两发）｜rc=0 `go vet ./internal/session/`（两发）｜
rc=0 `go test ./internal/memory/ -count=1`（导出树，邻包控制：本票那枚 ORDER BY 真被它跑着，⛔ 我动一字）｜
rc=0 `sh scripts/d22scan.sh`（clean，分母 bans#1-5 internal/=239、cmd/=41｜#6 frontend/=116｜#7 tools/=25｜#8 design/=39、frontend/=116、internal/=571、cmd/=138）｜
rc=0 `sh scripts/check-path-length-budget.sh --with-self-test`（在仓根；VERDICT GREEN，13453 枚 tracked 路径、越帽 60 全在册；
那句 `RED - over budget and NOT in the roster: …seeded-…` ＝种下去的 positive control，⛔ 红）｜
rc=0 `gofmt -l`＝列出 **0 枚**文件｜
`gofumpt`：`which gofumpt` rc=1（⛔ 在 PATH，派单那句在这一维成立）**但**盘上真在＝
`/d/work/base/gopath/bin/gofumpt.exe`（5,028,352 字节）⇒ 按绝对路径跑了那把更严的：rc=0、列出 **0 枚**文件
（⛔ 拿 gofmt 的 0 冒充 gofumpt 的 0，也⛔ 用"本机不可用"了事）。

## 4. 没做完／缺的读数（具名）

叙述 **AC#6 ⓐ＋ⓑ 两格本程⛔ 有"没做完"的格**：改前复现、改后 20/20、Search 同形、正控红＋还原绿，五发齐；
  台面名逐发写在旁边（01 §1.1 那一棵，全批发同一棵）。⛔ 新造第二种时钟助手 ⇒ 签名⛔ 扩 ⇒ 那条上报项空。
★一处**读数与派单前提不符**（具名，⛔ 我自己判成"顺手也修好了"）：派单/票面写"`:84`（`TestSessionSearchTitleAndBody`）
  ＝硬币"。实测⛔ 支持——把时钟再做平，那一枚 **20/20 照旧绿**（04 §4.4，raw_sha a8ff440eaa10a24c），
  读码原因＝它对结果的断言只有"命中 1 枚／0 枚"（现 :112/:119），⛔ 比两枚之间的先后。
  台账里它总共命中 1 处，且那处是"死计数"教训（`pending-and-issues.md:18583`），零红记录。
  ⇒ :84 那一处改步进时钟＝**防未来**（将来给 Search 加次序断言时⛔ 静默变硬币），⛔ 算"关掉第二枚硬币"；
  本程真正关掉硬币的只有 `:69`（`TestSessionListByWorkspace`）那一枚。
叙述 一处**台面实况与派单预告不符**（⛔ 记成我的凭据）：派单说 322-r1 在飞时全仓 `go build` 可能红——
  本程那两分钟 `go build ./...` 工作树与导出树都 rc=0（06 §6.4）＝我⛔ 赶上那一分钟，⛔ 据此说"winsec 没问题"。
叙述 本腿⛔ push（推送归编排者，只 `origin/dev`）；⛔ 勾票上任何 AC 框；票尾只追加一条 bullet。
