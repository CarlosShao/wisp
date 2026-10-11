# 306-v1 · 20 放宽断言尺／AC#3 行号核／门禁与名册（必答 2·5·6）

★本件所有"HEAD"＝**`17dd4a6f…`**（本腿写这件时第三枚重锚，见 `90-final.md` §1）。
本程⛔ 在共享工作树做任何写动作（⛔ `git checkout`、⛔ `stash`、⛔ 改文件再还原、⛔ push）。

## 1. ★必答 2：有没有顺手放宽断言（三把尺，本腿现跑，⛔ 腿自陈⛔ 转述）

### 1.1 尺 A＝`git diff 5480434f HEAD -- internal/audio`（`5480434f`＝派单前一枚编排者笔，＝实现腿的父锚区间左端）

件＝`logs/A0-rerun-at-head.txt`（对当前 HEAD 重跑）＋`logs/10-diffscope.txt`（同尺在 `HEAD=e4740e35` 时那一发，读数逐字相同）：

```
M  internal/audio/wave_format_float_300_windows_test.go     16  5
A  internal/audio/wavinjector_extensible_float_306_test.go  225 0
rc_diffquiet[internal/audio/wavinjector.go]=0        ← 产码 0 字节变化
rc_diffquiet[internal/audio/wavinjector_test.go]=0   ← 载 10 枚 writeWav 调用点那枚文件一字未动
```

⇒ **产码面（`wavinjector.go`）＝ 0 字节**；既有测试文件（`wavinjector_test.go`）＝未动；
`internal/audio` 在整个区间里⛔ 别家笔碰过（区间里夹着 `5f891575`／`53c854e3`／`52a5eb23`（`305-a1`）与 `c888420e`（编排者收腿笔），
名册只有上面那两行 ⇒ ⛔ "别人的改动被我算成放宽"这一形）。

### 1.2 尺 B＝逐笔 `git show --name-only`（8 笔，⛔ 区间尺）

件＝`logs/10-diffscope.txt`（笔笔全文件名）＋`logs/60-legclaims.txt` §V3（逐笔计数）。写面外枚数：

| 笔 | 该笔文件数 | 写面外 |
|---|---|---|
| `4e73d0a5` | 7 | **0** |
| `1b34e475` | 5 | **0** |
| `c0d5ec24` | 6 | **0** |
| `fbd150e3` | 35 | **0** |
| `71249b05` | 21 | **0** |
| `6284a489` | 3 | **0** |
| `53aa4b16` | 6 | **0** |
| `e4740e35` | 3 | **0** |

并集尺（显式列 8 枚 SHA，⛔ 任何 `/tmp` 宽 glob——那正是腿 `53aa4b16` 自己踩过的雷）＝
**80 枚**：`internal/audio/**` **2** ＋ `.scratch/wisp/probes/306/r1/**` **77** ＋ 票 306 那枚文件 **1** ＋ 写面外 **0**
（件＝`logs/V2-union.txt`，行数 80；三列各自 `grep -c`＝2／77／1，2+77+1=80 ⇒ 名册行数与"排除枚数"对得上）。

⚠ 与腿件不符（具名，见 `90-final.md` §5 第 4 条）：§7 那节写"六笔并集 75 枚＝2＋72＋1"，本腿八笔并集＝**80＝2＋77＋1**；
差的 5 枚全在 `probes/306/r1/**`（第 7、8 笔自己的更正件）。⛔ 矛盾（腿那句写的是六笔口径），只⛔ 它把枚数写死过。
⚠ 另：盘上 `r1` 目录 **78** 枚文件 vs tracked **77** 枚 ⇒ `logs/90-commit8.txt`（644 字节，那枚记录第 8 笔的件）**未提交、留在盘上**
（尺＝`find` 计数 vs `git ls-files` 计数，件＝`logs/70-rootaudit.txt` §W3/§W4）。本腿⛔ 动它，只报。

### 1.3 尺 C＝"只在 `//` 行"那一格，两把并排

