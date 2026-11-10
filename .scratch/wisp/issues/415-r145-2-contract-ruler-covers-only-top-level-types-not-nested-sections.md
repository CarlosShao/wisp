# 415 — `R-145-2`：合同尺只配顶层类型，嵌套 section 没有尺 ⇒ Go 加了 6 键、`panel.ts` 一枚没跟上，零红

- **Status**：ready-for-agent（编排者 2026-11-09 22:5x 立；起手锚点＝本票现读 `git rev-parse --short HEAD`）
- 来路＝票 145 验收腿 `145-accept-r1` 登记的 `R-145-2`；145 已于 `ff6b8a2c5` `-done`，本枚**不阻断**那张票。
- 关联：票 145（AC#6 六枚新键）· `internal/panel/composer_test.go:48` 的 `pairs` 表

## 现量的形状（票 145 验收腿原文，落地腿起手必须复跑）

- **缺口**：合同尺 `internal/panel/composer_test.go:48` 的 `pairs` **只配 6 个顶层类型**；AC#6 的 6 枚新键**全在嵌套 section 内** ⇒ 尺够不着。
- **★后果（现量）**：**Go 侧加了 6 键、`panel.ts` 一枚没跟上，零枚用例变红**。验收腿点名 `panel.ts` 里 `display`/`adapter`/`adapterThinkingLevels`/`argsJson`/`running` 等键缺失而**没有任何测试察觉**。
- ⇒ 这是**仪器缺口**，不是某张票的缺陷：合同尺的射程决定了"Go 与 TS 漂移能不能被发现"。

## 要建什么

- [ ] **AC#1 尺的射程伸到嵌套 section**：让 `pairs` 能覆盖 `Snapshot` 里**嵌套结构**的键（不只 6 个顶层类型）。★形状自选（递归遍历／显式嵌套表／反射比对），落地腿写明**为什么选这个形状**（按「主流 harness 怎么做」判，⛔ 按"改得少"判）。
- [ ] **AC#2 反向判据（本票的灵魂）**：造突变「Go 侧加一个嵌套键、`panel.ts` 不跟」⇒ **新尺必须变红并点名该键**；还原 ⇒ 变绿。**两个读数都贴**。
- [ ] **AC#3 不许把现状打红**：现有 6 枚新键在 `panel.ts` 里**本来就缺**（验收腿现量）。落地腿必须先**判定这 6 键该不该在 TS 侧存在**：该在 ⇒ 补 TS（⛔ 但 `frontend/**` 归外部 agent，本票只出清单）／不该在 ⇒ 尺要能表达"这键 Go-only"。★**⛔ 用"跳过这 6 键"让尺变绿**（那就是把缺口焊进尺里）。
- [ ] **AC#4 门禁**：gofumpt（绝对路径）空／`go vet` rc0／`go build` rc0／`go test ./internal/panel` 绿／`d22scan` clean／path-length GREEN。
- [ ] **AC#5 交付清单**：把"TS 侧缺哪些键"落成一份**具名清单**交外部 agent（本票写面不含 `frontend/**`）。

## 禁区

- 写面＝`internal/panel/**` ＋ 本票 Progress log ＋ `docs/evidence/s1/415-*.md`。
- ⛔ 碰 `frontend/**`／`design/**`（外部 agent）⛔ 碰 `cmd/wisp`／`internal/config`／`internal/ball`。
- ⛔ 动 `thresholds.go`／golden／`PLAN.md`／`specs/**`／`.github/**`。
- ⛔ 改本票 AC 原文；⛔ push。

## Progress log (append-only, newest last)

- [2026-11-09 22:5x +08] 编排者立票（来路 `R-145-2`；尺只配 6 顶层类型、6 新键全在嵌套内、零红逐条引；已写明"⛔ 用跳过让尺变绿"）。
- [2026-11-10 17:0x +08] 落地腿交付（commit `551d2a16d`，只动 `internal/panel/composer_test.go` +390 行，未 push；起手锚 `dd2974e73`）。AC#1＝双侧类型递归走查（`composerTSNameMap` 26 枚翻译＋exact-set 双向；形状＝类型派生，quicktype/ts-rs/OpenAPI 族做法，理由写在尺注释里）。AC#2＝`TestPlantedNestedContractDriftGoesRed` 四发全绿（删真嵌套键 toolsDropped⇒红 1 枚点名；幽灵键⇒红；账本吸新/藏旧⇒红；剥 GitWorktree 翻译⇒红点名；真树 0 违规）。AC#3＝11 枚 Go-only 键逐键判"该在 TS 侧"，落成 `composerTSGapLedger` exact-set 自翻红注册表（⛔ 非跳过：TS 声明即 stale 红、Go 加未注册键即红、Go 删注册键即红），逐键理由＋页侧清单＝`docs/evidence/s1/415-r145-2-nested-contract-ruler-impl-v1.md`。AC#4＝gofumpt 空／vet rc0／scoped build rc0（`go build ./...` rc1＝预存在未跟踪 gitignored `build/mut_rvw.go` 他腿遗留，本腿零触碰）／panel 包合同面全绿、唯一红＝预存在 HEAD 红 `TestC21DesignTokensFourWayAgree`（其输入文件 git 全干净，归口 65/239 族）／d22scan clean（452 件 0 findings）／path-length GREEN。AC#5＝11 键＋8 子接口具名清单已入证据件。诚实项：panel.ts 补键归外部 agent 未做；本件与证据件按派单"只 commit internal/panel/"留盘未 commit。