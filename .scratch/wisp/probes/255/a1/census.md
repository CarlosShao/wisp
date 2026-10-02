# 255-a1 只读取证 —— 「生效级别」词表 census（供编排者裁 AC#2 的 ⓐ／ⓑ）

**代号**：`255-a1` · **性质**：只读取证腿 · **唯一可写件**：本文件
**要答的那一个问题**：今天仓里「生效级别」到底有几套词表、每套谁生产谁消费、彼此有没有映射；
ⓐ（让 `Tier` 真被读者消费）与 ⓑ（具名登记「粒度只到段」+ 一把会响的尺）各要动哪几枚文件、撞不撞现成的钉。
**本格判据由编排者裁，本腿未裁。**

---

## §0 起手锚（逐字读数）

| 尺 | 逐字读数 |
|---|---|
| `date -Iseconds`（进场第一发） | `2026-10-02T10:43:11+08:00` |
| `git log -1 --format='%H %ad %s'` | `a4b221ac098d35c6cf81ac60d3063cf53b9d6b8c Fri Oct 2 10:41:23 2026 +0800 evidence(248-v1b) 骨架：接管腿 248-v1b 的 §0 起手锚落盘……` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain internal/config internal/panel cmd/wisp \| wc -l`（10:43:11 那一发） | `0` |
| 同上复跑（`date` 现取＝`2026-10-02T10:44:35+08:00`） | `1` |
| `git status --porcelain internal/config internal/panel cmd/wisp`（列出那一枚） | ` M internal/panel/config_handlers.go`（`od -c` 复认首列是空格＝**未暂存的修改**） |
| `git status --porcelain cmd/wisp internal/config \| wc -l` | `0` |

**⚠ 本腿取证期间落进工作树的一枚外部改动（不是我改的，我一字节未动产码）**：
脏的那一枚是 `internal/panel/config_handlers.go`，内容是一枚**植入变异**（正文逐字 `return true // MUTATION-M1 arbitrary key passthrough`，
把 `knownWritableField` 的表查询整段替成恒真）。`git diff --stat` 逐字：`1 file changed, 1 insertion(+), 6 deletions(-)`。
⇒ 归属判断（**射程判断，非内容引用**）：这正是对抗验收腿 `248-v1b` 的活儿，与票 255 排程节「此刻 `248-v1` 正在整包跑 `cmd/wisp`」同源。
⇒ **对本腿读数的影响（已逐条核过）**：变异落在 `:368` 之后，故本件引用的 `EffectiveTier` 定义区（`:193–202`、`:213`）
在脏树与 HEAD 里**行号全等**；`tierSentence` 一枚则不等——现尺：
`git cat-file blob HEAD:internal/panel/config_handlers.go | grep -n 'func tierSentence'` ⇒ **`428`**；
`grep -n 'func tierSentence' internal/panel/config_handlers.go` ⇒ **`423`**。
⇒ **本件凡引 `internal/panel/config_handlers.go` 的行号，一律以 HEAD blob 为准**（HEAD 是不动的锚，工作树在取证期间会漂）。

**取证面写面为空的复认**：`cmd/wisp` 与 `internal/config`（票 255 排程节点名的两枚写面）在本腿每一发尺时都是 0 行脏，
故 ⓐ／ⓑ 两形的代价表读的是**已提交形状**，不是谁的半成品。

**根目录纪律**：本件所有 `grep`／`find` 的根**逐条显式写作** `cmd internal tools scripts`（Go 取证面）或
`.scratch docs`（件面），⛔ 未用 `.` 当根，⛔ 未触及 `frontend/**`／`design/**`（两层禁令：不读不引）。

**禁跑尺的遵守**：本腿未跑 `go build`／`go vet`／`go test`／任何 `./...`；未跑 `go list`（依赖边全部由 `grep` 的
`import`／`pkg.` 前缀现读得出，不需依赖图，故按票面「用不上就别用」省掉）。

本节以下 §1／§2／§3 的表格在后续 commit 填写；本腿按票规**先写满 §4、§5 再回填**。

---

## §1 「生效级别」词表的现状表（每套一行，四问）

本节在 §4／§5 落盘后的下一发 commit 填写。

---

## §2 ⓐ 形代价：让 `Tier` 真被读者消费

本节在 §4／§5 落盘后的下一发 commit 填写。

---

## §3 ⓑ 形代价：具名登记「粒度只到段」+ 一把会响的尺

本节在 §4／§5 落盘后的下一发 commit 填写。

---

## §4 我可能判错的条目

## §5 判不动／没测到的地方

## §6 我推翻编排者哪一句

## §7 结论
