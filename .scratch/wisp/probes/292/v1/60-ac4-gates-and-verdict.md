# 292-v1 60 — AC#4 门禁＋越界＋★编排者 C2 拆分的复核＋总判语表

## ① 四把门禁（每把自己落一行 rc，量于 HEAD `bccafc30` 之上的我的 commit 链）

| 门禁 | 逐字尺 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** | 无输出 |
| d22scan | `sh scripts/d22scan.sh` | **0** | 末行逐字 `d22scan: clean - no D22 ban violations`（live scope：bans #1-5 internal/=229、cmd/=39，ban #6/#7/#8 齐全，`cmd/` 扫 117 枚 Go） |
| gofmt（工作树） | `gofmt -l cmd/wisp` | **0** | **3 枚**：`cmd\wisp\models.go`、`cmd\wisp\panel_inbound_guards_35r3_test.go`、`cmd\wisp\panel_transport_35r2_test.go` |
| gofumpt（GOPATH bin） | `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` | **0** | 同 **3 枚**（逐字同上） |
| 裸 `gofumpt` | `command -v goumpt` 之类 | — |  PATH 里只有 `gofmt`（`/d/work/base/go/bin/gofmt`），gofumpt **不在 PATH** ⇒ 票面那条"裸调用 rc=127 不算跳过"我按 `$(go env GOPATH)/bin/gofumpt.exe` 真跑了，**没跳过** |

## ② ★C2 拆分复核（编排者那把 HEAD-blob 尺，我自己重跑）

尺＝`for f in panel_inbound_guards_35r3_test.go panel_transport_35r2_test.go models.go subagent_carrier_197_test.go; do git show HEAD:cmd/wisp/$f > /tmp/v1-head/$f; done; gofmt -l /tmp/v1-head`
（临时件落 **仓外** `/tmp/v1-head`，⛔ 不落仓根、⛔ 不删）

⇒ 输出**只有 2 枚**：`…\v1-head\panel_inbound_guards_35r3_test.go`、`…\v1-head\panel_transport_35r2_test.go`。
⇒ **`models.go` 在 HEAD blob 上干净**、`subagent_carrier_197_test.go` 也干净。
⇒ 我再补一把编排者没做的：**gofumpt 对同一批 HEAD blob** `$(go env GOPATH)/bin/gofumpt.exe -l /tmp/v1-head` ⇒ **同样只有那 2 枚** ⇒ 拆分对两把工具都成立，不是 gofmt 独有现象。

CRLF 假枚四把尺（逐字）：
- `git config core.autocrlf` ⇒ `true`
- `git ls-files --eol cmd/wisp/models.go` ⇒ `i/lf    w/crlf  attr/text eol=lf      	cmd/wisp/models.go`
- `tr -cd '\r' < cmd/wisp/models.go | wc -c` ⇒ **334**
- `git status --short cmd/wisp/models.go` ⇒ **空**
- 我另加一把：`git ls-files --eol cmd/wisp | grep -c "w/crlf"` ⇒ **1**（整个 `cmd/wisp` 只有 `models.go` 一枚是工作树 CRLF）

⇒ **C2 的拆分我复核＝成立**：真格式债 2 枚（票 35 族，归票 298），第 3 枚 `models.go` 是**本机工作树 CRLF 的假枚**（`attr/text eol=lf` 已声明入库为 LF、`git status` 对它零输出 ⇒ 它不在任何人的改动账上）。

## ③ 票面「`gofmt -l` 空」按字面能不能满足？

**不能**：HEAD 上就带着那 2 枚真债（不是我改出来的、也不是本票的改动带进来的）。
⇒ 按字面判＝**永远不可满足的判据**（票面缺陷，不是腿的缺陷）；编排者 C2 追加读作"新增 0 枚"且**原句不改**，与票面自己"引用前先重跑"的口径不冲突 ⇒ **裁：订正方式正确**（不偷改原句）。
⇒ 腿按"新增 0 枚"交回＝**正确处置**：尺＝改后 3 枚 ⊖ 改前（`4405e0f4^`）同 3 枚 ⇒ 新增 **0**；本改动只动注释行，且 `subagent_carrier_197_test.go` 自己两把工具下都**不在脏名册里**（见 ② 的 HEAD-blob 尺）。

