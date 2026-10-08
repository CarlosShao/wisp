# 票 255 · 255-r6 交付件（落地腿，2026-10-08）

活＝把 `TestTicket255PanelHostBuildsItsWindowOptions`（ruler 3）里"全工程只许一处"那颗源读型
"恰好 1"断言的**射程扩到它名字宣称的范围**（派单 §2 的 (i) 支），并按 `A713` 定式用**盘上真种的突变**
证明它真会响（⛔ 全程未对源读尺用 overlay 作证）。起手锚与现量见 `00-anchor.md`。

---

## ① 射程选了哪个、为什么

**选定：`cmd/wisp` 包的产码源件（包目录全部 `.go` 剔 `_test.go`）**，census 形状照
`cmd/wisp/panel_locked_naming_33r11_windows_test.go:packageSourceFiles33r11`（`:176 os.ReadDir` ＋
`:207 os.ReadFile` ＋ 空名册＝Fatal ＋ 解析失败＝Fatal 不豁免）。

为什么不是"整仓"：本仓同 module 内（无独立 go.mod，已现量）存在**真的第二调用点**
`scripts/spike/webview2-latency/main.go:149` ⇒ 在整仓射程下"恰好 1"**基线即红**，而派单⛔不许动
期望值/阈值/断言方向 ⇒ 整仓这条只能靠放宽或改期望值"做绿"，那是假射程。为什么不是别的：
`255-v2` Q2 具名的最坏形状（"以后任何人在 `cmd/wisp` 别处加第二条尺寸路而 255 的尺一声不响"）
正落在**本包产码**这一层；`.scratch/wisp/probes/**` 下另有 7 枚 `NewWithOptions` 命中全是探针存档
（不在 `cmd/wisp` 目录内，整包射程天然不含）。测试件豁免与 33r11 同理：名册钉的是**发货码**，
否则守卫自己的夹具会把尺打死红（本腿的 M3 夹具即例）。

## ② 走了 (i)（扩射程），逐枚前后对照

写面内改动两枚，全在 `cmd/wisp/` 测试件：

1. **`cmd/wisp/panel_geometry_255_test.go`**（改 3 段，**净零增删行**：`wc -l` 前后均 **456**，
   numstat **11/11** 成对——A714 §3 的正确尺；行号一枚没挪，ruler 3 的 Fatalf 仍在 `:239`）：
   - `:29-33` 头注 ruler 3 一行摘要：`the ONE webview2.NewWithOptions call in / the host file must take`
     → `this package's production sources carry / exactly ONE ... call site, it lives in the host file, it takes`（同枚数行重写）。
   - `:201-206` ruler 3 doc：`that this file has exactly one NewWithOptions`
     → `that this package's production sources hold exactly one NewWithOptions site and that it is this host file's`（同枚数行重写）。
   - `:237-240` 断言块（原 4 枚→新 4 枚）：
     - 前：`if len(newWithOptions) != 1 { t.Fatalf("cmd/wisp/panel_host_windows.go has %d ... want exactly 1", len(newWithOptions)) }`（前有空行）
     - 后：`sites255 := newWithOptionsSitesInPackage255r6(t, fset, hostPkgDir255(t))` ＋
       `if len(newWithOptions) != 1 || len(sites255) != 1 { t.Fatalf(...) }`，红句同时印
       **整包枚数＋逐枚 `file:line` 名册＋host 枚数**（want 1 and 1）。
     - **方向、期望值、host 支全未动**：host 计数支原样保留（消费面 `newWithOptions[0]` 的
       WindowOptions/windowOptions 检查全部不碰），只把"恰好 1"的**扫描范围**并进来；
       `windowOptionsBodies`/int-literal/bringUp 三支**一字未动**（它们断的是宿主文件自身的消费形状，
       天生单文件——见 ⑤(e) 的名字诚实性：失败文案现在具名两级射程，不留 "this file" 冒充整包）。
