# 327-a1 §② `Cf` 逐枚可达 ＋ §③ 上游洗没洗 ＋ 必答题

腿 `327-a1`｜锚 `44f51fc8`｜2026-10-16｜⛔ 判语、⛔ 修法建议（`AC#1` 那张表归编排者裁）。决定性尺逐字转储＝`logs/readings.txt`（摘录，全量 497 行在仓外 `D:\tmp\327a1\logs\readings-raw.txt`；探针原始 stdout `D:\tmp\327a1\logs\cfprobe-02.txt`、stderr 0 字节）。

## 探针声明（自报跑了什么）

`go env GOVERSION GOROOT` → `go1.27.1` / `D:\work\base\go`（rc=0，非测试非构建）。
**唯一执行过的代码**＝仓外单文件纯 stdlib 探针 `D:\tmp\327a1\cfprobe\main.go`，命令 `go run ./cfprobe/main.go`（rc=0，stdout 7.6 KB 已转储进 `logs/`，stderr 0 字节）。它 import 的只有 `encoding/json fmt strings unicode unicode/utf8`，⛔ 本仓任何包。
本腿**未**跑 `go test`／`go build`／`go vet`／`gofumpt`／门禁（这波编译面归 `165-v2`）。
§① 与 §② 里凡"洗不洗"的结论都来自**读码（`sed -n`／`grep -n`）**；探针只量字符类别归属与 `encoding/json` 的转义集合，其中 `needsDisplayWash`／`scrubForDisplay`／`writeControlEscape` 三段是 `internal/panel/display_scrub_216.go:100–147,152–161` 的**逐字重打**（件头写明 RETYPED，⛔ 当作仓内函数跑过的证据）。正控同屏：`ESC`/`LF` 两枚 `Cc` 在同一探针里 `needsDisplayWash=true`、`scrubChanged=true`、`json` 转义＝`true` ⇒ 这把尺不是恒 0。

## §② 逐枚：今天真可达的码点（可达＝下面两枚具名入口之一；两枚入口都**不按字符筛**）

- 入口 A（磁盘真名 → `composer.workspace.canonical`）：`cmd/wisp/panel_pump.go:128` → `internal/panel/workspace.go:119,120`（`res.Canonical` 逐字）→ `display_scrub_216.go:169`（`IsControl` 筛，`Cf` 不响）→ 快照 `json:"canonical"`。
- 入口 B（模型正文 → 审批卡／回复行）：`internal/llm/events.go:239`（`ArgsDelta` 逐字累积）→ `internal/tools/bridge.go:382 Args: req.Args` → `cmd/wisp/panel_pump.go:68`＋`:297–311`（`argv`／`command`／整段原始 JSON）→ `internal/panel/approval.go:80 Args: append(…)`（逐字）→ 快照；或 `cmd/wisp/run.go:1331 stream.Append` → `internal/panel/pump.go:537 cur.Text += text` → 快照 `results[].text`。
- 筛子读数：`internal/risk/pathresolver.go` 全文 `ContainsAny`／`unicode.Is*`／字符类拒绝 **0 处**（尺 R56，命中 5 条全是注释与 `Actable`/`Resolve` 等函数签名行）⇒ 入口 A 对 `Cf` 不设卡。

