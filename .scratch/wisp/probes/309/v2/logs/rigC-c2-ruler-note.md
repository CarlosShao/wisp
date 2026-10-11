# `rigC.txt` 的 C2 那一发：尺口径自抓（裁决者⛔ 改原件，只追加这一枚更正件）

`logs/rigC-controls.log` 里 `=== C2: how many controls does the nail run at HEAD? ===` 那一行下面记的是 **0**。
**那枚 0 是尺错、⛔ 读数。** 成因逐字：脚本里写的是 `grep -c '\t\t{"'`，`\t` 在 GNU grep 的基本式里⛔ 展开成制表符
⇒ 名册永远 0 命中。同一发里的 `C2b`（`grep -n 'for _, c := range controls'`）⛔ 带转义符、命中真身＝钉的 `:359`＝那一行如实。
另：C2 那一段 python 只做了改名（`controls := []struct {` → `var controlsUnused = []struct {`）并且**从没调用过 `run`**
⇒ ⛔ 任何测试读数建立在那次突变上，`trap` 的 `restore` 当场把它还原（该文件末尾那行 `RESTORE cmp … match=YES` 就是它）。

## 换一把 tab-safe 的尺重量（工作树面＝HEAD blob 面，`git status --porcelain -- internal cmd`＝0 行）

```sh
awk '/controls := \[\]struct/,/^\t\}$/' internal/tools/toolrequest_nail_309_test.go | grep -c '^[[:space:]]*{'
```

读数＝**6**（＝控制名册六条：shape 1 限定名字面量／shape 2 声明＋逐字段赋值／shape 3 同文件别名／shape 1b 同包裸名／
guard 非 agent 包里的同名裸字面量 want 0／"⛔ 双计"那枚 `var rig = agent.ToolRequest{…}` want 1）。
该读数连同命令一起记在 `logs/s9-behaviour-invariance.log` 的头两行。

## 这条更正影响哪儿

只影响 `10-rigs-and-readings.md` 的 C2 那一行（那里写的就是 **6**，⛔ 引用那枚 0）。
`AC#1`／`AC#2`／`AC#3` 三格判语⛔ 用那枚 0 当凭据：判"控制⛔ 恒真"的是 **B11**（打瞎⇒红在 `control "shape 1…"`）
与 **C1**（catch-all⇒红在 `control "guard…"`）两发实测，两发都带 pristine 正控（B0／C0＝绿）。