- **腿那把（逐字重跑）**＝`git show c0d5ec24 -- internal/audio/wave_format_float_300_windows_test.go | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)|^[+-]//'`
  ⇒ 本腿现量 **零命中**（`rc_leg_ruler_c0d5ec24=1`，件＝`logs/20-ac3-diff.txt`）；同尺换成区间口径（`git diff 5480434f HEAD -- <那枚文件>`）⇒ 同样零命中（`rc_leg_ruler=1`）。
- **本腿另加一把更严的**＝同时豁免缩进的 `//` 行（`grep -vE '^[+-][[:space:]]*//'`）⇒ **仍零命中**（`rc_strict_ruler=1`）。
  ⇒ ⛔ 一枚改动行既⛔ 是行首 `//` 也⛔ 是缩进 `//`；`numstat` 16 加／5 删全在注释里。
- 期望值／断言侧的旁证：票 300 那两枚钉的形状本腿自己跑（⛔ 抄腿）＝
  `logs/Y3-ticket300-pins.txt`：`TestWaveFormatConstantsMatchMmregAuthority300` 顶层 1＋子测试 **6**、
  `TestParseWaveFormatSubFormatOffset300` 顶层 1＋子测试 **4**，`rc=0`、逐枚 `--- PASS` ⇒ ⛔ 放宽、⛔ 变绿成 `--- SKIP`。
- `waveFormatPCM` ⛔ 删（票面硬边界）：尺＝`grep -c waveFormatPCM <file>` ⇒ `wasapi_windows.go`=**1**、
  `wave_format_float_300_windows_test.go`=**7**（件＝`logs/90-misc2.txt` §Y1，与腿所报逐字同）。

### 1.4 禁改清单逐枚（本腿那把尺＝7 枚逐枚 `git diff --quiet 5480434f HEAD -- <f>`，件＝`logs/A0-rerun-at-head.txt` §R2）

`wavinjector.go`／`wasapi_windows.go`／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／`docs/PLAN.md`／
`internal/audio/parse_wave_format_300_windows_test.go` ⇒ **各 `rc=0`（0 字节）**；
`internal/audio/wave_format_float_300_windows_test.go` ⇒ `rc=1`＝**它确实变了**，变的只有 §1.3 那 16/5 枚注释行（＝`AC#3` 授权那一格）。
`git diff --name-only 5480434f HEAD -- frontend design third_party` 行数＝**0**（§R3）。

## 2. ★必答 5：`AC#3` 只动注释那一格＋三枚 SDK 头文件的行号（本腿逐枚去自己机器上核）

本腿机器上那份 SDK＝`C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/`（⛔ 第二枚树，尺＝`logs/20-sdk-find.txt`，
该 `Include` 目录下只装着 `10.0.26100.0` 一枚）。逐枚现量（件＝`logs/20-sdk-lines.txt`，都是 `awk NR` 取的原始行）：

| 引用（注释里写的号） | 盘上那一行的原文 | 判定 |
|---|---|---|
| `mmreg.h:2418` `WAVE_FORMAT_PCM` | `#define WAVE_FORMAT_PCM         1`（`:2418`） | **对** |
| `mmreg.h:2110` `WAVE_FORMAT_IEEE_FLOAT` | `#define  WAVE_FORMAT_IEEE_FLOAT  0x0003`（`:2110`） | **对** |
| `mmreg.h:2376` `WAVE_FORMAT_EXTENSIBLE` | `#define  WAVE_FORMAT_EXTENSIBLE   0xFFFE`（`:2376`） | **对** |
| `mmreg.h:2480-2484` 那份 GUID | `:2480 #if !defined( STATIC_… )` … `:2483 DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);` `:2484 #define … DEFINE_GUIDNAMED(…)` | **对** |
| 那句"GUID 的 `DEFINE_GUIDSTRUCT` 在 `:2483`" | `grep -n` 逐字命中 `mmreg.h:2483` ＋ `ksmedia.h:854` 两处 | **对** |
| `ksmedia.h:850-855`＋"由 `defined(_INC_MMREG)` 包着"＋`:854` | `:850 #if defined(_INC_MMREG)`；`:851 #if !defined( STATIC_KSDATAFORMAT_SUBTYPE_IEEE_FLOAT )`；`:854 DEFINE_GUIDSTRUCT("00000003-…", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);`；`:855 #define … DEFINE_GUIDNAMED(…)`；`:856 #endif` | **对**（一枚措辞级小账见下） |

