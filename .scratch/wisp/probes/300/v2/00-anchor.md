# 300-v2 — `00` 验收腿锚点与仪器前置（非实现者，攻 `300-r1` 的 `AC#2`）

票＝`.scratch/wisp/issues/300-parsewaveformat-subformat-offset-reads-two-bytes-past-the-guid.md`
被裁的两笔＝`028529fc`（产码＋用例＋证据件）／`41329475`（提交后复核）。

## 1. 起手两读（⛔ `TZ=` 覆盖，钟点＝`date` 原生 stdout）

```
$ git rev-parse --short HEAD
41329475
$ date
Sat Oct 10 13:45:29 CST 2026
$ date '+%Y-%m-%d %H:%M:%S %z'          # 同一把尺的 %z 形
2026-10-10 13:56:06 +0800               # （这是卫生那一发里现取的那一枚，⛔ 我换算）
$ git branch --show-current
dev
$ git status --porcelain | wc -l
814
$ git status --porcelain -- internal cmd | wc -l
0
```

- 起手锚＝**`413294755ea09ea8f18c94b8bb1276b5e702583a`**（短号 `41329475`），与派单给的号**逐字相同**。
- `git status --porcelain -- internal cmd`＝**0 行** ⇒ 起手时 `internal`＋`cmd` 名下工作树↔HEAD 无差（`300-r1` 的两枚都已入库）。
- ⚠ **HEAD 在我这一程中间漂过**（别人在飞）：`git reflog` 现读
  `849ce6e9 HEAD@{2026-10-10 13:52:33 +0800}: commit: 301-a2: AC#3 two-form cost tables …`
  ⇒ 我的**闭合锚＝`849ce6e9`**，而我所有导出树都按**显式号**建（`4eb29312`／`41329475`），⛔ 一棵是从"当时 HEAD"建的。
  尺＝`git diff --name-only 41329475..HEAD -- internal cmd` ⇒ **0 行**（那笔 301-a2 只动 `probes/301/**`＋票面）⇒ 我这把的**产码面**⛔ 因 HEAD 漂移而变。

## 2. 全仓脏面（⛔ 我造成、⛔ 我动；只是量下来免得判错红）

```
$ git status --porcelain | awk '{print $1}' | sort | uniq -c | sort -rn
    783 ??
     16 M
     16 D
```

- 16 枚 `D` 里含 **`design/assets/{base.css,icons.js,theme.js,tokens.css}`** 与十枚 `design/screens/*.html`、`design/index.html`
  ⇒ 派单点名的那枚本机工作树差**实测在盘**（本仓定式＝这类 `TestC21DesignTokensFourWayAgree` 一族只在本地红、CI 复现不上）。
  我撞上了就在红名册里具名标"本地树差、⛔ 本票"。
- 其余＝他腿未跟踪件／`.scratch/wisp/**`（517 枚）／`M .gitignore`。⛔ 一枚进我的射程、⛔ 一枚进我的 commit。

## 3. 仪器前置（`A802`／`A804` 那两把，跑任何带 harness 的包之前）

```
$ git ls-files third_party/sherpa-onnx | wc -l
0
$ grep -rln 'sherpa\|onnxruntime' internal/audio | wc -l
0
```

⇒ 三枚 DLL ⛔ 被跟踪 ⇒ `git archive` 的导出树里**连目录都不存在** ⇒ 跑 `cmd/wisp` 那档前必 `mkdir -p` ＋ `cp`。
⇒ `internal/audio` 与那三枚 DLL 无关（**我这把现量 0 命中**，⛔ 引用派单的断言），所以窄档（`-run` 只打 audio）⛔ 受这枚坑影响。

两棵导出树的 dll 前置尺（逐字）：

```
rc-git-archive=0 rc-tar=0            # 两棵同值（PIPESTATUS 未混：archive 与 tar 各一枚 rc）
treeA dllcount=3
treeH dllcount=3
```

## 4. 我建的树（⛔ 在共享工作树里做任何突变）

| 树 | 来源（显式号） | 用途 | 偏移那行逐字 |
|---|---|---|---|
| `/tmp/wisp300-v2/treeA` | `git archive 4eb29312`（＝`028529fc^`，改前） | 成对基线"改前"那一发 | `196:		sub := *(*windows.GUID)(unsafe.Add(p, 26))` |
| `/tmp/wisp300-v2/treeH` | `git archive 41329475`（＝起手锚） | 五枚偏移扫描＋成对基线"改后"那一发 | `211:		sub := *(*windows.GUID)(unsafe.Add(p, 24))` |

- 新用例存在性尺：`treeA` 里 `ls internal/audio/parse_wave_format_300_windows_test.go` ⇒ `No such file or directory`；`treeH` 里存在（`-rw-r--r-- … 7267`）。
- ⚠ **7267 字节**＝`ls -l` 那一把；`30-gates.md` §4 记的是 **7194 字节**（那一把尺＝`ls -l | awk '{print $5}'`，量的是它自己的 treeB）。
  两枚⛔ 同值⇒ 我这把⛔ 与它的读数对拉成同一枚尺（差 73 字节＝它后来那次"只改注释散文"的收紧）；
  **盘上工作树那枚＝`wc -c` 现量见 `10` 件 §0**。

## 5. 我这把的写面（⛔ 动任何产码／用例／票面框）

只新建 `.scratch/wisp/probes/300/v2/*.md`（证据＝原始读数嵌进正文，⛔ 新建 `.sh`／`.ps1`／`.out`）。
本程⛔ 翻过一枚 `- [ ]`、⛔ 改过 `internal/**` 或 `cmd/**` 一个字节、⛔ push。
