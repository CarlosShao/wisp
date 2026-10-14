# 票 114 验收裁决表 — 腿 `114-v2`（非实现者；只读＋仓外真种盘突变＋重跑；面＝internal/panel＋cmd/wisp 只读）

- 被裁：`114-r6` 产码 `3b6f170f`（单枚新件 `internal/panel/bound_dispatch_door_114_test.go` 530 行）＋交付 `c499adac`；钉主张＝第三装法两侧同名可达性（host `w.Bind(<包级字符串常量>,回调)` ⇔ 页 `window.<常量值>`，只许 `panel.ts` 读；reviewed 空即红；旧 `.postMessage(` 尺同树仍盲）。
- 锚：起手 HEAD `7e8ac3fa`（2026-10-14 16:5x +08）；窗口内他腿 `318-r1`（`internal/agent/approval`）在飞；5 枚关键件 sha256 起手＝收手（`45c1059b…`／`2b7a9fd0…`／`d2b2e848…`／`6924b8e1…`／`016b6943…`）。
- 方法：仓内零突变；牙全在仓外树 `D:\tmp\wisp114v2\tree`（`git archive HEAD` 真盘；A713：钉读盘、overlay 换不到 `os.ReadFile`，⛔ 不以 overlay 绿抵账）；logs＝`probes/114/v2/logs/01…04`。⛔ 框未翻、未动 `Status:`/票名、零 push。

## 1 判语（结论先行）

| 项 | 判语 | 一句依据 | 翻框建议 |
|---|---|---|---|
| **AC#10** | **成立**（重裁口径：钉认第三装法） | 第三装法被两侧同名读盘钉住（host 36 枚产码取 `:842 w.Bind(panelDispatchBinding,…)`、常量 `:80`＝`"wispDispatch"`；页 `panel.ts:350 window.wispDispatch` 且全 `frontend/` 只许 panel.ts）；我 9 发独立种形全按预期＝8 红（摘 Bind／拼写不一致／页读改注释／第二文件读门／reviewed 空×2／单参翻双参／根消失）＋1 绿对照（声明搬文件）；reviewed 空三形全红；旧尺同树两发 PASS＝盲区现场复现 | 可翻；留档句①②（§7） |

## 2 钉的射程边界核（四条，各抽一行现读；逐条以我自己的种形独立验过）

1. **skipped 单参调用**：钉 `bound_dispatch_door_114_test.go:367` `if len(call.Args) < 2 {` → `:372` `rep.skipped = append(…)`，不判；真站点＝`cmd/wisp/resident_task_source_windows.go:311` `ws.Bind(run.paths)`；live 绿句 `1 skipped non-door .Bind call(s)`。反向独立核（t7）：翻成双参 `ws.Bind(run.paths, nil)` ⇒ rc=1 逐字点名 `the w.Bind argument run.paths is not a plain name (a string literal or a package-level constant)…`（边界两侧都咬、fail-closed）。
2. **literal 不入 reviewed**：`:171` `if b.constBound && b.problem == ""`；live 绿句列 `:892 binds "wispProbeRT" via "wispProbeRT"`，门相关读里无它（literal 探针门不被要求页读）。正向独立核（t5）：把真门实参改成 literal `"wispDispatch"`（页读仍在、两 Bind 都在）⇒ reviewed 空、rc=1（literal 顶不了常量＝刻意窄口径）。
3. **注释剥离**：`:99 doorReadCodeLines`＋`:77 doorReadCommentOpener`（`//`、`/* */`、`<!-- -->`）；`:329` 的 prose 在 `:320–:346` 的 `/* */` 块内。独立核（t3）：把 `:350` 真读注释掉、prose 原样 ⇒ rc=1 `never reads window.wispDispatch`，读清单只剩 `:161`×2（prose 顶不了数）。
4. **only-panel.ts**：`:250` `if !rd.inPanel && r.bindsName(rd.name)`。独立核（t4）：往 `frontend/src/lib/config-view.ts` 追加一行 `window.wispDispatch` ⇒ rc=1 逐字点名 `config-view.ts:168 … is a second door`。

## 3 九发牙（仓外真树；红/绿句逐字在 logs/03；每发先还原 pristine 字节）

