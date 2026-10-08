# 255-v2 起手锚（非实现者验收腿）

派单给的 HEAD 号 `e05b8a2b` **已漂**，现量 HEAD（`git log --oneline -1`，rc=0）：
`6320a16a 182-c2 逐框六问（票 182 六枚未勾框 × Q1-Q6）...`
其下有 `e922e0d1 35-r7` / `e05b8a2b 111-c2`（派单那个号在 HEAD~2）。

派单三笔 commit 全部存在且为 `d37e75bd..18655988` 连续区间（`git show --stat`，rc=0）：
- `d37e75bdfffadae74d2fec3e365d43b9fe293aa5` 起手锚（只写 probes/255/r1/00-anchor.md，46 行）
- `8b32060b5570cc5c17990475a83030fc2518c29e` 实现＋仪器（5 文件 +702/-30）
- `18655988d567966996afb7ab5d07317dcaff40ae` 证据件＋票面追加（21 文件纯追加 2293/0）

## 现量差异（先记一笔，不采信实现腿自述）
`8b32060b --stat` 里有 **`cmd/wisp/panel_resident_windows.go | 6 +++---`**，
实现腿自报三件改动名册（甲形／改漂／仪器）里**没有这一枚**。
它改的是 `residentPanelGeometryNote` 的文案＋`state=per-create` → `per-create-and-per-reshow`。
⇒ 自报名册**漏计一枚产码文件**，逐条核在 Q1/Q3/Q4。

## 在飞改动（起手 `git status --porcelain`，rc=0，共 771 行，见 logs/01-worktree-start.txt）
工作树起手即脏，含 ` M .gitignore`、` D design/**`（14 枚）、` M .scratch/wisp/probes/**`、
`?? .scratch/.scratch/`、`?? -`（根目录一枚名为 `-` 的未跟踪文件）等，**均不属于本腿**。
本腿只写 `.scratch/wisp/probes/255/v2/**`，收尾再量一次作差；⛔ 不提交、不还原别人的件。
