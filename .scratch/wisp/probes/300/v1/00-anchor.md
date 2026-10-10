# 300-v1 — `00` 锚与起手闸门（非实现者验收腿，⛔ 盖章、⛔ 修产码）

派单编号：`300-v1`。腿种：**验收腿（攻击既有凭据）**，本程**零产码**、⛔ 改任何跟踪源件。
被审的腿＝`300-a1r`（凭据＝`.scratch/wisp/probes/300/a1r/{00,10,20,90}*.md`）。
票面＝`.scratch/wisp/issues/300-parsewaveformat-subformat-offset-reads-two-bytes-past-the-guid.md`。
本程射程＝**只新建** `.scratch/wisp/probes/300/v1/**`（正文 `.md`＋原始输出 `logs/*.txt`）＋票面 `Progress log` 追加一行。

## 起手闸门（第一条命令，逐字 stdout）

命令逐字（派单写死的第一发）：

```
date "+%Y-%m-%d %H:%M:%S %z" && git rev-parse --short HEAD && git status --porcelain -- internal cmd | wc -l
```

```
2026-10-10 12:36:25 +0800
efd85155
0
```

- **锚＝`efd85155`** ＝ 派单给的派单锚 ⇒ **没漂**，本程全部读数钉在 `efd85155`。
- ⚠ **锚与被审凭据的锚不是同一枚**：`300-a1r` 的锚＝`b2933dc9`（它的 `00-anchor.md` 逐字写着 `2026-10-10 12:18:27 +0800` / `b2933dc9`）。
  `efd85155` 与 `b2933dc9` 之间**有没有人动过 `internal/audio`**，我用两把尺现量（见下面 §环境 末条），⛔ 假设没动。
- 第三段尺＝`git status --porcelain -- internal cmd` 的 **stdout 行数**（`| wc -l`）⇒ **0 行**。
  ⇒ 起手时 `internal/**`＋`cmd/**` 两棵子树逐路径干净（与被审腿开工时同形）。

## 环境（第二条命令，逐字读数）

```
git rev-parse --abbrev-ref HEAD && go version && go env GOOS GOARCH GOFLAGS
```

由被审腿 `00-anchor.md` 记过的同一环境，我自己现跑一遍：
- 分支＝`dev`；`go version go1.27.1 windows/amd64`；`GOFLAGS` 未设；`GOOS=windows` `GOARCH=amd64`。
- 原始读数＝`logs/env.txt`（尺＝`wc -l`）。

## 两枚锚之间 `internal/audio` 有没有被动过（我自己量的，⛔ 引用别人）

```
git log --oneline b2933dc9..efd85155 -- internal/audio | wc -l     ⇒ 0
git diff --stat b2933dc9 efd85155 -- internal/audio               ⇒ 空
```

⇒ 两枚锚之间 `internal/audio/**` **零改动** ⇒ 被审腿的导出树（`b2933dc9`）与我这一程要复跑的内容
**在产码面上逐字节等价**，我的复跑读数可以拿来裁它的凭据，⛔ 因锚不同而失效。
原始输出＝`logs/anchor-drift.txt`。

## 口径纪律（写死在本件，后续每一件都照它）

1. **整棵工作树本来就脏约 805 枚**（`design/**`、`frontend/**` 有删除项，⛔ 我造成、⛔ 我复位）
   ⇒ 本程每一件**只**写"我动过的路径恢复原状"／"`git status --porcelain -- internal cmd` ＝ 0 行"，
   ⛔ 任何一处写"工作树干净"。
2. 枚数一律带**尺名＋射程目录＋量的还是算的**；词面尺逐字抄；格式／名册类写明"工作树 vs `git show` blob"口径
   （本仓实测 CRLF 会造假枚数，见被审腿 `90-hygiene.md` §3.1）。
3. 行号一律当快照；引用产码用**内容锚**。
4. ⛔ 证据件叫 `.out`（根 `.gitignore` 第 8 行逐字 `*.out`，会被静默跳过而 commit 仍回显成功）
   ⇒ 正文 `.md`、原始输出 `logs/*.txt`。
5. `/tmp/wisp300-a1r/tree`（被审腿的导出树）**留着⛔ 删**；我在里面造的突变形**同一条命令内**备份＋跑＋还原＋`cmp`。

## 我这一程要裁的四条必答（派单口径，⛔ 我新造判据）

| 条 | 要裁的事 | 我的件 |
|---|---|---|
| ① | `AC#0` 的权威凭据**有没有牙**：自己开 `mmreg.h` 逐字抄结构体＋packing＋两处 `cbSize = 22`，**自己算偏移**；并裁"票面 `:13` 后半句（`nValidSamples` 与 `dwChannelMask` 同一 union）谁对" | `10-ac0-authority.md` |
| ② | `AC#1` 那套台件**是不是循环论证**：夹具偏移必须来自 `AC#0` 权威、⛔ 来自仓里任一枚实现；四形是不是真"两种形状⛔ 都放行"；**自己重跑**＋突变形还原证明 | `20-ac1-rig-teeth.md` |
| ③ | 入库版（`AC#2`）要不要保留 **S4** 那一形（本格只裁形状，⛔ 写用例） | `30-harness-and-verdicts.md` |
| ④ | 成对导出树的 **harness 坑**：`third_party/sherpa-onnx/*.dll` 未跟踪 ⇒ `git archive` 树里那条 `PATH` 指到不存在的目录 ⇒ 上一程 `293-v1` 的 677／684 要不要打折；结论必写成三档之一 | `30-harness-and-verdicts.md` |

## 起手声明：我与派单／被审件的分歧（先登记，读数在后面每一件里给）

- ⚠ 派单让我"`grep -rln 'sherpa\|onnxruntime' internal/audio` 现找一枚已知依赖 DLL 的用例"——
  **我这把尺现跑＝ 0 命中**（`internal/audio/**` 里既无 `sherpa` 也无 `onnxruntime` 字样，且该包只 import `internal/observe`）。
  ⇒ 派单这条指引在 `internal/audio` 这一层**没有东西可找**；我改用**内容锚**另找真依赖 DLL 的那一枚
  （`cmd/wisp/dataroot_128_windows_test.go` 头顶注释逐字："the sherpa/onnxruntime DLLs - without that the child dies at
  load time with **0xC0000135** and zero output"；`scripts/portable-tests.sh:258` 逐字："the test binary dies at LOAD time
  without the sherpa DLLs on PATH"）⇒ ④ 的射程因此**落在 `cmd/wisp` 那一档**，与派单举的例子不同，**具名报回**，见 `30` 件。
- 被审腿 `10` 件 §A.2 那段 `#if !defined( RC_INVOKED ) ...` 的**行号范围标了 31–36、块里有 7 行**；
  我自己 `sed -n '28,40p'` 量到的是 `#pragma pack(1)`＝第 **33** 行、`#include "pshpack1.h"`＝第 **35** 行。
  ⛔ 影响结论（引文本身逐字对得上），只登记为它那一件的行号快照过期处。