| # | 形 | 读数 |
|---|---|---|
| t1 | 摘常量绑的 `w.Bind` 整条语句 | rc=1；floor 逐字点名（剩余站点列 `:886 wispProbeRT`） |
| t1b | 同树旧 `.postMessage(` 尺 | **PASS** rc=0（`113 files scanned, 2 host call sites`）＝旧射程盲 |
| t2 | 常量声明搬家到 `cmd/wisp/run.go` | rc=0 绿（跨文件包常量解析成立；Bind 站位移到 `:841`）＝诚实重构被放过 |
| t2b | 声明/调用拼写不一致（两边都在） | rc=1；两半各点名：`names panelDispatchBinding, which no package-level string constant in cmd/wisp declares…` ＋ floor 列 `binds "<unreadable>"` |
| t3 | 页真读改注释（prose 不动） | rc=1 `… never reads window.wispDispatch …`（注释剥离承重） |
| t3b | 同树旧尺 | **PASS** rc=0（页侧改动仍盲） |
| t4 | 第二文件读门 | rc=1 点名 `frontend/src/lib/config-view.ts:168 … is a second door` |
| t5 | reviewed 空形 A：实参换 literal（两 Bind 都在） | rc=1 floor（`w.Bind sites read:` 列两枚 literal 站点） |
| t6 | reviewed 空形 B：`.Bind` 全改 `BindDoor`（binds 全空） | rc=1 floor（`w.Bind sites read: (none)`） |
| t7 | 单参 skip 翻双参 | rc=1 点名 `:311`（边界 fail-closed） |
| t8 | `cmd/wisp` 整目录消失 | rc=1 `cmd/wisp is not walkable: …` ＋ floor `(none)` |

复跑驱动：`D:\tmp\wisp114v2\mut.py`（每发还原 pristine；终态 host sha256 `45c1059b…`＝live 同值）。

## 4 基线与对照

- live 真树新钉 **rc=0**（绿句逐字 logs/01：`36 cmd/wisp production file(s), 2 w.Bind door site(s) […:842 binds "wispDispatch" via panelDispatchBinding; …:892 binds "wispProbeRT"…], 1 skipped…, 27 window.* read site(s)…`）。
- live 整包 `225/144/5/0 rc=1`，5 红与 `114-r6`/`114-v1`/a5 记录**同集**；抽 2 名逐字：`composer_test.go:78: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`／`tokens_fourway_test.go:441: read design/assets/tokens.css: … The system cannot find the path specified.`（窗口内界面侧/既有红，非本波）。
- 同树 A/B（logs/02）：A（有新钉）`225/142/7` vs B（移走新钉）`223/141/7`，**FAIL 名集 diff rc=0＝SAME SET**；7 红＝live 5＋2 枚树工件（archive 导出 CRLF 而 fixture 尺块分隔符写字面 `" -->\n"`；树无 `.git`；A/B 同集即证与钉无关）。
- 旧 `.postMessage(` 尺在 t1/t3 两棵突变树上均 PASS＝旧射程盲、新钉有牙的现场复现（r6 主张的三树盲区我以自己的两发独立复现；其 m1/m2/m3 具体 span 与数值未复核）。

## 5 门禁（现跑，live；logs/04）

build `./internal/panel ./cmd/wisp` rc=0；两枚相关 `-run`（新钉＋旧尺两枚用例）rc=0；`sh scripts/d22scan.sh` clean rc=0（ban #8 internal/=556 含本件）；长度门 `VERDICT GREEN`（tracked=12876；over-budget 60＝roster 60，not-in-roster=0）rc=0。

## 6 与 114-v1 链的相容性（一句）

与 `114-v1` 链**相容**：v1 停手上报给出的现读链（真门＝`:80` const＋`:842 w.Bind`、旧钉正则 `composer_test.go:441` 只认 `.postMessage(`）正是 r6 按编排者重裁落地的对象；我复跑同一链同锚（sha256 三枚同值、绿句同名同址、`:892` 探针门不入 reviewed），无冲突。

## 7 留档句（翻框随附）

① 本钉 reviewed 的"常量"是刻意窄口径：literal 绑的门不入 reviewed（探针门例外）；代价＝将来"以常量绑的合法探针门/第二门"必须配 `panel.ts` 读或改口径——本钉可预见的唯一摩擦，登记备查。
② AC#10 原信"第一种（`AddHostObjectToScript`）要么显式禁用、要么进同一枚钉"：我复跑 `git grep AddHostObjectToScript -- '*.go' ':!*_test.go'` **零命中**（`frontend/src` 唯一命中＝`panel.ts:146` 注释）⇒"机制不存在"为真；但**本钉只在红句散文里陈述该事实、未把它做成断言**，按重裁"注释那半归前端"⇒该半支现状＝事实零调用点＋无门。若编排者要"零调用点"本身成门，属新仪器、不在本钉射程。

## 8 具名没做

- 未跑 `go test ./cmd/wisp`（314 死区，派单禁）／POSIX 三合包／真窗；未动 `cmd/wisp/**` 产码与 `frontend/**` 一字节；未翻框、未动 `Status:`/票名、零 push。
- 未复核：r6 的 m1 span `[42508,42619)` 与 m2/m3 具体偏移（其结论由我自己的九发独立复核，不复用其数值）；`panel.ts:146/329` 注释文本那半（重裁归前端）；junction 形（198/201 面）。

交件：本表＋`probes/114/v2/00-anchor.md`＋`logs/01…04`；票 114 尾 1 bullet。