| 码点 | 名 | 类别读数（探针） | `needsDisplayWash` | 可达判定 |
|---|---|---|---|---|
| `U+202A` LRE／`U+202B` RLE／`U+202C` PDF／`U+202D` LRO／**`U+202E` RLO** | 双向嵌入/覆盖/弹出 | `Cf`,`IsControl=false`,`IsPrint=false` | **不响**（探针逐枚 `false`） | **可达**（A＋B） |
| `U+2066` LRI／`U+2067` RLI／`U+2068` FSI／`U+2069` PDI | 双向隔离符一族 | 同上 | **不响** | **可达**（A＋B） |
| `U+200B` ZWSP／`U+200C` ZWNJ／`U+200D` ZWJ | 零宽一族 | 同上 | **不响** | **可达**（A＋B） |
| `U+2060` WJ／`U+FEFF` ZWNBSP/BOM | 零宽连接/BOM | 同上 | **不响** | **可达**（A＋B） |
| `U+00AD` SHY | 软连字符 | 同上 | **不响** | **可达**（A＋B） |
| `U+2061`–`U+2064`（FA/IT/IS/IP） | 数学作用符一族 | 同上 | **不响** | **可达**（A＋B）——票 327 `AC#1` 说这几枚"具名裁要不要"，本腿只报可达性 |
| `U+061C` ALM／`U+180E` MVS／`U+206A`（已弃用形） | `Cf` 类其余 | 同上 | **不响** | **可达**（A＋B）；`AC#1` 点名"逐枚要有名字"，这三枚在此具名 |
| `U+001B` ESC／`U+000A` LF（**正控**，非 `Cf`） | `Cc` | `Cc`,`IsControl=true` | **响** ⇒ 洗成 `\x1b`／`\n` | 可达且**今天已被洗**（同一入口，边界只对 `Cc` 有效——这条正是 `AC#2`"两把尺同一字符集"要钉的那对） |

⇒ **本票候选名单里没有一枚是"仅理论形"**：两枚入口都不按字符筛，所以 22 枚候选逐枚都给得出 `file:line`。唯一"仅理论形"的说法要留给**尚未落地的字段**——插件名／技能名／命令说明三枚今天零源（见 §① 表），那是"没有源"不是"源里有筛"。
⚠ OS 侧（NTFS/WebView2）对这些码点的行为本腿**未量**、⛔ 断言；上表只说 Go 侧这条路不设卡。

### AC#0 要求的三形读数（探针输出逐字，⛔ 我复述成句）

```
① 只含 U+202E：len(bytes)=19 runes=17  utf8=70 68 6F 74 6F E2 80 AE 64 6E 6F 6D 65 6C 2E 78 74 6E 70
   needsDisplayWash=false scrubChanged=false  → 快照里就是这 19 字节，E2 80 AE 原位不动
② 只含 U+200B：len(bytes)=16 runes=14  utf8=70 68 6F 74 6F E2 80 8B 6E 61 6D 65 2E 70 6E 67
   needsDisplayWash=false scrubChanged=false  → 同上，E2 80 8B 原位不动
③ 只差一枚 U+200B 的两枚名字：report.png = 72 65 70 6F 72 74 2E 70 6E 67（10 字节）
   report<U+200B>.png = 72 65 70 6F 72 74 E2 80 8B 2E 70 6E 67（13 字节）
   逐字节相同？ false（差 3 字节 E2 80 8B）；scrub 后仍 false：c2 after scrubForDisplay = "report\u200b.png"（那是 Go 的 %q 打印形，⛔  packet 里的字节）
   ⇒ 快照带的是两枚**不同**的字节串，屏幕上是同一枚名字。票面说的"看起来一模一样其实不是同一个"这半成立；
     "逐字节相同"那半**不成立**，本腿据实更正。（"是不是同一个名字"另有两处按串等判定：
     `attachments.go:201 ID:"att-"+shortHash(DisplayName+"|"+MIME)` 按字节哈希，两枚名字得两个 ID——⛔ 本腿不判它算不算缺陷）
```

## §③ 上游有没有已经替我们洗掉了一半——**没有**

逐条量（每条给尺，⛔ 混"判定侧忽略"与"显示侧去掉"两件事）：

