# 286-w1 / 10 改前改后三数并排（票 286 AC#4 ①②）

尺（两跑同一条命令、同一工作树、同一条计数规则）：

```
GOFLAGS= go build ./...
GOFLAGS= go test ./internal/agent/ ./internal/tools/ -count=1 -v
PASS= grep -c '^--- PASS\|^    --- PASS'
FAIL= grep -c '^--- FAIL\|^    --- FAIL'
SKIP= grep -c '^--- SKIP\|^    SKIP'   # 形如 '^--- SKIP\|^    --- SKIP'
```

计数规则说明：整包 `go test -v` 的多包一次跑，逐条 `--- PASS`/`--- FAIL`（含缩进的子测试）统一计数，
⛔ 不按包拆分（拆包需再切一次输出、两跑用同一把尺即可）。

## 三数并排

| 读数 | 改前（基线） | 改后 | 判 |
|---|---|---|---|
| `go build ./...` rc | `0` | `0` | 通过 |
| `go test ./internal/agent/ ./internal/tools/ -count=1` rc | `1` | `0` | 通过 |
| PASS 枚数 | **392** | **393** | **只增不减：成立（+1）** |
| FAIL 枚数 | **1** | **0** | 归零 |
| SKIP 枚数 | **0** | **0** | ⛔ 零新增 `t.Skip`（前后逐字同 0） |
| `internal/agent` 包行 | `FAIL github.com/CarlosShao/wisp/internal/agent 1.905s` | `ok github.com/CarlosShao/wisp/internal/agent 2.125s` | 通过 |
| `internal/tools` 包行 | `ok .../internal/tools 13.898s` | `ok .../internal/tools 14.249s` | 未受影响 |

+1 的来源＝`TestGoldenSingleToolCall` 本身由 FAIL 翻 PASS（FAIL 那一枚不占 PASS 计数）。
⛔ 无其它测试被删、被 Skip、被改名：改前那 1 枚 FAIL 就是本票目标，改后 FAIL=0。

## ① 改前红名册（具名，全文只有这一枚红）

```
--- FAIL: TestGoldenSingleToolCall (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/agent	1.905s
ok  	github.com/CarlosShao/wisp/internal/tools	13.898s
```

改前红句逐字（`logs/before.md:92`；task id 每次运行重新铸造，故与我引用的派单串不同值、同形）：

```
    loop_golden_test.go:70: call identity = {TaskID:6b81876c-cf16-4c6a-b57b-b3e4f9e45ed1 CorrelationID:6b81876c-cf16-4c6a-b57b-b3e4f9e45ed1#call_e1 CallID:call_e1 Name:echo Args:{"text":"22 摄氏度，晴"} Timeout:2s}, want task 6b81876c-cf16-4c6a-b57b-b3e4f9e45ed1
```

⚠ 行号确认：派单正文步骤里写了 `logs/before.txt`，同一派单文末「输出落盘规矩」写「只有 `*.md` 会被门禁放过，
读数一律写进 `.md` 件里」⇒ 按更靠后且更严的一条执行，落 `logs/before.md`、`logs/after.md`。⛔ 无 `.out`。

## 改后同一条测试的绿句（`logs/after.md:92`）

```
--- PASS: TestGoldenSingleToolCall (0.00s)
```

## ⛔ 未在射程内的其它红：无

改前整包（`internal/agent` + `internal/tools`）除目标那一枚外**没有第二枚红**，
所以本腿没有"具名报回但不顺手修"的别家红。`cmd/wisp` 的 boot Ctrl+C flake 与 winlive 红不在本腿跑的这两包里，
本腿未跑 `./cmd/...`，不据以作任何断言。
