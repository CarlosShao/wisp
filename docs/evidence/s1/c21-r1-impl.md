# 阶段 0-1 件三 · C21 夹具腿 c21-r1 · 实施证据

`s1/` 证据 · 写腿（实现者）· 2026-10-26 · 基线 `365d5d547`。裁定依据 `A1246⑦`：CSS 内容做成**测试夹具**进仓，测试读夹具（hermetic）。

## 0. 开工前核验（逐条现跑）

| 项 | 命令 | 读数 |
|---|---|---|
| 索引空 | `git diff --cached --name-status` | 空输出 ⇒ `INDEX_EMPTY=YES` |
| ① `HEAD` 无该件 | `git cat-file -e HEAD:design/assets/tokens.css` | `fatal: path 'design/assets/tokens.css' does not exist in 'HEAD'` |
| ① 删除提交 | `git show --name-status 72b65c614` | `D design/assets/tokens.css` |
| ② 盘上有 | `design/old/assets/tokens.css` | 存在，`12942` B |
| ② 就是被删那份 | 它的 sha256 vs `git show 72b65c614^:design/assets/tokens.css` | 同为 `8ee27dc5…5517`（**逐字节相同**） |
| ③ 未跟踪 | `git ls-files design/old` | `COUNT=0` |

三条结论**全部成立**。补充：`gen-tokens.mjs` 已不再读 `design/**`（见读数 5）⇒ 派生件 `tokens.generated.css` 不能当被测件。
读全文：`internal/ball/tokens_table_test.go` 共 `1571` 行；读路径在 `c21ParseTokensCSS`（1381–1436），常量 `c21TokensCSSPath` 在 1364 行。

## 1. 件一：夹具（新建）

`internal/ball/testdata/tokens_c21.css` —— `1450 B` 头注释 ＋ `12942 B` 源内容 ＝ `14392 B`，纯 LF，**零 emoji**。

头注释四项齐全：**出处**（`design/old/assets/tokens.css`）、**复制日期**（2026-10-26）、
**为什么进仓**（两个 C21 测试必须 hermetic；`design/` 非本树写面；读机主盘上文件的测试在干净 checkout 上永远红，
改成 SKIP 等于什么都没测）、**什么时候该更新**（规范 C21 tokens.css 变更、或 C21 表合法随之变更时；
替换头注释以下全部字节并重跑两个测试；夹具过期就该红，那不是削弱断言的理由）。

## 2. 件二：改测试读夹具（只改路径一行）

`internal/ball/tokens_table_test.go:1364`
`const c21TokensCSSPath = "design/assets/tokens.css"` ⇒ `"internal/ball/testdata/tokens_c21.css"`

断言、解析逻辑、`must never skip` 的 `t.Fatalf` 全部原样。

## 3. 读数

### 读数 1：目标测试 rc（改前 / 改后）

改前 `rc=1`：
```
=== RUN   TestC21TableColourRowsMatchTokensCSS
    tokens_table_test.go:1468: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified. - the CSS leg of this check must never skip
--- FAIL: TestC21TableColourRowsMatchTokensCSS (0.00s)
FAIL	github.com/CarlosShao/wisp/internal/ball	0.029s
```

改后 `rc=0`：
```
=== RUN   TestC21TableColourRowsMatchTokensCSS
    tokens_table_test.go:1530: three-way colour check: 78/78 cited rows agree across docs/evidence/s1/c21-native-tokens.md, internal/ball/testdata/tokens_c21.css and internal/ball/tokens.go (:root=129 decls, light=65)
--- PASS: TestC21TableColourRowsMatchTokensCSS (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/ball	0.026s
```
`78/78` 与 `:root=129 ＋ light=65` 同时对上 ⇒ 夹具被**真的**读进去并逐行对账。

### 读数 2：改前 `t.Fatalf` 逐字（证明原来真读不到）

```
tokens_table_test.go:1468: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified. - the CSS leg of this check must never skip
```
（`t.Fatalf` 在 1386 行，报错位置是调用点 1468 行。）

### 读数 3：夹具与源文件 sha256

```
402e32…ae496  internal/ball/testdata/tokens_c21.css      <- 整文件（含头注释）
8ee27d…5517   design/old/assets/tokens.css              <- 源
```
整文件哈希**必然不同**：件一强制要求加头注释。故按「去掉头注释后逐字节相同」证明：

```
tail -c +1451 internal/ball/testdata/tokens_c21.css > body
8ee27dc59bc3106901a971b533ed0f1e4e65ac7a62ce80c9ef4ea31646558517  body
8ee27dc59bc3106901a971b533ed0f1e4e65ac7a62ce80c9ef4ea31646558517  design/old/assets/tokens.css
cmp body design/old/assets/tokens.css            ⇒ CMP_BODY_EQ=YES
cmp fixture 与 (header ++ source) 的拼接件      ⇒ CMP_CONCAT_EQ=YES
```
源文件本身 `==` 被删的 `HEAD` 版本（同 `8ee27d…5517`）⇒ 夹具抄的是真件，不是重打的近似值。

### 读数 4：断言零 diff

