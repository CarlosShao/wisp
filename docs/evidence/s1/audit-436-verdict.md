# 436 只读核查裁决（audit-436-verdict）

- 核查腿：只读复跑（未改任何文件、未 commit、未 push）
- 时刻：2026-11-13 10:40 +08
- 票面：`.scratch/wisp/issues/436-the-247-test-finds-prototype-visuals-on-while-the-shipped-default-was-supposed-to-be-off-done.md`

verdict: HONEST

commit: `d1f12dd55` 存在（`git cat-file -t` ⇒ `commit`，CarlosShao，2026-11-12 23:11:54 +08），改了 3 件：`.scratch/wisp/issues/436-...-off.md`、`cmd/wisp/resident_audio_247_windows_test.go`、`docs/evidence/s1/436-247-prototype-visuals.md`（3 files，+222 −11）。旁证提交 `4b9efc889` 亦存在（票 68 翻转默认值那一枚）。

证据件: 在（`docs/evidence/s1/436-247-prototype-visuals.md`，10730B，2026-11-12 23:11）

AC 框: 勾了 5 个，没勾 0 个（AC#0／AC#1／AC#2／AC#3／AC#4 全 `[x]`）

复跑读数: `--- PASS: TestAC247HandingTheLevelToTheSeamIsNotVisibility (0.06s)`

越界: 没有（`git show d1f12dd55 --name-only` 三件命中 `docs/PLAN.md`／`docs/specs/`／`internal/observe/thresholds.go`／`frontend/`／`design/` 零项；`4b9efc889` 亦零命中）

不符: none
