# 票 67 AC#3（判据①）— `wisp providers` verdict 列改 ASCII `PASS`/`FAIL`

- 日期：2026-09-21（UTC）
- 代理：agent-ticket67-ac3
- 裁定依据：编排者 09:01 派单（A23 判据②）。**覆盖面扩展（判据①的 ban #8 扩到 `internal/`+`cmd/`）归票 70 落地之后**，本文件不涉。
- 基线：工作树 HEAD `98fa8ae`（票 70 在派单后又落了一条 commit；`38b3715..HEAD` 之间 **cmd/wisp/ 零改动**，`git log -- cmd/wisp/` 为空）。

## 1. 字形清点（改动前，全 U+2713/U+2717）

| 文件 | 行 | 性质 | 处理 |
|---|---|---|---|
| `cmd/wisp/providers.go` | `:7` 注释「声明 ✓ / 实测 ✗」 | 2 字形 | 改「声明 PASS / 实测 FAIL」 |
| `cmd/wisp/providers.go` | `:41` help 文本（**用户可见**，usage 打到 stderr） | 2 字形 | 同上 |
| `cmd/wisp/providers.go` | `:201/:203` verdict 字面量（**用户可见**，经 `:209` 打 stdout） | 2 字形 | `"FAIL"` / `"PASS"` |
| `cmd/wisp/providers_test.go` | `:4` 注释「声明 ✓ / 实测 ✗」 | 2 字形 | 同步改 ASCII（同目录、将来一并被扫） |

**providers.go 一次清 6 处（4 行），providers_test.go 清 2 处，合计 8 处。** 全文件按 ban #8 的完整区间
（U+2190–U+2BFF、U+1F300–U+1FAFF、U+FE0F）扫描：改后零残留（含变异试验后的复查）。

## 2. 测试改动

- 先 grep：`providers_test.go` **不存在**断言 `✓`/`✗` 字形的既有用例（既有断言是
  `Contains(out,"thinking")` / `Contains(errb,"能力实测不符")` 等）——**一条没删、一条没放宽**。
- 新增 `probeVerdicts` 解析器 + **正向**判据：verdict 列每一行必须**恰好等于** `PASS` 或 `FAIL`
  （不是"不含 ✓"），并钉两个可证点：
  - `TestProvidersProbeRecordsMeasuredThinkingFalse`：`thinking=FAIL`、`fc=PASS`（与持久化 `h.Probe.FC=true` 对口）。
  - `TestProvidersProbeRecordsMeasuredThinkingTrue`：`thinking=PASS` 且**全列 PASS**（fixture 声明四项全真，
    无 mismatch ⇒ 任何 FAIL 都会与 stderr 断言互相矛盾）。

### 变异自证（种子：verdict OK 分支临时改 `"✔"` U+2714，测完即撤）

```
--- FAIL: TestProvidersProbeRecordsMeasuredThinkingTrue (2.02s)
    providers_test.go:191: verdict for fc = "✔", want exactly PASS or FAIL
    providers_test.go:191: verdict for vision = "✔", want exactly PASS or FAIL
    providers_test.go:191: verdict for thinking = "✔", want exactly PASS or FAIL
```

⇒ 正向判据确实咬得住"塞第三个符号"，旧的"不含 ✓"式断言对 U+2714 是瞎的。撤销后复扫零残留。

## 3. 门禁（真实退出码）

| 命令 | 结果 |
|---|---|
| `go test -count=2 ./cmd/wisp/` | `ok github.com/CarlosShao/wisp/cmd/wisp 76.641s`，**rc=0** |
| `go vet ./cmd/wisp/` | **rc=0** |
| `gofmt -l cmd/wisp` | 空输出，**rc=0** |

⚠ 环境坑（既有属性，票 63 log 已记）：cmd/wisp 测试二进制链接 sherpa-onnx，**跑前必须把
`third_party/sherpa-onnx` 放进 PATH**，否则整个测试包以 `exit status 0xc0000135`（DLL 未找到）起不来。
本次实跑即先踩后修，上面 rc=0 是带 PATH 的。

## 4. 用户可见那张表的实跑（live mockllm，非测试内缓冲）

`go build -o $TEMP/.../mockllm.exe .`（`tools/mockllm`）→ `-addr 127.0.0.1:0 -print-addr` 起服务，
临时 data dir 放 fixture 同款 `config.toml`，然后：

```
$ WISP_ENV=test WISP_TEST_DATA_DIR=$TEMP/wisp67ac3/data go run ./cmd/wisp providers probe acme/mock-small
wisp providers: acme/mock-small 实测（p50 5ms，写入 true）
  fc 声明 实测 PASS
  vision 声明 实测 PASS
  thinking 声明 实测 PASS
rc=0
```

注入 `{"thinking":"broken"}`（`POST /__control/capability`，http=200）后同一命令：

```
wisp providers: 能力实测不符：acme/mock-small 声明支持 thinking，实测不可用（thinking probe: no reasoning delta ...）
wisp providers: acme/mock-small 实测（p50 6ms，写入 true）
  fc 声明 实测 PASS
  vision 声明 实测 PASS
  thinking 声明 实测 FAIL
    细节：thinking probe: no reasoning delta ...
wisp providers: 1 项声明与实测不符
rc=0
```

help 面（`go run ./cmd/wisp providers` 无子命令，usage 打 stderr）第 9 行现为
`「声明 PASS / 实测 FAIL」 mismatch`；`go run` 把子进程 exit 2 压成 1 是仓内已记录的既有行为
（票 67 AC#2 log 备注），非本次引入。

## 5. 覆盖面扩展前必须一起清的第二批（**只读未改**）

`internal/llm/probe_health.go` 注释里同形字 **3 行 / 6 字形**：`:17`、`:120`、`:203`（每行各含 ✓ 与 ✗）。
该文件归票 66/69 一族，本次未触碰。A23 里编排者的未决点仍然有效：**`emojiRe` 到底扫不扫注释尚未核实**，
票 70 扩面代理必须先实测这一点再决定这批的处理方式。
