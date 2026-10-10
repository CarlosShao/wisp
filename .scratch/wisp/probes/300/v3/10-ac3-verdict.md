# 300-v3 — `10` 判语：票 300 `AC#3` 那一格（判语归位，⛔ 越权结案、⛔ 翻框）

我这把＝非实现者（落地那程＝`300-r1`，件 `probes/300/r1/**`）⇒ 符合 AGENTS §0 第 3 句／D22 双角色。
锚＝起手 `dce0f133`（本程所有读数都在这枚号或其后的工作树上取，⛔ 跨锚引别腿的数：`A810`/`153` 那条"读数的世代"）。

## 0. 先落一枚顶回（派单与票面都要，⛔ 客气话）：`AC#3` 那句里**那个支号在盘上不唯一**

`AC#3` 原文逐字（尺＝`grep -n 'AC#3` 与票 297 那三形的关系' .scratch/wisp/issues/300-*.md`，⛔ 引行号）：

> `AC#1`／`AC#2` 成立⇒ 票 297 三支里的 **③（"切错字节"）**拿到一枚**可复现的形状**，但⛔ 本票不许把 297 的"是什么"判成已定案

而票 297 自己的**权威清单**（`AC#2` 那格，尺＝`grep -n '三形择一处置' -A 6 .scratch/wisp/issues/297-*.md`）逐字把号分给另一支：

- `**①声音根本没到麦**` ⇒ 仪器与流程问题
- `**②字节切错**（`convertPacket` 那一层的 frames/channels/bits 拼装）` ⇒ **产码缺陷**
- `**③设备本底就是这样**` ⇒ 要改判据形状、要摆机主一句话

⇒ 同一枚仓里"切错字节"这个标签**同时被叫作 ② 和 ③**：票 297 `AC#2` 清单＝②，票 297「编排者裁定」ⓐ 那句（逐字 `三支候选里 ② "设备本底真的高"仍只能靠真机那一对…，⛔ 因为 ③ 可测就宣布 ② 出局`）＝③，票 300 `AC#3` 跟的是 ⓐ 那把。
**我的处置＝按文字标签裁（"切错字节"那一支），⛔ 按号裁**；号本身是一枚该由编排者收的口径账（同 `A807`④／`A810` 那族"枚数必带单位名"，这次是**支号必带出处**）。

## 1. 三段判语（每段：凭据件路径＋尺逐字＋红/绿句逐字）

### ① **成立** —「`AC#1`／`AC#2` 成立 ⇒ 那一支拿到一枚**可复现的形状**」

- 凭据件：`internal/audio/parse_wave_format_300_windows_test.go`（tracked，四形 `S1..S4`）＋我这把的仓外突变台件（脚本入库＝`.scratch/wisp/probes/300/v3/logs/rig-mutation.txt`／`rig-mutation-2.txt`／`rig-mutation-3.txt`，原始 stdout 在仓外 `/tmp/wisp300-v3/logs/*.txt`，只建不删）。
- 尺（逐字，导出树＝`git archive dce0f133` 到 `/tmp/wisp300-v3/tree`，起手 `cmp` 两枚文件 ↔ `git show dce0f133:<path>` 皆 `IDENTICAL`）：
  `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/audio/ -run TestParseWaveFormatSubFormatOffset300 -v -count=1 -timeout 300s`
- **绿句逐字（pristine 正控）**：`@@@ tag=M0-pristine rc=0 top-PASS=1 subtest-PASS=4 FAIL=0 load0xc=0`（终态复跑同值＝`N3-final-recheck`，且两枚文件 `final … = IDENTICAL` ⇒ ⛔ 任何突变残留在被跟踪文件里）。
- **红句逐字（把产码偏移写回 26＝已知坏形，我这把那一代）**：`@@@ tag=M1-prod-offset-26 rc=1 top-PASS=0 subtest-PASS=0 FAIL=5` ＋四形逐字
  `tag = 0, want 3 (KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24)`／
  `tag = 0, want 1 (KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2475, Data1=0x00000001) at byte 24)`／
  `tag = 0, want 7 (positive control: …)`／`tag = 3, want 0 (positive control for offset sensitivity …)`。
- ★形状是**可复现**的，且⛔ 只在测试里复现——下游那枚后果也复现了：红句里逐字带着 `floating = false, want true (convertPacket would take the float32 path)`。

### ② **成立，但射程只到本机此刻** —「真机那枚混音格式到底是不是 extensible＋float」＝**是**

