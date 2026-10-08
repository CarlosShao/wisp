# 票 35 `:52`「No secret leakage scan across bridge payloads」— 35-r7 非实现者验收判定（v7）

裁决腿：验收腿（非实现者，SPEC-12 §4.3 #1/#3）。起手 HEAD `a67efeac` → 取证中共享树自前进至 `e05b8a2b` → 收尾 HEAD `0c9726f9`。日期 2026-10-08。
被审件：`internal/panel/inbound_raw_leak_35r7_test.go`（引入提交 `6038531c`）。
内容锚＝票面那句英文判据本身（`- [ ] No secret leakage scan across bridge payloads.`，现量落点行号 `:52`）。
⛔ 本腿未翻票面任何框、未改 `internal/panel/**` 任何字节（产码与测试件都没改）、零 push。全部突变在仓外 `D:/tmp/wisp35v7/mut/` 经 `go test -overlay` 完成。

---

## Q1 落地性

**结论一句**：文件真在盘上、278 行与自报逐字同、由单笔只含该文件的提交带入且那笔提交只碰这一枚文件（零产码改动可证）；票面现量 6 枚未勾／4 枚已勾，`:52` 仍是未勾框——本腿没动它。

尺与现量（件：`logs/q1-anchors.txt`、`logs/q1-gitshow.txt`）：
- `ls -la internal/panel/inbound_raw_leak_35r7_test.go` → `12158 Oct 8 11:27`，`rc=0`；`wc -l` → **278**，`rc=0`（＝实现腿自报 278，同数）。
- `git log --oneline -8 -- <该文件>` → 唯一一笔 `6038531c`，`rc=0`；`git show --stat 6038531c` → `1 file changed, 278 insertions(+)`，文件名只有 `internal/panel/inbound_raw_leak_35r7_test.go`。`git show --name-only` 同（件 `logs/q3-host-wiring.txt`）⇒ **那一笔带的就是显式 pathspec 的形状，顺带改动的产码＝0 枚**。
- `git merge-base --is-ancestor 6038531c HEAD` → `rc=0`（已在 dev 历史上）。
- 票面框数（⛔ 带 `[[:space:]]*` 的那把尺）：`grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/35-panel-bridge-c17.md` → **6**，`rc=0`；`- \[x\]` → **4**，`rc=0`；未勾六枚的行号现量 `:42 / :44 / :45 / :47 / :49 / :52`。
- 禁区起手态：`git status --porcelain -- cmd/wisp internal/panel docs .github/workflows frontend .scratch/wisp/issues` → **空输出，rc=0**。

读错了最坏放行什么：若那笔提交里夹带一枚产码改动而我只看行数，就会把「测试加了牙是因为产码改了」记成「仪器本身有牙」。

---

## Q2 这枚仪器有没有牙（本腿自建 8 发突变，⛔ 未复用实现腿自报的 M1/M2/M3 表）

**结论一句**：**有牙，且牙在三个不同层面上各自被验出**——凭据形状尺（正控）、写通道两枚防静音钉（`credCalls==1` / `lastRawSeen==raw`）、防漂钉都能单独红；同时具名三处恒真/无鉴别力面（下面 ①②③）。

逐发（件：`logs/mut-*.txt` 八枚，每件尾部自落 `[rc=N]`；驱动脚本存 `logs/mutation-driver.txt`）：

