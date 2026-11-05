# 407 — K5 插件面板：后端有读面、界面未接、无规格（182 普查 K5，无票认领）

**Status:** **open**（2026-11-05 编排者按机主 2026-11-05「立新票」决定建；来路＝182 普查 K5）。
**Type:** 功能缺口（插件功能缺一块界面 + 无规格）＋无票认领。
**Packages:** 前端 frontend/**（插件面板 UI，归机主外部 agent）；可能牵 internal/plugin 只读 bridge。
**Blocked by:** 前端 face 归外部 agent；Go 侧 bridge 若有缺口可另立子格。

## 现场
182 普查（docs/evidence/s1/182-census.md）K5：插件面板后台能读（internal/plugin 库有读面），但前端从未接、也无界面规格；普查时无任何 open 票负责。

## 建议
前端 UI 交机主外部 agent；先把「插件面板该显示什么」的规格写清（可用 C17/面板既有形态），再接线。Go 侧若缺 bridge 读口，另开一条本编队可做的子票。
