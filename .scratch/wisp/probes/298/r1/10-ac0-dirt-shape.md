# 298-r1 — AC#0 逐枚把"脏在哪"抄出来（⚠ 全部量在 **HEAD blob** 那一把，⛔ 不是工作树）

## 我用的是 blob 那一把（逐字命令，下一位可原样重跑）

```
B=$(mktemp -d)                       # 实测 /tmp/tmp.XnATepEG0K ＝ C:\Users\swq\AppData\Local\Temp\tmp.XnATepEG0K
mkdir -p "$B/blob"
for p in cmd/wisp/models.go cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go; do
  git show "HEAD:$p" > "$B/blob/$(basename $p)"     # 写出来的是 blob 原文＝LF，不经 core.autocrlf 归一化
done
gofmt -l "$B/blob"                                                            # ⇒ 2 枚（models.go 不在列）
gofmt -d "$B/blob/panel_inbound_guards_35r3_test.go" > "$B/d/…35r3….diff"      # rc=0，1 hunk，14 行
gofmt -d "$B/blob/panel_transport_35r2_test.go"      > "$B/d/…35r2….diff"      # rc=0，1 hunk，14 行
```
尺读数留档＝`logs/../00-anchor.md` G5.2；两份 diff 原件已入库：
`logs/blob-gofmt-d-35r3-inbound-guards.diff.txt`（1207 字节）／`logs/blob-gofmt-d-35r2-transport.diff.txt`（746 字节）。
⚠ 下面每段 diff 都带 blob 那把尺的落点行号（`@@ -63,8 +63,8 @@`／`@@ -804,8 +804,8 @@` 是相对 blob 的号）。

---

## 第 1 枚 `cmd/wisp/panel_inbound_guards_35r3_test.go`

- **hunk 枚数＝1**（`@@ -63,8 +63,8 @@`，落在 `const (` 块内），**改动行数＝2 换 2**，`git diff -U0` 实测 numstat `2 2`。
- **改动性质＝注释列对齐（尾随 `//` 注释的前导空格数），⛔ 零语义。**
  根因＝gofmt 的对齐算法按 **rune 计数**补齐列，而这三行原文是**按 CJK 显示宽度**手打的：
  反引号里是中文文案，多打的那 2/3 个空格在 gofmt 眼里就是"多出来了"。
- **落点形状逐段引（blob 尺的 `gofmt -d` 原文，未删一字）**：

```
@@ -63,8 +63,8 @@
 // cases may assert the same sentence (:63's "三支不许共用一句文案").
 const (
 	guard35r3Roster    = `不是面板 composer 通路的能力入口` // bridge.go:133 (name %q ... 不是面板 composer 通路的能力入口)
-	guard35r3Source    = `按伪造/串台拒绝`               // bridge.go:136-137 (来源 %q 不是 %q，按伪造/串台拒绝...)
-	guard35r3RequestID = `缺少 requestId`                // bridge.go:140-141 (缺少 requestId，无法与审计/卡片对齐...)
+	guard35r3Source    = `按伪造/串台拒绝`              // bridge.go:136-137 (来源 %q 不是 %q，按伪造/串台拒绝...)
+	guard35r3RequestID = `缺少 requestId`          // bridge.go:140-141 (缺少 requestId，无法与审计/卡片对齐...)
 )
```
- 逐字核对：三枚标识符名、三枚反引号字符串、三枚注释正文**一字未动**；只有 `//` 之前的空格数变（`:136-137` 那行少 2 个、`:140-141` 那行少 6 个）。
  ⇒ 注意 gofmt 把 `-` 那行换成 `+` 那行之后，第 3 行的列位置是**跟着最短那行对齐**的（gofmt 以块内最宽 rune 行为准），
  所以第 3 行的 `//` 视觉上离文案更近了——这正是"显示宽度 ≠ rune 数"的可见后果，仍⛔ 不是语义。
- **引入首笔（尺＝`git log --oneline -- cmd/wisp/panel_inbound_guards_35r3_test.go`，名册只 1 笔 ⇒ 首笔＝唯一一笔）**：
  `7f9d6e40` 「35-r3 交件（票35 :63 (d) 三支入向守卫红，只写测试）：新增 cmd/wisp/panel_inbound_guards_35r3_test.go（//go:build windows）…」
  ⇒ 带进来的落地腿＝**票 35 的 `35-r3`**；这枚件**自创建起就是 gofmt 脏的**（同尺在首笔 blob 上同样命中，见下"首笔即脏"证据）。

## 第 2 枚 `cmd/wisp/panel_transport_35r2_test.go`

