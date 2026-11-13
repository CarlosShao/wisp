# audit — 票 437 收口核查（只读腿，2026-11-13）

verdict: HONEST

commit: 两笔均在盘。写腿 `3fe4f7365` 存在（commit），改 3 件：`.scratch/wisp/issues/437-the-313-door-reports-seven-writable-fields-the-two-rosters-do-not-name.md`（+13）、`cmd/wisp/panel_config_313_test.go`（+3）、`docs/evidence/s1/437b-313-assertion-five-rosters.md`（+183），合计 +199/-0。取证腿 `018b3112d` 存在（commit），改 2 件：同票面（+4）、`docs/evidence/s1/437-313-door-union.md`（+227），合计 +231/-0。

证据件: 在。`docs/evidence/s1/437-313-door-union.md`（227 行）与 `docs/evidence/s1/437b-313-assertion-five-rosters.md`（183 行）均在盘，均在 `018b3112d`／`3fe4f7365` 内。

AC 框: 勾了 5 个（AC#0–AC#4 全为 `- [x]`），没勾 0 个（`- [ ]` 计数为 0）。

复跑读数: 逐字 `--- PASS: Test313WritableUnionEqualsTheTwoRosters (0.00s)`。全族 `^Test313` 另两行逐字 `--- PASS: Test313EmptyConfigToStepfunRowChainAndPair (0.15s)`、`--- PASS: Test313MultiElementCommaChainKeepsEveryElement (0.05s)`，末行 `PASS` ⇒ PASS=3 / FAIL=0。载具 `D:\work\soft\Git\bin\bash.exe -c "cd '/d/work/workspace/projects plans/Wisp' && scripts/wisp-cli-staged-tests.sh '-test.run' '^…$'"`（本机 PATH 的 bash 是 WSL 转发器，未用）。

越界: 没有。两笔 commit 的 `git show --name-only --format=''` 逐行扫过，`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／`frontend/**`／`design/**` 零命中；`internal/panel/**` 亦零改动（写腿只碰 `cmd/wisp/panel_config_313_test.go` 一处三行 append，已对盘上 249–252 行核过）。

不符: 票面标题与正文仍留旧数「门报 26 个／多出 6 个」，而实测为 27/7——但票面自己已自曝此数有误、取证腿已机械重数并改正，属票面已知笔误、非收口不诚实。
