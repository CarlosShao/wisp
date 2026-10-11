307-r1 收口补件：共享工作树脏面存档 + 节号更正

- 41-boundaries.md 新增 §3：工作树里躺着别家的未提交改动（.gitignore 的 .worktrees 段、
  go.mod 的 go-webview2 直接依赖改写、design/** 的十余枚删除、probes/152|161|242|268 与
  一堆根目录未跟踪件），本腿⛔ 写过一枚。尺＝本腿四笔的 git show --name-only 合并后
  grep -E '^(go\.mod|go\.sum|\.gitignore)$' | wc -l = 0 ⇒ ⛔ 一枚带进 commit。
  ⛔ 掩一处尺缺陷＝起手⛔ 量全工作树（只量 internal cmd），所以"进场时已脏"是推理⛔ 是实测。
- 欠读数那一节由 §3 renumber 成 §4（msg-04 里我写的是 §3，具名报回）。
- 逐笔越界复量：c6691756 / 51a0a0f2 / 183032c1 / fb6a171e 各 0 格越出自家两面写面。
- ⛔ 勾任何一格；本笔零 push。