1. **`internal/risk/taintmatch.go:39–45 isIgnorableTaintRune`**（`0x00AD,0x200B,0x200C,0x200D,0x2060,0xFEFF`＋`unicode.Mn`）**只作用于判定输入**：尺 R57 `grep -rn 'normalizeTaint\|isIgnorableTaintRune' --include=*.go internal cmd`（滤 `_test`）rc=0，调用者**全部落在 `internal/risk` 内**——`hostpath_185.go:67`、`provenance.go:559,578,841`，加定义处 `taintmatch.go:49,59`；`internal/panel/**` 与 `cmd/wisp/**` 命中 **0 枚**。⇒ 它会把那六枚**从判定用的归一化串里丢掉**，但字进屏幕的路**不经过它**。`U+202E` 与 `U+2066–U+2069` **不在**这份表里（尺＝`sed -n '25,60p'` 原文，转储 logs）。⛔ 动它（票 327 禁区第一条）。
2. **`attachments.go:301 ContainsAny(name, "\x00\r\n\t")`**＝只拒 4 枚 `Cc`，`Cf` 不在此列；且这条守卫今天**零生产调用者**（尺 R40：`.Ingest(` 19 命中全在 `_test.go` 与 `.scratch` 探针）。
3. **`sanitizeDisplayName`（`:310–312`）＝`TrimSpace`＋取末段**：探针末节量得 `TrimSpace` 对 `U+202E/U+200B/U+FEFF/U+00AD/U+2060` **一枚都不去掉**（`IsSpace=false`，逐行 `false`）。
4. **`cmd/wisp/approval_reply.go:609–619 sanitizeReplyText`**＝`IsControl`→换空格，`Cc`-only；调用者是 `:571/:575` 两枚 `reject/panelReject` 分支（尺 R60），去向是**进模型的工具结果**，不是出向显示链。
5. **`internal/risk/rules_shell.go:57 findMetachars`**＝只把 `shellMetachars` 里的 `\n\r\0` 写成可见形，产物进 `reason` 句子（`:37`）；它**不碰 argv 本身**，`Args` 仍逐字出向（`approval.go:80`）。
6. **`json.Marshal`（`pump.go:344`、`:359`）**：探针表末列逐枚 `escaped=false`——22 枚候选 `Cf` 全数原样进 JSON 字节；同一探针里 `ESC`→`\u001b`、`LF`→`\n` 是 `true`（正控）。⇒ 这层只转义 `Cc`，不转义 `Cf`。
7. **`StreamLog`（`pump.go:524–538`＋`:463–469 render`）**：按** rune** 做首尾截断，无字符筛；`Append` 是 `cur.Text += text` 逐字累积。

⇒ **§③ 的答案是"没有"**：今天没有任何一条生产路径在字进屏幕之前把那批字符物理去掉或转义。

## ★必答题："一枚模型输出的文件名带 `U+202E`，今天落到屏幕的整条路上，哪一环会去掉或转义它？"

**答：没有一环。** 每一环的 `file:line`（尺全在 `logs/readings.txt`）：