**票面／派单写的 `ksmedia.h:851-855` 与盘上差在哪：** 盘上那段的外层守卫是 **`:850`** 的 `#if defined(_INC_MMREG)`，
`:851` 是里层的 `#if !defined( STATIC_… )`；`DEFINE_GUIDSTRUCT` 那行在 **`:854`**。
⇒ 票面那个号段的**存在性判断⛔ 错**（`:851-855` 里确实装着那枚 GUID 定义），**守卫起点写晚了一行**；
腿按盘上写 `:850-855`＝**对的**。⚠ 本腿给这一段留一枚小账（⛔ 缺陷、⛔ 要改）：`:850-855` 那一段止于 `:855`，
外层 `#if defined(_INC_MMREG)` 的收口 `#endif` 在 **`:856`** ⇒ 注释那句"wrapped in a defined(_INC_MMREG) guard"
如果按号段读，区间⛔ 含守卫的收口。这属于号段粒度的选择（腿写的是"那份替代定义所在的行段"，`854`/`855` 都在段内），
**本腿判 AC#3 时⛔ 因此扣任何分**，只把它记成号段口径的一枚注脚。

兄弟件那一枚（票面 §现量 那句"其兄弟件 `:101` 引 `mmreg.h:2483` 是对的"）：本腿现量
`internal/audio/parse_wave_format_300_windows_test.go:101` 逐字＝`why: "KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24",`
⇒ **票面那一句对**（`:2483` 与本机 SDK 逐字一致，尺＝`logs/30-claims.txt`；⚠ 枚数口径：本腿那把尺是 `grep -n "2483"` 在该文件的**唯一**命中行）。

`AC#3` 那一格改后的正文（16 枚加行逐字在 `logs/20-ac3-diff.txt`）写足了三件事：两份定义各在何处、
那枚钉引用的出处＝**格式标签** `mmreg.h:2110`⛔ 是 GUID、以及"两者共享 3 是因为 GUID 由 `DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_IEEE_FLOAT)` 从标签拼出"。
本腿另核了一句：注释里"⛔ 读头文件"那句是**自我限制**（该文件⛔ 任何运行时读 SDK 的调用，`grep -n` 命中仅注释与常量声明）——
它把 SDK 号段写进注释而⛔ 依赖它，符合票面对 `AC#3` 的"最小面"要求。

## 3. ★必答 6：门禁三件套＋两把格式名册＋红名册（含交集尺盲区那问）

| 尺 | 台面 | rc | 读数 |
|---|---|---|---|
| `GOFLAGS= go vet ./internal/audio/` | **共享工作树**（本腿车道内） | **0** | 输出 0 行（`logs/g1-vet-worktree.txt`）⚠ 台面 caveat 见下 |
| `sh scripts/d22scan.sh` | 共享工作树（只读扫描＋它自己的 positive control 包） | **0** | 逐字 `clean - no D22 ban violations`；`runtests.sh: OK … top-level: PASS=35 FAIL=0 SKIP=0`（件＝`logs/g2-d22scan-worktree.txt` 23084 字节，⛔ 0 字节件） |
| `GOFLAGS= go build ./...`（票面 AC#4 那一枚） | **仓外导出树**（`git archive e4740e35`） | **0** | 0 行输出（`logs/g3-build-all-exporttree.txt`＝0 字节＝**读数**，`rc` 已单独落在 `ZZ-gates.txt`）⚠ 本腿⛔ 在共享工作树跑整包＝编排者车道 ⇒ 本格判语**附条件**，欠的那一发写在 `90-final.md` §6 |
| `GOFLAGS= go vet ./internal/audio/` | 导出树 | **0** | `logs/g3b-vet-exporttree.txt` |
| `GOOS=linux GOFLAGS= go vet ./internal/audio/`＋`go build` | 导出树 | **0／0** | 平台中立那一面仍绿（`logs/90-misc2.txt` §Y4）＝新文件确实⛔ 带 `//go:build`（另尺：`sed -n '1,2p'` 第一行逐字 `package audio` ⇒ 位置⛔ 放构建约束；`grep -c 'go:build'`＝1，那一枚命中在**注释正文**里，§Y2） |

