# 141 — `AC#3` / `AC#4` 补发的正向对照（编排者现量，快照在仓库外）

**时刻**：2026-09-25 07:21–07:24 `+0800`（每条读数旁边另标时刻）
**被验版本**：`471af50`（07:17 那次净扫描）与 `67192a5`（07:23 之后）—— 两支不同，已分别点名
**作者身份**：编排者。**本文件不是裁决表**——票 141 的非实现者裁决在 `141-q46c-accept-r1.md`（总裁＝附条件成立）。
本文件只补 `141-…md` 里那两处**"没人量过"**的读数：`AC#3` 示例那句到底响不响、`AC#4` 的
`frontend/` 与 `design/` 两条腿今天到底响不响。

**为什么这些读数能在快照里造而不破坏任何地界**：种子全部种在 `git archive` 解出来的
仓库外副本（`/d/tmp/d22-*`），**没有往 `frontend/**`、`design/**`、`internal/**` 真树里写一个字节**。

---

## 1. 门禁四数（`AC#4` 的"纯净快照"那一半）

```
$ git archive HEAD | tar -x -C /d/tmp/d22-clean-471af50     # 07:17，HEAD=471af50
$ cd /d/tmp/d22-clean-471af50 && sh scripts/d22scan.sh
runtests.sh: OK - packages=[./...] top-level: PASS=24 FAIL=0 SKIP=0, === RUN=64, '[no tests to run]'=0
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope bans #1-5 cmd/           examined  22 production Go files
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #7 internal/tools/   examined  18 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  40 text files
d22scan: scope ban #8 internal/         examined 405 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  39 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; ...
rc=0
```

**"各 scope 命中数不降"现在两表都有了**：本组 `design/16 frontend/40 internal/405 cmd/39`
与 `141-q46c-accept-r1.md:38-48`（`git archive bb61dc5`）逐枚同值 ⇒ 相对基线**零变化**。

**同一版扫描器、同一枚 HEAD，工作树 vs 净快照的 `frontend/` 枚数**：`43`（工作树，含被 gitignore 的
`frontend/dist/` 3 枚）vs `40`（净快照）。这是 `A207` 那枚膨胀的第三个读数，**不改判据，只是再钉实一次**。

---

## 2. `frontend/` 这条腿有没有分母（两形各一枚）

| # | 种子形状 | 落点 | rc | 仪器的话 |
|---|---|---|---|---|
| **A** | `≤`(U+2264) 放进 **`//` 行注释** | `frontend/src/App.tsx` 末行 | **0** | 无 finding ⇒ **文本 scope 的注释豁免是真的**（Q-46(c) 在 `.tsx` 上也生效） |
| **B** | `≤`(U+2264) 放进 **JSX 文本节点** `<div>≤</div>` | 同上文件 | **1** | `frontend/src/App.tsx:64: [emoji] ban #8 glyph in scope frontend/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)` |

〔形状交代〕A/B **同一枚快照目录**（`/d/tmp/d22-clean-471af50`）、**B 接在 A 之后累加**（那枚 `//` 行还在）。
所以严格说 B 读的是"注释行＋文本节点同时在"⇒ 红因只能归给文本节点（注释那一发在 A 里已单独证明不响）。
两发之间**没有重建快照**，这一点写在这里，别让下一位以为是两枚独立快照。

⇒ **`AC#4` 里"`frontend/` 那一支今天响不响"这一发读数从此存在**（此前 141 三张表零处记录）。
答：**响**，且只在非注释位置响。

## 3. `design/` 这条腿有没有分母

| # | 种子形状 | 落点 | rc | 仪器的话 |
|---|---|---|---|---|
| **C** | `≤`(U+2264) 单独一行（非 `<!-- -->`） | `design/screens/ball.html` | **1** | `design/screens/ball.html:580: [emoji] ban #8 glyph in scope design/ is banned (D23): ...` |

## 4. `AC#3` 示例那句到底响不响（这是本文件最要紧的一格）

`AC#3` 原文要的是："造一发**只含 `U+2190–U+25FF`、不含任何现有被扫字形**的种子（例如一枚 `→` 放进注释）"。
两枚快照、同一落点 `internal/secret/redact.go` 末行、只差一个字形：

```
$ printf 'var probe141Arrow = "\xe2\x86\x92"\n' >> <快照>/internal/secret/redact.go   # U+2192
$ go run . -root <快照>            # 07:24, HEAD=67192a5
→ redact.go **零 finding**（那一趟的 rc=1 只来自 §5 那枚 frontend/ 新红，与本种子无关）

$ printf 'var probe141Ge = "\xe2\x89\xa5"\n' >> <快照>/internal/secret/redact.go      # U+2265
$ go run . -root <快照>            # 07:24, 同一 HEAD
internal/secret/redact.go:14: [emoji] ban #8 glyph in scope internal/ is banned (D23): ...
```

**判**：
- `→`(U+2192) 在**字符串**里都不响 ⇒ `AC#3` 的示例是**双重不可满足**：既落在刻意不扫的箭头段、又被注释豁免覆盖。
  **"照字面执行这一发"永远拿不到红**，这不是取证没做，是判据结构上产不出它自己想量的读数。
- `≥`(U+2265) 在字符串里响 ⇒ 交付那一版**确实有会响的第五段种子**，且这一发不必靠读代码推断，今天现场复算过。
- 附带一条（原 `U8` 的形状）：先前一次误投种子到不存在的 `internal/config/config.go`（`>>` 当场创建了它），
  那枚文件同时报 `[unparseable]` **和** `[emoji]` 两条 finding ⇒ **go/ast 解析失败不会吞掉 ban #8**，
  它只会**另加**一条 `unparseable`。（这条纠正了我自己写在记忆里的"解析失败＝无豁免"那个含糊说法：真形状是"两条并报"。）

---

## 5. 本轮顺带量到的一枚新红（归因要说清，别记到我头上）

```
frontend/VENDORED.md:149: [emoji] ban #8 glyph in scope frontend/ is banned (D23)
```
- 字形＝`⚠`(**U+26A0**，落在**一直就扫**的 `2600–27BF`，**与我 `3c80352` 扩的 `2200–22FF` 无关**)。
- 引入者＝前端会话自己的 commit **`67192a5`**（07:23，U2 选 A 案那一枚）。
- 时间线：`07:17` 净快照 rc=0（`471af50`）→ `07:18` 台账 `d363e26` → `07:23` `67192a5` 带进 `⚠` → 现在红。
- 处置＝已具名发回前端会话（那是 owner 交出去的地界，我不进去改），要求**只改这一枚字符**并回 sha。
- ⇒ **§1 那枚 rc=0 只对 `471af50` 负责**，不许被读成"当前 HEAD 绿"。

## 6. 〔取证口径〕这一档里哪些是真跑、哪些不是

- **真跑**：§1 净快照四数＋rc；§2 A/B；§3 C；§4 两发种子对照。命令全部随文，快照目录在 `/d/tmp/d22-*`（只建不删）。
- **不是真跑**：`U+2200–U+22FF` 段内 256 枚码点里**只有 4 枚有会响的种子**（`scan_test.go:357-360`），
  其余 252 枚**靠类算术**，没有任何一发变异打过它们。这一条不是本文件能补的，属**已知残留**。
