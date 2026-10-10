# 301-r1 — `AC#1` 落地形（档清单 ＋ pin，同一笔 commit）

## 改动＝两枚 token，一枚文件

尺＝`git diff -- scripts/portable-tests.sh`（工作树 vs `74eb032c`），
`git diff --stat` ⇒ **`scripts/portable-tests.sh | 3 ++-`，1 file changed, 2 insertions(+), 1 deletion(-)**。

1. **`windows)` 档清单**（原 `:248-253` 块）：末行由
   `        ./internal/session/ ./internal/projctx/`
   改为
   `        ./internal/session/ ./internal/projctx/ ./internal/audio/`
2. **`win_pin`**（原 `:193` 起）：在 `github.com/CarlosShao/wisp/cmd/llmrecord` 与
   `github.com/CarlosShao/wisp/internal/ball` 之间插入 `github.com/CarlosShao/wisp/internal/audio`。

⛔ 别的一枚字符：⛔ 注释、⛔ 空行、⛔ 缩进、⛔ 别的档（`core`/`cli`/`winsec`）的清单或 pin。
⇒ 差分只有两行，可读作「一份 token 进档 ＋ 同一 token 进 pin」，这正是 GUARD D 的退码文本
（`:407`–`:423`，插 pin 后为 `:408`–`:424`）里 `:417`–`:418` 那两句要的形：
`:417` 尾 `… Pull the package` ／ `:418` 头 `into a named scope and update that tier's pin in the SAME commit, or`
（两行各以 `echo "portable-tests.sh: ` 起头 ⇒ ⛔ 引用成一行「逐字」）。

## 为什么 pin 的位置＝`cmd/llmrecord` 之后

GUARD C 的比较尺（现量 `:549-550`）＝
`want=$(printf '%s\n' "$pinned" | grep -v '^[[:space:]]*$' | sort -u)` ／ `got=$(cat "$resolved")`，
`$resolved` 是 `go list` 的 stdout（Go 按导入路径字典序输出）。⇒ `want` 被 `sort -u` 归一，
枚数与顺序由 `sort` 决定 ⇒ 位置**功能上无关**；本腿仍按 `win_pin` 既有的字典序插（`audio` < `ball`），
这样 diff 与那份 pin 的排序惯例同形，⛔ 制造一枚「看起来乱序」的读数。

## 一处先例确认（⛔ 自造规矩）

同一枚包同时被两档认领＝**本仓既有形状**，⛔ 本票新造：改动前的盘上现量
`core` 档清单（`:236-244`）里已经同时躺着 `./internal/config/...`、`./internal/secret/...`、
`./internal/risk/...`、`./internal/proc/...`、`./internal/ball/`、`./internal/perm/`、`./internal/plugin/`、
`./cmd/llmrecord/`、`./internal/session/`、`./internal/projctx/`——这十枚**同时**也在 `windows)` 档里。
⇒ `audio` 进 `windows)` 后变成「core ＋ windows 双认领」，与那十枚同形，
GUARD D 的判据（`$all` 里每枚**包**问「有没有任一档认领它」）⛔ 因此改变。

## ⛔ 动 `:593` 那行 ledger（＝派单逐字禁区第 1 句）

改动前该逐字行在 `:593`；本腿在 `win_pin` 插了一枚 token ⇒ **它现在在 `:594`**（行号漂 1，内容零漂）。
尺＝`git show HEAD:scripts/portable-tests.sh | grep 'TestLiveWasapiSmoke'` 与工作树
`grep 'TestLiveWasapiSmoke' scripts/portable-tests.sh` 两枚产物 `cmp` ⇒ **完全一致（退码 0）**。
⇒ 那枚活着的陈旧性钉（`:651` `go test -list "^${name}\$" "$pkg"` 按 ledger 行自己的包跑）
本腿一枚字节都⛔ 动。★**行号锚已因本腿的 pin 插入而漂 1**：以后引用「`:593`」的人要读 `:594`
（票面／派单都按 `:593` 写；这是锚腐烂那一族，先例＝票 255 `AC#6`，⛔ 本票新造判据，只做具名）。

## 格式与语法

`bash -n scripts/portable-tests.sh` ⇒ **rc=0**（改动前基线也是 rc=0）。
换行属性尺＝`git check-attr text eol -- scripts/portable-tests.sh` ⇒ `text: set` / `eol: lf`
⇒ 本腿⛔ 动 `.gitattributes`、⛔ 动任何 git 配置，插进去的两行是 LF，与该文件既有 3415 行同形。
