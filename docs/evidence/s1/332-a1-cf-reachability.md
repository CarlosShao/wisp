# 票 332 只读普查腿 `332-a1` 交件自证（`Cf` 全族 170 枚里显示边界外那 155 枚的可达性读数）

腿＝`332-a1`（只读普查：⛔ 产码 ⛔ 测试件 ⛔ `go test`/`go build`/`go vet` ⛔ 判语 ⛔ 翻框 ⛔ push）。
锚＝具名提交 `7535d83b`（dev）。料在两件里：
`.scratch/wisp/probes/332/a1/census-1-cf-155-roster.md`（AC#0 四把尺＋AC#1 名册）
`.scratch/wisp/probes/332/a1/census-2-carriers-and-readers.md`（AC#2 逐枚载体＋AC#3 正控＋腐烂清单）；
原始读数 15 枚 `.log` 同目录 `logs/`（逐枚已入库，枚枚一笔提交）。
**射程口径**（先钉死，⛔ 越过去读）：本件的"可达"＝**枚枚走到 `json.Marshal` 出的快照字节里**，
今天那一跳已进 WebView2 的 `Eval` 调用串（`cmd/wisp/panel_snapshot_deliver_329.go:105` → `cmd/wisp/panel_host_windows.go:336`），
但页面接收器仍 0 命中（`git grep -c wispOnPanelSnapshot 7535d83b -- frontend`＝**0 枚文件**，正控＝去掉目录限定同一把尺命中 **16** 行）
⇒ ⛔ 任何一句不构成"屏幕上已出现像素"。

## 一、交件自证 ①：逐枚提交号＋逐枚 `git show --name-only`

尺＝`git log --format='%h' 7535d83b..HEAD -- .scratch/wisp/probes/332/a1`（射程＝本腿写面），
逐枚展开已存 `logs/commits-roster.log`（65 行、16 个块）。**交件时点复量的两把尺**：
越格写面命中＝`grep -cE '^\s+(internal/|cmd/|frontend/|design/|go\.mod|go\.sum|\.gitignore|docs/PLAN\.md|docs/specs)' logs/commits-roster.log` → **0**。
本波**产码＝0 枚文件**（只读普查腿），16 笔全部只落
`.scratch/wisp/probes/332/a1/census-*.md`（2 枚）与 `.scratch/wisp/probes/332/a1/logs/*.log`（14 枚，含本件写完后再追加的一枚名册）
＋票 332 末尾**一条** bullet＋本枚 `docs/evidence/s1/` 件。

## 二、交件自证 ②：件路径＋`wc -c`＋行数＋占位符尺

