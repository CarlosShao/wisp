# pool-validity-2 / 批2 —— 零勾开放票有效性普查（号段 50–99，16 枚）

> 只读普查腿 `pool-validity-2`，接 `pool-validity-1`（死腿）的活。
> 起手 HEAD `949a5b92`（2026-10-06 09:03 +0800）。⚠ 换 HEAD 要重量。
> 本腿**零 go 命令**（并行腿 `231-r1` 此刻真在 `cmd/wisp` 跑定向用例与整包读数）。
> 工具全集＝`git log`／`ls`／`grep`／`wc`／`find`／Read。四档尺与本批名册的表头见 `batch1.md` §0。
> 票面零改动；唯一写入路径＝`.scratch/wisp/probes/pool-validity/2/**`；`pool-validity/1/**` 只读只引。
> ⛔ 零读零引 `frontend/**`、`design/**`、`cmd/wisp/**`、`scripts/**`、`.github/**`、
> `docs/reports/HANDOVER.md`、`docs/evidence/s1/**`、`.scratch/wisp/probes/{231,268,111,evidence-close}/**`。

## §1 名册与判定（16 枚）

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
★ 尺面统一为**提交面**（`git --no-pager grep -n <pat> HEAD -- <path>`）；落 `cmd/wisp/**` 的尺一律
`git ls-tree --name-only`（只取文件名与计数，⛔ 未读内容）。

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 50 | `50-tier1-manifest-plugins.md` | 144151ae 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD -- internal/plugin`；`git grep -n 'not built: plugin engine' HEAD -- internal/config/unwired.go`；`git grep -n 'ErrSlotNotLanded' HEAD -- internal/tools/registry.go` | `internal/plugin/` 提交面**只有 `disposal.go`＋`doc.go`**（DisposalScope 是 C11/票 09 的活，不是 manifest 引擎）；`unwired.go:141-146` 六枚 `plugins.*` 键**逐字挂着 "not built: plugin engine is ticket 50/51"**；`registry.go:27-31` 的 `ErrSlotNotLanded` 仍写 `manifest=ticket 50` ⇒ schema 校验器／安装流／篡改拒绝／net_allowlist 判定四格零落点 |
| 51 | `51-tier2-goja.md` | 6160bc8f 09-19 | **仍成立** | `git show HEAD:go.mod \| sed -n '/require/,/^)/p'`；`git grep -iln 'goja' HEAD -- internal/`；`git grep -in 'goja memory limit' HEAD -- internal/ docs/` | `go.mod` 两截 require 共 30 行、**`goja` 零命中**（依赖都没引）；`goja` 非测试命中全是注释面（`internal/plugin/doc.go`／`internal/tools/registry.go`／`internal/config/schema.go`）；AC#6 那把负向尺 `"goja memory limit"` 读数确为 **0 命中**——但负向尺为真≠这票交付了，中断／结果上限／C24 表面／panic 隔离／默认关闭五格（AC#1–#5）无处存在。备注列具名："AC#6 恰好为真，是这条禁令从没被触发过，不是有人做对了" |
| 52 | `52-d46-command-plugins.md` | 144151ae 09-19 | **仍成立** | `git grep -in 'shell\.exec' HEAD -- internal/ cmd/`（排 `_test`）；`git ls-tree -r --name-only HEAD \| grep -icE 'command-plugin\|d46'` | `shell.exec` 非测试命中**全在"它不存在"那几句里**：`unwired.go:63` 逐字 `missing: "no shell.exec tool is registered in internal/tools, so nothing ..."`，`:70`／`:76` 两处 "it lands with the shell.exec tool"；D46 命令插件载体文件 **0 枚**；批一 22/23 那把 `Name()` 全集尺里没有它 ⇒ 哈希钉／注入拒绝／抽取器／env 白名单／超时杀树五格无处可跑 |
| 53 | `53-larkcli-plugin-e2e.md` | 144151ae 09-19 | **量不到** | `git ls-tree -r --name-only HEAD \| grep -i lark` | 全仓 `lark` 命中**只有本票票文件自己那一枚**（其余零）。判据＝真插件打真账号 E2E＋四场失败预演＋`docs/evidence/s7/` 取证包（禁读）⇒ 归口：能跑账号级 E2E 的腿。本腿不替它判"能不能做"，只记**上游票 50 读数＝引擎未建**这一件事 |
| 54 | `54-s7-acceptance.md` | 144151ae 09-19 | **量不到** | 判据尺＝S7 切片验收（并发矩阵／失败预演／更新+回滚演练／开放决策定案） | 归口：运行期＋演练＋人工定案。上游 50/51/52/53 本腿读数＝三枚〔仍成立〕＋一枚量不到 ⇒ 这张卡今天不可能绿，但"绿不绿"不是静态尺读得到的 |
| 55 | `55-s8-macos-port.md` | e208efb1 09-21 | **仍成立** | `git ls-tree -r --name-only HEAD -- internal/ \| grep -icE '_darwin\|keychain'`；再数 `grep -ic 'darwin'` 全仓 | `_darwin` 文件与 Keychain 代码 **0 枚**；仓内 11 枚 `_other.go/_unix/_posix` 是"非 Windows 编译得动"的桩，不是 mac 落地 ⇒ 球体两版实现／SecretStore 移植矩阵／PathResolver mac 分支三格无载体。AC#5（P12 同步客户端须给 env/registry/config **级**证据）落真机 ⇒ 那一格量不到，整票成立 |
| 56 | `56-s8-signing-distribution-naming.md` | 144151ae 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD \| grep -icE 'winget\|scoop'`；`git ls-tree -r --name-only HEAD -- docs \| grep -icE 'docs/RELEASE'`；`git log -1 --format='%h %ad' --date=format:'%m-%d' -- docs/BUILD.md` | winget/Scoop 载体 **0 枚**；`docs/RELEASE.md` **不存在**（计数 0）；`docs/BUILD.md` 最后动过＝`10ed80d3 09-19`（票 01 冻结那一笔）此后未更新 ⇒ AC#2/#4/#5 静态即不成立。已落的一半：minisign 链在（`internal/models/minisign.go`＋`tools/signmodels/main.go`）；AC#1 SmartScreen 证据量不到 |
| 57 | `57-s8-plugin-sdk-registry.md` | 144151ae 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD \| grep -iE 'plugin-sdk\|sdk/' \| grep -v '\.scratch'`；叠在票 50 读数之上 | SDK／套件文件 **0 枚**；AC#1 判据是"第三方作者只靠文档＋套件出货 Tier1/Tier2 插件"，而两层引擎在 `unwired.go:141-146` 挂着 "not built"（票 50 尺）⇒ 套件没有可包的引擎。AC#2 registry install 落运行期 ⇒ 那一格量不到 |
| 58 | `58-s8-i18n-english.md` | 144151ae 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD \| grep -iE 'locale\|i18n\|en-US\|\.ftl' \| grep -v '\.scratch'` | 语言包／locale 框架文件 **0 枚** ⇒ AC#1"扫描零硬编码中文"连可扫的 locale 层都没有；AC#2 英文语音基线依赖的引擎面（批一票 15＝`internal/speech/` 只有 `doc.go`）也不在；AC#3 免重启切换＝运行期，量不到 |
| 59 | `59-p15-aec-spike.md` | 130c0943 09-19 | **仍成立** | `git grep -iln 'aec3\|erle' HEAD -- internal/ cmd/ third_party/`；`git ls-tree -r --name-only HEAD \| grep -iE 'aec\|erle\|echo'`；`git ls-tree -r --name-only HEAD -- scripts/spike`（只列目录名） | 内容尺在生产面 **0 命中**；文件名尺命中 3 枚＝本票票文件、一枚无关探针日志、`internal/agent/testdata/golden/repeat-echo.sse`（那枚 "echo" 是工具名 echo）；`scripts/spike/` 现有 spike 目录里**没有 aec 那一枚**（列名核对：common／goja-caps／model-residency／shell-baseline／speech-baseline／webview2-latency／xy-verdict）⇒ 许可裁决、绑定构建、回音质量台三格无产物。真机回音读数那一半量不到 |
| 60 | `60-c32-realtime-engine.md` | 5cba2d92 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD -- internal/llm \| grep -oE 'internal/llm/[a-z]+/' \| sort -u`；`git grep -n 'Realtime' HEAD -- internal/config/schema.go` | `internal/llm/` 适配器目录只有 `anthropic/`、`openaichat/`、`adaptertest/`、`golden/`（全文本 SSE 面）＋**零语音实时适配器**；`voice.realtime`（C32）在 `schema.go:232,258,329` **只有配置形状、无消费者**；引擎之家 `internal/speech/` 只有 `doc.go` ⇒ 多轮实时／barge-in／零工具证明／上下文卫生四格无处装配。AC#1 的 E2E 量不到 |
| 61 | `61-cloud-voice-providers-c9.md` | 5cba2d92 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD -- internal/ \| grep -iE 'speech\|voice' \| grep -v doc.go`；`git grep -n 'DEFERRED(engines)' HEAD -- internal/speech/doc.go` | 除 `doc.go` 外**没有任何 speech/voice 实现文件**；`doc.go:19` 那行 `DEFERRED(engines): implemented by ticket 15 … ticket 26 … ticket 41` 原样挂着，而同文件 `:3` 早写着 "plus the interface seams for cloud providers (C9, ticket 61)"⇒ 云端 ASR/TTS 接缝一行代码没有。链推进／配额定量／音频用量入账三格的**被测对象**不存在（配额消费侧见批一票 44＝`QuotaDailyMicro` 在 config 外 0 命中） |
| 65 | `65-ball-glass-quality-rework.md` | ed1e7c62 09-20 | **量不到** | — | 六格里五格判据＝`design/refs/*.png` 比对、并排差分产物、owner 第二次看 `-tour` 说"这次对了"⇒ 全落 `design/**`（禁读）＋真开窗＋人工判决。唯一半静态那格（渲染码无新增硬编码色值）要拿今树与返工后比，返工没发生 ⇒ 归口：视觉/tour 腿。⚠ 别读成"球不用返工了" |
| 86 | `86-resolvepercallbudget-wallclock-fragility.md` | 3d716a4b 09-21 | **仍成立** | `git grep -n 'resolveBudget\|testing.Short\|baseline' HEAD -- internal/risk/pathresolver_budget_norace_test.go`；`git log --format='%h %ad' --date=format:'%m-%d' -- <该文件>` | `:12 const resolveBudget = time.Millisecond` 原样在；`:36 if time.Duration(ns) > resolveBudget { t.Fatalf(...) }`＝**绝对墙钟断言未改**；`testing.Short`／自校准基线 **0 命中**；该文件最后动过＝`ec2c8799 09-20`（票 18 落地那笔），09-21 建票至今没人碰过它 ⇒ AC#1–#4 零进展。分布读数要跑 go ⇒ 那半格量不到，"没改"静态坐实 |
| 91 | `91-os-isolation-spike-restricted-token-vs-appcontainer.md` | 6114e3df 09-21 | **差翻勾** | 凭据尺：`git log --format='%h %ad %s' --date=format:'%m-%d' --diff-filter=A -- docs/evidence/s1/91-os-isolation-memo.md`；`git ls-tree -r --name-only HEAD -- docs/evidence/s1 \| grep '/91-'`；`grep -n '票 91' docs/reports/pending-and-issues.md \| tail -3`（⛔ 未读 evidence 内容，只按文件名与 commit 标题） | 四条 AC 的产物**两枚都在树上**：`docs/evidence/s1/91-isolation-options.md`（第一会话存盘点）＋`91-os-isolation-memo.md`（第二会话，`73cec78e 09-21`，commit 标题逐字 `docs(91,AC#1..AC#4): 第二会话补齐 OS 隔离实测`）；台账 `A70`（`pending-and-issues.md:2046`、`:2048`）逐字记"**票 91 的第二会话把三路都量了……现在的答案有读数了**"。⇒ 一枚"已裁成立、只差翻勾"的票，正是 pv1 要补那桶。**两处分叉本腿不裁**：①AC#1 点名的文件叫 `91-isolation-options.md`，产物另起一枚 `-memo.md`（票内注记与 A70① 都写着"别往 options 那枚上写、编排者最后合并"）；②A70① 自带档位"〔代理实测，我未逐格复现〕"⇒ 凭据种类＝票内注记＋台账 A##，**缺**"非实现者裁决表"那一层（本腿禁读 `docs/evidence/s1/**` 内容，只核到文件名与 commit 标题）。票还写着 DPAPI 那批读数"全不采信、AC 下能否解密仍未证" ⇒ 本腿只说"备忘录已交付"，不说"结论已复核" |
| 98 | `98-cmd-wisp-tests-never-run-on-this-host.md` | 92818af9 09-28 | **量不到** | 允许的尺只这一把：`git ls-tree -r --name-only HEAD -- cmd/wisp \| grep -c '_test.go'`（**数文件名，未读内容**） | **65 枚** `_test.go` 在 `cmd/wisp/` ⇒ AC#1(ii)"如果能跑会跑掉多少条"的分母材料在；但五格 AC 的判据全是运行期读数（本机 `go test ./cmd/wisp/` 原文与 `0xc0000135`、注入 DLL 路径后的真实 PASS/FAIL 数、守卫红没红）⇒ 本腿零 go 命令，**且 `cmd/wisp/**` 此刻是并行腿 `231-r1` 的地界**（它正在里面跑定向用例与整包读数）。归口写死：交 `231-r1`／后续能跑 go 的腿，⛔ 本腿不去量，免得互相洗 |

