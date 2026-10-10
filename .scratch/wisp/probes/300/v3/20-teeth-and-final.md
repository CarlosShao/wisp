# 300-v3 — `20` 恒真攻击／生产调用点／`AC#5` 那句话的有效期／终态锚（交回件）

原始 stdout 全在 `.scratch/wisp/probes/300/v3/logs/`（`rig-run.txt`／`rig-run-2.txt`／`rig-run-3.txt`／`final-anchor.txt`／`gates-*.txt`）＋仓外 `/tmp/wisp300-v3/logs/*.txt`（11 枚，只建不删，清理归编排者）。
**本件⛔ 重打任何"我在 stdout 里看过"的数**（`A810` 那一族：结论句与按节表必须逐枚对拉）⇒ 每行都写清是哪把尺、读数在哪个文件。

## 1. 必答②：攻恒真——`AC#1`／`AC#2` 的判据把**判据本体换成反形**之后还红不红

台件＝仓外导出树（`git archive dce0f133` → `/tmp/wisp300-v3/tree`；起手两枚文件 `cmp` ↔ `git show dce0f133:<path>` 皆 `IDENTICAL`；导出树里 sherpa DLL 前置尺＝**0 枚**，`A802` 那条，本包按 `A805` harness-inert、仍数了一发具名）。
每条突变都带三件硬要求：① 输出先落文件；② 跑完 `cp` 写回 + `cmp` 证还原；③ 自带正控（pristine 该绿＝`M0`、明知该红的要红＝`M1`）。**每发先过 `LANDS` 门（内容锚命中数＝期望数）才许取读数。**

| tag | 反形打在哪一侧 | rc | 尺名：顶层`--- PASS`／子测试`    --- PASS`／`--- FAIL` | 还原 |
|---|---|---|---|---|
| `M0-pristine` | ⛔（正控） | 0 | 1／4／0 | — |
| `M1-prod-offset-26` | 产码偏移 24→26 | 1 | 0／0／5 | `restore_cmp=IDENTICAL` |
| `M2b-fixture-S1-at-22` | **夹具侧**：只把 S1 的 `Data1` 写到 22，产码⛔ 动 | 1 | 0／3／2 | `restore_cmp=IDENTICAL` |
| `M3b-expect-S1-want-5` | **期望侧**：只把 S1 的 `wantTag` 写成 5，产码⛔ 动⛔ 夹具 | 1 | 0／3／2 | `restore_cmp=IDENTICAL` |
| `N1b-prod-constant-tag` | **产码语义侧**：把 `uint16(sub.Data1)` 换成恒等常量 3（那读仍在） | 1 | 0／1／4 | `restore_cmp=IDENTICAL` |
| `M6`（第一发） | 同上，但我删了 `sub` 的使用 | 1 | 0／0／0 | ⚠ **装置红**：`internal\audio\wasapi_windows.go:211:3: declared and not used: sub` ⇒ **该发⛔ 是判据红、读数作废**（`logs/M6-prod-constant-tag.txt`），改写成 `N1b` 才有效 |
| `N2-floatconst-wholepkg` | `waveFormatFloat` 3→1，跑**整包** | **0** | **41／6／0** | `restore_cmp=IDENTICAL` |
| `N0`／`N3` | 终态复跑 | 0 | 41／6／0（整包）／1／4／0（定向） | 两枚文件 `final … = IDENTICAL` |

**红句逐字（换形之后确实红了，各一句）**
- `M2b`（夹具侧反形）：`parse_wave_format_300_windows_test.go:138: tag = 0, want 3 (KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24)` ＋ `:155: floating = false, want true (convertPacket would take the float32 path)`
- `M3b`（期望侧反形）：`:138: tag = 3, want 5 (…)`（正确解析撞上写错的期望＝红 ⇒ 断言⛔ 在自证）
- `N1b`（语义反形）：`:138: tag = 3, want 7 (positive control: … the parser has to report the value it actually read, so this face may not go green on a constant)` ⇒ **S3 头顶那句 `why` 我这把当场兑现**：解析器一旦退化成常量，四形里三形红（`S2`/`S3`/`S4`），只有 `S1` 绿（它本来就期望 3）。
- 终态残尺（逐字，`logs/rig-run-3.txt`）：`prod24=1 prod26=0 tagassign=1 floatconst3=1 S1at24=3 want3=1` ⇒ **⛔ 一枚突变留在被跟踪文件里，也⛔ 一枚留在导出树里**。

