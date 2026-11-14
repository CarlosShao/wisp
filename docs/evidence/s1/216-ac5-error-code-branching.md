# 216 AC#5 证据件 — 错误提示按错误码分支（落地腿 216-r3）

- 腿：`216-r3`（写；写面＝`internal/panel/**` 两枚新件＋一枚既有件一行，⛔ `cmd/wisp/**`、`frontend/**`、`scripts/**`、`models/**`、阈值/golden/契约零字节）
- 起手锚：`git rev-parse --short HEAD` ＝ **`7bdf09042`**
- 日期：2026-11-14 +08

## 0. 起手三检（派单要求）

1. **AC 勾选现量**：票面判据区读数＝AC#1/#2/#3/#4/#6 已勾、**AC#5 未勾**（唯一待做格）。
2. **Progress log 全文**：80 行逐条读过。关键前账——`216-v2` 判 AC#5ⓐ「成立但没牙」、ⓑ「不成立」（钩子 `MessageKey()` 在 `internal/observe/errors.go:163`、生产消费者 0）；编排者裁定 ⓑ 归口 S6 DEFERRED、**本腿不碰那层**；ⓐ 的「没牙」立账 `A1010` 排队列尾。
3. **`git log --oneline --grep=216`**：已跑（读数存档在本节末尾的提交清单里），`216-r1`（`b7676daf`＋`92e8cc57`）、`216-a1/a2`、`216-v2/v3`、`216-r2` 全在盘——已落地部分（显示洗、两句话、算子钉、两臂差试）⛔ 本腿零重造。

## 1. 票面反面教材逐字复核（派单点名「落地前自己看原文」）

`ChatInput.tsx` **本仓不存在**（与 `216-a1`/`216-v2` 读数同向）。逐字复核走的是本仓的转写：
`docs/reports/survey-2026-09-28-composer-plus-menu.md:378-381` 逐字＝
`normalized.includes('payload too large') || normalized.includes('413') || normalized.includes('entity too large')`
⇒ 票面「靠 `413` 子串」是**三枚子串 OR**（`216-v2` 已具名该简化），本腿照转写原文复核、未改一字。

## 2. 错误码从哪来（既有源头，零新造）

- **分类法**：`observe.ErrorClass`——D37 的 17 枚枚举常量（`internal/observe/errors.go:19-36`），`Valid()` 钉死集合。
- **提取器**：`observe.ClassOf(err)`（`internal/observe/errors.go:234`）＝`errors.As` 沿错误链提取 `*observe.Error` 的 `Class`；未分类错误返回 `ok=false`。另有 `ClassOfOrInternal`（日志侧）本腿未用——显示层不发明归类。
- **分类错误的生产密度**：`observe.New/Wrap` 在 `internal/config`（写路由）、`internal/llm`、`internal/speech`、`internal/agent`、`internal/audio`、`internal/proc`、`internal/plugin` 等满仓使用（grep 现量 250+ 行）⇒ 分支键不愁「恒绿装饰」。
- **钩子**：`ErrorClass.MessageKey()`（`errors.go:167`）＝`216-v2` 判「未闭合」时点名的那枚零消费者钩子；本腿的分支键**逐类对齐**它的类常量（map 键就是 `observe.Class*` 常量本体，⛔ 错误串里的任何子串）。

## 3. 落了什么（三笔写面）

| 件 | 内容 |
|---|---|
| `internal/panel/error_copy_216.go`（新） | `errorCopy216`：17 枚 D37 类 → 一句人话头句的映射表；`errorCopy216Complete()`：表对 `AllClasses()` 的穷尽尺；`errorSentence216(err)`：**分支点**——`observe.ClassOf(err)` 提类、按类选头句、`err.Error()` 只作原因缀在分支**之后**（票 147「读不到不带名字」那句的纪律）；未分类错误逐字节保持原句。`userFacingErrorSentence216`：派发层拒绝句的成形入口。 |
| `internal/panel/bridge.go`（1 行＋注释） | `RefusedEnvelopeForUser`（panel 里**唯一**页面可见的拒绝句成形点，`composer_dispatch.go:202/:206` 是它的仅有的两个生产调用者）的 `%v` 改 `%s` 过分支：分类错误戴上自己类的头句、未分类错误**逐字节不变**。 |
| `internal/panel/error_copy_216_test.go`（新） | 四枚仪器，见 §4。 |

写面边界：`internal/panel/**` 三枚件，零产码删改（bridge.go 只换那一行格式动词＋补注释）；⛔ 未动 `MessageKey()` 一字（票 331 射程）；⛔ 未动 run_feed 的「human error sentence」（票 145 具名 NOT-landable）；⛔ 未新增 C17 方法、未新增名册键、未动任何既有测试断言。

## 4. 分支怎么测（四枚仪器，全绿）

1. **`TestErrorCodeSelectsTheSentence216`**（对照测试）：同一句错误文本 `upstream 413 payload too large` 喂 17 枚类，逐类断言 `errorSentence216` 恰返回 `errorCopy216[类] + err.Error()`；另钉 `errorCopy216Complete()`——将来 D37 加第 18 枚类而忘配文案 ⇒ 红。
2. **`TestUnclassifiedKeepsThePlainSentence216`**：未分类错误（裸＋`%w` 包裹两层）逐字节保持 `err.Error()` 原文——既有夹具（全部喂未分类错误）不移动的凭据。
3. **`TestErrorCopySentencesAreAllDistinct216`**：同 text 下 17 类产 17 句**互异**——这是结构判别尺。
4. **`TestRefusedEnvelopeCarriesTheClassSentence216`**（wire 正控）：真派发链 `ComposerDispatch.Handle` → 真门 `ConfigWriteHandler.write` → spy 腿返回 `%w` 包裹的 `observe.New(ClassConfig,…)`（与 `internal/config/sectionwrite.go` 生产形状同形）⇒ 页面可见拒绝句必须含 config 类头句且**底层原因仍在**；`applyCalls==1` 钉「真到腿」非装饰。