| 件（绝对路径） | 字节 | 行 |
|---|---|---|
| `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\332\a1\census-1-cf-155-roster.md` | 9136 | 79 |
| `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\332\a1\census-2-carriers-and-readers.md` | 15586 | 137 |
| `D:\work\workspace\projects plans\Wisp\docs\evidence\s1\332-a1-cf-reachability.md`（本件） | 133 行，字节数以最后一笔提交的 `wc -c` 现量为准（本腿在交件报文里逐字给出） | 133 |
| `...\probes\332\a1\logs\`（15 枚 `.log`，逐枚字节见 `logs/commits-roster.log` 之外的 `wc -c` 输出，总数 168207 字节） | | |

占位符尺（分词形）＝`grep -rniE 'TODO|TBD|FIXME|待填|占位|placeholder|XXXX|\.\.\.$' census-1*.md census-2*.md` → **0 命中**（rc=1＝真零）。
★过程中本腿自纠一处：件2 第 114 行原写着转义词汇名（四个 `X` 那一形），被这把尺扫到 ⇒ 追加一笔 `4f897c7ab` 改成逐字描述，**读数一字未动**（⛔ amend）。
0 字节文件尺＝`wc -c` 全部 ≥711 字节，无空件。

## 三、交件自证 ③：决定性那把尺的逐字输出（起手侧 ⇄ 交件侧，带 `file:line`）

**起手侧（名单，量的）**——尺＝`git show 7535d83b:internal/panel/display_scrub_216.go | grep -n 'case 0x'`：
```
166:	case 0x00AD,
```
★枚数**不是** 1：`grep -c 'case 0x'` 回 1 是因为 15 枚写在**同一条跨行 `case`**里（`:166-170`）。
逐枚枚出＝**15**＝`00AD 200B 200C 200D 2060 FEFF 202A 202B 202C 202D 202E 2066 2067 2068 2069`
（全文＝`logs/roster-display-scrub-216-at-7535d83b.log`，名单被 `:123` 写入侧与 `:214` 判据侧各调一次）。

**交件侧（`Cf` 全集，量的、独立复跑，⛔ 抄 `327-v1`）**——唯一执行过的代码＝
`D:\tmp\332a1\cfprobe\main.go`（`go run` rc=0、stderr 空、`go1.27.1 windows/amd64`、`GOTMPDIR=D:\tmp\332a1\gotmp` 先建后跑），逐字读数：
```
GO_UNICODE_VERSION=17.0.0
CF_TOTAL=170
CF_IN_ROSTER=15
CF_OUTSIDE_ROSTER=155
ROSTER_SIZE=15 ROSTER_NOT_CF=0
MISMATCH_COUNT=0
```
`MISMATCH_COUNT`＝本探针的自证列：170 枚逐枚送进**逐字重打**的 `scrubForDisplay`（重打自
`display_scrub_216.go:106-136/164-175/179-199/211-220`，函数体一字未改），比"名单说它洗不洗"⇄"它实际被不被换形"，**0 枚不一致**
⇒ 名册里"洗／穿"那一列是量出来的。两形正控都在盘上：`00AD ⇒ out="\xad" WASHED`、`0600 ⇒ out="\u0600" RAW`。
第二来源对合（算的）＝unicode.org 官方 `UnicodeData.txt`（17.0.0）里 `;Cf;` 行 **170** 枚，与 Go 那 170 枚逐枚比＝
`IN_GO_NOT_UD=[]`、`IN_UD_NOT_GO=[]`；逐枚底表（码位／官方名／Category／IN|OUT／WASHED|RAW）＝
`logs/ac1-cf-170-roster-with-unicode-names.log`（171 行＝1 表头＋170 枚，⛔ 抽样）。

**AC#2 的决定性尺（量的）**：三形入口没有任何一处按字符筛——
`git grep -n 'ContainsAny\|IndexAny\|ContainsRune' 7535d83b -- internal cmd | grep -v _test`＝**18 行**
（逐行原文＝`logs/ruler-no-charclass-screen-on-carriers.log`，字符集合枚枚只点 ASCII 标点／`\x00\r\n\t`／shell 元字符
`"|;&><`$()^\n\r\x00"`）；`internal/risk/pathresolver.go` 里字符类拒绝 **0 处**
（★成对正控＝同一把正则对 `display_scrub_216.go` 命中 **5** 行 ⇒ 那枚 0 不是尺失灵；逐字见 `logs/ruler-ingest-zero-with-control.log`）。
⇒ 155 枚分 11 族（求和 13+1+1+2+4+6+3+16+4+8+97＝**155**，LEFTOVER=[]、GHOST=[]）**逐族三形载体都＝有**，
载体枚枚带 `file:line`（件2 的表）；唯一写"无已知载体"的是**字段**不是字符＝`Attachments[].Name`
（`internal/panel/attachments.go:199 Ingest` 生产读者 **0** ⇄ 正控同尺形同射程 `json.Marshal` **37** 行；
`cmd/wisp/panel_inbound.go:301-302` `Attachment: nil`／`Message: nil`；`internal/panel/pump.go:275` 附件参传 `nil`）。

## 四、AC#3 两枚名单外字符的**真盘上正控**（⛔ "按代码读应该会出现"）

唯一执行过的代码＝`D:\tmp\332a1\carrier\main.go` 与 `D:\tmp\332a1\diskroster\main.go`（`go run` 各 rc=0、stderr 空）。
逐字读数（`logs/ac3-carrier-disk-and-marshal-dump.log`，26 行）：
```
DISK cp=U+2061 name="photo\u2061report.txt" created=true readdir_hit=1 raw_utf8_hex=70686F746FE281A17265706F72742E747874
DISK cp=U+206B name="photo\u206breport.txt" created=true readdir_hit=1 raw_utf8_hex=70686F746FE281AB7265706F72742E747874
BOUNDARY cp=U+2061 needsDisplayWash=false washed_changed=false
  MARSHAL_WASHED hex=2270686F746FE281A17265706F72742E74787422  SHA256 da1e771280e63fd92d26d37464bfc6cfd6611f8f8e7dab4dda7a6767b3af889b
