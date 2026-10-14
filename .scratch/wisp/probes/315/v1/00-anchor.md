# 票 315 验收腿 315-v1 —— 锚与读数名册（判语见 `docs/evidence/s1/ticket315-315-v1-accept-v1.md`）

- 腿身份＝**验收（非实现者）**；被裁＝产码 `bfc74ae9`＋交付 `b7a15297`（实现者 `315-r1`）。⛔ 未翻框、⛔ 未动 `Status:`/票名；⛔ 未 push。
- 锚（起手现读）：HEAD＝`b7a15297`；读数窗口 **2026-10-14 11:32–11:45 +08**。窗口内 HEAD 漂到 `e3788f08`
  （288-r1 停手上报）＋编排者两笔文档（`d286bd45`/`fcb0a2b2`/`68e235ff`）；逐枚 `--name-only` 核过，
  **与本腿判读面零相交**：`cmd/wisp/run.go` sha256 `71093661…` 全窗前后一致，`cmd/wisp/zz315a1probe_test.go` porcelain 空。
- 环境：`PATH=…/third_party/sherpa-onnx:$PATH`（MSYS 形）；`GOTMPDIR=C:\tmp\wisp315v1gotmp`；`GOFLAGS=`；
  全用 `-run` 小射程 `-count=1`；⛔ 未跑整包 `./cmd/wisp/`（314 死区）。窗口内他腿（236-r4/236-r5/288-r1）无编译相撞，无一发需重试。

## 判语速览（六格全**成立**；详表在证据件）

| # | 格 | 判语 |
|---|---|---|
| 1 | 选案正当性（ⓐ vs ⓑ/ⓒ） | 成立 |
| 2 | 读数独立复现（默认绿／自造反证红／留档核） | 成立 |
| 3 | Case C（project-instructions 种子） | 成立 |
| 4 | 去门控决定 | 成立 |
| 5 | 既有 `-race` 夹具红归因 | 成立（非本修引入） |
| 6 | 门禁抽核 | 成立 |

## 突变/反证件（全部在仓外 `C:\tmp\wisp315v1overlay\`，工作树零足迹）

| 件 | sha256 | 重建法 |
|---|---|---|
| 我自己的突变 `run.v1muted.go` | `d6a9add6…`（全文见证据件；对现读 `run.go` `71093661…`） | `sed 's/^\trt\.bridge\.OpenTask(taskID)$/\t\/\/ v1 counter-evidence: this open is deleted on purpose (acceptance leg 315-v1)/' cmd/wisp/run.go > C:/tmp/wisp315v1overlay/run.v1muted.go` |
| `v1-overlay.json`（我的） | `d3c1010e…` | `{"Replace":{"D:/work/…/cmd/wisp/run.go":"C:/tmp/wisp315v1overlay/run.v1muted.go"}}` |
| 门控复造 `zz315a1probe.gated.go` | `399e8d23…` | 对现读探针件在 `func TestZZ315A1SiblingChildFirstBridgeCall` 首行后插入 `if os.Getenv("WISP_315A1") != "1" { t.Skip(…) }`（sed 一行；原文见证据件 §4） |
| `v1-overlay-gated.json` | `fecc61a1…` | 同上，映射到 `zz315a1probe.gated.go` |
| 他们入库的 `mutation/run.mut-noopen.go`（交叉核） | `3b632e95…` | `diff` 对现读 `run.go`：**恰一行替换**（`rt.bridge.OpenTask(taskID)` → `// MUTATION (ticket 315 r1): …removed`）；其 `overlay-mut-noopen.json` 我复跑过（log `22`，B rc=1 31.21s 同红句） |

还原证：全窗结束后 `sha256sum cmd/wisp/run.go` 仍 `71093661…`；`git status --porcelain -- cmd/wisp/run.go cmd/wisp/zz315a1probe_test.go` 空。

## 读数名册（`.scratch/wisp/probes/315/v1/logs/`；★＝随本腿入库）