**判**：`AC#1`／`AC#2` 那枚入库判据⛔ 是恒真那一形——**三侧换形都红**（夹具侧／期望侧／产码语义侧），加上 `M1` 与历史三把（`a1r`/`v1`/`orch`）同值。恒真那三支里最要命的一形（"两种形状都放行"）在我这把⛔ 出现。

**但量到一枚真的、⛔ 在任何一枚腿件里出现过的缺口（两枚一体）**：
1. **代数**：那枚 `floating` 断言在测试里是 `gotFloat := got.tag == waveFormatFloat` / `wantFloat := tc.wantTag == waveFormatFloat` ⇒ 它⛔ 能独立于 `tag` 断言而红（`tag` 一等 ⇒ 两边同 ⇒ 恒绿）。＝一枚**派生断言**，形状上是"只在坏世界里跟着红"的哨兵，⛔ 是第二枚牙。
2. **实测**：`N2` 把 `waveFormatFloat` 从 **3** 改成 **1**（＝"哪个 tag 算浮点"这一枚语义常量整枚换掉），**整包 41 枚顶层用例全绿、`rc=0`** ⇒ `internal/audio` 里⛔ 有任一把尺钉住那枚常量的值。票面 `:196` 那句 `convertPacket branches on floating alone` 所依赖的正是它。
   ⇒ 归口＝**仪器票（票 301/302 那一族）射程**，⛔ 本票另开一格普查（派单第 4 条写死）；也⛔ 由我这把改产码或改断言——那是实现者的面。

## 2. 必答③：生产侧调用点枚数（现跑，内容锚 `grep -n`，⛔ 信行号）

尺＝`git grep -n "parseWaveFormat" -- internal cmd`／`… "GetMixFormat" …`／`… "convertPacket" …`／`… "wasapiOpener" …`（四把同发，stdout 见本件；命中逐枚分类见下）

- `parseWaveFormat`：名册 9 枚命中 ⇒ 1 枚函数声明 ＋ 2 枚注释（头顶 doc 与错误串）＋ **1 枚生产调用**（内容锚 `f, ferr := parseWaveFormat(wfex)`，在 `(*wasapiStream).init` 体内）＋ 5 枚在 `parse_wave_format_300_windows_test.go`（3 枚注释、1 枚调用、1 枚错误串）。
- `GetMixFormat`：6 枚命中 ⇒ **1 枚真实 COM 调用**（内容锚 `comCall(s.client, slotIAudioClientGetMixFormat, uintptr(unsafe.Pointer(&wfex)))`）＋ 1 枚常量声明 `slotIAudioClientGetMixFormat = 8` ＋ 2 枚注释 ＋ 2 枚错误文案串。
- ⇒ **那条边（`GetMixFormat` → `parseWaveFormat`）生产侧调用点＝1 枚**，⛔ 零。
- 上游可达性我也现跑了两把：`git grep -n "wasapiOpener" -- internal cmd` ⇒ 生产构造点 1 枚（内容锚 `return newWASAPIMicrophoneWith(newMMDeviceWatcher(), wasapiOpener{}, opts...)`）；`git grep -n "NewWASAPIMicrophone" -- cmd internal` ⇒ 生产侧 1 枚（`cmd/wisp/resident_audio_windows.go` 里逐字 `return audio.NewWASAPIMicrophone(opts...)`）。
- 下游同样⛔ 断头：`convertPacket` 生产调用 1 枚（内容锚 `out = append(out, convertPacket(data, int(frames), ch, floating, …)`），其外层 `Drain()` 的生产调用 1 枚（内容锚 `samples, err := stream.Drain()`）。
- **判**：本格⛔ 是"能力类判据只在测试里成立而生产零调用者"那一形——**边在产码里活着，四段（`NewWASAPIMicrophone`→`Open`→`init`/`GetMixFormat`→`parseWaveFormat`→`floating`→`convertPacket`）逐段有调用者**。
- ⚠ 一句别读大：**"生产有调用者"⛔ 等于"那一支今天真走到"**；后者靠的是 §1②（票 300 `AC#3`）那枚真机 dump（本机四枚端点全 extensible＋float），⛔ 任何静态尺能替代。