| 发 | 错了什么（overlay 替换，仓内零写入） | rc | 红的用例 | 判 |
|---|---|---|---|---|
| M0 | 无突变，只跑 4 名 35r7 用例 | 0 | 无 | 基线绿，突变面成立的前提 |
| **M1** | **中和判据本身**：`hits35r7` 的命中条件 `if re.Find(b) != nil {` → `&& false` | 1 | `TestInboundLeakRulerFiresOnPlantedShapes35r7`（`:211` "the mirrored ruler did not fire"） | 正控有牙；**但同跑里五路由泄密主测 5 枚子例全 PASS** ⇒ 见恒真面③ |
| **M2** | 镜像尺漂 1 字符：`sk-[A-Za-z0-9]{12,}` → `{11,}`（⛔ 实现腿表里没有这一形） | 1 | `TestDecisionRulesMirrorTicket248Ruler35r7`（`:272` "the generic key-shape rule drifted"） | 防漂钉是真钉，不是装饰 |
| M3 | 摘掉 `if reply != ""` 守卫 → `if true` | **0** | 无 | **恒真面①**：该守卫装饰性（空 blob 永不命中），摘了零信号 |
| **M4a** | 产码：`h.Store.StoreCredential(ctx, raw)` → `(ctx, "")`（跳接受写入但把原文吞掉） | 1 | `credential-write-accepted`（`:136` "did not really carry the raw (credCalls=1)"） | **`lastRawSeen==raw` 是真牙** |
| **M4b** | 产码：`res, err := h.Store.StoreCredential(ctx, raw)` → `SettingWriteResult{}, error(nil)`（处理器什么都不做） | 1 | `credential-write-accepted` | **`credCalls==1` 是真牙** |
| **M5** | 产码：凭据 accept 路由把入向 raw 拼进**返回给宿主的 receipt**（⛔ 实现腿三面全打在 audit/error 上，没打过「接受路由的 receipt 面」这一形） | 1 | `credential-write-accepted`：`leak surface: receipt matches sk-[A-Za-z0-9]{12,}` + `matches <canary-redacted>` | 页面可见面确实被盯着；命中原文按 248 纪律全程不回显（日志里只有 `<canary-redacted>`，本腿逐字核过） |
| **M6** | 产码：`dispatch` 两枚拒绝分支 `return err` → `return nil`（"让 handler 什么都不做"，即派单要求的第二形） | 1 | **只有** `refused-unattached-handler-carrying-content`（`:199`） | 未接处理器那枚是真牙；**`rosterMismatch` 被中和时全套绿** ⇒ 恒真面② |
| M7 | 产码：普通写路由把 `value` 拼进 receipt（`sk-` 形状就住在 value 里） | 1 | `plain-write-with-key-shaped-value`：`leak surface: receipt matches sk-[A-Za-z0-9]{12,}` | 第二枚内容字段面也有牙，不是只测哨兵 |

**恒真面具名（哪一发怎么坏都照绿）**：
- **①** `checkArtifacts35r7` 里的 `if reply != ""` 守卫＝装饰（M3 rc=0 全套绿）。真正干活的是 `hits35r7` 本身，删掉守卫不改任何判定。
- **②** 派发表 `default` 分支（`rosterMismatch`）**不被本件覆盖**：M6 同时中和 `rosterMismatch` 与 `unattached`，只有后者红。原因是那枚路由的拒绝来自 `ParseComposerRequest`（`bridge.go`），raw 根本走不到 `dispatch`（`composer_dispatch.go:154-156` 现量：parse 先返，再 dispatch）。⇒ 子用例那句 `t.Fatal("an unlisted method was accepted")` **对派发表 backstop 零鉴别力**，它实际重测的是 parse 的门（票 35 里已被 `:64`/`:76`/`:78` 那几格裁过的「同一道门测两遍」形状）。实现腿自述「未列方法名拒绝」这一路由**没有证明派发表的 fail-closed**，只证明 parse 的。
- **③** 五路由泄密主测**对盲尺无感**：M1 把尺摘掉后主测 5 枚子例照样 PASS（日志逐名可查）。整套「有牙」押在单独一枚 `TestInboundLeakRulerFiresOnPlantedShapes35r7` 上＝**单点防线**。这与票 248 原尺同形（其正控在 `:218`），不算违令，但必须具名：将来谁改那枚正控而没改尺，CI 只会红一处，泄密主测会一路绿。

⛔ 本腿未放宽任何断言、未摘尺当修尺：所有突变跑完即弃（`D:/tmp` 里的 overlay 副本），仓内被 overlay 的三枚文件跑后 `git hash-object` 与 `HEAD:<path>` **逐枚同**（`config_handlers.go` 19445546…／`composer_dispatch.go` 1757a487…／`inbound_raw_leak_35r7_test.go` d6b22579…，见 `logs/q4-diff-and-restore.txt`），`git status --porcelain -- internal/panel cmd/wisp` 收尾仍空。

---

## Q3 「No secret leakage across bridge payloads」这一格闭合了没有

**结论一句**：**未闭合**。本件把「入向那一跳产出的三面」这半句做到有牙（M1/M2/M4a/M4b/M5/M7 五形独立红可指），但票面判据是「across bridge payloads」＝**跨双向载荷**：①出向推载荷今天无被测物（本腿复跑实现腿那把尺＝同色，零命中）；②写腿是空桩，"raw 交给腿" 不等于 "腿存成凭据而非明文"；③派发表 backstop 面是恒真面②。⇒ `:52` 框**不该翻**。

**它防的结局形状**：用户从面板粘进去的一枚 key（或任意 `sk-` 形状串），被这一跳抄进本机能事后读到的面上——stderr/日志里的 `[audit]` 行、回到页面那句 receipt、返回给宿主的 error——于是密钥在 DPAPI 之外多长出若干份明文副本（日志文件、页面 DOM、事后被人 grep 的终端记录）。