BOUNDARY cp=U+206B needsDisplayWash=false washed_changed=false
  MARSHAL_WASHED hex=2270686F746FE281AB7265706F72742E74787422  SHA256 bb7a8ae1915fe249edc6b8b6165788e4b7140337f38f038966f90a352ffdd4db
SCREEN total_outside=155 scrub_leaves_raw=155 unicode_is_space=0 taint_ignores=0 unicode_is_control=0 json_escapes=0
```
⇒ `E2 81 A1`／`E2 81 AB` **在真盘名与 marshal 字节里都逐字节在位**，且 marshal 字节里⛔ 转义形（`needsDisplayWash` 返回 false ⇒ 快路返回同一枚串）。
**成对正控（本腿⛔ 产码 ⇒ 无"突变两发"，代之以同探针内的成对反向读数）**：
名单内 `U+202E` ⇒ `needsDisplayWash=true`、marshal 字节 `…5C5C7532303265…`＝可见转义在位、**`E2 80 AE` 0 命中**；
`Cc` 的 `ESC` ⇒ 洗成 `\x1b`，且**盘上写不出来**（`The filename, directory name, or volume label syntax is incorrect`）。
★把盘上尺扩到全部 155 枚（量的，`logs/ac2a-disk-legality-per-cf-codepoint.log` 170 行）：
155 枚逐枚写真名读回＝**155/155 `READBACK_OK`、`fail=0`**；8 枚对照里 `U+0000/0009/000A/000D/001B/001F`（`C0`）全 `WRITE_FAIL`、
`U+0085/U+009B`（`C1`）与 `U+00AD/U+200B/U+202E/U+FEFF`（名单内 `Cf`）全 `READBACK_OK`。
第三来源逐枚对合（`os.listdir` 尺）＝`P_COUNT=155`、`MISSING_ON_DISK=[]`、`EXTRA_ON_DISK=[]`、`CF_UNIQUE=True`；
`DISK_REMAINING=164` 的 9 枚差额本腿逐枚交代＝5 枚成功的对照行＋该目录里探针自己的 4 枚文件（`main.go`/`go.mod`/两份输出）。

## 五、派单必答题（两枚，⛔ 本腿盖章不算）

**必答题一＝判定层 `isIgnorableTaintRune` 的枚数与它为什么取向相反。**
枚数（量的）＝**6 枚显名**（`internal/risk/taintmatch.go:41`：`0x00AD, 0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF`）
**＋整个 `unicode.Mn` 族＝2059 枚**（`taintmatch.go:44`；尺＝`D:\tmp\332a1\mncount\main.go` `go run` rc=0：`MN_TOTAL=2059 MN_INTERSECT_CF=0`）
⇒ 表射程 `6+2059=2065` 枚，但**落在 `Cf` 轴上的只有那 6 枚**（`MN_INTERSECT_CF=0`，另有件2 §③那把筛检尺在 155 枚上逐枚复算＝`taint_ignores=0`）。
两表交集＝那 6 枚 ⇒ 显示名单 15＝判定 6＋方向族 9（`U+202A–U+202E`、`U+2066–U+2069` 不在判定表里）。
**为什么相反**（照它自己的注释逐字，`taintmatch.go:32-38`，出处＝对抗报告 `16.9#1`/`M-6`）：
"leaving them in place lets an exfiltrator cut a tainted fragment into sub-8-rune pieces with zero-width characters"
——判定轴要做的是把那批字符**从归一化串里删掉**（`normalizeTaint:59` 用它做删除），好让 C25 的 ≥8 字符片段匹配不被零宽切口打掉；
显示轴要做的是把看不见字符**显出来**（转义、⛔ 删，`display_scrub_216.go:71-79` 自陈"deletion splices"）。一删一显 ⇒ 同一枚码位在两轴上的正确动作相反，
故两份名单各写各的、分处两文件（`display_scrub_216.go:142-147` 亦自陈 "NOT internal/risk's isIgnorableTaintRune"）。★本腿对它一字未动。

