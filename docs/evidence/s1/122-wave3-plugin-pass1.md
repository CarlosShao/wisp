# 122-wave3-plugin-pass1 — 清账波 3：pass1 最后 2 条 → 0（`internal/plugin`）

- 日期：2026-11-13（+08，开工约 14:2x）
- 腿：写腿（派单 `.scratch/wisp/probes/orch/20261113/txt/122-wave3-plugin-dispatch.txt`）
- 开工 HEAD：`79890d2b3`（dev，A1834 的 ledger 落点）
- 工具版本（票 122 AC#7(b)/(c)：读数必须同报版本）：
  `staticcheck.exe 2026.2.2 (0.8.2)` · `gofumpt v0.12.0 (go1.27.1)` · `go version go1.27.1 windows/amd64`
- 写面：`internal/plugin/manifest.go`（1 处注释）＋ `internal/plugin/panic_isolation_51_test.go`
  （新 doc 块 ＋ 1 行 `//lint:ignore`）＋ 本证据件 ＋ 票 122 面 Progress log。零越界。

## 1. 起手复测（"前"读数，自己复跑）

`GOOS=linux CGO_ENABLED=0 staticcheck -checks=all,-U1000,-SA9009 ./...`：

```
internal\plugin\manifest.go:11:1: comment on exported const CurrentSchemaVersion should be of the form "CurrentSchemaVersion ..." (ST1022)
internal\plugin\panic_isolation_51_test.go:62:3: assignment to nil map (SA5000)
```

＝**2 条**，与派单逐字一致。附注：本机跑全仓时输出里另有一行
`-: build constraints exclude all Go files in .../sherpa-onnx-go-linux@v1.13.8 (compile)`——
staticcheck 源码（`lintcmd/cmd.go` `printDiagnostics`：`shouldExit["compile"]=true`，`numErrors>0 → exit 1`）
证实 compile 类无条件计入退出码 ⇒ **本机全仓跑 rc 恒为 1，与 finding 数无关**；
该行在开工基线里逐字存在、CI 口径单列为 crash-lines。清后用它判读必须连行数一起看（finding 行 0 条）。

## 2. 清后（"后"读数）

- 全仓同一命令：**finding 行 0 条**（仅剩上述开工即有的 1 行 sherpa compile crash-line；exit code 语义见 §1 附注）。
- scoped 复证：`GOOS=linux ... -checks=all,-U1000,-SA9009 ./internal/plugin/...`＝**零输出、rc=0**
  （该包不牵 sherpa，无 crash-line，rc 语义干净）。
- 判定：**pass1＝2 → 0**。零 `-checks` 放宽、零 CI 改动、零 staticcheck.conf；唯一抑制＝具名 `//lint:ignore SA5000`
  1 行（§4），ST1022 是真改注释。

## 3. ① `manifest.go:11` ST1022 — 真改注释（拆组，信息零丢失）

- 原：`:11-30` 是**一整块注释**挂在常量头上，首行是 "C14 plugin manifest — Tier-1 declarative base..."
  （文件/契约总述），`CurrentSchemaVersion` 那两句被埋在块尾，不是首句 ⇒ ST1022。
- 改（两刀）：
  1. **拆组**：总述块与常量句之间由 `//` 行改为**空行**，让 `// CurrentSchemaVersion ...` 成为常量
     自己的 doc 块（空行切开注释组，首句即以符号名开头）；
  2. **补真话**：常量 doc 追加一句真话——"It is the value validator.go requires, so a future schema
     bump must change this constant and the validator's refusal message together, not one of them
     alone."（依据：`validator.go:56-59` 现跑——`m.SchemaVersion != CurrentSchemaVersion` 拒载并把
     "this loader implements %d" 写进拒绝消息；该比较与消息内插**两处都写死引用这个常量**，句句可核）。
- **为什么没动总述块、没搬任何原有句子**：原有两句本来就是真话（"the only schema_version this loader
  understands" 与 "SPEC-07 §4 example: schema_version = 1"，后者现跑 `docs/specs/SPEC-07-tools-and-plugins.md:80`
  仍在），问题只是它们**长错了地方**（挂错块、非首句）。所以正确的刀是"拆+归位"，不是"改写"——
  总述块保留在文件类型区上方，原文一字不丢。
- 与 wave2 正面样本（`compress.go:42` 拆块）同性质：**拆开长错地方的注释组，不是模板填空**。

## 4. ② `panic_isolation_51_test.go:62` SA5000 — 独立复核采纳上腿判断；抑制落点与上腿给的不同

### 4.1 独立复核（自己重读代码，非照抄）

**结论：采纳**——夹具有意、断言有效、不该修。逐条依据（全部现读本文件与调用面）：

1. `crashHost` doc（现 `:40-43`）明写 "panicVal == nil selects the nil-map-write form: the classic
   nil-dereference host bug, which arrives as a runtime.Error"——nil-map 写是**被测输入**，不是夹具 bug。
2. 第二用例 doc（现 `:155-158`）明写它是"the other panic shape a host bug produces: a runtime.Error
   (nil map write / nil dereference) rather than an explicit panic"，且该用例夹具特意 `panicVal` 缺省
   （`:162` 注释 "nil panicVal => nil map write"）⇒ 两用例刻意覆盖**两种崩溃形状**；
   `panic("...")` 替换会**拆掉 runtime.Error 这半个** AC#4 断言。