| 环 | `file:line` | 对 `U+202E` 做了什么 |
|---|---|---|
| 1 模型 SSE 增量累积成 tool-call args | `internal/llm/events.go:239` | 逐字 `append(…, ev.ArgsDelta...)`，无筛 |
| 2 工具桥把 args 交给判定 | `internal/tools/bridge.go:382` | `Args: req.Args` 逐字带过 |
| 3 审批队列条目 → 面板读 | `cmd/wisp/panel_pump.go:68`＋`panelArgs` `:297–312` | 取 `argv`／`command`，取不到就 `[]string{raw}`（整段原始 JSON 逐字） |
| 4 渲成卡视图 | `internal/panel/approval.go:80`（`Args`）、`:83`（`Reason`） | `append` 逐字；`Reason` 注释自陈 verbatim；`:84 ReasonKnown` 只 `TrimSpace` 判空 |
| 5 组快照（显示边界在此） | `internal/panel/composer.go:145` | **只洗 `Composer` 段**（`display_scrub_216.go:168–179` 两处赋值），`Pending`/`Results` 不进这枚函数 |
| 6 就算它进了这枚函数 | `display_scrub_216.go:101,114,155` | 判据 `IsControl`，`Cf` **不响**（探针 22/22 `false`）⇒ `return s` |
| 7 编码出向 | `internal/panel/pump.go:344`／`:359` `json.Marshal` | `Cf` 原样（探针 `escaped=false`） |
| 8 最后一公里（Go → WebView） | `cmd/wisp/run.go:732 Out: rt.bookPanelSnapshot` → `cmd/wisp/panel_pump.go:360–366` | **快照只被存进 `rt.lastSnap/lastSnapBytes` 再加一条日志，不推给页面**；唯一读它的 `lastPanelSnapshot`（`:424`）**生产调用者 0 枚**（尺 R64：`grep -rn 'lastPanelSnapshot' --include=*.go cmd internal` rc=0，共 **17 行命中＝1 定义＋1 注释行＋15 枚全在 `cmd/wisp/*_test.go`**）。快照也没有"页面来问"的那扇门：入向方法常量共 7 枚（尺 R63a `grep -rn 'Method[A-Za-z0-9]* *= *"' --include=*.go internal/panel` rc=0＝`bridge.go:42,43,44,45,66,67,75`，与 `composer_dispatch.go:200,205,210,215,220,225,230` 的 7 个 `case` 逐枚对上，`grep -c 'case Method'`＝7），**没有一枚是快照读取**（尺 R63b `grep -rniE 'MethodSnapshot|snapshot.*= *"panel\.' --include=*.go internal cmd` **rc=1**；正控＝同族的 R63a 那把尺 8 行命中）；`grep -rn 'AddHostObjectToScript' --include=*.go cmd/wisp` **rc=1**（正控：同族尺对 `webview` 在 `panel_host_windows.go` 命中 45 行）。出向 Eval 那形今天只承载**RPC 回执**：尺 R66 `grep -rnE 'ExecuteScript\|CallJs\|PostWebMessage\|NavigateToString\|evalScript' --include=*.go cmd/wisp` rc=0 共 6 行＝1 行生产（`panel_control_windows.go:665 SetHtml→NavigateToString`）＋5 行测试（`panel_transport_35r2_test.go:1206/1446/1719/1887/1888`，其 1887 行的红句写的是库把回执发成 `window._rpc[0].resolve(...)`）。⚠ 两句自陈各归各的射程、⛔ 我读成一件事：`internal/panel/pump.go:366`"快照已构造但没有出口（最后一公里在票 33/35 名下）"射的是 `src.Out == nil` 分支（生产里 `Out` **不**是 nil）；`cmd/wisp/panel_pump.go:440` 的 "does not exist yet" 射的是推送频率的 merge/backpressure 规矩 |
| 9 页面侧渲染器（就位、进料缺） | `frontend/src/App.tsx:91 { snapshot = EMPTY }`、`approval-screen.tsx:103`、`l2-approval-card.tsx:305/152/160/164`、`composer.tsx:150/152`、`App.tsx:461–463`、`result-stream.tsx:37` | React 文本节点不做 bidi；`:160` 那枚 `JSON.stringify` 按 JS 规范只转引号/反斜杠/`U+0000–U+001F`。⛔ 本腿未在浏览器里量（超只读射程） |
| 模型正文那条支线 | `cmd/wisp/run.go:1328–1332` → `internal/panel/pump.go:524–538` | 逐字累积，无筛 |

⚠ **读数附带一句（不是判语，是名册的边界条件）**：环 8/9 那两枚读数说明"今天真走到**像素**上"这条路其实断在 Go 出向之后。所以本名册把"可达"定义为**到快照字节**（票 327 `AC#0` 的原话正是"进到快照里是什么字节"），并把"最后一公里未接"具名挂在环 8——它同时改变 `AC#0` 里"屏幕上哪一行"那一列的性质：那一列量的是**渲染器就位**，不是"像素已出现"。要不要因此重排加尺顺序＝编排者裁，⛔ 本腿排程。
