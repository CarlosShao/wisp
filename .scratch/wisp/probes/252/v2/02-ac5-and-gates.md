# 252-v2 AC#5 复读（dump 等级行）＋门禁抽核

## diff 事实（`git show d3cccfda`）

- 单文件 `cmd/wisp/ticket224_assembly_test.go`，`--numstat`＝`10 1`，**单 hunk**（`@@ -205,11 +205,20 @@` 一处）。全部改动＝`grantLinesOf` 上方 +9 行注释＋滤网条件加一枚针 `strings.Contains(l, "tools: call ")`；滤网函数之外零改动，非测试文件零改动（产码语义零动）。
- 现读函数体（＝`d3cccfda` blob）：条件行＝`strings.Contains(l, "REPLY ") || strings.Contains(l, "tools: call ")`。
- 4 处调用点（`:148`／`:192`／`:195`／`:333`）全部把返回值当 `%s` 参数写进 `t.Errorf`/`t.Fatalf` 的消息 ⇒ 加针只改 dump 内容，**无任何断言被放宽**（没有任何断言以 dump 字符串为成立条件）。

## 改前（旧滤网）grep 计数 0

- 发＝W1 改前红（见 01 件）；dump 段逐字只有两行（`[audit] wisp run: SESSION-MINT …`＋`[audit] session: GRANT-RECORD …`）；全日志 `grep -c "tools: call"`＝**0**、`grep -c in_allowlist_scope`＝**0**。全文在 `logs/20-ac5-dump-verbatim.txt`。

## 改后（本腿自造同形红，overlay）

- 发：`go test -count=1 -v -overlay='C:\tmp\wisp252v2\overlay-ac5-v2.json' -run '^TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking$' ./cmd/wisp/`（在 `tree-before` 上，短名 TMP；12:37:04→12:38:36；rc=1；`--- FAIL: … (83.11s)`）。
- 同一处 dump（交付行号 `:333`）逐字含两枚 `tools: call`（ask 与 hit），样行：
  `[audit] tools: call kind=approval-timeout task=t224r2-read corr=t224r2-ask tool=fs.write risk=L2 decision=timeout outcome=error rules_hit=[R1 R2] in_allowlist_scope=false grant_id=0 reason="R1: 工具声明为下界（L1）; R2: 目标路径在授权目录之外: C:\\tmp\\wisp252v2\\SHORT-~1\\…\\never-granted.txt"`
- 全日志 `grep -c "tools: call"`＝2、`in_allowlist_scope`＝2（改前 0/0）⇒ 判据"dump 那行等级"成立。

## 还原证（sha256；overlay 零写盘）

- 工作树 `cmd/wisp/ticket224_assembly_test.go`：发前＝`366da8f5…`、发后＝`366da8f5…`（同）；`tree-before` 同件：发前＝`8ac227a8…`、发后＝`8ac227a8…`（同）。overlay 仅构建期替换、零字节落盘；本腿零 Go 改动。

## 门禁（本腿现跑，工作树）

| 门 | 命令 | 起止（+08） | rc | 读数 |
|---|---|---|---|---|
| build | `go build ./...` | 12:42:18→12:42:22 | 0 | 0 B |
| d22scan | `sh scripts/d22scan.sh` | 12:42:29→12:42:56 | 0 | `d22scan: clean - no D22 ban violations`（23,758 B） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 12:42:57→12:42:58 | 0 | `VERDICT GREEN`；`over-budget=59 covered by roster=59 not in roster=0`（1,457 B） |

⚠ 门禁跑时 279-r1 的 `run.go`／`config_readers_255.go` 改动在场（12:44 由 `a1cf63b2` 落件）；三门禁全绿，不受其影响。

## 没做（具名）

- ⛔ 未跑整包 `./cmd/wisp/`；⛔ 未 push；⛔ 未翻任何 AC 框；⛔ 零 Go 文件改动。W2 的"等级行"不是从 W2 自己的日志直读（其断言路径不 dump 审计行）——直读证据来自同形 W1 发的 overlay 红（本件 §改后）；W2 侧成立链条＝红句（队列超时）＋封存绿（5.15s）＋同形分类读数，逐层具名在 01 件。