3. 该错误行为真实发生且被断言消费：`m["boom"]` 真的 panic（runtime.Error），被 recover 边界捕获后
   第二用例验了类（tool）、归因（`entry=C24FSWrite`）、审计行、"宿主存活"（`reg.Count()==0`、后续
   Execute 成功）、"无重试"（`host.Count()==1`）。**删除或防御（`m = map[string]string{}`）= 拆断言**。
4. 全文件只有这一个 nil-map 写；test 文件无独立构建面，无其他 GOOS 考量。

### 4.2 抑制落点：派单骨架写的 `:61` 行**从机制上不可能生效**（已证，非口味问题）

- 派单骨架与 wave2 证据件 §5②都说"在 `:61`（`var m` 声明行）加一行具名抑制"。**我落第一发就是按这个落的，
  staticcheck 直接证伪**：报 `this linter directive didn't match anything` ＋ SA5000 照报。
- 读 staticcheck 0.8.2 本机源码（module cache）拿到机制：
  - `lintcmd/lint.go` `lineIgnore.match`：`pos.Line != li.Line → false`——**诊断报告行必须等于指令
    `NodePosition` 所在行**；
  - `lintcmd/runner/runner.go` `serializeDirective`：`NodePosition`＝指令注释经 `ast.NewCommentMap`
    **所关联 AST 节点的位置**；
  - `analysis/lint/lint.go` `ParseDirectives`：关联由 `ast.NewCommentMap(fset, f, f.Comments)` 决定
    （doc 注释→声明；独立行注释→紧贴其后的节点）。
- SA5000 报在 `:62` 的**赋值语句**（`m["boom"] = ...`），而 `:61` 的指令注释关联到 `var m` 声明 ⇒ 差一行，
  永不匹配；挂函数 doc 又差得更远。staticcheck 自家 `cache.go:202` 的 `//lint:ignore SA4006` 正是
  **独立行、贴在被报语句正上方**——照这个形状落。
- **终态落法**（与 staticcheck 官方自家用法同形）：指令独立成行放在**赋值语句正上方**（现 `:66`），
  措辞按实情写：写明"这个 nil-map 写**就是被测对象**（票 51 AC#4 要求宿主扛住 runtime.Error 形状，
  不仅是显式 panic 形状；nil-map 写是最便宜且确定性的激发方式；删掉它＝拆掉半个断言）"，
  并登记出处（A1773 原始登记、A1833 波 2 独立复核）。**派单骨架里的 "A1834" 是票 70 的账，SA5000
  的原始登记是 A1773**，引用按真号写。
- 撤销口令：**删掉那一行 `//lint:ignore`（＋可选地删 CallC24 的 doc 块）即复原**，SA5000 立刻重现。
- 连带修正：wave2 证据件 §5②的落法描述（"在 ：61 行加"）按本节证据属**机制性不可用**，
  append-only 不回改，由本件与本条 Progress log 记录替代。

### 4.3 顺带的小诚实

`CallC24` 此前没有 doc。这次加的 doc 块说三句真话：本方法做什么、两种崩溃形状怎么选、SA5000 在哪
被抑制（"suppressed where it fires; the write is the subject, see there"）。它同时防一件事：
未来有人把指令挪位后，doc 里"suppressed where it fires"这句话就变成假话，会逼着挪的人回来对齐。

## 5. 门禁读数（五项，逐条如实）

| 门禁 | 读数 |
|---|---|
| `go build ./internal/plugin/` | rc=0 |
| `go vet ./internal/plugin/` | rc=0 |
| `go test -count=2 ./internal/plugin/ ./internal/plugin/manifestslot/` | `ok ... 13.296s` ＋ `ok ... 31.290s`（`-count=2` 不缓存；终态代码上跑） |
| `gofumpt -l internal/plugin` | **0 行**（v0.12.0） |
| `sh scripts/d22scan.sh` | **clean，rc=0**（481 个生产 Go 文件；全账本行逐条 PASS；日志副本 `/d/tmp/wisp/20261113/122-wave3/wisp_122w3_d22.log`） |
| pass1 复跑（全仓） | finding **0 条**（仅剩开工即有的 1 行 sherpa compile crash-line；§1 附注） |
| pass1 复跑（scoped `./internal/plugin/...`） | 零输出 **rc=0** |

- 测试过程读数（对每个改法版本都重跑过 `-count=2` 两包，全 ok；三个中间态见 §4.2）。
- 未跑 `go test ./cmd/wisp/`（票 98 豁免，派单明令别跑别报）。

## 6. 本波未动的共享树既有改动（留痕，防误领）

`git status` tracked 改动仍有 3 个他腿/外部件（开工即有，零触碰、零提交）：
`.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go`、
`.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt`、`frontend/fixtures/composer-states.html`。
（`build/mut_rvw.go` 本机残件在但被 .gitignore 忽略，未动未列；scoped build 不受影响。）

## 7. next=

- 编排者推送后取 CI `lint` job staticcheck 步两段自报读数：期望 **pass1＝0**（若 runner 的 go patch
  版不同，crash-lines 行数可能与本机不同，判读以 finding 行为准）。
- 本波后 pass1 分母归零；票 122 剩余工单面＝pass2 的 U1000 清账与 windows-only 14 条（归后续波次）＋ AC#6。