```
git diff --numstat -- internal/ball/tokens_table_test.go   ⇒ 1	1
git diff -U0 | grep '^@@'                                  ⇒ @@ -1364 +1364 @@
+/- 行数                                                   ⇒ 1 / 1
-const c21TokensCSSPath = "<design/assets/tokens.css>"
+const c21TokensCSSPath = "<internal/ball/testdata/tokens_c21.css>"
总行数                                                     ⇒ 改前 1571 / 改后 1571
```
单一 hunk、单加单减、总行数不变 ⇒ 断言段落（`c21ParseTokensCSS`、78 行对账、`--font-sans` 那条、
`checked == 0` 兜底）**零 diff**；`must never skip` 的 `t.Fatalf` 未动，未加任何 SKIP。

### 读数 5：`internal/panel` 判语 —— ★**同症状、不同根因 ⇒ 停手上报（未改）**

改前也红在同一行：
```
tokens_fourway_test.go:441: read design/assets/tokens.css: open ...\design\assets\tokens.css: The system cannot find the path specified.
```
缺件原因完全相同（同一常量、同一被删文件）。**但 P1 指到夹具不能让它绿，且绿不了的原因不是文件缺失**：

- P4 已不从 `design/**` 生成。`frontend/scripts/gen-tokens.mjs` 自述 `FOURTH GENERATION (owner pivot, 2026-09-25)`，
  输入是 `frontend/src` ＋ beautiful-ui 库主题（`TurboKach/ai-native-react-components` 的 `app/globals.css`，
  upstream `05dab2d2`）；全文**没有任何** `design/` 读取路径。
- 生成件结构已变：`:root` 43 条（含派生 `-color` 共 48）、`[data-theme="light"]` **0 条**（只有 `[data-theme="dark"]`），
  而四条腿的 `parseTokenBlocks` 只认 `:root` 与 `[data-theme="light"]`。
- 实测把 P1 指向夹具：parties 行 `135 dark decls / 78 rows / 40-38 fields / 48 dark decls`，
  `does not carry` 报 **270** 条。没有任何仓内文件能同时满足 `P2→P1`（C21 表记旧世代值，读数 1 已证 `78/78`
  只对得上旧件）与 `P1→P4`（新库主题）⇒ 这条测试**结构性不可满足**。
- 世代切换出处：`cc9ecf9ab`「真相源从 design/assets/tokens.css 换成 design/doubao/demo/styles.css（owner
  2026-09-25 原话…旧的那套他移进 design/old/ 就是不要了）」；`424ac84d9`「demo 路线废弃，收编库主题为 token 第四代」。

⇒ 让它绿只有两条路：**改 `frontend/**`**（本腿零写面）或**改这条测试的断言/删腿**（硬边界 2 明禁）。
故按「根因不同 ⇒ 停手上报」处理：`internal/panel/tokens_fourway_test.go` **原样未改**
（`git diff --quiet` 对它 `rc=0`，与 `HEAD` 逐字节相同；整包失败集合改前/改后同一 5 枚）。
**这个决定留给编排者/机主**：废弃 P1→P4 这条腿，还是让前端重新对齐 C21。

### 读数 6：双向验（证明它还在真的测东西）

```
把夹具改名到 D:\tmp\wisp\20261026\c21r1\fixture-hidden.bak
rc=1  tokens_table_test.go:1468: read internal/ball/testdata/tokens_c21.css:
      open ...tokens_c21.css: The system cannot find the file specified. - the CSS leg of this check must never skip
改名还原
rc=0  ok github.com/CarlosShao/wisp/internal/ball  0.040s
```
夹具缺失 ⇒ 红（**硬失败**，不是 SKIP）；还原 ⇒ 绿。

### 读数 7：四数（整包，不相减）

| 包 | 改前 rc | 改后 rc | 改后 PASS/FAIL/SKIP |
|---|---|---|---|
| `internal/ball` | `1` | `0` | 69 / 0 / 0 |
| `internal/panel` | `1` | `1` | 166 / 5 / 0 |

`internal/panel` 的 5 枚失败改前改后**同一集合**（`TestC21DesignTokensFourWayAgree` ＋ 4 枚既有无关红灯：
approval 键、composer 键、backpressure 两条）⇒ 本腿既未使 panel 变坏，也未擅自修它。

### 读数 8：格式甲乙两尺

```
尺甲  gofmt -l internal/ball internal/panel                                ⇒ 空输出，rc=0
尺乙  D:\work\base\gopath\bin\gofumpt.exe -l internal/ball internal/panel  ⇒ 空输出，rc=0
      D:\work\base\gopath\bin\gofumpt.exe -d internal/ball/tokens_table_test.go ⇒ 空输出，rc=0
      gofumpt v0.12.0 (go1.27.1)
```
（`gofumpt` 用绝对路径；PowerShell 下须写 `D:\…` 而非 `/d/…`。）

## 4. 写面自查 ＋ 停手

只动：`internal/ball/testdata/tokens_c21.css`（新建）、`internal/ball/tokens_table_test.go`（1 行）、本文件。
未碰 `design/**`（只读 `design/old/assets/tokens.css` 一个文件）、`frontend/**`、`internal/ball/tokens.go`、
`internal/panel` 非测试件、`cmd/**`、`.github/**`、`scripts/**`、票面/台账/`docs/specs/**`/`docs/PLAN.md`/`docs/audit/**`。
未加 SKIP、未改断言、未 push。

**停手一件**：读数 5 的 `internal/panel/tokens_fourway_test.go` 不属本腿可修（绿需改 `frontend/**` 或改断言），已原样留回，等裁定。