- **hunk 枚数＝1**（`@@ -804,8 +804,8 @@`，落在 `type jsParser struct` 之后的两枚单行方法），**改动行数＝2 换 2**，numstat `2 2`。
- **改动性质＝单行函数体的花括号前空格对齐（多余空格），⛔ 零语义。**
- **落点形状逐段引**：

```
@@ -804,8 +804,8 @@
 	i    int
 }
 
-func (p *jsParser) peek() jsTok    { return p.toks[p.i] }
-func (p *jsParser) next() jsTok    { t := p.toks[p.i]; p.i++; return t }
+func (p *jsParser) peek() jsTok { return p.toks[p.i] }
+func (p *jsParser) next() jsTok { t := p.toks[p.i]; p.i++; return t }
 func (p *jsParser) atPunct(s string) bool {
```
- 逐字核对：接收者、方法名 `peek`/`next`、返回类型 `jsTok`、函数体 `return p.toks[p.i]` 与
  `t := p.toks[p.i]; p.i++; return t` **一字未动**；只有 `{` 之前的空格从 4 个变 1 个。
- **引入首笔（尺＝`git log --oneline -- cmd/wisp/panel_transport_35r2_test.go`，名册 3 笔，新→旧 `3a343bc7`／`2fc5f5c9`／`286a7f30` ⇒ 首笔＝最老那笔）**：
  `286a7f30` 「35-r2 甲③ 落地＋判据换行为尺（票35 AC#6:52）… 新增 cmd/wisp/panel_transport_35r2_test.go＝内置 JS 子集解释器逐字执行…」
  ⇒ 带进来的落地腿＝**票 35 的 `35-r2`**。
- ⚠ **同一枚件上还有第二笔已知遗产**（本票⛔ 射程，具名留档不处理）：`2fc5f5c9` 的提交正文自己写着
  「gofmt 那 4 行漂移 HEAD 同处也在（`:775-776`→`:807-808` 位移，CR 两版皆 0）⇒ 本笔不新增未格式化面、那 4 行归 35-r2/票 275 族⛔ 不顺手改」。
  ⇒ 本票洗掉的**就是那 4 行**（现在的落点＝`:807-808`，本笔改后为 `:807-808` 的对齐形），
  也就是说这条债在票 35 的账上**已被具名挂过、一直没人认领**，本票是它的第一次真正清偿。

## "首笔即脏"（⛔ 不是推测，是已跑的尺；读数原件＝`logs/first-commit-blob-rulers.txt`）

尺＝对**首笔那棵树里的 blob** 再跑一次 blob 那把尺（逐字）：

```
N=$(mktemp -d); mkdir -p "$N/first"
git show "7f9d6e40:cmd/wisp/panel_inbound_guards_35r3_test.go" > "$N/first/panel_inbound_guards_35r3_test.go"
git show "286a7f30:cmd/wisp/panel_transport_35r2_test.go"      > "$N/first/panel_transport_35r2_test.go"
gofmt -l "$N/first"          # ⇒ 2 枚全命中，rc=0 count=2（落点＝仓外 $N/first，⛔ 不落仓内）
gofmt -d "$N/first"/*        # ⇒ 各 1 hunk：35r3 @@ -63,8 +63,8 @@ ／ 35r2 @@ -772,8 +772,8 @@
```

- **35r3**：首笔 blob 与本腿改前的 HEAD blob **逐字节 SAME**（9759 字节）⇒ `7f9d6e40` 创建那一刻就是脏的，此后无人动过这枚件。
- **35r2**：首笔 blob 脏形落在 `@@ -772,8 +772,8 @@`（首笔 1824 行／52759 字节），改前 HEAD 同一脏形漂到 `@@ -804,8 +804,8 @@`（2142 行／73476 字节），**位移 +32 行**。
  ⇒ 与 `2fc5f5c9` 提交正文自己记的那句**逐字对得上**：
  「gofmt 那 4 行漂移 HEAD 同处也在（`:775-776`→`:807-808` 位移，CR 两版皆 0）⇒ 本笔不新增未格式化面、那 4 行归 35-r2/票 275 族⛔ 不顺手改」。
  ⇒ **归因链闭合**：脏形自 `286a7f30`（35-r2 首笔）带入，后续两笔只搬位置⛔ 未修；本票 `298-r1` 是它的第一次真修。

⇒ 两枚**不是后来被别的腿改脏的**，是**创建那一笔就带着 gofmt 脏进 HEAD**；
配合票面现量节那句「本仓 CI 那两步 `go vet`／`gofmt denominator` 长期是 skipped＝从未求值」，
"落地时没人看见 gofmt 红"这件事在 CI 面上有解释 ⇒ ⛔ 不构成某位落地腿明知故犯。