**必答题二＝票 327 普查件（`census-1-roster.md`／`census-2-cf-and-upstream.md`）今天已腐烂的那一环。**
最硬的一环＝它那两句"最后一公里"读数：它写"快照字节今天**没有进 WebView 的生产路径**／读它的 `lastPanelSnapshot` 生产调用者 **0 枚**"
⇒ 现量生产读者 **1 枚**＝`cmd/wisp/panel_snapshot_deliver_329.go:95`，链条 `cmd/wisp/run.go:744 Out: rt.bookPanelSnapshot`
→ `cmd/wisp/panel_pump.go:447/463` → `panel_snapshot_deliver_329.go:105 h.evalPanelJS(...)` → `:66` 逐字包裹
`window.wispOnPanelSnapshot(<marshal 字节原文>)` → `cmd/wisp/panel_host_windows.go:331/:336 w.Dispatch(func(){ w.Eval(js) })`；
`internal/panel/pump.go:364-367` 那句"最后一公里在票 33/35 名下"今天只服务 `src.Out == nil` 那一形。
另三条同批更正：① 行号 `display_scrub_216.go:169/175` ⇒ 现量 `:228/:234`；② "`RequestWorkspaceSwitch` 非测试调用者 0 枚"
⇒ 现量 **1 枚**＝`internal/panel/workspace_handlers.go:125`，且常驻链从页面门可达（`panel_host_windows.go:948` → `composer_dispatch.go:209`）
⇒ 页面敲进去的 `Spelling` 一字不洗直接出向这一形**今天才真的通**；③ 它那族候选 22 枚、当时 `Cf` 洗 0 枚 ⇒ 今天全集 170、已洗 15、待答 155。
**未腐烂**（本腿复跑同读数）：`Attachments[].Name` 源与读者双 0；`results`/`pending`/`git.*`/`instructions.files[].path` 不过显示边界。

## 六、纪律自报（⛔ 漏一项就按派单退回）

- **本腿执行过的代码＝只有四枚仓外单文件纯 stdlib 探针**（逐字源码已入库 `logs/probe-sources.log`）：
  `D:\tmp\332a1\cfprobe\main.go`、`D:\tmp\332a1\carrier\main.go`、`D:\tmp\332a1\diskroster\main.go`、`D:\tmp\332a1\mncount\main.go`；
  `go run` 各 rc=0、stderr 逐枚为空；⛔ import 本仓任何包。四段显示边界函数＋一段判定层函数＝锚上原文的**逐字重打**（出处行号已在件2 与本件具名）。
  另自报跑过 `go version`（`go1.27.1 windows/amd64`）与四枚 `go mod init`（只写 `D:\tmp` 里探针自己的目录）；⛔ `go env`／`go list`。
- ⛔ `go test`/`go build`/`go vet` 一枚未跑（编译面与重包归 `327-r2`/`228-r3`/`216-r2`）；⛔ 门禁四数一枚未跑（派单："你这波⛔ gofumpt／⛔ d22scan"）
  ⇒ 因此**没有** `rc=` 可报，也没有改前改后四数（本腿零字节产码，颜色面不适用）；⛔ 突变两发（同上，代之以§四那组成对反向读数）。
- ⛔ 读工作树在飞字节：全部读数取自 `git show 7535d83b:<路径>`／`git grep … 7535d83b`。
  ★由此带来的已知射程限制如实交：本腿的逐字重打只对锚 `7535d83b` 有效；若 `327-r2` 正在改 `internal/panel` 的字节，
  其结果对本件的影响＝名单本身可能变长（那正是编排者与裁决腿要裁的那一纸），⛔ 本腿替它担保。
- ⛔ 翻 AC 框／⛔ 改票名／`Status:`／⛔ `-done` 后缀；票 332 只追加一条 bullet（`git show --name-only` 可核＝只那一枚文件、纯增行）。
- ⛔ push；`git add` 与 `git commit` 两处逐枚显式 pathspec（⛔ `add -A`／⛔ `-- .`／⛔ amend／⛔ reset／⛔ stash／⛔ checkout ．／⛔ clean）。
- 活实例 `wisp.exe` **PID 21892 完好未动**（交件时点现量＝`tasklist /FI "PID eq 21892"` 命中该行，Mem 65,376 K）。
- ⛔ 结论／建议／"建议纳入名单"那一类判语：本件通篇只有读数与尺；155 枚的三形载体齐平这一事实**不**被本腿写成任何"该扩／不该扩"。