## 5. 突变正控（换回子串匹配必红）——含一次自抓

第一发突变（sed 式替换）**失败**：突变体缺 `strings` import ⇒ build 失败、rc=1 是假红，⛔ 不算数；还原步又因「文件未提交 `git show HEAD` 取不到」把文件截成 1 字节——**本腿自抓、当场纠正**：改为整文件仓外副本（`D:\wisp216mut-clean.go.txt`）存底＋整体替换法。

第二发（合规）读数：
- 突变＝在 `errorSentence216` 里**按反面教材原形插入** `if strings.Contains(err.Error(), "413") { return "上游限流了，稍等再试：" + err.Error() }`（三枚子串 OR 简化成单枚同形），突变体 `Contains` 落在 `:71`。
- 跑 `TestErrorCodeSelectsTheSentence216` 等三枚 ⇒ **FAIL**（红句逐字）：
  - `class injection: sentence = "上游限流了，稍等再试：injection: upstream 413 payload too large", want "这一步被判定为疑似注入，已按安全规则拦下：…"`
  - `class internal: sentence = "上游限流了，稍等再试：…", want "程序内部这一步出了错，…"`
  - `class resource: sentence = "上游限流了，稍等再试：…", want "资源这一步超限了，…"`
- 同发里 `TestErrorCopySentencesAreAllDistinct216` **仍绿**——如实记：distinct 尺对「单枚子串抢先 return」这形**不敏感**（其余 16 类走 fallback 各自头句仍互异）；**对子串形真正断牙的是 table test**（对照测试）＋wire test。这个不敏感已在仪器注释里如实写明（仪器文件头 §2 说的是「substring 形塌成同一分支」那一族；对**单枚抢先 return** 这一形，牙在 table test）。
- 还原＝仓外副本整体拷回，SHA256 前后逐字一致：`1DC7BD0FBA9E65DD7D1271C327BADA9E4DDBB2742DCEB0A56BC88FA393936909` ＝ 零残留。

## 6. 门禁四数（逐条贴读数）

| 门 | 读数 |
|---|---|
| `CGO=1 go build ./...` | ⛔ 整树 rc=1，但红**全部**在 `build/`——git-ignore 目录、未跟踪的别腿突变残件（`mut_rvw.go` 等引用未定义符号，与本腿零关联）。改跑 `go build ./internal/... ./cmd/... ./tools/...`（全部跟踪包）⇒ **rc=0**。build/ 那枚红在起手锚同色，⛔ 本腿引入。 |
| `go vet`（改的包） | `go vet ./internal/panel/ ./internal/observe/` ⇒ **rc=0**，零输出。 |
| `sh scripts/d22scan.sh` | **clean**（rc=0）：`no D22 ban violations`，射程读数 internal=1007/cmd=278/design=39/frontend=123 逐枚在场。 |
| `gofumpt -l`（改的目录） | `go run mvdan.cc/gofumpt@latest -l internal/panel` ⇒ **0 枚文件列出**、rc=0。 |
| 相关包 `go test -count=1` 四数 | 四发各自 ok：`./internal/panel/` **ok 49.162s**；`./internal/observe/` **ok 7.845s**；`./internal/config/` **ok 6.708s**；`./internal/risk/` **ok 17.481s**。四发全 rc=0（四数＝四发 RUN/PASS/FAIL/SKIP 口径下 FAIL=0、SKIP=0；panel 整包无 -v 故四数以 ok/FAIL 汇总行记，红 0 枚）。 |
| 新仪器单飞 | 4 枚 `-run` 名册单发 ⇒ **4 PASS / 0 FAIL / 0 SKIP**、ok 0.272s。 |

## 7. 残余（如实具名，⛔ 抹平）

1. **distinct 尺对单枚抢先 return 子串形不敏感**（§5 如实记）——要它对那一形也断牙，可以把尺收紧成「断言 17 类头句前缀逐类等于表值」（table test 已是这形）＋「未分类路径逐字节」（已有）。本腿判 table test 已覆盖该形，distinct 尺保留作族尺；不再加第四枚仪器（收益零）。
2. **AC#5ⓑ 那半（「按码选句」铺满全仓用户可见文案）**＝编排者已裁归口 S6 DEFERRED（账 `A1010` 同族），本腿只落 panel 显示层这一寸的**分支点＋消费点＋仪器**，⛔ 自称铺满。
3. **`d22scan` 无「错误串分支」禁令**（`A1010` 的仪器工程排队列尾）——本腿的防线是测试面那四枚，不是静态禁令；后来者要绕开测试写子串形，仪器面才拦得住。
4. `build/` 目录的未跟踪突变残件（`go build ./...` 整树红的真身）——非本腿字节，归编排者处置；本腿按「仓内不删东西」零触碰。

## 8. 交付清单

- `internal/panel/error_copy_216.go`（新，85 行）
- `internal/panel/error_copy_216_test.go`（新，125 行）
- `internal/panel/bridge.go`（+6/−2）
- 本件＋票尾一条 bullet；commit 显式 pathspec、⛔ push。