2. **`cmd/wisp/panel_geometry_255r6_range_windows_test.go`**（新，103 行，测试件）：
   - `packageProductionSources255r6(t, dir)`：ReadDir→剔目录→留 `.go`→剔 `_test.go`→排序；
     目录读不到＝Fatal；**空名册＝Fatal**（没读到件的普查什么也不证明）。
   - `newWithOptionsSitesInPackage255r6(t, fset, dir)`：逐枚 `os.ReadFile`＋`parser.ParseFile(fset, ..., SkipObjectResolution)`
     **进调用方同一 fset**；读/解析失败＝Fatal 不豁免；命中形状与原尺逐字同形
     （`*ast.CallExpr` 的 `Fun` 是 `*ast.SelectorExpr` 且 `Sel.Name == "NewWithOptions"`——
     **不比 receiver 字面名**，所以 import 别名躲不掉，也天然不会因别名而假红；`A713` 顶回过的那种
     "词面钉"脆性在本尺射程外）。返回 `"file:line"` 名册，红句直接点名。

测试函数**名字保留**（`TestTicket255PanelHostBuildsItsWindowOptions`）：它现在量的不比名字少——
"PanelHost" 指的就是这个宿主包，失败文案具名两级射程。⛔ 没动票面、⛔ 没动台账、⛔ 没碰 `:377` 的
`retired` 名册（复跑后它照旧绿，见 `logs/10`）。

## ③ 正控清单（每种红指名到哪一支）

| 正控 | 种的形 | 预期/实测 | 红在哪一支（可分辨凭据） |
|---|---|---|---|
| PC-1 射程正控 | M1：兄弟产码件 `panel_resident_windows.go` 尾部种第二枚**标准形**调用点 | 旧尺（HEAD 形件回插同突变）**PASS＝盲**（`logs/11`）；新尺同突变 **FAIL rc=1**（`logs/12`，红句见 ④M1） | 红句印 `carry 2 ... [panel_host_windows.go:386 panel_resident_windows.go:558] and ... itself carries 1, want 1 and 1` ⇒ host 支 1、**整包支 2** ⇒ 红只能来自新增的整包普查支 |
| PC-2 名册/范围形 | M2：同件再种 `wv2.NewWithOptions`（**别名 import** 形） | **FAIL rc=1**（`logs/13`），名册点名 `panel_resident_windows.go:560` | 名册逐枚 `file:line` ⇒ 能分辨红在哪一枚文件哪一行；证明"只认 Sel.Name"的匹配面把别名形收进射程（⛔ 非只判"有红"） |
| PC-3 豁免支名册 | M3：新临时 `_test.go` 件里的调用点 | **PASS rc=0**（`logs/14`，纯形：resident 已还原并 hash 拉平后单跑） | 绿＝具名豁免成立（census 剔 `_test.go`）；这是 ⑤(a) 恒真面的实证，也是"红在哪一支"的另一半分辨力 |

## ④ 突变表（全部盘上种，⛔ 零 overlay 用作证据）

| 突变 | 盘上动作 | 结果 | 还原与拉平凭据 |
|---|---|---|---|
| M1 | `panel_resident_windows.go` 尾追 `probe255r6SecondCreateSite()`（6 行；追加在文件尾，不动 207 等名册行） | 见 PC-1 | `git show HEAD:... > 件` 后 `git hash-object` = `0f8912fdd3336c0f5f07d2bfef2b411066ecdbc9` = `git rev-parse HEAD:...` ✓ |
| M1b | 同一 M1 下把 `panel_geometry_255_test.go` 换回 HEAD 基线副本（旧尺盲测） | PASS（旧尺量不到兄弟件） | 交付副本从 `mutations/delivered-panel_geometry_255_test.go` 换回；此后未再动 |
| M2 | resident 再种：import 双名（`webview2` ＋ `wv2` 同路径，Go 合法）＋ `probe255r6AliasedCreateSite()` | 见 PC-2（红） | 同 M1 还原，hash 拉平 ✓（④M1 行同凭，收尾复量） |
| M3 | 新建 `cmd/wisp/zz_255r6_m3_windows_test.go`（测试件内调用点） | 见 PC-3（绿，具名豁免） | 内容存档 `mutations/m3-zz_255r6_m3_windows_test.go`，cmd/wisp 内件已移除（共享工作树⛔不留突变夹具）；`git status -- cmd/wisp` 收尾只见本腿两件 |
| 现场失误具名 | M3 首跑时 resident 的 M2 尚未还原 ⇒ 那次红因是 M2 非 M3，日志未采；随即还原＋hash 拉平＋**纯形重跑**（`logs/14`） | 已在流程中纠正 | 上面 M2 行的还原凭据即那次纠正的实测 |