## 3. 必答④：windows-tagged × `core` 认领 ⇒ 票面上哪句话要说成⛔ 持续防护

★**先顶回派单给我的那枚前提——它在盘上已经过期**（具名，⛔ 凑数）：
派单原文：「`./internal/audio/` 已被 `core`(ubuntu) 档认领 ⇒ 那枚新用例在 CI 上**从未被执行过一次**」。
- 尺＝`grep -n "internal/audio" scripts/portable-tests.sh` ⇒ 现量三处认领：`core_pin`／`win_pin` 各 1 枚导入路径，**外加 `windows)` 清单里逐字那枚 `./internal/audio/`**；尺＝`git log --oneline -1 -- scripts/portable-tests.sh` ⇒ 来路＝`0a0f62ef`（14:34:30，标题逐字 `ticket 301 AC#1: claim internal/audio in the windows) tier so its //go:build windows tests can be evaluated by a CI job at all`）；尺＝`git merge-base --is-ancestor 0a0f62ef dce0f133` ⇒ **YES-ancestor**（已在我起手锚里）。
- ⇒ "**从未被执行过一次**"只对**那一笔之前的历史**成立；盘上现状＝windows 档已认领该包。**CI 究竟有没有真跑过一发、色是什么＝⛔ 我这把能裁**（票 301 `AC#2`/`AC#4b` 的射程，且我刚量到 `303-a1` 在飞、⛔ 我取整包级 CI 读数）⇒ 记〔欠，归票 301，等编队空〕。

★**票面因此要说成⛔ 持续防护的句子＝两处"一进门就绿"**（尺＝`grep -n '一进门就绿' .scratch/wisp/issues/300-*.md`）：
- 「编排者补格」里那句 `⇒ 真正的跟踪测试随 **AC#2**（修完）同批入库，那时它应当**一进门就绿**`；
- 「收 `300-v2`」里判语那行 `★**判语＝AC#2 成立**（四件事各有尺：…入库用例一进门就绿…）`。
这两处的"门"在盘上＝**本机那枚窄档 `-run`**（我这把 `N3` 顶上同一事实的另一代：`rc=0`、顶层 1／子测试 4／FAIL 0），⛔ 一枚 CI 档的求值 ⇒ "一进门"三字⛔ 许被读成"CI 会替我们守着"。
- ★票面**已经**自己写过这一句（尺＝`grep -n 'CI 兜着' .scratch/wisp/issues/300-*.md` ⇒ 现量在「编排者现量 12:5x」那节：`⛔ 任何人以为"用例落进 internal/audio 就有 CI 兜着"是读反了`）⇒ **⛔ 需要新加一句"⛔ 持续防护"**；该写的是**那一句的有效期**：它当时的射程靠"windows 档⛔ 认领 audio"这一枚前提，而那枚前提被 `0a0f62ef` 换掉了。**改字归编排者（票面⛔ 我这把动）**。
- ⚠⛔ 另开普查格（那是票 301/302 的射程）：我只裁票面措辞，⛔ 裁"全仓还有多少枚 tagged 用例进不了执行面"。

## 4. 顶回清单（本程量到的，⛔ 凑数）

1. **支号冲突**：`AC#3` 原文的"③（切错字节）"与票 297 `AC#2` 权威清单（"②字节切错"）互斥，票 297 ⓐ 又用另一把 ⇒ 盘上三处编号⛔ 自洽。详见 `10-ac3-verdict.md` §0。**我按标签裁。**
2. **派单前提过期**：`./internal/audio/` 现已被 `windows)` 档认领（`0a0f62ef`），"从未被执行过一次"只覆盖历史。见本件 §3。
3. **`gofmt -l cmd internal scripts tools` 名册 5 枚**（`logs/gates-gofmt.txt`，含 `cmd\wisp\models.go`）＝⛔ 一枚在我射程（scoped porcelain 全程 0 行）；但 `300-r1`/`A810` 那本账点过的是 4 枚 internal 存量 ⇒ 名册比我预期的多一枚 cmd 侧的，**⛔ 我动、⛔ 顺手格式化**，只具名报回。
4. **`go vet ./internal/audio/` 的 stdout 是 0 字节而 `rc=0`**：本件把两件事分开写（`10` 件 §3 同），⛔ 用"0 字节＝没交"那把尺误判这一格——与 `300-v2` 必答题 3 同判语。

