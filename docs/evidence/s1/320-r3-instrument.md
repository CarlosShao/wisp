# 320 AC#4 仪器收口 — 非实现者裁决（编排者，2026-11-04 14:0x +08）

> 来路：落地腿 `320-r3`（workflow 覆盖 magpie/cline-pass/deepseek-v4.1-flash-1017）把两文件改到位、测试全绿，但**死前未 commit、未出证据件、未写 Progress log、report=null**。编排者按「收件只认盘上三把尺＋自己复跑」把它从工作树捞起、补自证、出本裁语。裁决者（编排者）＝该腿的非实现者。

## 一、改了什么（两文件）
- `cmd/wisp/panel_models_320.go`
  - `modelsFetchBudget320`：`const` → `var`（加注释：让「慢供应商」这一形可测；生产唯一写者是本票测试，跑完复原）。
  - 删 `modelsListErrorCap320` 常量与 `boundedFirstLine320` 函数（死码）。
  - **错误句（AC#4 第④口径，裁定 A1483 ruling 4）**：由 `boundedFirstLine320(err.Error(), 200)`（会把上游首行带进页面句）改为**只带 D37 分类**：`fmt.Errorf("拉取失败（供应商 %s；D37 分类 %s）：没能取到它的模型列表，检查基址与密钥后再试", provider, observe.ClassOfOrInternal(err))` —— **不含 `err.Error()` 的上游正文**，与文件头 `:19-22` 自述「does not relay the upstream error body verbatim」一致。
- `cmd/wisp/panel_models_320_test.go`：新增 AC#4 四条有牙测试
  - `Test320FetchLegWritesTheKeyNowhereOnDisk`（不落盘明文 a4m5；含**正控**：植密钥断言探测器能看见）
  - `Test320FetchLegTimesOutWithinItsBudget`（超时 a4m3；慢夹具＋缩预算验机制）
  - `Test320FetchLegClassifiesNon200WithoutRelayingTheBody`（错误分类＋不转述）
  - `Test320FetchLegCarriesIDsAndNoUpstreamTextOrKey`（id 列表＋成功/失败回执零密钥零上游串）

## 二、门禁读数（编排者现跑，诚实树 @ 工作树）
- `go build ./cmd/wisp ./internal/panel` rc=0；`go vet ./cmd/wisp ./internal/panel` rc=0；`gofumpt.exe -l` 两文件**空**；`sh scripts/d22scan.sh` rc=0 clean（internal/=848、cmd/=196 为被扫文件数）。
- `go test -count=2 -v ./cmd/wisp -run Test320` ⇒ **rc=0 RUN=24 PASS=24 FAIL=0 SKIP=0**（12 函数×2 趟，含 AC#1-3/5-7 既有腿 + AC#4 四条新腿）。

## 三、变异自证（不转述 tooth，编排者亲跑 -overlay 仓外 mutant，⛔ 未改工作树产码）
- 变异：把错误句改回 `fmt.Errorf("拉取失败（%s）：%s", provider, err.Error())`（relay 上游）＋保留 `_ = observe.ClassOfOrInternal(err)` 保活 import；`go build -overlay` rc=0。
- 结果：`Test320FetchLegClassifiesNon200WithoutRelayingTheBody` **RED**，红线逐字：
  `the page-visible refusal relays the upstream body's own words: "…拉取失败（deepseek）：provider (<nil>): HTTP 500 … {\"error\":{\"message\":\"UPSTREAM_LEAK_320XYZ: the provider's own words\"}}"`
- 对照诚实树同一测试 **PASS**。⇒ 该尺对「把上游正文塞进页面句」这一形**有真牙**，我对第④口径的代码修正确有效。
- ⚠ 首发 0xc0000135＝sherpa DLL 未 staging、测试二进制没起（平台陷阱#15），**不作为红读数**，已带 PATH staging 重跑，上为有效读数。

## 四、其余两条缺口的牙（说明性，未另做变异）
- **不落盘明文**：`Test320FetchLegWritesTheKeyNowhereOnDisk` 自带正控——把随机密钥写盘后断言探测器「恰好找到 1 处且就是那枚」⇒ 探测器有眼；同尺遍历配置根与 CWD 断言「零命中」⇒ 处理腿若把密钥写文件即被逮（正是 a4m5 那一形）。
- **超时**：慢夹具（server 挂起至 cancel）＋把预算临时缩到 120ms，断言被分类为 `observe.ClassNetwork`、elapsed<2s、endpoint 恰命中 1 次。⚠ 具名：该测试**临时覆写** `modelsFetchBudget320`，故只验「超时机制被强制执行」这一层，**生产默认 8s 的常数本身不被此尺钉住**（≤ A972③ 10s 属读码层）。

## 五、判
AC#4 四条（超时／错误分类／不回显密钥／不落盘明文）＋返回 id 列表，de **有牙仪器证明**齐；判 **AC#4 成立**，翻勾。320 八格全 [x]，转 -done（规格/退可「撤 320」）。