`panel_host_windows.go` 起手/收尾 `git hash-object` 均 `58e2b155dbcc8d380793e90bffc5514ab00743f1`
＝HEAD blob ⇒ 产码零留痕。`panel_geometry_255_test.go` 现值 `ee94bbc4…` ≠ HEAD＝本腿交付改动（预期）。

## ⑤ 恒真面具名（哪一形怎么坏都照绿）

- (a) **`_test.go` 内的第二调用点照绿**——设计内豁免，名册有实证（PC-3），doc 与头注都具名。
- (b) **本包外的产码件照绿**（如 `scripts/spike/.../main.go:149` 今天就是绿的第二枚）——
  范围名已写进失败文案与测试件头注（"this package's production sources"），⛔ 不留 "全工程" 冒充。
- (c) **间接形照绿（推演未实种）**：`f := webview2.NewWithOptions; f(...)` 的第二枚调用是裸 Ident，
  不匹配 SelectorExpr 形 ⇒ 绿。未实种（种它需要再动产码件），诚实记为推演。
- (d) **census 的"解析失败＝Fatal"支在 `go test` 下不可达**（编译门在前，解析不出的件根本构建不出
  测试二进制）⇒ 该支是防"以读不到为豁免"的保险，不是可实证的牙。
- (e) **天生单文件没扩的三支**：`windowOptionsBodies`/int-literal/bringUp 仍只读宿主件——它们的断言
  主语本来就是"这个宿主文件的消费形状"，文案未冒充全包；第二枚 `windowOptions` 声明在包内会是
  编译错（同名方法冲突）或消费形检查咬不到——这一形今天无盘上证据，具名未验。

## ⑥ 门禁逐条 rc

| 门禁 | rc | 凭据/数 |
|---|---|---|
| `go vet ./cmd/wisp/`（带 DLL harness） | **0**（零输出） | 直打 rc，无管道 |
| `sh scripts/d22scan.sh` | **0** clean | `logs/16`；ban #8 明写扫 `cmd/` **112 Go files, comments and _test.go included** ⇒ 本腿两件在被扫面内且过 |
| `gofmt -l 两件` | **0**（零输出） | 前置 CR 尺：`tr -cd '\r' < 件 \| wc -c` 工作树 **0**／HEAD 基线 **0**（新件 0）⇒ 无幻影可辩 |
| 靶向族 `go test -run 'TestTicket255' -v` 交付后 | **0**，34 PASS / 0 FAIL | `logs/10` |
| 名册逐枚复跑（票 255 `file:LINE [token]`，源＝`cmd/wisp/config_readers_255.go`，现量 **14 枚**） | **14 HIT / 0 MISS** | `logs/15`；A713 的 14/14 vs 2-MISS 未裁读数，本腿现量＝14/14 |
| 整包 `go test -count=1 ./cmd/wisp`（harness，跑到终态，⛔ 非 `-run` 单跑） | **1**（预期：本仓有在册红） | `logs/20`：836 行，`FAIL cmd/wisp 525.468s`，**零 `panic:`／零 `0xc0000135`** ⇒ 用例真跑了 |

