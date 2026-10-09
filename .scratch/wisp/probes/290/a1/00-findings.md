# 票 290 AC#0 + AC#1 — 只读普查件（腿 `290-a1`）

**生成腿**：`290-a1`（只读普查；AC#0 契约普查 + AC#1 三形代价表）
**锚点**：分支 `dev`，HEAD `34a012599118516ae41cc56ff9453216df875779`（2026-10-09T14:03:23+08:00）
**性质**：⛔ **本格零产码改动、零默认值改动、零 AC 翻框**。
**纪律自报**：
- ⛔ **我没跑 `go build` / `go vet` / `go test` 任何一枚**（Go 编译/测试面此刻由写腿 `289-r1` 独占，树须保持干净 HEAD）。只跑过 `go list -deps`（只出不入）。
- ⛔ 零读零写 `frontend/**` 与 `design/**`（本件不引它们任何结论）。
- 文件存在性一律只认 git 对象层（`git cat-file -e` / `git ls-tree`）；行号锚全部现跑（下表每把尺带命令原文，可重跑）。
- 本件写点唯一＝`.scratch/wisp/probes/290/a1/00-findings.md`；票面只在 `## Progress log` 末尾追加一行，⛔ 不改票面任何原句。

**脏度前置尺**：`git status --porcelain -- internal/audio/gate.go cmd/wisp/resident_audio_windows.go cmd/wisp/resident_ball_windows.go internal/ball/ball_windows.go internal/statemachine/table.go internal/statemachine/events.go internal/config/schema.go internal/config/manager.go docs/PLAN.md docs/specs` → **输出为空** ⇒ 上面这些文件的工作树 == HEAD，故本件引用的行号既是工作树行号也是 HEAD 行号。

---

## 0. 我对票面的复量（先自证，别照抄编排者的现量）

（待填：`SetMuted` / `SetSpeaking` 非测试调用者枚数尺）

## 1. AC#0 契约普查：「麦克风的开／关由谁发起」

### 1.1 搜过的字样名册（尺与命中数）
### 1.2 逐字引文（带 `文件:行`）
### 1.3 三问结论 ⓐ / ⓑ / ⓒ

## 2. AC#1 三形代价表（⛔ 不落地）

### 2.1 甲形：静音热键那支接到门上
### 2.2 乙形：面板／托盘开关
### 2.3 丙形：说话起始／唤醒词
### 2.4 依赖边作差尺（`go list -deps` 名册）

## 3. 够不到的地方（具名）

## 4. 与票面/代码原文的出入（本仓铁律：以原文为准）
