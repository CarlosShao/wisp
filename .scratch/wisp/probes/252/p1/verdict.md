# 票 252 · AC#1 读数 · 探针腿 `252-p1`

> 本件只交读数，不交修法，不勾 AC。判语归非实现者（`SPEC-12 §4.3` #1/#3、D22 双角色）。

## 0. 起手锚（只读闸门，逐字）

- `date -Iseconds`：`2026-10-02T09:37:00+08:00`
- `git log -1`：`b5b8a09e ledger(A517)：224-c2 结档——它推翻我派单那句根因（"会话授权没接进装配"判为不成立）＋新立票 252（同一路径两种拼法在允许根判定里分家 ⇒ 允许根内新建被判越界 R2 升 L2）＋票 195 我裁完写进票面`
- 分支：`dev`
- `git status --porcelain internal/tools internal/risk`：（逐字，空输出＝两包此刻干净）

```
$ git status --porcelain internal/tools internal/risk
（无输出）
```

起手名册（本腿终态必须与之逐枚具名相等）：**空**。本腿自己的新增件在下面 §6 具名登记。

## 0.1 派单里抄来的尺（非本腿读数，逐条自验前不作数）

- `internal/tools/paths.go:133` `InAllowlist`：第一段 `paths.go:140`、第二段 `paths.go:152-153`。
- `internal/risk/pathresolver.go:127-139`：存在⇒折长；不存在⇒`res.Canonical = p`（`:138`）。
- `internal/risk/rules_gateway.go:45` `if !ctx.canon.InAllowlist(canonical) {` ⇒ 出 R2／L2。
- `internal/tools/bridge.go:334` `if !sil.Silenced && sil.Level == risk.L1 {`。
- 本机长短名：`REPO_LONG=D:\work\workspace\projects plans\Wisp` ／ `REPO_SHORT=D:\work\WORKSP~1\PROJEC~1\Wisp`。

以上四枚行号本腿已用 `grep -n` 复核，复核结果写进 §4。

## 1. AC#1 四问读数

（骨架：本节由探针 `internal/tools/paths_shortname_252_probe_test.go` 的输出填，逐字读数在 `.scratch/wisp/probes/252/p1/probe-run.log`。）

## 2. 短名从哪来（生产链路）

（骨架：本节由 §1 读数之后的码面追写。）

## 3. 门禁读数

（骨架：本节与 §5 在第 100 轮硬预算闸门之前必须写满并 commit。）

## 4. 本腿推翻编排者哪一句

（骨架：逐条列，不许空本节交件。）

## 5. 判不动的地方

（骨架：甲＝谁能补＋命令＋期望读数；乙＝补不上，明写不做。）

## 6. 本腿写面（越界即缺陷）

- 新增 `internal/tools/paths_shortname_252_probe_test.go`
- 新增 `.scratch/wisp/probes/252/p1/verdict.md`（本件）与 `.scratch/wisp/probes/252/p1/probe-run.log`
- `.scratch/wisp/issues/252-*.md`：只在末尾追加一节「AC#1 读数」，不碰任何 AC 勾选框。