| 件 | 字节 | sha256 前缀 | 窗口（+08） | rc | 是什么 |
|---|---|---|---|---|---|
| ★ `10-my-default-three.txt` | 22,536 | `81a223f8` | 11:32:58–11:33:06 | 0 | 默认三形（无门控）A 0.15s／B 0.20s／C 0.14s；`[确认 L2` 命中 0 |
| `20-my-overlay-caseB-red.txt` | 17,931 | `063fa909` | 11:33:47–11:34:28 | 1 | 我的 overlay 反证 B：31.22s 红（原件留盘） |
| `21-my-overlay-caseC-red.txt` | 12,494 | `2e378e05` | 11:34:48–11:35:30 | 1 | 我的 overlay 反证 C：31.17s 红（原件留盘） |
| ★ `20-my-overlay-B-and-C-red.txt` | 30,915 | `4f061f54` | （合订 11:47） | — | 上面两枚的**合订入库件**（原始两枚留盘） |
| `22-cross-run-their-overlay-caseB-red.txt` | 17,861 | `dbe529e5` | 11:35:39–11:36:18 | 1 | 复跑他们 overlay 件：B 31.21s 红 |
| ★ `30-my-strict-runner-degated.txt` | 22,685 | `c8d46ac8` | 11:39:26–11:39:35 | 0 | `runtests.sh ./cmd/wisp/ -run '^TestZZ315A1'`：PASS=3 FAIL=0 SKIP=0 |
| ★ `31-my-degate-premise-gated-skip-fatal.txt` | 925 | `096274af` | 11:40:02–11:40:13 | 1 | 门控复造件在 strict runner：`SKIP is not a pass` rc=1 |
| `40-my-race-197-current-tree.txt` | 37,732 | `adee6bd7` | 11:40:43–11:41:06 | 1 | 197 自批 `-race`（当前树）：7 warning、FAIL 13.99s（原件留盘） |
| `41-my-race-197-overlay-noopen.txt` | 38,020 | `6d3be686` | 11:41:18–11:41:50 | 1 | 同命令 overlay 无本修：7 warning、FAIL 13.95s（原件留盘） |
| ★ `40-my-race-197-both-trees.txt` | 76,132 | `af3cd49b` | （合订 11:47） | — | 上面两枚的**合订入库件** |
| `42-my-197-family-norace.txt` | 4,361 | `398216a2` | 11:42:41–11:43:03 | 0 | 无 `-race` 两枚 PASS（13.66/0.00s） |
| `50-my-race-zz315.txt` | 22,556 | `ddcfaf91` | 11:43:13–11:43:26 | 0 | `-race` 三形 rc=0（ok 5.395s） |
| `60-my-build.txt` | 103 | `5e0f07f6` | 11:43:35–11:43:42 | 0 | `GOFLAGS= go build ./...` 空输出 |
| `61-my-d22scan.txt` | 23,069 | `b8483e91` | 11:43:42–11:44:17 | 0 | `clean`；cmd/ examined 125 |
| `62-my-pathlen.txt` | 1,575 | `bfcaeab7` | 11:44:17–11:44:18 | 0 | `VERDICT GREEN`；over=59／roster=59（59=59） |
| `63-my-test222-race.txt` | 707 | `d1faf8e7` | 11:44:56–11:45:13 | 0 | `Test222 -race` 3 PASS（ok 4.464s） |

## 具名没做

- ⛔ 未跑整包 `./cmd/wisp/`；⛔ 未单独量"并发兄弟"形（写腿已登记候选；机制上 admission 按任务、修后 run 里两枚孩子 close 全 `was_open=true`，本腿不扩修）。
- ⛔ 未修既有 `-race` 夹具红（射程外、测试夹具面）；⛔ 未动 C19/C25 语义/阈值/golden/allowlist/PLAN/specs；⛔ 未动票面 `Status:`/框；仅票尾追加 1 条 bullet。
