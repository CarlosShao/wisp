# 216 AC#5 · 两枚无牙尺装牙（`w216-fix-rulers`）

日期：2026-11-14 ／ 写面：`internal/panel/error_copy_216.go`、`internal/panel/error_copy_216_test.go` ／ 触发＝编排者突变取证（ticket 216 AC#5 re-run，见票尾 2026-11-14 前一条与 `216-ac5-rulers-audit.md`）。

## 1. 被证伪的事实（继承自编排者，本腿复核成立）

`internal/observe/errors.go` 的 `func (e *Error) Error() string` 第一句就是 `fmt.Fprintf(&b, "%s", e.Class)` ⇒ 类名在 `err.Error()` 里。于是旧尺两处对表内容不敏感：

- 旧 `:49` `want := errorCopy216[c] + err.Error()`：期望值与产码同源（同一张表），改表两边同变 ⇒ 同义反复。
- 旧互异尺比 `sentence + err.Error()`：`err.Error()` 含类名，每类必异 ⇒ 两格同句也比不出。

编排者突变（`ClassCancelled` 句改成 `ClassConfig` 句）旧尺 `ok 0.102s` rc=0，该红没红。

## 2. 改动

### 2.1 产码（`error_copy_216.go`）

- 新增 `errorHead216(err error) string`：只返分类映射头句；`ClassOf` 不中或表缺行返空串；`err == nil` 返 `""`。
- `errorSentence216` 改为 `errorHead216(err) + err.Error()`（nil 短路返 `""`）。行为与旧形逐字节一致：有头接头，没头原样 `err.Error()`；未分类路径 `"" + err.Error() == err.Error()`。
- `userFacingErrorSentence216` 一字未动。表 `errorCopy216` 本体、`errorCopy216Complete` 未动。

### 2.2 测试（`error_copy_216_test.go`）

- 新增 `wantHead216`：**测试自己的** 17 条字面量期望头句（逐类一行，Go string literal，⛔ 不从 `errorCopy216` 取）。产码表与测试表是两份独立工件：改产码句不同步改测试 ⇒ 必红。
- `TestErrorCodeSelectsTheSentence216`：逐类断言 `errorHead216(observe.New(c, sameDetail)) == wantHead216[c]`（literal 对照）；另保留整句对照 `errorSentence216(err) == wantHead216[c] + err.Error()`；`errorCopy216Complete` 穷尽尺保留。
- `TestErrorCopySentencesAreAllDistinct216`：改比 17 个 `errorHead216` 两两互异，**并加前缀自检**：任一头句是另一头句的前缀（或相同）即红——堵「整句含类名蒙混」与「头句互为前缀靠原因缀凑异」两形。
- 其余测试（未分类透传、RefusedEnvelope wire）未动。

## 3. 正控：重打编排者那发突变（可编译）

突变＝`error_copy_216.go` 表中 `observe.ClassCancelled` 的句子字面量改成 `observe.ClassConfig` 那句（只改 string 字面量，编译通过）。`go test -count=1 ./internal/panel -run 216` ⇒ **rc=1，两尺都红**，红句逐字：

**表测（literal 对照）红：**

```
--- FAIL: TestErrorCodeSelectsTheSentence216 (0.00s)
    error_copy_216_test.go:80: class cancelled: head = "配置这一步出了问题，先看下面这句具体原因：", want literal "这一步被取消了，所以没有继续："
    error_copy_216_test.go:85: class cancelled: sentence = "配置这一步出了问题，先看下面这句具体原因：cancelled: upstream 413 payload too large", want head+cause "这一步被取消了，所以没有继续：cancelled: upstream 413 payload too large"
```

**互异尺（头句两两＋前缀自检）红：**

```
--- FAIL: TestErrorCopySentencesAreAllDistinct216 (0.00s)
    error_copy_216_test.go:134: classes config and cancelled produce the same head "配置这一步出了问题，先看下面这句具体原因：" - the branch no longer separates codes
    error_copy_216_test.go:137: head of cancelled is a prefix of head of config ("配置这一步出了问题，先看下面这句具体原因：") - the heads do not separate on their own
    error_copy_216_test.go:140: head of config is a prefix of head of cancelled ("配置这一步出了问题，先看下面这句具体原因：") - the heads do not separate on their own
```

（行号 80/85/134/137/140 为突变时刻行号；同一突变打在旧尺上曾得 `ok 0.102s` rc=0——本腿交付后同突变 rc=1，牙已装上。）

## 4. 还原与复绿

- 逐字节还原突变行后：`error_copy_216.go` SHA256 ＝ `90E6436F4A2D06E9BDE45B1BEE878FE5D976A82B4D05998395CE608AAF3C3ADF`（＝改后未突变态；表内容行 `git diff` 0 行）。
- 复绿：`go test -count=1 ./internal/panel -run 216` ⇒ `ok 0.128s` rc=0。

## 5. 门禁四数

| 门 | 读数 |
|---|---|
| `go test -count=1 ./internal/panel` | `ok 23.923s` rc=0；`-v`：PASS=233 / FAIL=0 / SKIP=0（顶层 ===RUN 349 含子测试，无 FAIL 无 SKIP） |
| `go vet ./internal/panel` | rc=0 |
| `sh scripts/d22scan.sh` | rc=0，`d22scan: clean - no D22 ban violations` |
| `gofumpt -l internal/panel` | 空列表（`go run mvdan.cc/gofumpt@latest`） |

## 6. 边界与残余

- 未动：`bridge.go`、`cmd/wisp/**`、`scripts/**`、`frontend/**`、`models/**`、阈值/golden/契约；未放宽任何断言；期望值不再取自产码表。
- `errorHead216` 对「已分类但表缺行」返空串——该形被 `errorCopy216Complete` 尺钉住不可达。
- AC#5 框的翻勾仍归非实现者裁决腿；本腿只交实现＋正控证据。
