# 197-a2 §D — 一程能做完的最小增量切法（**建议，⛔ 拍板归编排者**）

> 现量约束：`cmd/wisp` 写面正排 `294-r1 → 313-r1a → 186-r1`（派单存档 `8e4f6f83` 逐字）；Go 测试面串行。
> ⛔ 本腿不派活、不挑支，只给线路图。

## D0 为什么不能「一程把 197 收口」

- 五枚未勾格里，**没有一格欠新产码**——实体/流/载体/正控全在盘上（§A）。缺的是**裁决料**。
- 真正欠 Go 产码的只有**三块**，且**两块撞 `cmd/wisp`**：C4（`197-r3` 尺修）与 C3（票 220 的 `panel_pump.go`）；一块在 `internal/`：C2（票 229 名字）。
- ⇒ 最小切法＝**一枚 internal-only 落地腿 ＋ 一枚只读裁决腿**，把「能不碰 `cmd/wisp` 就推进」的那两样先推掉。

## D1｜leg ①（落地·`internal/` only，**不站 `cmd/wisp` 队列**）：票 229 甲形

- **做什么**：把 `internal/tools/subagent_197.go:383-384` 那枚收尾协程名 `"subagent-finish-"+bg.ID` 换成 D38 名册内**已有的形状**（每任务前缀族，如 `agent-task-finish-<id>` 一类，具体名以名册三张表现读为准）。
- **为什么是它**：① 写在 `internal/tools/`，**零 `cmd/wisp` 写面**，不被队列挡；② 是 AC#6 链上唯一「名册外协程」实体缺陷；③ 票 229 自带判据齐（AC#1 名册现读／AC#3 正控：种假名字必须仍打 WARN／AC#4 全仓同类普查）。
- **⛔ 红线**：不许走乙形（动 D38b 表＝人工批准）；不许把 WARN 调轻/删。
- **交付面**：`internal/tools/subagent_197.go`（2 处）＋ 票 229 那枚正控/普查件。
- **验收**：非实现者 `229-v1`（须跑整包 ⇒ **排 `cmd/wisp` 队列之后**，或错峰在写腿空窗）。
- 尺：`sed -n '383,384p' internal/tools/subagent_197.go`；`grep -c "subagent" internal/observe/goroutine.go`。

## D2｜leg ②（只读·裁决，写面只在 `docs/evidence/s1/`）：`197-v2` 五格裁决表

- **做什么**：一枚非实现者裁决腿，逐格判 `AC#0`／`AC#1`／`AC#2`／`AC#3`／`AC#6`「盘上形状是否＝票面声称」，给逐格正控（不是只跑绿）。
- **落点**：`docs/evidence/s1/197-v2-*.md`（**零 Go 产码**）。
- **要 `cmd/wisp` 的只是读数面**：须带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 跑 `go test ./cmd/wisp ./internal/...` ⇒ **排在 `294-r1`／`313-r1a`／`186-r1` 之后**，并自证没有别的 `go test` 在飞（先例：`197-r1` §⑤ 让出窗口）。
- **每格的最小判据**（本腿已把落点量清）：
  - AC#0＝票面 §AC#0 那段（S5 为家/S7 为闸）逐字读回＋`PLAN.md:1544`／`:1447` 原文仍在。
  - AC#1＝`run.go:780` 生产调用者＋常驻侧 `resident_task_source_windows.go:462` 共用 `run.execute`。
  - AC#2＝`subagent_197.go:99-104` 四枚态逐枚＝D43 名；`subagent_roster_197_test.go:246` 那枚 D43 外态不得上线。
  - AC#3＝`subagent_stream_197_test.go:57 TestNPlusOneSubagentsEachGetTheirOwnStream` ＋ `94ea50f7` 的 M1/M2 正控复跑。
  - AC#6＝`subagent_197.go:327` 不级联＋`:401` 写给模型看＋`task.go:742` 停止出口（票 221 已落）；D1 若已落则补那枚名字。
- ⛔ 不许把该表当「翻框依据」由本腿自己用工单框——翻框归编排者。

## D3｜腿序与队列（建议，不拍板）

```
现在可起：  D1 落地腿（internal/tools only，零 cmd/wisp 写面）
cmd/wisp 队列退空后：D1 的 229-v1 验收  →  D2 的 197-v2 裁决
更后（撞 cmd/wisp，别并进上面两腿）：
            197-r3 尺修（C1/C2/C3，cmd/wisp/subagent_selfapproval_197_test.go）
            票 220（panel_pump.go join＋gate.go L1 枚举口，须先落票 194 批准面）
```

- **为什么要这么切**：D1 是唯一能在 `cmd/wisp` 被占时推进的落地；D2 是五格唯一缺的东西，且它天生只读、串行门禁窗口的代价最小。
- **⛔ 不要并进任何一腿的**：C5 界面那一跳（`frontend/**` 归 owner）；AC#5 的 C1／C2／C3（撞 `cmd/wisp`，须单开 `197-r3` 并同批重跑 P1／M1 真形／D1 三发正控）。

## D4｜切法不覆盖的（具名，别当已收口）

- **AC#4 后半**（票 220）＝**不在上面两腿内**，它要动 `cmd/wisp/panel_pump.go` 且可能撞票 194 批准面。
- **AC#5 四条附条件**里的 C1／C2／C3＝**不在上面两腿内**（单开 `197-r3`）。
- **界面红**（`PanelSnapshot` 缺两键）＝**不在上面两腿内**（owner 侧）。