**它看得见什么（三面里有两面是真物，本腿现量，件 `logs/q3-production-wiring.txt`／`logs/q3-host-wiring.txt`）**：
- **receipt 面＝真页面可见**：`cmd/wisp/panel_host_windows.go:814-822` 的 `dispatchRaw` 直接 `return m.disp.Handle(ctx, raw)`，宿主 bind 闭包把这返回值原样交回页面；`cmd/wisp/panel_inbound.go:163` 也把 `reply` 逐行打到 stdout（`:169` 那句 `"第 %d 行被拒绝：%s"`）。⇒ M5/M7 的红是真泄漏形状，不是纸面形状。
- **audit 面＝真落盘面（stderr）**：`cmd/wisp/panel_inbound.go:142-144` `auditf := func(format string, args ...any){ fmt.Fprintln(s.stderr, "[audit] "+fmt.Sprintf(format, args...)) }`，并在 `:268/:280` 逐字接到 `ConfigWriteHandler.Audit` 与 `ComposerDispatch.Audit`。测试闭包用的是**同一个 format+args 的 Sprintf**，捕获＝生产会写出的那一行。⇒ 这里**不是**把守卫绕过去。
- **白盒直调绕过了什么**：本件是 `Handle` 直调，⛔ 没过 WebView2 收信回调（`panelPostMessageForwardInit` 那侧）、没过回写页面那一跳（H10）、没过真 ConfigStore。parse 那道门（伪造来源/未列名）**确实在 Handle 内被真走过**——这点由 M6 定位清楚。

**最坏会被它放行什么真 bug（具体到谁把什么写到哪里、盘上剩什么）**：
1. **凭据值被写腿落成明文**：`cfg248Leg.StoreCredential` 返回罐头 `CredentialSummary{Ref: "dpapi:35r7refnotablob"}`，本件只断言 `lastRawSeen==raw`（raw 交出去了）。若 `cmd/wisp/panel_inbound.go:267` 的 `newConfigStore(mgr, secrets, …)` 那条真腿把 key 写进 `%APPDATA%\wisp\config.toml` 或 `wisp.db` 的明文列，**本件照绿**。这就是本项目反复出现的那形（空桩把出向吞成 nil ⇒ 包绿证明不了对端收到）的变体：**"包绿" 证明不了 "对端存对了"**。面在 `cmd/wisp`/`internal/config`，本程写面禁 ⇒ 属具名残差，不是本件缺陷。
2. **出向推载荷泄密**：Go→页主动推的那一载荷（事件/快照/resync）不在本件被检面里，且今天**没有被测物**——本腿现量 `grep -rn "PostWebMessage\|EvaluateScript\|\.Eval(" --include=*.go cmd internal | grep -v _test.go` → **空**（件 `logs/q3-carrier.txt`，那行 `rc` 由 grep 自带）。⇒ 「across bridge payloads」的另半边是**欠物**不是欠尺，实现腿的残差②判定本腿复核成立。
3. **派发表 backstop 侧的抄写**：恒真面② ⇒ 若将来 `bridge.go` 名册加一枚方法而 `dispatch` 忘了 case，那一路由的拒绝句改由 `rosterMismatch` 产出，本件今天对它零鉴别力，抄了 raw 也不红（audit 抄 raw 那一形要等实现腿 M1 那种突变在**新分支**上重跑才看得见）。

---

## Q4 不新增红 + 卫生

**结论一句**：整包前后**逐名作差新增红＝0**（同名 4 枚既有红，全是前端/token 在飞面）；`go vet` rc=0 且**零输出（0 字节，已具名自证）**；`sh scripts/d22scan.sh` rc=0 clean；三枚被 overlay 的文件 hash 与 HEAD 逐枚同，仓内零写入。