## §2 本批计数

| 档 | 枚数 | 票号 |
|---|---|---|
| 仍成立 | **11** | 50 51 52 55 56 57 58 59 60 61 86 |
| 已失效／已被别人做掉 | **0** | —（无一枚过得了硬门） |
| 差翻勾 | **1** | 91 |
| 量不到 | **4** | 53 54 65 98 |
| 合计 | **16** | ✓ 与 §0 分母对上 |

## §3 与 pv1 的对照 ＋ 本批两处具名读数

pv1 批二 16 行判定**全是"未判"占位**（`grep '^| [0-9]' .scratch/wisp/probes/pool-validity/1/batch2.md \| grep -v 未判` → 0 命中）⇒ **本批无对判对象**，"它判 X／我判 Y"这一栏为空。

两处值得单独上报：

1. **票 91 是四批里第一枚落进〔差翻勾〕的**，凭据＝票内注记＋台账 `A70`，⛔ 不是"缺陷字面今树 0 命中"。
   它同时把 pv1 §2 那个合并会造成的错演出来一遍：若按 pv1 分档法，这枚会被记成"已失效／已被别人做掉"，
   凭空多出一枚"被别的 commit 做掉的能力票"——而它本来就是只读评估票，活干完、框没翻。
   ⇒ 编排者那句"还剩多少真活"如果按 pv1 的三档读，会把 91 从"真活"里扣错两次（扣成失效、又算成别人做过）。
2. **票 98 是本腿唯一一枚"归口点名另一条活腿"的**：五格全要 go 读数，而 `cmd/wisp` 正被 `231-r1` 占。
   本腿按边界只数了文件名（65 枚），⛔ 未碰包内容、⛔ 未跑 go 尺。
