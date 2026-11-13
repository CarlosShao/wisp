# 432 只读核查裁决

verdict: HONEST
commit: 780fa6ea4 存在（`git cat-file -t 780fa6ea4` = commit；`git show --stat` 见下）；改了 4 件：`.scratch/wisp/issues/432-the-three-258-hotkey-tests-red-on-first-real-execution-because-the-old-binding-and-live-slot-count-premises-moved.md`、`cmd/wisp/resident_hotkey_258_windows_test.go`、`cmd/wisp/resident_hotkey_v1probe_test.go`、`docs/evidence/s1/432-hotkey-258-premises.md`
证据件: 在（`docs/evidence/s1/432-hotkey-258-premises.md`，293 行，含三枚 Test258 测试名）
AC 框: 勾了 5 个，没勾 0 个
复跑读数:
--- PASS: Test258BridgeMutationNoSrcKeepsOldBinding (2.47s)
--- PASS: Test258OccupiedCombinationNamesTheNewValue (1.53s)
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (1.58s)
越界: 没有（`git show 780fa6ea4 --name-only --format=''` 四件里无 `docs/PLAN.md`／`docs/specs/`／`internal/observe/thresholds.go`／`frontend/`／`design/` 任一项）
不符: none