尺与现量（件：`logs/q4-gates.txt`、`logs/q4-test-after.txt`、`logs/q4-baseline-full-package.txt`、`logs/q4-vet.txt`、`logs/q4-d22scan.txt`、`logs/q4-diff-and-restore.txt`；每条 rc 均由命令自身落一行，`echo $?` 前⛔ 无管道）：
- `go test ./internal/panel/ -count=1` 起手（写入被审件之前那一次）→ `rc=1`，红 4 枚：`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestC21DesignTokensFourWayAgree`。
- 同命令收尾再跑 → `rc=1`，红**同名 4 枚**；`diff` 两份 FAIL 名册 → **rc=0（空 diff）＝新增红 0 枚**。⚠ 这 4 枚不归本程账（`frontend/**`/token 族在飞），本腿一枚没放宽。
- `go vet ./internal/panel/` → **rc=0**（11:47 首发 0 字节输出；11:53 收尾复跑同 rc=0，件 `logs/q4-vet.txt`／名册 `logs/q4-final-refix.txt`）。⚠ **具名一次瞬时 rc=1**：11:52 那一发 vet 报 rc=1，而本腿当时把输出重定向进了 `/dev/null`（自己踩了「门禁件必须自带输出」那枚坑）**⇒ 无红句可查、不能归因**；11:53 同命令复跑 rc=0，同刻禁区 `git status --porcelain` 空。共享工作树里最可能是别的腿瞬时脏面，本腿不把它记进 35-r7 账，也不掩盖，全程留件 `logs/q4-vet-attest.txt`。
- `sh scripts/d22scan.sh` → **rc=0**，尾行 `d22scan: clean - no D22 ban violations`（⛔ 未用仓根 `go run ./tools/d22scan`，那是独立 module，命令必失败）。
- `gofmt` 面：本腿未新增任何仓内 Go 文件，突变件全在 `D:/tmp` ⇒ 不适用，具名不装。
- 收尾 `git status --porcelain -- cmd/wisp internal/panel docs .github/workflows frontend .scratch/wisp/issues` → **空，rc=0**；全仓 porcelain 行数 756（起手）→ 771（收尾，增量＝本腿 `probes/35/v7/**` 自己的证据件）。

---

## 我没量到的／不成立的派单前提

**顶回派单（逐条具名，均为本腿现量）**：
1. **⚠ 派单前提不成立**：派单 §3 要求我具名说明「`cmd/wisp/panel_host_windows.go` 那枚 M 是别人的在飞改动」。实测**今天没有这枚 M**：起手 `git status --porcelain` 全仓 756 行里 `grep -c "cmd/wisp"` → **0（grep rc=1）**，定向 `git status --porcelain -- cmd/wisp internal/panel docs .github/workflows frontend .scratch/wisp/issues` → 空、rc=0（件 `logs/q1-anchors.txt`）。收尾同空。⇒ 我没有提交也没有还原任何它，因为它不在脏面上；禁区写权本腿全程未碰。若编排者的"在飞"账记的是另一枚锚点时刻的状态，请以本件现量为准更正那条账。
2. **行号锚与文件内注释不一致（不冲突，但会漂）**：被审件第 3 行自述 `Ticket 35 AC#6 (票面 :51)`，而票面那句英文判据现量在 `:52`（`:51` 是普查表那一行）。内容锚判＝一致，本件确实是给 `:52` 那格造的仪器；但注释里的行号引用已经是错的，下次插行会更错。
3. **实现腿自述的「五路由」里有一枚不裁它声称的东西**：「未列方法名拒绝」这一路由的拒绝实际发生在 `ParseComposerRequest`，`dispatch` 的 `default`/`rosterMismatch` 分支从未被执行（M6 现量：中和 `rosterMismatch` 全套绿）。所以「五路由扫三面」这句在第四枚上应读作「parse 门＋它的 audit 句」，不是「派发表 fail-closed 被扫过」。
4. **`composer_dispatch.go` 头注过期**：`:48` 写「it has NO production caller yet」，现量 `cmd/wisp/panel_inbound.go:163` 与 `cmd/wisp/panel_host_windows.go:821` 两处非 test 调用。非 35-r7 责任（本程禁改产码），具名上报，别按那行注释推断"面还没接上"。

**本腿没量到的（不装成量过）**：
- 真 WebView2 收信/回写那一跳（需 `winlive` 构建面）——本腿无平台面，未取证。
- `wisp panel-inbound` 之外的宿主装配里 `auditf` 是否被重定向进日志文件（`cmd/wisp/run.go:679` 那侧本腿只看到接线名，未跟到底）。⇒ 「audit＝落盘面」在 CLI 装配里现量是 **stderr 面**，"落盘"这半句本腿只对到 stderr→日志文件那一环存在与否未证。
- 票 248 原尺（`cmd/wisp/panel_config_248_test.go`）的六面本身：本腿只靠 M2 那枚防漂钉间接对拉（canary 字面＋两枚 MustCompile 计数），未逐行重裁 248 族。
- 恒真面②给出的后果（新方法进名册后 `rosterMismatch` 抄 raw）今天不可达，本腿**没有**构造 bridge+dispatch 双改的复合突变去证明它将来可达。
