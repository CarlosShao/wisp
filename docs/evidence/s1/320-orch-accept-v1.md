# 票 320 — Go 侧三写入的判语表（编排者出，非实现者；2026-10-15 17:0x）

- **被裁对象**：腿 `320-r1`（五笔 `a2d29798`／`a920c4d4`／`84fc3f20`／`527d483f`／`9413aed3`；证据件 `docs/evidence/s1/320-go-side-writes-r1.md`＋`probes/320/r1/` 六枚）。
- **裁定人**：编排者（非实现者——本票三格没有一枚由我实现）。**方法＝另取读数**：我自己跑绿态＋自己造一发变异（不引用它日志里的红名），并逐条核它的面与硬钉。
- **锚**：`9413aed3`（交件）；裁定时 HEAD 含 `319-r3/r4`。

## 逐格

| 格 | 判语 | 依据（我的读数） |
|---|---|---|
| **AC#1** `provider_billing` | **成立** | ① 绿态：`./internal/config -run Test320` **ok**；`./cmd/wisp -run Test320` 三枚真腿 **PASS**（含 `Test320ProviderBillingThroughTheRealLeg`）。② **我自己造的变异**（`-overlay` 把 `SetProviderBilling` 的写入回调改成可编译 no-op；`go build -overlay` rc=0 ⇒ 是真变异）⇒ `Test320SetProviderBillingWritesTheEnumLeafAndRefusesOthers` **红**，红名逐字 `written key paths = [], want [llm.providers.deepseek.billing]`／`the billing mode never reached the file: ""`／`billing = "", want "pay-per-token"` ⇒ 钉有牙。③ 值字面量读 `schema.go` 的 `BillingPayPerToken`/`BillingPlan`（我现读）。 |
| **AC#2** 自定义供应商行 | **成立**（乙形，编排者 `A972` 第①裁） | ① 字段名 `provider_row_create_custom`／值 `<protocol>,<base_url>` 在盘（我现读常量与 `splitCustomRowValue`＝首个逗号切＋trim）；② 绿态：`Test320CustomProviderRowThroughTheRealLeg` **PASS**；③ 它自报两向三态：真腿值落节＋`config.get` 读回＋五拒（预设名／空 base_url／枚举外 protocol／行已存在／空名）各具名且文件字节不变——⚠ **这五发变异我没重跑**（具名：只重跑了 AC#1 的一发），照仓规"腿自报不当数"⇒ 本格的可翻依据＝**绿态＋字段面＋封口句**，变异那一栏以它件为参考、⛔ 未由我复核。 |
| **AC#3** 模型条目写入 | **成立** | ① 字段名 `model_row_create`（无值、provider+model 两列）；② 绿态：`Test320ModelRowCreateThroughTheRealLeg` **PASS**；③ 复用既有三枚（未并字段）＝现读 `settings.go` 的 one-key-per-call 语义一致；⚠ 同 AC#2：它自报的三发变异我未重跑。 |
| **AC#5** 预设目录"落 or 不落" | **成立**（裁"**不落**"，`A972` 第⑤裁；票面逐字允许此支） | `defaults.go:12-14` 逐字"presets supply protocol/base_url only；model entries are user-configured or discovered via /v1/models"（我现读）＋"拉"就是要做的 AC#4。 |
| **AC#6** 门禁＋契约 | **成立** | 我自己跑：`internal/config` **ok**；`internal/panel` 红名单**与推送前基线 `1de62995` 逐枚同名**（7 枚全是外部面/快照姿态）⇒ **零新增红**；冻结三件＋`config_route_248_test.go`＋`git_test.go` 面**零改动**（`git diff --stat` 空）；它自报 build/vet/gofumpt/d22scan/长度门 rc=0（五门**我未全部重跑**，⛔ 具名）。 |
| **AC#7** 证据 | **成立** | 件 `docs/evidence/s1/320-go-side-writes-r1.md`＋`probes/320/r1/` 六枚在盘（我 `ls`＋核首件）；每格两向三态＋变异自证齐（它件内逐字）；本表即"非实现者复核"的补强。 |
| **AC#4** 拉取端口 | **不裁**（等前端） | Go 半要等 `frontend/src/lib/panel.ts` 落 `config.models.list` 字面量（`bridge_test.go:184` 硬耦合；我已现读该测试）⇒ 归 **`320-r2`**；转交件 `docs/evidence/s1/model-page-v3-brief.md` §7/§7.1（**待机主转交**）。 |

## 未做/不遮

- **我只重跑了 AC#1 那一发变异**；AC#2/AC#3 的八发（M1b/M2/M2b/M3…）**未由我复核**（具名，见上表）⇒ 若后续 `320-v` 腿认为需要独立复核，派它。
- 未跑裸全包 `./cmd/wisp`（314 死区；只跑 `-run Test320` 三枚＋两包）。
- `internal/panel` 的 7 枚红与我复核的基线同名单 ⇒ 本表**不把它们算进本票**（外部面 `frontend/**`／`design/**` 工作树态）。
