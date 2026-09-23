# 票 133 AC#1 / AC#2 —— 非实现者验收方的裁决表

**执行方**：`acceptor-ticket133-r1`（**非实现者**，只读被验代码；不复用实现方的任何判据）
**被验版本锚定**：`git rev-parse --short HEAD` = **`52cf311`**（本验收方开工时刻的工作树 HEAD，
仅供对照）；**被验 sha = `c5f140c`**（= `c5f140c7a1c2be66ea16f1fcb09ae56737ea1ccf`）。
**纯净树取法**（不在仓库内建 worktree、不 checkout）：

```
$ mkdir -p /d/tmp/wisp133-acc-r1
$ git archive c5f140c | tar -x -C /d/tmp/wisp133-acc-r1        # 有仪器的被验树
$ mkdir -p /d/tmp/wisp133-acc-r1-noins
$ git archive c720494 | tar -x -C /d/tmp/wisp133-acc-r1-noins  # 无仪器的"①栏"树（前任的基线树）
$ cp third_party/sherpa-onnx/*.dll <树>/third_party/sherpa-onnx/   # 三枚 DLL，缺即 0xc0000135
```

`git diff --stat c720494..c5f140c -- cmd/wisp/` ⇒ 只有 `leg_dispatch_gate_133_test.go` +1377 与
`main.go`/`panel_assets.go`/`slo_windows.go` 的 +19/+8/+6 行（裁决注释），与前任 §4.5 同形。

开工时刻 `date -u` = 2026-09-23 12:11:01z（本地 20:11 +08）。

## 0. 骨架（本节以下逐节填，每裁一格 commit 一次）

| 格 | 判据出处 | 判定 |
| --- | --- | --- |
| AC#1 五发（N-3／X4／X8／X12／X14 两拍，门关着仍红） | 票面 AC#1 + 归属硬线 | （待填） |
| AC#2 覆盖面主张（清单式 vs 图式＋"删掉它哪条用例会红"） | 票面 AC#2 | （待填） |

## 1. 独立基线（未变异）

（待填）

## 2. AC#1 五发独立复跑（①无仪器／②门开着／③门关着）

（待填）

## 3. 首要攻击：这枚仪器自己的洞（同名方法／helper 名）

（待填）

## 4. 同形排查（first-match／按下标配对／只认前缀／只数不核）

（待填）

## 5. §2.0.1 语义裁定（票面"先证明现在不红"还剩几发）

（待填）

## 6. AC#2 覆盖面主张裁定＋"删掉它哪条用例会红"

（待填）

## 7. 四数账与门禁抽查复算

（待填）

## 8. R-133-x 清单

（待填）

## 9. 总判

（待填）