## ④ 越界与名册（尺＝`git show --name-status` 三笔逐枚，整族穷举）

- `4405e0f4`：`M cmd/wisp/subagent_carrier_197_test.go` ＋ `A probes/292/r1/{10-comment-diff.md, logs/ac2-changed-lines.txt, logs/ac2-hunks.txt, msg-02-ac2.txt, run-pass.sh}` ⇒ **全在票面允许名册内**（`那枚件＋probes/292/**`）。
- `28292c05`：只有 `A probes/292/r1/{00-anchor-and-scales.md, msg-01-ac1.txt}` ⇒ 界内。
- `e390fdd9`：`M 票面` ＋ `A probes/292/r1/{20-runs.md, 30-gates.md, logs/gates.txt, msg-03-close.txt, who-owns-wisp.ps1}` ⇒ 界内；票面那一枚的 diff **只有一行新增**（`+- [2026-10-09 17:51:57 +08] agent=292-r1 did=…`）⇒ **没翻任何框、没动裁定节**。
- 冻结面：三笔里没有 `frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／`internal/**` 任何一枚 ⇒ **零字节，成立**。
- 票面四格在 HEAD 上仍是 `- [ ]` ×4（尺＝`grep -nE '^(- \[[ x]\])' 票面` ⇒ 22/23/24/25 行全 `[ ]`）。
- ⚠ 两处我要具名的"字面合规、精神存疑"：腿在 `probes/292/**` 里放了 `run-pass.sh` 与 `who-owns-wisp.ps1` 两枚**可执行脚本**。票面 AC#4 的名册按字面允许（`.scratch/**` 不在 d22scan/gofumpt 的 Go 分母里，我 ① 那两把尺的分母没被它们挪动）⇒ **不记越界**；但派单模板若统一"⛔ 不建 .sh/.ps1"，这两枚属**规约不一致**，归编排者定夺（不是 292 的账）。
- ⚠ 编排者 C4① 记腿的 `rm -f` 违反 `issues/README` 规则 8 ⇒ 我从盘上**无法复核**（文件已不在，`probes/292/r1/logs/` 现存名册里也没有仓根那 4 枚日志的影子）⇒ 该格判 **【不可判(欠：仓根那 4 枚日志的原件)】**，但那是**纪律账不是语义账**，与四格判语无关。

## ⑤ 与票面／裁定节的不符之处（我这一路的清单）

1. **票面现量 ⓒ 的尺口径**：票面写"尺＝`grep -rn "L1Windows:" --include=*.go .` ⇒ **1 命中**"。我两把都跑了：`grep -rn "L1Windows:" --include=*.go . | wc -l` ＝ **1**（票面这把复现 ✓），`grep -rn "L1Windows" --include=*.go . | wc -l`（**少个冒号**）＝ **26**。⇒ 票面结论对、**尺的写法必须连冒号一起抄**；照抄"去掉冒号"的人拿到的是 26 不是 1。⚠ 轻微、不改判语。
2. **票面 ⓑ 引的 `loop.go:603` 行号**：现量 `:603 func callCorr(...)` ✓ 未漂；但票面现量节说"用在 `:676`" ✓ 也未漂。**零处不符**。
3. **票面现量节引 roster 的 `:212/:214/:219/:220/:236`**：全部逐字对上 ✓。**零处不符**。
4. **`internal/panel/subagent_roster_197.go:180` 那节注释自己引 `gate.go:275, :471`**：我量到的 `orDefaultText` 调用在 **`gate.go:282`**（不是 275）⇒ **行号漂 7**，但那是**票 197/220 族那枚件的注释**，不在本票写面 ⇒ **不记 292 的账**，建议归口到下一枚碰那件的内务票。
5. **裁定节 C1**：说"我这二发不带 `-v` ⇒ `^--- PASS` 恒 0＝尺的口径"。我复核 `grep -cE '^\s*--- PASS' c1-post-1.txt`＝**0** ✓ 成立；红名册两发各 5 枚 ✓ 成立。**零处不符。**
6. **裁定节 C2**：拆分 ✓ 我独立复核成立（见 ②），并且我补的一点是**gofumpt 同形**。**零处不符。**
7. **裁定节 C3**（`:598-600` 三行内改写、允许 `:601` 承接）：腿确实没动 `:601`（单 hunk `@@ -598,3 +598,3 @@`＋改前后 `:601` 逐字相同）✓。**零处不符。**
8. **裁定节对 AC#3 的裁**（"本格现由我按读数裁为成立"）：⚠ **这一条我要具名保留意见**——票面 AC#3 自己写了"主证＝语义面零、名册作差＝辅证"，编排者按辅证读数裁成立**方向没错**；但"本格……现由我裁为成立"写在**翻框那一节之前**（同节末句又是"四格一枚不翻"），读起来像"读数已定案"。⇒ 我的处置：**按票面口径判 AC#3 成立（辅证＋主证都过）**，把"谁有权说'成立'二字"留给编排者落笔时定；**不改变判语**。

## ⑥ 四格总判语（我 = 292-v1，非实现者）

| AC | 判语 | 关键凭据（文件:行） | 还欠什么 |
|---|---|---|---|
| AC#1 | **成立** | 10 件：ⓐ`subagent_roster_197.go:214/:219-220/:236` ⓑ`loop.go:603-607`＋调用点唯一`:676` ⓒ`run.go` L1Windows 0 命中＋`subagent_blocked_220_test.go:81` ⓓ`internal/agent/*.go` ToolChoice 0 命中＋`tools/mockllm/chat.go:172` | 无（另记两节射程沉默：`indexWaitingKey:267-268`、`gate.go:282 orDefaultText`——都不是假话） |
| AC#2 | **成立** | 40 件：`git show --numstat`＝`3 3 cmd/wisp/subagent_carrier_197_test.go`；`-U0` 全量 6 行纯注释、非注释 0；单 hunk `@@ -598,3 +598,3 @@`；S1-S5 逐节真；`:601` 未动；断言零触及 | 无 |
| AC#3 | **成立** | 50 件：两发 382/5/2 与 383/4/2，红名册并集＝基线那 5 枚，新增红 0；两枚外因红在无探针态两发皆绿 | 无（主证在 AC#2；本格按票面口径是辅证） |
| AC#4 | **成立（带两笔归口别记）** | 60 件 ①②④：build/d22scan rc=0；gofmt/gofumpt 工作树各 3 枚、HEAD blob 各 2 枚 ⇒ 新增 0；`models.go` 假枚四把尺坐实；三笔名册全在允许面内、冻结面零字节、四格未翻 | 票面"`gofmt -l` 空"按字面**不可满足**（票面缺陷，编排者已追加订正）；C4① 那 4 枚被删日志＝**不可判**（纪律账，不影响本格） |

**总判**：本改动**没有把注释改成新的假话**；三行现在说的每一节我今天都能独立量到。⛔ 我不翻框、⛔ 不 push。

## ⑦ 没攻完的（诚实交底）

- 编排者 C4① 的"仓根 4 枚临时日志被 `rm -f`"＝原件不在盘上，我只能按裁定节自述读，**判不可判**。
- 我没做第 3、4 发改后整包（票面要求每态 ≥2 发，改前那两发用的是腿/编排者已有日志，我自己只跑了改后两发＝票面对"改后"的最低要求；⛔ 我没跑"改前"态的两发 ⇒ 改前态是**引用他人日志**，已在 50 件具名）。
- 我没去证"真 provider 的 callID 会不会是纯数字"（②/10 件 ⓐ 那节射程沉默的外部世界问题），只量了本仓内 mockllm 的 id 形状。