## 5. 欠的读数（具名，⛔ 静默）

- 〔欠，归票 301，等编队空〕`--scope=windows` 的改前／改后名册与 CI 那一条真跑色（`0a0f62ef` 之后本包 tagged 用例到底被求值过没被求值过）。⛔ 我取：`303-a1` 在飞、派单写死⛔ 整包级读数⛔ 起窗。
- 〔欠，⛔ 归我这把〕票 297 `AC#0` 那枚判别仪器（真包前 16 枚样本值）＋ `AC#1` 已知幅度对拉：需要真设备那一发，票面 `:31` 与派单都写死〔仅本机可量、⛔ 归腿〕。
- 〔⛔ 欠，但⛔ 闭合〕票 300 `AC#4`／`AC#5` 两格仍 `- [ ]`：本格判语⛔ 覆盖它们（`10` 件 §2 末行）。

## 6. 终态锚（与起手同形；逐字读数在 `logs/final-anchor.txt`，本件只抄我确实见过的几枚）

| 尺 | 起手（`00` 件） | 终态 |
|---|---|---|
| `date` | `Sat Oct 10 16:56:52 CST 2026` | `Sat Oct 10 17:20:24 CST 2026` |
| `git log -1 --format=%H` | `dce0f133b5123967e527ffbbcf9f3fb3ed4ce14c` | `48ee655edc8ba6ebd903df23a539dc857cbe98d1`（漂过：中间 `303-a1` 等腿的笔） |
| `git status --porcelain \| wc -l` | 819 | 848 |
| `git status --porcelain -- cmd internal scripts tools \| wc -l` | **0** | **0**（⇒ 跨锚证产码面⛔ 漂移，`300-v2` 那把自加尺我这把照做；逐字读数在 `logs/final-anchor.txt`） |
| `tasklist` `wisp.exe`／`balldebug.exe` | 0／0 | 0／0（全程⛔ 起任何窗、⛔ 占麦） |
| 卫生 | — | `sh scripts/d22scan.sh` `rc=0`／`go vet ./internal/audio/` `rc=0`／`gofmt -l` `rc=0`（名册见 §4-3） |

**过程痕迹（具名⛔ 抹）**：本程中途锚（约 15 次调用那一笔）⛔ 单独成件，而是并到 §6 与本件 §1 的时间戳里——我把调用预算花在突变台件上，结果是我的第一笔提交 `f76091cd` 之后一直到交回才再落笔（违反"⛔ 攒到最后一笔"的字面形状，**只建不删、读数全在盘**）。
三发失败/作废具名：`M2`/`M3` 被自己的 `LANDS` 门拦下（锚与 S4 那形撞名）、`M6` 编译失败＝装置红、`rig-mutation-2` 末行一次 `tag: unbound variable`（`local` 同行展开）。⇒ 并回定式候选给编排者：**突变锚要先在本树数一遍命中数再落笔**（我这把是靠门拦住的，⛔ 是靠运气）。

纪律自陈：全程⛔ push、⛔ `git add -A`/`.`、⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`；commit 用显式 pathspec；写面只 `.scratch/wisp/probes/300/v3/**` 新建 `.md`/`.txt`（⛔ `.go`、⛔ `.out`、⛔ 删）；⛔ 改票面/台账/HANDOVER/产码/任何 `*_test.go` 的断言；⛔ 读 secrets 与机主 config 值；⛔ 碰 `thresholds.go`/`allowlist.txt`/冻结件/`frontend/src/**`/`design/**`。

**本格判语（交回，⛔ 宣布翻勾）**：见 `10-ac3-verdict.md` §2 —— `AC#3` 本格判据＝**闭合（成立）**，带两枚限定（"可复现"只到合成夹具＋本机一次 dump；支号③/②盘上⛔ 唯一）；票 297 那支的"是什么"⛔ 定案、那一格仍归票 297 自己的 `AC#0`/`AC#1`；两票各自留格、⛔ 合并结案。

rc=0