**ban #8 射程（本腿自己现量，HEAD＝`17dd4a6f`）：**`internal/`=**527** Go files（含注释与 `_test.go`）、`cmd/`=**119**，
逐字两行在 `logs/g2-d22scan-worktree.txt:249-250`，`ZZ-gates.txt` 也贴了一份，另有仪器自己的自证行
`scan_test.go:1261: verified ban #8 internal/: 527 files`／`verified ban #8 cmd/: 119 files`（＝⛔ 我数出来一个数就报一个数，仪器自己也在钉这个数）。
⇒ 与派单快照 526 差 1 枚＝本票那枚新测试文件进册；腿那句"它⛔ 豁免于注释、全 ASCII"与本腿读数一致（ban #8 `clean`）。

**格式名册（两把并排，各自射程）：**

| 尺 | 射程 | 未格式枚数 |
|---|---|---|
| `gofmt -l internal/audio`（工作树，递归） | 工作树 `internal/audio` | **0**（`logs/g4-gofmt-worktree-internalaudio.txt`＝0 字节，`rc=0`） |
| `git archive <sha> internal/audio \| tar -x` 后 `gofmt -l <目录>` | **HEAD blob**（`e4740e35`，22 枚 `.go`） | **0**（`logs/g5-gofmt-blob-e4740e35.txt`） |
| 同上，基准锚 `5480434f`（21 枚 `.go`） | **blob** | **0**（`logs/g5-gofmt-blob-5480434f.txt`） |

⇒ **新增未格式 0 枚**，且分母对得上腿的读数（21→22 枚＝那枚新文件）。
⚠ 全仓那把尺（`gofmt -l .`，射程＝工作树全仓）**本腿量到 `rc=2`、名册 44 行**，其中 `internal/audio/` 命中 **0** 行
（尺＝`grep -cE '^internal/audio/'`，件＝`logs/g4b-gofmt-worktree-fullrepo.txt`）。`rc=2`＝有些文件 gofmt 根本读不动／⛔ 是 Go；
那 44 行的形状（反斜杠档名、`.scratch/wisp/probes/**`、`internal\agent\tools.go` 等）＝票 303 收腿笔里已经具名过的 CRLF 假枚与别家夹具，
**⛔ 本票射程、⛔ 归因给本腿、也⛔ 由本腿算成"新增"** ⇒ 我只写"自家射程 0 枚"，⛔ 写"全仓 0 枚"。

**红名册（★交集尺的结构性盲区那一问，本腿⛔ 用交集尺）：**

本腿的尺＝**单发整包名册逐枚并排**：`go test ./internal/audio/ -count=1 -v`（导出树，HEAD 台面）共 **5 发**
＝`logs/sens/B1-fullpkg-shot1.txt`／`B2-fullpkg-shot2.txt`／`sens/B3…B5-fullpkg-shot{3,4,5}.txt`：

```
每发：rc=0  ^--- FAIL=0  ^--- PASS=44  ^--- SKIP=1  ^    --- PASS=12
名册（^--- FAIL|^--- SKIP，剥时长，排序）＝ 1 行：--- SKIP: TestLiveWasapiSmoke
逐对 diff：B1 vs B2 rc=0 ／ B1 vs B3 rc=0 ／ B1 vs B4 rc=0 ／（五发名册逐字同一枚行）
```