整包逐名作差（`logs/20` vs 名册现抽 `/tmp/roster111.txt`，源＝`.scratch/wisp/probes/111/c2/logs/03-roster-rerun-vs-33v4.md`，本腿自抽 16 枚唯一名）：

| 本发红的（5 枚） | 在名册？ | 255-v2 在 HEAD 的耗时 | 本发耗时 |
|---|---|---|---|
| `TestPanelHostRealWindowHopAndLifecycle` | 在 | 5.08s | 5.10s |
| `TestAC4FocusReturnToPriorWindowGap33r5` | 在 | 5.15s | 5.17s |
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | 在 | 20.01s | 20.02s |
| `TestAC14AwaitedBindingReplyReachesThePage` | 在 | 20.01s | 20.01s |
| `TestAC14GoSideEvalPushReachesThePage` | 在 | 20.01s | 20.02s |

⇒ **5/5 ⊂ 名册、0 枚名册外、`--- FAIL: TestTicket255*` 零枚**（255 一族在整包跑无一红；其全绿另有 `logs/10` 靶向 `-v` 34 PASS/0 FAIL 正证）⇒ **新增红＝0 成立**。`internal/panel` 的 2 枚有意保留红灯不在本写面射程（⛔ 未跑 `./internal/...`，本腿整包只 `./cmd/wisp`）。

## ⑦ 这一格今天闭合没有

**按本腿锚闭合**：`TestTicket255PanelHostBuildsItsWindowOptions` 的"恰好 1"现在量的是它名字宣称的
整包产码射程；同一颗突变（M1）下**旧形尺照绿、新形尺指名红**（`logs/11` vs `logs/12`），红句带逐枚
`file:line` 名册可分辨红在哪一支；别名形（M2）在射程内、测试件豁免（M3）有实证；全部突变盘上种、
盘上还原、`git hash-object` 与 HEAD 逐枚拉平。⛔ 不闭合的部分本件没含糊：`255-v2` Q2 点名的
**reshow hint 尺（`TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne`）同族单文件盲**
——不在派单七枚行锚的射程里，本腿没动它，那一格**今天没闭合**（见 ⑧ 顶回第 3 条）。

## ⑧ 派单没让我量到/我主动具名搁置的

- **`255-v2` Q2 点名的另一枚 "the only one" 尺**（`panel_reshow_255r1_windows_test.go` 的
  HintNone/HintFixed resize 普查）**不在本腿锚内**：派单七枚行锚逐枚复认落在 `panel_geometry_255_test.go`
  的 NewWithOptions 尺（锚为准，§0 定式）；HEAD 台账 `A714` 那句 "把那把 the only one AST 尺扩到整包"
  按**字面名**更像指 reshow 那枚——**名与锚在本仓同时指两格**，我按内容锚交付了 geometry 尺，
  reshow 尺的扩程**今天仍没做**，需编排者另派或翻本件顶回。
- M3 首跑的混种失误已具名（④ 末行）；`census 解析失败支` 不可实证已具名（⑤(d)）。
- `-race`／`-tags winlive` 档：整包 harness 不含，本腿未跑（与 255-v2 同格口径）。

## 顶回编排者（具名）

1. **HEAD 已漂**：派单转述 `c4b8e2eb`，起手现量 `040e424d`（未采转述号，锚件已记）。
2. **"现量枚数以你自己跑的为准"兑现**：票 255 名册我从 `config_readers_255.go` 自抽，现量 **14 枚**、
   **14/14 HIT**；这不是抄"14"——A713 那两种读数在我这把只有前者复现。
3. **名/锚歧义顶回**（⑧ 第一条）：若编排者本意是 reshow 的 hint 尺，本腿交付**没覆盖它**，
   别把本件当两枚都扩完了。
4. 派单里"把注释改成双空格形"那条已撤回的错派：本腿未动 `:377` 名册一字（`logs/10` 复绿为凭）。
