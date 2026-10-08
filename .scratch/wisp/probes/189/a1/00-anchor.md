# 189-a1（10-08 程）起手锚 — 只读普查腿

> 本程零 go 命令、零产码写点。写点只有：本目录 `.scratch/wisp/probes/189/a1/**` ＋
> 工单 `189-*.md` 追加的一节（不改写既有行、不勾任何框）。

## 锚（本程自己现量，不是"应为"）

| 项 | 值 | 尺原文 |
|---|---|---|
| branch | `dev` | `git rev-parse --abbrev-ref HEAD` |
| HEAD | `914177e60aed5d768988e74f180be42ec29857e8` | `git log -1 --format='%H %ci %s'` |
| HEAD 提交时刻 | `2026-10-08 10:31:36 +0800` | 同上 |
| 本腿起手时刻 | `2026-10-08 10:33+0800` | `date '+%Y-%m-%d %H:%M%z'` |

派单说"起手 HEAD 不早于 `914177e6`"——现量恰等于它，但**本件按"现量值"记，不按"应为"记**
（本项目已作废"应为 <sha>"写法；并行腿会推进 HEAD，后续读数若与本锚行号不符，以现量为准）。

## 工单 189 文件名一枚（不是"一堆"）

派单原文写「它票池里 **`.scratch/wisp/issues/189-*.md`** 就是这一堆」——现量：**只有 1 枚文件**。

```
尺：ls -1 .scratch/wisp/issues/ | grep -i '^189'
命中：1 行
  189-the-review-stack-needs-uncommitted-count-and-per-file-diff-but-go-side-has-zero-source.md
rc=0
旁尺：find .scratch -maxdepth 3 -iname '*189*'
命中：2 行（工单 1 枚 ＋ 派单件 .scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md）
rc=0
```

⇒ 没有 `-done` 后缀件、没有 189 的分裂票；"这一堆"在本仓＝**一枚 189 工单 ＋ 它引用的两枚既有交付件**
（见 §下一节），不是多枚待办票。

## 框数（派单指定的带缩进那把尺）

```
F=189-the-review-stack-…-zero-source.md   （cd .scratch/wisp/issues）
未勾：grep -cE '^[[:space:]]*- \[ \]' $F  ⇒ 6      rc=0
已勾：grep -cE '^[[:space:]]*- \[x\]'  $F  ⇒ 0      rc=1
松口径（不限缩进/前缀）：grep -c '\[x\]' $F ⇒ 0      rc=1
其它框字符：grep -oE '^[[:space:]]*- \[[^ ]\]' $F | sort | uniq -c ⇒ 无输出  rc=0
```

⇒ **6 未勾 / 0 已勾**（AC#1–AC#6 六格全开）。与 `.scratch/wisp/probes/182/c1/100-census.md` K6 那格
"六框全未勾"合一，无冲突。本件尺带 `[[:space:]]*`，故与 `^- [ ]`（漏缩进那把）不同分母；本票面 AC
六行**恰好都在行首无缩进**，所以两把尺在本件上读数相同——这是巧合不是通则，别据此换尺。

## 本程不重造的两把已量过的尺（只引用）

- `.scratch/wisp/probes/182/c1/100-census.md`（锚 `17a54338`，10-08 同日）：
  K6 那格已给 `grep -rn "git diff\|diff --numstat" --include=*.go internal/ cmd/` ⇒ **0**（rc=1，正控已打）
  与 `grep -rn "uncommitted\|numstat\|Uncommitted" internal/panel internal/tools cmd/wisp --include=*.go | grep -v _test.go` ⇒ **0**。
- `.scratch/wisp/probes/167/c2/census.md`（锚 `0ce6cb91`）：占用条／排队序号／停止三堆。

⇒ "有没有源"这一问**本程不重跑**，只按派单要求改答"要做要动哪几面、贵在哪、碰不碰契约面"。

## 一个前置事实（会影响本程定位，先记）

工单 189 的 Progress log `09-28 15:2x / 17:5x / 18:1x` 三行说明：**同名 `189-a1` 的只读设计核早在 09-28 已交过一次**，
件在 `docs/evidence/s1/189-uncommitted-count-design-a1.md`，编排者 09-28 18:1x 已收并**照准路线与落点**（账 `A395`）。
⇒ 本程（10-08 的 `189-a1`）与派单动机"再证一次有没有源"无关，也**不是从零做设计核**；
本程的实际增量＝① 复核 09-28 那件的**前提是否被十月产码改动弄过期**，② 出"最少落点表"，③ 雷区复核。
派单那句"要的是最少要动哪几面"与 09-28 件的 §②「落点」有重叠，重叠部分本程只复核不重造。