⇒ **新增红 0 枚**这一句在本腿这里⛔ 是"两把尺取交集"得到的，而是"5 发单发名册两两 `diff` 全空"；
被"只红过一次"藏起来那一形在本腿这 5 发里**没有出现**（每一发都恰好 0 枚 `--- FAIL`，⛔ 任何一枚时有时无）。
唯一名册内那枚 SKIP 的**世代**本腿钉了＝`git grep -c "func TestLiveWasapiSmoke" 5480434f -- internal/audio`
⇒ 命中在 `5480434f:internal/audio/hotplug_test.go`＝**改前就有、⛔ 本腿所造、⛔ 别家新造**（件＝`logs/ZZ-controls.txt` §H2）。

定向（改⇒红）那一侧的成对名册，本腿也各取了一发整包 `-v`：
`D6`（`:196` 上下界同移，夹具在）＝`rc=1`／FAIL 恰好 2 枚（两枚指名用例）／顶层 PASS 仍 **42**；
`D7`（`:218` `3`→`2`）＝同形；`B4`（`:192`）＝同形 ⇒ **⛔ 连带红**（其余 42 枚逐字绿，`red roster` 只有那 2 枚）。
"无夹具"那一侧 5 发（`D1`-`D5`）＝`rc=0`／PASS 42／FAIL 0／`panic` 0 ⇒ 三枚权威值（含 `:196` 两形）
在本票夹具之外**确实⛔ 任何执行者**＝票面 §现量 那两发 `rc=0` 的**同台面复跑＋补量第三枚**。

⚠ 台账级 caveat（⛔ 影响判语、⛔ 归本腿）：本腿在共享工作树跑的那两把尺（`go vet`／`d22scan`）落在
**`go.mod` 脏面**上（别家改动：`github.com/jchv/go-webview2` 从 indirect 挪进 direct，尺＝`git diff --stat -- go.mod`＝1/1，件＝`logs/40-misc.txt` §M3）。
`internal/audio` ⛔ 依赖它；且 `git diff --quiet 5480434f HEAD -- go.mod go.sum` ＝ **`rc=0`** ⇒ 实现腿的读数与本腿的读数用的都是**已提交**的 `go.mod`，
导出树那几把（整包 build／vet／5 发红名册）用的是 `e4740e35` 的 `go.mod`。`.gitignore` 也有别家 +9 行脏面（§M4），本腿⛔ 碰。

## 4. 交件套件尺（票面 `AC#4` 那句"0 字节件⛔ 算交"，本腿自查＋量腿）

- 实现腿：`find .scratch/wisp/probes/306/r1 -type f -size 0`＝**0 枚**、`-name '*.out'`＝**0 枚**；
  `.md`＝**6** 枚、`logs/**`＝**72** 枚（件＝`logs/60-legclaims.txt` §V1）
  ⇒ 腿 §1 那句"6＋66"是它自己第 7、8 笔之前的世代读数（它 §7 已经自己把枚数改成"⛔ 写死、尺在句子里"，本腿认同那句）。
- 本腿：所有 `rc=` 走纯重定向 `> 件 2>&1; rc=$?`（驱动脚本 `logs/*-driver.sh` 全文可查）；
  本件里**具名自报**一枚自己踩到的同类雷＝驱动 2 里 `diff a b | sed … >> SUM` 之后那句 `rc=$?` 取的是 `sed` 的退码
  （⇒ 那一行 `rc_B1_B2_roster_diff=0` ⛔ 是 diff 的退码，那把尺作废）；本腿另起一把干净的
  （`diff … > 件 2>&1; rc=$?`，读数＝`rc_diff_B1_B2=0`，件＝`logs/ZZ-d9redo.txt`），⛔ 删坏件、就地保留 `logs/D9redo-roster-diff.txt` 与原始错件 `logs/D9-b1-b2-roster-diff.txt` 两枚并存。
  ⚠ 同一次事故还暴露一枚真读数：第一次那把尺找错目录（roster 件在 `logs/sens/`⛔ 在 `logs/`）⇒ 那份 `D9` 件里那 240 字节是 `diff` 的报错，⛔ 名册差。