- 凭据件：`probes/300/orch/2026-10-10-ac3-getmixformat-hexdump.md`（跑＝编排者本人，载具在仓外导出树、⛔ 入库；件末 `rc=0`）。
- 那把尺（逐字，件 §2）：`go test ./internal/audio/ -run TestOrch300MixFormatHexdump -v -timeout 180s` ⇒ `rc=0`／`--- PASS`；跑的是哪一版码由 `cmp` 导出树 ↔ `git show HEAD:internal/audio/wasapi_windows.go` ＝`IDENTICAL`、两边同一内容锚都是 `unsafe.Add(p, 24)` 钉住。
- 关键两枚字节（件 §5 解码表逐字）：偏移 `24` = `03 00 00 00` = `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT`，偏移 `26` = `00 00`；四枚默认端点（render/capture × console/multimedia）**同一份 40 字节**、`cbSize=22`、2ch/48000/32bit、`mask=3`。
- ⛔ 说成别的机器／别的驱动／别的采样率、⛔ 说成采集链路跑得通：件 §1 自己写死了这两条边界（"只这一台、只这一刻、四枚端点"）。

### ③ **⛔ 成立，而且本格⛔ 要求它⛔ 成立** —「票 297 那枚 `0.3677` 的底**就是**这一支造成的」

- 尺（逐字）＝`grep -n '^- \[ \]\|^- \[x\]' .scratch/wisp/issues/297-quiet-room-level-sits-at-a-near-constant-0-37-so-loud-is-not-separable.md` ⇒ **5 枚 `- [ ]`／0 枚 `- [x]`**：`AC#0`（那枚判别仪器：真包前 16 枚样本值逐枚打印）／`AC#1`（已知幅度注入面对拉）／`AC#2`（三形择一）**全部未落**。
- 尺（逐字）＝`ls .scratch/wisp/probes/297/` ⇒ 只有 `a1/`（只读普查腿；它自己那条 Progress log 逐字写着 `⛔ 零 go build/vet/test/-list`）⇒ **盘上没有任一件交出过"喂进电平尺的那串字节长什么样"**。
- 我这把也⛔ 能补这一发：票面 `:31` 与 `AC#3` 都把真机那一发写死〔仅本机可量、⛔ 归腿〕，而我派单里另有一条⛔ 起任何窗（`303-a1` 在飞）。
- ⇒ 这一支**未闭合＝符合本格要求**（票面 `AC#3` 逐字：`⛔ 本票不许把 297 的"是什么"判成已定案`；`⇒ 两票各自留格，⛔ 合并结案`）。它欠的是票 297 的 `AC#0`/`AC#1` 两格，⛔ 本票的账。

## 2. 综合判语（交回编排者，⛔ 我翻框）

**`AC#3` 那一格的判据＝闭合（成立）**，但要带着两枚限定一起翻才诚实：

1. **"可复现的形状"这五个字的射程**＝① 合成夹具（本包四形）＋② 本机四枚端点那一次 dump；⛔ 第三枚独立形状（真包样本值＝票 297 `AC#0`）盘上⛔ 存在。
2. **支号③/②在盘上互斥**（本件 §0）⇒ 今后引用这一格要带**标签**（"切错字节"）⛔ 带号。

另两格⛔ 由本格顺带结案：票 297 `AC#2` 的"三形择一"仍⛔ 开（我这把⛔ 判它）；票 300 `AC#4`／`AC#5` 仍 `- [ ]`（尺＝同 §1③ 那把：票 300 现 **3 枚 `- [ ]`**＝`AC#3`／`AC#4`／`AC#5`，`AC#0`..`AC#2` 已勾）。

## 3. 本格的卫生（⛔ 一把合计数，每把带 rc）

| 尺（逐字） | 读数 |
|---|---|
| `git status --porcelain -- cmd internal scripts tools \| wc -l`（跑门禁前／后各一发） | **0 行**／**0 行** |
| `sh scripts/d22scan.sh` | `rc=0`，末行 `d22scan: clean - no D22 ban violations`；`logs/gates-d22scan.txt` |
| `gofmt -l cmd internal scripts tools` | `rc=0`、**名册 5 枚**＝`cmd\wisp\models.go`／`internal\agent\pending_read.go`…（见 `logs/gates-gofmt.txt`）⇒ ⛔ 一枚在我射程（我这把⛔ 动产码，scoped porcelain 0 行），＝`300-r1` 具名过的既有存量＋`298` 的范围外残留，**⛔ 顺手格式化** |
| `go vet ./internal/audio/` | `rc=0`、stdout **0 字节**（`logs/gates-vet.txt`；这一形按 `300-v2` 必答题 3＝"⛔ 0 字节＝没交"在这格⛔ 落地，本件把 `rc=0` 与"0 字节"两件事都写出来） |
| 我这把的写面 | 只 `.scratch/wisp/probes/300/v3/**` 新建 `.md`/`.txt`；⛔ 新建 `.go`、⛔ `.out`、⛔ 删任何件 |

rc=0
